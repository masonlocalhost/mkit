package server

import (
	"context"
	"mkit/example/grpcapp/pkg/api/go/nanoid/v1"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) ListTechnologyVendors(ctx context.Context, req *nanoid.ListTechnologyVendorsRequest) (*nanoid.ListTechnologyVendorsResponse, error) {
	vendors, err := s.technologyService.FindTechnologyVendors(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	return &nanoid.ListTechnologyVendorsResponse{
		Vendors: vendors,
	}, nil
}
