// The purpose of this file to provide an API similar to the flag package for parsing command-line arguments
// without impacting the testing package(see https://github.com/golang/go/issues/31859 and https://github.com/golang/go/issues/39093).

package options

import (
	"fmt"
	"os"
	"strings"
)

// flagger defines a type agnostic interface to parse out flags.
type flagger interface {
	Name() string
	Usage() string
	Parse(string) error
}

var flags []flagger

// stringVar replaces flag.StringVar, but without the default value.
// That functionality is provided by the rest of the options package.
func stringVar(p *string, name, usage string) {
	// TODO: Test
	flags = append(flags, StringFlag{
		name:    name,
		usage:   usage,
		pointer: p,
	})
}

// parseCL parses the command-line arguments into the registered options.
// Replaces flag.Parse.
func parseCL() {
	var args []string = os.Args[1:]

	for k, v := range args {
		for _, f := range flags {
			if strings.Replace(v, "--", "", 1) == f.Name() {
				if i := k + 1; i >= len(args) {
					logger.Critical(fmt.Sprintf("Can not parse %v, EOL reached", v))
				} else {
					f.Parse(args[k+1])
				}
			}
		}
	}
}

type StringFlag struct { // TODO: Test struct methods
	name    string
	usage   string
	pointer *string
}

func (f StringFlag) Parse(s string) error {
	*f.pointer = s
	return nil
}

func (f StringFlag) Name() string {
	return f.name
}

func (f StringFlag) Usage() string {
	return f.usage
}
