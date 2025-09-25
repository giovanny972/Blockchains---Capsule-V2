package keeper

import (
	"context"
	"fmt"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/timecapsule/types"
)

// NotificationService handles intelligent notifications for capsule events
type NotificationService struct {
	keeper *Keeper
}

// NewNotificationService creates a new notification service
func NewNotificationService(keeper *Keeper) *NotificationService {
	return &NotificationService{
		keeper: keeper,
	}
}

// NotificationEvent represents different types of notifications
type NotificationEvent struct {
	Type        string                 `json:"type"`
	CapsuleID   uint64                 `json:"capsule_id"`
	User        string                 `json:"user"`
	Message     string                 `json:"message"`
	Priority    string                 `json:"priority"` // low, medium, high, urgent
	Timestamp   time.Time              `json:"timestamp"`
	Metadata    map[string]interface{} `json:"metadata"`
	ActionItems []string               `json:"action_items,omitempty"`
}

// NotificationRule defines when and how to send notifications
type NotificationRule struct {
	EventType     string        `json:"event_type"`
	AdvanceTime   time.Duration `json:"advance_time"`
	RepeatInterval time.Duration `json:"repeat_interval,omitempty"`
	MaxRepeats    int           `json:"max_repeats"`
	Priority      string        `json:"priority"`
	Template      string        `json:"template"`
}

// Default notification rules
var DefaultNotificationRules = []NotificationRule{
	{
		EventType:   "unlock_soon",
		AdvanceTime: 24 * time.Hour,
		Priority:    "medium",
		Template:    "Your capsule '{title}' will unlock in {time_remaining}",
	},
	{
		EventType:   "unlock_very_soon",
		AdvanceTime: 1 * time.Hour,
		Priority:    "high",
		Template:    "⚠️ Your capsule '{title}' will unlock in {time_remaining}!",
	},
	{
		EventType:   "unlock_imminent",
		AdvanceTime: 10 * time.Minute,
		Priority:    "urgent",
		Template:    "🚨 URGENT: Your capsule '{title}' unlocks in {time_remaining}!",
	},
	{
		EventType:   "dead_mans_switch_warning",
		AdvanceTime: 7 * 24 * time.Hour, // 7 days before
		Priority:    "high",
		Template:    "⚠️ Dead Man's Switch: Your capsule '{title}' will activate for recipient in {time_remaining} unless you update activity",
	},
	{
		EventType:   "capsule_unlocked",
		AdvanceTime: 0,
		Priority:    "high",
		Template:    "🎉 Your capsule '{title}' has been unlocked and is ready for access!",
	},
	{
		EventType:   "transfer_received",
		AdvanceTime: 0,
		Priority:    "medium",
		Template:    "📥 You've received a capsule transfer: '{title}' from {sender}",
	},
}

// ProcessNotifications checks for pending notifications and processes them
func (ns *NotificationService) ProcessNotifications(ctx context.Context) error {
	currentTime := sdk.UnwrapSDKContext(ctx).BlockTime()
	
	// Process all active capsules for notifications
	return ns.keeper.capsules.Walk(ctx, nil, func(id uint64, capsule types.TimeCapsule) (bool, error) {
		if capsule.Status != types.CapsuleStatus_ACTIVE {
			return false, nil // Skip non-active capsules
		}
		
		// Check each notification rule
		for _, rule := range DefaultNotificationRules {
			if ns.shouldNotify(capsule, rule, currentTime) {
				event := ns.createNotificationEvent(capsule, rule, currentTime)
				if err := ns.sendNotification(ctx, event); err != nil {
					// Log error but continue processing
					fmt.Printf("Failed to send notification: %v\n", err)
				}
			}
		}
		
		return false, nil // Continue iteration
	})
}

