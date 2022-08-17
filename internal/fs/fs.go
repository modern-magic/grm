package fs

type FS interface {
	ReadFile(path string) (fileContent string, canonicalError error, originalError error)
	MkDir(path string, option MakeDirectoryOptions) (canonicalError error, originalError error)
	OuputFile(path string, content []byte) (canonicalError error, originalError error)
	ReadYAML(path string, out interface{}) (content interface{}, canonicalError error, originalError error)
	WriteYAML(path string, content interface{}) (canonicalError error, originalError error)
}
