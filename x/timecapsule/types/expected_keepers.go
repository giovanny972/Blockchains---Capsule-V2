package types

import (
	"context"
	"cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/types"
)

// AccountKeeper defines the expected account keeper used for simulations (noalias)
type AccountKeeper interface {
	GetAccount(ctx context.Context, addr sdk.AccAddress) sdk.AccountI
	HasAccount(ctx context.Context, addr sdk.AccAddress) bool
	SetAccount(ctx context.Context, acc sdk.AccountI)
	NewAccountWithAddress(ctx context.Context, addr sdk.AccAddress) sdk.AccountI
	
	GetModuleAddress(name string) sdk.AccAddress
	GetModuleAccount(ctx context.Context, name string) types.ModuleAccountI
	SetModuleAccount(ctx context.Context, macc types.ModuleAccountI)
}

// BankKeeper defines the expected interface needed to retrieve account balances.
type BankKeeper interface {
	SpendableCoins(ctx context.Context, addr sdk.AccAddress) sdk.Coins
	GetBalance(ctx context.Context, addr sdk.AccAddress, denom string) sdk.Coin
	GetAllBalances(ctx context.Context, addr sdk.AccAddress) sdk.Coins
	
	SendCoins(ctx context.Context, fromAddr sdk.AccAddress, toAddr sdk.AccAddress, amt sdk.Coins) error
	SendCoinsFromModuleToAccount(ctx context.Context, senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coins) error
	SendCoinsFromAccountToModule(ctx context.Context, senderAddr sdk.AccAddress, recipientModule string, amt sdk.Coins) error
	SendCoinsFromModuleToModule(ctx context.Context, senderModule, recipientModule string, amt sdk.Coins) error
	
	MintCoins(ctx context.Context, name string, amt sdk.Coins) error
	BurnCoins(ctx context.Context, name string, amt sdk.Coins) error
}

// StakingKeeper expected staking keeper (noalias)
type StakingKeeper interface {
	// Methods to get validators and their information
	GetValidator(ctx context.Context, addr sdk.ValAddress) (validator StakingValidator, found bool)
	GetAllValidators(ctx context.Context) (validators []StakingValidator)
	GetBondedValidatorsByPower(ctx context.Context) []StakingValidator
}

// StakingValidator is a subset of the staking Validator type
type StakingValidator interface {
	GetOperator() sdk.ValAddress
	GetConsPubKey() (sdk.ConsAddress, error)
	IsJailed() bool
	IsBonded() bool
	GetStatus() BondStatus
	GetTokens() math.Int
	GetDelegatorShares() math.LegacyDec
	GetMoniker() string
}

// BondStatus is the status of a validator
type BondStatus int32

const (
	// BondStatusUnbonded defines a validator that is not bonded
	BondStatusUnbonded BondStatus = 1
	// BondStatusUnbonding defines a validator that is unbonding
	BondStatusUnbonding BondStatus = 2
	// BondStatusBonded defines a validator that is bonded
	BondStatusBonded BondStatus = 3
)

// GovKeeper expected governance keeper (for condition contracts)
type GovKeeper interface {
	GetProposal(ctx context.Context, proposalID uint64) (Proposal, bool)
	GetProposals(ctx context.Context) (proposals []Proposal)
}

// Proposal defines a governance proposal
type Proposal interface {
	GetId() uint64
	GetStatus() ProposalStatus
	GetFinalTallyResult() TallyResult
}

// ProposalStatus defines the status of a governance proposal
type ProposalStatus int32

const (
	StatusNil           ProposalStatus = 0
	StatusDepositPeriod ProposalStatus = 1
	StatusVotingPeriod  ProposalStatus = 2
	StatusPassed        ProposalStatus = 3
	StatusRejected      ProposalStatus = 4
	StatusFailed        ProposalStatus = 5
)

// TallyResult defines a standard tally for a governance proposal
type TallyResult struct {
	YesCount        math.Int
	AbstainCount    math.Int
	NoCount         math.Int
	NoWithVetoCount math.Int
}

// OracleKeeper defines expected oracle keeper interface
type OracleKeeper interface {
	GetPrice(ctx context.Context, symbol string) (math.LegacyDec, error)
	GetData(ctx context.Context, key string) ([]byte, error)
	IsDataAvailable(ctx context.Context, key string) bool
}

// DistributionKeeper expected distribution keeper
type DistributionKeeper interface {
	FundCommunityPool(ctx context.Context, amount sdk.Coins, sender sdk.AccAddress) error
	DistributeFromFeePool(ctx context.Context, amount sdk.Coins, receiveAddr sdk.AccAddress) error
}