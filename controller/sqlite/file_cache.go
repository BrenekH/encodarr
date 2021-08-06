package sqlite

import (
	"encoding/json"
	"time"

	"github.com/BrenekH/encodarr/controller"
)

// NewFileCacheAdapter returns an instantiated FileCacheAdapter.
func NewFileCacheAdapter(db *Database) FileCacheAdapter {
	return FileCacheAdapter{db: db}
}

// FileCacheAdapter satisfies the controller.FilesCacheDataStorer interface by turning interface
// requests into SQL requests that are passed on to an underlying SQLiteDatabase.
type FileCacheAdapter struct {
	db *Database
}

// Modtime uses a SQL SELECT statement to obtain the modtime associated with the provided path.
func (a *FileCacheAdapter) Modtime(path string) (time.Time, error) {
	row := a.db.Client.QueryRow("SELECT modtime FROM files WHERE path = $1;", path)

	var storedModtime time.Time

	err := row.Scan(&storedModtime)
	if err != nil {
		return time.Now(), err
	}

	return storedModtime, nil
}

// Metadata uses a SQL SELECT statement to obtain the metadata associated with the provided path.
func (a *FileCacheAdapter) Metadata(path string) (controller.FileMetadata, error) {
	row := a.db.Client.QueryRow("SELECT metadata FROM files WHERE path = $1;", path)

	var storedMetadataBytes []byte

	err := row.Scan(&storedMetadataBytes)
	if err != nil {
		return controller.FileMetadata{}, err
	}

	var storedMetadata controller.FileMetadata

	err = json.Unmarshal(storedMetadataBytes, &storedMetadata)
	if err != nil {
		return controller.FileMetadata{}, err
	}

	return storedMetadata, nil
}

// SaveModtime uses the UPSERT syntax to update the modtime that is associated with the provided path in the database.
func (a *FileCacheAdapter) SaveModtime(path string, t time.Time) error {
	_, err := a.db.Client.Exec("INSERT INTO files (path, modtime) VALUES ($1, $2) ON CONFLICT(path) DO UPDATE SET path=$1, modtime=$2;",
		path,
		t,
	)
	if err != nil {
		return err
	}

	return nil
}

// SaveMetadata uses the UPSERT syntax to update the metadata that is associated with the provided path in the database.
func (a *FileCacheAdapter) SaveMetadata(path string, f controller.FileMetadata) error {
	b, err := json.Marshal(f)
	if err != nil {
		return err
	}

	_, err = a.db.Client.Exec("INSERT INTO files (path, metadata) VALUES ($1, $2) ON CONFLICT(path) DO UPDATE SET path=$1, metadata=$2;",
		path,
		b,
	)
	if err != nil {
		return err
	}

	return nil
}
