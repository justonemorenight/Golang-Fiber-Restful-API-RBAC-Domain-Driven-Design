package models

type User struct {
	BaseModel
	Name     string `json:"name" validate:"required,min=3"`
	Email    string `json:"email" validate:"required,email" gorm:"unique"`
	Password string `json:"password" validate:"required,min=6"`
}
