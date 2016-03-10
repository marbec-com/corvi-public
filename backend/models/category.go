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

func (c *Category) Equal(a *Category) bool {
	return c.ID == a.ID && c.Name == a.Name && c.CreatedAt.Equal(a.CreatedAt)
}
