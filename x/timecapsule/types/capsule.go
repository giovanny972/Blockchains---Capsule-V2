package types

import (
	"fmt"
	"time"
	
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// CapsuleType defines the type of time capsule
type CapsuleType int32

const (
	CapsuleType_UNKNOWN         CapsuleType = 0
	CapsuleType_SAFE           CapsuleType = 1 // Safe storage capsule
	CapsuleType_TIME_LOCK      CapsuleType = 2 // Time-locked capsule  
	CapsuleType_CONDITIONAL    CapsuleType = 3 // Condition-based capsule
	CapsuleType_MULTI_SIG      CapsuleType = 4 // Multi-signature capsule
	CapsuleType_DEAD_MANS_SWITCH CapsuleType = 5 // Dead man's switch capsule
)

// String returns the string representation of CapsuleType
func (t CapsuleType) String() string {
	switch t {
	case CapsuleType_UNKNOWN:
		return "UNKNOWN"
	case CapsuleType_SAFE:
		return "SAFE"
	case CapsuleType_TIME_LOCK:
		return "TIME_LOCK"
	case CapsuleType_CONDITIONAL:
		return "CONDITIONAL"
	case CapsuleType_MULTI_SIG:
		return "MULTI_SIG"
	case CapsuleType_DEAD_MANS_SWITCH:
		return "DEAD_MANS_SWITCH"
	default:
		return "UNKNOWN"
	}
}

// CapsuleStatus defines the status of a capsule
type CapsuleStatus int32

const (
	CapsuleStatus_UNKNOWN    CapsuleStatus = 0
	CapsuleStatus_ACTIVE     CapsuleStatus = 1 // Capsule is active and locked
	CapsuleStatus_UNLOCKED   CapsuleStatus = 2 // Capsule has been unlocked
	CapsuleStatus_EXPIRED    CapsuleStatus = 3 // Capsule has expired
	CapsuleStatus_CANCELLED  CapsuleStatus = 4 // Capsule has been cancelled
)

// String returns the string representation of CapsuleStatus
func (s CapsuleStatus) String() string {
	switch s {
	case CapsuleStatus_UNKNOWN:
		return "UNKNOWN"
	case CapsuleStatus_ACTIVE:
		return "ACTIVE"
	case CapsuleStatus_UNLOCKED:
		return "UNLOCKED"
	case CapsuleStatus_EXPIRED:
		return "EXPIRED"
	case CapsuleStatus_CANCELLED:
		return "CANCELLED"
	default:
		return "UNKNOWN"
	}
}

// TimeCapsule represents a secure data container with conditional access
type TimeCapsule struct {
	ID            uint64        `json:"id"`
	Owner         string        `json:"owner"`
	Recipient     string        `json:"recipient,omitempty"`
	CapsuleType   CapsuleType   `json:"capsule_type"`
	Status        CapsuleStatus `json:"status"`
	
	// Data and encryption
	EncryptedData   []byte `json:"encrypted_data,omitempty"`  // For small data < 1MB
	DataHash        string `json:"data_hash"`                 // SHA-256 hash of original data
	EncryptionAlgo  string `json:"encryption_algo"`           // e.g., "AES-256-GCM"
	
	// IPFS Storage (for large data)
	IPFSHash        string `json:"ipfs_hash,omitempty"`       // IPFS hash for large data
	DataSize        int64  `json:"data_size"`                 // Original data size in bytes
	StorageType     string `json:"storage_type"`              // "blockchain" or "ipfs"
	
	// Access conditions
	UnlockTime      *time.Time `json:"unlock_time,omitempty"`      // For time-locked capsules
	ConditionContract string   `json:"condition_contract,omitempty"` // Smart contract address
	RequiredSigs    uint32     `json:"required_sigs,omitempty"`    // For multi-sig capsules
	
	// Key management (Shamir's Secret Sharing)
	Threshold       uint32   `json:"threshold"`         // Minimum shares needed
	TotalShares     uint32   `json:"total_shares"`      // Total shares created
	ShareHolders    []string `json:"share_holders"`     // Addresses holding shares
	
	// Metadata
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
	ExpiresAt       *time.Time `json:"expires_at,omitempty"`
	
	// Dead man's switch specific
	LastActivity    *time.Time `json:"last_activity,omitempty"`
	InactivityPeriod uint64    `json:"inactivity_period,omitempty"` // In seconds
	
	// Additional metadata
	Title           string            `json:"title,omitempty"`
	Description     string            `json:"description,omitempty"`
	Tags            []string          `json:"tags,omitempty"`
	Metadata        map[string]string `json:"metadata,omitempty"`
}

// KeyShare represents a Shamir secret share
type KeyShare struct {
	CapsuleID   uint64 `json:"capsule_id"`
	ShareIndex  uint32 `json:"share_index"`
	NodeID      string `json:"node_id"`      // Masternode holding the share
	EncryptedShare []byte `json:"encrypted_share"` // The actual encrypted share
	CreatedAt   time.Time `json:"created_at"`
}

// ConditionContract represents a smart contract that defines access conditions
type ConditionContract struct {
	Address     string            `json:"address"`
	Type        string            `json:"type"`        // e.g., "time", "oracle", "multisig"
	Parameters  map[string]string `json:"parameters"`
	CreatedBy   string            `json:"created_by"`
	CreatedAt   time.Time         `json:"created_at"`
}

// CapsuleAccess represents an access attempt or successful access
type CapsuleAccess struct {
	CapsuleID   uint64    `json:"capsule_id"`
	Accessor    string    `json:"accessor"`
	AccessTime  time.Time `json:"access_time"`
	Success     bool      `json:"success"`
	Reason      string    `json:"reason,omitempty"`
}

// TransferHistory represents the history of transfers for a capsule
type TransferHistory struct {
	CapsuleID     uint64    `json:"capsule_id"`
	TransferID    string    `json:"transfer_id"`
	FromOwner     string    `json:"from_owner"`
	ToOwner       string    `json:"to_owner"`
	TransferType  string    `json:"transfer_type"`  // "direct", "batch", "approved"
	Status        string    `json:"status"`         // "pending", "completed", "rejected"
	TransferTime  time.Time `json:"transfer_time"`
	ApprovalTime  *time.Time `json:"approval_time,omitempty"`
	Message       string    `json:"message,omitempty"`
	TransferFee   sdk.Coins `json:"transfer_fee,omitempty"`
	BlockHeight   int64     `json:"block_height"`
	TxHash        string    `json:"tx_hash,omitempty"`
}

// PendingTransfer represents a transfer awaiting approval
type PendingTransfer struct {
	TransferID      string    `json:"transfer_id"`
	CapsuleID       uint64    `json:"capsule_id"`
	FromOwner       string    `json:"from_owner"`
	ToOwner         string    `json:"to_owner"`
	RequestTime     time.Time `json:"request_time"`
	ExpiryTime      time.Time `json:"expiry_time"`
	Message         string    `json:"message,omitempty"`
	RequireApproval bool      `json:"require_approval"`
	Status          string    `json:"status"` // "pending", "approved", "rejected", "expired"
}

// TransferStats represents statistics about capsule transfers
type TransferStats struct {
	TotalTransfers     uint64    `json:"total_transfers"`
	PendingTransfers   uint64    `json:"pending_transfers"`
	CompletedTransfers uint64    `json:"completed_transfers"`
	RejectedTransfers  uint64    `json:"rejected_transfers"`
	BatchTransfers     uint64    `json:"batch_transfers"`
	LastTransferTime   *time.Time `json:"last_transfer_time,omitempty"`
	TotalFeesCollected sdk.Coins `json:"total_fees_collected"`
}

// Validate validates the time capsule
func (tc *TimeCapsule) Validate() error {
	if tc.ID == 0 {
		return fmt.Errorf("capsule ID cannot be zero")
	}
	
	if _, err := sdk.AccAddressFromBech32(tc.Owner); err != nil {
		return fmt.Errorf("invalid owner address: %w", err)
	}
	
	if tc.Recipient != "" {
		if _, err := sdk.AccAddressFromBech32(tc.Recipient); err != nil {
			return fmt.Errorf("invalid recipient address: %w", err)
		}
	}
	
	if len(tc.EncryptedData) == 0 {
		return fmt.Errorf("encrypted data cannot be empty")
	}
	
	if tc.DataHash == "" {
		return fmt.Errorf("data hash cannot be empty")
	}
	
	if tc.Threshold > tc.TotalShares {
		return fmt.Errorf("threshold cannot be greater than total shares")
	}
	
	if tc.TotalShares == 0 {
		return fmt.Errorf("total shares must be greater than zero")
	}
	
	// Validate capsule type specific requirements
	switch tc.CapsuleType {
	case CapsuleType_TIME_LOCK:
		if tc.UnlockTime == nil {
			return fmt.Errorf("time-locked capsule must have unlock time")
		}
		if tc.UnlockTime.Before(time.Now()) {
			return fmt.Errorf("unlock time must be in the future")
		}
	case CapsuleType_CONDITIONAL:
		if tc.ConditionContract == "" {
			return fmt.Errorf("conditional capsule must have condition contract")
		}
	case CapsuleType_MULTI_SIG:
		if tc.RequiredSigs == 0 {
			return fmt.Errorf("multi-sig capsule must have required signatures")
		}
	case CapsuleType_DEAD_MANS_SWITCH:
		if tc.InactivityPeriod == 0 {
			return fmt.Errorf("dead man's switch capsule must have inactivity period")
		}
	}
	
	return nil
}

// IsUnlockable checks if the capsule can be unlocked based on current conditions
func (tc *TimeCapsule) IsUnlockable(ctx sdk.Context) bool {
	if tc.Status != CapsuleStatus_ACTIVE {
		return false
	}
	
	switch tc.CapsuleType {
	case CapsuleType_SAFE:
		return true // Always unlockable by owner
		
	case CapsuleType_TIME_LOCK:
		if tc.UnlockTime == nil {
			return false
		}
		return ctx.BlockTime().After(*tc.UnlockTime)
		
	case CapsuleType_DEAD_MANS_SWITCH:
		if tc.LastActivity == nil || tc.InactivityPeriod == 0 {
			return false
		}
		inactivityDuration := time.Duration(tc.InactivityPeriod) * time.Second
		return ctx.BlockTime().After(tc.LastActivity.Add(inactivityDuration))
		
	case CapsuleType_CONDITIONAL:
		// This would require checking the smart contract condition
		// Implementation depends on the specific condition logic
		return false // Placeholder
		
	case CapsuleType_MULTI_SIG:
		// This would require checking if enough signatures have been provided
		// Implementation depends on signature collection mechanism
		return false // Placeholder
	}
	
	return false
}

// UpdateActivity updates the last activity timestamp for dead man's switch capsules
func (tc *TimeCapsule) UpdateActivity(blockTime time.Time) {
	if tc.CapsuleType == CapsuleType_DEAD_MANS_SWITCH {
		tc.LastActivity = &blockTime
		tc.UpdatedAt = blockTime
	}
}

// CapsuleStats represents comprehensive statistics about capsules
type CapsuleStats struct {
	TotalCapsules       uint64            `json:"total_capsules"`
	ActiveCapsules      uint64            `json:"active_capsules"`
	UnlockedCapsules    uint64            `json:"unlocked_capsules"`
	ExpiredCapsules     uint64            `json:"expired_capsules"`
	MyCapsulesCount     uint64            `json:"my_capsules_count"`
	MyActiveCapsules    uint64            `json:"my_active_capsules"`
	TotalDataStored     string            `json:"total_data_stored"`
	AverageUnlockTime   int64             `json:"average_unlock_time"`
	MostUsedType        string            `json:"most_used_type"`
	TypeDistribution    map[string]uint64 `json:"type_distribution"`
	StatusDistribution  map[string]uint64 `json:"status_distribution"`
}

// BatchOpenRequest represents a request to open a capsule in a batch operation
type BatchOpenRequest struct {
	CapsuleID uint64      `json:"capsule_id"`
	Accessor  string      `json:"accessor"`
	KeyShares []*KeyShare `json:"key_shares"`
}

// BatchOpenResponse represents the response for a batch open operation
type BatchOpenResponse struct {
	CapsuleID uint64 `json:"capsule_id"`
	Success   bool   `json:"success"`
	Data      []byte `json:"data,omitempty"`
	Error     string `json:"error,omitempty"`
}

// NetworkHealth represents the health status of the network
type NetworkHealth struct {
	BlockchainStatus   string  `json:"blockchain_status"`
	BlockHeight        uint64  `json:"block_height"`
	AverageBlockTime   float64 `json:"average_block_time"`
	ConnectedNodes     int32   `json:"connected_nodes"`
	IPFSStatus         string  `json:"ipfs_status"`
	IPFSNodes          int32   `json:"ipfs_nodes"`
	TotalTransactions  uint64  `json:"total_transactions"`
	CapsuleCount       uint64  `json:"capsule_count"`
	NetworkLatency     int64   `json:"network_latency"`
}

// CapsuleMetrics represents detailed metrics for monitoring
type CapsuleMetrics struct {
	CreationRate       float64 `json:"creation_rate"`       // Capsules per hour
	OpeningRate        float64 `json:"opening_rate"`        // Openings per hour
	AverageDataSize    float64 `json:"average_data_size"`   // In bytes
	StorageEfficiency  float64 `json:"storage_efficiency"`  // Percentage
	SecurityScore      float64 `json:"security_score"`      // 0-100
	UptimePercentage   float64 `json:"uptime_percentage"`   // 0-100
}

// SmartOpenCondition represents intelligent opening conditions
type SmartOpenCondition struct {
	Type        string                 `json:"type"`
	Parameters  map[string]interface{} `json:"parameters"`
	Description string                 `json:"description"`
	IsActive    bool                   `json:"is_active"`
	Priority    int32                  `json:"priority"`
}

// OptimizedCapsuleView provides a lightweight view for listing
type OptimizedCapsuleView struct {
	ID              uint64    `json:"id"`
	Title           string    `json:"title"`
	Type            string    `json:"type"`
	Status          string    `json:"status"`
	Owner           string    `json:"owner"`
	Recipient       string    `json:"recipient,omitempty"`
	DataSize        int64     `json:"data_size"`
	StorageType     string    `json:"storage_type"`
	UnlockTime      *time.Time `json:"unlock_time,omitempty"`
	CreatedAt       time.Time `json:"created_at"`
	IsUnlockable    bool      `json:"is_unlockable"`
	TimeRemaining   *int64    `json:"time_remaining,omitempty"` // Seconds until unlock
}

// GetOptimizedView returns a lightweight view of the capsule
func (tc *TimeCapsule) GetOptimizedView(currentTime time.Time) *OptimizedCapsuleView {
	view := &OptimizedCapsuleView{
		ID:           tc.ID,
		Title:        tc.Title,
		Type:         tc.CapsuleType.String(),
		Status:       tc.Status.String(),
		Owner:        tc.Owner,
		Recipient:    tc.Recipient,
		DataSize:     tc.DataSize,
		StorageType:  tc.StorageType,
		UnlockTime:   tc.UnlockTime,
		CreatedAt:    tc.CreatedAt,
		IsUnlockable: false,
	}

	// Calculate if unlockable and time remaining
	switch tc.CapsuleType {
	case CapsuleType_SAFE:
		view.IsUnlockable = (tc.Status == CapsuleStatus_ACTIVE)
	case CapsuleType_TIME_LOCK:
		if tc.UnlockTime != nil {
			if currentTime.After(*tc.UnlockTime) {
				view.IsUnlockable = (tc.Status == CapsuleStatus_ACTIVE)
			} else {
				remaining := int64(tc.UnlockTime.Sub(currentTime).Seconds())
				view.TimeRemaining = &remaining
			}
		}
	case CapsuleType_DEAD_MANS_SWITCH:
		if tc.LastActivity != nil && tc.InactivityPeriod > 0 {
			inactivityDuration := time.Duration(tc.InactivityPeriod) * time.Second
			deadlineTime := tc.LastActivity.Add(inactivityDuration)
			if currentTime.After(deadlineTime) {
				view.IsUnlockable = (tc.Status == CapsuleStatus_ACTIVE)
			} else {
				remaining := int64(deadlineTime.Sub(currentTime).Seconds())
				view.TimeRemaining = &remaining
			}
		}
	}

	return view
}

// IsExpiringSoon checks if capsule will unlock/expire within specified hours
func (tc *TimeCapsule) IsExpiringSoon(hours int) bool {
	if tc.Status != CapsuleStatus_ACTIVE {
		return false
	}

	threshold := time.Now().Add(time.Duration(hours) * time.Hour)
	
	switch tc.CapsuleType {
	case CapsuleType_TIME_LOCK:
		return tc.UnlockTime != nil && tc.UnlockTime.Before(threshold)
	case CapsuleType_DEAD_MANS_SWITCH:
		if tc.LastActivity != nil && tc.InactivityPeriod > 0 {
			inactivityDuration := time.Duration(tc.InactivityPeriod) * time.Second
			deadlineTime := tc.LastActivity.Add(inactivityDuration)
			return deadlineTime.Before(threshold)
		}
	}
	
	return false
}

// EmergencyAction represents an emergency action taken on a capsule
type EmergencyAction struct {
	ID               string    `json:"id"`
	CapsuleID        uint64    `json:"capsule_id"`
	Creator          string    `json:"creator"`
	ActionType       string    `json:"action_type"`       // "contract_deletion", "force_unlock", etc.
	Reason           string    `json:"reason"`
	ConfirmationCode string    `json:"confirmation_code"`
	ActionTime       time.Time `json:"action_time"`
	BlockHeight      int64     `json:"block_height"`
	IsReversible     bool      `json:"is_reversible"`
	ReversedAt       *time.Time `json:"reversed_at,omitempty"`
	ReversedBy       string    `json:"reversed_by,omitempty"`
}