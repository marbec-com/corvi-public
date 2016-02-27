package models

type Category struct {
	ID   uint
	Name string
}

func NewCategory() *Category {
	return &Category{}
}
