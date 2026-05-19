package rbac

type Role struct {
	ID         string  `json:"id"`
	Key        string  `json:"key"`
	Name       string  `json:"name"`
	ProductKey *string `json:"product_key,omitempty"`
}

type Permission struct {
	ID          string  `json:"id"`
	Key         string  `json:"key"`
	Description string  `json:"description"`
	ProductKey  *string `json:"product_key,omitempty"`
}

type RolePermission struct {
	RoleID       string `json:"role_id"`
	PermissionID string `json:"permission_id"`
}
