package navaudit_test

import (
	"testing"

	"github.com/yueli-official/foundation/go/audit"

	"platform/products/nav/api/internal/navaudit"
)

func TestDefinitionCompilesStableConsumerActions(t *testing.T) {
	catalog, err := audit.Compile(navaudit.Definition())
	if err != nil {
		t.Fatal(err)
	}
	if catalog.Digest() == "" {
		t.Fatal("compiled definition has no digest")
	}
}
