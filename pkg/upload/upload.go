package upload

import (
	"io"
	"mime/multipart"
	"os"
	"path"
	"strings"

	"gin-example/pkg/app"
	"gin-example/pkg/setting"
)

type FileType int

const (
	// 图片类型
	TypeImage FileType = iota + 1
)

func GetStoragePath() string {
	return setting.AppSetting.UploadSavePath
}

func EncryptionFileName(name string) string {
	extension := path.Ext(name)

	md5Name := app.EncodeMD5(strings.TrimSuffix(name, extension))
	return md5Name + extension
}

func CheckExist(dst string) bool {
	_, err := os.Stat(dst)
	return os.IsNotExist(err)
}

func CheckContainExt(t FileType, name string) bool {
	ext := strings.ToUpper(path.Ext(name))
	switch t {
	case TypeImage:
		for _, allowExt := range setting.AppSetting.UploadImageAllowExts {
			if strings.ToUpper(allowExt) == ext {
				return true
			}
		}
	}
	return false
}

func CheckMaxSize(t FileType, f multipart.File) bool {
	size := 0
	buffer := make([]byte, 1024)

	for true {
		n, err := io.ReadFull(f, buffer)
		if err != nil {
			size += n
			break
		}
		size += n
	}

	switch t {
	case TypeImage:
		if size <= setting.AppSetting.UploadImageMaxSize*1024*1024 {
			return true
		}
	}
	return false
}

func CheckPermission(dst string) bool {
	_, err := os.Stat(dst)
	return os.IsPermission(err)
}

func CreateStoragePath(storagePath string, permission os.FileMode) error {
	if err := os.MkdirAll(storagePath, permission); err != nil {
		return err
	}
	return nil
}

func StorageFile(multipartFileHeader *multipart.FileHeader, storageDestination string) error {
	src, err := multipartFileHeader.Open()
	if err != nil {
		return err
	}
	defer func() {
		_ = src.Close()
	}()

	dst, err := os.Create(storageDestination)
	if err != nil {
		return err
	}
	defer func() {
		_ = dst.Close()
	}()

	_, err = io.Copy(dst, src)
	return err
}
