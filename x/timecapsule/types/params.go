package types

import (
	"fmt"
	"time"
	"cosmossdk.io/math"
	
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Parameter keys
var (
	KeyMaxDataSize          = []byte("MaxDataSize")
	KeyMaxCapsuleDuration   = []byte("MaxCapsuleDuration")
	KeyMinThreshold         = []byte("MinThreshold")
	KeyMaxShares            = []byte("MaxShares")
	KeyCreationFee          = []byte("CreationFee")
	KeyMaintenanceFee       = []byte("MaintenanceFee")
	KeyMinInactivityPeriod  = []byte("MinInactivityPeriod")
	KeyMaxInactivityPeriod  = []byte("MaxInactivityPeriod")
	KeyAllowedCapsuleTypes  = []byte("AllowedCapsuleTypes")
	KeyMasterNodeMinStake   = []byte("MasterNodeMinStake")
)

// Default parameter values
const (
	DefaultMaxDataSize         = 1024 * 1024 // 1MB
	DefaultMaxCapsuleDuration  = 365 * 24 * time.Hour // 1 year
	DefaultMinThreshold        = uint32(2)
	DefaultMaxShares           = uint32(10)
	DefaultMinInactivityPeriod = uint64(30 * 24 * 60 * 60) // 30 days in seconds
	DefaultMaxInactivityPeriod = uint64(365 * 24 * 60 * 60) // 365 days in seconds
	DefaultMasterNodeMinStake  = "1000000" // 1 million base units
)

// Default creation and maintenance fees
var (
	DefaultCreationFee    = sdk.NewCoins(sdk.NewCoin("stake", math.NewInt(100000)))   // 0.1 stake
	DefaultMaintenanceFee = sdk.NewCoins(sdk.NewCoin("stake", math.NewInt(10000)))    // 0.01 stake per day
)

// Default allowed capsule types
var DefaultAllowedCapsuleTypes = []CapsuleType{
	CapsuleType_SAFE,
	CapsuleType_TIME_LOCK,
	CapsuleType_CONDITIONAL,
	CapsuleType_MULTI_SIG,
	CapsuleType_DEAD_MANS_SWITCH,
}

// Params defines the parameters for the time capsule module
type Params struct {
	MaxDataSize          uint64        `json:"max_data_size"`
	MaxCapsuleDuration   time.Duration `json:"max_capsule_duration"`
	MinThreshold         uint32        `json:"min_threshold"`
	MaxShares            uint32        `json:"max_shares"`
	CreationFee          sdk.Coins     `json:"creation_fee"`
	MaintenanceFee       sdk.Coins     `json:"maintenance_fee"`
	MinInactivityPeriod  uint64        `json:"min_inactivity_period"`
	MaxInactivityPeriod  uint64        `json:"max_inactivity_period"`
	AllowedCapsuleTypes  []CapsuleType `json:"allowed_capsule_types"`
	MasterNodeMinStake   math.Int      `json:"master_node_min_stake"`
}

// NewParams creates a new Params object
func NewParams(
	maxDataSize uint64,
	maxCapsuleDuration time.Duration,
	minThreshold uint32,
	maxShares uint32,
	creationFee sdk.Coins,
	maintenanceFee sdk.Coins,
	minInactivityPeriod uint64,
	maxInactivityPeriod uint64,
	allowedCapsuleTypes []CapsuleType,
	masterNodeMinStake math.Int,
) Params {
	return Params{
		MaxDataSize:         maxDataSize,
		MaxCapsuleDuration:  maxCapsuleDuration,
		MinThreshold:        minThreshold,
		MaxShares:           maxShares,
		CreationFee:         creationFee,
		MaintenanceFee:      maintenanceFee,
		MinInactivityPeriod: minInactivityPeriod,
		MaxInactivityPeriod: maxInactivityPeriod,
		AllowedCapsuleTypes: allowedCapsuleTypes,
		MasterNodeMinStake:  masterNodeMinStake,
	}
}

// DefaultParams returns the default parameters for the time capsule module
func DefaultParams() Params {
	return NewParams(
		DefaultMaxDataSize,
		DefaultMaxCapsuleDuration,
		DefaultMinThreshold,
		DefaultMaxShares,
		DefaultCreationFee,
		DefaultMaintenanceFee,
		DefaultMinInactivityPeriod,
		DefaultMaxInactivityPeriod,
		DefaultAllowedCapsuleTypes,
		math.MustNewIntFromString(DefaultMasterNodeMinStake),
	)
}


// Validate validates the set of params
func (p Params) Validate() error {
	if err := validateMaxDataSize(p.MaxDataSize); err != nil {
		return err
	}
	if err := validateMaxCapsuleDuration(p.MaxCapsuleDuration); err != nil {
		return err
	}
	if err := validateMinThreshold(p.MinThreshold); err != nil {
		return err
	}
	if err := validateMaxShares(p.MaxShares); err != nil {
		return err
	}
	if err := validateCreationFee(p.CreationFee); err != nil {
		return err
	}
	if err := validateMaintenanceFee(p.MaintenanceFee); err != nil {
		return err
	}
	if err := validateMinInactivityPeriod(p.MinInactivityPeriod); err != nil {
		return err
	}
	if err := validateMaxInactivityPeriod(p.MaxInactivityPeriod); err != nil {
		return err
	}
	if err := validateAllowedCapsuleTypes(p.AllowedCapsuleTypes); err != nil {
		return err
	}
	if err := validateMasterNodeMinStake(p.MasterNodeMinStake); err != nil {
		return err
	}
	
	// Cross-field validation
	if p.MinThreshold > p.MaxShares {
		return fmt.Errorf("min threshold (%d) cannot be greater than max shares (%d)", p.MinThreshold, p.MaxShares)
	}
	
	if p.MinInactivityPeriod > p.MaxInactivityPeriod {
		return fmt.Errorf("min inactivity period (%d) cannot be greater than max inactivity period (%d)", 
			p.MinInactivityPeriod, p.MaxInactivityPeriod)
	}
	
	return nil
}


// Validation functions

func validateMaxDataSize(i interface{}) error {
	v, ok := i.(uint64)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	
	if v == 0 {
		return fmt.Errorf("max data size must be positive")
	}
	
	// Max 10MB for practical reasons
	if v > 10*1024*1024 {
		return fmt.Errorf("max data size cannot exceed 10MB")
	}
	
	return nil
}

func validateMaxCapsuleDuration(i interface{}) error {
	v, ok := i.(time.Duration)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	
	if v <= 0 {
		return fmt.Errorf("max capsule duration must be positive")
	}
	
	// Max 100 years for practical reasons
	if v > 100*365*24*time.Hour {
		return fmt.Errorf("max capsule duration cannot exceed 100 years")
	}
	
	return nil
}

func validateMinThreshold(i interface{}) error {
	v, ok := i.(uint32)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	
	if v == 0 {
		return fmt.Errorf("min threshold must be positive")
	}
	
	return nil
}

func validateMaxShares(i interface{}) error {
	v, ok := i.(uint32)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	
	if v == 0 {
		return fmt.Errorf("max shares must be positive")
	}
	
	// Reasonable upper limit
	if v > 100 {
		return fmt.Errorf("max shares cannot exceed 100")
	}
	
	return nil
}

func validateCreationFee(i interface{}) error {
	v, ok := i.(sdk.Coins)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	
	return v.Validate()
}

func validateMaintenanceFee(i interface{}) error {
	v, ok := i.(sdk.Coins)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	
	return v.Validate()
}

func validateMinInactivityPeriod(i interface{}) error {
	v, ok := i.(uint64)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	
	// Minimum 1 hour
	if v < 3600 {
		return fmt.Errorf("min inactivity period cannot be less than 1 hour")
	}
	
	return nil
}

func validateMaxInactivityPeriod(i interface{}) error {
	v, ok := i.(uint64)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	
	// Maximum 10 years
	if v > 10*365*24*3600 {
		return fmt.Errorf("max inactivity period cannot exceed 10 years")
	}
	
	return nil
}

func validateAllowedCapsuleTypes(i interface{}) error {
	v, ok := i.([]CapsuleType)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	
	if len(v) == 0 {
		return fmt.Errorf("allowed capsule types cannot be empty")
	}
	
	// Check for duplicates
	seen := make(map[CapsuleType]bool)
	for _, capsuleType := range v {
		if seen[capsuleType] {
			return fmt.Errorf("duplicate capsule type: %v", capsuleType)
		}
		seen[capsuleType] = true
		
		// Validate each type
		switch capsuleType {
		case CapsuleType_SAFE, CapsuleType_TIME_LOCK, CapsuleType_CONDITIONAL, 
		     CapsuleType_MULTI_SIG, CapsuleType_DEAD_MANS_SWITCH:
			// Valid types
		default:
			return fmt.Errorf("invalid capsule type: %v", capsuleType)
		}
	}
	
	return nil
}

func validateMasterNodeMinStake(i interface{}) error {
	v, ok := i.(math.Int)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	
	if v.IsNegative() {
		return fmt.Errorf("master node min stake cannot be negative")
	}
	
	return nil
}