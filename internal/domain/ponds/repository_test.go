package ponds

import (
	"context"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/nmluci/da-farm-be/internal/core/errs"
)

func TestShouldGetPondWithResult(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open stub connection, err: %s", err)
	}
	defer db.Close()

	pondRepo := NewRepository(sqlx.NewDb(db, "sqlmock"))

	// expected queries
	rows := sqlmock.NewRows([]string{"p.id", "farm_id", "p.name", "farm_name"}).
		AddRow(1, 1, "Pond A", "Farm A").
		AddRow(2, 1, "Pond B", "Farm A")

	mock.ExpectQuery(regexp.QuoteMeta("SELECT p.id, f.id farm_id, p.name, f.name farm_name FROM ponds p LEFT JOIN farms f on p.farm_id = f.id WHERE (p.farm_id = $1 AND f.deleted_at IS NULL AND p.deleted_at IS NULL) LIMIT 100 OFFSET 0")).
		WithArgs(1).
		WillReturnRows(rows)

	pondRepo.GetAll(context.Background(), &pondQuery{FarmID: 1, Limit: 100, Page: 1})

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("%s", err)
	}
}

func TestShouldCountPondAboveZero(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open stub connection, err: %s", err)
	}
	defer db.Close()

	pondRepo := NewRepository(sqlx.NewDb(db, "sqlmock"))

	// expected queries
	rows := sqlmock.NewRows([]string{"count(*)"}).
		AddRow(1)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT count(*) FROM ponds p LEFT JOIN farms f on p.farm_id = f.id WHERE (p.farm_id = $1 AND f.deleted_at IS NULL AND p.deleted_at IS NULL)")).
		WithArgs(1).
		WillReturnRows(rows)

	pondRepo.Count(context.Background(), &pondQuery{FarmID: 1, Limit: 100, Page: 1})

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("%s", err)
	}
}

func TestShouldCountPondLessThanZero(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open stub connection, err: %s", err)
	}
	defer db.Close()

	pondRepo := NewRepository(sqlx.NewDb(db, "sqlmock"))

	// expected queries
	rows := sqlmock.NewRows([]string{"count(*)"}).
		AddRow(0)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT count(*) FROM ponds p LEFT JOIN farms f on p.farm_id = f.id WHERE (p.farm_id = $1 AND f.deleted_at IS NULL AND p.deleted_at IS NULL)")).
		WithArgs(1).
		WillReturnRows(rows)

	pondRepo.Count(context.Background(), &pondQuery{FarmID: 1, Limit: 100, Page: 1})

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("%s", err)
	}
}

func TestShouldGetPond(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open stub connection, err: %s", err)
	}
	defer db.Close()

	pondRepo := NewRepository(sqlx.NewDb(db, "sqlmock"))

	// expected queries
	rows := sqlmock.NewRows([]string{"p.id", "farm_id", "p.name", "farm_name"}).
		AddRow(1, 1, "Pond A", "Farm A")

	mock.ExpectQuery(regexp.QuoteMeta("SELECT p.id, f.id farm_id, p.name, f.name farm_name FROM ponds p LEFT JOIN farms f on p.farm_id = f.id WHERE (p.id = $1 AND p.farm_id = $2 AND f.deleted_at IS NULL AND p.deleted_at IS NULL)")).
		WithArgs(1, 1).
		WillReturnRows(rows)

	pondRepo.GetOne(context.Background(), &pondQuery{ID: 1, FarmID: 1})

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("%s", err)
	}
}

