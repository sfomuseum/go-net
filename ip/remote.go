package ip

// Disclosure: The following was "vibe-coded" using gpt-oss-20b (20260303)
// This code has been reviewed and doesn't raise any obvious red flags.

import (
	"net"
	"net/http"
	"strings"
)

// RemoteIP extracts the client’s IP address from an http.Request,
// checking the most common proxy / service headers before falling back
// to r.RemoteAddr.  It returns an empty string if no valid IP can
// be determined.
//
//   r.RemoteAddr   – usually "ip:port"
//   X-Forwarded-For – "client, proxy1, proxy2"
//   X-Real-IP      – "client"
//   Forwarded      – RFC7239 (for=client;by=proxy)
//   CloudFront-Viewer-Address – "client:port"
//   CF-Connecting-IP, True-Client-IP, X-Client-IP – "client"
//
// NOTE: These headers can be spoofed by a malicious client.  Only
// trust them if you know the request came through a trusted
// reverse‑proxy (e.g., your own ALB, CloudFront, etc.).
func RemoteIP(r *http.Request) string {
	// Ordered list of headers to inspect, from most specific to most general.
	const (
		headerXFF                = "X-Forwarded-For"
		headerXRealIP            = "X-Real-IP"
		headerCFConnectingIP     = "CF-Connecting-IP"
		headerTrueClientIP       = "True-Client-IP"
		headerXClientIP          = "X-Client-IP"
		headerForwarded          = "Forwarded"
		headerCloudFrontViewerIP = "CloudFront-Viewer-Address"
	)

	// Helper to try a header value and return the first valid IP it contains.
	tryHeader := func(name string) string {
		val := strings.TrimSpace(r.Header.Get(name))
		if val == "" {
			return ""
		}

		// CloudFront adds a port, so strip it.
		if name == headerCloudFrontViewerIP {
			if h, _, err := net.SplitHostPort(val); err == nil {
				val = h
			}
		}

		// X-Forwarded-For may be a comma‑separated list.
		if name == headerXFF {
			for _, part := range strings.Split(val, ",") {
				if ip := firstValidIP(strings.TrimSpace(part)); ip != "" {
					return ip
				}
			}
			return ""
		}

		// RFC7239 Forwarded header may contain multiple comma‑separated
		// items: "for=client;proto=http, for=proxy1".
		if name == headerForwarded {
			for _, item := range strings.Split(val, ",") {
				if ip := forwardedFor(item); ip != "" {
					return ip
				}
			}
			return ""
		}

		// All other headers contain a single IP.
		return firstValidIP(val)
	}

	// Order matters – we want the most reliable source first.
	if ip := tryHeader(headerXFF); ip != "" {
		return ip
	}
	if ip := tryHeader(headerXRealIP); ip != "" {
		return ip
	}
	if ip := tryHeader(headerCFConnectingIP); ip != "" {
		return ip
	}
	if ip := tryHeader(headerTrueClientIP); ip != "" {
		return ip
	}
	if ip := tryHeader(headerXClientIP); ip != "" {
		return ip
	}
	if ip := tryHeader(headerForwarded); ip != "" {
		return ip
	}
	if ip := tryHeader(headerCloudFrontViewerIP); ip != "" {
		return ip
	}

	// Fallback to the TCP remote address.
	if host, _, err := net.SplitHostPort(r.RemoteAddr); err == nil {
		return host
	}
	return r.RemoteAddr // probably already just an IP
}

// firstValidIP returns s if it parses as a valid IP address, otherwise "".
func firstValidIP(s string) string {
	if ip := net.ParseIP(s); ip != nil {
		return s
	}
	return ""
}

// forwardedFor extracts the "for=" value from a single RFC7239 Forwarded item.
// It handles quoted values and IPv6 addresses in brackets.
func forwardedFor(item string) string {
	for _, part := range strings.Split(strings.TrimSpace(item), ";") {
		part = strings.TrimSpace(part)
		if strings.HasPrefix(strings.ToLower(part), "for=") {
			val := strings.TrimPrefix(part, "for=")
			val = strings.Trim(val, "\"") // strip optional quotes
			// IPv6 addresses may be wrapped in brackets.
			if strings.HasPrefix(val, "[") && strings.HasSuffix(val, "]") {
				val = val[1 : len(val)-1]
			}
			if ip := net.ParseIP(val); ip != nil {
				return val
			}
		}
	}
	return ""
}
