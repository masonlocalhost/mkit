package converter

import (
	"mkit/example/ginapp/internal/model"
	"mkit/example/ginapp/pkg/dto"
)

func TechnologyToDTO(t *model.Technology) *dto.Technology {
	if t == nil {
		return nil
	}

	var technology = &dto.Technology{
		ID:           t.ID,
		Name:         t.Name,
		Version:      t.Version,
		Vendor:       t.Vendor,
		CPEType:      t.CPEType,
		ThumbnailUrl: t.ThumbnailUrl,
		CreatedAt:    t.CreatedAt,
		UpdatedAt:    t.UpdatedAt,
	}

	return technology
}

func TechnologiesToDTOs(techs []*model.Technology) []*dto.Technology {
	if techs == nil {
		return nil
	}

	result := make([]*dto.Technology, 0, len(techs))
	for _, t := range techs {
		result = append(result, TechnologyToDTO(t))
	}

	return result
}
