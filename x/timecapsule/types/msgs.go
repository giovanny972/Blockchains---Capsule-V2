package types

import (
	"time"
	"fmt"
	
	sdk "github.com/cosmos/cosmos-sdk/types"
	"cosmossdk.io/errors"
)

// Message types for time capsule operations
const (
	TypeMsgCreateCapsule     = "create_capsule"
	TypeMsgOpenCapsule       = "open_capsule"
	TypeMsgUpdateActivity    = "update_activity"
	TypeMsgCancelCapsule     = "cancel_capsule"
	TypeMsgTransferCapsule   = "transfer_capsule"
	TypeMsgBatchTransferCapsules = "batch_transfer_capsules"
	TypeMsgApproveTransfer   = "approve_transfer"
	TypeMsgEmergencyDeleteContract = "emergency_delete_contract"
)

// MsgCreateCapsule defines the message to create a new time capsule
type MsgCreateCapsule struct {
	Creator           string            `json:"creator"`
	Recipient         string            `json:"recipient,omitempty"`
	Data              []byte            `json:"data"`
	CapsuleType       CapsuleType       `json:"capsule_type"`
	Threshold         uint32            `json:"threshold"`
	TotalShares       uint32            `json:"total_shares"`
	UnlockTime        *time.Time        `json:"unlock_time,omitempty"`
	ConditionContract string            `json:"condition_contract,omitempty"`
	RequiredSigs      uint32            `json:"required_sigs,omitempty"`
	InactivityPeriod  uint64            `json:"inactivity_period,omitempty"`
	Title             string            `json:"title,omitempty"`
	Description       string            `json:"description,omitempty"`
	Metadata          map[string]string `json:"metadata,omitempty"`
}

// NewMsgCreateCapsule creates a new MsgCreateCapsule
func NewMsgCreateCapsule(
	creator string,
	recipient string,
	data []byte,
	capsuleType CapsuleType,
	threshold uint32,
	totalShares uint32,
) *MsgCreateCapsule {
	return &MsgCreateCapsule{
		Creator:     creator,
		Recipient:   recipient,
		Data:        data,
		CapsuleType: capsuleType,
		Threshold:   threshold,
		TotalShares: totalShares,
	}
}

// Route implements the sdk.Msg interface
func (msg *MsgCreateCapsule) Route() string {
	return RouterKey
}

// Type implements the sdk.Msg interface
func (msg *MsgCreateCapsule) Type() string {
	return TypeMsgCreateCapsule
}

// GetSigners implements the sdk.Msg interface
func (msg *MsgCreateCapsule) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

// GetSignBytes implements the sdk.Msg interface
func (msg *MsgCreateCapsule) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// ValidateBasic implements the sdk.Msg interface
func (msg *MsgCreateCapsule) ValidateBasic() error {
	// Validate creator address
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return errors.Wrapf(ErrInvalidAddress, "invalid creator address (%s)", err)
	}

	// Validate recipient address if provided
	if msg.Recipient != "" {
		_, err := sdk.AccAddressFromBech32(msg.Recipient)
		if err != nil {
			return errors.Wrapf(ErrInvalidAddress, "invalid recipient address (%s)", err)
		}
	}

	// Validate data
	if len(msg.Data) == 0 {
		return errors.Wrap(ErrInvalidCapsule, "data cannot be empty")
	}

	// Check data size limits (1MB max)
	maxDataSize := 1024 * 1024
	if len(msg.Data) > maxDataSize {
		return errors.Wrapf(ErrDataTooLarge, "data size %d exceeds maximum %d", len(msg.Data), maxDataSize)
	}

	// Validate threshold and total shares
	if msg.Threshold == 0 {
		return errors.Wrap(ErrInvalidThreshold, "threshold must be greater than 0")
	}

	if msg.TotalShares == 0 {
		return errors.Wrap(ErrInvalidThreshold, "total shares must be greater than 0")
	}

	if msg.Threshold > msg.TotalShares {
		return errors.Wrap(ErrInvalidThreshold, "threshold cannot exceed total shares")
	}

	// Validate capsule type specific requirements
	switch msg.CapsuleType {
	case CapsuleType_TIME_LOCK:
		if msg.UnlockTime == nil {
			return errors.Wrap(ErrInvalidTimelock, "time-locked capsule must have unlock time")
		}
		if msg.UnlockTime.Before(time.Now()) {
			return errors.Wrap(ErrInvalidTimelock, "unlock time must be in the future")
		}

	case CapsuleType_CONDITIONAL:
		if msg.ConditionContract == "" {
			return errors.Wrap(ErrConditionNotMet, "conditional capsule must have condition contract")
		}

	case CapsuleType_MULTI_SIG:
		if msg.RequiredSigs == 0 {
			return errors.Wrap(ErrInvalidSignature, "multi-sig capsule must have required signatures")
		}

	case CapsuleType_DEAD_MANS_SWITCH:
		if msg.InactivityPeriod == 0 {
			return errors.Wrap(ErrInvalidCapsule, "dead man's switch capsule must have inactivity period")
		}
		if msg.Recipient == "" {
			return errors.Wrap(ErrInvalidRecipient, "dead man's switch capsule must have recipient")
		}
	}

	return nil
}

