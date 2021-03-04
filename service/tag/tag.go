package tagsvc

import (
	"gin-example/models"
	"gin-example/pkg/errcode"
)

type Tag struct {
	ID         int
	Name       string
	CreatedBy  string
	ModifiedBy string
	State      int

	PageNumber int
	PageSize   int
}

type TagList struct {
	Tags       []models.Tag
	TotalCount uint
}

// for swagger show Response
type TagListResponse struct {
	errcode.ErrorMessage
	Data TagList
}

func (t *Tag) getMaps() map[string]interface{} {
	maps := make(map[string]interface{})

	if t.Name != "" {
		maps["name"] = t.Name
	}
	if t.State >= 0 {
		maps["state"] = t.State
	}

	return maps
}

func (t *Tag) ExistByName() (bool, error) {
	return models.ExistTagByName(t.Name)
}

func (t *Tag) ExistByID() (bool, error) {
	return models.ExistTagByID(t.ID)
}

func (t *Tag) Add() error {
	return models.AddTag(t.Name, t.State, t.CreatedBy)
}

func (t *Tag) Edit() error {
	data := make(map[string]interface{})

	data["modified_by"] = t.ModifiedBy
	data["name"] = t.Name
	data["state"] = t.State

	return models.EditTag(t.ID, data)
}

func (t *Tag) Delete() error {
	return models.DeleteTag(t.ID)
}

func (t *Tag) GetTags() (*TagList, error) {
	tags, err := models.GetTags(t.PageNumber, t.PageSize, t.getMaps())
	if err != nil {
		return nil, err
	}
	tagList := &TagList{Tags: tags}

	count, err := models.GetTagTotal(t.getMaps())
	if err != nil {
		return nil, err
	}
	tagList.TotalCount = count

	return tagList, nil
}
