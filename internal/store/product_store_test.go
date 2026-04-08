package store // Указываем, что тесты относятся к пакету store.

import "testing" // Подключаем стандартный пакет testing для написания тестов.

func TestProductStoreCreate(t *testing.T) { // Тест проверяет создание нового товара.
	s := NewProductStore() // Создаём новое тестовое хранилище.

	product := s.Create("Ячмень", "зерновые", 1450.25) // Создаём новый товар через store.

	if product.ID == 0 { // Проверяем, что товар получил ID.
		t.Fatal("expected product ID to be assigned") // Если ID не назначен, тест падает.
	}

	if product.Name != "Ячмень" { // Проверяем, что имя сохранилось правильно.
		t.Fatalf("expected name to be Ячмень, got %s", product.Name) // Если имя не совпало, тест падает.
	}

	if product.Category != "зерновые" { // Проверяем, что категория сохранилась правильно.
		t.Fatalf("expected category to be зерновые, got %s", product.Category) // Если категория не совпала, тест падает.
	}

	if product.Price != 1450.25 { // Проверяем, что цена сохранилась правильно.
		t.Fatalf("expected price to be 1450.25, got %f", product.Price) // Если цена не совпала, тест падает.
	}
}

func TestProductStoreUpdate(t *testing.T) { // Тест проверяет обновление существующего товара.
	s := NewProductStore() // Создаём новое тестовое хранилище.

	created := s.Create("Пшеница", "зерновые", 1200.50) // Сначала сами создаём товар, который потом будем обновлять.

	updated, err := s.UpdateName(created.ID, "Овес", "зерновые", 1600.00) // Обновляем товар по реальному ID, который получили после создания.
	if err != nil {                                                       // Если store вернул ошибку...
		t.Fatalf("expected no error, got %v", err) // ...тест падает.
	}

	if updated.Name != "Овес" { // Проверяем, что имя действительно обновилось.
		t.Fatalf("expected updated name to be Овес, got %s", updated.Name) // Если имя не обновилось, тест падает.
	}

	if updated.Category != "зерновые" { // Проверяем, что категория обновилась корректно.
		t.Fatalf("expected updated category to be зерновые, got %s", updated.Category) // Если категория не совпала, тест падает.
	}

	if updated.Price != 1600.00 { // Проверяем, что цена обновилась корректно.
		t.Fatalf("expected updated price to be 1600.00, got %f", updated.Price) // Если цена не совпала, тест падает.
	}
}

func TestProductStoreDelete(t *testing.T) { // Тест проверяет удаление существующего товара.
	s := NewProductStore() // Создаём новое тестовое хранилище.

	created := s.Create("Подсолнечник", "масличные", 2100.75) // Сначала создаём товар, который потом будем удалять.

	err := s.Delete(created.ID) // Удаляем товар по реальному ID созданного товара.
	if err != nil {             // Если store вернул ошибку...
		t.Fatalf("expected no error, got %v", err) // ...тест падает.
	}

	_, ok := s.GetByID(created.ID) // Пытаемся снова получить удалённый товар.
	if ok {                        // Если товар всё ещё существует...
		t.Fatal("expected product to be deleted") // ...тест падает.
	}
}
