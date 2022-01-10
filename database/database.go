package database

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Database struct {
	Db *gorm.DB
}

var models = []interface{}{
	&Post{},
	&Comment{},
}

func New() (*Database, error) {
	db, err := gorm.Open(sqlite.Open("db.sqlite"), &gorm.Config{})

	self := &Database{Db: db}

	if err != nil {
		return self, err
	}

	return self, nil
}

func (db *Database) AutoMigrate() error {
	for _, model := range models {
		err := db.Db.AutoMigrate(model)

		if err != nil {
			return err
		}
	}

	return nil
}

func (db *Database) Truncate() {
	for _, model := range models {
		db.Db.Where("1 = 1").Delete(model)
	}
}
