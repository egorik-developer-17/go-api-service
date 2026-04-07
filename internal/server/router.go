package server

import (
	"net/http"

	"github.com/egorik-developer-17/go-api-service/internal/handler"
)

func NewRouter(productHandler *handler.ProductHandler) http.Handler {
	mux := http.NewServeMux()

	staticFiles := http.FileServer(http.Dir("./web/static"))
	mux.Handle("GET /static/", http.StripPrefix("/static/", staticFiles))

	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}

		http.ServeFile(w, r, "./web/static/index.html")
	})

	mux.HandleFunc("GET /health", productHandler.Health)
	mux.HandleFunc("GET /products", productHandler.ListProducts)
	mux.HandleFunc("GET /products/{id}", productHandler.GetProduct)
	mux.HandleFunc("POST /products", productHandler.CreateProduct)
	mux.HandleFunc("PUT /products/{id}", productHandler.UpdateProduct)
	mux.HandleFunc("DELETE /products/{id}", productHandler.DeleteProduct)

	return mux
}
