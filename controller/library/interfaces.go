package library

import (
	"io/fs"

	"github.com/BrenekH/encodarr/controller"
)

type MetadataReader interface {
	Read(path string) (controller.FileMetadata, error)
}

type CommandDecider interface {
	Decide(m controller.FileMetadata, cmdDeciderSettings string) (cmd []string, err error)
	DefaultSettings() string
}

// stater is an interface that allows for the mocking of os.Stat for testing.
type stater interface {
	Stat(name string) (fs.FileInfo, error)
}

// videoFileser is an interface that allows for the mocking of GetVideoFilesFromDir for testing.
type videoFileser interface {
	VideoFiles(dir string) ([]string, error)
}

type fileRemover interface {
	Remove(path string) error
}

type fileMover interface {
	Move(from string, to string) error
}

type fileStater interface {
	Stat(path string) (fs.FileInfo, error)
}
