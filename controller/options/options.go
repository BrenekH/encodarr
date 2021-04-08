// Package options is a centralized location for all supported command-line/environment variable options for Encodarr
package options

import (
	"fmt"
	"log"
	"os"

	"github.com/BrenekH/logange"
)

type optionConst struct {
	EnvVar  string
	CmdLine string
}

var portConst optionConst = optionConst{"ENCODARR_PORT", "port"}
var port string = "8123"

var configDirConst optionConst = optionConst{"ENCODARR_CONFIG_DIR", "config-dir"}
var configDir string = ""

var inputsParsed bool = false

var logger logange.Logger

func init() {
	cDir, err := os.UserConfigDir()
	if err != nil {
		log.Fatalln(err)
	}
	configDir = cDir + "/encodarr/controller/config"

	logger = logange.NewLogger("options")
}

// parseInputs parses the command line and environment variables into Golang variables
func parseInputs() {
	if inputsParsed {
		return
	}

	// HTTP Server port
	stringVarFromEnv(&port, portConst.EnvVar)
	stringVar(&port, portConst.CmdLine, "")

	// Config directory
	stringVarFromEnv(&configDir, configDirConst.EnvVar)
	stringVar(&configDir, configDirConst.CmdLine, "")

	makeConfigDir()

	parseCL()

	inputsParsed = true
}

// stringVarFromEnv applies the string value found from environment variables to the passed variable
// but only if the returned value is not an empty string
func stringVarFromEnv(s *string, key string) {
	v := os.Getenv(key)
	logger.Debug(fmt.Sprintf("Got `%v` from `%v`", v, key))
	if v != "" {
		*s = v
	}
}

// Port returns the parsed HTTP server port
func Port() string {
	parseInputs()
	return port
}

// ConfigDir returns the passed config directory
func ConfigDir() string {
	parseInputs()
	return configDir
}

// makeConfigDir creates the options.configDir
func makeConfigDir() {
	err := os.MkdirAll(configDir, 0644)
	if err != nil {
		fmt.Println(err)
		logger.Critical(fmt.Sprintf("Failed to create config directory '%v' because of error: %v", configDir, err.Error()))
	}
}
