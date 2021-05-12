package sqlite

import (
	_ "modernc.org/sqlite"
)

type SQLiteDataStore struct{}

func NewSQLiteDataStore(file string) SQLiteDataStore {
	return SQLiteDataStore{}
}
