package keeper

import (
	"context"
	"encoding/json"
	"fmt"

	"cosmossdk.io/collections"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/timecapsule/crypto"
	"github.com/cosmos/cosmos-sdk/x/timecapsule/types"
)

// MsgServer implements the module's message server
type MsgServer struct {
	keeper Keeper
}

// NewMsgServerImpl returns an implementation of the module's MsgServer interface
func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &MsgServer{keeper: keeper}
}

var _ types.MsgServer = MsgServer{}

// CreateCapsule creates a new time capsule
func (ms MsgServer) CreateCapsule(goCtx context.Context, msg *types.MsgCreateCapsule) (*types.MsgCreateCapsuleResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Charge creation fee from creator
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return nil, err
	}

	// Get module parameters
	params, err := ms.keeper.GetParams(ctx)
	if err != nil {
		return nil, err
	}

	// Validate capsule type is allowed
	allowed := false
	for _, allowedType := range params.AllowedCapsuleTypes {
		if allowedType == msg.CapsuleType {
			allowed = true
			break
		}
	}
	if !allowed {
		return nil, types.ErrInvalidCapsuleType.Wrapf("capsule type %s not allowed", msg.CapsuleType.String())
	}

	// Validate data size
	if uint64(len(msg.Data)) > params.MaxDataSize {
		return nil, types.ErrDataTooLarge.Wrapf("data size %d exceeds maximum %d", len(msg.Data), params.MaxDataSize)
	}

	// Validate threshold and shares
	if msg.Threshold < params.MinThreshold {
		return nil, types.ErrInvalidThreshold.Wrapf("threshold %d below minimum %d", msg.Threshold, params.MinThreshold)
	}
	if msg.TotalShares > params.MaxShares {
		return nil, types.ErrInvalidThreshold.Wrapf("total shares %d exceeds maximum %d", msg.TotalShares, params.MaxShares)
	}

	// Charge creation fee
	if !params.CreationFee.IsZero() {
		if err := ms.keeper.bankKeeper.SendCoinsFromAccountToModule(
			ctx, creator, types.ModuleName, params.CreationFee,
		); err != nil {
			return nil, err
		}
	}

	// Prepare metadata
	metadata := make(map[string]string)
	if msg.Title != "" {
		metadata["title"] = msg.Title
	}
	if msg.Description != "" {
		metadata["description"] = msg.Description
	}
	for k, v := range msg.Metadata {
		metadata[k] = v
	}

	// Create the capsule
	capsule, err := ms.keeper.CreateCapsule(
		ctx,
		msg.Creator,
		msg.Recipient,
		msg.Data,
		msg.CapsuleType,
		msg.Threshold,
		msg.TotalShares,
		msg.UnlockTime,
		msg.ConditionContract,
		metadata,
	)
	if err != nil {
		return nil, err
	}

	return &types.MsgCreateCapsuleResponse{
		CapsuleId: capsule.ID,
	}, nil
}

// OpenCapsule opens a time capsule and retrieves its data
func (ms MsgServer) OpenCapsule(goCtx context.Context, msg *types.MsgOpenCapsule) (*types.MsgOpenCapsuleResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Parse key shares if provided
	var shares []*crypto.Share
	if len(msg.KeyShares) > 0 {
		shares = make([]*crypto.Share, len(msg.KeyShares))
		for i, shareStr := range msg.KeyShares {
			var share crypto.Share
			if err := json.Unmarshal([]byte(shareStr), &share); err != nil {
				return nil, types.ErrInvalidKeyShare.Wrapf("invalid key share at index %d: %s", i, err)
			}
			shares[i] = &share
		}
	}

	// Prepare condition parameters
	conditionParams := make(map[string]interface{})
	if len(msg.Signatures) > 0 {
		conditionParams["signatures"] = msg.Signatures
	}
	for k, v := range msg.ConditionProof {
		conditionParams[k] = v
	}

	// Open the capsule
	data, err := ms.keeper.OpenCapsule(ctx, msg.CapsuleID, msg.Accessor, shares, conditionParams)
	if err != nil {
		return nil, err
	}

	return &types.MsgOpenCapsuleResponse{
		Data: data,
	}, nil
}

// UpdateActivity updates the last activity timestamp for dead man's switch capsules
func (ms MsgServer) UpdateActivity(goCtx context.Context, msg *types.MsgUpdateActivity) (*types.MsgUpdateActivityResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	err := ms.keeper.UpdateLastActivity(ctx, msg.CapsuleID, msg.Owner)
	if err != nil {
		return nil, err
	}

	return &types.MsgUpdateActivityResponse{}, nil
}

