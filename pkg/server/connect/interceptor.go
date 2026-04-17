package connect

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"mkit/pkg/log"
	"runtime/debug"
	"time"

	"buf.build/go/protovalidate"
	"connectrpc.com/connect"
	"github.com/google/uuid"
	"google.golang.org/protobuf/proto"
)

func UnaryLogger(logger *slog.Logger) connect.UnaryInterceptorFunc {
	return func(next connect.UnaryFunc) connect.UnaryFunc {
		return func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
			start := time.Now()

			reqID := req.Header().Get("x-request-id")
			if reqID == "" {
				id, _ := uuid.NewV7()
				reqID = id.String()
			}

			entry := logger.With("request_id", reqID)
			ctx = log.WithLogger(ctx, entry)

			resp, err := next(ctx, req)

			code := "ok"
			if err != nil {
				var ce *connect.Error
				if errors.As(err, &ce) {
					code = ce.Code().String()
				} else {
					code = "unknown"
				}
			}

			entry.InfoContext(ctx, "Incoming ConnectRPC request",
				"status", code,
				"duration", time.Since(start),
				"procedure", req.Spec().Procedure,
				"peer", req.Peer().Addr,
			)

			return resp, err
		}
	}
}

func UnaryValidation(validator protovalidate.Validator) connect.UnaryInterceptorFunc {
	return func(next connect.UnaryFunc) connect.UnaryFunc {
		return func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
			if msg, ok := req.Any().(proto.Message); ok {
				if err := validator.Validate(msg); err != nil {
					return nil, connect.NewError(connect.CodeInvalidArgument, err)
				}
			}
			return next(ctx, req)
		}
	}
}

func UnaryPanicInterceptor(logger *slog.Logger) connect.UnaryInterceptorFunc {
	return func(next connect.UnaryFunc) connect.UnaryFunc {
		return func(ctx context.Context, req connect.AnyRequest) (resp connect.AnyResponse, err error) {
			defer func() {
				if r := recover(); r != nil {
					stack := string(debug.Stack())
					logger.ErrorContext(ctx, fmt.Sprintf("panic recovered in ConnectRPC handler: %v", r),
						"stack", stack,
						"procedure", req.Spec().Procedure,
					)
					err = connect.NewError(connect.CodeInternal, fmt.Errorf("internal error"))
				}
			}()
			return next(ctx, req)
		}
	}
}
