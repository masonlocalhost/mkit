package server

import (
	"context"
	"mkit/example/grpcapp/pkg/api/go/nanoid/v1"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) ListTechnologyCPETypes(ctx context.Context, req *nanoid.ListTechnologyCPETypesRequest) (*nanoid.ListTechnologyCPETypesResponse, error) {
	types, err := s.technologyService.FindTechnologyCPETypes(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	return &nanoid.ListTechnologyCPETypesResponse{
		CpeTypes: types,
	}, nil
}
