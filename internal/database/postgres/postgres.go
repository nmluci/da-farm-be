package postgres

import (
	"fmt"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog"
)

// New return sqlx.DB of newly established PostgreSQL connection
func New(logger zerolog.Logger, conf *PostgresConfig) (db *sqlx.DB, err error) {
	datasource := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable",
		conf.Username, conf.Password, conf.Address, conf.DB,
	)

	db, err = sqlx.Connect("pgx", datasource)
	if err != nil {
		logger.Error().Err(err).Msg("failed to connect to DB")
		return
	}

	logger.Info().Msg("db init successfully")

	// dbMigrate, err := migrate.New("file://migrations", datasource)
	// if err != nil {
	// 	logger.Error().Err(err).Msg("failed to connect to migration engine")
	// 	return
	// }

	// if err = dbMigrate.Up(); err != nil && err != migrate.ErrNoChange {
	// 	logger.Error().Err(err).Msg("failed to perform migrations")
	// 	return
	// }

	// rev, isDirty, err := dbMigrate.Version()
	// if err != nil && err != migrate.ErrNilVersion {
	// 	logger.Error().Err(err).Msg("failed to fetch migration version")
	// 	return
	// }

	// if isDirty {
	// 	logger.Warn().Msg("db migration is dirty")
	// }

	// if err == migrate.ErrNilVersion {
	// 	logger.Info().Msg("db migration Version: None")
	// } else {
	// 	logger.Info().Msgf("db migration Version: %d", rev)
	// }

	return
}
