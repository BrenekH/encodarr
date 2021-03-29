package db

import (
	"database/sql"

	"github.com/BrenekH/encodarr/controller/options"
	"github.com/BrenekH/logange"

	// Load SQLite driver
	_ "github.com/mattn/go-sqlite3"
)

// Client is a "database/sql" DB pointer for access to the SQLite database
var Client *sql.DB

var logger logange.Logger

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
	queued integer
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

func init() {
	// Setup logger
	logger = logange.NewLogger("database")

	// Setup SQLite database
	var err error
	Client, err = sql.Open("sqlite3", options.ConfigDir()+"/data.db")
	if err != nil {
		logger.Critical(err.Error())
	}

	Client.SetMaxOpenConns(1) // We set the max connections to 1, to prevent any "database is locked" errors

	_, err = Client.Exec(schemaStmt)
	if err != nil {
		logger.Critical(err.Error())
	}

	logger.Debug("Database setup and ready")
}
