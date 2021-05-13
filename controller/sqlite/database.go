package sqlite

import (
	"database/sql"

	_ "modernc.org/sqlite"
)

type SQLiteDatabase struct {
	Client *sql.DB
}

func NewSQLiteDatabase(configDir string) (SQLiteDatabase, error) {
	client, err := sql.Open("sqlite", configDir+"/data.db")
	if err != nil {
		return SQLiteDatabase{Client: client}, err
	}

	client.SetMaxOpenConns(1) // Set max connections to 1 to prevent "database is locked" errors

	// Ensure all tables are in the database
	_, err = client.Exec(schemaStmt)
	if err != nil {
		return SQLiteDatabase{Client: client}, err
	}

	return SQLiteDatabase{Client: client}, nil
}

var schemaStmt string = `
CREATE TABLE IF NOT EXISTS libraries (
	ID integer PRIMARY KEY,
	folder text,
	priority integer,
	fs_check_interval text,
	pipeline binary,
	queue binary,
	file_cache binary,
	path_masks binary
);
CREATE TABLE IF NOT EXISTS files (
	path text,
	modtime timestamp,
	mediainfo binary
);
CREATE TABLE IF NOT EXISTS history (
	time_completed timestamp,
	filename text,
	warnings binary,
	errors binary
);
CREATE TABLE IF NOT EXISTS dispatched_jobs (
	uuid text NOT NULL UNIQUE,
	job binary,
	status binary,
	runner text,
	last_updated timestamp
);
`
