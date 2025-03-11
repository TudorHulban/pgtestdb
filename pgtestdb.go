package pgtestdb

import (
	"database/sql"
	"fmt"
	"io/fs"
	"strings"
	"testing"
	"time"

	"github.com/TudorHulban/pgtestdb/migration"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/stretchr/testify/require"
)

// Connection in form of "postgres://%s:%s@%s:%s/%s?", db user, db password, host, port, db name.
type PGTestDB struct {
	ConnectionURL                string
	RegexValidationMigrationFile sql.NullString // default is `^V\d{4}__[a-zA-Z0-9_]+\.sql$`
	MigrationsTableName          sql.NullString // overrides default name

	TemplateRenderFunction func(string) (string, error)
	TemplateFilePaths      []string

	MigrationDirectories []fs.FS
	MigrationFilePaths   []string

	T *testing.T
}

// Execute performs migrations and returns created db name and cleanup function
// that should be run on test exit.
func (pg *PGTestDB) Execute() (string, func()) {
	dbCreate, errOpenCurrentConnection := sql.Open("pgx", pg.ConnectionURL)
	require.NoError(pg.T, errOpenCurrentConnection)
	require.NotNil(pg.T, dbCreate)

	dbName := fmt.Sprintf(
		`t%d__%s`,

		time.Now().Unix(),
		strings.ToLower(pg.T.Name()),
	)

	_, errCreateDB := dbCreate.Exec(
		fmt.Sprintf(
			`create database %s;`,
			dbName,
		),
	)
	require.NoError(pg.T, errCreateDB)

	connectionNewDB := updateDBNameInConnection(
		pg.ConnectionURL,
		dbName,
	)

	dbTest, errOpenNewDB := sql.Open("pgx", connectionNewDB)
	require.NoError(pg.T, errOpenNewDB)
	require.NotNil(pg.T, dbCreate)

	defer dbTest.Close()

	pgMigrator := migration.NewPGMigrator(
		&migration.ParamsNewPGMigrator{
			MigrationsTableName:          pg.MigrationsTableName,
			RegexValidationMigrationFile: pg.RegexValidationMigrationFile,

			TemplateFilePaths:      pg.TemplateFilePaths,
			TemplateRenderFunction: pg.TemplateRenderFunction,

			Directories: pg.MigrationDirectories,
			FilePaths:   pg.MigrationFilePaths,

			T: pg.T,
		},
	)
	require.NotNil(pg.T, pgMigrator)

	pgMigrator.Migrate(dbTest)

	return dbName,
		func() {
			_, errDropDB := dbCreate.Exec(
				fmt.Sprintf(
					`drop database %s;`,
					dbName,
				),
			)
			require.NoError(pg.T, errDropDB)

			dbCreate.Close()
		}
}
