package main

import (
    "log"
    "go-server/models"
	"gorm.io/gorm"
	"gorm.io/driver/postgres"
)

func main() {
    dsn := "host=localhost user=gorm password=gorm dbname=gorm port=5432 sslmode=disable"
    db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
    if err != nil {
        log.Fatal("Не удалось подключиться к БД:", err)
    }

    // Пример создания пользователя
    user := models.User{
        Username: "John Doe",
        Email:    "ololo@gmail.com",
        Password: "ololo123",
    }
    if err := db.Create(&user).Error; err != nil {
        log.Fatal("Ошибка создания пользователя:", err)
    }

    // Обновление пользователя
    if err := db.Model(&user).Updates(models.User{
        Username: "Mussolini",
        Email:    "Italy@gmail.com",
        Password: "123",
    }).Error; err != nil {
        log.Fatal("Ошибка обновления:", err)
    }

    // Удаление пользователя
    if err := db.Delete(&user).Error; err != nil {
        log.Fatal("Ошибка удаления:", err)
    }
}