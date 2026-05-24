package admin

import "strings"

func maskIP(value string) string {
	parts := strings.Split(value, ".")
	if len(parts) == 4 {
		return parts[0] + "." + parts[1] + ".x.x"
	}
	if value == "" {
		return value
	}

	return "masked"
}
