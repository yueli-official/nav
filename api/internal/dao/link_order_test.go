package dao

import "testing"

func TestLinkOrderUsesWhitelistedColumnsAndDirection(t *testing.T) {
	tests := []struct {
		name      string
		sort      string
		direction string
		want      string
	}{
		{name: "default ascending", sort: "default", direction: "asc", want: "featured DESC, sort_order ASC, title ASC"},
		{name: "default descending", sort: "default", direction: "desc", want: "featured DESC, sort_order DESC, title DESC"},
		{name: "recently updated", sort: "updated", direction: "desc", want: "updated_at DESC, title ASC"},
		{name: "title", sort: "title", direction: "asc", want: "title ASC"},
		{name: "published date", sort: "published", direction: "desc", want: "published_at DESC NULLS LAST, title ASC"},
		{name: "unknown values fall back", sort: "drop table", direction: "sideways", want: "featured DESC, sort_order ASC, title ASC"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := linkOrder(test.sort, test.direction); got != test.want {
				t.Fatalf("linkOrder(%q, %q) = %q, want %q", test.sort, test.direction, got, test.want)
			}
		})
	}
}
