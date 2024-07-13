package ponds

import (
	"context"
	"database/sql"

	"github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"github.com/nmluci/da-farm-be/internal/core/errs"
	"github.com/rs/zerolog"
)

type PondRepository interface {
	GetAll(context.Context, *pondQuery) ([]*PondFarmType, error)
	Count(context.Context, *pondQuery) (uint64, error)
	GetOne(context.Context, *pondQuery) (*PondFarmType, error)
	Store(context.Context, *PondType) error
	Upsert(context.Context, *PondType) error
	Delete(context.Context, *pondQuery) error
}

type pondRepository struct {
	db *sqlx.DB
}

// NewRepository return an instance of pondRepository containing interface to DB layer
func NewRepository(db *sqlx.DB) PondRepository {
	return &pondRepository{db: db}
}

type pondQuery struct {
	ID, FarmID  int64
	Keyword     string
	Limit, Page uint64
}

var pgSquirrel = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

func (repo *pondRepository) GetAll(ctx context.Context, params *pondQuery) (res []*PondFarmType, err error) {
	logger := zerolog.Ctx(ctx)

	stmt, args, _ := pgSquirrel.Select("p.id", "f.id farm_id", "p.name", "f.name farm_name").From("ponds p").
		LeftJoin("farms f on p.farm_id = f.id").
		Where(squirrel.And{
			squirrel.Eq{"p.farm_id": params.FarmID},
			squirrel.Eq{"f.deleted_at": nil},
			squirrel.Eq{"p.deleted_at": nil},
		}).
		Limit(params.Limit).
		Offset((params.Page - 1) * params.Limit).ToSql()

	res = []*PondFarmType{}

	rows, err := repo.db.QueryxContext(ctx, stmt, args...)
	if err != nil {
		logger.Error().Err(err).Msg("failed to fetch data")
		return
	}

	for rows.Next() {
		col := &PondFarmType{}

		if err = rows.StructScan(col); err != nil {
			logger.Error().Err(err).Msg("failed to map row")
			return
		}

		res = append(res, col)
	}

	return
}

func (repo *pondRepository) Count(ctx context.Context, params *pondQuery) (res uint64, err error) {
	logger := zerolog.Ctx(ctx)

	stmt, args, _ := pgSquirrel.Select("count(*)").From("ponds p").
		LeftJoin("farms f on p.farm_id = f.id").
		Where(squirrel.And{
			squirrel.Eq{"p.farm_id": params.FarmID},
			squirrel.Eq{"f.deleted_at": nil},
			squirrel.Eq{"p.deleted_at": nil},
		}).ToSql()

	err = repo.db.QueryRowxContext(ctx, stmt, args...).Scan(&res)
	if err != nil && err != sql.ErrNoRows {
		logger.Error().Err(err).Msg("failed to fetch data")
		return
	} else if err == sql.ErrNoRows {
		return 0, nil
	}

	return
}

func (repo *pondRepository) GetOne(ctx context.Context, params *pondQuery) (res *PondFarmType, err error) {
	logger := zerolog.Ctx(ctx)

	stmt, args, _ := pgSquirrel.Select("p.id", "f.id farm_id", "p.name", "f.name farm_name").From("ponds p").
		LeftJoin("farms f on p.farm_id = f.id").
		Where(squirrel.And{
			squirrel.Eq{"p.id": params.ID},
			squirrel.Eq{"p.farm_id": params.FarmID},
			squirrel.Eq{"f.deleted_at": nil},
			squirrel.Eq{"p.deleted_at": nil},
		}).ToSql()

	err = repo.db.QueryRowxContext(ctx, stmt, args...).Scan(&res)
	if err != nil && err != sql.ErrNoRows {
		logger.Error().Err(err).Msg("failed to fetch data")
		return
	} else if err == sql.ErrNoRows {
		return nil, nil
	}

	return
}

func (repo *pondRepository) Store(ctx context.Context, payload *PondType) (err error) {
	logger := zerolog.Ctx(ctx)

	tx, err := repo.db.BeginTxx(ctx, nil)
	if err != nil {
		logger.Error().Err(err).Msg("failed to initialize transaction")
		return
	}
	defer tx.Rollback()

	var stmt string
	var args []any
	var count int64

	// check for farm existence
	stmt, args, _ = pgSquirrel.Select("count(*)").From("farms").Where(squirrel.And{
		squirrel.Eq{"id": payload.FarmID},
		squirrel.Eq{"deleted_at": nil},
	}).ToSql()

	if err = tx.QueryRowxContext(ctx, stmt, args...).Scan(&count); err != nil && err != sql.ErrNoRows { // make sure it's not an err from non-existing result
		logger.Error().Err(err).Msg("failed to validate farm data existence")
		return
	}

	// if selected farm doesn't exists, bail out from here
	if count == 0 {
		return errs.ErrNotFound
	}

	// check for duplicated name existence
	stmt, args, _ = pgSquirrel.Select("count(*)").From("ponds").Where(squirrel.And{
		squirrel.Eq{"name": payload.Name},
		squirrel.Eq{"deleted_at": nil},
	}).ToSql()

	if err = tx.QueryRowxContext(ctx, stmt, args...).Scan(&count); err != nil && err != sql.ErrNoRows { // make sure it's not an err from non-existing result
		logger.Error().Err(err).Msg("failed to validate duplicated data existence")
		return
	}

	// if active (non-deleted) row exist with such name, return duplicated err
	if count != 0 {
		return errs.ErrDuplicatedResources
	}

	stmt, args, _ = pgSquirrel.Insert("ponds").Columns("farm_id", "name").Values(payload.FarmID, payload.Name).ToSql()

	_, err = tx.ExecContext(ctx, stmt, args...)
	if err != nil {
		logger.Error().Err(err).Msg("failed to save data")
		return
	}

	if err = tx.Commit(); err != nil {
		logger.Error().Err(err).Msg("failed to commit transaction")
		return
	}

	return
}

