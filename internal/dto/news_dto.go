package dto


type CreateNewsDTO struct {
	Authors []string `json:"authors" validate:"required"`
	Title string `json:"title"`
	Subtitle string `json:"subtitle"`
	Topic string `json:"topic"`
	News []CreateNewsDTO `json:"news"`
}