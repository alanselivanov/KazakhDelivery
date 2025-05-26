package fixtures

type ProductFixture struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	CategoryID  string  `json:"category_id"`
	Stock       int     `json:"stock"`
}

func GetTestProducts() []ProductFixture {
	return []ProductFixture{
		{
			ID:          "prod-1",
			Name:        "Laptop",
			Description: "High-end laptop for developers",
			Price:       1299.99,
			CategoryID:  "cat-electronics",
			Stock:       50,
		},
		{
			ID:          "prod-2",
			Name:        "Keyboard",
			Description: "Mechanical gaming keyboard",
			Price:       129.99,
			CategoryID:  "cat-electronics",
			Stock:       100,
		},
		{
			ID:          "prod-3",
			Name:        "Mouse",
			Description: "Wireless gaming mouse",
			Price:       79.99,
			CategoryID:  "cat-electronics",
			Stock:       150,
		},
		{
			ID:          "prod-4",
			Name:        "Monitor",
			Description: "32-inch 4K monitor",
			Price:       399.99,
			CategoryID:  "cat-electronics",
			Stock:       30,
		},
		{
			ID:          "prod-5",
			Name:        "Headphones",
			Description: "Noise-cancelling headphones",
			Price:       249.99,
			CategoryID:  "cat-electronics",
			Stock:       75,
		},
	}
}

func GetProductByID(id string) *ProductFixture {
	for _, product := range GetTestProducts() {
		if product.ID == id {
			return &product
		}
	}
	return nil
}
