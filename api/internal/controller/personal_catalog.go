package controller

// PersonalAPIDeclaration exports the PAT permissions and route allowlist used
// by the running product.
func PersonalAPIDeclaration() map[string]any {
	permissions := make([]map[string]any, 0, len(personalPermissions))
	for _, permission := range personalPermissions {
		permissions = append(permissions, map[string]any{
			"key":         permission.Key,
			"label":       permission.Label,
			"description": permission.Description,
			"capability":  string(permission.Capability),
		})
	}
	routes := make([]map[string]any, 0, len(personalRoutes))
	for _, route := range personalRoutes {
		authentication := "pat"
		scopeAnyOf := [][]string{{string(route.capability)}}
		if route.capability == "" {
			authentication = "public-or-pat"
			scopeAnyOf = [][]string{{}}
		}
		routes = append(routes, map[string]any{
			"method": route.method, "path": route.path,
			"scopeAnyOf": scopeAnyOf, "authentication": authentication,
		})
	}
	return map[string]any{
		"schemaVersion": 1,
		"product": map[string]string{
			"id": "nav", "name": "Nav",
			"guide": "../developer-tokens.md", "openapi": "../../contracts/openapi/nav.json",
		},
		"permissions": permissions, "routes": routes, "mediaProfiles": []map[string]any{},
	}
}
