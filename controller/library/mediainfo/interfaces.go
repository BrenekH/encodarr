package mediainfo

// Commander is an interface that allows for mocking out the os/exec package for testing.
type Commander interface {
	Command(name string, args ...string) Cmder
}

// Cmder is an interface for mocking out the exec.Cmd struct.
type Cmder interface {
	Output() ([]byte, error)
}
