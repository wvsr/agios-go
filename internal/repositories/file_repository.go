package repositories

import (
	"io"
	"os"
)

type FileRepository interface {
	SaveFile(dstPath string, src io.Reader) error
	CreateDirectory(dirPath string) error
}

type FileRepositoryImpl struct{}

func NewFileRepository() FileRepository {
	return &FileRepositoryImpl{}
}

func (r *FileRepositoryImpl) SaveFile(dstPath string, src io.Reader) error {
	dst, err := os.Create(dstPath)
	if err != nil {
		return err
	}
	defer dst.Close()

	_, err = io.Copy(dst, src)
	return err
}

func (r *FileRepositoryImpl) CreateDirectory(dirPath string) error {
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		return os.MkdirAll(dirPath, 0755)
	}
	return nil
}
