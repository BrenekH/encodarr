package sqlite

import (
	_ "modernc.org/sqlite"
)

type SQLiteDataStore struct{}

func NewSQLiteDataStore(configDir string) SQLiteDataStore {
	return SQLiteDataStore{}
}
