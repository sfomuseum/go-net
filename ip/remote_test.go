package ip

// Disclosure: The following was "vibe-coded" using gpt-oss-20b (20260303)
// This code has been reviewed and doesn't raise any obvious red flags.

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// helper that builds a request with a given remote address and a set of headers
func newTestRequest(remoteAddr string, headers map[string]string) *http.Request {
	req := httptest.NewRequest(http.MethodGet, "http://example.com", nil)
	req.RemoteAddr = remoteAddr
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	return req
}

func TestRemoteIP(t *testing.T) {
	// a set of IPs that we will use in the tests
	const (
		clientIP  = "192.0.2.1"
		proxyIP   = "198.51.100.2"
		remoteIP  = "203.0.113.5:12345" // RemoteAddr format
		ipv6Client = "2001:db8::1"
	)

	tests := []struct {
		name       string
		remoteAddr string
		headers    map[string]string
		wantIP     string
	}{
		{
			name:       "X-Forwarded-For single IP",
			remoteAddr: remoteIP,
			headers: map[string]string{
				"X-Forwarded-For": clientIP,
			},
			wantIP: clientIP,
		},
		{
			name:       "X-Forwarded-For multiple IPs",
			remoteAddr: remoteIP,
			headers: map[string]string{
				"X-Forwarded-For": clientIP + ", " + proxyIP,
			},
			wantIP: clientIP,
		},
		{
			name:       "X-Real-IP overrides when XFF absent",
			remoteAddr: remoteIP,
			headers: map[string]string{
				"X-Real-IP": clientIP,
			},
			wantIP: clientIP,
		},
		{
			name:       "CF-Connecting-IP",
			remoteAddr: remoteIP,
			headers: map[string]string{
				"CF-Connecting-IP": clientIP,
			},
			wantIP: clientIP,
		},
		{
			name:       "True-Client-IP",
			remoteAddr: remoteIP,
			headers: map[string]string{
				"True-Client-IP": clientIP,
			},
			wantIP: clientIP,
		},
		{
			name:       "X-Client-IP",
			remoteAddr: remoteIP,
			headers: map[string]string{
				"X-Client-IP": clientIP,
			},
			wantIP: clientIP,
		},
		{
			name:       "Forwarded header single item",
			remoteAddr: remoteIP,
			headers: map[string]string{
				"Forwarded": "for=" + clientIP + ";proto=http;by=" + proxyIP,
			},
			wantIP: clientIP,
		},
		{
			name:       "Forwarded header multiple items",
			remoteAddr: remoteIP,
			headers: map[string]string{
				"Forwarded": "for=" + clientIP + ", for=" + proxyIP,
			},
			wantIP: clientIP,
		},
		{
			name:       "Forwarded header quoted IPv6",
			remoteAddr: remoteIP,
			headers: map[string]string{
				"Forwarded": `for="[2001:db8::1]";proto=https`,
			},
			wantIP: ipv6Client,
		},
		{
			name:       "CloudFront-Viewer-Address with port",
			remoteAddr: remoteIP,
			headers: map[string]string{
				"CloudFront-Viewer-Address": clientIP + ":443",
			},
			wantIP: clientIP,
		},
		{
			name:       "RemoteAddr fallback",
			remoteAddr: remoteIP,
			headers:    map[string]string{},
			wantIP:     "203.0.113.5",
		},
		{
			name:       "Invalid XFF value falls back to next header",
			remoteAddr: remoteIP,
			headers: map[string]string{
				"X-Forwarded-For": "invalid-ip",
				"X-Real-IP":       clientIP,
			},
			wantIP: clientIP,
		},
		{
			name:       "All headers invalid, fallback to RemoteAddr",
			remoteAddr: remoteIP,
			headers: map[string]string{
				"X-Forwarded-For": "not-an-ip",
				"CF-Connecting-IP": "also-not-an-ip",
			},
			wantIP: "203.0.113.5",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req := newTestRequest(tc.remoteAddr, tc.headers)
			got := RemoteIP(req)
			if got != tc.wantIP {
				t.Fatalf("RemoteIP() = %q, want %q", got, tc.wantIP)
			}
		})
	}
}
