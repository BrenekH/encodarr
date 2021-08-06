// Package globals is the location of read-only constants such as Version, which is set at build time for release binaries.
package globals

// Version is a read-only constant that specifies the software version.
// Using ldflags, Version can be set at build time. If it is not set using ldflags, its value will be 'develop'.
var Version string = "develop"
