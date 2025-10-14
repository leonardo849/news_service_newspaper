package dto


type CreateNewsDTO struct {
	Authors []string `json:"authors" validate:"omitempty,max=6"`
	Title string `json:"title"`
	Subtitle string `json:"subtitle"`
	Topic string `json:"topic"`
	Blocks []CreateBlockDTO `json:"blocks" validate:"min=1,max=20"`
}