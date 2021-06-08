package sqlite

import "github.com/BrenekH/encodarr/controller"

func NewUserInterfacerAdapter(db *SQLiteDatabase) UserInterfacerAdapter {
	return UserInterfacerAdapter{db: db}
}

type UserInterfacerAdapter struct {
	db *SQLiteDatabase
}

func (u *UserInterfacerAdapter) DispatchedJobs() ([]controller.DispatchedJob, error) {
	// TODO: Implement
	return []controller.DispatchedJob{}, nil
}

func (u *UserInterfacerAdapter) HistoryEntries() ([]controller.History, error) {
	// TODO: Implement
	return []controller.History{}, nil
}

func (u *UserInterfacerAdapter) DeleteLibrary(id int) error {
	// TODO: Implement
	return nil
}
