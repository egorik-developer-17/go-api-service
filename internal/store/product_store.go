package store

import (
	"errors"
	"sync"

	"github.com/egorik-developer-17/go-api-service/internal/model"
)

var ErrProductNotFound = errors.New("product not found")

type ProductStore struct {
	mu       sync.RWMutex    //  mutex защищает данные от одновременного доступа
	products []model.Product // слайс для хранения продуктов
	nextID   int             // для генерации уникальных ID
}

func NewProductStore() *ProductStore { // конструктор для создания нового экземпляра ProductStore
	return &ProductStore{ // возвращаем указатель на новый ProductStore
		products: make([]model.Product, 0), // инициализируем слайс продуктов
		nextID:   1,                        // начинаем ID с 1
	}
}

func (s *ProductStore) List() []model.Product { // метод для получения списка всех продуктов
	s.mu.RLock()         // блокируем для чтения
	defer s.mu.RUnlock() // разблокируем после чтения

	result := make([]model.Product, len(s.products)) // создаем новый слайс для результатов
	copy(result, s.products)                         // копируем данные из основного слайса в результат
	return result                                    // возвращаем результат
}

func (s *ProductStore) GetByID(id int) (model.Product, bool) { // метод для получения продукта по ID
	s.mu.RLock()         // блокируем для чтения
	defer s.mu.RUnlock() // разблокируем после чтения

	for _, product := range s.products {
		if product.ID == id {
			return product, true
		}
	}
	return model.Product{}, false
}

func (s *ProductStore) Create(name string, category string, price float64) model.Product { // метод для создания нового продукта
	s.mu.Lock()               // блокируем для записи
	defer s.mu.Unlock()       // разблокируем после записи
	product := model.Product{ // создаем новый продукт
		ID:       s.nextID, // присваиваем уникальный ID
		Name:     name,     // устанавливаем имя продукта
		Category: category,
		Price:    price,
	}
	s.products = append(s.products, product) // добавляем новый продукт в слайс
	s.nextID++
	return product // возвращаем созданный продукт
}

func (s *ProductStore) UpdateName(id int, name string, category string, price float64) (model.Product, error) { // метод для обновления имени продукта по ID
	s.mu.Lock()         // блокируем для записи
	defer s.mu.Unlock() // разблокируем после записи

	for i, product := range s.products { // проходим по слайсу продуктов
		if product.ID == id { // если ID совпало
			s.products[i].Name = name // обновляем имя продукта
			s.products[i].Category = category
			s.products[i].Price = price
			return s.products[i], nil // возвращаем обновленный продукт
		}
	}
	return model.Product{}, ErrProductNotFound // если продукт не найден, возвращаем ошибку
}

func (s *ProductStore) Delete(id int) error { // Метод Delete удаляет товар по ID.
	s.mu.Lock()         // Ставим lock, потому что будем менять срез.
	defer s.mu.Unlock() // После завершения функции снимаем блокировку.

	for i, product := range s.products { // Проходим по всем товарам.
		if product.ID == id { // Если нашли товар с нужным ID...
			s.products = append(s.products[:i], s.products[i+1:]...) // ...вырезаем его из среза.
			return nil                                               // Возвращаем nil, то есть удаление прошло успешно.
		}
	}

	return ErrProductNotFound // Если товар не найден, возвращаем ошибку.
}