// CancelCapsule cancels a time capsule
func (ms MsgServer) CancelCapsule(goCtx context.Context, msg *types.MsgCancelCapsule) (*types.MsgCancelCapsuleResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Get the capsule
	capsule, err := ms.keeper.GetCapsule(ctx, msg.CapsuleID)
	if err != nil {
		return nil, err
	}

	// Verify ownership
	if capsule.Owner != msg.Owner {
		return nil, types.ErrUnauthorized.Wrap("only owner can cancel capsule")
	}

	// Check if capsule can be cancelled
	if capsule.Status != types.CapsuleStatus_ACTIVE {
		return nil, types.ErrInvalidCapsule.Wrapf("cannot cancel capsule with status %s", capsule.Status.String())
	}

	// Update status to cancelled
	capsule.Status = types.CapsuleStatus_CANCELLED
	capsule.UpdatedAt = ctx.BlockTime()

	if err := ms.keeper.capsules.Set(ctx, msg.CapsuleID, *capsule); err != nil {
		return nil, fmt.Errorf("failed to update capsule status: %w", err)
	}

	// Emit event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			"capsule_cancelled",
			sdk.NewAttribute("capsule_id", fmt.Sprintf("%d", msg.CapsuleID)),
			sdk.NewAttribute("owner", msg.Owner),
			sdk.NewAttribute("reason", msg.Reason),
		),
	)

	return &types.MsgCancelCapsuleResponse{}, nil
}

// TransferCapsule transfers ownership of a capsule
func (ms MsgServer) TransferCapsule(goCtx context.Context, msg *types.MsgTransferCapsule) (*types.MsgTransferCapsuleResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Get the capsule
	capsule, err := ms.keeper.GetCapsule(ctx, msg.CapsuleID)
	if err != nil {
		return nil, err
	}

	// Verify current ownership
	if capsule.Owner != msg.CurrentOwner {
		return nil, types.ErrUnauthorized.Wrap("only current owner can transfer capsule")
	}

	// Check if capsule can be transferred
	if capsule.Status != types.CapsuleStatus_ACTIVE {
		return nil, types.ErrInvalidCapsule.Wrapf("cannot transfer capsule with status %s", capsule.Status.String())
	}

	// Update ownership
	capsule.Owner = msg.NewOwner
	capsule.UpdatedAt = ctx.BlockTime()

	if err := ms.keeper.capsules.Set(ctx, msg.CapsuleID, *capsule); err != nil {
		return nil, fmt.Errorf("failed to update capsule ownership: %w", err)
	}

	// Update user index
	if err := ms.keeper.userCapsules.Remove(ctx, collections.Join(msg.CurrentOwner, msg.CapsuleID)); err != nil {
		return nil, fmt.Errorf("failed to remove old user index: %w", err)
	}

	if err := ms.keeper.userCapsules.Set(ctx, collections.Join(msg.NewOwner, msg.CapsuleID)); err != nil {
		return nil, fmt.Errorf("failed to add new user index: %w", err)
	}

	// Emit event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			"capsule_transferred",
			sdk.NewAttribute("capsule_id", fmt.Sprintf("%d", msg.CapsuleID)),
			sdk.NewAttribute("from", msg.CurrentOwner),
			sdk.NewAttribute("to", msg.NewOwner),
		),
	)

	return &types.MsgTransferCapsuleResponse{}, nil
}

