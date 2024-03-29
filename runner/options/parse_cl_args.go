// The purpose of this file to provide an API similar to the flag package for parsing command-line arguments
// without impacting the testing package(see https://github.com/golang/go/issues/31859 and https://github.com/golang/go/issues/39093).

package options

import (
	"fmt"
	"os"
	"runtime"
	"strings"
)

// flagger defines a type agnostic interface to parse out flags.
type flagger interface {
	Name() string
	Description() string
	Usage() string
	Parse(string) error
}

var flags []flagger

// stringVar replaces flag.StringVar, but without the default value.
// That functionality is provided by the rest of the options package.
func stringVar(p *string, name, description, usage string) {
	flags = append(flags, stringFlag{
		name:        name,
		description: description,
		usage:       usage,
		pointer:     p,
	})
}

// parseCL parses the command-line arguments into the registered options.
// Replaces flag.Parse.
func parseCL() {
	var args []string = os.Args[1:]

	for k, v := range args {
		if v == "--help" {
			helpStr := fmt.Sprintf("Encodarr Runner %v Help\n\n", Version)

			for _, f := range flags {
				helpStr += fmt.Sprintf(" --%v - %v\n   Usage: \"%v\"\n\n",
					f.Name(),
					f.Description(),
					f.Usage(),
				)
			}

			fmt.Println(strings.TrimRight(helpStr, "\n"))
			os.Exit(0)
		} else if v == "--version" {
			fmt.Printf("Encodarr Runner %v %v/%v", Version, runtime.GOOS, runtime.GOARCH)
			os.Exit(0)
		}

		for _, f := range flags {
			if strings.Replace(v, "--", "", 1) == f.Name() {
				if i := k + 1; i >= len(args) {
					fmt.Printf("Can not parse %v, EOL reached", v)
					logger.Critical(fmt.Sprintf("Can not parse %v, EOL reached", v))
				} else {
					f.Parse(args[k+1])
				}
			}
		}
	}
}

type stringFlag struct {
	name        string
	description string
	usage       string
	pointer     *string
}

func (f stringFlag) Parse(s string) error {
	*f.pointer = s
	return nil
}

func (f stringFlag) Description() string {
	return f.description
}

func (f stringFlag) Name() string {
	return f.name
}

func (f stringFlag) Usage() string {
	return f.usage
}