// MsgOpenCapsule defines the message to open a time capsule
type MsgOpenCapsule struct {
	Accessor        string                 `json:"accessor"`
	CapsuleID       uint64                 `json:"capsule_id"`
	KeyShares       []string               `json:"key_shares,omitempty"`       // Serialized key shares
	Signatures      []string               `json:"signatures,omitempty"`       // For multi-sig capsules
	ConditionProof  map[string]interface{} `json:"condition_proof,omitempty"`  // Proof that conditions are met
}

// NewMsgOpenCapsule creates a new MsgOpenCapsule
func NewMsgOpenCapsule(accessor string, capsuleID uint64) *MsgOpenCapsule {
	return &MsgOpenCapsule{
		Accessor:  accessor,
		CapsuleID: capsuleID,
	}
}

// Route implements the sdk.Msg interface
func (msg *MsgOpenCapsule) Route() string {
	return RouterKey
}

// Type implements the sdk.Msg interface
func (msg *MsgOpenCapsule) Type() string {
	return TypeMsgOpenCapsule
}

// GetSigners implements the sdk.Msg interface
func (msg *MsgOpenCapsule) GetSigners() []sdk.AccAddress {
	accessor, err := sdk.AccAddressFromBech32(msg.Accessor)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{accessor}
}

// GetSignBytes implements the sdk.Msg interface
func (msg *MsgOpenCapsule) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// ValidateBasic implements the sdk.Msg interface
func (msg *MsgOpenCapsule) ValidateBasic() error {
	// Validate accessor address
	_, err := sdk.AccAddressFromBech32(msg.Accessor)
	if err != nil {
		return errors.Wrapf(ErrInvalidAddress, "invalid accessor address (%s)", err)
	}

	// Validate capsule ID
	if msg.CapsuleID == 0 {
		return errors.Wrap(ErrCapsuleNotFound, "capsule ID cannot be zero")
	}

	return nil
}

// MsgUpdateActivity defines the message to update last activity for dead man's switch
type MsgUpdateActivity struct {
	Owner     string `json:"owner"`
	CapsuleID uint64 `json:"capsule_id"`
}

// NewMsgUpdateActivity creates a new MsgUpdateActivity
func NewMsgUpdateActivity(owner string, capsuleID uint64) *MsgUpdateActivity {
	return &MsgUpdateActivity{
		Owner:     owner,
		CapsuleID: capsuleID,
	}
}

// Route implements the sdk.Msg interface
func (msg *MsgUpdateActivity) Route() string {
	return RouterKey
}

// Type implements the sdk.Msg interface
func (msg *MsgUpdateActivity) Type() string {
	return TypeMsgUpdateActivity
}

// GetSigners implements the sdk.Msg interface
func (msg *MsgUpdateActivity) GetSigners() []sdk.AccAddress {
	owner, err := sdk.AccAddressFromBech32(msg.Owner)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{owner}
}

