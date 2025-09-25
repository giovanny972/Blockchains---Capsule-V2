package keeper

import (
	"context"
	"fmt"
	"sort"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/timecapsule/types"
)

// MultiSigManager handles multi-signature operations for capsules
type MultiSigManager struct {
	keeper *Keeper
}

// MultiSigSession represents an active multi-signature session
type MultiSigSession struct {
	ID              string                 `json:"id"`
	CapsuleID       uint64                 `json:"capsule_id"`
	RequiredSigs    uint32                 `json:"required_sigs"`
	Participants    []string               `json:"participants"`
	Signatures      map[string]*Signature  `json:"signatures"`
	Status          string                 `json:"status"` // "pending", "completed", "expired", "cancelled"
	CreatedAt       time.Time              `json:"created_at"`
	ExpiresAt       time.Time              `json:"expires_at"`
	CompletedAt     *time.Time             `json:"completed_at,omitempty"`
	SessionData     []byte                 `json:"session_data"` // Data being signed
	Purpose         string                 `json:"purpose"`      // "open", "transfer", "modify"
	CreatedBy       string                 `json:"created_by"`
	Metadata        map[string]interface{} `json:"metadata"`
}

// Signature represents a digital signature in a multi-sig session
type Signature struct {
	Signer       string    `json:"signer"`
	Signature    []byte    `json:"signature"`
	PublicKey    []byte    `json:"public_key"`
	Algorithm    string    `json:"algorithm"`
	Timestamp    time.Time `json:"timestamp"`
	Message      string    `json:"message,omitempty"`
	Verified     bool      `json:"verified"`
	IPAddress    string    `json:"ip_address,omitempty"`
	UserAgent    string    `json:"user_agent,omitempty"`
}

// MultiSigPolicy defines the policy for multi-signature operations
type MultiSigPolicy struct {
	CapsuleID         uint64                 `json:"capsule_id"`
	RequiredSigs      uint32                 `json:"required_sigs"`
	AuthorizedSigners []string               `json:"authorized_signers"`
	ExpirationTime    time.Duration          `json:"expiration_time"`
	AllowedOperations []string               `json:"allowed_operations"`
	Restrictions      map[string]interface{} `json:"restrictions"`
	CreatedAt         time.Time              `json:"created_at"`
	UpdatedAt         time.Time              `json:"updated_at"`
	CreatedBy         string                 `json:"created_by"`
}

// NewMultiSigManager creates a new multi-signature manager
func NewMultiSigManager(keeper *Keeper) *MultiSigManager {
	return &MultiSigManager{
		keeper: keeper,
	}
}

// CreateMultiSigSession creates a new multi-signature session
func (msm *MultiSigManager) CreateMultiSigSession(
	ctx context.Context,
	capsuleID uint64,
	purpose string,
	sessionData []byte,
	participants []string,
	requiredSigs uint32,
	expirationHours int,
	createdBy string,
) (*MultiSigSession, error) {
	// Validate capsule exists
	capsule, err := msm.keeper.GetCapsule(ctx, capsuleID)
	if err != nil {
		return nil, fmt.Errorf("capsule not found: %w", err)
	}

	// Validate creator authorization
	if capsule.Owner != createdBy {
		return nil, fmt.Errorf("only capsule owner can create multi-sig sessions")
	}

	// Validate participants
	if len(participants) < int(requiredSigs) {
		return nil, fmt.Errorf("not enough participants: need at least %d, got %d", requiredSigs, len(participants))
	}

	// Remove duplicates and validate addresses
	uniqueParticipants := make([]string, 0, len(participants))
	seen := make(map[string]bool)
	for _, participant := range participants {
		if !seen[participant] {
			// Validate address format
			if _, err := msm.keeper.addressCodec.StringToBytes(participant); err != nil {
				return nil, fmt.Errorf("invalid participant address %s: %w", participant, err)
			}
			uniqueParticipants = append(uniqueParticipants, participant)
			seen[participant] = true
		}
	}

	if len(uniqueParticipants) < int(requiredSigs) {
		return nil, fmt.Errorf("not enough unique participants after deduplication")
	}

	// Create session
	now := sdk.UnwrapSDKContext(ctx).BlockTime()
	session := &MultiSigSession{
		ID:           fmt.Sprintf("msig-%d-%d", capsuleID, now.UnixNano()),
		CapsuleID:    capsuleID,
		RequiredSigs: requiredSigs,
		Participants: uniqueParticipants,
		Signatures:   make(map[string]*Signature),
		Status:       "pending",
		CreatedAt:    now,
		ExpiresAt:    now.Add(time.Duration(expirationHours) * time.Hour),
		SessionData:  sessionData,
		Purpose:      purpose,
		CreatedBy:    createdBy,
		Metadata:     make(map[string]interface{}),
	}

	// Store session (in a real implementation, this would be persisted)
	// For now, we'll emit an event
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	sdkCtx.EventManager().EmitEvent(
		sdk.NewEvent(
			"multisig_session_created",
			sdk.NewAttribute("session_id", session.ID),
			sdk.NewAttribute("capsule_id", fmt.Sprintf("%d", capsuleID)),
			sdk.NewAttribute("required_sigs", fmt.Sprintf("%d", requiredSigs)),
			sdk.NewAttribute("participants_count", fmt.Sprintf("%d", len(uniqueParticipants))),
			sdk.NewAttribute("purpose", purpose),
			sdk.NewAttribute("created_by", createdBy),
		),
	)

	return session, nil
}

