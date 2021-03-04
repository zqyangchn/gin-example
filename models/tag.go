package models

import (
	"gorm.io/gorm"

	"gin-example/pkg/app"
	"gin-example/pkg/database"
)

type Tag struct {
	gorm.Model

	Name       string `json:"name"`
	CreatedBy  string `json:"created_by"`
	ModifiedBy string `json:"modified_by"`
	State      int    `json:"state"`
}

func GetTagTotal(maps map[string]interface{}) (uint, error) {
	var count int64

	if err := database.GetGormDB().Model(&Tag{}).Where(maps).Count(&count).Error; err != nil {
		return 0, err
	}

	return uint(count), nil
}

// ExistTagByName checks if there is a tag with the same name
func ExistTagByName(name string) (bool, error) {
	var tag Tag
	err := database.GetGormDB().Select("id").Where("name = ?", name).First(&tag).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return false, err
	}

	if tag.ID > 0 {
		return true, nil
	}

	return false, nil
}

// ExistTagByID determines whether a Tag exists based on the ID
func ExistTagByID(id int) (bool, error) {
	var tag Tag
	err := database.GetGormDB().Select("id").Where("id = ?", id).First(&tag).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return false, err
	}
	if tag.ID > 0 {
		return true, nil
	}
	return false, nil
}

// AddTag Add a Tag
func AddTag(name string, state int, createdBy string) error {
	tag := Tag{
		Name:      name,
		State:     state,
		CreatedBy: createdBy,
	}
	db := database.GetGormDB()
	if err := db.Create(&tag).Error; err != nil {
		return err
	}

	return nil
}

// EditTag modify a single tag
func EditTag(id int, data interface{}) error {
	if err := database.GetGormDB().Model(&Tag{}).Where("id = ?", id).Updates(data).Error; err != nil {
		return err
	}

	return nil
}

// DeleteTag delete a tag
func DeleteTag(id int) error {
	if err := database.GetGormDB().Where("id = ?", id).Delete(&Tag{}).Error; err != nil {
		return err
	}

	/*
		真实删除
		tag := Tag{
			Model: Model{
				ID: id,
			},
		}
		if err := db.Unscoped().Delete(&tag).Error; err != nil {
			return err
		}
	*/

	return nil
}

func GetTags(pageNumber, pageSize int, maps interface{}) ([]Tag, error) {
	var tags []Tag
	pageOffset := app.GetPageOffset(pageNumber, pageSize)

	if err := database.GetGormDB().Offset(pageOffset).Limit(pageSize).Where(maps).Find(&tags).Error; err != nil {
		return nil, err
	}

	return tags, nil
}
