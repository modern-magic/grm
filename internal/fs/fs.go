package fs

type FS interface {
	MkDir(path string, option MakeDirectoryOptions) error
	OuputFile(path string, content []byte) error
	Rm(path string) error
}