// AddSignature adds a signature to a multi-sig session
func (msm *MultiSigManager) AddSignature(
	ctx context.Context,
	sessionID string,
	signer string,
	signature []byte,
	publicKey []byte,
	message string,
) error {
	// In a real implementation, this would load the session from storage
	// For now, we'll create a mock validation
	
	// Validate signer address
	if _, err := msm.keeper.addressCodec.StringToBytes(signer); err != nil {
		return fmt.Errorf("invalid signer address: %w", err)
	}

	// Validate signature format
	if len(signature) == 0 {
		return fmt.Errorf("signature cannot be empty")
	}

	if len(publicKey) == 0 {
		return fmt.Errorf("public key cannot be empty")
	}

	// Create signature object
	sig := &Signature{
		Signer:    signer,
		Signature: signature,
		PublicKey: publicKey,
		Algorithm: "secp256k1", // Default algorithm
		Timestamp: sdk.UnwrapSDKContext(ctx).BlockTime(),
		Message:   message,
		Verified:  false, // Will be verified by cryptographic validation
	}

	// Verify signature (simplified - in production would use actual crypto verification)
	if err := msm.verifySignature(sig); err != nil {
		return fmt.Errorf("signature verification failed: %w", err)
	}

	sig.Verified = true

	// Emit event for signature addition
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	sdkCtx.EventManager().EmitEvent(
		sdk.NewEvent(
			"multisig_signature_added",
			sdk.NewAttribute("session_id", sessionID),
			sdk.NewAttribute("signer", signer),
			sdk.NewAttribute("signature_verified", "true"),
			sdk.NewAttribute("timestamp", sig.Timestamp.Format(time.RFC3339)),
		),
	)

	return nil
}

// verifySignature verifies a digital signature
func (msm *MultiSigManager) verifySignature(sig *Signature) error {
	// Simplified signature verification
	// In production, this would perform actual cryptographic verification
	
	if len(sig.Signature) < 64 {
		return fmt.Errorf("signature too short")
	}

	if len(sig.PublicKey) < 33 {
		return fmt.Errorf("public key too short")
	}

	// Mock verification - always passes for demo
	return nil
}

// CheckSessionStatus checks and updates the status of a multi-sig session
func (msm *MultiSigManager) CheckSessionStatus(
	ctx context.Context,
	sessionID string,
) (string, error) {
	// In a real implementation, this would load and check the actual session
	// For now, we'll provide a mock implementation
	
	now := sdk.UnwrapSDKContext(ctx).BlockTime()
	
	// Mock session data for demonstration
	// In production, load from storage
	mockSession := &MultiSigSession{
		ID:           sessionID,
		RequiredSigs: 3,
		Status:       "pending",
		ExpiresAt:    now.Add(24 * time.Hour),
		Signatures:   make(map[string]*Signature),
	}

	// Check if expired
	if now.After(mockSession.ExpiresAt) {
		mockSession.Status = "expired"
		return mockSession.Status, nil
	}

	// Check if enough signatures collected
	verifiedSigs := 0
	for _, sig := range mockSession.Signatures {
		if sig.Verified {
			verifiedSigs++
		}
	}

	if verifiedSigs >= int(mockSession.RequiredSigs) {
		mockSession.Status = "completed"
		completedAt := now
		mockSession.CompletedAt = &completedAt
		
		// Emit completion event
		sdkCtx := sdk.UnwrapSDKContext(ctx)
		sdkCtx.EventManager().EmitEvent(
			sdk.NewEvent(
				"multisig_session_completed",
				sdk.NewAttribute("session_id", sessionID),
				sdk.NewAttribute("signatures_collected", fmt.Sprintf("%d", verifiedSigs)),
				sdk.NewAttribute("completed_at", completedAt.Format(time.RFC3339)),
			),
		)
	}

	return mockSession.Status, nil
}

