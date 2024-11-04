package dtos

type Inventory struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

func NewInventory(id int64, name string) *Inventory {
	return &Inventory{
		ID:   id,
		Name: name,
	}
}
