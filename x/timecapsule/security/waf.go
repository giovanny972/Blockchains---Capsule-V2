package security

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// WAF (Web Application Firewall) provides comprehensive security protection
type WAF struct {
	mutex                sync.RWMutex
	ipWhitelist         map[string]bool
	ipBlacklist         map[string]bool
	rateLimiter         *RateLimiter
	attackDetector      *AttackDetector
	securityRules       []SecurityRule
	enabled             bool
	maxRequestSize      int64
	bannedUserAgents    []string
	suspiciousPatterns  []string
	geoBlockList        map[string]bool // Country codes to block
}

// SecurityRule defines a security rule for the WAF
type SecurityRule struct {
	ID          string                           `json:"id"`
	Name        string                           `json:"name"`
	Pattern     string                           `json:"pattern"`
	Action      string                           `json:"action"` // "block", "log", "rate_limit"
	Severity    string                           `json:"severity"` // "low", "medium", "high", "critical"
	Enabled     bool                             `json:"enabled"`
	Condition   func(*SecurityContext) bool     `json:"-"`
	Response    func(*SecurityContext) error    `json:"-"`
	Description string                           `json:"description"`
	LastTriggered *time.Time                    `json:"last_triggered,omitempty"`
	TriggerCount  int64                         `json:"trigger_count"`
}

// SecurityContext contains request context for security evaluation
type SecurityContext struct {
	IP            string
	UserAgent     string
	Path          string
	Method        string
	Headers       map[string]string
	Body          []byte
	BodySize      int64
	Timestamp     time.Time
	User          string
	TxHash        string
	BlockHeight   int64
	RequestID     string
	GeoLocation   *GeoLocation
}

// GeoLocation represents geographical location data
type GeoLocation struct {
	Country     string `json:"country"`
	CountryCode string `json:"country_code"`
	Region      string `json:"region"`
	City        string `json:"city"`
	ISP         string `json:"isp"`
	Threat      bool   `json:"threat"`
}

// RateLimiter manages request rate limiting
type RateLimiter struct {
	mutex         sync.RWMutex
	requests      map[string][]time.Time
	limits        map[string]RateLimit
	cleanupTicker *time.Ticker
}

// RateLimit defines rate limiting rules
type RateLimit struct {
	MaxRequests int           `json:"max_requests"`
	Window      time.Duration `json:"window"`
	BurstSize   int           `json:"burst_size"`
}

// AttackDetector identifies potential security attacks
type AttackDetector struct {
	mutex             sync.RWMutex
	suspiciousIPs     map[string]*SuspiciousActivity
	attackPatterns    []AttackPattern
	anomalyDetector   *AnomalyDetector
	enabled           bool
}

// SuspiciousActivity tracks suspicious behavior from IPs
type SuspiciousActivity struct {
	IP              string        `json:"ip"`
	FirstSeen       time.Time     `json:"first_seen"`
	LastSeen        time.Time     `json:"last_seen"`
	RequestCount    int64         `json:"request_count"`
	FailedRequests  int64         `json:"failed_requests"`
	AttackAttempts  int64         `json:"attack_attempts"`
	ThreatLevel     string        `json:"threat_level"` // "low", "medium", "high", "critical"
	BlockedUntil    *time.Time    `json:"blocked_until,omitempty"`
	Patterns        []string      `json:"patterns"`
	UserAgents      []string      `json:"user_agents"`
	GeoLocations    []GeoLocation `json:"geo_locations"`
}

// AttackPattern defines patterns that indicate attacks
type AttackPattern struct {
	Name        string `json:"name"`
	Pattern     string `json:"pattern"`
	Type        string `json:"type"` // "sql_injection", "xss", "path_traversal", "command_injection"
	Severity    string `json:"severity"`
	Enabled     bool   `json:"enabled"`
	Description string `json:"description"`
}

// AnomalyDetector detects abnormal behavior patterns
type AnomalyDetector struct {
	baseline         map[string]float64
	thresholds       map[string]float64
	anomalies        []Anomaly
	enabled          bool
	learningMode     bool
	learningPeriod   time.Duration
}

// Anomaly represents detected abnormal behavior
type Anomaly struct {
	Type        string    `json:"type"`
	Value       float64   `json:"value"`
	Expected    float64   `json:"expected"`
	Deviation   float64   `json:"deviation"`
	Severity    string    `json:"severity"`
	Timestamp   time.Time `json:"timestamp"`
	Context     string    `json:"context"`
	Resolved    bool      `json:"resolved"`
}

