package main

import (
	"encoding/json"
	"os"
	"reflect"
	"sort"
	"strings"
	"testing"
)

type adminAccessContract struct {
	Roles           []AdminRole                     `json:"roles"`
	AdminRoles      []AdminRole                     `json:"adminRoles"`
	Permissions     []AdminPermission               `json:"permissions"`
	RolePermissions map[AdminRole][]AdminPermission `json:"rolePermissions"`
}

func TestAdminAccessContractMatchesGoPolicy(t *testing.T) {
	contract := loadAdminAccessContract(t)

	for _, role := range contract.Roles {
		if !validAdminRole(role) {
			t.Fatalf("contract role %q is not accepted by Go policy", role)
		}
	}

	for _, role := range contract.AdminRoles {
		got := permissionsByRole(role)
		want := append([]AdminPermission{}, contract.RolePermissions[role]...)
		sort.Slice(got, func(i, j int) bool { return got[i] < got[j] })
		sort.Slice(want, func(i, j int) bool { return want[i] < want[j] })
		if !reflect.DeepEqual(got, want) {
			t.Fatalf("permissions for role %q drifted\nwant: %#v\ngot:  %#v", role, want, got)
		}
	}
}

func TestAdminAccessContractMatchesSQLSeed(t *testing.T) {
	contract := loadAdminAccessContract(t)
	content, err := os.ReadFile("migrations/004_auth_signup_smtp_notifications.sql")
	if err != nil {
		t.Fatal(err)
	}
	sql := string(content)

	for _, permission := range contract.Permissions {
		if !strings.Contains(sql, "'"+string(permission)+"'") {
			t.Fatalf("SQL seed is missing permission %q", permission)
		}
	}
	for role, permissions := range contract.RolePermissions {
		for _, permission := range permissions {
			tuple := "('" + string(role) + "', '" + string(permission) + "')"
			if !strings.Contains(sql, tuple) {
				t.Fatalf("SQL seed is missing role-permission tuple %s", tuple)
			}
		}
	}
}

func loadAdminAccessContract(t *testing.T) adminAccessContract {
	t.Helper()

	content, err := os.ReadFile("../../../packages/config/admin-access.json")
	if err != nil {
		t.Fatal(err)
	}

	var contract adminAccessContract
	if err := json.Unmarshal(content, &contract); err != nil {
		t.Fatal(err)
	}

	return contract
}
