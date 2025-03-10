package handlers

import "gorm.io/gorm"

type Handlers struct {
	db *gorm.DB
}

func New(db *gorm.DB) *Handlers {
	return &Handlers{
		db: db,
	}
}
