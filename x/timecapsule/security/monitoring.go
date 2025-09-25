package security

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// SecurityMonitor provides comprehensive security monitoring and alerting
type SecurityMonitor struct {
	mutex               sync.RWMutex
	enabled             bool
	alertRules          []AlertRule
	metrics             *SecurityMetrics
	eventCollector      *EventCollector
	alertQueue          chan Alert
	logQueue            chan SecurityLog
	notificationService *NotificationService
	auditTrail          *AuditTrail
	thresholds          map[string]float64
	anomalyDetector     *RealTimeAnomalyDetector
}

// AlertRule defines conditions for triggering security alerts
type AlertRule struct {
	ID            string                    `json:"id"`
	Name          string                    `json:"name"`
	Description   string                    `json:"description"`
	Severity      string                    `json:"severity"` // "info", "warning", "error", "critical"
	Enabled       bool                      `json:"enabled"`
	Condition     func(*SecurityEvent) bool `json:"-"`
	Action        func(*SecurityEvent) error `json:"-"`
	Cooldown      time.Duration             `json:"cooldown"`
	LastTriggered *time.Time                `json:"last_triggered,omitempty"`
	TriggerCount  int64                     `json:"trigger_count"`
	Tags          []string                  `json:"tags"`
}

// SecurityMetrics tracks security-related metrics
type SecurityMetrics struct {
	mutex                    sync.RWMutex
	TotalRequests           int64     `json:"total_requests"`
	BlockedRequests         int64     `json:"blocked_requests"`
	SuspiciousActivity      int64     `json:"suspicious_activity"`
	FailedAuthentications   int64     `json:"failed_authentications"`
	RateLimitViolations     int64     `json:"rate_limit_violations"`
	AttackAttempts          int64     `json:"attack_attempts"`
	DataBreachAttempts      int64     `json:"data_breach_attempts"`
	UnauthorizedAccess      int64     `json:"unauthorized_access"`
	EncryptionFailures      int64     `json:"encryption_failures"`
	IntegrityViolations     int64     `json:"integrity_violations"`
	AverageResponseTime     float64   `json:"average_response_time"`
	SecurityScore           float64   `json:"security_score"`
	ThreatLevel             string    `json:"threat_level"`
	LastSecurityScan        time.Time `json:"last_security_scan"`
	VulnerabilitiesFound    int32     `json:"vulnerabilities_found"`
	SecurityPatchesApplied  int32     `json:"security_patches_applied"`
	ComplianceScore         float64   `json:"compliance_score"`
}

// SecurityEvent represents a security-related event
type SecurityEvent struct {
	ID          string                 `json:"id"`
	Type        string                 `json:"type"`
	Source      string                 `json:"source"`
	Target      string                 `json:"target"`
	Severity    string                 `json:"severity"`
	Timestamp   time.Time              `json:"timestamp"`
	Message     string                 `json:"message"`
	Details     map[string]interface{} `json:"details"`
	IP          string                 `json:"ip,omitempty"`
	UserAgent   string                 `json:"user_agent,omitempty"`
	User        string                 `json:"user,omitempty"`
	Action      string                 `json:"action"`
	Resource    string                 `json:"resource,omitempty"`
	Outcome     string                 `json:"outcome"` // "success", "failure", "blocked"
	RiskScore   float64                `json:"risk_score"`
	Indicators  []ThreatIndicator      `json:"indicators,omitempty"`
	Context     *SecurityContext       `json:"context,omitempty"`
}

// ThreatIndicator represents an indicator of compromise
type ThreatIndicator struct {
	Type        string    `json:"type"` // "ip", "domain", "hash", "pattern"
	Value       string    `json:"value"`
	Source      string    `json:"source"`
	Confidence  float64   `json:"confidence"`
	FirstSeen   time.Time `json:"first_seen"`
	LastSeen    time.Time `json:"last_seen"`
	Description string    `json:"description"`
	Tags        []string  `json:"tags"`
}