// shouldNotify determines if a notification should be sent based on the rule
func (ns *NotificationService) shouldNotify(capsule types.TimeCapsule, rule NotificationRule, currentTime time.Time) bool {
	switch rule.EventType {
	case "unlock_soon", "unlock_very_soon", "unlock_imminent":
		if capsule.CapsuleType == types.CapsuleType_TIME_LOCK && capsule.UnlockTime != nil {
			timeUntilUnlock := capsule.UnlockTime.Sub(currentTime)
			return timeUntilUnlock > 0 && timeUntilUnlock <= rule.AdvanceTime
		}
		
	case "dead_mans_switch_warning":
		if capsule.CapsuleType == types.CapsuleType_DEAD_MANS_SWITCH {
			if capsule.LastActivity != nil && capsule.InactivityPeriod > 0 {
				inactivityDuration := time.Duration(capsule.InactivityPeriod) * time.Second
				deadlineTime := capsule.LastActivity.Add(inactivityDuration)
				timeUntilDeadline := deadlineTime.Sub(currentTime)
				return timeUntilDeadline > 0 && timeUntilDeadline <= rule.AdvanceTime
			}
		}
		
	case "capsule_unlocked":
		// This would be triggered by the unlock event, not by time
		return false
	}
	
	return false
}

// createNotificationEvent creates a notification event from a capsule and rule
func (ns *NotificationService) createNotificationEvent(capsule types.TimeCapsule, rule NotificationRule, currentTime time.Time) *NotificationEvent {
	event := &NotificationEvent{
		Type:      rule.EventType,
		CapsuleID: capsule.ID,
		User:      capsule.Owner,
		Priority:  rule.Priority,
		Timestamp: currentTime,
		Metadata:  make(map[string]interface{}),
	}
	
	// Calculate time remaining
	var timeRemaining time.Duration
	switch rule.EventType {
	case "unlock_soon", "unlock_very_soon", "unlock_imminent":
		if capsule.UnlockTime != nil {
			timeRemaining = capsule.UnlockTime.Sub(currentTime)
		}
	case "dead_mans_switch_warning":
		if capsule.LastActivity != nil && capsule.InactivityPeriod > 0 {
			inactivityDuration := time.Duration(capsule.InactivityPeriod) * time.Second
			deadlineTime := capsule.LastActivity.Add(inactivityDuration)
			timeRemaining = deadlineTime.Sub(currentTime)
		}
	}
	
	// Build message from template
	message := rule.Template
	message = ns.replacePlaceholder(message, "{title}", capsule.Title)
	message = ns.replacePlaceholder(message, "{time_remaining}", ns.formatDuration(timeRemaining))
	event.Message = message
	
	// Add metadata
	event.Metadata["capsule_type"] = capsule.CapsuleType.String()
	event.Metadata["time_remaining_seconds"] = int64(timeRemaining.Seconds())
	event.Metadata["storage_type"] = capsule.StorageType
	
	// Add action items based on event type
	switch rule.EventType {
	case "unlock_soon", "unlock_very_soon", "unlock_imminent":
		event.ActionItems = []string{
			"Prepare to access your capsule",
			"Ensure you have the necessary key shares",
			"Check your wallet connection",
		}
	case "dead_mans_switch_warning":
		event.ActionItems = []string{
			"Update your activity to prevent activation",
			"Review recipient settings",
			"Extend inactivity period if needed",
		}
		event.User = capsule.Owner // Notify owner, not recipient
	case "capsule_unlocked":
		event.ActionItems = []string{
			"Access your unlocked capsule",
			"Download your data",
			"Verify data integrity",
		}
	}
	
	return event
}

// sendNotification sends the notification through various channels
func (ns *NotificationService) sendNotification(ctx context.Context, event *NotificationEvent) error {
	// Emit blockchain event
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	sdkCtx.EventManager().EmitEvent(
		sdk.NewEvent(
			"capsule_notification",
			sdk.NewAttribute("type", event.Type),
			sdk.NewAttribute("capsule_id", fmt.Sprintf("%d", event.CapsuleID)),
			sdk.NewAttribute("user", event.User),
			sdk.NewAttribute("priority", event.Priority),
			sdk.NewAttribute("message", event.Message),
		),
	)
	
	// In production, this would also:
	// - Send email notifications
	// - Send push notifications
	// - Store in notification history
	// - Integrate with external notification services
	
	return nil
}

// replacePlaceholder replaces a placeholder in a template string
func (ns *NotificationService) replacePlaceholder(template, placeholder, value string) string {
	// Simple replacement - in production would use a proper template engine
	result := ""
	placeholderLen := len(placeholder)
	i := 0
	
	for i < len(template) {
		if i+placeholderLen <= len(template) && template[i:i+placeholderLen] == placeholder {
			result += value
			i += placeholderLen
		} else {
			result += string(template[i])
			i++
		}
	}
	
	return result
}