// GetSignBytes implements the sdk.Msg interface
func (msg *MsgUpdateActivity) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// ValidateBasic implements the sdk.Msg interface
func (msg *MsgUpdateActivity) ValidateBasic() error {
	// Validate owner address
	_, err := sdk.AccAddressFromBech32(msg.Owner)
	if err != nil {
		return errors.Wrapf(ErrInvalidAddress, "invalid owner address (%s)", err)
	}

	// Validate capsule ID
	if msg.CapsuleID == 0 {
		return errors.Wrap(ErrCapsuleNotFound, "capsule ID cannot be zero")
	}

	return nil
}

// MsgCancelCapsule defines the message to cancel a time capsule
type MsgCancelCapsule struct {
	Owner     string `json:"owner"`
	CapsuleID uint64 `json:"capsule_id"`
	Reason    string `json:"reason,omitempty"`
}

// NewMsgCancelCapsule creates a new MsgCancelCapsule
func NewMsgCancelCapsule(owner string, capsuleID uint64, reason string) *MsgCancelCapsule {
	return &MsgCancelCapsule{
		Owner:     owner,
		CapsuleID: capsuleID,
		Reason:    reason,
	}
}

// Route implements the sdk.Msg interface
func (msg *MsgCancelCapsule) Route() string {
	return RouterKey
}

// Type implements the sdk.Msg interface
func (msg *MsgCancelCapsule) Type() string {
	return TypeMsgCancelCapsule
}

// GetSigners implements the sdk.Msg interface
func (msg *MsgCancelCapsule) GetSigners() []sdk.AccAddress {
	owner, err := sdk.AccAddressFromBech32(msg.Owner)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{owner}
}

// GetSignBytes implements the sdk.Msg interface
func (msg *MsgCancelCapsule) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// ValidateBasic implements the sdk.Msg interface
func (msg *MsgCancelCapsule) ValidateBasic() error {
	// Validate owner address
	_, err := sdk.AccAddressFromBech32(msg.Owner)
	if err != nil {
		return errors.Wrapf(ErrInvalidAddress, "invalid owner address (%s)", err)
	}

	// Validate capsule ID
	if msg.CapsuleID == 0 {
		return errors.Wrap(ErrCapsuleNotFound, "capsule ID cannot be zero")
	}

	return nil
}

// MsgTransferCapsule defines the message to transfer capsule ownership
type MsgTransferCapsule struct {
	CurrentOwner string `json:"current_owner"`
	NewOwner     string `json:"new_owner"`
	CapsuleID    uint64 `json:"capsule_id"`
}

// NewMsgTransferCapsule creates a new MsgTransferCapsule
func NewMsgTransferCapsule(currentOwner, newOwner string, capsuleID uint64) *MsgTransferCapsule {
	return &MsgTransferCapsule{
		CurrentOwner: currentOwner,
		NewOwner:     newOwner,
		CapsuleID:    capsuleID,
	}
}

// Route implements the sdk.Msg interface
func (msg *MsgTransferCapsule) Route() string {
	return RouterKey
}

// Type implements the sdk.Msg interface
func (msg *MsgTransferCapsule) Type() string {
	return TypeMsgTransferCapsule
}

// GetSigners implements the sdk.Msg interface
func (msg *MsgTransferCapsule) GetSigners() []sdk.AccAddress {
	owner, err := sdk.AccAddressFromBech32(msg.CurrentOwner)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{owner}
}

// GetSignBytes implements the sdk.Msg interface
func (msg *MsgTransferCapsule) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// ValidateBasic implements the sdk.Msg interface
func (msg *MsgTransferCapsule) ValidateBasic() error {
	// Validate current owner address
	_, err := sdk.AccAddressFromBech32(msg.CurrentOwner)
	if err != nil {
		return errors.Wrapf(ErrInvalidAddress, "invalid current owner address (%s)", err)
	}

	// Validate new owner address
	_, err = sdk.AccAddressFromBech32(msg.NewOwner)
	if err != nil {
		return errors.Wrapf(ErrInvalidAddress, "invalid new owner address (%s)", err)
	}

	// Check that addresses are different
	if msg.CurrentOwner == msg.NewOwner {
		return errors.Wrap(ErrInvalidRequest, "current and new owner must be different")
	}

	// Validate capsule ID
	if msg.CapsuleID == 0 {
		return errors.Wrap(ErrCapsuleNotFound, "capsule ID cannot be zero")
	}

	return nil
}

