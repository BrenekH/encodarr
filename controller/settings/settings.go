package settings

import (
	"bytes"
	"encoding/json"
	"io"
	"os"
	"time"

	"github.com/BrenekH/encodarr/controller"
)

type ReadWriteSeekCloser interface {
	io.ReadWriteCloser
	io.Seeker
	Truncate(size int64) error
}

type SettingsStore struct {
	healthCheckInterval uint64
	healthCheckTimeout  uint64
	logVerbosity        string

	file   ReadWriteSeekCloser
	closed bool
}

// settings a marshaling struct used for converting between a slice of bytes and the parsed values
// in SettingsStore.
type settings struct {
	HealthCheckInterval uint64
	HealthCheckTimeout  uint64
	LogVerbosity        string
}

func (s *SettingsStore) Load() error {
	if s.closed {
		return controller.ErrClosed
	}

	s.file.Seek(0, io.SeekStart)
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

func (s *SettingsStore) Save() error {
	if s.closed {
		return controller.ErrClosed
	}

	// Erase current contents
	s.file.Truncate(0)

	// Move file pointer to start
	s.file.Seek(0, io.SeekStart)

	se := settings{
		HealthCheckInterval: s.healthCheckInterval,
		HealthCheckTimeout:  s.healthCheckTimeout,
		LogVerbosity:        s.logVerbosity,
	}
	b, err := json.MarshalIndent(se, "", "\t")
	if err != nil {
		return err
	}

	io.Copy(s.file, bytes.NewReader(b))

	return nil
}

func (s *SettingsStore) Close() error {
	s.closed = true
	return s.file.Close()
}

// SettingsStore Getters and Setters

func (s *SettingsStore) HealthCheckInterval() uint64 {
	return s.healthCheckInterval
}

func (s *SettingsStore) SetHealthCheckInterval(n uint64) {
	s.healthCheckInterval = n
}

func (s *SettingsStore) HealthCheckTimeout() uint64 {
	return s.healthCheckTimeout
}

func (s *SettingsStore) SetHealthCheckTimeout(n uint64) {
	s.healthCheckTimeout = n
}

func (s *SettingsStore) LogVerbosity() string {
	return s.logVerbosity
}

func (s *SettingsStore) SetLogVerbosity(n string) {
	s.logVerbosity = n
}

func NewSettingsStore(configDir string) (SettingsStore, error) {
	// Setup a SettingsStore struct with sensible defaults
	s := defaultSettings()

	var err error
	s.file, err = os.OpenFile(configDir+"/settings.json", os.O_RDWR|os.O_CREATE, 0755)
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
func defaultSettings() SettingsStore {
	return SettingsStore{
		healthCheckInterval: uint64(1 * time.Minute),
		healthCheckTimeout:  uint64(1 * time.Hour),
		logVerbosity:        "INFO",
	}
}
