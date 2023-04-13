package fs

type FS interface {
	ReadFile(path string) (string, error)
	MkDir(path string, option MakeDirectoryOptions) error
	OuputFile(path string, content []byte) error
}
