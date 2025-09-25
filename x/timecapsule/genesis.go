package timecapsule

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/timecapsule/keeper"
	"github.com/cosmos/cosmos-sdk/x/timecapsule/types"
)

// GenesisState defines the time capsule module's genesis state.
type GenesisState struct {
	Params         types.Params           `json:"params"`
	Capsules       []types.TimeCapsule    `json:"capsules"`
	KeyShares      []types.KeyShare       `json:"key_shares"`
	CapsuleCounter uint64                 `json:"capsule_counter"`
	ConditionContracts []types.ConditionContract `json:"condition_contracts"`
}

// DefaultGenesis returns the default time capsule genesis state
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		Params:             types.DefaultParams(),
		Capsules:           []types.TimeCapsule{},
		KeyShares:          []types.KeyShare{},
		CapsuleCounter:     0,
		ConditionContracts: []types.ConditionContract{},
	}
}

// ValidateGenesis validates the time capsule module's genesis state
func ValidateGenesis(genState *GenesisState) error {
	// Validate parameters
	if err := genState.Params.Validate(); err != nil {
		return fmt.Errorf("invalid params: %w", err)
	}

	// Validate capsules
	capsuleIDs := make(map[uint64]bool)
	for i, capsule := range genState.Capsules {
		if err := capsule.Validate(); err != nil {
			return fmt.Errorf("invalid capsule at index %d: %w", i, err)
		}
		
		// Check for duplicate capsule IDs
		if capsuleIDs[capsule.ID] {
			return fmt.Errorf("duplicate capsule ID %d", capsule.ID)
		}
		capsuleIDs[capsule.ID] = true
		
		// Validate that capsule ID doesn't exceed counter
		if capsule.ID > genState.CapsuleCounter {
			return fmt.Errorf("capsule ID %d exceeds counter %d", capsule.ID, genState.CapsuleCounter)
		}
	}

	// Validate key shares
	for i, keyShare := range genState.KeyShares {
		if keyShare.CapsuleID == 0 {
			return fmt.Errorf("key share at index %d has invalid capsule ID", i)
		}
		
		if keyShare.NodeID == "" {
			return fmt.Errorf("key share at index %d has empty node ID", i)
		}
		
		if len(keyShare.EncryptedShare) == 0 {
			return fmt.Errorf("key share at index %d has empty encrypted share", i)
		}
		
		// Validate that key share references existing capsule
		if !capsuleIDs[keyShare.CapsuleID] {
			return fmt.Errorf("key share references non-existent capsule ID %d", keyShare.CapsuleID)
		}
	}

	// Validate condition contracts
	contractAddresses := make(map[string]bool)
	for i, contract := range genState.ConditionContracts {
		if contract.Address == "" {
			return fmt.Errorf("condition contract at index %d has empty address", i)
		}
		
		if contractAddresses[contract.Address] {
			return fmt.Errorf("duplicate condition contract address %s", contract.Address)
		}
		contractAddresses[contract.Address] = true
		
		if contract.Type == "" {
			return fmt.Errorf("condition contract at index %d has empty type", i)
		}
		
		if contract.CreatedBy == "" {
			return fmt.Errorf("condition contract at index %d has empty creator", i)
		}
	}

	return nil
}

// InitGenesis initializes the time capsule module's state from a provided genesis state.
func InitGenesis(ctx sdk.Context, k keeper.Keeper, genState *GenesisState) {
	// Set parameters
	if err := k.SetParams(ctx, genState.Params); err != nil {
		panic(fmt.Errorf("failed to set params: %w", err))
	}

	// Initialize capsule counter
	if err := k.SetCapsuleCounter(ctx, genState.CapsuleCounter); err != nil {
		panic(fmt.Errorf("failed to set capsule counter: %w", err))
	}

	// Initialize capsules
	for _, capsule := range genState.Capsules {
		if err := k.SetCapsule(ctx, &capsule); err != nil {
			panic(fmt.Errorf("failed to set capsule %d: %w", capsule.ID, err))
		}
		
		// Index by user
		if err := k.IndexUserCapsule(ctx, capsule.Owner, capsule.ID); err != nil {
			panic(fmt.Errorf("failed to index user capsule: %w", err))
		}
	}

	// Initialize key shares
	for _, keyShare := range genState.KeyShares {
		if err := k.SetKeyShare(ctx, &keyShare); err != nil {
			panic(fmt.Errorf("failed to set key share for capsule %d: %w", keyShare.CapsuleID, err))
		}
	}

	// Initialize condition contracts
	for _, contract := range genState.ConditionContracts {
		if err := k.SetConditionContract(ctx, &contract); err != nil {
			panic(fmt.Errorf("failed to set condition contract %s: %w", contract.Address, err))
		}
	}

	k.Logger(ctx).Info("Time capsule module genesis initialized",
		"capsules", len(genState.Capsules),
		"key_shares", len(genState.KeyShares),
		"condition_contracts", len(genState.ConditionContracts),
		"counter", genState.CapsuleCounter,
	)
}

// ExportGenesis returns the time capsule module's exported genesis
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *GenesisState {
	genesis := DefaultGenesis()
	
	// Export parameters
	params, err := k.GetParams(ctx)
	if err != nil {
		panic(fmt.Errorf("failed to get params: %w", err))
	}
	genesis.Params = params

	// Export capsule counter
	counter, err := k.GetCapsuleCounter(ctx)
	if err != nil {
		panic(fmt.Errorf("failed to get capsule counter: %w", err))
	}
	genesis.CapsuleCounter = counter

	// Export capsules
	capsules, err := k.GetAllCapsules(ctx)
	if err != nil {
		panic(fmt.Errorf("failed to get all capsules: %w", err))
	}
	genesis.Capsules = capsules

	// Export key shares
	keyShares, err := k.GetAllKeyShares(ctx)
	if err != nil {
		panic(fmt.Errorf("failed to get all key shares: %w", err))
	}
	genesis.KeyShares = keyShares

	// Export condition contracts
	contracts, err := k.GetAllConditionContracts(ctx)
	if err != nil {
		panic(fmt.Errorf("failed to get all condition contracts: %w", err))
	}
	genesis.ConditionContracts = contracts

	return genesis
}