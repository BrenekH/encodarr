package sqlite

import "github.com/BrenekH/encodarr/controller"

func NewLibraryManagerAdapter(db *SQLiteDatabase, logger controller.Logger) LibraryManagerAdapter {
	return LibraryManagerAdapter{
		db:     db,
		logger: logger,
	}
}

type LibraryManagerAdapter struct {
	db     *SQLiteDatabase
	logger controller.Logger
}

func (l *LibraryManagerAdapter) Libraries() (libSlice []controller.Library) {
	l.logger.Critical("Not implemented")
	// TODO: Implement
	return
}
