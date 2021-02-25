package config

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/BrenekH/project-redcedar-controller/options"
)

// The config struct could contain a channel to update parameters on the fly.
// An additional method to check and apply any changes would probably be a good idea for this change.

// Settings is used to represent how settings are saved to a file
type Settings struct {
	FileSystemCheckInterval int
	HealthCheckInterval     int
	HealthCheckTimeout      int
	LogVerbosity            string
}

// Global is an instance of ControllerConfiguration that can be accessed anywhere in the program
var Global Settings

// LoadSettings returns a Settings struct that has been instantiated with the values found in options.ConfigDir()/settings.json
func LoadSettings() (Settings, error) {
	b, err := os.ReadFile(fmt.Sprintf("%v/settings.json", options.ConfigDir()))
	if err != nil {
		return Settings{}, err
	}

	settings := Settings{}
	err = json.Unmarshal(b, &settings)
	if err != nil {
		return Settings{}, err
	}

	return settings, nil
}

// SaveSettings saves the provided settings struct to options.ConfigDir()/settings.json
func SaveSettings(s Settings) error {
	b, err := json.MarshalIndent(s, "", "\t")
	if err != nil {
		return err
	}

	err = os.WriteFile(fmt.Sprintf("%v/settings.json", options.ConfigDir()), b, 0664)
	if err != nil {
		return err
	}

	return nil
}

// SaveGlobal saves the global settings variable by passing it to SaveSettings
func SaveGlobal() error {
	return SaveSettings(Global)
}

// DefaultSettings returns a "constant" settings struct with sensible defaults
func DefaultSettings() Settings {
	return Settings{
		FileSystemCheckInterval: int(15 * time.Minute),
		HealthCheckInterval:     int(1 * time.Minute),
		HealthCheckTimeout:      int(1 * time.Hour),
		LogVerbosity:            "INFO",
	}
}