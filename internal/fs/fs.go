package fs

type FS interface {
	ReadFile(path string) (string, error)
	MkDir(path string, option MakeDirectoryOptions) error
	OuputFile(path string, content []byte) error
	ReadYAML(path string, out interface{}) (content interface{}, err error)
	WriteYAML(path string, content interface{}) error
}
