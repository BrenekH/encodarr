package library

import "github.com/BrenekH/encodarr/controller"

type MetadataReader interface {
	Read(path string) controller.FileMetadata
}

type CommandDecider interface {
	Decide(m controller.FileMetadata, cmdDeciderSettings []byte) (runCmd bool, cmd []string)
}