func (repo *pondRepository) Upsert(ctx context.Context, payload *PondType) (err error) {
	logger := zerolog.Ctx(ctx)

	tx, err := repo.db.BeginTxx(ctx, nil)
	if err != nil {
		logger.Error().Err(err).Msg("failed to initialize transaction")
		return
	}
	defer tx.Rollback()

	var stmt string
	var args []any
	var count int64

	// check for farm existence
	stmt, args, _ = pgSquirrel.Select("count(*)").From("farms").Where(squirrel.And{
		squirrel.Eq{"id": payload.FarmID},
		squirrel.Eq{"deleted_at": nil},
	}).ToSql()

	if err = tx.QueryRowxContext(ctx, stmt, args...).Scan(&count); err != nil && err != sql.ErrNoRows { // make sure it's not an err from non-existing result
		logger.Error().Err(err).Msg("failed to validate farm data existence")
		return
	}

	// if selected farm doesn't exists, bail out from here
	if count == 0 {
		return errs.ErrNotFound
	}

	// check for duplicated name existence
	stmt, args, _ = pgSquirrel.Select("count(*)").From("ponds").Where(squirrel.And{
		squirrel.NotEq{"id": payload.ID},
		squirrel.Eq{"name": payload.Name},
		squirrel.Eq{"deleted_at": nil},
	}).ToSql()

	if err = tx.QueryRowxContext(ctx, stmt, args...).Scan(&count); err != nil && err != sql.ErrNoRows { // make sure it's not an err from non-existing result
		logger.Error().Err(err).Msg("failed to validate duplicated data existence")
		return
	}

	// if active (non-deleted) row exist with such name, return duplicated err
	if count != 0 {
		return errs.ErrDuplicatedResources
	}

	// check for ponds existence
	stmt, args, _ = pgSquirrel.Select("count(*)").From("ponds").Where(squirrel.And{
		squirrel.Eq{"id": payload.ID},
		squirrel.Eq{"deleted_at": nil},
	}).ToSql()

	if err = tx.QueryRowxContext(ctx, stmt, args...).Scan(&count); err != nil && err != sql.ErrNoRows { // make sure it's not an err from non-existing result
		logger.Error().Err(err).Msg("failed to validate pond data existence")
		return
	}

	switch count {
	case 0:
		stmt, args, _ = pgSquirrel.Insert("ponds").Columns("farm_id", "name").Values(payload.FarmID, payload.Name).ToSql()
	default:
		stmt, args, _ = pgSquirrel.Update("ponds").SetMap(map[string]interface{}{
			"name":       payload.Name,
			"updated_at": squirrel.Expr("NOW()"),
		}).Where(squirrel.And{
			squirrel.Eq{"id": payload.ID},
		}).ToSql()
	}

	_, err = tx.ExecContext(ctx, stmt, args...)
	if err != nil {
		logger.Error().Err(err).Msg("failed to update data")
		return
	}

	if err = tx.Commit(); err != nil {
		logger.Error().Err(err).Msg("failed to commit transaction")
		return
	}

	return
}

func (repo *pondRepository) Delete(ctx context.Context, params *pondQuery) (err error) {
	logger := zerolog.Ctx(ctx)

	tx, err := repo.db.BeginTxx(ctx, nil)
	if err != nil {
		logger.Error().Err(err).Msg("failed to initialize a transaction")
		return
	}
	defer tx.Rollback()

	var stmt string
	var args []any

	// check for row existence
	stmt, args, _ = pgSquirrel.Select("count(*)").From("ponds").Where(squirrel.And{
		squirrel.Eq{"id": params.ID},
		squirrel.Eq{"deleted_at": nil},
	}).ToSql()

	var count int64
	if err = tx.QueryRowxContext(ctx, stmt, args...).Scan(&count); err != nil && err != sql.ErrNoRows {
		logger.Error().Err(err).Msg("failed to fetch data")
		return
	}

	if count == 0 {
		err = errs.ErrNotFound
		logger.Error().Err(err).Msg("pond doesn't exists")
		return
	}

	stmt, args, _ = pgSquirrel.Update("ponds").SetMap(map[string]interface{}{
		"updated_at": squirrel.Expr("NOW()"),
		"deleted_at": squirrel.Expr("NOW()"),
	}).Where(squirrel.Eq{"id": params.ID}).ToSql()

	_, err = tx.ExecContext(ctx, stmt, args...)
	if err != nil {
		logger.Error().Err(err).Msg("failed to delete data")
		return
	}

	if err = tx.Commit(); err != nil {
		logger.Error().Err(err).Msg("failed to commit transaction")
		return
	}

	return
}
