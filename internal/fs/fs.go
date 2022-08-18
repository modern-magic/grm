package fs

type FS interface {
	ReadFile(path string) (fileContent string, originalError error)
	MkDir(path string, option MakeDirectoryOptions) (originalError error)
	OuputFile(path string, content []byte) (originalError error)
	ReadYAML(path string, out interface{}) (content interface{}, originalError error)
	WriteYAML(path string, content interface{}) (originalError error)
}
