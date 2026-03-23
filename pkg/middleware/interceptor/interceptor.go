package interceptor

import (
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
	"runtime/debug"
	"time"
)

// UnaryPanicInterceptor recovers from panics in unary RPCs
func UnaryPanicInterceptor(logger *logrus.Logger) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (resp interface{}, err error) {
		defer func() {
			if r := recover(); r != nil {
				stack := string(debug.Stack())

				logger.WithFields(logrus.Fields{
					"stack":  stack,
					"method": info.FullMethod,
				}).Error(fmt.Sprintf("panic recovered in Gin handler: %v", r))

				err = status.Errorf(codes.Internal, "internal error")
			}
		}()

		return handler(ctx, req)
	}
}

func UnaryLogger(logger *logrus.Logger) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		start := time.Now()

		// Call the handler
		resp, err := handler(ctx, req)

		duration := time.Since(start)
		st := status.Convert(err)

		// Get client IP if available
		var ip string
		if p, ok := peer.FromContext(ctx); ok {
			ip = p.Addr.String()
		}

		logger.WithFields(logrus.Fields{
			"method":   info.FullMethod,
			"status":   st.Code().String(),
			"duration": duration,
			"ip":       ip,
			"error":    st.Message(),
		}).Info("Incoming gRPC request")

		return resp, err
	}
}
