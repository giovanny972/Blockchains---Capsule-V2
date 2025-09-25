package keeper

import (
	"context"
	"fmt"

	"cosmossdk.io/collections"
	"cosmossdk.io/log"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/timecapsule/crypto"
	"github.com/cosmos/cosmos-sdk/x/timecapsule/types"
)

// Helper methods for the keeper

// GetParams retrieves the module parameters
func (k Keeper) GetParams(ctx context.Context) (types.Params, error) {
	// For now, return default parameters
	// In a full implementation, this would retrieve from a params store
	return types.DefaultParams(), nil
}

// SetParams sets the module parameters
func (k Keeper) SetParams(ctx context.Context, params types.Params) error {
	// For now, this is a no-op
	// In a full implementation, this would store in a params store
	return nil
}

// GetCapsuleCounter retrieves the current capsule counter
func (k Keeper) GetCapsuleCounter(ctx context.Context) (uint64, error) {
	return k.capsuleCounter.Peek(ctx)
}

// SetCapsuleCounter sets the capsule counter
func (k Keeper) SetCapsuleCounter(ctx context.Context, counter uint64) error {
	return k.capsuleCounter.Set(ctx, counter)
}

// SetCapsule stores a capsule
func (k Keeper) SetCapsule(ctx context.Context, capsule *types.TimeCapsule) error {
	return k.capsules.Set(ctx, capsule.ID, *capsule)
}

// IndexUserCapsule adds a capsule to the user index
func (k Keeper) IndexUserCapsule(ctx context.Context, owner string, capsuleID uint64) error {
	return k.userCapsules.Set(ctx, collections.Join(owner, capsuleID))
}

// SetKeyShare stores a key share
func (k Keeper) SetKeyShare(ctx context.Context, keyShare *types.KeyShare) error {
	key := collections.Join(keyShare.CapsuleID, keyShare.ShareIndex)
	return k.keyShares.Set(ctx, key, *keyShare)
}

// SetConditionContract stores a condition contract
func (k Keeper) SetConditionContract(ctx context.Context, contract *types.ConditionContract) error {
	return k.conditionContracts.Set(ctx, contract.Address, *contract)
}

// GetAllCapsules retrieves all capsules
func (k Keeper) GetAllCapsules(ctx context.Context) ([]types.TimeCapsule, error) {
	var capsules []types.TimeCapsule
	
	err := k.capsules.Walk(ctx, nil, func(key uint64, capsule types.TimeCapsule) (bool, error) {
		capsules = append(capsules, capsule)
		return false, nil // Continue iteration
	})
	
	return capsules, err
}

// GetAllKeyShares retrieves all key shares
func (k Keeper) GetAllKeyShares(ctx context.Context) ([]types.KeyShare, error) {
	var keyShares []types.KeyShare
	
	err := k.keyShares.Walk(ctx, nil, func(key collections.Pair[uint64, uint32], share types.KeyShare) (bool, error) {
		keyShares = append(keyShares, share)
		return false, nil // Continue iteration
	})
	
	return keyShares, err
}

// GetAllConditionContracts retrieves all condition contracts
func (k Keeper) GetAllConditionContracts(ctx context.Context) ([]types.ConditionContract, error) {
	var contracts []types.ConditionContract
	
	err := k.conditionContracts.Walk(ctx, nil, func(key string, contract types.ConditionContract) (bool, error) {
		contracts = append(contracts, contract)
		return false, nil // Continue iteration
	})
	
	return contracts, err
}

// Logger returns a module-specific logger
func (k Keeper) Logger(ctx context.Context) log.Logger {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	return sdkCtx.Logger().With("module", "x/"+types.ModuleName)
}

// BeginBlocker processes module logic at the beginning of each block
func (k Keeper) BeginBlocker(ctx context.Context) error {
	// Check for expired capsules and update their status
	return k.processExpiredCapsules(ctx)
}

// EndBlocker processes module logic at the end of each block  
func (k Keeper) EndBlocker(ctx context.Context) error {
	// Process any end-block logic here
	return nil
}

// processExpiredCapsules checks for and processes expired capsules
func (k Keeper) processExpiredCapsules(ctx context.Context) error {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	currentTime := sdkCtx.BlockTime()
	
	return k.capsules.Walk(ctx, nil, func(key uint64, capsule types.TimeCapsule) (bool, error) {
		// Check if capsule should be expired
		if capsule.Status == types.CapsuleStatus_ACTIVE {
			shouldExpire := false
			
			// Check expiration conditions based on capsule type
			switch capsule.CapsuleType {
			case types.CapsuleType_TIME_LOCK:
				if capsule.UnlockTime != nil && currentTime.After(*capsule.UnlockTime) {
					// For time-locked capsules, they don't expire but become unlockable
					// This is handled in the IsUnlockable method
				}
			case types.CapsuleType_DEAD_MANS_SWITCH:
				if capsule.ExpiresAt != nil && currentTime.After(*capsule.ExpiresAt) {
					shouldExpire = true
				}
			}
			
			// Update status if expired
			if shouldExpire {
				capsule.Status = types.CapsuleStatus_EXPIRED
				capsule.UpdatedAt = currentTime
				
				if err := k.capsules.Set(ctx, key, capsule); err != nil {
					return false, fmt.Errorf("failed to update expired capsule %d: %w", key, err)
				}
				
				// Emit expiration event
				sdkCtx.EventManager().EmitEvent(
					sdk.NewEvent(
						"capsule_expired",
						sdk.NewAttribute("capsule_id", fmt.Sprintf("%d", key)),
						sdk.NewAttribute("owner", capsule.Owner),
					),
				)
			}
		}
		
		return false, nil // Continue iteration
	})
}

// shareToBytes is a helper method to make the method public
func (k Keeper) shareToBytes(share *crypto.Share) []byte {
	return crypto.ShareToBytes(share)
}