package requestmeta

import (
	"net"
	"net/http"
	"net/netip"
	"strings"
)

type Resolver struct {
	trustedProxyCIDRs []netip.Prefix
}

func NewResolver(trustedProxyCIDRs []string) (Resolver, error) {
	resolver := Resolver{trustedProxyCIDRs: make([]netip.Prefix, 0, len(trustedProxyCIDRs))}
	for _, rawCIDR := range trustedProxyCIDRs {
		rawCIDR = strings.TrimSpace(rawCIDR)
		if rawCIDR == "" {
			continue
		}
		prefix, err := netip.ParsePrefix(rawCIDR)
		if err != nil {
			return Resolver{}, err
		}
		resolver.trustedProxyCIDRs = append(resolver.trustedProxyCIDRs, prefix)
	}
	return resolver, nil
}

func (r Resolver) ClientIPAddress(request *http.Request) string {
	remoteIP := remoteIPAddress(request.RemoteAddr)
	if remoteIP != "" && r.isTrustedProxy(remoteIP) {
		if forwarded := forwardedIPAddress(request); forwarded != "" {
			return forwarded
		}
	}

	if remoteIP != "" {
		return remoteIP
	}

	return request.RemoteAddr
}

func (r Resolver) isTrustedProxy(ipAddress string) bool {
	ip, err := netip.ParseAddr(ipAddress)
	if err != nil {
		return false
	}
	for _, prefix := range r.trustedProxyCIDRs {
		if prefix.Contains(ip) {
			return true
		}
	}
	return false
}

func ClientIPAddress(r *http.Request) string {
	return Resolver{}.ClientIPAddress(r)
}

func forwardedIPAddress(r *http.Request) string {
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

	return ""
}

func remoteIPAddress(remoteAddr string) string {
	host, _, err := net.SplitHostPort(remoteAddr)
	if err == nil && host != "" {
		return host
	}
	if net.ParseIP(remoteAddr) != nil {
		return remoteAddr
	}

	return ""
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
