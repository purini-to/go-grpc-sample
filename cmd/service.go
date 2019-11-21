package cmd

import (
	"context"
	"fmt"
	"net"

	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"

	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"

	"github.com/purini-to/go-grpc-sample/pkg/cat"
	"google.golang.org/grpc"

	"github.com/pkg/errors"

	"github.com/spf13/cobra"
)

// ServiceStartOptions are the command flags
type ServiceStartOptions struct {
	port uint
}

// NewServiceStartCmd creates a new http Service command
func NewServiceStartCmd(ctx context.Context) *cobra.Command {
	opts := &ServiceStartOptions{}

	cmd := &cobra.Command{
		Use:   "service",
		Short: "Starts a Go gRPC sample gRPC Service",
		RunE: func(cmd *cobra.Command, args []string) error {
			return RunServiceStart(ctx, opts)
		},
	}

	cmd.PersistentFlags().UintVarP(&opts.port, "port", "p", 6565, "Service gRPC port")

	return cmd
}

func RunServiceStart(ctx context.Context, opts *ServiceStartOptions) error {
	logger, err := initLog()
	if err != nil {
		return errors.Wrap(err, "could not initialize log")
	}
	grpc_zap.ReplaceGrpcLoggerV2(logger)

	logger.Info("Go gRPC sample service starting...")

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", opts.port))
	if err != nil {
		return errors.Wrapf(err, "could not listen of port %d", opts.port)
	}

	server := grpc.NewServer(
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			grpc_ctxtags.UnaryServerInterceptor(),
			grpc_zap.UnaryServerInterceptor(logger),
			grpc_recovery.UnaryServerInterceptor(),
		)),
	)
	service := cat.NewCatService()
	cat.RegisterCatServer(server, service)
	err = server.Serve(listener)
	if err != nil {
		return errors.Wrapf(err, "could not serve gRPC server. service: cat")
	}

	return nil
}