// Alert represents a security alert
type Alert struct {
	ID          string                 `json:"id"`
	RuleID      string                 `json:"rule_id"`
	Severity    string                 `json:"severity"`
	Title       string                 `json:"title"`
	Message     string                 `json:"message"`
	Timestamp   time.Time              `json:"timestamp"`
	Event       *SecurityEvent         `json:"event"`
	Status      string                 `json:"status"` // "new", "acknowledged", "resolved", "false_positive"
	AssignedTo  string                 `json:"assigned_to,omitempty"`
	Actions     []AlertAction          `json:"actions,omitempty"`
	Context     map[string]interface{} `json:"context"`
	Escalated   bool                   `json:"escalated"`
	ResolvedAt  *time.Time             `json:"resolved_at,omitempty"`
}

// AlertAction represents an automated action taken in response to an alert
type AlertAction struct {
	Type        string                 `json:"type"` // "block_ip", "notify", "escalate", "quarantine"
	Timestamp   time.Time              `json:"timestamp"`
	Status      string                 `json:"status"` // "pending", "completed", "failed"
	Details     map[string]interface{} `json:"details"`
	Error       string                 `json:"error,omitempty"`
}

// SecurityLog represents a detailed security log entry
type SecurityLog struct {
	ID          string                 `json:"id"`
	Level       string                 `json:"level"` // "debug", "info", "warn", "error", "fatal"
	Module      string                 `json:"module"`
	Function    string                 `json:"function"`
	Message     string                 `json:"message"`
	Timestamp   time.Time              `json:"timestamp"`
	Data        map[string]interface{} `json:"data,omitempty"`
	StackTrace  string                 `json:"stack_trace,omitempty"`
	RequestID   string                 `json:"request_id,omitempty"`
	UserID      string                 `json:"user_id,omitempty"`
	SessionID   string                 `json:"session_id,omitempty"`
	IP          string                 `json:"ip,omitempty"`
	Tags        []string               `json:"tags,omitempty"`
}

// EventCollector aggregates and processes security events
type EventCollector struct {
	mutex         sync.RWMutex
	events        []SecurityEvent
	maxEvents     int
	retentionTime time.Duration
	filters       []EventFilter
	processors    []EventProcessor
}

// EventFilter defines filtering criteria for events
type EventFilter struct {
	Name      string                         `json:"name"`
	Enabled   bool                           `json:"enabled"`
	Condition func(*SecurityEvent) bool      `json:"-"`
}

// EventProcessor processes events for analysis
type EventProcessor struct {
	Name    string                           `json:"name"`
	Enabled bool                             `json:"enabled"`
	Process func(*SecurityEvent) error       `json:"-"`
}

// AuditTrail maintains an immutable record of security events
type AuditTrail struct {
	mutex     sync.RWMutex
	entries   []AuditEntry
	enabled   bool
	retention time.Duration
}

// AuditEntry represents an immutable audit log entry
type AuditEntry struct {
	ID          string                 `json:"id"`
	Timestamp   time.Time              `json:"timestamp"`
	Actor       string                 `json:"actor"`
	Action      string                 `json:"action"`
	Resource    string                 `json:"resource"`
	Outcome     string                 `json:"outcome"`
	Details     map[string]interface{} `json:"details"`
	IPAddress   string                 `json:"ip_address"`
	UserAgent   string                 `json:"user_agent"`
	Signature   string                 `json:"signature"` // Cryptographic signature for integrity
	Hash        string                 `json:"hash"`      // Hash of the entry for verification
}

// RealTimeAnomalyDetector detects anomalies in real-time
type RealTimeAnomalyDetector struct {
	mutex           sync.RWMutex
	enabled         bool
	models          map[string]*AnomalyModel
	threshold       float64
	learningPeriod  time.Duration
	detectionWindow time.Duration
	anomalies       []Anomaly
}

