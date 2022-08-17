package fs

import (
	"fmt"
	"os"
	"path"
	"runtime"
	"syscall"

	"gopkg.in/yaml.v3"
)

type FsImpl struct {
}

func IsWindows() bool {
	return runtime.GOOS == "windows"
}

func (fs *FsImpl) canonicalizeError(err error) error {
	if pathError, ok := err.(*os.PathError); ok {
		err = pathError.Unwrap()
	}
	const ERROR_INVALID_NAME syscall.Errno = 123
	if IsWindows() && err == ERROR_INVALID_NAME {
		err = syscall.ENOENT
	}
	if err == syscall.ENOTDIR {
		err = syscall.ENOENT
	}
	return err
}

type MakeDirectoryOptions struct {
	recursive bool
	mode      os.FileMode
}

func (fs *FsImpl) MkDir(file string, option MakeDirectoryOptions) (canonicalError error, originalError error) {

	if option.mode == 0 {
		option.mode = os.ModePerm.Perm()
	}
	if option.recursive {
		originalError := os.MkdirAll(file, option.mode)
		canonicalError = fs.canonicalizeError(originalError)
		return canonicalError, originalError
	}
	originalError = os.Mkdir(file, option.mode)
	canonicalError = fs.canonicalizeError(originalError)
	return canonicalError, originalError
}

func (fs *FsImpl) ReadFile(file string) (filecontent string, canonicalError error, originalError error) {
	buffer, originalError := os.ReadFile(file)
	canonicalError = fs.canonicalizeError(originalError)
	fileContent := string(buffer)
	return fileContent, canonicalError, originalError
}

func (fs *FsImpl) OuputFile(file string, content []byte) (canonicalError error, originalError error) {

	dirName := path.Dir(file)
	_, err := os.Stat(dirName)
	if os.IsNotExist(err) {
		canonicalError, originalError = fs.MkDir(dirName, MakeDirectoryOptions{
			recursive: true,
		})
		if originalError != nil {
			return canonicalError, originalError
		}
	}
	originalError = os.WriteFile(file, content, 0644)
	canonicalError = fs.canonicalizeError(originalError)
	return canonicalError, originalError
}

func (fs *FsImpl) ReadYAML(file string, out interface{}) (content interface{}, canonicalError error, originalError error) {
	ext := path.Ext(file)
	if ext != ".yaml" && ext != ".yml" {
		err := fmt.Errorf("%s%s%s\n", "Error with", file, "Please pass a yaml file.")
		return nil, err, err
	}

	fileContent, canonicalError, originalError := fs.ReadFile(file)
	if canonicalError != nil || originalError != nil {
		return nil, canonicalError, originalError
	}
	originalError = yaml.Unmarshal([]byte(fileContent), out)
	return out, originalError, originalError
}

func (fs *FsImpl) WriteYAML(file string, content interface{}) (canonicalError error, originalError error) {
	out, originalError := yaml.Marshal(content)
	if originalError != nil {
		return originalError, originalError
	}
	originalError = os.WriteFile(file, out, 0644)
	canonicalError = fs.canonicalizeError(originalError)
	return canonicalError, originalError
}
