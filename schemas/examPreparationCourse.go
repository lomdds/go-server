package schemas

type ExamPreparationCourse struct {
	ID                uint      `json:"id"`
	Subject           string    `json:"subject" validate:"required"`
	UserID            uint      `json:"user_id" validate:"required"`
	Relevance         int       `json:"relevance" validate:"required"`
	NumberOfClasses   int       `json:"number_of_classes" validate:"gte=1"`
	ContactTheTeacher bool      `json:"contact_the_teacher"`
	Individuality     bool      `json:"individuality"`
	Price             int       `json:"price" validate:"required,gt=0"`
}