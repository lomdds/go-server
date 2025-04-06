package schemas

type ArticleCreate struct {
    Title   string `json:"title" validate:"required,min=5"`
    Content string `json:"content" validate:"required,min=10"`
    UserID  uint   `json:"user_id" validate:"required"`
}