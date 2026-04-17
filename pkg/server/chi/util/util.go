package chiutil

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"mkit/pkg/error/serviceerror"
	"mkit/pkg/log"
	"mkit/pkg/restful/rpc"
	"mkit/pkg/sqlutil"
	"net/http"

	"github.com/go-playground/form/v4"
)

var decoder = form.NewDecoder()

func BindQuery(r *http.Request, dst any) error {
	if err := r.ParseForm(); err != nil {
		return err
	}
	return decoder.Decode(dst, r.Form)
}

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

func HandleError(w http.ResponseWriter, r *http.Request, err error) {
	var (
		sErr   *serviceerror.Error
		logger = log.GetLogger(r.Context())
	)

	if errors.Is(err, context.Canceled) {
		w.WriteHeader(499)
		return
	}

	if errors.As(err, &sErr) {
		if sErr.StatusCode == http.StatusInternalServerError {
			logger.ErrorContext(r.Context(), "internal service error", "error", err)
		}
		JSON(w, sErr.StatusCode, rpc.NewResponse(false).SetMessage(sErr.Message).SetError(err))
		return
	}

	if errors.Is(err, serviceerror.ErrNotFound) {
		JSON(w, http.StatusNotFound, rpc.NewResponse(false).SetMessage("Not found.").SetError(err))
		return
	}
	if errors.Is(err, serviceerror.ErrInvalidArgument) {
		JSON(w, http.StatusBadRequest, rpc.NewResponse(false).SetMessage("Invalid argument.").SetError(err))
		return
	}
	if errors.Is(err, serviceerror.ErrPermissionDenied) {
		JSON(w, http.StatusForbidden, rpc.NewResponse(false).SetMessage("Permission denied.").SetError(err))
		return
	}
	if errors.Is(err, serviceerror.ErrUnauthenticated) {
		JSON(w, http.StatusUnauthorized, rpc.NewResponse(false).SetMessage("Unauthenticated.").SetError(err))
		return
	}

	logger.ErrorContext(r.Context(), "internal error", "error", err)
	JSON(w, http.StatusInternalServerError, rpc.NewResponse(false).SetMessage("Internal server error.").SetError(err))
}

func HandleSuccess(w http.ResponseWriter, r *http.Request, data any, paginationMetas ...int) {
	res := rpc.NewResponse(true).SetData(data)
	if len(paginationMetas) == 3 {
		res.SetMeta(paginationMetas[0], paginationMetas[1], paginationMetas[2])
	}
	JSON(w, http.StatusOK, res)
}

func JSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}

func MustGetValue[T any](ctx context.Context, key any) T {
	value := ctx.Value(key)
	if value == nil {
		panic(fmt.Sprintf("context value not exist for key: %v", key))
	}
	v, ok := value.(T)
	if !ok {
		panic(fmt.Sprintf("context value exist but type not match for key: %v", key))
	}
	return v
}
