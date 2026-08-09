package iconcontract

import (
	_ "embed"
	"encoding/json"
	"regexp"
	"slices"
)

//go:embed category-tabler.v1.json
var categoryTablerJSON []byte

var tablerName = regexp.MustCompile(`^i-tabler-[a-z0-9]+(?:-[a-z0-9]+)*$`)

type categoryTablerContract struct {
	Version int      `json:"version"`
	Icons   []string `json:"icons"`
}

var categoryTablerIcons = loadCategoryTablerIcons()

func loadCategoryTablerIcons() []string {
	var contract categoryTablerContract
	if err := json.Unmarshal(categoryTablerJSON, &contract); err != nil {
		panic("decode embedded category Tabler icon contract: " + err.Error())
	}
	if contract.Version != 1 || len(contract.Icons) == 0 {
		panic("category Tabler icon contract must be non-empty version 1")
	}
	icons := append([]string(nil), contract.Icons...)
	slices.Sort(icons)
	icons = slices.Compact(icons)
	for _, icon := range icons {
		if !tablerName.MatchString(icon) {
			panic("invalid category Tabler icon contract value: " + icon)
		}
	}
	return icons
}

func CategoryTablerIcons() []string {
	return append([]string(nil), categoryTablerIcons...)
}

func IsCategoryTablerIcon(value string) bool {
	_, found := slices.BinarySearch(categoryTablerIcons, value)
	return found
}