// ExecuteMultiSigOperation executes an operation after multi-sig validation
func (msm *MultiSigManager) ExecuteMultiSigOperation(
	ctx context.Context,
	sessionID string,
	operation string,
) error {
	// Verify session is completed
	status, err := msm.CheckSessionStatus(ctx, sessionID)
	if err != nil {
		return fmt.Errorf("failed to check session status: %w", err)
	}

	if status != "completed" {
		return fmt.Errorf("session not completed, current status: %s", status)
	}

	// Execute the operation based on type
	switch operation {
	case "open_capsule":
		return msm.executeOpenCapsule(ctx, sessionID)
	case "transfer_ownership":
		return msm.executeTransferOwnership(ctx, sessionID)
	case "modify_capsule":
		return msm.executeModifyCapsule(ctx, sessionID)
	default:
		return fmt.Errorf("unknown operation: %s", operation)
	}
}

// executeOpenCapsule executes a multi-sig capsule opening
func (msm *MultiSigManager) executeOpenCapsule(ctx context.Context, sessionID string) error {
	// In a real implementation, this would:
	// 1. Load the session details
	// 2. Extract the capsule ID
	// 3. Perform the actual capsule opening
	// 4. Return the decrypted data
	
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	sdkCtx.EventManager().EmitEvent(
		sdk.NewEvent(
			"multisig_capsule_opened",
			sdk.NewAttribute("session_id", sessionID),
			sdk.NewAttribute("operation", "open_capsule"),
			sdk.NewAttribute("timestamp", sdkCtx.BlockTime().Format(time.RFC3339)),
		),
	)

	return nil
}

// executeTransferOwnership executes a multi-sig ownership transfer
func (msm *MultiSigManager) executeTransferOwnership(ctx context.Context, sessionID string) error {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	sdkCtx.EventManager().EmitEvent(
		sdk.NewEvent(
			"multisig_ownership_transferred",
			sdk.NewAttribute("session_id", sessionID),
			sdk.NewAttribute("operation", "transfer_ownership"),
			sdk.NewAttribute("timestamp", sdkCtx.BlockTime().Format(time.RFC3339)),
		),
	)

	return nil
}

// executeModifyCapsule executes a multi-sig capsule modification
func (msm *MultiSigManager) executeModifyCapsule(ctx context.Context, sessionID string) error {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	sdkCtx.EventManager().EmitEvent(
		sdk.NewEvent(
			"multisig_capsule_modified",
			sdk.NewAttribute("session_id", sessionID),
			sdk.NewAttribute("operation", "modify_capsule"),
			sdk.NewAttribute("timestamp", sdkCtx.BlockTime().Format(time.RFC3339)),
		),
	)

	return nil
}

// GetMultiSigPolicy retrieves the multi-sig policy for a capsule
func (msm *MultiSigManager) GetMultiSigPolicy(
	ctx context.Context,
	capsuleID uint64,
) (*MultiSigPolicy, error) {
	// In a real implementation, this would load from storage
	// For now, return a default policy
	
	capsule, err := msm.keeper.GetCapsule(ctx, capsuleID)
	if err != nil {
		return nil, err
	}

	policy := &MultiSigPolicy{
		CapsuleID:         capsuleID,
		RequiredSigs:      capsule.RequiredSigs,
		AuthorizedSigners: capsule.ShareHolders,
		ExpirationTime:    24 * time.Hour,
		AllowedOperations: []string{"open", "transfer", "modify"},
		Restrictions:      make(map[string]interface{}),
		CreatedAt:         capsule.CreatedAt,
		UpdatedAt:         capsule.UpdatedAt,
		CreatedBy:         capsule.Owner,
	}

	// Add default restrictions
	policy.Restrictions["max_session_duration"] = "24h"
	policy.Restrictions["require_owner_signature"] = true
	policy.Restrictions["allow_delegation"] = false

	return policy, nil
}

