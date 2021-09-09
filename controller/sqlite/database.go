package sqlite

import (
	"database/sql"
	"embed"
	"errors"
	"io"
	"os"

	_ "modernc.org/sqlite" // The SQLite database driver

	"github.com/BrenekH/encodarr/controller"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/sqlite" // Add the sqlite database source to golang-migrate
	_ "github.com/golang-migrate/migrate/v4/source/file"     // Add the file migrations source to golang-migrate
)

//go:embed migrations
var migrations embed.FS

const targetMigrationVersion uint = 2

// Database is a wrapper around the database driver client
type Database struct {
	Client *sql.DB
}

// NewDatabase returns an instantiated SQLiteDatabase.
func NewDatabase(configDir string, logger controller.Logger) (Database, error) {
	dbFilename := configDir + "/data.db"
	dbBackupFilename := configDir + "/data.db.backup"

	client, err := sql.Open("sqlite", dbFilename)
	if err != nil {
		return Database{Client: client}, err
	}

	// Set max connections to 1 to prevent "database is locked" errors
	client.SetMaxOpenConns(1)

	err = gotoDBVer(dbFilename, targetMigrationVersion, configDir, dbBackupFilename, logger)

	return Database{Client: client}, err
}

// gotoDBVer uses github.com/golang-migrate/migrate to move the db version up or down to the passed target version.
func gotoDBVer(dbFilename string, targetVersion uint, configDir string, backupFilename string, logger controller.Logger) error {
	// Instead of directly using the embedded files, write them out to {configDir}/migrations. This allows the files for downgrading the
	// database to be present even when the executable doesn't contain them.
	fsMigrationsDir := configDir + "/migrations"

	if err := os.MkdirAll(fsMigrationsDir, 0777); err != nil {
		return err
	}

	dirEntries, err := migrations.ReadDir("migrations")
	if err != nil {
		return err
	}

	var copyErred bool
	for _, v := range dirEntries {
		f, err := os.Create(fsMigrationsDir + "/" + v.Name())
		if err != nil {
			logger.Error("%v", err)
			copyErred = true
			continue
		}

		embeddedFile, err := migrations.Open("migrations/" + v.Name())
		if err != nil {
			logger.Error("%v", err)
			copyErred = true
			f.Close()
			continue
		}

		if _, err := io.Copy(f, embeddedFile); err != nil {
			logger.Error("%v", err)
			copyErred = true
			// Don't continue right here so that the files are closed before looping again
		}

		f.Close()
		embeddedFile.Close()
	}
	if copyErred {
		return errors.New("error(s) while copying migrations, check logs for more details")
	}

	mig, err := migrate.New("file://"+configDir+"/migrations", "sqlite://"+dbFilename)
	if err != nil {
		return err
	}
	defer mig.Close()

	currentVer, _, err := mig.Version()
	if err != nil {
		if err == migrate.ErrNilVersion {
			// DB is likely before golang-migrate was introduced. Upgrade to new version
			logger.Warn("Database does not have a schema version. Attempting to migrate up.")
			err = backupFile(dbFilename, backupFilename, logger)
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

	err = backupFile(dbFilename, backupFilename, logger)
	if err != nil {
		return err
	}

	logger.Info("Migrating database to schema version %v.", targetVersion)
	return mig.Migrate(targetVersion)
}

// backupFile backups a file to an io.Writer and logs about it.
func backupFile(from, to string, logger controller.Logger) error {
	fromReader, err := os.Open(from)
	if err != nil {
		return err
	}

	toWriter, err := os.Create(to)
	if err != nil {
		return err
	}

	logger.Info("Backing up database.")
	_, err = io.Copy(toWriter, fromReader)
	return err
}
