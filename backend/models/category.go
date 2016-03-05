package models

import (
	"time"
)

type Category struct {
	ID        uint
	Name      string
	CreatedAt time.Time
}

func NewCategory() *Category {
	return &Category{
		CreatedAt: time.Now(),
	}
}
