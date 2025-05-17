package models

import "gorm.io/gorm"

type ExamPreparationCourse struct {
	gorm.Model
	Subject           string `gorm:"not null"`
	UserID            uint   `gorm:"not null"`
	User              User   `gorm:"foreignKey:UserID"`
	Relevance         int    `gorm:"not null"`
	NumberOfClasses   int    `gorm:"not null"`
	ContactTheTeacher bool	 
	Individuality     bool	 
	Price 			  int    `gorm:"not null"`
}