// NewWAF creates a new Web Application Firewall
func NewWAF() *WAF {
	waf := &WAF{
		ipWhitelist:        make(map[string]bool),
		ipBlacklist:        make(map[string]bool),
		rateLimiter:        NewRateLimiter(),
		attackDetector:     NewAttackDetector(),
		securityRules:      []SecurityRule{},
		enabled:            true,
		maxRequestSize:     10 * 1024 * 1024, // 10MB
		bannedUserAgents:   []string{},
		suspiciousPatterns: []string{},
		geoBlockList:       make(map[string]bool),
	}

	// Initialize default security rules
	waf.initializeDefaultRules()

	return waf
}

// NewRateLimiter creates a new rate limiter
func NewRateLimiter() *RateLimiter {
	rl := &RateLimiter{
		requests: make(map[string][]time.Time),
		limits: map[string]RateLimit{
			"default": {MaxRequests: 100, Window: time.Minute, BurstSize: 10},
			"create":  {MaxRequests: 10, Window: time.Minute, BurstSize: 2},
			"open":    {MaxRequests: 50, Window: time.Minute, BurstSize: 5},
		},
		cleanupTicker: time.NewTicker(5 * time.Minute),
	}

	// Start cleanup goroutine
	go rl.cleanup()

	return rl
}

// NewAttackDetector creates a new attack detector
func NewAttackDetector() *AttackDetector {
	ad := &AttackDetector{
		suspiciousIPs:   make(map[string]*SuspiciousActivity),
		attackPatterns:  []AttackPattern{},
		anomalyDetector: &AnomalyDetector{
			baseline:       make(map[string]float64),
			thresholds:     make(map[string]float64),
			anomalies:      []Anomaly{},
			enabled:        true,
			learningMode:   true,
			learningPeriod: 24 * time.Hour,
		},
		enabled: true,
	}

	// Initialize default attack patterns
	ad.initializeAttackPatterns()

	return ad
}

