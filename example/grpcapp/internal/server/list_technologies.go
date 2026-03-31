package server

import (
	"context"
	"mkit/example/grpcapp/pkg/api/go/nanoid/v1"
	"mkit/example/grpcapp/pkg/converter"
	"mkit/pkg/sqlutil"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) ListTechnologies(ctx context.Context, req *nanoid.ListTechnologiesRequest) (*nanoid.ListTechnologiesResponse, error) {
	sorts, err := buildSorts(req.GetSortBy(), req.GetSortOrder())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}

	var createdFrom, createdTo time.Time
	if req.CreatedFrom != nil {
		createdFrom = req.CreatedFrom.AsTime()
	}
	if req.CreatedTo != nil {
		createdTo = req.CreatedTo.AsTime()
	}

	techs, total, err := s.technologyService.FindByFilters(
		ctx,
		req.GetVendors(),
		req.GetCpeTypes(),
		createdFrom,
		createdTo,
		req.GetSearch(),
		int(req.GetLimit()),
		int(req.GetOffset()),
		true,
		sorts,
	)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	return &nanoid.ListTechnologiesResponse{
		Technologies: converter.TechnologiesToPbs(techs),
		Total:        int64(total),
		Offset:       req.GetOffset(),
		Limit:        req.GetLimit(),
	}, nil
}

func buildSorts(sortBy, sortOrder string) ([]sqlutil.SortItem, error) {
	if sortBy == "" {
		return nil, nil
	}

	order := "asc"
	if sortOrder != "" {
		order = sortOrder
	}

	return []sqlutil.SortItem{
		{
			Field:     sortBy,
			SortValue: sqlutil.SortValue(order),
		},
	}, nil
}
