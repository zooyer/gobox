package types

type OS struct {
	Setenv  func(key string, value string) error
	Environ func() []string
	Getpid  func() int
}

type Syscall struct {
	OS OS
}