// MsgBatchTransferCapsules defines the message to transfer multiple capsules in a single transaction
type MsgBatchTransferCapsules struct {
	CurrentOwner  string            `json:"current_owner"`
	Transfers     []CapsuleTransfer `json:"transfers"`
	TransferFee   sdk.Coins         `json:"transfer_fee,omitempty"`
	RequireApproval bool            `json:"require_approval,omitempty"`
}

// CapsuleTransfer represents a single capsule transfer in a batch
type CapsuleTransfer struct {
	CapsuleID uint64 `json:"capsule_id"`
	NewOwner  string `json:"new_owner"`
	Message   string `json:"message,omitempty"`
}

// NewMsgBatchTransferCapsules creates a new MsgBatchTransferCapsules
func NewMsgBatchTransferCapsules(currentOwner string, transfers []CapsuleTransfer) *MsgBatchTransferCapsules {
	return &MsgBatchTransferCapsules{
		CurrentOwner: currentOwner,
		Transfers:    transfers,
	}
}

// Route implements the sdk.Msg interface
func (msg *MsgBatchTransferCapsules) Route() string {
	return RouterKey
}

// Type implements the sdk.Msg interface
func (msg *MsgBatchTransferCapsules) Type() string {
	return TypeMsgBatchTransferCapsules
}

// GetSigners implements the sdk.Msg interface
func (msg *MsgBatchTransferCapsules) GetSigners() []sdk.AccAddress {
	owner, err := sdk.AccAddressFromBech32(msg.CurrentOwner)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{owner}
}

// GetSignBytes implements the sdk.Msg interface
func (msg *MsgBatchTransferCapsules) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// ValidateBasic implements the sdk.Msg interface
func (msg *MsgBatchTransferCapsules) ValidateBasic() error {
	// Validate current owner address
	_, err := sdk.AccAddressFromBech32(msg.CurrentOwner)
	if err != nil {
		return errors.Wrapf(ErrInvalidAddress, "invalid current owner address (%s)", err)
	}

	// Validate transfers
	if len(msg.Transfers) == 0 {
		return errors.Wrap(ErrInvalidRequest, "at least one transfer must be specified")
	}

	if len(msg.Transfers) > 100 { // Limite raisonnable pour éviter le spam
		return errors.Wrap(ErrInvalidRequest, "cannot transfer more than 100 capsules in a single transaction")
	}

	for i, transfer := range msg.Transfers {
		// Validate new owner address
		_, err := sdk.AccAddressFromBech32(transfer.NewOwner)
		if err != nil {
			return errors.Wrapf(ErrInvalidAddress, "invalid new owner address at index %d (%s)", i, err)
		}

		// Check that addresses are different
		if msg.CurrentOwner == transfer.NewOwner {
			return errors.Wrapf(ErrInvalidRequest, "current and new owner must be different at index %d", i)
		}

		// Validate capsule ID
		if transfer.CapsuleID == 0 {
			return errors.Wrapf(ErrCapsuleNotFound, "capsule ID cannot be zero at index %d", i)
		}
	}

	// Validate transfer fee if provided
	if !msg.TransferFee.IsZero() && !msg.TransferFee.IsValid() {
		return errors.Wrap(ErrInvalidCoins, "invalid transfer fee")
	}

	return nil
}

// MsgApproveTransfer defines the message to approve a pending transfer
type MsgApproveTransfer struct {
	Approver     string `json:"approver"`
	TransferID   string `json:"transfer_id"`
	CapsuleID    uint64 `json:"capsule_id"`
	Approved     bool   `json:"approved"`
}

