package handler // Объявляем пакет handler. Здесь будет HTTP-логика.z

import ( // Подключаем пакеты, которые нужны для работы с HTTP и JSON.

	"encoding/json" // Нужен для чтения JSON из запроса и записи JSON в ответ.
	"errors"        // Нужен для проверки типов ошибок через errors.Is.
	"net/http"      // Главный пакет стандартной библиотеки для HTTP-сервера.
	"strconv"       // Нужен, чтобы преобразовать ID из строки в число.
	"strings"       // Нужен для очистки строки от лишних пробелов.

	"github.com/egorik-developer-17/go-api-service/internal/model" // Подключаем модели данных.
	"github.com/egorik-developer-17/go-api-service/internal/store" // Подключаем хранилище товаров.
)

type ProductHandler struct { // ProductHandler — структура, которая знает, как обрабатывать HTTP-запросы.
	store *store.ProductStore // store — ссылка на наше хранилище товаров.
}

func NewProductHandler(store *store.ProductStore) *ProductHandler { // Конструктор ProductHandler.
	return &ProductHandler{store: store} // Создаём handler и сохраняем в него ссылку на store.
}

func (h *ProductHandler) Health(w http.ResponseWriter, r *http.Request) { // Простой health-check endpoint.
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok", "service": "go-api-service", "version": "1.0.0"}) // Возвращаем JSON {"status":"ok"} со статусом 200.
}

func (h *ProductHandler) ListProducts(w http.ResponseWriter, r *http.Request) { // Обработчик получения всего каталога.
	products := h.store.List()            // Забираем список товаров из store.
	writeJSON(w, http.StatusOK, products) // Отправляем список клиенту со статусом 200 OK.
}

func (h *ProductHandler) GetProduct(w http.ResponseWriter, r *http.Request) { // Обработчик получения одного товара по ID.
	id, err := parseID(r) // Пытаемся достать ID из URL и преобразовать его в число.
	if err != nil {       // Если ID не удалось распарсить...
		writeError(w, http.StatusBadRequest, "invalid product id") // ...возвращаем ошибку 400 Bad Request.
		return                                                     // Прерываем выполнение функции.
	}

	product, ok := h.store.GetByID(id) // Ищем товар по ID.
	if !ok {                           // Если товар не найден...
		writeError(w, http.StatusNotFound, "product not found") // ...возвращаем ошибку 404 Not Found.
		return                                                  // Прерываем выполнение функции.
	}

	writeJSON(w, http.StatusOK, product) // Если всё хорошо, возвращаем товар со статусом 200.
}

func (h *ProductHandler) CreateProduct(w http.ResponseWriter, r *http.Request) { // Обработчик создания нового товара.
	var req model.CreateProductRequest // Создаём переменную, в которую будем читать тело запроса.

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil { // Пытаемся декодировать JSON из тела запроса в req.
		writeError(w, http.StatusBadRequest, "invalid request body") // Если JSON неверный — возвращаем 400.
		return                                                       // Прерываем выполнение функции.
	}

	req.Name = strings.TrimSpace(req.Name) // Удаляем пробелы в начале и в конце строки.
	req.Category = strings.TrimSpace(req.Category)

	if req.Name == "" { // Если после удаления пробелов название оказалось пустым...
		writeError(w, http.StatusBadRequest, "name is required") // ...возвращаем 400 и сообщение об ошибке.
		return                                                   // Прерываем выполнение функции.
	}
	if req.Category == "" { // Если после удаления пробелов название оказалось пустым...
		writeError(w, http.StatusBadRequest, "category is required") // ...возвращаем 400 и сообщение об ошибке.
		return                                                       // Прерываем выполнение функции.
	}

	if req.Price <= 0 {
		writeError(w, http.StatusBadRequest, "price must be greater than zero")
		return // прерывание фуекции
	}

	product := h.store.Create(req.Name, req.Category, req.Price) // Создаём новый товар в store.
	writeJSON(w, http.StatusCreated, product)                    // Возвращаем созданный товар со статусом 201 Created.
}

