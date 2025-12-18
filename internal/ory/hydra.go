package ory

import "strings"

func HydraAdminURL(base, path string) string {
	trimmed := strings.TrimSuffix(strings.TrimSpace(base), "/")
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}
	if strings.HasSuffix(trimmed, "/admin") {
		return trimmed + path
	}
	return trimmed + "/admin" + path
}
