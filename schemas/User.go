package schemas

type UserCreate struct {
    Username string `json:"username" validate:"required,min=3"`
    Email    string `json:"email" validate:"required,email"`
    Password string `json:"password" validate:"required,min=6"`
}

type UserUpdate struct {
    Username string `json:"username" validate:"omitempty,min=3"`
    Email    string `json:"email" validate:"omitempty,email"`
    Password string `json:"password" validate:"omitempty,min=6"`
}