package model

type Product struct {
	ID       int     `json:"id"`
	Name     string  `json:"name"`
	Category string  `json:"category"`
	Price    float64 `jsone:"price"`
}

type CreateProductRequest struct {
	Name     string  `json:"name" validate:"required"`
	Category string  `json:"category"`
	Price    float64 `jsone:"price"`
}

type UpdateProductRequest struct {
	Name     string  `json:"name"`
	Category string  `json:"category"`
	Price    float64 `jsone:"price"`
}
