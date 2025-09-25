package cli

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/spf13/cobra"

	"cosmossdk.io/core/address"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	
	"github.com/cosmos/cosmos-sdk/x/timecapsule/types"
)

// GetTxCmd returns the transaction commands for this module
func GetTxCmd(ac address.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("%s transactions subcommands", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(
		CmdCreateCapsule(ac),
		CmdOpenCapsule(ac),
		CmdUpdateActivity(ac),
		CmdCancelCapsule(ac),
		CmdTransferCapsule(ac),
		CmdBatchTransferCapsules(ac),
		CmdApproveTransfer(ac),
	)

	return cmd
}

// CmdCreateCapsule returns a CLI command for creating a time capsule
func CmdCreateCapsule(ac address.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-capsule [data-file] [capsule-type] [threshold] [total-shares]",
		Short: "Create a new time capsule",
		Long: `Create a new time capsule with encrypted data.

Capsule types:
- safe: Basic secure storage
- time_lock: Unlocks at specific time
- conditional: Unlocks based on conditions
- multi_sig: Requires multiple signatures
- dead_mans_switch: Unlocks after inactivity period

Example:
$ simd tx timecapsule create-capsule ./data.json time_lock 2 3 \
  --unlock-time="2025-12-31T23:59:59Z" \
  --recipient="cosmos1..." \
  --from=alice`,
		Args: cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			// Parse arguments
			dataFile := args[0]
			capsuleTypeStr := args[1]
			thresholdStr := args[2]
			totalSharesStr := args[3]

			// Read data from file
			data, err := readDataFile(dataFile)
			if err != nil {
				return fmt.Errorf("failed to read data file: %w", err)
			}

			// Parse capsule type
			capsuleType, err := parseCapsuleType(capsuleTypeStr)
			if err != nil {
				return err
			}

			// Parse threshold and shares
			threshold, err := strconv.ParseUint(thresholdStr, 10, 32)
			if err != nil {
				return fmt.Errorf("invalid threshold: %w", err)
			}

			totalShares, err := strconv.ParseUint(totalSharesStr, 10, 32)
			if err != nil {
				return fmt.Errorf("invalid total shares: %w", err)
			}

			// Get flags
			recipient, _ := cmd.Flags().GetString("recipient")
			unlockTimeStr, _ := cmd.Flags().GetString("unlock-time")
			conditionContract, _ := cmd.Flags().GetString("condition-contract")
			requiredSigs, _ := cmd.Flags().GetUint32("required-sigs")
			inactivityPeriod, _ := cmd.Flags().GetUint64("inactivity-period")
			title, _ := cmd.Flags().GetString("title")
			description, _ := cmd.Flags().GetString("description")

			// Parse unlock time if provided
			var unlockTime *time.Time
			if unlockTimeStr != "" {
				parsedTime, err := time.Parse(time.RFC3339, unlockTimeStr)
				if err != nil {
					return fmt.Errorf("invalid unlock time format (use RFC3339): %w", err)
				}
				unlockTime = &parsedTime
			}

			// Create message
			msg := &types.MsgCreateCapsule{
				Creator:           clientCtx.GetFromAddress().String(),
				Recipient:         recipient,
				Data:              data,
				CapsuleType:       capsuleType,
				Threshold:         uint32(threshold),
				TotalShares:       uint32(totalShares),
				UnlockTime:        unlockTime,
				ConditionContract: conditionContract,
				RequiredSigs:      requiredSigs,
				InactivityPeriod:  inactivityPeriod,
				Title:             title,
				Description:       description,
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	cmd.Flags().String("recipient", "", "Recipient address for the capsule")
	cmd.Flags().String("unlock-time", "", "Unlock time in RFC3339 format (e.g., 2025-12-31T23:59:59Z)")
	cmd.Flags().String("condition-contract", "", "Address of the condition contract")
	cmd.Flags().Uint32("required-sigs", 0, "Required signatures for multi-sig capsules")
	cmd.Flags().Uint64("inactivity-period", 0, "Inactivity period in seconds for dead man's switch")
	cmd.Flags().String("title", "", "Capsule title")
	cmd.Flags().String("description", "", "Capsule description")
	
	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

// CmdOpenCapsule returns a CLI command for opening a time capsule
func CmdOpenCapsule(ac address.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "open-capsule [capsule-id]",
		Short: "Open a time capsule",
		Long: `Open a time capsule and retrieve its decrypted data.

Example:
$ simd tx timecapsule open-capsule 1 \
  --key-shares="share1.json,share2.json,share3.json" \
  --from=recipient`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			// Parse capsule ID
			capsuleID, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid capsule ID: %w", err)
			}

			// Get flags
			keyShareFiles, _ := cmd.Flags().GetStringSlice("key-shares")
			signatureFiles, _ := cmd.Flags().GetStringSlice("signatures")

			// Read key shares
			var keyShares []string
			for _, file := range keyShareFiles {
				data, err := readDataFile(file)
				if err != nil {
					return fmt.Errorf("failed to read key share file %s: %w", file, err)
				}
				keyShares = append(keyShares, string(data))
			}

			// Read signatures
			var signatures []string
			for _, file := range signatureFiles {
				data, err := readDataFile(file)
				if err != nil {
					return fmt.Errorf("failed to read signature file %s: %w", file, err)
				}
				signatures = append(signatures, string(data))
			}

			// Create message
			msg := &types.MsgOpenCapsule{
				Accessor:  clientCtx.GetFromAddress().String(),
				CapsuleID: capsuleID,
				KeyShares: keyShares,
				Signatures: signatures,
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	cmd.Flags().StringSlice("key-shares", []string{}, "Key share files (JSON format)")
	cmd.Flags().StringSlice("signatures", []string{}, "Signature files for multi-sig capsules")
	
	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

// CmdUpdateActivity returns a CLI command for updating capsule activity
func CmdUpdateActivity(ac address.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-activity [capsule-id]",
		Short: "Update last activity timestamp for dead man's switch capsule",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			capsuleID, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid capsule ID: %w", err)
			}

			msg := &types.MsgUpdateActivity{
				Owner:     clientCtx.GetFromAddress().String(),
				CapsuleID: capsuleID,
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

// CmdCancelCapsule returns a CLI command for cancelling a time capsule
func CmdCancelCapsule(ac address.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "cancel-capsule [capsule-id]",
		Short: "Cancel a time capsule",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			capsuleID, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid capsule ID: %w", err)
			}

			reason, _ := cmd.Flags().GetString("reason")

			msg := &types.MsgCancelCapsule{
				Owner:     clientCtx.GetFromAddress().String(),
				CapsuleID: capsuleID,
				Reason:    reason,
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	cmd.Flags().String("reason", "", "Reason for cancellation")
	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

// CmdTransferCapsule returns a CLI command for transferring capsule ownership
func CmdTransferCapsule(ac address.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "transfer-capsule [capsule-id] [new-owner]",
		Short: "Transfer ownership of a time capsule",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			capsuleID, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid capsule ID: %w", err)
			}

			newOwner := args[1]

			msg := &types.MsgTransferCapsule{
				CurrentOwner: clientCtx.GetFromAddress().String(),
				NewOwner:     newOwner,
				CapsuleID:    capsuleID,
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

// Helper functions

func readDataFile(filename string) ([]byte, error) {
	// In a real implementation, this would read from the file system
	// For now, return dummy data
	dummyData := map[string]interface{}{
		"message": "This is encrypted capsule data",
		"timestamp": time.Now().Unix(),
		"file": filename,
	}
	return json.Marshal(dummyData)
}

func parseCapsuleType(typeStr string) (types.CapsuleType, error) {
	switch typeStr {
	case "safe":
		return types.CapsuleType_SAFE, nil
	case "time_lock":
		return types.CapsuleType_TIME_LOCK, nil
	case "conditional":
		return types.CapsuleType_CONDITIONAL, nil
	case "multi_sig":
		return types.CapsuleType_MULTI_SIG, nil
	case "dead_mans_switch":
		return types.CapsuleType_DEAD_MANS_SWITCH, nil
	default:
		return types.CapsuleType_UNKNOWN, fmt.Errorf("unknown capsule type: %s", typeStr)
	}
}

// CmdBatchTransferCapsules returns a CLI command for batch transferring capsules
func CmdBatchTransferCapsules(ac address.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "batch-transfer-capsules [transfers-file]",
		Short: "Transfer multiple time capsules in a single transaction",
		Long: `Transfer multiple time capsules in a single transaction.

The transfers file should be a JSON file with the following format:
[
  {
    "capsule_id": 1,
    "new_owner": "cosmos1...",
    "message": "Transfer to Alice"
  },
  {
    "capsule_id": 2, 
    "new_owner": "cosmos1...",
    "message": "Transfer to Bob"
  }
]`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			// Read transfers from file
			transfersFile := args[0]
			transfersData, err := os.ReadFile(transfersFile)
			if err != nil {
				return fmt.Errorf("failed to read transfers file: %w", err)
			}

			var transfers []types.CapsuleTransfer
			if err := json.Unmarshal(transfersData, &transfers); err != nil {
				return fmt.Errorf("failed to parse transfers file: %w", err)
			}

			transferFeeStr, _ := cmd.Flags().GetString("transfer-fee")
			var transferFee sdk.Coins
			if transferFeeStr != "" {
				transferFee, err = sdk.ParseCoinsNormalized(transferFeeStr)
				if err != nil {
					return fmt.Errorf("invalid transfer fee: %w", err)
				}
			}

			requireApproval, _ := cmd.Flags().GetBool("require-approval")

			msg := &types.MsgBatchTransferCapsules{
				CurrentOwner:    clientCtx.GetFromAddress().String(),
				Transfers:       transfers,
				TransferFee:     transferFee,
				RequireApproval: requireApproval,
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	cmd.Flags().String("transfer-fee", "", "Optional transfer fee to pay (e.g., '100stake')")
	cmd.Flags().Bool("require-approval", false, "Require approval from recipients")
	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

// CmdApproveTransfer returns a CLI command for approving a transfer
func CmdApproveTransfer(ac address.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "approve-transfer [transfer-id] [capsule-id] [approved]",
		Short: "Approve or reject a pending capsule transfer",
		Long: `Approve or reject a pending capsule transfer.

Example:
  simd tx timecapsule approve-transfer transfer123 1 true --from alice
  simd tx timecapsule approve-transfer transfer456 2 false --from bob`,
		Args: cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			transferID := args[0]

			capsuleID, err := strconv.ParseUint(args[1], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid capsule ID: %w", err)
			}

			approved, err := strconv.ParseBool(args[2])
			if err != nil {
				return fmt.Errorf("invalid approved value: %w", err)
			}

			msg := &types.MsgApproveTransfer{
				Approver:   clientCtx.GetFromAddress().String(),
				TransferID: transferID,
				CapsuleID:  capsuleID,
				Approved:   approved,
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}