// NewMsgApproveTransfer creates a new MsgApproveTransfer
func NewMsgApproveTransfer(approver string, transferID string, capsuleID uint64, approved bool) *MsgApproveTransfer {
	return &MsgApproveTransfer{
		Approver:   approver,
		TransferID: transferID,
		CapsuleID:  capsuleID,
		Approved:   approved,
	}
}

// Route implements the sdk.Msg interface
func (msg *MsgApproveTransfer) Route() string {
	return RouterKey
}

// Type implements the sdk.Msg interface
func (msg *MsgApproveTransfer) Type() string {
	return TypeMsgApproveTransfer
}

// GetSigners implements the sdk.Msg interface
func (msg *MsgApproveTransfer) GetSigners() []sdk.AccAddress {
	approver, err := sdk.AccAddressFromBech32(msg.Approver)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{approver}
}

// GetSignBytes implements the sdk.Msg interface
func (msg *MsgApproveTransfer) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// ValidateBasic implements the sdk.Msg interface
func (msg *MsgApproveTransfer) ValidateBasic() error {
	// Validate approver address
	_, err := sdk.AccAddressFromBech32(msg.Approver)
	if err != nil {
		return errors.Wrapf(ErrInvalidAddress, "invalid approver address (%s)", err)
	}

	// Validate transfer ID
	if msg.TransferID == "" {
		return errors.Wrap(ErrInvalidRequest, "transfer ID cannot be empty")
	}

	// Validate capsule ID
	if msg.CapsuleID == 0 {
		return errors.Wrap(ErrCapsuleNotFound, "capsule ID cannot be zero")
	}

	return nil
}

// MsgEmergencyDeleteContract defines the message to emergency delete a capsule's smart contract
type MsgEmergencyDeleteContract struct {
	Creator        string `json:"creator"`
	CapsuleID      uint64 `json:"capsule_id"`
	EmergencyReason string `json:"emergency_reason"`
	ConfirmationCode string `json:"confirmation_code"` // Additional security layer
}

// NewMsgEmergencyDeleteContract creates a new MsgEmergencyDeleteContract
func NewMsgEmergencyDeleteContract(creator string, capsuleID uint64, reason string, confirmationCode string) *MsgEmergencyDeleteContract {
	return &MsgEmergencyDeleteContract{
		Creator:          creator,
		CapsuleID:        capsuleID,
		EmergencyReason:  reason,
		ConfirmationCode: confirmationCode,
	}
}

// Route implements the sdk.Msg interface
func (msg *MsgEmergencyDeleteContract) Route() string {
	return RouterKey
}

// Type implements the sdk.Msg interface
func (msg *MsgEmergencyDeleteContract) Type() string {
	return TypeMsgEmergencyDeleteContract
}

// GetSigners implements the sdk.Msg interface
func (msg *MsgEmergencyDeleteContract) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

// GetSignBytes implements the sdk.Msg interface
func (msg *MsgEmergencyDeleteContract) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// ValidateBasic implements the sdk.Msg interface
func (msg *MsgEmergencyDeleteContract) ValidateBasic() error {
	// Validate creator address
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return errors.Wrapf(ErrInvalidAddress, "invalid creator address (%s)", err)
	}

	// Validate capsule ID
	if msg.CapsuleID == 0 {
		return errors.Wrap(ErrCapsuleNotFound, "capsule ID cannot be zero")
	}

	// Validate emergency reason
	if msg.EmergencyReason == "" {
		return errors.Wrap(ErrInvalidRequest, "emergency reason cannot be empty")
	}

	if len(msg.EmergencyReason) > 500 {
		return errors.Wrap(ErrInvalidRequest, "emergency reason too long (max 500 characters)")
	}

	// Validate confirmation code
	if msg.ConfirmationCode == "" {
		return errors.Wrap(ErrInvalidRequest, "confirmation code cannot be empty")
	}

	// Simple validation for confirmation code format
	if len(msg.ConfirmationCode) < 8 {
		return errors.Wrap(ErrInvalidRequest, "confirmation code must be at least 8 characters")
	}

	return nil
}