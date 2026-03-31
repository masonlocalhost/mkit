package server

import (
	"context"
	"mkit/example/grpcapp/pkg/api/go/nanoid/v1"
	"mkit/example/grpcapp/pkg/converter"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) GetTechnology(ctx context.Context, req *nanoid.GetTechnologyRequest) (*nanoid.GetTechnologyResponse, error) {
	tech, err := s.technologyService.FirstByID(ctx, req.GetId())
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	return &nanoid.GetTechnologyResponse{
		Technology: converter.TechnologyToPb(tech),
	}, nil
}
