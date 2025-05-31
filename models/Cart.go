package models

import "gorm.io/gorm"

type Cart struct {
	gorm.Model
	UserID uint `gorm:"not null;uniqueIndex"`
	User User `gorm:"foreignKey:UserID"`
	CartItems []CartItem `gorm:"foreignKey:CartID"`
}