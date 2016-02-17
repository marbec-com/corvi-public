package models

type Category struct {
	ID   uint
	Name string
}

func (c *Category) Boxes() []*Box {

	boxes := make([]*Box, 0)

	// SELECT * FROM Boxes Where CategoryID = c.ID

	return boxes
}