func TestShouldNOTGetPond(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open stub connection, err: %s", err)
	}
	defer db.Close()

	pondRepo := NewRepository(sqlx.NewDb(db, "sqlmock"))

	// expected queries
	rows := sqlmock.NewRows([]string{"p.id", "farm_id", "p.name", "farm_name"})

	mock.ExpectQuery(regexp.QuoteMeta("SELECT p.id, f.id farm_id, p.name, f.name farm_name FROM ponds p LEFT JOIN farms f on p.farm_id = f.id WHERE (p.id = $1 AND p.farm_id = $2 AND f.deleted_at IS NULL AND p.deleted_at IS NULL)")).
		WithArgs(1, 2).
		WillReturnRows(rows)

	pondRepo.GetOne(context.Background(), &pondQuery{ID: 1, FarmID: 2})

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("%s", err)
	}
}

func TestShouldStorePond(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open stub connection, err: %s", err)
	}
	defer db.Close()

	pondRepo := NewRepository(sqlx.NewDb(db, "sqlmock"))

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta("SELECT count(*) FROM farms WHERE (id = $1 AND deleted_at IS NULL)")).WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"count(*)"}).AddRow(1))
	mock.ExpectQuery(regexp.QuoteMeta("SELECT count(*) FROM ponds WHERE (name = $1 AND deleted_at IS NULL)")).WithArgs("Pond A").
		WillReturnRows(sqlmock.NewRows([]string{"count(*)"}).AddRow(0))
	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO ponds (farm_id,name) VALUES ($1,$2)")).WithArgs(1, "Pond A").
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	pondRepo.Store(context.Background(), &PondType{FarmID: 1, Name: "Pond A"})

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("%s", err)
	}
}

func TestShouldNOTStorePondDueFarmNotExisted(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open stub connection, err: %s", err)
	}
	defer db.Close()

	pondRepo := NewRepository(sqlx.NewDb(db, "sqlmock"))

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta("SELECT count(*) FROM farms WHERE (id = $1 AND deleted_at IS NULL)")).WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"count(*)"}).AddRow(0)).WillReturnError(errs.ErrNotFound)
	mock.ExpectRollback()

	pondRepo.Store(context.Background(), &PondType{FarmID: 1, Name: "Pond A"})

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("%s", err)
	}
}

func TestShouldNOTStorePondDuePondDuplicated(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open stub connection, err: %s", err)
	}
	defer db.Close()

	pondRepo := NewRepository(sqlx.NewDb(db, "sqlmock"))

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta("SELECT count(*) FROM farms WHERE (id = $1 AND deleted_at IS NULL)")).WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"count(*)"}).AddRow(1))
	mock.ExpectQuery(regexp.QuoteMeta("SELECT count(*) FROM ponds WHERE (name = $1 AND deleted_at IS NULL)")).WithArgs("Pond A").
		WillReturnRows(sqlmock.NewRows([]string{"count(*)"}).AddRow(1)).WillReturnError(errs.ErrDuplicatedResources)
	mock.ExpectRollback()

	pondRepo.Store(context.Background(), &PondType{FarmID: 1, Name: "Pond A"})

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("%s", err)
	}
}

func TestShouldInsertOnUpdatePond(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open stub connection, err: %s", err)
	}
	defer db.Close()

	pondRepo := NewRepository(sqlx.NewDb(db, "sqlmock"))

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta("SELECT count(*) FROM farms WHERE (id = $1 AND deleted_at IS NULL)")).WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"count(*)"}).AddRow(1))
	mock.ExpectQuery(regexp.QuoteMeta("SELECT count(*) FROM ponds WHERE (id <> $1 AND name = $2 AND deleted_at IS NULL)")).WithArgs(1, "Pond A").
		WillReturnRows(sqlmock.NewRows([]string{"count(*)"}).AddRow(0))
	mock.ExpectQuery(regexp.QuoteMeta("SELECT count(*) FROM ponds WHERE (id = $1 AND deleted_at IS NULL)")).WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"count(*)"}).AddRow(0))
	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO ponds (farm_id,name) VALUES ($1,$2)")).WithArgs(1, "Pond A").
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	pondRepo.Upsert(context.Background(), &PondType{ID: 1, FarmID: 1, Name: "Pond A"})

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("%s", err)
	}
}

