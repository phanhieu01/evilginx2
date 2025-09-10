package core

import (
	"net/http"
	"testing"
)

func TestTokenModifier_NewTokenModifier(t *testing.T) {
	tm := NewTokenModifier()
	if tm == nil {
		t.Fatal("NewTokenModifier returned nil")
	}
	if tm.IsEnabled() {
		t.Error("Token modifier should be disabled by default")
	}
}

func TestTokenModifier_EnableDisable(t *testing.T) {
	tm := NewTokenModifier()
	
	// Test enable
	tm.Enable()
	if !tm.IsEnabled() {
		t.Error("Token modifier should be enabled after calling Enable()")
	}
	
	// Test disable
	tm.Disable()
	if tm.IsEnabled() {
		t.Error("Token modifier should be disabled after calling Disable()")
	}
}

func TestTokenModifier_AddRule(t *testing.T) {
	tm := NewTokenModifier()
	
	err := tm.AddRule("^Bearer\\s+(.+)$", "prefix", "modified_", "Test Bearer rule", TokenLocationHeader, "Authorization")
	if err != nil {
		t.Fatalf("AddRule failed: %v", err)
	}
	
	rules := tm.ListRules()
	found := false
	for _, rule := range rules {
		if rule.Description == "Test Bearer rule" {
			found = true
			break
		}
	}
	if !found {
		t.Error("Added rule not found in list")
	}
}

func TestTokenModifier_ModifyBearerToken(t *testing.T) {
	tm := NewTokenModifier()
	tm.Enable()
	
	// Add rule to prefix Bearer tokens
	err := tm.AddRule("^Bearer\\s+(.+)$", "replace", "Bearer modified_$1", "Test Bearer rule", TokenLocationHeader, "Authorization")
	if err != nil {
		t.Fatalf("AddRule failed: %v", err)
	}
	
	// Create test request with Bearer token
	req, _ := http.NewRequest("GET", "http://example.com", nil)
	req.Header.Set("Authorization", "Bearer abc123token")
	
	// Modify request
	tm.ModifyRequest(req)
	
	// Check if token was modified
	authHeader := req.Header.Get("Authorization")
	expected := "Bearer modified_abc123token"
	if authHeader != expected {
		t.Errorf("Expected %s, got %s", expected, authHeader)
	}
}

func TestTokenModifier_ModifyJWTToken(t *testing.T) {
	tm := NewTokenModifier()
	tm.Enable()
	
	// Add JWT prefix rule
	err := tm.SetJWTPrefix("modified_")
	if err != nil {
		t.Fatalf("SetJWTPrefix failed: %v", err)
	}
	
	// Create test request with JWT token
	req, _ := http.NewRequest("GET", "http://example.com", nil)
	req.Header.Set("Authorization", "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c")
	
	// Modify request
	tm.ModifyRequest(req)
	
	// Check if token was modified
	authHeader := req.Header.Get("Authorization")
	if authHeader == "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c" {
		t.Error("JWT token was not modified")
	}
}

func TestTokenModifier_ModifyQueryParams(t *testing.T) {
	tm := NewTokenModifier()
	tm.Enable()
	
	// Add rule for query parameter token
	err := tm.AddRule("^(.+)$", "prefix", "modified_", "Test query token", TokenLocationQueryParam, "token")
	if err != nil {
		t.Fatalf("AddRule failed: %v", err)
	}
	
	// Create test request with query parameter
	req, _ := http.NewRequest("GET", "http://example.com?token=abc123", nil)
	
	// Modify request
	tm.ModifyRequest(req)
	
	// Check if query parameter was modified
	queryValues := req.URL.Query()
	tokenValue := queryValues.Get("token")
	expected := "modified_abc123"
	if tokenValue != expected {
		t.Errorf("Expected %s, got %s", expected, tokenValue)
	}
}

func TestTokenModifier_ModifyCookies(t *testing.T) {
	tm := NewTokenModifier()
	tm.Enable()
	
	// Add rule for cookies
	err := tm.AddRule("^(.+)$", "prefix", "modified_", "Test cookie rule", TokenLocationCookie, "")
	if err != nil {
		t.Fatalf("AddRule failed: %v", err)
	}
	
	// Create test request with cookie
	req, _ := http.NewRequest("GET", "http://example.com", nil)
	req.AddCookie(&http.Cookie{Name: "session", Value: "abc123"})
	
	// Modify request
	tm.ModifyRequest(req)
	
	// Check if cookie was modified
	cookies := req.Cookies()
	found := false
	for _, cookie := range cookies {
		if cookie.Name == "session" && cookie.Value == "modified_abc123" {
			found = true
			break
		}
	}
	if !found {
		t.Error("Cookie value was not modified correctly")
	}
}

func TestTokenModifier_DisabledRules(t *testing.T) {
	tm := NewTokenModifier()
	tm.Enable()
	
	// Add rule but keep it disabled
	err := tm.AddRule("^Bearer\\s+(.+)$", "replace", "Bearer modified_$1", "Disabled rule", TokenLocationHeader, "Authorization")
	if err != nil {
		t.Fatalf("AddRule failed: %v", err)
	}
	
	// Disable the rule
	tm.DisableRule("Disabled rule")
	
	// Create test request
	req, _ := http.NewRequest("GET", "http://example.com", nil)
	originalToken := "Bearer abc123token"
	req.Header.Set("Authorization", originalToken)
	
	// Modify request
	tm.ModifyRequest(req)
	
	// Check that token was NOT modified (rule is disabled)
	authHeader := req.Header.Get("Authorization")
	if authHeader != originalToken {
		t.Errorf("Token should not be modified when rule is disabled. Expected %s, got %s", originalToken, authHeader)
	}
}