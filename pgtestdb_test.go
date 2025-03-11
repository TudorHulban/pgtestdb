package pgtestdb_test

import (
	"bytes"
	"fmt"
	"io/fs"
	"os"
	"testing"
	"text/template"

	pgtestdb "github.com/TudorHulban/pgtestdb"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/stretchr/testify/require"
)

func TestPGTestDB(t *testing.T) {
	f := func(filePath string) (string, error) {
		tpl, errParse := template.ParseFiles(filePath)
		if errParse != nil {
			return "", errParse
		}

		var buf bytes.Buffer

		if errExecute := tpl.Execute(
			&buf,
			map[string]string{
				"TableName": "tableD",
			},
		); errExecute != nil {
			return "", errExecute
		}

		return buf.String(), nil
	}

	pgTest := pgtestdb.PGTestDB{
		ConnectionURL: fmt.Sprintf(
			"postgres://%s:%s@%s:%s/%s?",
			"postgres",
			"password",
			"localhost",
			"5471",
			"tara_crm",
		),

		MigrationDirectories: []fs.FS{
			os.DirFS("./migrations2"),
			os.DirFS("./migrations1"),
		},

		MigrationFilePaths: []string{
			"pgmigrator_test.sql",
		},

		TemplateRenderFunction: f,

		TemplateFilePaths: []string{
			"pgmigrator_template_test.sql",
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
	// defer cleanUp() - up to when test becomes stable

	t.Log(dbName)

	// run test operations

	t.Cleanup(
		func() {
			cleanUp()
		},
	)
}
