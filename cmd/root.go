package cmd

import (
	"context"

	"github.com/spf13/cobra"
)

// NewRootCmd creates a new instance of the root command
func NewRootCmd() *cobra.Command {
	ctx := context.Background()

	cmd := &cobra.Command{
		Use:   "grpc-sample",
		Short: "Go gRPC sample is a sample to realize micro service with grpc.",
		Long: `
Go gRPC sample is a sample to realize micro service with grpc.`,
	}

	cmd.AddCommand(NewServerStartCmd(ctx))
	cmd.AddCommand(NewServiceStartCmd(ctx))

	return cmd
}