func (h *ProductHandler) UpdateProduct(w http.ResponseWriter, r *http.Request) { // Обработчик обновления названия товара.
	id, err := parseID(r) // Достаём ID товара из URL.
	if err != nil {       // Если ID неправильный...
		writeError(w, http.StatusBadRequest, "invalid product id") // ...возвращаем 400.
		return                                                     // Прерываем выполнение функции.
	}

	var req model.UpdateProductRequest // Создаём переменную для JSON тела запроса.

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil { // Пытаемся прочитать JSON тела запроса.
		writeError(w, http.StatusBadRequest, "invalid request body") // Если JSON плохой — возвращаем 400.
		return                                                       // Прерываем выполнение функции.
	}

	req.Name = strings.TrimSpace(req.Name) // Убираем лишние пробелы из названия.
	req.Category = strings.TrimSpace(req.Category)

	if req.Name == "" { // Если название пустое...
		writeError(w, http.StatusBadRequest, "name is required") // ...возвращаем 400.
		return                                                   // Прерываем выполнение функции.
	}
	if req.Category == "" { // Если название пустое...
		writeError(w, http.StatusBadRequest, "category is required") // ...возвращаем 400.
		return                                                       // Прерываем выполнение функции.
	}
	if req.Price <= 0 {
		writeError(w, http.StatusBadRequest, "price must be greater than zero")
		return
	}

	product, err := h.store.UpdateName(id, req.Name, req.Category, req.Price) // Пытаемся обновить название товара в store.
	if err != nil {                                                           // Если произошла ошибка...
		if errors.Is(err, store.ErrProductNotFound) { // ...и эта ошибка означает, что товар не найден...
			writeError(w, http.StatusNotFound, "product not found") // ...возвращаем 404.
			return                                                  // Прерываем выполнение функции.
		}

		writeError(w, http.StatusInternalServerError, "internal server error") // Для всех прочих ошибок возвращаем 500.
		return                                                                 // Прерываем выполнение функции.
	}

	writeJSON(w, http.StatusOK, product) // Если обновление прошло успешно, возвращаем обновлённый товар.
}

func (h *ProductHandler) DeleteProduct(w http.ResponseWriter, r *http.Request) { // Обработчик удаления товара.
	id, err := parseID(r) // Достаём ID из URL.
	if err != nil {       // Если ID неверный...
		writeError(w, http.StatusBadRequest, "invalid product id") // ...возвращаем 400.
		return                                                     // Прерываем выполнение функции.
	}

	if err := h.store.Delete(id); err != nil { // Пытаемся удалить товар из store.
		if errors.Is(err, store.ErrProductNotFound) { // Если store сказал, что товар не найден...
			writeError(w, http.StatusNotFound, "product not found") // ...возвращаем 404.
			return                                                  // Прерываем выполнение функции.
		}

		writeError(w, http.StatusInternalServerError, "internal server error") // Если ошибка другая — возвращаем 500.
		return                                                                 // Прерываем выполнение функции.
	}

	w.WriteHeader(http.StatusNoContent) // Если удаление прошло успешно, возвращаем статус 204 No Content без тела ответа.
}

func parseID(r *http.Request) (int, error) { // Вспомогательная функция для чтения ID из URL.
	id, err := strconv.Atoi(r.PathValue("id")) // Берём значение параметра {id} из пути и превращаем его в int.
	if err != nil {                            // Если строку не удалось преобразовать в число...
		return 0, err // ...возвращаем 0 и ошибку.
	}

	return id, nil // Если всё хорошо, возвращаем готовый числовой ID.
}

func writeJSON(w http.ResponseWriter, status int, data any) { // Универсальная функция для записи JSON-ответа.
	w.Header().Set("Content-Type", "application/json") // Сообщаем клиенту, что ответ будет в формате JSON.
	w.WriteHeader(status)                              // Устанавливаем HTTP-статус ответа.

	if err := json.NewEncoder(w).Encode(data); err != nil { // Пытаемся сериализовать data в JSON и отправить клиенту.
		http.Error(w, "failed to encode response", http.StatusInternalServerError) // Если сериализация не удалась, отправляем 500.
	}
}

func writeError(w http.ResponseWriter, status int, message string) { // Вспомогательная функция для единообразных ошибок.
	writeJSON(w, status, map[string]string{"error": message}) // Возвращаем JSON вида {"error":"текст ошибки"}.
}
