// Package options is a centralized location for all supported command-line/environment variable options for RedCedar
package options

import (
	"flag"
	"fmt"
	"os"

	"github.com/BrenekH/logange"
)

type optionConst struct {
	EnvVar  string
	CmdLine string
}

var portConst optionConst = optionConst{"REDCEDAR_PORT", "port"}
var port string = "8123"

var configDirConst optionConst = optionConst{"REDCEDAR_CONFIG_DIR", "config-dir"}
var configDir string = "/redcedar/config"

var searchDirConst optionConst = optionConst{"REDCEDAR_SEARCH_DIR", "search-dir"}
var searchDir string = ""

var inputsParsed bool = false

var logger logange.Logger

func init() {
	cwd, _ := os.Getwd()
	searchDir = cwd

	logger = logange.NewLogger("options")
}

// parseInputs parses the command line and environment variables into Golang variables
func parseInputs() {
	if inputsParsed {
		return
	}

	// HTTP Server port
	stringVarFromEnv(&port, portConst.EnvVar)
	flag.StringVar(&port, portConst.CmdLine, port, "")

	// Config directory
	stringVarFromEnv(&configDir, configDirConst.EnvVar)
	flag.StringVar(&configDir, configDirConst.CmdLine, configDir, "")

	// Search directory
	stringVarFromEnv(&searchDir, searchDirConst.EnvVar)
	flag.StringVar(&searchDir, searchDirConst.CmdLine, searchDir, "")

	flag.Parse()

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

// SearchDir returns the passed search directory
func SearchDir() string {
	parseInputs()
	return searchDir
}
