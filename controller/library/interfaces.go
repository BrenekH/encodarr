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

// stater is an interface that allows for the mocking of os.Stat for testing.
type stater interface {
	Stat(name string) (fs.FileInfo, error)
}

// videoFileser is an interface that allows for the mocking of GetVideoFilesFromDir for testing.
type videoFileser interface {
	VideoFiles(dir string) ([]string, error)
}
