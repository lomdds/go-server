package models

import "gorm.io/gorm"

//создание таблицы статей
type Article struct {
    gorm.Model
    Title string `gorm:"not null"`
    Content string `gorm:"not null"`
    UserID uint `gorm:"not null"`
    User User `gorm:"foreignKey:UserID"`
}