// Package options is a centralized location for all supported command-line/environment variable options for the Encodarr Controller
package options

import (
	"fmt"
	"log"
	"os"
)

type optionConst struct {
	EnvVar      string
	CmdLine     string
	Description string
	Usage       string
}

var portConst optionConst = optionConst{"ENCODARR_PORT", "port", "Sets the port of the HTTP server.", "--port <port>"}
var port string = "8123"

var configDirConst optionConst = optionConst{"ENCODARR_CONFIG_DIR", "config-dir", "Sets the location that configuration files are saved to.", "--config-dir <directory>"}
var configDir string = ""

var inputsParsed bool = false

func init() {
	cDir, err := os.UserConfigDir()
	if err != nil {
		log.Fatalln(err)
	}
	configDir = cDir + "/encodarr/controller/config"
}

// parseInputs parses the command line and environment variables into Golang variables
func parseInputs() {
	if inputsParsed {
		return
	}

	// HTTP Server port
	stringVarFromEnv(&port, portConst.EnvVar)
	stringVar(&port, portConst.CmdLine, portConst.Description, portConst.Usage)

	// Config directory
	stringVarFromEnv(&configDir, configDirConst.EnvVar)
	stringVar(&configDir, configDirConst.CmdLine, configDirConst.Description, configDirConst.Usage)

	makeConfigDir()

	parseCL()

	inputsParsed = true
}

// stringVarFromEnv applies the string value found from environment variables to the passed variable
// but only if the returned value is not an empty string
func stringVarFromEnv(s *string, key string) {
	v := os.Getenv(key)
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
	err := os.MkdirAll(configDir, 0777)
	if err != nil {
		log.Fatalln(fmt.Sprintf("Failed to create config directory '%v' because of error: %v", configDir, err.Error()))
	}
}
