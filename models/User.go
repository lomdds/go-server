package models

import "gorm.io/gorm"

// создание таблицы пользователей
type User struct {
    gorm.Model
    Username     string        `gorm:"not null;unique"`
    Email        string        `gorm:"not null;unique"`
    Password     string        `gorm:"not null"`
    ProductCards []ProductCard `gorm:"foreignKey:UserID"`
}
