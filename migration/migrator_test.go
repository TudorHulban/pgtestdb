package migration

import (
	"database/sql"
	"fmt"
	"io/fs"
	"os"
	"testing"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/stretchr/testify/require"
)

type paramsDBConnection struct {
	DBHost     string
	DBPort     string
	DBName     string
	DBUser     string
	DBPassword string
}

func TestMigrator(t *testing.T) {
	pgMigrator := NewPGMigrator(
		&ParamsNewPGMigrator{
			Directories: []fs.FS{
				os.DirFS("../migrations2"),
				os.DirFS("../migrations1"),
			},
			FilePaths: []string{
				"../pgmigrator_test.sql",
			},

			T: t,
		},
	)
	require.NotNil(t, pgMigrator)
	require.Len(t, pgMigrator.migrations, 3)

	params := paramsDBConnection{
		DBHost: "localhost",
		DBPort: "5471",

		DBUser:     "postgres",
		DBPassword: "password",
	}

	db, errOpen := sql.Open(
		"pgx",
		fmt.Sprintf(
			"postgres://%s:%s@%s:%s/%s?",
			params.DBUser,
			params.DBPassword,
			params.DBHost,
			params.DBPort,
			params.DBName,
		),
	)
	require.NoError(t, errOpen)
	require.NotNil(t, db)

	defer db.Close()

	pgMigrator.Migrate(db)
}