// AnomalyModel represents a machine learning model for anomaly detection
type AnomalyModel struct {
	Name          string             `json:"name"`
	Type          string             `json:"type"` // "statistical", "ml", "behavioral"
	Trained       bool               `json:"trained"`
	Accuracy      float64            `json:"accuracy"`
	LastUpdated   time.Time          `json:"last_updated"`
	Parameters    map[string]float64 `json:"parameters"`
	TrainingData  []DataPoint        `json:"training_data"`
}

// DataPoint represents a data point for training
type DataPoint struct {
	Timestamp  time.Time              `json:"timestamp"`
	Features   map[string]float64     `json:"features"`
	Label      string                 `json:"label"`
	Metadata   map[string]interface{} `json:"metadata"`
}

// NotificationChannel represents a notification delivery channel
type NotificationChannel struct {
	ID          string                 `json:"id"`
	Type        string                 `json:"type"` // "email", "slack", "webhook", "sms"
	Enabled     bool                   `json:"enabled"`
	Config      map[string]interface{} `json:"config"`
	Filters     []string               `json:"filters"` // Severity levels to notify
	RateLimit   *RateLimit             `json:"rate_limit,omitempty"`
}

// NewSecurityMonitor creates a new security monitor
func NewSecurityMonitor() *SecurityMonitor {
	monitor := &SecurityMonitor{
		enabled:         true,
		alertRules:      []AlertRule{},
		metrics:         &SecurityMetrics{},
		eventCollector:  NewEventCollector(),
		alertQueue:      make(chan Alert, 1000),
		logQueue:        make(chan SecurityLog, 1000),
		auditTrail:      NewAuditTrail(),
		thresholds:      make(map[string]float64),
		anomalyDetector: NewRealTimeAnomalyDetector(),
	}

	// Initialize default alert rules
	monitor.initializeDefaultAlertRules()

	// Initialize default thresholds
	monitor.initializeDefaultThresholds()

	// Start background workers
	go monitor.processAlerts()
	go monitor.processLogs()

	return monitor
}

// NewEventCollector creates a new event collector
func NewEventCollector() *EventCollector {
	return &EventCollector{
		events:        []SecurityEvent{},
		maxEvents:     10000,
		retentionTime: 30 * 24 * time.Hour, // 30 days
		filters:       []EventFilter{},
		processors:    []EventProcessor{},
	}
}

// NewAuditTrail creates a new audit trail
func NewAuditTrail() *AuditTrail {
	return &AuditTrail{
		entries:   []AuditEntry{},
		enabled:   true,
		retention: 365 * 24 * time.Hour, // 1 year
	}
}

// NewRealTimeAnomalyDetector creates a new real-time anomaly detector
func NewRealTimeAnomalyDetector() *RealTimeAnomalyDetector {
	return &RealTimeAnomalyDetector{
		enabled:         true,
		models:          make(map[string]*AnomalyModel),
		threshold:       0.8,
		learningPeriod:  7 * 24 * time.Hour, // 7 days
		detectionWindow: time.Hour,
		anomalies:       []Anomaly{},
	}
}

// CollectEvent collects a security event for processing
func (sm *SecurityMonitor) CollectEvent(event *SecurityEvent) error {
	if !sm.enabled {
		return nil
	}

	// Add to event collector
	if err := sm.eventCollector.AddEvent(event); err != nil {
		return fmt.Errorf("failed to add event to collector: %w", err)
	}

	// Check alert rules
	sm.checkAlertRules(event)

	// Update metrics
	sm.updateMetrics(event)

	// Create audit entry
	if err := sm.auditTrail.AddEntry(event); err != nil {
		sm.LogError("failed to add audit entry", map[string]interface{}{
			"event_id": event.ID,
			"error":    err.Error(),
		})
	}

	// Check for anomalies
	if err := sm.anomalyDetector.ProcessEvent(event); err != nil {
		sm.LogWarning("anomaly detection failed", map[string]interface{}{
			"event_id": event.ID,
			"error":    err.Error(),
		})
	}

	return nil
}

