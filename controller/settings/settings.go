package settings

import (
	"bytes"
	"encoding/json"
	"io"
	"os"
	"time"

	"github.com/BrenekH/encodarr/controller"
)

type readWriteSeekCloser interface {
	io.ReadWriteCloser
	io.Seeker
	Truncate(size int64) error
}

// Store satisfies the controller.SettingsStorer interface using a JSON file.
type Store struct {
	healthCheckInterval uint64
	healthCheckTimeout  uint64
	logVerbosity        string

	file   readWriteSeekCloser
	closed bool
}

// settings a marshaling struct used for converting between a slice of bytes and the parsed values
// in SettingsStore.
type settings struct {
	HealthCheckInterval uint64
	HealthCheckTimeout  uint64
	LogVerbosity        string
}

// Load loads the settings from the file.
func (s *Store) Load() error {
	if s.closed {
		return controller.ErrClosed
	}

	if _, err := s.file.Seek(0, io.SeekStart); err != nil {
		return err
	}

	b, err := io.ReadAll(s.file)
	if err != nil {
		return err
	}

	se := settings{}
	err = json.Unmarshal(b, &se)
	if err != nil {
		return err
	}

	s.healthCheckInterval = se.HealthCheckInterval
	s.healthCheckTimeout = se.HealthCheckTimeout
	s.logVerbosity = se.LogVerbosity

	return nil
}

// Save saves the settings to the file.
func (s *Store) Save() error {
	if s.closed {
		return controller.ErrClosed
	}

	// Erase current contents
	if err := s.file.Truncate(0); err != nil {
		return err
	}

	// Move file pointer to start
	if _, err := s.file.Seek(0, io.SeekStart); err != nil {
		return err
	}

	se := settings{
		HealthCheckInterval: s.healthCheckInterval,
		HealthCheckTimeout:  s.healthCheckTimeout,
		LogVerbosity:        s.logVerbosity,
	}
	b, err := json.MarshalIndent(se, "", "\t")
	if err != nil {
		return err
	}

	if _, err := io.Copy(s.file, bytes.NewReader(b)); err != nil {
		return err
	}

	return nil
}

// Close closes the underlying file
func (s *Store) Close() error {
	s.closed = true
	return s.file.Close()
}

// SettingsStore Getters and Setters

// HealthCheckInterval returns the currently set health check interval.
func (s *Store) HealthCheckInterval() uint64 {
	return s.healthCheckInterval
}

// SetHealthCheckInterval sets the health check interval to the provided value.
func (s *Store) SetHealthCheckInterval(n uint64) {
	s.healthCheckInterval = n
}

// HealthCheckTimeout returns the currently set health check timeout value.
func (s *Store) HealthCheckTimeout() uint64 {
	return s.healthCheckTimeout
}

// SetHealthCheckTimeout sets the health check timeout to the provided value.
func (s *Store) SetHealthCheckTimeout(n uint64) {
	s.healthCheckTimeout = n
}

// LogVerbosity returns the currently set log verbosity value.
func (s *Store) LogVerbosity() string {
	return s.logVerbosity
}

// SetLogVerbosity sets the log verbosity to the provided value.
func (s *Store) SetLogVerbosity(n string) {
	s.logVerbosity = n
}

// NewStore returns an instantiated SettingsStore.
func NewStore(configDir string) (Store, error) {
	// Setup a SettingsStore struct with sensible defaults
	s := defaultSettings()

	var err error
	s.file, err = os.OpenFile(configDir+"/settings.json", os.O_RDWR|os.O_CREATE, 0777)
	if err != nil {
		return s, err
	}

	// Save to the file if the file is empty
	var b []byte
	b, err = io.ReadAll(s.file)
	if err != nil {
		return s, err
	}
	if len(b) == 0 {
		err = s.Save()
		if err != nil {
			return s, err
		}
	}

	err = s.Load()
	return s, err
}

// defaultSettings returns a SettingsStore struct with sensible defaults applied.
func defaultSettings() Store {
	return Store{
		healthCheckInterval: uint64(1 * time.Minute),
		healthCheckTimeout:  uint64(1 * time.Hour),
		logVerbosity:        "INFO",
	}
}
