package farms

import (
	"context"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/nmluci/da-farm-be/internal/core/errs"
)

func TestShouldGetFarmWithResult(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open stub connection, err: %s", err)
	}
	defer db.Close()

	farmRepo := NewRepository(sqlx.NewDb(db, "sqlmock"))

	// expected queries
	rows := sqlmock.NewRows([]string{"id", "name"}).
		AddRow(1, "Farm A").
		AddRow(2, "Farm B")

	mock.ExpectQuery("^SELECT (.+) FROM farms WHERE").WillReturnRows(rows)

	farmRepo.GetAll(context.Background(), &farmQuery{Limit: 100, Page: 1})

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("%s", err)
	}
}

func TestShouldCountFarmAboveZero(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open stub connection, err: %s", err)
	}
	defer db.Close()

	farmRepo := NewRepository(sqlx.NewDb(db, "sqlmock"))

	// expected queries
	rows := sqlmock.NewRows([]string{"count(*)"}).AddRow(0)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT count(*) FROM farms WHERE (deleted_at IS NULL)")).WillReturnRows(rows)

	farmRepo.Count(context.Background(), &farmQuery{Limit: 100, Page: 1})

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("%s", err)
	}
}

func TestShouldCountFarmLessThanZero(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open stub connection, err: %s", err)
	}
	defer db.Close()

	farmRepo := NewRepository(sqlx.NewDb(db, "sqlmock"))

	rows := sqlmock.NewRows([]string{"count(*)"}).AddRow(10)
	mock.ExpectQuery(regexp.QuoteMeta("SELECT count(*) FROM farms WHERE (deleted_at IS NULL)")).WillReturnRows(rows)

	farmRepo.Count(context.Background(), &farmQuery{Limit: 100, Page: 1})

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("%s", err)
	}
}

func TestShouldGetOneFarmWithID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open stub connection, err: %s", err)
	}
	defer db.Close()

	farmRepo := NewRepository(sqlx.NewDb(db, "sqlmock"))

	// expected queries
	rows := sqlmock.NewRows([]string{"id", "name"}).
		AddRow(1, "Farm A")

	mock.ExpectQuery(regexp.QuoteMeta("SELECT id, name FROM farms WHERE (id = $1 AND deleted_at IS NULL)")).WithArgs(1).WillReturnRows(rows)

	farmRepo.GetOne(context.Background(), &farmQuery{ID: 1})

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("%s", err)
	}
}

func TestShouldNOTGetOneFarmWithID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open stub connection, err: %s", err)
	}
	defer db.Close()

	farmRepo := NewRepository(sqlx.NewDb(db, "sqlmock"))

	// expected queries
	rows := sqlmock.NewRows([]string{"id", "name"})

	mock.ExpectQuery(regexp.QuoteMeta("SELECT id, name FROM farms WHERE (id = $1 AND deleted_at IS NULL)")).WithArgs(2).WillReturnRows(rows)

	farmRepo.GetOne(context.Background(), &farmQuery{ID: 2})

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("%s", err)
	}
}

