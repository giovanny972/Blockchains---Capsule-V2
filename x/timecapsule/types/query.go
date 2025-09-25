package types

import "context"

// Query request and response types

// QueryParamsRequest is the request type for the Query/Params RPC method
type QueryParamsRequest struct{}

// QueryParamsResponse is the response type for the Query/Params RPC method
type QueryParamsResponse struct {
	Params Params `json:"params"`
}

// QueryCapsuleRequest is the request type for the Query/Capsule RPC method
type QueryCapsuleRequest struct {
	CapsuleId uint64 `json:"capsule_id"`
}

// QueryCapsuleResponse is the response type for the Query/Capsule RPC method
type QueryCapsuleResponse struct {
	Capsule *TimeCapsule `json:"capsule"`
}

// QueryCapsulesRequest is the request type for the Query/Capsules RPC method
type QueryCapsulesRequest struct {
	// Pagination parameters can be added here
	Limit  uint64 `json:"limit,omitempty"`
	Offset uint64 `json:"offset,omitempty"`
}

// QueryCapsulesResponse is the response type for the Query/Capsules RPC method
type QueryCapsulesResponse struct {
	Capsules []TimeCapsule `json:"capsules"`
}

// QueryUserCapsulesRequest is the request type for the Query/UserCapsules RPC method
type QueryUserCapsulesRequest struct {
	Owner string `json:"owner"`
}

// QueryUserCapsulesResponse is the response type for the Query/UserCapsules RPC method
type QueryUserCapsulesResponse struct {
	Capsules []*TimeCapsule `json:"capsules"`
	Owner    string         `json:"owner"`
}

// QueryCapsulesByTypeRequest is the request type for the Query/CapsulesByType RPC method
type QueryCapsulesByTypeRequest struct {
	CapsuleType CapsuleType `json:"capsule_type"`
}

// QueryCapsulesByTypeResponse is the response type for the Query/CapsulesByType RPC method
type QueryCapsulesByTypeResponse struct {
	Capsules    []TimeCapsule `json:"capsules"`
	CapsuleType CapsuleType   `json:"capsule_type"`
}

// QueryCapsulesByStatusRequest is the request type for the Query/CapsulesByStatus RPC method
type QueryCapsulesByStatusRequest struct {
	Status CapsuleStatus `json:"status"`
}

// QueryCapsulesByStatusResponse is the response type for the Query/CapsulesByStatus RPC method
type QueryCapsulesByStatusResponse struct {
	Capsules []TimeCapsule `json:"capsules"`
	Status   CapsuleStatus `json:"status"`
}

// QueryStatsRequest is the request type for the Query/Stats RPC method
type QueryStatsRequest struct{}

// ModuleStats represents statistics about the time capsule module
type ModuleStats struct {
	TotalCapsules     uint64 `json:"total_capsules"`
	ActiveCapsules    uint64 `json:"active_capsules"`
	OpenedCapsules    uint64 `json:"opened_capsules"`
	ExpiredCapsules   uint64 `json:"expired_capsules"`
	CancelledCapsules uint64 `json:"cancelled_capsules"`
	TotalDataSize     uint64 `json:"total_data_size,omitempty"`
	TotalKeyShares    uint64 `json:"total_key_shares,omitempty"`
}

// QueryStatsResponse is the response type for the Query/Stats RPC method
type QueryStatsResponse struct {
	Stats *ModuleStats `json:"stats"`
}

// QueryKeySharesRequest is the request type for the Query/KeyShares RPC method
type QueryKeySharesRequest struct {
	CapsuleId uint64 `json:"capsule_id"`
}

// QueryKeySharesResponse is the response type for the Query/KeyShares RPC method
type QueryKeySharesResponse struct {
	KeyShares []KeyShare `json:"key_shares"`
	CapsuleId uint64     `json:"capsule_id"`
}

// QueryConditionContractRequest is the request type for the Query/ConditionContract RPC method
type QueryConditionContractRequest struct {
	Address string `json:"address"`
}

// QueryConditionContractResponse is the response type for the Query/ConditionContract RPC method
type QueryConditionContractResponse struct {
	Contract *ConditionContract `json:"contract"`
}

// QueryConditionContractsRequest is the request type for the Query/ConditionContracts RPC method
type QueryConditionContractsRequest struct{}