// formatDuration formats a duration in a human-readable way
func (ns *NotificationService) formatDuration(d time.Duration) string {
	if d <= 0 {
		return "now"
	}
	
	days := int(d.Hours()) / 24
	hours := int(d.Hours()) % 24
	minutes := int(d.Minutes()) % 60
	
	if days > 0 {
		if hours > 0 {
			return fmt.Sprintf("%d days, %d hours", days, hours)
		}
		return fmt.Sprintf("%d days", days)
	}
	
	if hours > 0 {
		if minutes > 0 {
			return fmt.Sprintf("%d hours, %d minutes", hours, minutes)
		}
		return fmt.Sprintf("%d hours", hours)
	}
	
	if minutes > 0 {
		return fmt.Sprintf("%d minutes", minutes)
	}
	
	seconds := int(d.Seconds())
	return fmt.Sprintf("%d seconds", seconds)
}

// GetUserNotifications retrieves pending notifications for a user
func (ns *NotificationService) GetUserNotifications(ctx context.Context, user string, limit int) ([]*NotificationEvent, error) {
	var notifications []*NotificationEvent
	currentTime := sdk.UnwrapSDKContext(ctx).BlockTime()
	
	err := ns.keeper.capsules.Walk(ctx, nil, func(id uint64, capsule types.TimeCapsule) (bool, error) {
		// Only check capsules owned by or sent to the user
		if capsule.Owner != user && capsule.Recipient != user {
			return false, nil
		}
		
		// Check each notification rule
		for _, rule := range DefaultNotificationRules {
			if ns.shouldNotify(capsule, rule, currentTime) {
				event := ns.createNotificationEvent(capsule, rule, currentTime)
				notifications = append(notifications, event)
				
				if len(notifications) >= limit {
					return true, nil // Stop when limit reached
				}
			}
		}
		
		return false, nil
	})
	
	return notifications, err
}

// NotifyCapsulesExpiringSoon sends notifications for capsules expiring within specified hours
func (ns *NotificationService) NotifyCapsulesExpiringSoon(ctx context.Context, hours int) error {
	expiringSoon, err := ns.keeper.GetExpiringSoonCapsules(ctx, hours)
	if err != nil {
		return err
	}
	
	for _, capsuleView := range expiringSoon {
		// Get full capsule data
		capsule, err := ns.keeper.GetCapsule(ctx, capsuleView.ID)
		if err != nil {
			continue // Skip if can't get capsule
		}
		
		// Create notification event
		event := &NotificationEvent{
			Type:      "unlock_reminder",
			CapsuleID: capsule.ID,
			User:      capsule.Owner,
			Priority:  "medium",
			Timestamp: sdk.UnwrapSDKContext(ctx).BlockTime(),
			Message:   fmt.Sprintf("Your capsule '%s' will unlock soon", capsule.Title),
			Metadata:  map[string]interface{}{
				"time_remaining_seconds": capsuleView.TimeRemaining,
				"capsule_type":          capsuleView.Type,
			},
		}
		
		// Send notification
		if err := ns.sendNotification(ctx, event); err != nil {
			fmt.Printf("Failed to send expiring soon notification for capsule %d: %v\n", capsule.ID, err)
		}
	}
	
	return nil
}

// ScheduleSmartNotifications sets up intelligent notification scheduling
func (ns *NotificationService) ScheduleSmartNotifications(ctx context.Context, capsuleID uint64, customRules []NotificationRule) error {
	capsule, err := ns.keeper.GetCapsule(ctx, capsuleID)
	if err != nil {
		return err
	}
	
	// Validate custom rules
	for i, rule := range customRules {
		if rule.EventType == "" {
			return fmt.Errorf("rule %d: event type cannot be empty", i)
		}
		if rule.Priority == "" {
			customRules[i].Priority = "medium"
		}
		if rule.Template == "" {
			return fmt.Errorf("rule %d: template cannot be empty", i)
		}
	}
	
	// Store custom rules (in production, would persist to database)
	// For now, emit an event that external systems can listen to
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	sdkCtx.EventManager().EmitEvent(
		sdk.NewEvent(
			"smart_notifications_scheduled",
			sdk.NewAttribute("capsule_id", fmt.Sprintf("%d", capsuleID)),
			sdk.NewAttribute("owner", capsule.Owner),
			sdk.NewAttribute("rules_count", fmt.Sprintf("%d", len(customRules))),
		),
	)
	
	return nil
}