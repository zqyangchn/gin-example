package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	"gin-example/pkg/app"
	"gin-example/pkg/convert"
	"gin-example/pkg/errcode"
	"gin-example/service/tag"
)

// curl -X GET "http://127.0.0.1:8000/api/v1/tags?pageNumber=1&pageSize=10&name=zqyangchn&state=1"
// 接口校验
type GetTagForm struct {
	Name       string `form:"name" binding:"max=100"`
	State      int    `form:"state,default=1" binding:"oneof=0 1"`
	PageNumber int    `form:"pageNumber,default=1" binding:"min=1"`
	PageSize   int    `form:"pageSize,default=1" binding:"min=1,max=10"`
}

// @Summary 获取多个标签
// @Produce json
// @Param name query string false "标签名称" maxlength(100)
// @Param state query int false "状态" Enums(0,1) default(1)
// @Param pageNumber query int false "页码"
// @Param PageSize query int false "每页数量"
// @Success 200 {object} tagsvc.TagListResponse
// @Failure 500 {object} app.Response
// @Router /api/v1/tags [get]
func GetTags(c *gin.Context) {
	appG := app.Gin{Context: c}

	query := GetTagForm{}
	if err := app.BindAndValid(c, &query); err != nil {
		appG.Response(http.StatusBadRequest, errcode.InvalidParamsError.WithDetails(err.Error()), struct{}{})
		return
	}

	tagService := tagsvc.Tag{
		Name:       query.Name,
		State:      query.State,
		PageNumber: query.PageNumber,
		PageSize:   query.PageSize,
	}

	tagList, err := tagService.GetTags()
	if err != nil {
		appG.Response(http.StatusInternalServerError, errcode.GetTagError.WithDetails(err.Error()), nil)
		return
	}
	appG.Response(http.StatusOK, errcode.Success, tagList)
}

/*
	curl -X POST "http://127.0.0.1:8000/api/v1/tags" -H "accept: application/json" -H "Content-Type: application/json" -d '
	{
	"name": "zqyangchn",
	"createdBy": "zqyangchn",
	"state": 1
	}'
*/
// go get -u github.com/go-playground/validator/v10
type AddTagForm struct {
	Name      string `form:"name" binding:"required,min=3,max=100"`
	CreatedBy string `form:"createdBy" binding:"required,min=3,max=100"`
	State     int    `form:"state,default=1" binding:"oneof=0 1"`
}

// @Summary 添加标签
// @Produce json
// @Param name body string true "Name" minlength(3) maxlength(100)
// @Param createdBy body string true "CreatedBy" minlength(3) maxlength(100)
// @Param state body int false "State" Enums(0,1) default(1)
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /api/v1/tags [post]
func AddTag(c *gin.Context) {
	var (
		appG = app.Gin{Context: c}
		form AddTagForm
	)

	if err := app.BindAndValid(c, &form); err != nil {
		appG.Response(http.StatusBadRequest, errcode.InvalidParamsError.WithDetails(err.Error()), struct{}{})
		return
	}

	tagService := tagsvc.Tag{
		Name:      form.Name,
		CreatedBy: form.CreatedBy,
		State:     form.State,
	}
	exists, err := tagService.ExistByName()
	if err != nil {
		appG.Response(http.StatusInternalServerError, errcode.ServerError.WithDetails(err.Error()), struct{}{})
		return
	}
	if exists {
		appG.Response(http.StatusOK, errcode.CreateTagError.WithDetails("Tag Name Exist"), struct{}{})
		return
	}

	err = tagService.Add()
	if err != nil {
		appG.Response(http.StatusInternalServerError, errcode.ServerError.WithDetails(err.Error()), struct{}{})
		return
	}

	appG.Response(http.StatusOK, errcode.Success, struct{}{})
}

/*
	curl -X PUT "http://127.0.0.1:8000/api/v1/tags/1" -H "accept: application/json" -H "Content-Type: application/json" -d '
	{
		"name": "zqyangchn",
		"modifiedBy": "zqyang",
		"state": 1
	}'
*/
type EditTagForm struct {
	ID         int    `form:"id" binding:"required,min=1"`
	Name       string `form:"name" binding:"required,max=100"`
	ModifiedBy string `form:"modifiedBy" binding:"required,max=100"`
	State      int    `form:"state,default=1" binding:"oneof=0 1"`
}

// @Summary 更新标签
// @Produce json
// @Param id path int true "标签id"
// @Param name body string true "标签名称" minlength(3) maxlength(100)
// @Param state body int false "状态" Enums(0,1) default(0)
// @Param modified_by body string true "修改者" minlength(3) maxlength(100)
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /api/v1/tags/{id} [put]
func EditTag(c *gin.Context) {
	appG := app.Gin{Context: c}
	form := EditTagForm{ID: convert.StrTo(c.Param("id")).MustInt()}

	if err := app.BindAndValid(c, &form); err != nil {
		appG.Response(http.StatusBadRequest, errcode.InvalidParamsError.WithDetails(err.Error()), struct{}{})
		return
	}

	tagService := tagsvc.Tag{
		ID:         form.ID,
		Name:       form.Name,
		ModifiedBy: form.ModifiedBy,
		State:      form.State,
	}

	exists, err := tagService.ExistByID()
	if err != nil {
		appG.Response(http.StatusInternalServerError, errcode.ServerError.WithDetails(err.Error()), struct{}{})
		return
	}
	if !exists {
		appG.Response(http.StatusOK, errcode.EditTagError.WithDetails("Tag id not exist"), struct{}{})
		return
	}

	if err := tagService.Edit(); err != nil {
		appG.Response(http.StatusInternalServerError, errcode.EditTagError.WithDetails(err.Error()), struct{}{})
		return
	}
	appG.Response(http.StatusOK, errcode.Success, struct{}{})
}

// curl -X DELETE "http://127.0.0.1:8000/api/v1/tags/2"

// @Summary 删除标签
// @Produce json
// @Param id path int true "标签id"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /api/v1/tags/{id} [delete]
func DeleteTag(c *gin.Context) {
	appG := app.Gin{Context: c}
	id := convert.StrTo(c.Param("id")).MustInt()
	if err := validator.New().Var(id, "gte=1"); err != nil {
		appG.Response(http.StatusBadRequest, errcode.DeleteTagError.WithDetails(err.Error()), struct{}{})
		return
	}

	tagService := tagsvc.Tag{ID: id}
	exists, err := tagService.ExistByID()
	if err != nil {
		appG.Response(http.StatusInternalServerError, errcode.ServerError.WithDetails(err.Error()), struct{}{})
		return
	}
	if !exists {
		appG.Response(http.StatusOK, errcode.DeleteTagError.WithDetails("Tag id not exist"), struct{}{})
		return
	}

	if err := tagService.Delete(); err != nil {
		appG.Response(http.StatusInternalServerError, errcode.DeleteTagError.WithDetails(err.Error()), struct{}{})
		return
	}
	appG.Response(http.StatusOK, errcode.Success, struct{}{})
}