func TestShouldNOTInsertOnUpdatePondDueFarmNotExisted(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open stub connection, err: %s", err)
	}
	defer db.Close()

	pondRepo := NewRepository(sqlx.NewDb(db, "sqlmock"))

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta("SELECT count(*) FROM farms WHERE (id = $1 AND deleted_at IS NULL)")).WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"count(*)"}).AddRow(0)).WillReturnError(errs.ErrNotFound)
	mock.ExpectRollback()

	pondRepo.Upsert(context.Background(), &PondType{ID: 1, FarmID: 1, Name: "Pond A"})

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("%s", err)
	}
}

func TestShouldNOTInsertOnUpdatePondDuePondDuplicated(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open stub connection, err: %s", err)
	}
	defer db.Close()

	pondRepo := NewRepository(sqlx.NewDb(db, "sqlmock"))

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta("SELECT count(*) FROM farms WHERE (id = $1 AND deleted_at IS NULL)")).WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"count(*)"}).AddRow(1))
	mock.ExpectQuery(regexp.QuoteMeta("SELECT count(*) FROM ponds WHERE (id <> $1 AND name = $2 AND deleted_at IS NULL)")).WithArgs(1, "Pond A").
		WillReturnRows(sqlmock.NewRows([]string{"count(*)"}).AddRow(1))
	mock.ExpectRollback()

	pondRepo.Upsert(context.Background(), &PondType{ID: 1, FarmID: 1, Name: "Pond A"})

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("%s", err)
	}
}

func TestShouldUpdatePond(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open stub connection, err: %s", err)
	}
	defer db.Close()

	pondRepo := NewRepository(sqlx.NewDb(db, "sqlmock"))

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta("SELECT count(*) FROM farms WHERE (id = $1 AND deleted_at IS NULL)")).WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"count(*)"}).AddRow(1))
	mock.ExpectQuery(regexp.QuoteMeta("SELECT count(*) FROM ponds WHERE (id <> $1 AND name = $2 AND deleted_at IS NULL)")).WithArgs(1, "Pond A").
		WillReturnRows(sqlmock.NewRows([]string{"count(*)"}).AddRow(0))
	mock.ExpectQuery(regexp.QuoteMeta("SELECT count(*) FROM ponds WHERE (id = $1 AND deleted_at IS NULL)")).WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"count(*)"}).AddRow(1))
	mock.ExpectExec(regexp.QuoteMeta("UPDATE ponds SET name = $1, updated_at = NOW() WHERE (id = $2)")).WithArgs("Pond A", 1).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	pondRepo.Upsert(context.Background(), &PondType{ID: 1, FarmID: 1, Name: "Pond A"})

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("%s", err)
	}
}

func TestShouldDeletePond(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open stub connection, err: %s", err)
	}
	defer db.Close()

	pondRepo := NewRepository(sqlx.NewDb(db, "sqlmock"))

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta("SELECT count(*) FROM ponds WHERE (id = $1 AND deleted_at IS NULL)")).WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"count(*)"}).AddRow(1))
	mock.ExpectExec(regexp.QuoteMeta("UPDATE ponds SET deleted_at = NOW(), updated_at = NOW() WHERE id = $1")).WithArgs(1).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	pondRepo.Delete(context.Background(), &pondQuery{ID: 1})

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("%s", err)
	}
}

func TestShouldNOTDeletePondDuePondNotExisted(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open stub connection, err: %s", err)
	}
	defer db.Close()

	pondRepo := NewRepository(sqlx.NewDb(db, "sqlmock"))

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta("SELECT count(*) FROM ponds WHERE (id = $1 AND deleted_at IS NULL)")).WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"count(*)"}).AddRow(0)).WillReturnError(errs.ErrNotFound)
	mock.ExpectRollback()

	pondRepo.Delete(context.Background(), &pondQuery{ID: 1})

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("%s", err)
	}
}
