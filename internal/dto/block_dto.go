package dto

type CreateBlockDTO struct {
	Content string `validate:"required,max=1000"`
	Position uint `validate:"required"`
	Images []CreateImageDTO `json:"images"`
}