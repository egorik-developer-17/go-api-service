package model

type Product struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type CreateProductRequest struct {
	Name string `json:"name" validate:"required"`
}

type UpdateProductRequest struct {
	Name string `json:"name"`
}
