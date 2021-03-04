package models

import "gin-example/pkg/database"

func Setup() error {
	if err := autoMigrateAll(); err != nil {
		return err
	}
	return nil
}

func autoMigrateAll() error {
	toMigrate := []interface{}{
		&User{},
		&Tag{},
	}
	if err := database.GetGormDB().AutoMigrate(toMigrate...); err != nil {
		return err
	}
	return nil
}