// AddEvent adds an event to the collector
func (ec *EventCollector) AddEvent(event *SecurityEvent) error {
	ec.mutex.Lock()
	defer ec.mutex.Unlock()

	// Apply filters
	for _, filter := range ec.filters {
		if filter.Enabled && filter.Condition != nil {
			if !filter.Condition(event) {
				return nil // Event filtered out
			}
		}
	}

	// Add event
	ec.events = append(ec.events, *event)

	// Trim if exceeding max events
	if len(ec.events) > ec.maxEvents {
		ec.events = ec.events[len(ec.events)-ec.maxEvents:]
	}

	// Process event
	for _, processor := range ec.processors {
		if processor.Enabled && processor.Process != nil {
			if err := processor.Process(event); err != nil {
				// Log error but continue processing
				fmt.Printf("Event processor %s failed: %v\n", processor.Name, err)
			}
		}
	}

	return nil
}

// checkAlertRules checks if any alert rules are triggered by the event
func (sm *SecurityMonitor) checkAlertRules(event *SecurityEvent) {
	sm.mutex.RLock()
	defer sm.mutex.RUnlock()

	for i := range sm.alertRules {
		rule := &sm.alertRules[i]
		if rule.Enabled && rule.Condition != nil {
			// Check cooldown period
			if rule.LastTriggered != nil {
				if time.Since(*rule.LastTriggered) < rule.Cooldown {
					continue
				}
			}

			if rule.Condition(event) {
				// Create alert
				alert := Alert{
					ID:        fmt.Sprintf("alert-%d-%s", time.Now().UnixNano(), rule.ID),
					RuleID:    rule.ID,
					Severity:  rule.Severity,
					Title:     rule.Name,
					Message:   rule.Description,
					Timestamp: time.Now(),
					Event:     event,
					Status:    "new",
					Context:   make(map[string]interface{}),
				}

				// Send to alert queue
				select {
				case sm.alertQueue <- alert:
					now := time.Now()
					rule.LastTriggered = &now
					rule.TriggerCount++

					// Execute rule action if defined
					if rule.Action != nil {
						if err := rule.Action(event); err != nil {
							sm.LogError("alert rule action failed", map[string]interface{}{
								"rule_id": rule.ID,
								"error":   err.Error(),
							})
						}
					}
				default:
					sm.LogWarning("alert queue full, dropping alert", map[string]interface{}{
						"rule_id": rule.ID,
						"event_id": event.ID,
					})
				}
			}
		}
	}
}

// updateMetrics updates security metrics based on the event
func (sm *SecurityMonitor) updateMetrics(event *SecurityEvent) {
	sm.metrics.mutex.Lock()
	defer sm.metrics.mutex.Unlock()

	sm.metrics.TotalRequests++

	switch event.Type {
	case "request_blocked":
		sm.metrics.BlockedRequests++
	case "suspicious_activity":
		sm.metrics.SuspiciousActivity++
	case "authentication_failed":
		sm.metrics.FailedAuthentications++
	case "rate_limit_violation":
		sm.metrics.RateLimitViolations++
	case "attack_attempt":
		sm.metrics.AttackAttempts++
	case "data_breach_attempt":
		sm.metrics.DataBreachAttempts++
	case "unauthorized_access":
		sm.metrics.UnauthorizedAccess++
	case "encryption_failure":
		sm.metrics.EncryptionFailures++
	case "integrity_violation":
		sm.metrics.IntegrityViolations++
	}

	// Update threat level based on recent activity
	sm.updateThreatLevel()

	// Calculate security score
	sm.calculateSecurityScore()
}

