package adminaccess

import (
	"encoding/json"
	"os"
	"reflect"
	"testing"
)

type contract struct {
	Roles           []string            `json:"roles"`
	AdminRoles      []string            `json:"adminRoles"`
	Permissions     []string            `json:"permissions"`
	RolePermissions map[string][]string `json:"rolePermissions"`
}

func TestGeneratedAdminAccessMatchesContract(t *testing.T) {
	content, err := os.ReadFile("../config/admin-access.json")
	if err != nil {
		t.Fatal(err)
	}

	var want contract
	if err := json.Unmarshal(content, &want); err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(Roles, want.Roles) {
		t.Fatalf("roles mismatch: got %v want %v", Roles, want.Roles)
	}
	if !reflect.DeepEqual(AdminRoles, want.AdminRoles) {
		t.Fatalf("admin roles mismatch: got %v want %v", AdminRoles, want.AdminRoles)
	}
	if !reflect.DeepEqual(Permissions, want.Permissions) {
		t.Fatalf("permissions mismatch: got %v want %v", Permissions, want.Permissions)
	}
	if !reflect.DeepEqual(RolePermissions, want.RolePermissions) {
		t.Fatalf("role permissions mismatch: got %v want %v", RolePermissions, want.RolePermissions)
	}
}
