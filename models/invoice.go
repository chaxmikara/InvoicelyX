package models

type Item struct {
	Name     string  `bson:"name" json:"name"`
	Price    float64 `bson:"price" json:"price"`
	Quantity int     `bson:"quantity" json:"quantity"`
}

// Invoice represents the complete invoice structure
type Invoice struct {
	ID        string  `bson:"_id,omitempty" json:"id,omitempty"`
	Customer  string  `bson:"customer" json:"customer"`
	Email     string  `bson:"email" json:"email"`
	Items     []Item  `bson:"items" json:"items"`
	Total     float64 `bson:"total" json:"total"`
	CreatedAt string  `bson:"created_at" json:"created_at"`
}
