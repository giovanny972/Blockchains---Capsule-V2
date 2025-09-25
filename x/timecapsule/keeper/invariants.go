package keeper

import (
	"fmt"
	"time"

	"cosmossdk.io/collections"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/timecapsule/types"
)

// RegisterInvariants registers all module invariants
func RegisterInvariants(ir sdk.InvariantRegistry, k Keeper) {
	ir.RegisterRoute(types.ModuleName, "capsule-key-shares", CapsuleKeySharesInvariant(k))
	ir.RegisterRoute(types.ModuleName, "capsule-user-index", CapsuleUserIndexInvariant(k))
	ir.RegisterRoute(types.ModuleName, "capsule-status-consistency", CapsuleStatusConsistencyInvariant(k))
}

// CapsuleKeySharesInvariant checks that each capsule has the correct number of key shares
func CapsuleKeySharesInvariant(k Keeper) sdk.Invariant {
	return func(ctx sdk.Context) (string, bool) {
		var (
			broken bool
			msg    string
		)

		// Check each capsule has the expected number of key shares
		err := k.capsules.Walk(ctx, nil, func(capsuleID uint64, capsule types.TimeCapsule) (bool, error) {
			// Count key shares for this capsule
			shareCount := uint32(0)
			prefix := collections.NewPrefixedPairRange[uint64, uint32](capsuleID)
			
			walkErr := k.keyShares.Walk(ctx, prefix, func(key collections.Pair[uint64, uint32], share types.KeyShare) (bool, error) {
				shareCount++
				return false, nil // Continue counting
			})
			
			if walkErr != nil {
				return false, walkErr
			}

			// Check if share count matches expected
			if shareCount != capsule.TotalShares {
				broken = true
				msg += fmt.Sprintf("capsule %d: expected %d key shares, found %d\n", 
					capsuleID, capsule.TotalShares, shareCount)
			}

			return false, nil // Continue iteration
		})

		if err != nil {
			broken = true
			msg += fmt.Sprintf("error walking capsules: %s\n", err.Error())
		}

		return sdk.FormatInvariant(types.ModuleName, "capsule-key-shares", msg), broken
	}
}

// CapsuleUserIndexInvariant checks that the user index is consistent with capsule ownership
func CapsuleUserIndexInvariant(k Keeper) sdk.Invariant {
	return func(ctx sdk.Context) (string, bool) {
		var (
			broken bool
			msg    string
		)

		// Track capsules found in user index
		userIndexedCapsules := make(map[uint64]string)

		// Walk through user capsule index
		err := k.userCapsules.Walk(ctx, nil, func(key collections.Pair[string, uint64]) (bool, error) {
			owner := key.K1()
			capsuleID := key.K2()
			
			// Check if capsule exists and has correct owner
			capsule, err := k.GetCapsule(ctx, capsuleID)
			if err != nil {
				broken = true
				msg += fmt.Sprintf("user index references non-existent capsule %d for owner %s\n", 
					capsuleID, owner)
				return false, nil
			}

			if capsule.Owner != owner {
				broken = true
				msg += fmt.Sprintf("user index inconsistency: capsule %d indexed under %s but owned by %s\n", 
					capsuleID, owner, capsule.Owner)
			}

			userIndexedCapsules[capsuleID] = owner
			return false, nil // Continue iteration
		})

		if err != nil {
			broken = true
			msg += fmt.Sprintf("error walking user capsule index: %s\n", err.Error())
		}

		// Check that all capsules are properly indexed
		err = k.capsules.Walk(ctx, nil, func(capsuleID uint64, capsule types.TimeCapsule) (bool, error) {
			indexedOwner, found := userIndexedCapsules[capsuleID]
			if !found {
				broken = true
				msg += fmt.Sprintf("capsule %d owned by %s is not in user index\n", 
					capsuleID, capsule.Owner)
			} else if indexedOwner != capsule.Owner {
				broken = true
				msg += fmt.Sprintf("capsule %d: index shows owner %s, actual owner %s\n", 
					capsuleID, indexedOwner, capsule.Owner)
			}

			return false, nil // Continue iteration
		})

		if err != nil {
			broken = true
			msg += fmt.Sprintf("error walking capsules for index check: %s\n", err.Error())
		}

		return sdk.FormatInvariant(types.ModuleName, "capsule-user-index", msg), broken
	}
}

