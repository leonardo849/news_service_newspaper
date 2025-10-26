package dto

type CreateBlockDTO struct {
	Content string `validate:"required,max=1000" json:"content"`
	Position uint `validate:"required" json:"position"`
	Images []CreateImageDTO `json:"images"`
}

type CreateBlockDTOJSON struct {
	Blocks []CreateBlockDTO `json:"blocks"`
}

type FindBlockDTO struct {
	ID string `json:"id"`
	Position uint `json:"position"`
	Content string `json:"content"`
	Images []FindImageDTO `json:"images"`
}