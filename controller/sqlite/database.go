package sqlite

import (
	"database/sql"
	"embed"
	"io"
	"os"

	_ "modernc.org/sqlite"

	"github.com/BrenekH/encodarr/controller"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/sqlite"
	"github.com/golang-migrate/migrate/v4/source/iofs"
)

//go:embed migrations
var migrations embed.FS

const targetMigrationVersion uint = 2

type SQLiteDatabase struct {
	Client *sql.DB
}

func NewSQLiteDatabase(configDir string, logger controller.Logger) (SQLiteDatabase, error) {
	dbFile := configDir + "/data.db"
	dbBckpFile := configDir + "/data.db.backup"

	client, err := sql.Open("sqlite", dbFile)
	if err != nil {
		return SQLiteDatabase{Client: client}, err
	}

	// Set max connections to 1 to prevent "database is locked" errors
	client.SetMaxOpenConns(1)

	dbBackup, err := os.Create(dbBckpFile)
	if err != nil {
		return SQLiteDatabase{Client: client}, err
	}

	err = gotoDBVer(dbFile, targetMigrationVersion, dbBackup, logger)

	return SQLiteDatabase{Client: client}, err
}

// gotoDBVer uses github.com/golang-migrate/migrate to move the db version up or down to the passed target version.
func gotoDBVer(dbFile string, targetVersion uint, backupWriter io.Writer, logger controller.Logger) error {
	// TODO: Solve issue where embed won't include newer down files for downgrading
	migrationsSource, err := iofs.New(migrations, "migrations")
	if err != nil {
		return err
	}
	defer migrationsSource.Close()

	mig, err := migrate.NewWithSourceInstance("file://migrations", migrationsSource, "sqlite://"+dbFile)
	if err != nil {
		return err
	}
	defer mig.Close()

	currentVer, _, err := mig.Version()
	if err != nil {
		if err == migrate.ErrNilVersion {
			// DB is likely before golang-migrate was introduced. Upgrade to new version
			logger.Warn("Database does not have a schema version. Attempting to migrate up.")
			err = backupFile(dbFile, backupWriter, logger)
			if err != nil {
				return err
			}

			return mig.Migrate(targetVersion)
		}
		return err
	}

	if currentVer == targetVersion {
		return nil
	}

	err = backupFile(dbFile, backupWriter, logger)
	if err != nil {
		return err
	}

	logger.Info("Migrating database to schema version %v", targetVersion)
	return mig.Migrate(targetVersion)
}

func backupFile(from string, to io.Writer, logger controller.Logger) error {
	fromReader, err := os.Open(from)
	if err != nil {
		return err
	}

	logger.Warn("Backing up database")
	_, err = io.Copy(to, fromReader)
	return err
}
