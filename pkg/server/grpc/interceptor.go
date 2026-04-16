package grpc

import (
	"context"
	"fmt"
	"log/slog"
	"mkit/pkg/log"
	"runtime/debug"
	"time"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
)

// UnaryPanicInterceptor recovers from panics in unary RPCs.
func UnaryPanicInterceptor(logger *slog.Logger) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (resp interface{}, err error) {
		defer func() {
			if r := recover(); r != nil {
				stack := string(debug.Stack())

				logger.ErrorContext(ctx, fmt.Sprintf("panic recovered in gRPC handler: %v", r),
					"stack", stack,
					"method", info.FullMethod,
				)

				err = status.Errorf(codes.Internal, "internal error")
			}
		}()

		return handler(ctx, req)
	}
}

func UnaryLogger(logger *slog.Logger) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {

		start := time.Now()

		var reqID string
		if md, ok := metadata.FromIncomingContext(ctx); ok {
			if v := md.Get("x-request-id"); len(v) > 0 {
				reqID = v[0]
			}
		}
		if reqID == "" {
			id, _ := uuid.NewV7()
			reqID = id.String()
		}

		var ip string
		if p, ok := peer.FromContext(ctx); ok {
			ip = p.Addr.String()
		}

		entry := logger.With("request_id", reqID)
		ctx = log.WithLogger(ctx, entry)

		resp, err := handler(ctx, req)

		st := status.Convert(err)
		entry.InfoContext(ctx, "Incoming gRPC request",
			"status", st.Code().String(),
			"duration", time.Since(start),
			"method", info.FullMethod,
			"ip", ip,
		)

		return resp, err
	}
}
