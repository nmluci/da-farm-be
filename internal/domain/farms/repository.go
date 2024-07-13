package farms

import (
	"context"
	"database/sql"
	"errors"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jmoiron/sqlx"
	"github.com/nmluci/da-farm-be/internal/core/errs"
	"github.com/rs/zerolog"
)

// FarmRepository contain contract that defined all necessary public function available to be interact with
type FarmRepository interface {
	GetAll(context.Context, *farmQuery) ([]*FarmType, error)
	Count(context.Context, *farmQuery) (uint64, error)
	GetOne(context.Context, *farmQuery) (*FarmType, error)
	Store(context.Context, *FarmType) error
	Upsert(context.Context, *FarmType) error
	Delete(context.Context, *farmQuery) error
}

type farmRepository struct {
	db *sqlx.DB
}

// NewRepository return a instance of farmRepository containing interface to DB layer
func NewRepository(db *sqlx.DB) FarmRepository {
	return &farmRepository{
		db: db,
	}
}

type farmQuery struct {
	ID          int64
	Keyword     string
	Limit, Page uint64
}

var pgSquirrel = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

func (repo *farmRepository) GetAll(ctx context.Context, params *farmQuery) (res []*FarmType, err error) {
	logger := zerolog.Ctx(ctx)

	stmt, args, _ := pgSquirrel.Select("id", "name").From("farms").
		Where(squirrel.And{
			squirrel.Eq{"deleted_at": nil},
		}).
		Limit(params.Limit).
		Offset((params.Page - 1) * params.Limit).ToSql()

	res = []*FarmType{}
	rows, err := repo.db.QueryxContext(ctx, stmt, args...)
	if err != nil {
		logger.Error().Err(err).Msg("failed to fetch data")
		return
	}

	for rows.Next() {
		col := &FarmType{}

		if err = rows.StructScan(col); err != nil {
			logger.Error().Err(err).Msg("failed to mapped row")
			return
		}

		res = append(res, col)
	}

	return
}

func (repo *farmRepository) Count(ctx context.Context, params *farmQuery) (res uint64, err error) {
	logger := zerolog.Ctx(ctx)

	stmt, args, _ := pgSquirrel.Select("count(*)").From("farms").
		Where(squirrel.And{
			squirrel.Eq{"deleted_at": nil},
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

func (repo *farmRepository) GetOne(ctx context.Context, params *farmQuery) (res *FarmType, err error) {
	logger := zerolog.Ctx(ctx)

	stmt, args, _ := pgSquirrel.Select("id", "name").From("farms").
		Where(squirrel.And{
			squirrel.Eq{"id": params.ID},
			squirrel.Eq{"deleted_at": nil},
		}).ToSql()

	res = &FarmType{}
	err = repo.db.QueryRowxContext(ctx, stmt, args...).StructScan(res)
	if err != nil && err != sql.ErrNoRows {
		logger.Error().Err(err).Msg("failed to fetch data")
		return
	} else if err == sql.ErrNoRows {
		return nil, nil
	}

	return
}

func (repo *farmRepository) Store(ctx context.Context, payload *FarmType) (err error) {
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

	// check for duplicated name existence
	stmt, args, _ = pgSquirrel.Select("count(*)").From("farms").Where(squirrel.And{
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

	stmt, args, _ = pgSquirrel.Insert("farms").
		Columns("name").
		Values(payload.Name).ToSql()

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

// Upsert will update an column if any matched, otherwise create a new one
func (repo *farmRepository) Upsert(ctx context.Context, payload *FarmType) (err error) {
	logger := zerolog.Ctx(ctx)

	tx, err := repo.db.BeginTxx(ctx, nil)
	if err != nil {
		logger.Error().Err(err).Msg("failed to initialize a transaction")
		return
	}
	defer tx.Rollback()

	var stmt string
	var args []any
	var count int64

	// check for duplicated name existence
	stmt, args, _ = pgSquirrel.Select("count(*)").From("farms").Where(squirrel.And{
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

	// check for farm existence
	stmt, args, _ = pgSquirrel.Select("count(*)").From("farms").Where(squirrel.And{
		squirrel.Eq{"id": payload.ID},
		squirrel.Eq{"deleted_at": nil},
	}).ToSql()

	if err = tx.QueryRowxContext(ctx, stmt, args...).Scan(&count); err != nil && err != sql.ErrNoRows { // make sure it's not an err from non-existing result
		logger.Error().Err(err).Msg("failed to validate current farm existence data")
		return
	}

	switch count {
	case 0:
		stmt, args, _ = pgSquirrel.Insert("farms").Columns("name").Values(payload.Name).ToSql()
	default:
		stmt, args, _ = pgSquirrel.Update("farms").SetMap(map[string]interface{}{
			"name":       payload.Name,
			"updated_at": squirrel.Expr("NOW()"),
		}).Where(squirrel.Eq{"id": payload.ID}).ToSql()
	}

	_, err = tx.ExecContext(ctx, stmt, args...)
	if err != nil {
		// if DB return an Unique Violation err, then there's duplicated data
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			err = errs.ErrDuplicatedResources
			logger.Error().Err(err).Msg("failed to update due to duplicated data")
			return
		}

		logger.Error().Err(err).Msg("failed to update data")
		return
	}

	if err = tx.Commit(); err != nil {
		logger.Error().Err(err).Msg("failed to commit transaction")
		return
	}

	return
}

func (repo *farmRepository) Delete(ctx context.Context, payload *farmQuery) (err error) {
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
	stmt, args, _ = pgSquirrel.Select("count(*)").From("farms").Where(squirrel.And{
		squirrel.Eq{"id": payload.ID},
		squirrel.Eq{"deleted_at": nil},
	}).ToSql()

	var count int64
	if err = tx.QueryRowxContext(ctx, stmt, args...).Scan(&count); err != nil && err != sql.ErrNoRows {
		logger.Error().Err(err).Msg("failed to fetch data")
		return
	}

	if count == 0 {
		err = errs.ErrNotFound
		logger.Error().Err(err).Msg("farm doesn't exists")
		return
	}

	stmt, args, _ = pgSquirrel.Update("farms").SetMap(map[string]interface{}{
		"updated_at": squirrel.Expr("NOW()"),
		"deleted_at": squirrel.Expr("NOW()"),
	}).Where(squirrel.Eq{"id": payload.ID}).ToSql()

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