// QueryConditionContractsResponse is the response type for the Query/ConditionContracts RPC method
type QueryConditionContractsResponse struct {
	Contracts []ConditionContract `json:"contracts"`
}

// Message response types

// MsgCreateCapsuleResponse is the response type for MsgCreateCapsule
type MsgCreateCapsuleResponse struct {
	CapsuleId uint64 `json:"capsule_id"`
}

// MsgOpenCapsuleResponse is the response type for MsgOpenCapsule  
type MsgOpenCapsuleResponse struct {
	Data []byte `json:"data"`
}

// MsgUpdateActivityResponse is the response type for MsgUpdateActivity
type MsgUpdateActivityResponse struct{}

// MsgCancelCapsuleResponse is the response type for MsgCancelCapsule
type MsgCancelCapsuleResponse struct{}

// MsgTransferCapsuleResponse is the response type for MsgTransferCapsule
type MsgTransferCapsuleResponse struct{}

// MsgBatchTransferCapsulesResponse is the response type for MsgBatchTransferCapsules
type MsgBatchTransferCapsulesResponse struct {
	TransferredCapsules []uint64        `json:"transferred_capsules"`
	FailedTransfers     []FailedTransfer `json:"failed_transfers"`
}

// FailedTransfer represents a failed transfer in a batch
type FailedTransfer struct {
	CapsuleID uint64 `json:"capsule_id"`
	Reason    string `json:"reason"`
}

// MsgApproveTransferResponse is the response type for MsgApproveTransfer
type MsgApproveTransferResponse struct {
	Approved bool `json:"approved"`
}

// MsgEmergencyDeleteContractResponse is the response type for MsgEmergencyDeleteContract
type MsgEmergencyDeleteContractResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// Interface definitions for gRPC services

// QueryServer defines the gRPC querier service
type QueryServer interface {
	// Params returns the module parameters
	Params(context.Context, *QueryParamsRequest) (*QueryParamsResponse, error)
	
	// Capsule returns details of a specific capsule
	Capsule(context.Context, *QueryCapsuleRequest) (*QueryCapsuleResponse, error)
	
	// Capsules returns a list of capsules
	Capsules(context.Context, *QueryCapsulesRequest) (*QueryCapsulesResponse, error)
	
	// UserCapsules returns all capsules owned by a user
	UserCapsules(context.Context, *QueryUserCapsulesRequest) (*QueryUserCapsulesResponse, error)
	
	// CapsulesByType returns capsules filtered by type
	CapsulesByType(context.Context, *QueryCapsulesByTypeRequest) (*QueryCapsulesByTypeResponse, error)
	
	// CapsulesByStatus returns capsules filtered by status
	CapsulesByStatus(context.Context, *QueryCapsulesByStatusRequest) (*QueryCapsulesByStatusResponse, error)
	
	// Stats returns module statistics
	Stats(context.Context, *QueryStatsRequest) (*QueryStatsResponse, error)
	
	// KeyShares returns key shares for a capsule
	KeyShares(context.Context, *QueryKeySharesRequest) (*QueryKeySharesResponse, error)
	
	// ConditionContract returns a condition contract
	ConditionContract(context.Context, *QueryConditionContractRequest) (*QueryConditionContractResponse, error)
	
	// ConditionContracts returns all condition contracts
	ConditionContracts(context.Context, *QueryConditionContractsRequest) (*QueryConditionContractsResponse, error)
}

// MsgServer defines the gRPC message service
type MsgServer interface {
	// CreateCapsule creates a new time capsule
	CreateCapsule(context.Context, *MsgCreateCapsule) (*MsgCreateCapsuleResponse, error)
	
	// OpenCapsule opens a time capsule
	OpenCapsule(context.Context, *MsgOpenCapsule) (*MsgOpenCapsuleResponse, error)
	
	// UpdateActivity updates last activity for dead man's switch
	UpdateActivity(context.Context, *MsgUpdateActivity) (*MsgUpdateActivityResponse, error)
	
	// CancelCapsule cancels a time capsule
	CancelCapsule(context.Context, *MsgCancelCapsule) (*MsgCancelCapsuleResponse, error)
	
	// TransferCapsule transfers capsule ownership
	TransferCapsule(context.Context, *MsgTransferCapsule) (*MsgTransferCapsuleResponse, error)
}