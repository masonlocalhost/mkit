package grpc

import (
	"context"
	"fmt"
	"mkit/pkg/log"
	"runtime/debug"
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
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

				logger.WithContext(ctx).WithFields(logrus.Fields{
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

		// Request ID (from metadata if present)
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

		// Client IP
		var ip string
		if p, ok := peer.FromContext(ctx); ok {
			ip = p.Addr.String()
		}

		entry := logrus.NewEntry(logger).WithFields(logrus.Fields{
			"request_id": reqID,
		})
		ctx = log.WithLogger(ctx, entry)

		// Call handler with enriched context
		resp, err := handler(ctx, req)

		// Status + duration
		duration := time.Since(start)
		st := status.Convert(err)

		entry.WithFields(logrus.Fields{
			"status":   st.Code().String(),
			"duration": duration,
			"method":   info.FullMethod,
			"ip":       ip,
		}).Info("Incoming gRPC request")

		return resp, err
	}
}
