package dto


type CreateNewsDTO struct {
	Authors []string `json:"authors" validate:"required,min=1,max=15"`
	Title string `json:"title"`
	Subtitle string `json:"subtitle"`
	Topic string `json:"topic"`
	Blocks []CreateBlockDTO `json:"news" validate:"min=1,max=20"`
}