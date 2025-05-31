package models

import "gorm.io/gorm"

type CartItem struct {
	gorm.Model
	CartID uint `gorm:"not null"`
	ProductCardID uint `gorm:"not null"`
	Quantity int `gorm:"not null;default:1"`
	Cart Cart `gorm:"foreignKey:CartID"`
	ProductCard ProductCard `gorm:"foreignKey:ProductCardID`
}