// BatchTransferCapsules transfers multiple capsules in a single transaction
func (ms MsgServer) BatchTransferCapsules(goCtx context.Context, msg *types.MsgBatchTransferCapsules) (*types.MsgBatchTransferCapsulesResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	var transferredCapsules []uint64
	var failedTransfers []types.FailedTransfer

	for _, transfer := range msg.Transfers {
		// Get the capsule
		capsule, err := ms.keeper.GetCapsule(ctx, transfer.CapsuleID)
		if err != nil {
			failedTransfers = append(failedTransfers, types.FailedTransfer{
				CapsuleID: transfer.CapsuleID,
				Reason:    fmt.Sprintf("capsule not found: %s", err),
			})
			continue
		}

		// Verify current ownership
		if capsule.Owner != msg.CurrentOwner {
			failedTransfers = append(failedTransfers, types.FailedTransfer{
				CapsuleID: transfer.CapsuleID,
				Reason:    "unauthorized: not the owner",
			})
			continue
		}

		// Check if capsule can be transferred
		if capsule.Status != types.CapsuleStatus_ACTIVE {
			failedTransfers = append(failedTransfers, types.FailedTransfer{
				CapsuleID: transfer.CapsuleID,
				Reason:    fmt.Sprintf("cannot transfer capsule with status %s", capsule.Status.String()),
			})
			continue
		}

		// Perform transfer
		if err := ms.keeper.TransferCapsuleOwnership(ctx, transfer.CapsuleID, msg.CurrentOwner, transfer.NewOwner, "batch", transfer.Message); err != nil {
			failedTransfers = append(failedTransfers, types.FailedTransfer{
				CapsuleID: transfer.CapsuleID,
				Reason:    fmt.Sprintf("transfer failed: %s", err),
			})
			continue
		}

		transferredCapsules = append(transferredCapsules, transfer.CapsuleID)
	}

	// Charge transfer fees if specified
	if !msg.TransferFee.IsZero() {
		creator, _ := sdk.AccAddressFromBech32(msg.CurrentOwner)
		if err := ms.keeper.bankKeeper.SendCoinsFromAccountToModule(ctx, creator, types.ModuleName, msg.TransferFee); err != nil {
			return nil, fmt.Errorf("failed to charge transfer fee: %w", err)
		}
	}

	// Emit batch transfer event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			"batch_capsules_transferred",
			sdk.NewAttribute("from", msg.CurrentOwner),
			sdk.NewAttribute("transferred_count", fmt.Sprintf("%d", len(transferredCapsules))),
			sdk.NewAttribute("failed_count", fmt.Sprintf("%d", len(failedTransfers))),
		),
	)

	return &types.MsgBatchTransferCapsulesResponse{
		TransferredCapsules: transferredCapsules,
		FailedTransfers:     failedTransfers,
	}, nil
}

// ApproveTransfer approves or rejects a pending transfer
func (ms MsgServer) ApproveTransfer(goCtx context.Context, msg *types.MsgApproveTransfer) (*types.MsgApproveTransferResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Get the pending transfer
	pendingTransfer, err := ms.keeper.GetPendingTransfer(ctx, msg.TransferID)
	if err != nil {
		return nil, fmt.Errorf("pending transfer not found: %w", err)
	}

	// Check if transfer has expired
	if ctx.BlockTime().After(pendingTransfer.ExpiryTime) {
		pendingTransfer.Status = "expired"
		ms.keeper.SetPendingTransfer(ctx, msg.TransferID, *pendingTransfer)
		return nil, types.ErrInvalidTransfer.Wrap("transfer has expired")
	}

	// Verify approver is the recipient
	if pendingTransfer.ToOwner != msg.Approver {
		return nil, types.ErrUnauthorized.Wrap("only the recipient can approve the transfer")
	}

	// Update transfer status
	if msg.Approved {
		pendingTransfer.Status = "approved"
		
		// Execute the transfer
		if err := ms.keeper.TransferCapsuleOwnership(ctx, msg.CapsuleID, pendingTransfer.FromOwner, pendingTransfer.ToOwner, "approved", pendingTransfer.Message); err != nil {
			return nil, fmt.Errorf("failed to execute approved transfer: %w", err)
		}
	} else {
		pendingTransfer.Status = "rejected"
	}

	// Update pending transfer
	ms.keeper.SetPendingTransfer(ctx, msg.TransferID, *pendingTransfer)

	// Emit approval event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			"transfer_approval",
			sdk.NewAttribute("transfer_id", msg.TransferID),
			sdk.NewAttribute("capsule_id", fmt.Sprintf("%d", msg.CapsuleID)),
			sdk.NewAttribute("approver", msg.Approver),
			sdk.NewAttribute("approved", fmt.Sprintf("%t", msg.Approved)),
		),
	)

	return &types.MsgApproveTransferResponse{
		Approved: msg.Approved,
	}, nil
}

// EmergencyDeleteContract permanently deletes a capsule's smart contract in emergency situations
func (ms MsgServer) EmergencyDeleteContract(goCtx context.Context, msg *types.MsgEmergencyDeleteContract) (*types.MsgEmergencyDeleteContractResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Execute emergency deletion with comprehensive validation
	err := ms.keeper.EmergencyDeleteContract(
		ctx,
		msg.CapsuleID,
		msg.Creator,
		msg.EmergencyReason,
		msg.ConfirmationCode,
	)
	if err != nil {
		return nil, err
	}

	return &types.MsgEmergencyDeleteContractResponse{
		Success: true,
		Message: "Smart contract deleted permanently",
	}, nil
}