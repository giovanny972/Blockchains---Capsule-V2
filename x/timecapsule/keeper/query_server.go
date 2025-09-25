package keeper

import (
	"context"

	"cosmossdk.io/collections"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/timecapsule/types"
)

// QueryServer implements the module query server
type QueryServer struct {
	keeper Keeper
}

// NewQueryServerImpl returns an implementation of the module QueryServer interface
func NewQueryServerImpl(keeper Keeper) types.QueryServer {
	return &QueryServer{keeper: keeper}
}

var _ types.QueryServer = QueryServer{}

// Params returns the module parameters
func (qs QueryServer) Params(c context.Context, req *types.QueryParamsRequest) (*types.QueryParamsResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	
	params, err := qs.keeper.GetParams(ctx)
	if err != nil {
		return nil, err
	}

	return &types.QueryParamsResponse{Params: params}, nil
}

// Capsule returns details of a specific capsule
func (qs QueryServer) Capsule(c context.Context, req *types.QueryCapsuleRequest) (*types.QueryCapsuleResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	if req.CapsuleId == 0 {
		return nil, types.ErrCapsuleNotFound.Wrap("capsule ID cannot be zero")
	}

	capsule, err := qs.keeper.GetCapsule(ctx, req.CapsuleId)
	if err != nil {
		return nil, err
	}

	return &types.QueryCapsuleResponse{Capsule: capsule}, nil
}

// Capsules returns a list of capsules with pagination
func (qs QueryServer) Capsules(c context.Context, req *types.QueryCapsulesRequest) (*types.QueryCapsulesResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	var capsules []types.TimeCapsule
	
	err := qs.keeper.capsules.Walk(ctx, nil, func(key uint64, capsule types.TimeCapsule) (bool, error) {
		capsules = append(capsules, capsule)
		return false, nil // Continue iteration
	})
	if err != nil {
		return nil, err
	}

	return &types.QueryCapsulesResponse{
		Capsules: capsules,
	}, nil
}

// UserCapsules returns all capsules owned by a specific user
func (qs QueryServer) UserCapsules(c context.Context, req *types.QueryUserCapsulesRequest) (*types.QueryUserCapsulesResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	if req.Owner == "" {
		return nil, types.ErrInvalidRecipient.Wrap("owner address cannot be empty")
	}

	// Validate owner address
	_, err := sdk.AccAddressFromBech32(req.Owner)
	if err != nil {
		return nil, types.ErrInvalidRecipient.Wrapf("invalid owner address: %s", err)
	}

	capsules, err := qs.keeper.ListUserCapsules(ctx, req.Owner)
	if err != nil {
		return nil, err
	}

	// Convert to response format
	var responseCapsules []*types.TimeCapsule
	for _, capsule := range capsules {
		responseCapsules = append(responseCapsules, capsule)
	}

	return &types.QueryUserCapsulesResponse{
		Capsules: responseCapsules,
		Owner:    req.Owner,
	}, nil
}

// CapsulesByType returns capsules filtered by type
func (qs QueryServer) CapsulesByType(c context.Context, req *types.QueryCapsulesByTypeRequest) (*types.QueryCapsulesByTypeResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	var capsules []types.TimeCapsule
	
	err := qs.keeper.capsules.Walk(ctx, nil, func(key uint64, capsule types.TimeCapsule) (bool, error) {
		if capsule.CapsuleType == req.CapsuleType {
			capsules = append(capsules, capsule)
		}
		return false, nil // Continue iteration
	})
	if err != nil {
		return nil, err
	}

	return &types.QueryCapsulesByTypeResponse{
		Capsules:    capsules,
		CapsuleType: req.CapsuleType,
	}, nil
}

// CapsulesByStatus returns capsules filtered by status
func (qs QueryServer) CapsulesByStatus(c context.Context, req *types.QueryCapsulesByStatusRequest) (*types.QueryCapsulesByStatusResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	var capsules []types.TimeCapsule
	
	err := qs.keeper.capsules.Walk(ctx, nil, func(key uint64, capsule types.TimeCapsule) (bool, error) {
		if capsule.Status == req.Status {
			capsules = append(capsules, capsule)
		}
		return false, nil // Continue iteration
	})
	if err != nil {
		return nil, err
	}

	return &types.QueryCapsulesByStatusResponse{
		Capsules: capsules,
		Status:   req.Status,
	}, nil
}

// Stats returns statistics about the time capsule module
func (qs QueryServer) Stats(c context.Context, req *types.QueryStatsRequest) (*types.QueryStatsResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	stats := &types.ModuleStats{
		TotalCapsules:   0,
		ActiveCapsules:  0,
		OpenedCapsules:  0,
		ExpiredCapsules: 0,
		CancelledCapsules: 0,
	}

	// Count capsules by status
	err := qs.keeper.capsules.Walk(ctx, nil, func(key uint64, capsule types.TimeCapsule) (bool, error) {
		stats.TotalCapsules++
		
		switch capsule.Status {
		case types.CapsuleStatus_ACTIVE:
			stats.ActiveCapsules++
		case types.CapsuleStatus_UNLOCKED:
			stats.OpenedCapsules++
		case types.CapsuleStatus_EXPIRED:
			stats.ExpiredCapsules++
		case types.CapsuleStatus_CANCELLED:
			stats.CancelledCapsules++
		}
		
		return false, nil // Continue iteration
	})
	if err != nil {
		return nil, err
	}

	return &types.QueryStatsResponse{Stats: stats}, nil
}

// KeyShares returns key shares for a specific capsule
func (qs QueryServer) KeyShares(c context.Context, req *types.QueryKeySharesRequest) (*types.QueryKeySharesResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	if req.CapsuleId == 0 {
		return nil, types.ErrCapsuleNotFound.Wrap("capsule ID cannot be zero")
	}

	var keyShares []types.KeyShare
	
	// Iterate through key shares for this capsule
	prefix := collections.NewPrefixedPairRange[uint64, uint32](req.CapsuleId)
	err := qs.keeper.keyShares.Walk(ctx, prefix, func(key collections.Pair[uint64, uint32], share types.KeyShare) (bool, error) {
		keyShares = append(keyShares, share)
		return false, nil // Continue iteration
	})
	if err != nil {
		return nil, err
	}

	return &types.QueryKeySharesResponse{
		KeyShares: keyShares,
		CapsuleId: req.CapsuleId,
	}, nil
}

// ConditionContract returns details of a condition contract
func (qs QueryServer) ConditionContract(c context.Context, req *types.QueryConditionContractRequest) (*types.QueryConditionContractResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	if req.Address == "" {
		return nil, types.ErrInvalidCapsule.Wrap("contract address cannot be empty")
	}

	contract, err := qs.keeper.conditionContracts.Get(ctx, req.Address)
	if err != nil {
		if err.Error() == "not found" {
			return nil, types.ErrContractExecution.Wrapf("condition contract %s not found", req.Address)
		}
		return nil, err
	}

	return &types.QueryConditionContractResponse{
		Contract: &contract,
	}, nil
}

// ConditionContracts returns all condition contracts
func (qs QueryServer) ConditionContracts(c context.Context, req *types.QueryConditionContractsRequest) (*types.QueryConditionContractsResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	var contracts []types.ConditionContract
	
	err := qs.keeper.conditionContracts.Walk(ctx, nil, func(key string, contract types.ConditionContract) (bool, error) {
		contracts = append(contracts, contract)
		return false, nil // Continue iteration
	})
	if err != nil {
		return nil, err
	}

	return &types.QueryConditionContractsResponse{
		Contracts: contracts,
	}, nil
}