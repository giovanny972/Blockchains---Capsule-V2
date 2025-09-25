package types

import (
	"fmt"
	
	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
)

// RegisterCodec registers the necessary x/timecapsule interfaces and concrete types
// on the provided Amino codec. These types are used for Amino JSON serialization.
func RegisterCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(&MsgCreateCapsule{}, "timecapsule/MsgCreateCapsule", nil)
	cdc.RegisterConcrete(&MsgOpenCapsule{}, "timecapsule/MsgOpenCapsule", nil)
	cdc.RegisterConcrete(&MsgUpdateActivity{}, "timecapsule/MsgUpdateActivity", nil)
	cdc.RegisterConcrete(&MsgCancelCapsule{}, "timecapsule/MsgCancelCapsule", nil)
	cdc.RegisterConcrete(&MsgTransferCapsule{}, "timecapsule/MsgTransferCapsule", nil)
}

// RegisterInterfaces registers the x/timecapsule interfaces types with the
// interface registry.
func RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgCreateCapsule{},
		&MsgOpenCapsule{},
		&MsgUpdateActivity{},
		&MsgCancelCapsule{},
		&MsgTransferCapsule{},
	)

	// msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc) // TODO: implement when protobuf is generated
}

var (
	// Amino is the legacy amino codec
	Amino = codec.NewLegacyAmino()
	// ModuleCdc is the module codec
	ModuleCdc = codec.NewProtoCodec(cdctypes.NewInterfaceRegistry())
)

func init() {
	RegisterCodec(Amino)
	Amino.Seal()
}

// Stub functions for gRPC service registration (to be replaced when protobuf is generated)
func RegisterMsgServer(server interface{}, impl interface{}) {
	// TODO: implement when protobuf is generated
}

func RegisterQueryServer(server interface{}, impl interface{}) {
	// TODO: implement when protobuf is generated
}

// NewQueryClient creates a new query client stub
func NewQueryClient(clientCtx interface{}) QueryClient {
	return &queryClient{}
}

// QueryClient interface stub
type QueryClient interface {
	Params(ctx interface{}, req *QueryParamsRequest) (*QueryParamsResponse, error)
	Capsule(ctx interface{}, req *QueryCapsuleRequest) (*QueryCapsuleResponse, error)
	Capsules(ctx interface{}, req *QueryCapsulesRequest) (*QueryCapsulesResponse, error)
	UserCapsules(ctx interface{}, req *QueryUserCapsulesRequest) (*QueryUserCapsulesResponse, error)
	CapsulesByType(ctx interface{}, req *QueryCapsulesByTypeRequest) (*QueryCapsulesByTypeResponse, error)
	CapsulesByStatus(ctx interface{}, req *QueryCapsulesByStatusRequest) (*QueryCapsulesByStatusResponse, error)
	Stats(ctx interface{}, req *QueryStatsRequest) (*QueryStatsResponse, error)
	KeyShares(ctx interface{}, req *QueryKeySharesRequest) (*QueryKeySharesResponse, error)
	ConditionContract(ctx interface{}, req *QueryConditionContractRequest) (*QueryConditionContractResponse, error)
	ConditionContracts(ctx interface{}, req *QueryConditionContractsRequest) (*QueryConditionContractsResponse, error)
}

// queryClient stub implementation
type queryClient struct{}

func (q *queryClient) Params(ctx interface{}, req *QueryParamsRequest) (*QueryParamsResponse, error) {
	return nil, fmt.Errorf("query client not implemented - please use protobuf generation")
}

func (q *queryClient) Capsule(ctx interface{}, req *QueryCapsuleRequest) (*QueryCapsuleResponse, error) {
	return nil, fmt.Errorf("query client not implemented - please use protobuf generation")
}

func (q *queryClient) Capsules(ctx interface{}, req *QueryCapsulesRequest) (*QueryCapsulesResponse, error) {
	return nil, fmt.Errorf("query client not implemented - please use protobuf generation")
}

func (q *queryClient) UserCapsules(ctx interface{}, req *QueryUserCapsulesRequest) (*QueryUserCapsulesResponse, error) {
	return nil, fmt.Errorf("query client not implemented - please use protobuf generation")
}

func (q *queryClient) CapsulesByType(ctx interface{}, req *QueryCapsulesByTypeRequest) (*QueryCapsulesByTypeResponse, error) {
	return nil, fmt.Errorf("query client not implemented - please use protobuf generation")
}

func (q *queryClient) CapsulesByStatus(ctx interface{}, req *QueryCapsulesByStatusRequest) (*QueryCapsulesByStatusResponse, error) {
	return nil, fmt.Errorf("query client not implemented - please use protobuf generation")
}

func (q *queryClient) Stats(ctx interface{}, req *QueryStatsRequest) (*QueryStatsResponse, error) {
	return nil, fmt.Errorf("query client not implemented - please use protobuf generation")
}

func (q *queryClient) KeyShares(ctx interface{}, req *QueryKeySharesRequest) (*QueryKeySharesResponse, error) {
	return nil, fmt.Errorf("query client not implemented - please use protobuf generation")
}

func (q *queryClient) ConditionContract(ctx interface{}, req *QueryConditionContractRequest) (*QueryConditionContractResponse, error) {
	return nil, fmt.Errorf("query client not implemented - please use protobuf generation")
}

func (q *queryClient) ConditionContracts(ctx interface{}, req *QueryConditionContractsRequest) (*QueryConditionContractsResponse, error) {
	return nil, fmt.Errorf("query client not implemented - please use protobuf generation")
}