// CapsuleStatusConsistencyInvariant checks that capsule status is consistent with conditions
func CapsuleStatusConsistencyInvariant(k Keeper) sdk.Invariant {
	return func(ctx sdk.Context) (string, bool) {
		var (
			broken bool
			msg    string
		)

		currentTime := ctx.BlockTime()

		err := k.capsules.Walk(ctx, nil, func(capsuleID uint64, capsule types.TimeCapsule) (bool, error) {
			// Check time-locked capsules
			if capsule.CapsuleType == types.CapsuleType_TIME_LOCK && capsule.UnlockTime != nil {
				if capsule.Status == types.CapsuleStatus_ACTIVE && currentTime.After(*capsule.UnlockTime) {
					// This is not necessarily broken - the capsule becomes unlockable but stays active
					// until someone actually opens it
				}
			}

			// Check dead man's switch capsules
			if capsule.CapsuleType == types.CapsuleType_DEAD_MANS_SWITCH {
				if capsule.InactivityPeriod > 0 && capsule.LastActivity != nil {
					if capsule.Status == types.CapsuleStatus_ACTIVE {
						// Check if enough time has passed for the switch to trigger
						switchTime := capsule.LastActivity.Add(time.Duration(capsule.InactivityPeriod) * time.Second)
						if currentTime.After(switchTime) {
							// Capsule should be unlockable by recipient
							// This is handled by the IsUnlockable method
						}
					}
				}
			}

			// Check that cancelled/expired capsules don't have inconsistent states
			if capsule.Status == types.CapsuleStatus_CANCELLED || capsule.Status == types.CapsuleStatus_EXPIRED {
				// These should not be unlockable
				if capsule.IsUnlockable(ctx) {
					broken = true
					msg += fmt.Sprintf("capsule %d has status %s but is still unlockable\n", 
						capsuleID, capsule.Status.String())
				}
			}

			// Validate threshold vs total shares
			if capsule.Threshold > capsule.TotalShares {
				broken = true
				msg += fmt.Sprintf("capsule %d: threshold %d exceeds total shares %d\n", 
					capsuleID, capsule.Threshold, capsule.TotalShares)
			}

			// Check that required fields are set for each capsule type
			switch capsule.CapsuleType {
			case types.CapsuleType_TIME_LOCK:
				if capsule.UnlockTime == nil {
					broken = true
					msg += fmt.Sprintf("time-locked capsule %d missing unlock time\n", capsuleID)
				}
			case types.CapsuleType_CONDITIONAL:
				if capsule.ConditionContract == "" {
					broken = true
					msg += fmt.Sprintf("conditional capsule %d missing condition contract\n", capsuleID)
				}
			case types.CapsuleType_MULTI_SIG:
				if capsule.RequiredSigs == 0 {
					broken = true
					msg += fmt.Sprintf("multi-sig capsule %d missing required signatures\n", capsuleID)
				}
			case types.CapsuleType_DEAD_MANS_SWITCH:
				if capsule.InactivityPeriod == 0 {
					broken = true
					msg += fmt.Sprintf("dead man's switch capsule %d missing inactivity period\n", capsuleID)
				}
				if capsule.Recipient == "" {
					broken = true
					msg += fmt.Sprintf("dead man's switch capsule %d missing recipient\n", capsuleID)
				}
			}

			return false, nil // Continue iteration
		})

		if err != nil {
			broken = true
			msg += fmt.Sprintf("error walking capsules: %s\n", err.Error())
		}

		return sdk.FormatInvariant(types.ModuleName, "capsule-status-consistency", msg), broken
	}
}