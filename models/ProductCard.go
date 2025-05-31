package models

import "gorm.io/gorm"

type ProductCard struct {
    gorm.Model
    UserID          uint   `gorm:"not null;index"`
    User            User   `gorm:"foreignKey:UserID"`
    Brand           string `gorm:"not null"`
    BikeModel       string `gorm:"not null"`
    EngineCapacity  int    `gorm:"not null"`
    Power           int    `gorm:"not null"`
    Color           string `gorm:"not null"`     
    Price           int    `gorm:"not null"`
}