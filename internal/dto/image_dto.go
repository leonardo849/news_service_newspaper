package dto


type CreateImageDTO struct {
	URL string `json:"url" validate:"required,url"`
	BlockID string `json:"block_id" validate:"required,uuid"`
}