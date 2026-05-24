package requestmeta

import (
	"net"
	"net/http"
	"strings"
)

func ClientIPAddress(r *http.Request) string {
	for _, header := range []string{"X-Forwarded-For", "X-Real-IP"} {
		value := r.Header.Get(header)
		if value == "" {
			continue
		}
		ip := strings.TrimSpace(strings.Split(value, ",")[0])
		if net.ParseIP(ip) != nil {
			return ip
		}
	}

	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err == nil && host != "" {
		return host
	}

	return r.RemoteAddr
}

func CountryFromRequest(r *http.Request) string {
	for _, header := range []string{"CF-IPCountry", "X-CodeLinks-Country"} {
		value := strings.ToUpper(strings.TrimSpace(r.Header.Get(header)))
		if len(value) == 2 {
			return value
		}
	}

	return "ZZ"
}

func CityFromRequest(r *http.Request) string {
	return strings.TrimSpace(r.Header.Get("X-CodeLinks-City"))
}

func DeviceFromUserAgent(userAgent string) (string, string) {
	lower := strings.ToLower(userAgent)
	browser := "Unknown"
	switch {
	case strings.Contains(lower, "edg/"):
		browser = "Edge"
	case strings.Contains(lower, "chrome/"):
		browser = "Chrome"
	case strings.Contains(lower, "firefox/"):
		browser = "Firefox"
	case strings.Contains(lower, "safari/"):
		browser = "Safari"
	}

	operatingSystem := "Unknown"
	switch {
	case strings.Contains(lower, "windows"):
		operatingSystem = "Windows"
	case strings.Contains(lower, "mac os") || strings.Contains(lower, "macintosh"):
		operatingSystem = "macOS"
	case strings.Contains(lower, "linux"):
		operatingSystem = "Linux"
	case strings.Contains(lower, "android"):
		operatingSystem = "Android"
	case strings.Contains(lower, "iphone") || strings.Contains(lower, "ipad"):
		operatingSystem = "iOS"
	}

	return browser, operatingSystem
}