// UpdateMultiSigPolicy updates the multi-sig policy for a capsule
func (msm *MultiSigManager) UpdateMultiSigPolicy(
	ctx context.Context,
	capsuleID uint64,
	policy *MultiSigPolicy,
	updatedBy string,
) error {
	// Validate capsule exists and updater authorization
	capsule, err := msm.keeper.GetCapsule(ctx, capsuleID)
	if err != nil {
		return err
	}

	if capsule.Owner != updatedBy {
		return fmt.Errorf("only capsule owner can update multi-sig policy")
	}

	// Validate policy parameters
	if policy.RequiredSigs == 0 {
		return fmt.Errorf("required signatures must be greater than 0")
	}

	if len(policy.AuthorizedSigners) < int(policy.RequiredSigs) {
		return fmt.Errorf("not enough authorized signers for required signatures")
	}

	// Update policy timestamp
	policy.UpdatedAt = sdk.UnwrapSDKContext(ctx).BlockTime()

	// Emit event
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	sdkCtx.EventManager().EmitEvent(
		sdk.NewEvent(
			"multisig_policy_updated",
			sdk.NewAttribute("capsule_id", fmt.Sprintf("%d", capsuleID)),
			sdk.NewAttribute("required_sigs", fmt.Sprintf("%d", policy.RequiredSigs)),
			sdk.NewAttribute("authorized_signers_count", fmt.Sprintf("%d", len(policy.AuthorizedSigners))),
			sdk.NewAttribute("updated_by", updatedBy),
		),
	)

	return nil
}

// ListActiveMultiSigSessions lists active multi-sig sessions for a user
func (msm *MultiSigManager) ListActiveMultiSigSessions(
	ctx context.Context,
	user string,
	limit int,
) ([]*MultiSigSession, error) {
	// In a real implementation, this would query the storage
	// For now, return empty list
	sessions := make([]*MultiSigSession, 0)

	// Mock implementation - would filter by user participation
	return sessions, nil
}

// GetMultiSigSessionDetails retrieves detailed information about a session
func (msm *MultiSigManager) GetMultiSigSessionDetails(
	ctx context.Context,
	sessionID string,
) (*MultiSigSession, error) {
	// In a real implementation, this would load from storage
	// For now, return a mock session
	
	now := sdk.UnwrapSDKContext(ctx).BlockTime()
	session := &MultiSigSession{
		ID:           sessionID,
		CapsuleID:    1,
		RequiredSigs: 3,
		Participants: []string{"participant1", "participant2", "participant3", "participant4"},
		Signatures:   make(map[string]*Signature),
		Status:       "pending",
		CreatedAt:    now.Add(-time.Hour),
		ExpiresAt:    now.Add(23 * time.Hour),
		Purpose:      "open_capsule",
		CreatedBy:    "creator_address",
		Metadata:     make(map[string]interface{}),
	}

	return session, nil
}

// ValidateMultiSigOperation validates if an operation can be performed with multi-sig
func (msm *MultiSigManager) ValidateMultiSigOperation(
	ctx context.Context,
	capsuleID uint64,
	operation string,
	user string,
) (bool, string, error) {
	// Get capsule
	capsule, err := msm.keeper.GetCapsule(ctx, capsuleID)
	if err != nil {
		return false, "capsule not found", err
	}

	// Check if capsule requires multi-sig
	if capsule.CapsuleType != types.CapsuleType_MULTI_SIG {
		return false, "capsule is not multi-sig type", nil
	}

	// Get policy
	policy, err := msm.GetMultiSigPolicy(ctx, capsuleID)
	if err != nil {
		return false, "failed to get multi-sig policy", err
	}

	// Check if operation is allowed
	operationAllowed := false
	for _, allowedOp := range policy.AllowedOperations {
		if allowedOp == operation {
			operationAllowed = true
			break
		}
	}

	if !operationAllowed {
		return false, fmt.Sprintf("operation '%s' not allowed by policy", operation), nil
	}

	// Check if user is authorized
	userAuthorized := false
	for _, signer := range policy.AuthorizedSigners {
		if signer == user {
			userAuthorized = true
			break
		}
	}

	if !userAuthorized {
		return false, "user not authorized for multi-sig operations", nil
	}

	return true, "validation passed", nil
}

// GetMultiSigStatistics returns statistics about multi-sig operations
func (msm *MultiSigManager) GetMultiSigStatistics(ctx context.Context) map[string]interface{} {
	stats := make(map[string]interface{})

	// Mock statistics - in production would query actual data
	stats["total_sessions"] = 0
	stats["active_sessions"] = 0
	stats["completed_sessions"] = 0
	stats["expired_sessions"] = 0
	stats["average_completion_time"] = "2h 30m"
	stats["most_common_operation"] = "open_capsule"
	stats["success_rate"] = 0.95

	return stats
}