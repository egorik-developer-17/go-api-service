package main

import (
	"log"
	"net/http"
	"os"

	"github.com/egorik-developer-17/go-api-service/internal/handler"
	"github.com/egorik-developer-17/go-api-service/internal/server"
	"github.com/egorik-developer-17/go-api-service/internal/store"
)

func main() { // Главная функция. С неё начинается выполнение программы.
	productStore := store.NewProductStore()                   // Создаём хранилище товаров в памяти.
	productHandler := handler.NewProductHandler(productStore) // Создаём handler и передаём ему store.
	router := server.NewRouter(productHandler)                // Собираем маршруты API.

	port := os.Getenv("HTTP_PORT") // Пытаемся прочитать порт из переменной окружения HTTP_PORT.
	if port == "" {                // Если переменная окружения не задана...
		port = "8080" // ...используем стандартный порт 8080.
	}

	log.Printf("server started on :%s", port) // Пишем в лог, что сервер стартует на выбранном порту.

	if err := http.ListenAndServe(":"+port, router); err != nil { // Запускаем HTTP-сервер на порту и передаём ему router.
		log.Fatal(err) // Если сервер завершился с ошибкой, выводим её и завершаем программу.
	}
}
