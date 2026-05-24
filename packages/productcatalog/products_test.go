package productcatalog

import (
	"encoding/json"
	"os"
	"reflect"
	"sort"
	"strings"
	"testing"
)

type productConfig struct {
	ID string `json:"id"`
}

func TestProductCatalogMatchesTypeScriptConfig(t *testing.T) {
	content, err := os.ReadFile("../config/products.json")
	if err != nil {
		t.Fatal(err)
	}

	var products []productConfig
	if err := json.Unmarshal(content, &products); err != nil {
		t.Fatal(err)
	}

	want := make([]string, 0, len(products))
	for _, product := range products {
		want = append(want, product.ID)
	}
	got := append([]string{}, IDs...)
	sort.Strings(want)
	sort.Strings(got)

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("product catalog drifted\nwant: %#v\ngot:  %#v", want, got)
	}
}

func TestAuthSchemaProductConstraintMatchesCatalog(t *testing.T) {
	content, err := os.ReadFile("../../apps/auth-service/backend/internal/store/postgres/migrations/001_auth_schema.sql")
	if err != nil {
		t.Fatal(err)
	}
	schema := string(content)

	for _, id := range IDs {
		if !strings.Contains(schema, "'"+id+"'") {
			t.Fatalf("auth schema product constraint missing %q", id)
		}
	}
}