// updateThreatLevel updates the overall threat level
func (sm *SecurityMonitor) updateThreatLevel() {
	// Simple algorithm based on recent attack attempts
	recentAttacks := sm.metrics.AttackAttempts + sm.metrics.SuspiciousActivity + sm.metrics.UnauthorizedAccess

	if recentAttacks > 100 {
		sm.metrics.ThreatLevel = "critical"
	} else if recentAttacks > 50 {
		sm.metrics.ThreatLevel = "high"
	} else if recentAttacks > 20 {
		sm.metrics.ThreatLevel = "medium"
	} else {
		sm.metrics.ThreatLevel = "low"
	}
}

// calculateSecurityScore calculates an overall security score
func (sm *SecurityMonitor) calculateSecurityScore() {
	if sm.metrics.TotalRequests == 0 {
		sm.metrics.SecurityScore = 100.0
		return
	}

	// Calculate based on various factors
	blockRate := float64(sm.metrics.BlockedRequests) / float64(sm.metrics.TotalRequests)
	attackRate := float64(sm.metrics.AttackAttempts) / float64(sm.metrics.TotalRequests)
	failureRate := float64(sm.metrics.FailedAuthentications) / float64(sm.metrics.TotalRequests)

	// Score starts at 100 and decreases based on negative events
	score := 100.0
	score -= blockRate * 20.0       // Blocked requests reduce score
	score -= attackRate * 30.0      // Attack attempts reduce score more
	score -= failureRate * 15.0     // Failed authentications reduce score

	if score < 0 {
		score = 0
	}

	sm.metrics.SecurityScore = score
}

// processAlerts processes alerts from the queue
func (sm *SecurityMonitor) processAlerts() {
	for alert := range sm.alertQueue {
		// Process the alert (send notifications, log, etc.)
		sm.processAlert(&alert)
	}
}

// processAlert processes a single alert
func (sm *SecurityMonitor) processAlert(alert *Alert) {
	// Log the alert
	sm.LogWarning("Security alert triggered", map[string]interface{}{
		"alert_id":   alert.ID,
		"rule_id":    alert.RuleID,
		"severity":   alert.Severity,
		"title":      alert.Title,
		"event_type": alert.Event.Type,
	})

	// Send notifications based on severity
	if sm.notificationService != nil {
		sm.notificationService.SendAlert(alert)
	}

	// Auto-escalate critical alerts
	if alert.Severity == "critical" {
		alert.Escalated = true
		sm.LogError("Critical security alert auto-escalated", map[string]interface{}{
			"alert_id": alert.ID,
			"title":    alert.Title,
		})
	}
}

// processLogs processes logs from the queue
func (sm *SecurityMonitor) processLogs() {
	for log := range sm.logQueue {
		// Process the log (store, analyze, etc.)
		sm.processLog(&log)
	}
}

// processLog processes a single log entry
func (sm *SecurityMonitor) processLog(log *SecurityLog) {
	// In a real implementation, this would:
	// - Store logs in a persistent storage system
	// - Index logs for searching
	// - Apply log retention policies
	// - Send to external log aggregation systems

	// For now, just print to console for critical errors
	if log.Level == "error" || log.Level == "fatal" {
		fmt.Printf("[%s] %s: %s\n", log.Level, log.Module, log.Message)
	}
}

// LogInfo logs an informational message
func (sm *SecurityMonitor) LogInfo(message string, data map[string]interface{}) {
	sm.logMessage("info", message, data)
}

// LogWarning logs a warning message
func (sm *SecurityMonitor) LogWarning(message string, data map[string]interface{}) {
	sm.logMessage("warn", message, data)
}

// LogError logs an error message
func (sm *SecurityMonitor) LogError(message string, data map[string]interface{}) {
	sm.logMessage("error", message, data)
}

// logMessage creates and queues a log message
func (sm *SecurityMonitor) logMessage(level, message string, data map[string]interface{}) {
	log := SecurityLog{
		ID:        fmt.Sprintf("log-%d", time.Now().UnixNano()),
		Level:     level,
		Module:    "security_monitor",
		Message:   message,
		Timestamp: time.Now(),
		Data:      data,
	}

	select {
	case sm.logQueue <- log:
	default:
		// Log queue is full, print to console as fallback
		fmt.Printf("[%s] %s: %s\n", level, "security_monitor", message)
	}
}

