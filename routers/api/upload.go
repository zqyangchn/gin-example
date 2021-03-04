package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"gin-example/pkg/app"
	"gin-example/pkg/errcode"
	"gin-example/pkg/logging"
	"gin-example/pkg/upload"
	"gin-example/service"
)

// curl -X POST http://127.0.0.1:8000/upload/file -F file=@mianyin.jpg -F type=1
type UploadFileForm struct {
	Type int `form:"type,default=1" binding:"oneof=1"`
}

// @Summary 上传文件
// @Produce multipart/form-data
// @Param type query string false "文件类型" Enums(1) default(1)
// @Param file query string false "文件"
// @Success 200 {object} service.UploadFileResponse
// @Failure 500 {object} app.Response
// @Router /upload/file [post]
func UploadFile(c *gin.Context) {
	appG := app.Gin{Context: c}
	uploadFileForm := UploadFileForm{}

	if err := app.BindAndValid(c, &uploadFileForm); err != nil {
		appG.Response(http.StatusBadRequest, errcode.UploadFileError.WithDetails(err.Error()), struct{}{})
		return
	}

	file, fileHeader, err := c.Request.FormFile("file")
	if err != nil {
		appG.Response(http.StatusBadRequest, errcode.UploadFileError.WithDetails(err.Error()), struct{}{})
		return
	}
	if fileHeader == nil {
		appG.Response(http.StatusBadRequest, errcode.UploadFileError.WithDetails("Header is nil"), struct{}{})
		return
	}

	fileInformation, err := service.UploadFile(upload.FileType(uploadFileForm.Type), file, fileHeader)
	if err != nil {
		logging.Logger.Error("UploadFile error", zap.Error(err))
		appG.Response(http.StatusInternalServerError, errcode.UploadFileError.WithDetails(err.Error()), struct{}{})
	}

	appG.ResponseSuccess(http.StatusOK,
		&service.UploadFileResponse{
			ErrorMessage: errcode.Success,
			AccessUrl:    fileInformation.AccessUrl,
		})
}
