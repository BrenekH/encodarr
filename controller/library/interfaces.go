package library

import (
	"io/fs"

	"github.com/BrenekH/encodarr/controller"
)

type MetadataReader interface {
	Read(path string) (controller.FileMetadata, error)
}

type CommandDecider interface {
	Decide(m controller.FileMetadata, cmdDeciderSettings string) (runCmd bool, cmd []string)
}

type stater interface {
	Stat(name string) (fs.FileInfo, error)
}
