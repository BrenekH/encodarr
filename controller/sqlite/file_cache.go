package sqlite

import (
	"encoding/json"
	"time"

	"github.com/BrenekH/encodarr/controller"
)

func NewFileCacheAdapter(db *SQLiteDatabase, logger controller.Logger) FileCacheAdapter {
	return FileCacheAdapter{db: db, logger: logger}
}

// FileCacheAdapter satisfies the controller.FilesCacheDataStorer interface by turning interface
// requests into SQL requests that are passed on to an underlying SQLiteDatabase.
type FileCacheAdapter struct {
	db     *SQLiteDatabase
	logger controller.Logger
}

func (a *FileCacheAdapter) Modtime(path string) (time.Time, error) {
	row := a.db.Client.QueryRow("SELECT modtime FROM files WHERE path = $1;", path)

	var storedModtime time.Time

	err := row.Scan(&storedModtime)
	if err != nil {
		a.logger.Error("%v", err)
		return time.Now(), err
	}

	return storedModtime, nil
}

func (a *FileCacheAdapter) Metadata(path string) (controller.FileMetadata, error) {
	row := a.db.Client.QueryRow("SELECT metadata FROM files WHERE path = $1;", path)

	var storedMetadataBytes []byte

	err := row.Scan(&storedMetadataBytes)
	if err != nil {
		a.logger.Error("%v", err)
		return controller.FileMetadata{}, err
	}

	var storedMetadata controller.FileMetadata

	err = json.Unmarshal(storedMetadataBytes, &storedMetadata)
	if err != nil {
		a.logger.Error("%v", err)
		return controller.FileMetadata{}, err
	}

	return storedMetadata, nil
}

func (a *FileCacheAdapter) SaveModtime(path string, t time.Time) error {
	_, err := a.db.Client.Exec("INSERT INTO files (path, modtime) VALUES ($1, $2) ON CONFLICT(path) DO UPDATE SET path=$1, modtime=$2;",
		path,
		t,
	)
	if err != nil {
		a.logger.Error(err.Error())
		return err
	}

	return nil
}

func (a *FileCacheAdapter) SaveMetadata(path string, f controller.FileMetadata) error {
	_, err := a.db.Client.Exec("INSERT INTO files (path, metadata) VALUES ($1, $2) ON CONFLICT(path) DO UPDATE SET path=$1, metadata=$2;",
		path,
		f,
	)
	if err != nil {
		a.logger.Error(err.Error())
		return err
	}

	return nil
}
