package converter

import (
	"mkit/example/grpcapp/internal/model"
	"mkit/example/grpcapp/pkg/api/go/nanoid/v1"

	"google.golang.org/protobuf/types/known/timestamppb"
)

func TechnologyToPb(t *model.Technology) *nanoid.Technology {
	if t == nil {
		return nil
	}

	return &nanoid.Technology{
		Id:           t.ID,
		Name:         t.Name,
		Version:      t.Version,
		Vendor:       t.Vendor,
		CpeType:      t.CPEType,
		ThumbnailUrl: t.ThumbnailUrl,
		CreatedAt:    timestamppb.New(t.CreatedAt),
		UpdatedAt:    timestamppb.New(t.UpdatedAt),
	}
}

func TechnologiesToPbs(list []*model.Technology) []*nanoid.Technology {
	res := make([]*nanoid.Technology, 0, len(list))
	for _, t := range list {
		res = append(res, TechnologyToPb(t))
	}
	return res
}
