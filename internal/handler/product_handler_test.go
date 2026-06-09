package handler_test

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/egorik-developer-17/go-api-service/internal/handler"
	"github.com/egorik-developer-17/go-api-service/internal/model"
	"github.com/egorik-developer-17/go-api-service/internal/server"
	"github.com/egorik-developer-17/go-api-service/internal/store"
)

type errorResponse struct {
	Error string `json:"error"`
}

func newTestRouter() (*store.ProductStore, http.Handler) {
	s := store.NewProductStore()
	h := handler.NewProductHandler(s)
	return s, server.NewRouter(h)
}

func performRequest(t *testing.T, router http.Handler, method, path string, body any) *httptest.ResponseRecorder {
	t.Helper()

	var reader io.Reader

	switch v := body.(type) {
	case nil:
		reader = nil
	case string:
		reader = strings.NewReader(v)
	default:
		data, err := json.Marshal(v)
		if err != nil {
			t.Fatalf("marshal request body: %v", err)
		}
		reader = bytes.NewReader(data)
	}

	req := httptest.NewRequest(method, path, reader)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	return rec
}

func decodeBody[T any](t *testing.T, rec *httptest.ResponseRecorder) T {
	t.Helper()

	var out T
	if err := json.Unmarshal(rec.Body.Bytes(), &out); err != nil {
		t.Fatalf("decode response body: %v; body=%s", err, rec.Body.String())
	}

	return out
}

func TestCreateProduct_Success(t *testing.T) {
	s, router := newTestRouter()

	rec := performRequest(t, router, http.MethodPost, "/products", model.CreateProductRequest{
		Name:     "  Молоко  ",
		Category: "  dairy  ",
		Price:    99.90,
	})

	if rec.Code != http.StatusCreated {
		t.Fatalf("expected status %d, got %d; body=%s", http.StatusCreated, rec.Code, rec.Body.String())
	}

	product := decodeBody[model.Product](t, rec)

	if product.ID != 1 {
		t.Fatalf("expected ID 1, got %d", product.ID)
	}
	if product.Name != "Молоко" {
		t.Fatalf("expected trimmed name, got %q", product.Name)
	}
	if product.Category != "dairy" {
		t.Fatalf("expected trimmed category, got %q", product.Category)
	}
	if product.Price != 99.90 {
		t.Fatalf("expected price 99.90, got %v", product.Price)
	}

	stored, ok := s.GetByID(product.ID)
	if !ok {
		t.Fatal("expected product to be stored")
	}
	if stored != product {
		t.Fatalf("expected stored product %+v, got %+v", product, stored)
	}
}

func TestCreateProduct_Validation(t *testing.T) {
	tests := []struct {
		name       string
		body       any
		wantStatus int
		wantError  string
	}{
		{
			name:       "invalid json",
			body:       `{"name":`,
			wantStatus: http.StatusBadRequest,
			wantError:  "invalid request body",
		},
		{
			name: "empty name",
			body: model.CreateProductRequest{
				Name:     "   ",
				Category: "food",
				Price:    100,
			},
			wantStatus: http.StatusBadRequest,
			wantError:  "name is required",
		},
		{
			name: "empty category",
			body: model.CreateProductRequest{
				Name:     "Bread",
				Category: "   ",
				Price:    100,
			},
			wantStatus: http.StatusBadRequest,
			wantError:  "category is required",
		},
		{
			name: "invalid price",
			body: model.CreateProductRequest{
				Name:     "Bread",
				Category: "food",
				Price:    0,
			},
			wantStatus: http.StatusBadRequest,
			wantError:  "price must be greater than zero",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s, router := newTestRouter()

			rec := performRequest(t, router, http.MethodPost, "/products", tt.body)

			if rec.Code != tt.wantStatus {
				t.Fatalf("expected status %d, got %d; body=%s", tt.wantStatus, rec.Code, rec.Body.String())
			}

			errResp := decodeBody[errorResponse](t, rec)
			if errResp.Error != tt.wantError {
				t.Fatalf("expected error %q, got %q", tt.wantError, errResp.Error)
			}

			if got := len(s.List()); got != 0 {
				t.Fatalf("expected store to stay empty, got %d products", got)
			}
		})
	}
}

func TestUpdateProduct_Success(t *testing.T) {
	s, router := newTestRouter()
	created := s.Create("Старое имя", "old-category", 100)

	rec := performRequest(t, router, http.MethodPut, "/products/1", model.UpdateProductRequest{
		Name:     "  Новое имя  ",
		Category: "  new-category  ",
		Price:    150.50,
	})

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d; body=%s", http.StatusOK, rec.Code, rec.Body.String())
	}

	updated := decodeBody[model.Product](t, rec)

	if updated.ID != created.ID {
		t.Fatalf("expected ID %d, got %d", created.ID, updated.ID)
	}
	if updated.Name != "Новое имя" {
		t.Fatalf("expected trimmed name, got %q", updated.Name)
	}
	if updated.Category != "new-category" {
		t.Fatalf("expected trimmed category, got %q", updated.Category)
	}
	if updated.Price != 150.50 {
		t.Fatalf("expected price 150.50, got %v", updated.Price)
	}

	stored, ok := s.GetByID(created.ID)
	if !ok {
		t.Fatal("expected updated product to exist in store")
	}
	if stored != updated {
		t.Fatalf("expected stored product %+v, got %+v", updated, stored)
	}
}

func TestUpdateProduct_Validation(t *testing.T) {
	tests := []struct {
		name       string
		path       string
		body       any
		wantStatus int
		wantError  string
	}{
		{
			name:       "invalid id",
			path:       "/products/abc",
			body:       model.UpdateProductRequest{Name: "Test", Category: "food", Price: 100},
			wantStatus: http.StatusBadRequest,
			wantError:  "invalid product id",
		},
		{
			name:       "invalid json",
			path:       "/products/1",
			body:       `{"name":`,
			wantStatus: http.StatusBadRequest,
			wantError:  "invalid request body",
		},
		{
			name:       "empty name",
			path:       "/products/1",
			body:       model.UpdateProductRequest{Name: "   ", Category: "food", Price: 100},
			wantStatus: http.StatusBadRequest,
			wantError:  "name is required",
		},
		{
			name:       "empty category",
			path:       "/products/1",
			body:       model.UpdateProductRequest{Name: "Test", Category: "   ", Price: 100},
			wantStatus: http.StatusBadRequest,
			wantError:  "category is required",
		},
		{
			name:       "invalid price",
			path:       "/products/1",
			body:       model.UpdateProductRequest{Name: "Test", Category: "food", Price: 0},
			wantStatus: http.StatusBadRequest,
			wantError:  "price must be greater than zero",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s, router := newTestRouter()
			original := s.Create("Пшеница", "зерновые", 1000)

			rec := performRequest(t, router, http.MethodPut, tt.path, tt.body)

			if rec.Code != tt.wantStatus {
				t.Fatalf("expected status %d, got %d; body=%s", tt.wantStatus, rec.Code, rec.Body.String())
			}

			errResp := decodeBody[errorResponse](t, rec)
			if errResp.Error != tt.wantError {
				t.Fatalf("expected error %q, got %q", tt.wantError, errResp.Error)
			}

			stored, ok := s.GetByID(original.ID)
			if !ok {
				t.Fatal("expected original product to remain in store")
			}
			if stored != original {
				t.Fatalf("product changed after invalid update: want %+v, got %+v", original, stored)
			}
		})
	}
}
