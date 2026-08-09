package iconcontract

import "testing"

func TestCategoryTablerContractIsFiniteAndDefensive(t *testing.T) {
	icons := CategoryTablerIcons()
	if len(icons) < 8 || !IsCategoryTablerIcon("i-tabler-palette") || IsCategoryTablerIcon("i-tabler-brand-github") {
		t.Fatalf("unexpected category icon contract: %#v", icons)
	}
	icons[0] = "i-tabler-mutated"
	if IsCategoryTablerIcon("i-tabler-mutated") {
		t.Fatal("callers must not be able to mutate the embedded contract")
	}
}
