package dto

type CreateBlockDTO struct {
	Content string `validate:"required,max=1000"`
	Position uint `validate:"required"`
	NewsID string `validate:"required,uuid"`
	Images []CreateImageDTO `json:"images"`
}