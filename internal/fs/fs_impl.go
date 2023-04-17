package fs

import (
	"os"
	"path"
)

type fsImpl struct {
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

func (fs *fsImpl) Rm(file string) error {
	return os.Remove(file)
}
