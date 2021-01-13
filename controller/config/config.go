package config

// The config struct could contain a channel to update parameters on the fly.
// An additional method to check and apply any changes would probably be a good idea for this change.

// ControllerConfiguration contains all of the config options relating to the Controller. Use UpdateChan to update values on the fly.
type ControllerConfiguration struct {
	UpdateChan              *chan string
	SearchDir               string
	FileSystemCheckInterval int
	HealthCheckInterval     int
}

// ProcessUpdates reads the update channel and applies any configuration updates that have been sent through
func (c *ControllerConfiguration) ProcessUpdates() {
}
