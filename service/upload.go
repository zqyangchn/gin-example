package service

import (
	"errors"
	"mime/multipart"
	"os"
	"strings"

	"gin-example/pkg/errcode"
	"gin-example/pkg/upload"
)

type FileInformation struct {
	Name      string
	AccessUrl string
}

// Response struct
type UploadFileResponse struct {
	*errcode.ErrorMessage
	AccessUrl string
}

// upload file to local path
func UploadFile(fileType upload.FileType, multipartFile multipart.File, multipartFileHeader *multipart.FileHeader) (*FileInformation, error) {
	encryptionName := upload.EncryptionFileName(multipartFileHeader.Filename)
	storagePath := upload.GetStoragePath()
	storageDestination := strings.Join([]string{storagePath, encryptionName}, "/")

	// 检查文件后缀
	if !upload.CheckContainExt(fileType, storageDestination) {
		return nil, errors.New("file suffix is not supported")
	}
	// 检查存储目录是否存在, 不存在, 创建
	if upload.CheckExist(storagePath) {
		err := upload.CreateStoragePath(storagePath, os.ModePerm)
		if err != nil {
			return nil, errors.New("failed to create storage directory")
		}
	}
	// 检查存储目录权限
	if upload.CheckPermission(storagePath) {
		return nil, errors.New("insufficient file permissions")
	}
	// 检查文件大小
	if !upload.CheckMaxSize(fileType, multipartFile) {
		return nil, errors.New("exceeded maximum file limit")
	}
	// 上传文件
	if err := upload.StorageFile(multipartFileHeader, storageDestination); err != nil {
		return nil, err
	}

	return &FileInformation{Name: encryptionName, AccessUrl: strings.Join([]string{"/static", encryptionName}, "/")}, nil
}
