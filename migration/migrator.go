package migration

import (
	"database/sql"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

const _MigrationsTableName = "aaa_migrations"
const _MsgErrorRollback = "rollback executed"
const _DefaultValidationRegex = `^V\d{4}__[a-zA-Z0-9_]+\.sql$`

type PGMigrator struct {
	migrationsTableName          string
	regexValidationMigrationFile sql.NullString
	migrations

	T *testing.T
}

type ParamsNewPGMigrator struct {
	MigrationsTableName          sql.NullString
	RegexValidationMigrationFile sql.NullString

	Directories []fs.FS
	FilePaths   []string

	T *testing.T
}

func NewPGMigrator(params *ParamsNewPGMigrator) *PGMigrator {
	var migrations migrations

	regex := _DefaultValidationRegex

	if params.RegexValidationMigrationFile.Valid {
		regex = params.RegexValidationMigrationFile.String
	}

	regexCompiled := regexp.MustCompile(regex)

	for _, directory := range params.Directories {
		buf, errLoad := load(
			directory,
			regexCompiled,
		)
		require.NoError(params.T, errLoad)

		migrations = append(migrations, buf...)
	}

	migrations.SortByID()

	for _, filePath := range params.FilePaths {
		content, errRead := os.ReadFile(filePath)
		require.NoError(
			params.T,
			errRead,
			fmt.Errorf("failed to read file: %w", errRead),
		)

		migrations = append(migrations,
			migration{
				ID:  filepath.Base(filePath),
				SQL: string(content),
			},
		)
	}

	migrationsTableName := _MigrationsTableName

	if params.MigrationsTableName.Valid {
		migrationsTableName = params.MigrationsTableName.String
	}

	return &PGMigrator{
		migrationsTableName: migrationsTableName,
		migrations:          migrations,
		T:                   params.T,
	}
}

type paramsApplyMigrationSQLDB struct {
	Tx        *sql.Tx
	Item      migration
	StartedAt time.Time
}

func (m *PGMigrator) applyMigration(params *paramsApplyMigrationSQLDB) error {
	if _, errExecSQL := params.Tx.Exec(params.Item.SQL); errExecSQL != nil {
		require.NoError(
			m.T,
			params.Tx.Rollback(),
			fmt.Sprintf(
				"execution failed for %s",
				params.Item.ID,
			),
		)

		return errors.New(
			_MsgErrorRollback,
		)
	}

	if _, errExecAudit := params.Tx.Exec(
		`insert into `+m.migrationsTableName+`(checksum,script,miliseconds_execution_time,applied_at,success) values($1,$2,$3,$4,$5);`,

		params.Item.MD5(),
		params.Item.ID,
		time.Since(params.StartedAt).Milliseconds(),
		time.Now().UTC(),
		true,
	); errExecAudit != nil {
		require.NoError(
			m.T,
			params.Tx.Rollback(),
			fmt.Sprintf(
				"execution failed audit insert for %s",
				params.Item.ID,
			),
		)

		return errors.New(
			_MsgErrorRollback,
		)
	}

	require.NoError(
		m.T,
		params.Tx.Commit(),
	)

	return nil
}

func (m *PGMigrator) ddlMigrationsTable() string {
	return `create table if not exists ` + m.migrationsTableName +
		` (id smallint generated by default as identity primary key,` +
		`script text not null,` +
		`checksum text not null,` +
		`miliseconds_execution_time bigint not null,` +
		`applied_at timestamptz not null,` +
		`success bool not null default false);`
}

func (m *PGMigrator) Migrate(db *sql.DB) {
	_, errExec := db.Exec(
		m.ddlMigrationsTable(),
	)
	require.NoError(
		m.T,
		errExec,

		fmt.Sprintf(
			"create migratons table: %s",
			m.migrationsTableName,
		),
	)

	for _, item := range m.migrations {
		startedAt := time.Now().UTC()

		tx, errTx := db.Begin()
		require.NoError(m.T, errTx)

		errMigration := m.applyMigration(
			&paramsApplyMigrationSQLDB{
				Tx:        tx,
				Item:      item,
				StartedAt: startedAt,
			},
		)
		if errMigration != nil {
			_, errAuditFailedOperation := db.Exec(
				`insert into `+m.migrationsTableName+`(checksum,script,miliseconds_execution_time,applied_at) values($1,$2,$3,$4);`,

				item.MD5(),
				item.ID,
				time.Since(startedAt).Milliseconds(),
				time.Now().UTC(),
			)

			require.NoError(m.T, errAuditFailedOperation)
		}

		require.NoError(
			m.T,

			errMigration,
		)
	}
}
