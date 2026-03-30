package technology

import (
	"context"
	"errors"
	"fmt"
	"mkit/example/ginapp/internal/model"
	"mkit/pkg/error/repoerror"
	"mkit/pkg/sqlutil"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) GetDB() *gorm.DB {
	return r.db
}

func (r *Repository) UpsertWithTx(tx *gorm.DB, name, version string) (*model.Technology, error) {
	var tech = &model.Technology{
		Name:    name,
		Version: version,
	}

	err := tx.Clauses(
		clause.OnConflict{
			Columns: []clause.Column{
				{Name: model.Technology_NAME_COLUMN},
				{Name: model.Technology_VERSION_COLUMN},
			},
			UpdateAll: true,
		},
		clause.Returning{},
	).Create(tech).Error
	if err != nil {
		return nil, err
	}

	return tech, nil
}

func (r *Repository) UpsertBatchWithTx(tx *gorm.DB, techs []*model.Technology) ([]*model.Technology, error) {
	err := tx.Clauses(
		clause.OnConflict{
			Columns: []clause.Column{
				{Name: model.Technology_NAME_COLUMN},
				{Name: model.Technology_VERSION_COLUMN},
			},
			UpdateAll: true,
		},
		clause.Returning{},
	).Create(&techs).Error
	if err != nil {
		return nil, err
	}

	return techs, nil
}

func (r *Repository) FirstTechnologyByID(
	ctx context.Context, id string, //isPreloadCollections bool,
) (*model.Technology, error) {
	var (
		tech model.Technology
		tx   = r.db.WithContext(ctx)
	)

	// if isPreloadCollections {
	// 	tx = tx.Preload(fmt.Sprintf(
	// 		"%s.%s", model.Technology_TECHNOLOGY_COLLECTIONS_PRELOAD, model.TechnologyCollection_COLLECTION_PRELOAD,
	// 	))
	// }

	if err := tx.Where(sqlutil.EqualClause(model.Technology_ID_COLUMN), id).First(&tech).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, repoerror.ErrNotFound
		}

		return nil, err
	}

	return &tech, nil
}

func (r *Repository) FindTechnologiesByFilters(
	ctx context.Context, vendors, cpeTypes []string, createdAtFrom, createdAtTo time.Time, search string,
	limit, offset int, isCount bool, sorts []sqlutil.SortItem, //isPreloadCollections bool,
) ([]*model.Technology, int, error) {
	var (
		technologies []*model.Technology
		total        int64
	)
	db := r.db.WithContext(ctx)
	tx := db.WithContext(ctx).Model(&model.Technology{})

	// if isPreloadCollections {
	// 	tx = tx.Preload(fmt.Sprintf(
	// 		"%s.%s", model.Technology_TECHNOLOGY_COLLECTIONS_PRELOAD, model.TechnologyCollection_COLLECTION_PRELOAD,
	// 	))
	// }

	// Filter by collectionIDs via join table if needed
	// if len(collectionIDs) > 0 {
	// 	subQuery := db.Model(&model.TechnologyCollection{}).
	// 		Select("1").
	// 		Where(
	// 			fmt.Sprintf("%s.%s = %s.%s",
	// 				model.TechnologyCollectionsTableName,
	// 				model.TechnologyCollection_TECHNOLOGY_ID_COLUMN,
	// 				model.TechnologyTableName,
	// 				model.Technology_ID_COLUMN),
	// 		).Where(
	// 		sqlutil.InClauseWithTable(model.TechnologyCollectionsTableName, model.TechnologyCollection_COLLECTION_ID_COLUMN),
	// 		collectionIDs,
	// 	)

	// 	tx = tx.Where(
	// 		sqlutil.ExistClause(), subQuery,
	// 	)
	// }

	if len(cpeTypes) > 0 {
		tx = tx.Where(sqlutil.InClauseWithTable(model.TechnologyTableName, model.Technology_CPE_TYPE_COLUMN), cpeTypes)
	}
	if len(vendors) > 0 {
		tx = tx.Where(sqlutil.InClauseWithTable(model.TechnologyTableName, model.Technology_VENDOR_COLUMN), vendors)
	}
	if !createdAtFrom.IsZero() && !createdAtTo.IsZero() {
		tx = tx.Where(sqlutil.BetweenClauseWithTable(model.TechnologyTableName, model.Technology_CREATED_AT_COLUMN), createdAtFrom, createdAtTo)
	} else if !createdAtTo.IsZero() {
		tx = tx.Where(sqlutil.LessThanOrEqualClauseWithTable(model.TechnologyTableName, model.Technology_CREATED_AT_COLUMN), createdAtTo)
	} else if !createdAtFrom.IsZero() {
		tx = tx.Where(sqlutil.GreaterThanOrEqualClauseWithTable(model.TechnologyTableName, model.Technology_CREATED_AT_COLUMN), createdAtFrom)
	}

	// Basic search on name or version
	if search != "" {
		searchLike := "%" + search + "%"
		tx = tx.Where(
			sqlutil.ConcatClauses([]*sqlutil.ConcatClause{
				{
					Clause:   sqlutil.ILikeClauseWithTable(model.TechnologyTableName, model.Technology_NAME_COLUMN),
					Operator: sqlutil.LogicOperator_OR,
				},
				{
					Clause:   sqlutil.ILikeClauseWithTable(model.TechnologyTableName, model.Technology_VERSION_COLUMN),
					Operator: sqlutil.LogicOperator_NONE,
				},
			}),
			searchLike, searchLike,
		)
	}

	// Count total
	if isCount {
		if err := tx.Count(&total).Error; err != nil {
			return nil, 0, err
		}
	}

	// Sorting
	if len(sorts) > 0 {
		for _, item := range sorts {
			tx = tx.Order(fmt.Sprintf("%s %s", item.Field, item.SortValue))
		}
	}

	// Pagination
	if limit > 0 {
		tx = tx.Limit(limit).Offset(offset)
	}

	if err := tx.Find(&technologies).Error; err != nil {
		return nil, 0, err
	}

	return technologies, int(total), nil
}

func (r *Repository) FindTechnologyColumnValues(
	ctx context.Context, column string, sorts []sqlutil.SortItem, //collectionIDs []string,
) ([]string, error) {
	var values []string

	db := r.db.WithContext(ctx)
	tx := db.Model(&model.Technology{}).
		Select(fmt.Sprintf("DISTINCT %s", column)).
		Where(sqlutil.NotEqualClause(column), "")

	// if len(collectionIDs) > 0 {
	// 	subQuery := db.Model(&model.TechnologyCollection{}).
	// 		Select("1").
	// 		Where(
	// 			fmt.Sprintf("%s.%s = %s.%s",
	// 				model.TechnologyCollectionsTableName,
	// 				model.TechnologyCollection_TECHNOLOGY_ID_COLUMN,
	// 				model.TechnologyTableName,
	// 				model.Technology_ID_COLUMN),
	// 		).Where(
	// 		sqlutil.InClauseWithTable(model.TechnologyCollectionsTableName, model.TechnologyCollection_COLLECTION_ID_COLUMN),
	// 		collectionIDs,
	// 	)

	// 	tx = tx.Where(
	// 		sqlutil.ExistClause(), subQuery,
	// 	)
	// }

	if len(sorts) > 0 {
		for _, item := range sorts {
			tx = tx.Order(fmt.Sprintf("%s %s", item.Field, item.SortValue))
		}
	} else {
		tx = tx.Order(fmt.Sprintf("%s %s", column, sqlutil.SortValue_ASC))
	}

	if err := tx.Pluck(column, &values).Error; err != nil {
		return nil, err
	}

	return values, nil
}
