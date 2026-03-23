package grpc

import (
	"fmt"
	"mkit/pkg/config"

	"buf.build/go/protovalidate"
	"connectrpc.com/vanguard"
	"connectrpc.com/vanguard/vanguardgrpc"
	pvinterceptor "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/protovalidate"
	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/encoding"
	"google.golang.org/grpc/health"
	healthv1 "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/encoding/protojson"
)

func New(
	cfg *config.App, logger *logrus.Logger,
) (*grpc.Server, *health.Server, error) {
	validator, err := protovalidate.New()
	if err != nil {
		return nil, nil, fmt.Errorf("cant setup validator: %w", err)
	}

	var (
		grpcCfg = cfg.GRPC
		opts    []grpc.ServerOption
	)

	if cfg.Tracing.Enabled {
		opts = append(opts, grpc.StatsHandler(otelgrpc.NewServerHandler()))
	}

	opts = append(opts, grpc.ChainUnaryInterceptor(
		pvinterceptor.UnaryServerInterceptor(validator),
		UnaryLogger(logger),
		UnaryPanicInterceptor(logger),
	), grpc.MaxRecvMsgSize(104857600)) // 100MB

	if grpcCfg.JsonTranscodeEnabled {
		encoding.RegisterCodec(vanguardgrpc.NewCodec(&vanguard.JSONCodec{
			MarshalOptions:   protojson.MarshalOptions{EmitUnpopulated: true},
			UnmarshalOptions: protojson.UnmarshalOptions{DiscardUnknown: true},
		}))
	}
	healthServer := health.NewServer()
	srv := grpc.NewServer(opts...)
	healthv1.RegisterHealthServer(srv, healthServer)

	if grpcCfg.ReflectionEnabled {
		reflection.Register(srv)
	}

	return srv, healthServer, nil
}
