package pgtestdb_test

import (
	"fmt"
	"io/fs"
	"os"
	"testing"

	pgtestdb "github.com/TudorHulban/pgtestdb"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/stretchr/testify/require"
)

func TestPGTestDB(t *testing.T) {
	pgTest := pgtestdb.PGTestDB{
		ConnectionURL: fmt.Sprintf(
			"postgres://%s:%s@%s:%s/%s?",
			"postgres",
			"password",
			"localhost",
			"5471",
			"",
		),

		MigrationDirectories: []fs.FS{
			os.DirFS("./migrations2"),
			os.DirFS("./migrations1"),
		},

		MigrationFilePaths: []string{
			"pgmigrator_test.sql",
		},

		T: t,

		// RegexValidationMigrationFile: sql.NullString{
		// 	Valid:  true,
		// 	String: `^V\d{4}_[a-zA-Z0-9]+\.sql$`,
		// },
	}

	// pgTest.Execute()
	dbName, cleanUp := pgTest.Execute()
	require.NotZero(t, dbName)
	defer cleanUp()

	// run test operations

	// debug onlu - cleanUp()
}
