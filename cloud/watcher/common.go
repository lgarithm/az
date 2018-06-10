package watcher

import "strings"

// ID2Name extracts the resource name from resource ID.
func ID2Name(id string) string {
	parts := strings.Split(id, "/")
	if len(parts) > 8 {
		return parts[8]
	}
	return ""
}

// ID2Group extracts the resource group name from resource ID.
func ID2Group(id string) string {
	parts := strings.Split(id, "/")
	if len(parts) > 4 {
		return parts[4]
	}
	return ""
}
