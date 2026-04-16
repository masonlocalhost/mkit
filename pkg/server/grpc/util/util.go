package grpcutil

import (
	"context"
	"errors"
	"mkit/pkg/error/serviceerror"
	"mkit/pkg/log"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func HandleError(ctx context.Context, err error) error {
	var (
		logger = log.GetLogger(ctx)
	)

	if errors.Is(err, context.Canceled) {
		return status.Error(codes.Canceled, err.Error())
	}
	if errors.Is(err, serviceerror.ErrNotFound) {
		return status.Error(codes.NotFound, err.Error())
	}
	if errors.Is(err, serviceerror.ErrInvalidArgument) {
		return status.Error(codes.InvalidArgument, err.Error())
	}
	if errors.Is(err, serviceerror.ErrPermissionDenied) {
		return status.Error(codes.PermissionDenied, err.Error())
	}
	if errors.Is(err, serviceerror.ErrUnauthenticated) {
		return status.Error(codes.Unauthenticated, err.Error())
	}

	logger.ErrorContext(ctx, "internal error", "error", err)
	// Error internal

	return status.Error(codes.Internal, err.Error())
}
