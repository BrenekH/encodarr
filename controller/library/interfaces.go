package library

import "github.com/BrenekH/encodarr/controller"

type MetadataReader interface {
	Read(path string) (controller.FileMetadata, error)
}

type CommandDecider interface {
	Decide(m controller.FileMetadata, cmdDeciderSettings string) (runCmd bool, cmd []string)
}
