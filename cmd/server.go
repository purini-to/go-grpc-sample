package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	grpc_retry "github.com/grpc-ecosystem/go-grpc-middleware/retry"

	"github.com/go-chi/chi/middleware"

	"github.com/purini-to/go-grpc-sample/pkg/cat"

	"google.golang.org/grpc"

	"go.uber.org/zap"

	"github.com/pkg/errors"

	"github.com/go-chi/chi"

	"github.com/spf13/cobra"
)

// ServerStartOptions are the command flags
type ServerStartOptions struct {
	port               uint
	catServiceEndpoint string
}

// NewServerStartCmd creates a new http server command
func NewServerStartCmd(ctx context.Context) *cobra.Command {
	opts := &ServerStartOptions{}

	cmd := &cobra.Command{
		Use:   "start",
		Short: "Starts a Go gRPC sample web server",
		RunE: func(cmd *cobra.Command, args []string) error {
			return RunServerStart(ctx, opts)
		},
	}

	cmd.PersistentFlags().UintVarP(&opts.port, "port", "p", 8080, "Service HTTP port")
	cmd.PersistentFlags().StringVarP(&opts.catServiceEndpoint, "cat-endpoint", "", "127.0.0.1:6565", "Cat service endpoint")

	return cmd
}

func RunServerStart(ctx context.Context, opts *ServerStartOptions) error {
	logger, err := initLog()
	if err != nil {
		return errors.Wrap(err, "could not initialize log")
	}

	logger.Info("Go gRPC sample starting...", zap.String("version", "v0.0.9"))

	conn, err := grpc.Dial(opts.catServiceEndpoint,
		grpc.WithInsecure(),
		grpc.WithDefaultServiceConfig(`{"loadBalancingConfig": [{"round_robin":{}}]}`),
		grpc.WithUnaryInterceptor(grpc_retry.UnaryClientInterceptor(
			grpc_retry.WithMax(3),
			grpc_retry.WithBackoff(grpc_retry.BackoffExponentialWithJitter(50*time.Millisecond, 0.10)),
		)),
	)
	if err != nil {
		return errors.Wrap(err, "could not gRPC connection")
	}
	defer conn.Close()

	client := cat.NewCatClient(conn)

	r := chi.NewRouter()
	r.Use(middleware.Recoverer)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("welcome"))
	})

	r.Get("/cat/{catName}", func(w http.ResponseWriter, r *http.Request) {
		defer func(t time.Time) {
			logger.Info(
				"Access log",
				zap.String("uri", r.URL.String()),
				zap.String("method", r.Method),
				zap.Duration("duration", time.Since(t)),
			)
		}(time.Now())
		catName := chi.URLParam(r, "catName")

		message := &cat.GetMyCatMessage{TargetCat: catName}
		res, err := client.GetMyCat(r.Context(), message)
		if err != nil {
			logger.Error("Failed GetMyCat service.", zap.Error(err), zap.String("TargetCat", catName))
			panic("Failed GetMyCat service")
		}

		json.NewEncoder(w).Encode(res)
	})

	logger.Info(fmt.Sprintf("Listen starting on port %d", opts.port), zap.Uint("port", opts.port))
	err = http.ListenAndServe(fmt.Sprintf(":%d", opts.port), r)
	if err != nil {
		return errors.Wrapf(err, "could not listen of port %d", opts.port)
	}

	return nil
}
