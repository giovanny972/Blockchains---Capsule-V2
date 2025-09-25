package types

import (
	"cosmossdk.io/errors"
)

// DONTCOVER

// x/timecapsule module sentinel errors
var (
	ErrInvalidCapsule        = errors.Register(ModuleName, 2, "invalid capsule")
	ErrCapsuleNotFound       = errors.Register(ModuleName, 3, "capsule not found")
	ErrUnauthorized          = errors.Register(ModuleName, 4, "unauthorized access")
	ErrCapsuleLocked         = errors.Register(ModuleName, 5, "capsule is locked")
	ErrCapsuleExpired        = errors.Register(ModuleName, 6, "capsule has expired")
	ErrInvalidKeyShare       = errors.Register(ModuleName, 7, "invalid key share")
	ErrInsufficientShares    = errors.Register(ModuleName, 8, "insufficient key shares")
	ErrInvalidEncryption     = errors.Register(ModuleName, 9, "invalid encryption")
	ErrConditionNotMet       = errors.Register(ModuleName, 10, "access condition not met")
	ErrInvalidTimelock       = errors.Register(ModuleName, 11, "invalid timelock")
	ErrCapsuleAlreadyOpened  = errors.Register(ModuleName, 12, "capsule already opened")
	ErrInvalidRecipient      = errors.Register(ModuleName, 13, "invalid recipient")
	ErrInvalidThreshold      = errors.Register(ModuleName, 14, "invalid threshold")
	ErrNodeNotFound          = errors.Register(ModuleName, 15, "masternode not found")
	ErrKeyShareExists        = errors.Register(ModuleName, 16, "key share already exists")
	ErrInvalidSignature      = errors.Register(ModuleName, 17, "invalid signature")
	ErrContractExecution     = errors.Register(ModuleName, 18, "contract execution failed")
	ErrInvalidCapsuleType    = errors.Register(ModuleName, 19, "invalid capsule type")
	ErrDataTooLarge          = errors.Register(ModuleName, 20, "data too large")
	ErrInvalidMetadata       = errors.Register(ModuleName, 21, "invalid metadata")
	ErrInvalidTransfer       = errors.Register(ModuleName, 22, "invalid transfer")
	ErrTransferExpired       = errors.Register(ModuleName, 23, "transfer has expired")
	ErrTransferNotFound      = errors.Register(ModuleName, 24, "transfer not found")
	ErrDataStorageFailed     = errors.Register(ModuleName, 25, "data storage failed")
	ErrDataRetrievalFailed   = errors.Register(ModuleName, 26, "data retrieval failed")
	ErrInvalidAddress        = errors.Register(ModuleName, 27, "invalid address")
	ErrInvalidRequest        = errors.Register(ModuleName, 28, "invalid request")
	ErrInvalidCoins          = errors.Register(ModuleName, 29, "invalid coins")
)