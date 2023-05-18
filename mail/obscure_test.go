package mail

import (
	"testing"
)

func TestObscureAddress(t *testing.T) {

	tests := map[string]string{
		"bob@example.com":            "b***b@e***e.com",
		"susan.taylor@company.co.uk": "s***r@c***y.***.uk",
	}

	for raw, expected := range tests {

		obscured := ObscureAddress(raw)

		if obscured != expected {
			t.Fatalf("Failed to obscure '%s'. Expected '%s' but got '%s'", raw, expected, obscured)
		}
	}
}
