package dao

import "gorm.io/gorm"

type TestDAO struct {
	db *gorm.DB
}

func NewTestDAO(db *gorm.DB) *TestDAO {
	return &TestDAO{db: db}
}
