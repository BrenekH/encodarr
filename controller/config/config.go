package config

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/BrenekH/encodarr/controller/options"
)

// Settings is used to represent how settings are saved to a file
type Settings struct {
	HealthCheckInterval uint64
	HealthCheckTimeout  uint64
	LogVerbosity        string
	SmallerFiles        bool
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

	err = os.WriteFile(fmt.Sprintf("%v/settings.json", options.ConfigDir()), b, 0666)
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
		HealthCheckInterval: uint64(1 * time.Minute),
		HealthCheckTimeout:  uint64(1 * time.Hour),
		LogVerbosity:        "INFO",
		SmallerFiles:        false,
	}
}
