package authuser

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type Config struct {
	BaseURL string
	Client  *http.Client
}

type Client struct {
	baseURL    *url.URL
	httpClient *http.Client
}

func NewClient(config Config) (*Client, error) {
	if strings.TrimSpace(config.BaseURL) == "" {
		return nil, errors.New("BaseURL is required")
	}
	baseURL, err := url.Parse(strings.TrimRight(config.BaseURL, "/"))
	if err != nil {
		return nil, err
	}
	if baseURL.Scheme == "" || baseURL.Host == "" {
		return nil, errors.New("BaseURL must be absolute")
	}
	if config.Client == nil {
		config.Client = &http.Client{Timeout: 5 * time.Second}
	}
	return &Client{baseURL: baseURL, httpClient: config.Client}, nil
}

type UserSnapshot struct {
	ID            string   `json:"id"`
	Email         string   `json:"email"`
	Name          string   `json:"name"`
	Status        string   `json:"status"`
	EmailVerified bool     `json:"emailVerified"`
	Licenses      []string `json:"licenses"`
	Roles         []string `json:"roles"`
	Permissions   []string `json:"permissions"`
	CreatedAt     string   `json:"createdAt"`
	UpdatedAt     string   `json:"updatedAt"`
}

type UserCard struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func (c *Client) Me(ctx context.Context, bearerToken string) (UserSnapshot, error) {
	var snapshot UserSnapshot
	if err := c.do(ctx, http.MethodGet, "/api/userinfo/me", bearerToken, nil, &snapshot); err != nil {
		return UserSnapshot{}, err
	}
	return snapshot, nil
}

func (c *Client) Lookup(ctx context.Context, bearerToken string, userIDs []string) ([]UserCard, error) {
	var cards []UserCard
	if err := c.do(ctx, http.MethodPost, "/api/userinfo/lookup", bearerToken, lookupRequest{UserIDs: userIDs}, &cards); err != nil {
		return nil, err
	}
	return cards, nil
}

type lookupRequest struct {
	UserIDs []string `json:"userIds"`
}

func (c *Client) do(ctx context.Context, method string, path string, bearerToken string, body any, target any) error {
	if strings.TrimSpace(bearerToken) == "" {
		return errors.New("bearer token is required")
	}

	var reader *bytes.Reader
	if body == nil {
		reader = bytes.NewReader(nil)
	} else {
		content, err := json.Marshal(body)
		if err != nil {
			return err
		}
		reader = bytes.NewReader(content)
	}

	endpoint := c.baseURL.ResolveReference(&url.URL{Path: path})
	request, err := http.NewRequestWithContext(ctx, method, endpoint.String(), reader)
	if err != nil {
		return err
	}
	request.Header.Set("Authorization", "Bearer "+bearerToken)
	if body != nil {
		request.Header.Set("Content-Type", "application/json")
	}

	response, err := c.httpClient.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode < 200 || response.StatusCode >= 300 {
		var errorResponse struct {
			Error string `json:"error"`
		}
		_ = json.NewDecoder(response.Body).Decode(&errorResponse)
		if errorResponse.Error == "" {
			errorResponse.Error = response.Status
		}
		return fmt.Errorf("auth user %s %s: %s", method, path, errorResponse.Error)
	}

	return json.NewDecoder(response.Body).Decode(target)
}
