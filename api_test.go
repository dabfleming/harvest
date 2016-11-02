package harvest

import (
	"testing"
)

func TestNewBasicAuthAPI(t *testing.T) {
	a1 := NewBasicAuthAPI("example", "user@example.com", "password")
	if a1.BaseURL != "https://example.harvestapp.com" {
		t.Errorf("Incorrect domain name '%s'.", a1.BaseURL)
	}
}
