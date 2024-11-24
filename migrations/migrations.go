package migrations

import (
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"log"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func RunMigrations(dbURL string) {
	m, err := migrate.New(
		"file://migrations", // Путь к папке с миграциями
		dbURL,               // URL подключения к базе данных
	)
	if err != nil {
		log.Fatalf("Ошибка инициализации миграции: %v", err)
	}

	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		log.Fatalf("Ошибка выполнения миграции: %v", err)
	}

	fmt.Println("Миграции успешно применены")
}
