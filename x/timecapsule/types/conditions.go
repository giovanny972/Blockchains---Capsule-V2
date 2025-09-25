package types

import (
	"encoding/json"
	"fmt"
	"time"
	
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// ConditionType defines the type of access condition
type ConditionType string

const (
	ConditionType_TIME        ConditionType = "time"
	ConditionType_ORACLE      ConditionType = "oracle"
	ConditionType_MULTISIG    ConditionType = "multisig"
	ConditionType_COMPOSITE   ConditionType = "composite"
	ConditionType_EXTERNAL    ConditionType = "external"
	ConditionType_DEATH       ConditionType = "death"
	ConditionType_INACTIVITY  ConditionType = "inactivity"
)

// Condition represents a general access condition interface
type Condition interface {
	GetType() ConditionType
	Validate() error
	Evaluate(ctx sdk.Context, params map[string]interface{}) (bool, error)
	GetMetadata() map[string]interface{}
}

// TimeCondition represents a time-based access condition
type TimeCondition struct {
	UnlockTime time.Time `json:"unlock_time"`
	Timezone   string    `json:"timezone,omitempty"`
}

func (tc *TimeCondition) GetType() ConditionType {
	return ConditionType_TIME
}

func (tc *TimeCondition) Validate() error {
	if tc.UnlockTime.IsZero() {
		return fmt.Errorf("unlock time cannot be zero")
	}
	if tc.UnlockTime.Before(time.Now()) {
		return fmt.Errorf("unlock time must be in the future")
	}
	return nil
}

func (tc *TimeCondition) Evaluate(ctx sdk.Context, params map[string]interface{}) (bool, error) {
	return ctx.BlockTime().After(tc.UnlockTime), nil
}

func (tc *TimeCondition) GetMetadata() map[string]interface{} {
	return map[string]interface{}{
		"unlock_time": tc.UnlockTime,
		"timezone":    tc.Timezone,
	}
}

// MultiSigCondition represents a multi-signature access condition
type MultiSigCondition struct {
	RequiredSignatures uint32   `json:"required_signatures"`
	Signers           []string `json:"signers"`
	Signatures        []string `json:"signatures,omitempty"`
}

func (msc *MultiSigCondition) GetType() ConditionType {
	return ConditionType_MULTISIG
}

func (msc *MultiSigCondition) Validate() error {
	if msc.RequiredSignatures == 0 {
		return fmt.Errorf("required signatures must be greater than 0")
	}
	if len(msc.Signers) == 0 {
		return fmt.Errorf("signers list cannot be empty")
	}
	if msc.RequiredSignatures > uint32(len(msc.Signers)) {
		return fmt.Errorf("required signatures cannot exceed number of signers")
	}
	
	// Validate signer addresses
	for i, signer := range msc.Signers {
		if _, err := sdk.AccAddressFromBech32(signer); err != nil {
			return fmt.Errorf("invalid signer address at index %d: %w", i, err)
		}
	}
	
	return nil
}

func (msc *MultiSigCondition) Evaluate(ctx sdk.Context, params map[string]interface{}) (bool, error) {
	// Check if enough valid signatures are provided
	validSigCount := uint32(0)
	
	signatures, ok := params["signatures"].([]string)
	if !ok {
		signatures = msc.Signatures
	}
	
	// This is a simplified check - in practice, you'd verify cryptographic signatures
	for _, sig := range signatures {
		if sig != "" && msc.isValidSignature(sig) {
			validSigCount++
		}
	}
	
	return validSigCount >= msc.RequiredSignatures, nil
}

func (msc *MultiSigCondition) isValidSignature(signature string) bool {
	// Placeholder for signature verification logic
	// In practice, this would verify the cryptographic signature
	return len(signature) > 0
}

func (msc *MultiSigCondition) GetMetadata() map[string]interface{} {
	return map[string]interface{}{
		"required_signatures": msc.RequiredSignatures,
		"signers":            msc.Signers,
		"current_signatures": len(msc.Signatures),
	}
}

// OracleCondition represents an oracle-based access condition
type OracleCondition struct {
	OracleAddress string                 `json:"oracle_address"`
	Query         string                 `json:"query"`
	ExpectedValue interface{}            `json:"expected_value"`
	Operator      string                 `json:"operator"` // "eq", "gt", "lt", "gte", "lte", "ne"
	Metadata      map[string]interface{} `json:"metadata,omitempty"`
}

func (oc *OracleCondition) GetType() ConditionType {
	return ConditionType_ORACLE
}

func (oc *OracleCondition) Validate() error {
	if oc.OracleAddress == "" {
		return fmt.Errorf("oracle address cannot be empty")
	}
	if oc.Query == "" {
		return fmt.Errorf("oracle query cannot be empty")
	}
	if oc.ExpectedValue == nil {
		return fmt.Errorf("expected value cannot be nil")
	}
	
	validOperators := map[string]bool{
		"eq": true, "gt": true, "lt": true,
		"gte": true, "lte": true, "ne": true,
	}
	if !validOperators[oc.Operator] {
		return fmt.Errorf("invalid operator: %s", oc.Operator)
	}
	
	return nil
}

func (oc *OracleCondition) Evaluate(ctx sdk.Context, params map[string]interface{}) (bool, error) {
	// Placeholder for oracle query logic
	// In practice, this would query the specified oracle
	oracleValue, exists := params["oracle_value"]
	if !exists {
		return false, fmt.Errorf("oracle value not provided")
	}
	
	return oc.compareValues(oracleValue, oc.ExpectedValue, oc.Operator)
}

func (oc *OracleCondition) compareValues(actual, expected interface{}, operator string) (bool, error) {
	switch operator {
	case "eq":
		return actual == expected, nil
	case "ne":
		return actual != expected, nil
	// Add more comparison logic as needed
	default:
		return false, fmt.Errorf("unsupported operator: %s", operator)
	}
}

func (oc *OracleCondition) GetMetadata() map[string]interface{} {
	metadata := map[string]interface{}{
		"oracle_address": oc.OracleAddress,
		"query":          oc.Query,
		"expected_value": oc.ExpectedValue,
		"operator":       oc.Operator,
	}
	
	for k, v := range oc.Metadata {
		metadata[k] = v
	}
	
	return metadata
}

// CompositeCondition represents a combination of multiple conditions
type CompositeCondition struct {
	Operator   string      `json:"operator"` // "AND", "OR", "NOT"
	Conditions []Condition `json:"conditions"`
}

func (cc *CompositeCondition) GetType() ConditionType {
	return ConditionType_COMPOSITE
}

func (cc *CompositeCondition) Validate() error {
	validOperators := map[string]bool{"AND": true, "OR": true, "NOT": true}
	if !validOperators[cc.Operator] {
		return fmt.Errorf("invalid composite operator: %s", cc.Operator)
	}
	
	if len(cc.Conditions) == 0 {
		return fmt.Errorf("composite condition must have at least one sub-condition")
	}
	
	if cc.Operator == "NOT" && len(cc.Conditions) != 1 {
		return fmt.Errorf("NOT operator requires exactly one sub-condition")
	}
	
	// Validate all sub-conditions
	for i, condition := range cc.Conditions {
		if err := condition.Validate(); err != nil {
			return fmt.Errorf("invalid sub-condition at index %d: %w", i, err)
		}
	}
	
	return nil
}

func (cc *CompositeCondition) Evaluate(ctx sdk.Context, params map[string]interface{}) (bool, error) {
	switch cc.Operator {
	case "AND":
		for _, condition := range cc.Conditions {
			result, err := condition.Evaluate(ctx, params)
			if err != nil {
				return false, err
			}
			if !result {
				return false, nil
			}
		}
		return true, nil
		
	case "OR":
		for _, condition := range cc.Conditions {
			result, err := condition.Evaluate(ctx, params)
			if err != nil {
				continue // Skip failed conditions in OR
			}
			if result {
				return true, nil
			}
		}
		return false, nil
		
	case "NOT":
		if len(cc.Conditions) != 1 {
			return false, fmt.Errorf("NOT operator requires exactly one condition")
		}
		result, err := cc.Conditions[0].Evaluate(ctx, params)
		if err != nil {
			return false, err
		}
		return !result, nil
		
	default:
		return false, fmt.Errorf("unsupported composite operator: %s", cc.Operator)
	}
}

func (cc *CompositeCondition) GetMetadata() map[string]interface{} {
	subMetadata := make([]map[string]interface{}, len(cc.Conditions))
	for i, condition := range cc.Conditions {
		subMetadata[i] = condition.GetMetadata()
	}
	
	return map[string]interface{}{
		"operator":      cc.Operator,
		"conditions":    subMetadata,
		"condition_count": len(cc.Conditions),
	}
}

// InactivityCondition represents a dead man's switch condition
type InactivityCondition struct {
	InactivityPeriod uint64     `json:"inactivity_period"` // In seconds
	LastActivity     *time.Time `json:"last_activity,omitempty"`
	GracePeriod      uint64     `json:"grace_period,omitempty"` // Additional time buffer
}

func (ic *InactivityCondition) GetType() ConditionType {
	return ConditionType_INACTIVITY
}

func (ic *InactivityCondition) Validate() error {
	if ic.InactivityPeriod == 0 {
		return fmt.Errorf("inactivity period must be greater than 0")
	}
	return nil
}

func (ic *InactivityCondition) Evaluate(ctx sdk.Context, params map[string]interface{}) (bool, error) {
	if ic.LastActivity == nil {
		// If no activity recorded, use current time as baseline
		now := ctx.BlockTime()
		ic.LastActivity = &now
		return false, nil
	}
	
	totalPeriod := time.Duration(ic.InactivityPeriod+ic.GracePeriod) * time.Second
	threshold := ic.LastActivity.Add(totalPeriod)
	
	return ctx.BlockTime().After(threshold), nil
}

func (ic *InactivityCondition) GetMetadata() map[string]interface{} {
	metadata := map[string]interface{}{
		"inactivity_period": ic.InactivityPeriod,
		"grace_period":      ic.GracePeriod,
	}
	
	if ic.LastActivity != nil {
		metadata["last_activity"] = ic.LastActivity
		metadata["unlock_time"] = ic.LastActivity.Add(time.Duration(ic.InactivityPeriod+ic.GracePeriod) * time.Second)
	}
	
	return metadata
}

// ConditionFactory creates conditions from JSON data
type ConditionFactory struct{}

func NewConditionFactory() *ConditionFactory {
	return &ConditionFactory{}
}

func (cf *ConditionFactory) CreateCondition(conditionType ConditionType, data []byte) (Condition, error) {
	switch conditionType {
	case ConditionType_TIME:
		var condition TimeCondition
		if err := json.Unmarshal(data, &condition); err != nil {
			return nil, fmt.Errorf("failed to unmarshal time condition: %w", err)
		}
		return &condition, nil
		
	case ConditionType_MULTISIG:
		var condition MultiSigCondition
		if err := json.Unmarshal(data, &condition); err != nil {
			return nil, fmt.Errorf("failed to unmarshal multisig condition: %w", err)
		}
		return &condition, nil
		
	case ConditionType_ORACLE:
		var condition OracleCondition
		if err := json.Unmarshal(data, &condition); err != nil {
			return nil, fmt.Errorf("failed to unmarshal oracle condition: %w", err)
		}
		return &condition, nil
		
	case ConditionType_INACTIVITY:
		var condition InactivityCondition
		if err := json.Unmarshal(data, &condition); err != nil {
			return nil, fmt.Errorf("failed to unmarshal inactivity condition: %w", err)
		}
		return &condition, nil
		
	default:
		return nil, fmt.Errorf("unsupported condition type: %s", conditionType)
	}
}