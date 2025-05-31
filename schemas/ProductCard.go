package schemas

type ProductCardCreate struct {
    Brand           string `json:"brand" validate:"required,min=2"`
    BikeModel       string `json:"bikeModel" validate:"required,min=1"`
    EngineCapacity  int    `json:"engineCapacity" validate:"required,gt=0"`
    Power           int    `json:"power" validate:"required,gt=0"`
    Color           string `json:"color" validate:"required,min=3"`
    Price           int    `json:"price" validate:"required,gt=0"`
}

type ProductCardUpdate struct {
    Brand           string `json:"brand" validate:"omitempty,min=2"`
    BikeModel       string `json:"bikeModel" validate:"omitempty,min=1"`
    EngineCapacity  int    `json:"engineCapacity" validate:"omitempty,gt=0"`
    Power           int    `json:"power" validate:"omitempty,gt=0"`
    Color           string `json:"color" validate:"omitempty,min=3"`
    Price           int    `json:"price" validate:"omitempty,gt=0"`
}