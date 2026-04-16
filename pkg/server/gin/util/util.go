package ginutil

import (
	"context"
	"errors"
	"fmt"
	"mkit/pkg/error/serviceerror"
	"mkit/pkg/log"
	"mkit/pkg/restful/rpc"
	"mkit/pkg/sqlutil"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetSortArray(sortBy, sortOrder string, orderableMap map[string]bool) ([]sqlutil.SortItem, error) {
	if sortBy != "" && (sortOrder == "" || !orderableMap[sortBy]) {
		return nil, serviceerror.NewInvalidArgument(fmt.Errorf("sort params not valid")).
			SetMessage("Input validation error.")
	}

	if sortBy == "" {
		return nil, nil
	}

	return []sqlutil.SortItem{{
		Field:     sortBy,
		SortValue: sqlutil.SortValue(sortOrder),
	}}, nil
}

func HandleError(c *gin.Context, err error) {
	var (
		sErr   *serviceerror.Error
		logger = log.GetLogger(c.Request.Context())
	)

	if errors.Is(err, context.Canceled) {
		c.Status(499) // status 499 client closed

		return
	}

	if errors.As(err, &sErr) {
		if sErr.StatusCode == http.StatusInternalServerError {
			logger.ErrorContext(c.Request.Context(), "internal service error", "error", err)
		}

		c.JSON(
			sErr.StatusCode,
			rpc.NewResponse(false).SetMessage(sErr.Message).SetError(err),
		)

		return
	}

	if errors.Is(err, serviceerror.ErrNotFound) {
		c.JSON(
			http.StatusNotFound,
			rpc.NewResponse(false).SetMessage("Not found.").SetError(err),
		)

		return
	}
	if errors.Is(err, serviceerror.ErrInvalidArgument) {
		c.JSON(
			http.StatusBadRequest,
			rpc.NewResponse(false).SetMessage("Invalid argument.").SetError(err),
		)

		return
	}
	if errors.Is(err, serviceerror.ErrPermissionDenied) {
		c.JSON(
			http.StatusForbidden,
			rpc.NewResponse(false).SetMessage("Permission denied.").SetError(err),
		)

		return
	}
	if errors.Is(err, serviceerror.ErrUnauthenticated) {
		c.JSON(
			http.StatusUnauthorized,
			rpc.NewResponse(false).SetMessage("Unauthenticated.").SetError(err),
		)

		return
	}

	logger.ErrorContext(c.Request.Context(), "internal error", "error", err)
	// Error internal
	c.JSON(
		http.StatusInternalServerError,
		rpc.NewResponse(false).SetMessage("Internal server error.").SetError(err),
	)

	return
}

func HandleSuccess(c *gin.Context, data any, paginationMetas ...int) {
	res := rpc.NewResponse(true).SetData(data)
	if len(paginationMetas) == 3 {
		res.SetMeta(paginationMetas[0], paginationMetas[1], paginationMetas[2])
	}

	c.JSON(
		http.StatusOK,
		res,
	)
}

func MustGetValue[T any](c *gin.Context, key string) T {
	value, exists := c.Get(key)
	if !exists {
		panic(fmt.Sprintf("gin context value not exist for key: %s", key))
	}

	v, ok := value.(T)
	if !ok {
		panic(fmt.Sprintf("gin context value exist but type not match for key: %s", key))
	}

	return v
}
