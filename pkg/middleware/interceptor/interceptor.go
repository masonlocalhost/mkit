package interceptor

import (
	"context"
	"fmt"
	"log/slog"
	"runtime/debug"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
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

		resp, err := handler(ctx, req)

		duration := time.Since(start)
		st := status.Convert(err)

		var ip string
		if p, ok := peer.FromContext(ctx); ok {
			ip = p.Addr.String()
		}

		logger.InfoContext(ctx, "Incoming gRPC request",
			"method", info.FullMethod,
			"status", st.Code().String(),
			"duration", duration,
			"ip", ip,
			"error", st.Message(),
		)

		return resp, err
	}
}