// GetSecurityMetrics returns current security metrics
func (sm *SecurityMonitor) GetSecurityMetrics() *SecurityMetrics {
	sm.metrics.mutex.RLock()
	defer sm.metrics.mutex.RUnlock()

	// Create a copy to avoid race conditions
	metricsCopy := *sm.metrics
	return &metricsCopy
}

// initializeDefaultAlertRules sets up default alert rules
func (sm *SecurityMonitor) initializeDefaultAlertRules() {
	// High number of failed authentication attempts
	sm.alertRules = append(sm.alertRules, AlertRule{
		ID:          "high_failed_auth",
		Name:        "High Failed Authentication Attempts",
		Description: "Multiple failed authentication attempts detected",
		Severity:    "high",
		Enabled:     true,
		Cooldown:    5 * time.Minute,
		Condition: func(event *SecurityEvent) bool {
			return event.Type == "authentication_failed" && 
				   sm.metrics.FailedAuthentications > 10
		},
		Action: func(event *SecurityEvent) error {
			// Could implement automatic IP blocking here
			return nil
		},
		Tags: []string{"authentication", "brute_force"},
	})

	// Critical attack detected
	sm.alertRules = append(sm.alertRules, AlertRule{
		ID:          "critical_attack",
		Name:        "Critical Attack Detected",
		Description: "A critical security attack has been detected",
		Severity:    "critical",
		Enabled:     true,
		Cooldown:    time.Minute,
		Condition: func(event *SecurityEvent) bool {
			return event.Type == "attack_attempt" && event.Severity == "critical"
		},
		Action: func(event *SecurityEvent) error {
			// Could implement emergency response here
			return nil
		},
		Tags: []string{"attack", "critical"},
	})

	// Data breach attempt
	sm.alertRules = append(sm.alertRules, AlertRule{
		ID:          "data_breach_attempt",
		Name:        "Data Breach Attempt",
		Description: "Unauthorized attempt to access sensitive data",
		Severity:    "critical",
		Enabled:     true,
		Cooldown:    time.Minute,
		Condition: func(event *SecurityEvent) bool {
			return event.Type == "data_breach_attempt"
		},
		Action: func(event *SecurityEvent) error {
			// Could implement data protection measures here
			return nil
		},
		Tags: []string{"data_breach", "unauthorized_access"},
	})
}

// initializeDefaultThresholds sets up default security thresholds
func (sm *SecurityMonitor) initializeDefaultThresholds() {
	sm.thresholds["failed_auth_rate"] = 0.1    // 10% failed authentication rate
	sm.thresholds["attack_rate"] = 0.05        // 5% attack rate
	sm.thresholds["block_rate"] = 0.15         // 15% block rate
	sm.thresholds["response_time"] = 1000.0    // 1 second response time
	sm.thresholds["anomaly_score"] = 0.8       // 80% anomaly confidence
}

// AddEntry adds an entry to the audit trail
func (at *AuditTrail) AddEntry(event *SecurityEvent) error {
	if !at.enabled {
		return nil
	}

	at.mutex.Lock()
	defer at.mutex.Unlock()

	entry := AuditEntry{
		ID:        fmt.Sprintf("audit-%d", time.Now().UnixNano()),
		Timestamp: event.Timestamp,
		Actor:     event.User,
		Action:    event.Action,
		Resource:  event.Resource,
		Outcome:   event.Outcome,
		Details:   event.Details,
		IPAddress: event.IP,
		UserAgent: event.UserAgent,
	}

	// Calculate hash for integrity
	entryData, _ := json.Marshal(entry)
	entry.Hash = fmt.Sprintf("%x", entryData) // Simplified hash

	at.entries = append(at.entries, entry)

	// Apply retention policy
	at.applyRetentionPolicy()

	return nil
}

