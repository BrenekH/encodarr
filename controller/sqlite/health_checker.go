package sqlite

import "github.com/BrenekH/encodarr/controller"

func NewHealthCheckerAdapater(db *SQLiteDatabase) HealthCheckerAdapter {
	return HealthCheckerAdapter{db: db}
}

// HealthCheckerAdapter satisfies the controller.HealthCheckerDataStorer interface by turning interface
// requests into SQL requests that are passed on to an underlying SQLiteDatabase.
type HealthCheckerAdapter struct {
	db *SQLiteDatabase
}

func (h *HealthCheckerAdapter) DispatchedJobs() (d []controller.DispatchedJob) {
	// TODO: Implement
	return
}

func (h *HealthCheckerAdapter) DeleteJob(uuid controller.UUID) {
	// TODO: Implement
}
