package notify

import (
	"strconv"
	"strings"
)

func parseReceiverTokens(raw string) (map[string][]string, []string) {
	typed := map[string][]string{}
	legacy := []string{}
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return typed, legacy
	}
	for _, item := range strings.Split(raw, ",") {
		item = strings.TrimSpace(item)
		if item == "" {
			continue
		}
		if strings.Contains(item, ":") {
			parts := strings.SplitN(item, ":", 2)
			if len(parts) != 2 {
				continue
			}
			key := strings.TrimSpace(parts[0])
			val := strings.TrimSpace(parts[1])
			if key == "" || val == "" {
				continue
			}
			typed[key] = append(typed[key], val)
			continue
		}
		legacy = append(legacy, item)
	}
	return typed, legacy
}

func containsWildcard(values []string) bool {
	for _, v := range values {
		v = strings.TrimSpace(v)
		if v == "*" || v == "all" {
			return true
		}
	}
	return false
}

func toIntSet(values []string) map[int]struct{} {
	set := make(map[int]struct{}, len(values))
	for _, v := range values {
		v = strings.TrimSpace(v)
		if v == "" {
			continue
		}
		if n, err := strconv.Atoi(v); err == nil {
			set[n] = struct{}{}
		}
	}
	return set
}
