package fileservice

import "github.com/dmitrovia/passkeeper/internal/server/storage"

type FileService struct {
	repository storage.FileStorage
}

func NewFileService(
	rep storage.UserStorage,
) *FileService {
	return &FileService{
		repository: rep,
	}
}