// ValidateRequest performs comprehensive request validation
func (w *WAF) ValidateRequest(ctx context.Context, secCtx *SecurityContext) error {
	if !w.enabled {
		return nil
	}

	w.mutex.RLock()
	defer w.mutex.RUnlock()

	// Check IP whitelist/blacklist
	if err := w.checkIPRestrictions(secCtx); err != nil {
		return err
	}

	// Check rate limits
	if err := w.rateLimiter.CheckLimit(secCtx); err != nil {
		return err
	}

	// Check request size
	if secCtx.BodySize > w.maxRequestSize {
		return fmt.Errorf("request size %d exceeds maximum %d", secCtx.BodySize, w.maxRequestSize)
	}

	// Check user agent restrictions
	if err := w.checkUserAgent(secCtx); err != nil {
		return err
	}

	// Check for attack patterns
	if err := w.attackDetector.AnalyzeRequest(secCtx); err != nil {
		return err
	}

	// Apply security rules
	for _, rule := range w.securityRules {
		if rule.Enabled && rule.Condition != nil {
			if rule.Condition(secCtx) {
				if err := rule.Response(secCtx); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

// checkIPRestrictions validates IP against whitelist/blacklist
func (w *WAF) checkIPRestrictions(secCtx *SecurityContext) error {
	ip := secCtx.IP

	// Check blacklist first
	if w.ipBlacklist[ip] {
		return fmt.Errorf("IP %s is blacklisted", ip)
	}

	// If whitelist exists and IP is not in it
	if len(w.ipWhitelist) > 0 && !w.ipWhitelist[ip] {
		return fmt.Errorf("IP %s is not whitelisted", ip)
	}

	// Check geo-blocking
	if secCtx.GeoLocation != nil {
		if w.geoBlockList[secCtx.GeoLocation.CountryCode] {
			return fmt.Errorf("requests from country %s are blocked", secCtx.GeoLocation.Country)
		}
	}

	return nil
}

// checkUserAgent validates user agent
func (w *WAF) checkUserAgent(secCtx *SecurityContext) error {
	userAgent := strings.ToLower(secCtx.UserAgent)

	for _, banned := range w.bannedUserAgents {
		if strings.Contains(userAgent, strings.ToLower(banned)) {
			return fmt.Errorf("user agent contains banned pattern: %s", banned)
		}
	}

	return nil
}

// CheckLimit validates rate limits for a request
func (rl *RateLimiter) CheckLimit(secCtx *SecurityContext) error {
	rl.mutex.Lock()
	defer rl.mutex.Unlock()

	key := secCtx.IP
	now := time.Now()

	// Determine which limit to apply
	limitKey := "default"
	if strings.Contains(secCtx.Path, "/create") {
		limitKey = "create"
	} else if strings.Contains(secCtx.Path, "/open") {
		limitKey = "open"
	}

	limit, exists := rl.limits[limitKey]
	if !exists {
		limit = rl.limits["default"]
	}

	// Get or create request history for this key
	if rl.requests[key] == nil {
		rl.requests[key] = []time.Time{}
	}

	// Remove old requests outside the time window
	var validRequests []time.Time
	cutoff := now.Add(-limit.Window)
	for _, reqTime := range rl.requests[key] {
		if reqTime.After(cutoff) {
			validRequests = append(validRequests, reqTime)
		}
	}

	// Check if limit would be exceeded
	if len(validRequests) >= limit.MaxRequests {
		return fmt.Errorf("rate limit exceeded: %d requests in %v window", len(validRequests), limit.Window)
	}

	// Add current request
	validRequests = append(validRequests, now)
	rl.requests[key] = validRequests

	return nil
}

// AnalyzeRequest analyzes request for attack patterns
func (ad *AttackDetector) AnalyzeRequest(secCtx *SecurityContext) error {
	if !ad.enabled {
		return nil
	}

	ad.mutex.Lock()
	defer ad.mutex.Unlock()

	// Track suspicious activity
	activity := ad.trackSuspiciousActivity(secCtx)

	// Check for attack patterns
	for _, pattern := range ad.attackPatterns {
		if pattern.Enabled {
			if ad.matchesPattern(secCtx, pattern) {
				activity.AttackAttempts++
				activity.Patterns = append(activity.Patterns, pattern.Name)

				if pattern.Severity == "critical" || pattern.Severity == "high" {
					return fmt.Errorf("attack pattern detected: %s", pattern.Name)
				}
			}
		}
	}

	// Update threat level
	ad.updateThreatLevel(activity)

	// Check if IP should be blocked
	if activity.ThreatLevel == "critical" || activity.AttackAttempts > 10 {
		blockUntil := time.Now().Add(time.Hour)
		activity.BlockedUntil = &blockUntil
		return fmt.Errorf("IP %s temporarily blocked due to suspicious activity", secCtx.IP)
	}

	return nil
}

// trackSuspiciousActivity tracks and updates suspicious activity for an IP
func (ad *AttackDetector) trackSuspiciousActivity(secCtx *SecurityContext) *SuspiciousActivity {
	ip := secCtx.IP
	now := time.Now()

	activity, exists := ad.suspiciousIPs[ip]
	if !exists {
		activity = &SuspiciousActivity{
			IP:           ip,
			FirstSeen:    now,
			LastSeen:     now,
			RequestCount: 0,
			ThreatLevel:  "low",
			Patterns:     []string{},
			UserAgents:   []string{},
		}
		ad.suspiciousIPs[ip] = activity
	}

	activity.LastSeen = now
	activity.RequestCount++

	// Add user agent if not already tracked
	userAgentExists := false
	for _, ua := range activity.UserAgents {
		if ua == secCtx.UserAgent {
			userAgentExists = true
			break
		}
	}
	if !userAgentExists {
		activity.UserAgents = append(activity.UserAgents, secCtx.UserAgent)
	}

	// Add geo location if available
	if secCtx.GeoLocation != nil {
		locationExists := false
		for _, loc := range activity.GeoLocations {
			if loc.CountryCode == secCtx.GeoLocation.CountryCode {
				locationExists = true
				break
			}
		}
		if !locationExists {
			activity.GeoLocations = append(activity.GeoLocations, *secCtx.GeoLocation)
		}
	}

	return activity
}

// matchesPattern checks if request matches an attack pattern
func (ad *AttackDetector) matchesPattern(secCtx *SecurityContext, pattern AttackPattern) bool {
	searchText := fmt.Sprintf("%s %s %s", secCtx.Path, secCtx.UserAgent, string(secCtx.Body))
	return strings.Contains(strings.ToLower(searchText), strings.ToLower(pattern.Pattern))
}

// updateThreatLevel updates the threat level based on activity
func (ad *AttackDetector) updateThreatLevel(activity *SuspiciousActivity) {
	if activity.AttackAttempts > 20 {
		activity.ThreatLevel = "critical"
	} else if activity.AttackAttempts > 10 {
		activity.ThreatLevel = "high"
	} else if activity.AttackAttempts > 5 {
		activity.ThreatLevel = "medium"
	} else {
		activity.ThreatLevel = "low"
	}
}

// cleanup removes old rate limiting data
func (rl *RateLimiter) cleanup() {
	for range rl.cleanupTicker.C {
		rl.mutex.Lock()
		now := time.Now()
		for key, requests := range rl.requests {
			var validRequests []time.Time
			cutoff := now.Add(-time.Hour) // Keep last hour of data
			for _, reqTime := range requests {
				if reqTime.After(cutoff) {
					validRequests = append(validRequests, reqTime)
				}
			}
			if len(validRequests) == 0 {
				delete(rl.requests, key)
			} else {
				rl.requests[key] = validRequests
			}
		}
		rl.mutex.Unlock()
	}
}

// initializeDefaultRules sets up default security rules
func (w *WAF) initializeDefaultRules() {
	// SQL Injection protection
	w.securityRules = append(w.securityRules, SecurityRule{
		ID:       "sql_injection",
		Name:     "SQL Injection Protection",
		Pattern:  `(?i)(union|select|insert|update|delete|drop|create|alter)\s`,
		Action:   "block",
		Severity: "high",
		Enabled:  true,
		Condition: func(ctx *SecurityContext) bool {
			content := string(ctx.Body) + ctx.Path + ctx.UserAgent
			return strings.Contains(strings.ToLower(content), "union") ||
				strings.Contains(strings.ToLower(content), "select") ||
				strings.Contains(strings.ToLower(content), "drop")
		},
		Response: func(ctx *SecurityContext) error {
			return fmt.Errorf("SQL injection attempt detected from IP: %s", ctx.IP)
		},
		Description: "Protects against SQL injection attacks",
	})

	// XSS protection
	w.securityRules = append(w.securityRules, SecurityRule{
		ID:       "xss_protection",
		Name:     "XSS Protection",
		Pattern:  `(?i)<script|javascript:|on\w+\s*=`,
		Action:   "block",
		Severity: "medium",
		Enabled:  true,
		Condition: func(ctx *SecurityContext) bool {
			content := string(ctx.Body) + ctx.Path
			return strings.Contains(strings.ToLower(content), "<script") ||
				strings.Contains(strings.ToLower(content), "javascript:")
		},
		Response: func(ctx *SecurityContext) error {
			return fmt.Errorf("XSS attempt detected from IP: %s", ctx.IP)
		},
		Description: "Protects against cross-site scripting attacks",
	})

	// Excessive request size
	w.securityRules = append(w.securityRules, SecurityRule{
		ID:       "large_request",
		Name:     "Large Request Protection",
		Action:   "block",
		Severity: "medium",
		Enabled:  true,
		Condition: func(ctx *SecurityContext) bool {
			return ctx.BodySize > 50*1024*1024 // 50MB
		},
		Response: func(ctx *SecurityContext) error {
			return fmt.Errorf("request size too large: %d bytes from IP: %s", ctx.BodySize, ctx.IP)
		},
		Description: "Blocks excessively large requests",
	})
}

// initializeAttackPatterns sets up default attack patterns
func (ad *AttackDetector) initializeAttackPatterns() {
	ad.attackPatterns = []AttackPattern{
		{
			Name:        "SQL Injection",
			Pattern:     "union select",
			Type:        "sql_injection",
			Severity:    "high",
			Enabled:     true,
			Description: "Detects SQL injection attempts",
		},
		{
			Name:        "XSS Attack",
			Pattern:     "<script>",
			Type:        "xss",
			Severity:    "medium",
			Enabled:     true,
			Description: "Detects cross-site scripting attempts",
		},
		{
			Name:        "Path Traversal",
			Pattern:     "../",
			Type:        "path_traversal",
			Severity:    "high",
			Enabled:     true,
			Description: "Detects directory traversal attempts",
		},
		{
			Name:        "Command Injection",
			Pattern:     "; rm -rf",
			Type:        "command_injection",
			Severity:    "critical",
			Enabled:     true,
			Description: "Detects command injection attempts",
		},
	}
}

// AddIPToWhitelist adds an IP to the whitelist
func (w *WAF) AddIPToWhitelist(ip string) {
	w.mutex.Lock()
	defer w.mutex.Unlock()
	w.ipWhitelist[ip] = true
}

// AddIPToBlacklist adds an IP to the blacklist
func (w *WAF) AddIPToBlacklist(ip string) {
	w.mutex.Lock()
	defer w.mutex.Unlock()
	w.ipBlacklist[ip] = true
}

// GetSecurityStats returns current security statistics
func (w *WAF) GetSecurityStats() map[string]interface{} {
	w.mutex.RLock()
	defer w.mutex.RUnlock()

	stats := make(map[string]interface{})
	stats["enabled"] = w.enabled
	stats["whitelisted_ips"] = len(w.ipWhitelist)
	stats["blacklisted_ips"] = len(w.ipBlacklist)
	stats["security_rules"] = len(w.securityRules)
	stats["suspicious_ips"] = len(w.attackDetector.suspiciousIPs)

	// Count blocked IPs
	blockedCount := 0
	for _, activity := range w.attackDetector.suspiciousIPs {
		if activity.BlockedUntil != nil && time.Now().Before(*activity.BlockedUntil) {
			blockedCount++
		}
	}
	stats["currently_blocked_ips"] = blockedCount

	return stats
}

// EnableWAF enables the WAF
func (w *WAF) EnableWAF() {
	w.mutex.Lock()
	defer w.mutex.Unlock()
	w.enabled = true
}

// DisableWAF disables the WAF
func (w *WAF) DisableWAF() {
	w.mutex.Lock()
	defer w.mutex.Unlock()
	w.enabled = false
}