package fs

import (
	"fmt"
	"os"
	"path"
	"runtime"

	"gopkg.in/yaml.v3"
)

type fsImpl struct {
}

func IsWindows() bool {
	return runtime.GOOS == "windows"
}

func NewFS() FS {
	return &fsImpl{}
}

type MakeDirectoryOptions struct {
	recursive bool
	mode      os.FileMode
}

func (fs *fsImpl) MkDir(file string, option MakeDirectoryOptions) error {

	if option.mode == 0 {
		option.mode = os.ModePerm.Perm()
	}
	if option.recursive {
		err := os.MkdirAll(file, option.mode)
		return err
	}
	err := os.Mkdir(file, option.mode)
	return err
}

func (fs *fsImpl) ReadFile(file string) (string, error) {
	buffer, err := os.ReadFile(file)
	fileContent := string(buffer)
	return fileContent, err
}

func (fs *fsImpl) OuputFile(file string, content []byte) error {

	dirName := path.Dir(file)
	_, err := os.Stat(dirName)
	if os.IsNotExist(err) {
		err := fs.MkDir(dirName, MakeDirectoryOptions{
			recursive: true,
		})
		if err != nil {
			return err
		}
	}
	err = os.WriteFile(file, content, 0644)
	return err
}

func (fs *fsImpl) ReadYAML(file string, out interface{}) (interface{}, error) {
	ext := path.Ext(file)
	if ext != ".yaml" && ext != ".yml" {
		err := fmt.Errorf("%s%s%s\n", "Error with", file, "Please pass a yaml file.")
		return err, err
	}
	_, err := os.Stat(file)
	if err != nil {
		if !os.IsExist(err) {
			return err, err
		}
	}
	fileContent, err := fs.ReadFile(file)
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal([]byte(fileContent), out)
	return out, err
}

func (fs *fsImpl) WriteYAML(file string, content interface{}) error {
	out, err := yaml.Marshal(content)
	if err != nil {
		return err
	}
	err = fs.OuputFile(file, out)
	return err
}
