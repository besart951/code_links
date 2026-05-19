package middleware

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

const productKeyInfraLink = "infra_link"

type authorizeRequest struct {
	UserID     string `json:"user_id"`
	TenantID   string `json:"tenant_id"`
	ProductKey string `json:"product_key"`
	FeatureKey string `json:"feature_key"`
}

type authorizeResponse struct {
	Allowed bool    `json:"allowed"`
	Reason  *string `json:"reason,omitempty"`
}

func RequireEntitlement(featureKey string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if !entitlementsRequired() {
			c.Next()
			return
		}

		userID, ok := GetUserID(c)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			c.Abort()
			return
		}

		tenantID := strings.TrimSpace(c.GetHeader("X-Tenant-ID"))
		if tenantID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "tenant_id_required"})
			c.Abort()
			return
		}

		allowed, reason, err := authorizeEntitlement(c.Request.Context(), authorizeRequest{
			UserID:     userID.String(),
			TenantID:   tenantID,
			ProductKey: productKeyInfraLink,
			FeatureKey: featureKey,
		})
		if err != nil {
			c.JSON(http.StatusBadGateway, gin.H{"error": "platform_unavailable"})
			c.Abort()
			return
		}
		if !allowed {
			if reason == "" {
				reason = "entitlement_required"
			}
			c.JSON(http.StatusForbidden, gin.H{"error": reason})
			c.Abort()
			return
		}

		c.Next()
	}
}

func entitlementsRequired() bool {
	return strings.EqualFold(os.Getenv("PLATFORM_ENTITLEMENTS_REQUIRED"), "true")
}

func authorizeEntitlement(ctx context.Context, payload authorizeRequest) (bool, string, error) {
	platformURL := strings.TrimRight(os.Getenv("PLATFORM_URL"), "/")
	if platformURL == "" {
		platformURL = "http://platform-api:8080"
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return false, "", err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, platformURL+"/internal/authorize", bytes.NewReader(body))
	if err != nil {
		return false, "", err
	}
	req.Header.Set("content-type", "application/json")
	req.Header.Set("accept", "application/json")
	if token := os.Getenv("PLATFORM_INTERNAL_TOKEN"); token != "" {
		req.Header.Set("X-Internal-Token", token)
	}

	client := http.Client{Timeout: 5 * time.Second}
	response, err := client.Do(req)
	if err != nil {
		return false, "", err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return false, "", nil
	}

	var result authorizeResponse
	if err := json.NewDecoder(response.Body).Decode(&result); err != nil {
		return false, "", err
	}
	if result.Reason == nil {
		return result.Allowed, "", nil
	}
	return result.Allowed, *result.Reason, nil
}
