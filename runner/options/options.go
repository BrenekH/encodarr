// options is a centralized location to read all supported command-line/environment variable options for Encodarr Runner
package options

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/BrenekH/logange"
)

// Version represents the current version of the Encodarr Runner
var Version string = "development"

type optionConst struct {
	EnvVar      string
	CmdLine     string
	Description string
	Usage       string
}

var configDirConst optionConst = optionConst{"ENCODARR_CONFIG_DIR", "config-dir", "Sets the location of the directory where configuration files are saved.", "--config-dir <directory>"}
var configDir string = ""

var tempDirConst optionConst = optionConst{"ENCODARR_TEMP_DIR", "temp-dir", "Sets the location where media files are saved when they are being worked on.", "--temp-dir <directory>"}
var tempDir string = os.TempDir()

var logLevelConst optionConst = optionConst{"ENCODARR_LOG_LEVEL", "log-level", "Sets the level of logging to both the console and the logging file.", "--log-level <trace|debug|info|warn|warning|error|critical>"}
var logLevel string = "INFO"

var runnerNameConst optionConst = optionConst{"ENCODARR_RUNNER_NAME", "name", "Sets the name shown in the Web UI.", "--name <runner-name>"}
var runnerName string = ""

var controllerIPConst optionConst = optionConst{"ENCODARR_RUNNER_CONTROLLER_IP", "controller-ip", "Sets the IP for connecting to a Controller.", "--controller-ip <ip>"}
var controllerIP string = "localhost"

var controllerPortConst optionConst = optionConst{"ENCODARR_RUNNER_CONTROLLER_PORT", "controller-port", "Sets the port for connecting to a Controller.", "--controller-port <port>"}
var controllerPort string = "8123"

var inTestMode bool = strings.HasSuffix(os.Args[0], ".test") || strings.HasSuffix(os.Args[0], ".test.exe")

var inputsParsed bool = false

var logger logange.Logger

func init() {
	//! Logange can't be used in this function because it requires the config location

	// Initialize default config location
	cDir, err := os.UserConfigDir()
	if err != nil {
		log.Fatalln(err)
	}
	configDir = cDir + fmt.Sprintf("/encodarr/runner/%v/config", time.Now().Format("2006-01-02-15-04-05.000"))

	// Initialize default Runner name
	hostname, err := os.Hostname()
	if err != nil {
		hostname = "Runner"
	}
	runnerName = fmt.Sprintf("%v-%v", hostname, rand.Intn(999))

	logger = logange.NewLogger("options")
}

// parseInputs parses the command line and environment variables into Golang variables
func parseInputs() {
	if inputsParsed {
		return
	}

	// Config directory
	stringVarFromEnv(&configDir, configDirConst.EnvVar)
	stringVar(&configDir, configDirConst.CmdLine, configDirConst.Description, configDirConst.Usage)

	// Temporary directory
	stringVarFromEnv(&tempDir, tempDirConst.EnvVar)
	stringVar(&tempDir, tempDirConst.CmdLine, tempDirConst.Description, tempDirConst.Usage)

	// Log level
	stringVarFromEnv(&logLevel, logLevelConst.EnvVar)
	stringVar(&logLevel, logLevelConst.CmdLine, logLevelConst.Description, logLevelConst.Usage)

	// Runner name
	stringVarFromEnv(&runnerName, runnerNameConst.EnvVar)
	stringVar(&runnerName, runnerNameConst.CmdLine, runnerNameConst.Description, runnerNameConst.Usage)

	// Controller IP
	stringVarFromEnv(&controllerIP, controllerIPConst.EnvVar)
	stringVar(&controllerIP, controllerIPConst.CmdLine, controllerIPConst.Description, controllerIPConst.Usage)

	// Controller Port
	stringVarFromEnv(&controllerPort, controllerPortConst.EnvVar)
	stringVar(&controllerPort, controllerPortConst.CmdLine, controllerPortConst.Description, controllerPortConst.Usage)

	if !inTestMode {
		makeConfigDir()
	}

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

// ConfigDir returns the passed config directory
func ConfigDir() string {
	parseInputs()
	return configDir
}

func TempDir() string {
	parseInputs()
	return tempDir
}

func LogLevel() logange.Level {
	parseInputs()

	switch strings.ToLower(logLevel) {
	case "trace":
		return logange.LevelTrace
	case "debug":
		return logange.LevelDebug
	case "info":
		return logange.LevelInfo
	case "warn", "warning":
		return logange.LevelWarn
	case "error":
		return logange.LevelError
	case "critical":
		return logange.LevelCritical
	}

	// I'm using fmt.Printf instead of logger.Warn here because
	// I can't be sure that the logger is properly instantiated yet.
	fmt.Printf("Invalid log level: `%v`. Default to INFO.\n", logLevel)
	return logange.LevelInfo
}

func RunnerName() string {
	parseInputs()
	return runnerName
}

func ControllerIP() string {
	parseInputs()
	return controllerIP
}

func ControllerPort() string {
	parseInputs()
	return controllerPort
}

func InTestMode() bool {
	return inTestMode
}

// makeConfigDir creates the options.configDir
func makeConfigDir() {
	err := os.MkdirAll(configDir, 0777)
	if err != nil {
		fmt.Printf("options.makeConfigDir: %v\n", err)
		logger.Critical(fmt.Sprintf("Failed to create config directory '%v' because of error: %v", configDir, err.Error()))
	}
}