// applyRetentionPolicy removes old entries based on retention time
func (at *AuditTrail) applyRetentionPolicy() {
	cutoff := time.Now().Add(-at.retention)
	var validEntries []AuditEntry

	for _, entry := range at.entries {
		if entry.Timestamp.After(cutoff) {
			validEntries = append(validEntries, entry)
		}
	}

	at.entries = validEntries
}

// ProcessEvent processes an event for anomaly detection
func (rad *RealTimeAnomalyDetector) ProcessEvent(event *SecurityEvent) error {
	if !rad.enabled {
		return nil
	}

	rad.mutex.Lock()
	defer rad.mutex.Unlock()

	// Extract features from the event
	features := rad.extractFeatures(event)

	// Check each model for anomalies
	for _, model := range rad.models {
		if model.Trained {
			score := rad.calculateAnomalyScore(features, model)
			if score > rad.threshold {
				anomaly := Anomaly{
					Type:      "behavioral",
					Value:     score,
					Expected:  rad.threshold,
					Deviation: score - rad.threshold,
					Severity:  rad.determineSeverity(score),
					Timestamp: event.Timestamp,
					Context:   fmt.Sprintf("Model: %s, Event: %s", model.Name, event.Type),
					Resolved:  false,
				}

				rad.anomalies = append(rad.anomalies, anomaly)
			}
		}
	}

	return nil
}

// extractFeatures extracts numerical features from a security event
func (rad *RealTimeAnomalyDetector) extractFeatures(event *SecurityEvent) map[string]float64 {
	features := make(map[string]float64)

	features["risk_score"] = event.RiskScore
	features["hour_of_day"] = float64(event.Timestamp.Hour())
	features["day_of_week"] = float64(event.Timestamp.Weekday())

	// Extract features from event details
	if val, ok := event.Details["request_size"]; ok {
		if size, ok := val.(float64); ok {
			features["request_size"] = size
		}
	}

	return features
}

// calculateAnomalyScore calculates an anomaly score for given features
func (rad *RealTimeAnomalyDetector) calculateAnomalyScore(features map[string]float64, model *AnomalyModel) float64 {
	// Simplified anomaly scoring - in production would use actual ML models
	score := 0.0
	count := 0

	for feature, value := range features {
		if baseline, exists := model.Parameters[feature+"_mean"]; exists {
			if stddev, exists := model.Parameters[feature+"_stddev"]; exists && stddev > 0 {
				zscore := (value - baseline) / stddev
				score += zscore * zscore // Squared Z-score
				count++
			}
		}
	}

	if count > 0 {
		return score / float64(count)
	}

	return 0.0
}

// determineSeverity determines the severity of an anomaly based on its score
func (rad *RealTimeAnomalyDetector) determineSeverity(score float64) string {
	if score > 2.0 {
		return "critical"
	} else if score > 1.5 {
		return "high"
	} else if score > 1.0 {
		return "medium"
	}
	return "low"
}

// GetCurrentStatus returns the current security monitoring status
func (sm *SecurityMonitor) GetCurrentStatus() map[string]interface{} {
	status := make(map[string]interface{})

	status["enabled"] = sm.enabled
	status["metrics"] = sm.GetSecurityMetrics()
	status["alert_rules_count"] = len(sm.alertRules)
	status["events_collected"] = len(sm.eventCollector.events)
	status["audit_entries"] = len(sm.auditTrail.entries)
	status["anomalies_detected"] = len(sm.anomalyDetector.anomalies)

	// Count active alerts
	activeAlerts := 0
	select {
	case <-time.After(time.Millisecond):
		// Non-blocking check of alert queue
		activeAlerts = len(sm.alertQueue)
	default:
	}
	status["active_alerts"] = activeAlerts

	return status
}