func TestShouldStoreFarm(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open stub connection, err: %s", err)
	}
	defer db.Close()

	farmRepo := NewRepository(sqlx.NewDb(db, "sqlmock"))

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM farms WHERE (name = $1 AND deleted_at IS NULL)`)).WithArgs("Farm A").
		WillReturnRows(sqlmock.NewRows([]string{"count(*)"}).AddRow(0))
	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO farms (name) VALUES ($1)")).WithArgs("Farm A").
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	farmRepo.Store(context.Background(), &FarmType{Name: "Farm A"})

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("%s", err)
	}
}

func TestShouldInsertOnUpdateFarm(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open stub connection, err: %s", err)
	}
	defer db.Close()

	farmRepo := NewRepository(sqlx.NewDb(db, "sqlmock"))

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM farms WHERE (id <> $1 AND name = $2 AND deleted_at IS NULL)`)).WithArgs(1, "Farm A").
		WillReturnRows(sqlmock.NewRows([]string{"count(*)"}).AddRow(0))

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM farms WHERE (id = $1 AND deleted_at IS NULL)`)).WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"count(*)"}).AddRow(0))

	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO farms (name) VALUES ($1)")).WithArgs("Farm A").
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectCommit()

	farmRepo.Upsert(context.Background(), &FarmType{ID: 1, Name: "Farm A"})

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("%s", err)
	}
}

func TestShouldNOTInsertFarmOnUpdateDueDuplicated(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open stub connection, err: %s", err)
	}
	defer db.Close()

	farmRepo := NewRepository(sqlx.NewDb(db, "sqlmock"))

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM farms WHERE (id <> $1 AND name = $2 AND deleted_at IS NULL)`)).WithArgs(1, "Farm A").
		WillReturnRows(sqlmock.NewRows([]string{"count(*)"}).AddRow(0)).WillReturnError(errs.ErrDuplicatedResources)
	mock.ExpectRollback()

	farmRepo.Upsert(context.Background(), &FarmType{ID: 1, Name: "Farm A"})

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("%s", err)
	}
}

func TestShouldUpdateFarm(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open stub connection, err: %s", err)
	}
	defer db.Close()

	farmRepo := NewRepository(sqlx.NewDb(db, "sqlmock"))

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM farms WHERE (id <> $1 AND name = $2 AND deleted_at IS NULL)`)).WithArgs(1, "Farm A").
		WillReturnRows(sqlmock.NewRows([]string{"count(*)"}).AddRow(0))

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM farms WHERE (id = $1 AND deleted_at IS NULL)`)).WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"count(*)"}).AddRow(1))

	mock.ExpectExec(regexp.QuoteMeta("UPDATE farms SET name = $1, updated_at = NOW() WHERE id = $2")).WithArgs("Farm A", 1).
		WillReturnResult(sqlmock.NewResult(0, 1))

	mock.ExpectCommit()

	farmRepo.Upsert(context.Background(), &FarmType{ID: 1, Name: "Farm A"})

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("%s", err)
	}
}

func TestShouldNOTUpdateFarmDuplicated(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open stub connection, err: %s", err)
	}
	defer db.Close()

	farmRepo := NewRepository(sqlx.NewDb(db, "sqlmock"))

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM farms WHERE (id <> $1 AND name = $2 AND deleted_at IS NULL)`)).WithArgs(1, "Farm A").
		WillReturnRows(sqlmock.NewRows([]string{"count(*)"}).AddRow(1)).WillReturnError(errs.ErrDuplicatedResources)
	mock.ExpectRollback()

	farmRepo.Upsert(context.Background(), &FarmType{ID: 1, Name: "Farm A"})

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("%s", err)
	}
}

func TestShouldDeleteFarm(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open stub connection, err: %s", err)
	}
	defer db.Close()

	farmRepo := NewRepository(sqlx.NewDb(db, "sqlmock"))

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM farms WHERE (id = $1 AND deleted_at IS NULL)`)).WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"count(*)"}).AddRow(1))
	mock.ExpectExec(regexp.QuoteMeta("UPDATE farms SET deleted_at = NOW(), updated_at = NOW() WHERE id = $1")).WithArgs(1).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	farmRepo.Delete(context.Background(), &farmQuery{ID: 1})

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("%s", err)
	}
}

func TestShouldNOTDeleteFarmDueNotExisted(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open stub connection, err: %s", err)
	}
	defer db.Close()

	farmRepo := NewRepository(sqlx.NewDb(db, "sqlmock"))

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM farms WHERE (id = $1 AND deleted_at IS NULL)`)).WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"count(*)"}).AddRow(0)).WillReturnError(errs.ErrNotFound)
	mock.ExpectRollback()

	farmRepo.Delete(context.Background(), &farmQuery{ID: 1})

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("%s", err)
	}

}
