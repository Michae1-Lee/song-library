package main

import (
	"database/sql"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	httpSwagger "github.com/swaggo/http-swagger"
	"log"
	"net/http"
	"os"
	"song-library/api"
	"song-library/controller"
	_ "song-library/docs"
	"song-library/migrations"
	"song-library/repository"
	"song-library/service"
)

// @title			Song Library API
// @version		1.0
// @description	API для управления библиотекой песен
// @termsOfService	http://swagger.io/terms/
func main() {
	// Загрузка .env файла
	if err := godotenv.Load(); err != nil {
		log.Printf("Не удалось загрузить .env файл: %v", err)
	}

	// Чтение конфигурационных переменных из окружения
	dbURL := os.Getenv("DB_URL")
	apiBaseURL := os.Getenv("API_BASE_URL")
	appPort := os.Getenv("APP_PORT")
	if appPort == "" {
		appPort = "8080" // Значение по умолчанию
	}

	// Подключение к базе данных
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Ошибка подключения к базе данных: %v", err)
	}
	defer db.Close()

	// Выполнение миграций
	migrations.RunMigrations(dbURL)

	// Логгер
	logger := log.New(os.Stdout, "SONG-APP: ", log.LstdFlags|log.Lshortfile)

	// Репозиторий, сервис и контроллер
	repo := repository.NewSongRepository(db, logger)
	songService := service.NewSongService(repo, logger, apiBaseURL)
	songController := controller.NewSongController(songService)
	infoController := api.NewInfoController(songService)

	// Настройка маршрутов
	mux := http.NewServeMux()
	mux.HandleFunc("GET /library", songController.GetLibraryHandler)         // Получение библиотеки с фильтрацией и пагинацией
	mux.HandleFunc("GET /song/{id}/text", songController.GetSongTextHandler) // Получение текста песни с пагинацией по куплетам
	mux.HandleFunc("DELETE /song/{id}", songController.DeleteSongHandler)    // Удаление песни
	mux.HandleFunc("PUT /song/{id}", songController.UpdateSongHandler)       // Изменение данных песни
	mux.HandleFunc("POST /song", songController.AddSongHandler)              // Добавление новой песни
	mux.HandleFunc("GET /swagger/", httpSwagger.WrapHandler)
	// Внешний API
	mux.HandleFunc("GET /info", infoController.InfoHandler)

	// Запуск сервера
	server := http.Server{
		Addr:    ":" + appPort,
		Handler: mux,
	}
	log.Printf("Сервер запущен на порту %s", appPort)
	log.Fatal(server.ListenAndServe())
}
