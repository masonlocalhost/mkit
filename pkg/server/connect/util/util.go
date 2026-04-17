package connectutil

import (
	"context"
	"errors"
	"mkit/pkg/error/serviceerror"
	"mkit/pkg/log"

	"connectrpc.com/connect"
)

func HandleError(ctx context.Context, err error) error {
	logger := log.GetLogger(ctx)

	if errors.Is(err, context.Canceled) {
		return connect.NewError(connect.CodeCanceled, err)
	}
	if errors.Is(err, serviceerror.ErrNotFound) {
		return connect.NewError(connect.CodeNotFound, err)
	}
	if errors.Is(err, serviceerror.ErrInvalidArgument) {
		return connect.NewError(connect.CodeInvalidArgument, err)
	}
	if errors.Is(err, serviceerror.ErrPermissionDenied) {
		return connect.NewError(connect.CodePermissionDenied, err)
	}
	if errors.Is(err, serviceerror.ErrUnauthenticated) {
		return connect.NewError(connect.CodeUnauthenticated, err)
	}

	var sErr *serviceerror.Error
	if errors.As(err, &sErr) {
		logger.ErrorContext(ctx, "internal service error", "error", err)
		return connect.NewError(connect.CodeInternal, err)
	}

	logger.ErrorContext(ctx, "internal error", "error", err)
	return connect.NewError(connect.CodeInternal, err)
}
