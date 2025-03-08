package main

import (
	"gorm.io/gorm"
	"gorm.io/driver/postgres"
)

//создание таблицы пользователей
type User struct {
    gorm.Model
    Username string `gorm:"not null;unique"`
    Email string `gorm:"not null;unique"`
    Password string `gorm:"not null"`
    Articles []Article
}

//создание таблицы статей
type Article struct {
    gorm.Model
    Title string `gorm:"not null"`
    Content string `gorm:"not null"`
    UserID uint `gorm:"not null"`
    User User `gorm:"foreignKey:UserID"`
}

func main() {
    dsn := "host=localhost user=gorm password=gorm dbname=gorm port=5432 sslmode=disable"
    db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
    if err != nil {
        panic("Не удалось присоебиниться к бд")
    }

    user := User{Username: "John Doe", Email: "ololo@gmail.com", Password: "ololo123"}
    result := db.Create(&user)

    if result.Error != nil {
        panic("Ошибка")
    }

    db.Model(&user).Updates(User{Username: "Mussolini", Email: "Italy@gmail.com", Password: "123"})

    db.Delete(&user, 1)

    db.AutoMigrate(&User{}, &Article{})

}