package facility

import (
	"strings"

	"github.com/besart951/go_infra_link/backend/internal/handler/facility/alarm"
	"github.com/besart951/go_infra_link/backend/internal/handler/facility/fielddevice"
	"github.com/besart951/go_infra_link/backend/internal/handler/facility/hierarchy"
	"github.com/besart951/go_infra_link/backend/internal/handler/facility/internal/routing"
	"github.com/besart951/go_infra_link/backend/internal/handler/facility/objectdata"
	"github.com/besart951/go_infra_link/backend/internal/handler/facility/reference"
	"github.com/besart951/go_infra_link/backend/internal/handler/middleware"
	"github.com/gin-gonic/gin"
)

type routeDefinition = routing.Definition

func RegisterRoutes(protectedV1 *gin.RouterGroup, handlers *Handlers, authChecker middleware.AuthorizationChecker) {
	facility := protectedV1.Group("/facility")
	registerRoutes(facility, authChecker, routeDefinitions(handlers))
}

func registerRoutes(group *gin.RouterGroup, authChecker middleware.AuthorizationChecker, routes []routeDefinition) {
	for _, route := range routes {
		handlers := []gin.HandlerFunc{
			middleware.RequirePermission(authChecker, route.Permission),
		}
		if route.Entitlement != "" {
			handlers = append([]gin.HandlerFunc{middleware.RequireEntitlement(route.Entitlement)}, handlers...)
		}
		handlers = append(handlers, route.Handler)

		group.Handle(
			route.Method,
			route.Path,
			handlers...,
		)
	}
}

func routeDefinitions(handlers *Handlers) []routeDefinition {
	routes := make([]routeDefinition, 0, 96)
	routes = append(routes, hierarchy.Routes(hierarchyHandlers(handlers))...)
	routes = append(routes, reference.Routes(referenceHandlers(handlers))...)
	routes = append(routes, fielddevice.Routes(fieldDeviceHandlers(handlers))...)
	routes = append(routes, objectdata.Routes(objectDataHandlers(handlers))...)
	routes = append(routes, alarm.Routes(alarmHandlers(handlers))...)
	return withInferredEntitlements(routes)
}

func withInferredEntitlements(routes []routeDefinition) []routeDefinition {
	for index := range routes {
		routes[index].Entitlement = inferredEntitlement(routes[index].Path)
	}
	return routes
}

func inferredEntitlement(path string) string {
	lowerPath := strings.ToLower(path)
	switch {
	case strings.Contains(lowerPath, "bacnet"), strings.Contains(lowerPath, "object-data"), strings.Contains(lowerPath, "alarm"):
		return "infra.module_bacnet"
	case strings.Contains(lowerPath, "sps"):
		return "infra.module_sps"
	case strings.Contains(lowerPath, "field-device"):
		return "infra.module_field_devices"
	default:
		return ""
	}
}
