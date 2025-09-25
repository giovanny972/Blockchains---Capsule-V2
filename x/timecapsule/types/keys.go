package types

import (
	"cosmossdk.io/collections"
)

const (
	// ModuleName defines the module name
	ModuleName = "timecapsule"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "mem_timecapsule"

	// RouterKey defines the module's message routing key
	RouterKey = ModuleName
)

// KVStore keys
var (
	// CapsuleKeyPrefix is the prefix for capsule storage keys
	CapsuleKeyPrefix = collections.NewPrefix(0)
	
	// UserCapsulesKeyPrefix is the prefix for user capsules index
	UserCapsulesKeyPrefix = collections.NewPrefix(1)
	
	// ConditionContractsKeyPrefix is the prefix for condition contracts
	ConditionContractsKeyPrefix = collections.NewPrefix(9)
	
	// KeySharesKeyPrefix is the prefix for key shares storage
	KeySharesKeyPrefix = collections.NewPrefix(3)
	
	// CapsuleCounterKey is the key for the global capsule counter
	CapsuleCounterKey = collections.NewPrefix(4)
	
	// NodeKeysKeyPrefix is the prefix for masternode keys
	NodeKeysKeyPrefix = collections.NewPrefix(5)
	
	// TransferHistoryKeyPrefix is the prefix for transfer history storage
	TransferHistoryKeyPrefix = collections.NewPrefix(6)
	
	// PendingTransfersKeyPrefix is the prefix for pending transfers storage
	PendingTransfersKeyPrefix = collections.NewPrefix(7)
	
	// TransferStatsKey is the key for transfer statistics
	TransferStatsKey = collections.NewPrefix(8)
	
	// EmergencyActionsKeyPrefix is the prefix for emergency actions storage
	EmergencyActionsKeyPrefix = collections.NewPrefix(10)
)

// Event types
const (
	EventTypeCapsuleCreated = "capsule_created"
	EventTypeCapsuleOpened  = "capsule_opened"
	EventTypeCapsuleUpdated = "capsule_updated"
	EventTypeKeyShareDistributed = "key_share_distributed"
	EventTypeEmergencyContractDeleted = "emergency_contract_deleted"
)

// Event attributes
const (
	AttributeKeyCapsuleID    = "capsule_id"
	AttributeKeyOwner        = "owner"
	AttributeKeyRecipient    = "recipient"
	AttributeKeyCapsuleType  = "capsule_type"
	AttributeKeyUnlockTime   = "unlock_time"
	AttributeKeyDataHash     = "data_hash"
	AttributeKeyNodeID       = "node_id"
	AttributeKeyShareIndex   = "share_index"
	AttributeKeyEmergencyAction = "emergency_action"
	AttributeKeyDeletionID   = "deletion_id"
	AttributeKeyEmergencyReason = "emergency_reason"
)