package productcatalog

import (
	"encoding/json"
	"os"
	"reflect"
	"sort"
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
