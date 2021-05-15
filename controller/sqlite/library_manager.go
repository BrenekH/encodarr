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

func (l *LibraryManagerAdapter) IsPathDispatched(path string) (b bool) {
	// Loops through the dispatched_jobs table to determine if any jobs with the provided path have already been dispatched
	l.logger.Critical("Not implemented")
	return
}

func (l *LibraryManagerAdapter) FileEntryByPath(path string) (f controller.File) {
	l.logger.Critical("Not implemented")
	return
}

func (l *LibraryManagerAdapter) SaveFileEntry(f controller.File) {
	l.logger.Critical("Not implemented")
}
