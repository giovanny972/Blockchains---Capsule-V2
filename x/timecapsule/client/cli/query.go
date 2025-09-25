package cli

import (
	"context"
	"fmt"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	
	"github.com/cosmos/cosmos-sdk/x/timecapsule/types"
)

// GetQueryCmd returns the cli query commands for this module
func GetQueryCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("Querying commands for the %s module", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(
		CmdQueryParams(),
		CmdQueryCapsule(),
		CmdQueryCapsules(),
		CmdQueryUserCapsules(),
		CmdQueryCapsulesByType(),
		CmdQueryCapsulesByStatus(),
		CmdQueryStats(),
		CmdQueryKeyShares(),
		CmdQueryConditionContract(),
		CmdQueryConditionContracts(),
	)

	return cmd
}

// CmdQueryParams implements the params query command
func CmdQueryParams() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "params",
		Args:  cobra.NoArgs,
		Short: "Query the current time capsule parameters",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)
			res, err := queryClient.Params(context.Background(), &types.QueryParamsRequest{})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(&res.Params)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// CmdQueryCapsule implements the capsule query command
func CmdQueryCapsule() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "capsule [capsule-id]",
		Short: "Query details of a specific time capsule",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			capsuleID, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid capsule ID: %w", err)
			}

			queryClient := types.NewQueryClient(clientCtx)
			res, err := queryClient.Capsule(context.Background(), &types.QueryCapsuleRequest{
				CapsuleId: capsuleID,
			})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// CmdQueryCapsules implements the capsules query command
func CmdQueryCapsules() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "capsules",
		Short: "Query all time capsules",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)
			res, err := queryClient.Capsules(context.Background(), &types.QueryCapsulesRequest{})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// CmdQueryUserCapsules implements the user-capsules query command
func CmdQueryUserCapsules() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "user-capsules [owner-address]",
		Short: "Query all capsules owned by a specific address",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			owner := args[0]

			queryClient := types.NewQueryClient(clientCtx)
			res, err := queryClient.UserCapsules(context.Background(), &types.QueryUserCapsulesRequest{
				Owner: owner,
			})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// CmdQueryCapsulesByType implements the capsules-by-type query command
func CmdQueryCapsulesByType() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "capsules-by-type [capsule-type]",
		Short: "Query capsules filtered by type",
		Long: `Query capsules filtered by type.

Valid types: safe, time_lock, conditional, multi_sig, dead_mans_switch`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			capsuleType, err := parseCapsuleType(args[0])
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)
			res, err := queryClient.CapsulesByType(context.Background(), &types.QueryCapsulesByTypeRequest{
				CapsuleType: capsuleType,
			})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// CmdQueryCapsulesByStatus implements the capsules-by-status query command
func CmdQueryCapsulesByStatus() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "capsules-by-status [status]",
		Short: "Query capsules filtered by status",
		Long: `Query capsules filtered by status.

Valid statuses: active, unlocked, expired, cancelled`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			status, err := parseCapsuleStatus(args[0])
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)
			res, err := queryClient.CapsulesByStatus(context.Background(), &types.QueryCapsulesByStatusRequest{
				Status: status,
			})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// CmdQueryStats implements the stats query command
func CmdQueryStats() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "stats",
		Short: "Query time capsule module statistics",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)
			res, err := queryClient.Stats(context.Background(), &types.QueryStatsRequest{})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// CmdQueryKeyShares implements the key-shares query command
func CmdQueryKeyShares() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "key-shares [capsule-id]",
		Short: "Query key shares for a specific capsule",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			capsuleID, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid capsule ID: %w", err)
			}

			queryClient := types.NewQueryClient(clientCtx)
			res, err := queryClient.KeyShares(context.Background(), &types.QueryKeySharesRequest{
				CapsuleId: capsuleID,
			})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// CmdQueryConditionContract implements the condition-contract query command
func CmdQueryConditionContract() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "condition-contract [address]",
		Short: "Query details of a condition contract",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			address := args[0]

			queryClient := types.NewQueryClient(clientCtx)
			res, err := queryClient.ConditionContract(context.Background(), &types.QueryConditionContractRequest{
				Address: address,
			})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// CmdQueryConditionContracts implements the condition-contracts query command
func CmdQueryConditionContracts() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "condition-contracts",
		Short: "Query all condition contracts",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)
			res, err := queryClient.ConditionContracts(context.Background(), &types.QueryConditionContractsRequest{})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// Helper functions

func parseCapsuleStatus(statusStr string) (types.CapsuleStatus, error) {
	switch statusStr {
	case "active":
		return types.CapsuleStatus_ACTIVE, nil
	case "unlocked":
		return types.CapsuleStatus_UNLOCKED, nil
	case "expired":
		return types.CapsuleStatus_EXPIRED, nil
	case "cancelled":
		return types.CapsuleStatus_CANCELLED, nil
	default:
		return types.CapsuleStatus_UNKNOWN, fmt.Errorf("unknown capsule status: %s", statusStr)
	}
}