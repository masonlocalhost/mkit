package technology

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"mkit/example/ginapp/internal/model"
	"mkit/example/ginapp/internal/repository/technology"
	"mkit/pkg/error/repoerror"
	"mkit/pkg/error/serviceerror"
	"mkit/pkg/sqlutil"
	"time"
)

type Service struct {
	logger         *slog.Logger
	technologyRepo *technology.Repository
}

func NewService(logger *slog.Logger, technologyRepo *technology.Repository) *Service {
	return &Service{
		logger:         logger,
		technologyRepo: technologyRepo,
	}
}

func (s *Service) FindByFilters(
	ctx context.Context, vendors, cpeTypes []string, createdAtFrom, createdAtTo time.Time, search string,
	limit, offset int, isCount bool, sorts []sqlutil.SortItem,
) ([]*model.Technology, int, error) {
	techs, total, err := s.technologyRepo.FindTechnologiesByFilters(
		ctx, vendors, cpeTypes, createdAtFrom, createdAtTo, search, limit, offset,
		isCount, sorts,
	)
	if err != nil {
		return nil, 0, serviceerror.NewInternal(fmt.Errorf("cannot find technologies: %v", err))
	}

	return techs, total, nil
}

func (s *Service) FirstByID(ctx context.Context, id string) (*model.Technology, error) {
	tech, err := s.technologyRepo.FirstTechnologyByID(ctx, id)
	if err != nil {
		if errors.Is(err, repoerror.ErrNotFound) {
			return nil, serviceerror.NewNotFound(err).SetMessage("Technology not found")
		}

		return nil, serviceerror.NewInternal(fmt.Errorf("cannot find technology by id: %v", err))
	}

	return tech, nil
}

func (s *Service) FindTechnologyVendors(ctx context.Context) ([]string, error) {
	return s.technologyRepo.FindTechnologyColumnValues(ctx, model.Technology_VENDOR_COLUMN, nil)
}

func (s *Service) FindTechnologyCPETypes(ctx context.Context) ([]string, error) {
	return s.technologyRepo.FindTechnologyColumnValues(ctx, model.Technology_CPE_TYPE_COLUMN, nil)
}
