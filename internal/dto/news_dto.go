package dto

import "time"

type CreateNewsDTO struct {
	Authors  []string         `json:"authors" validate:"omitempty,max=6"`
	Title    string           `json:"title"`
	Subtitle string           `json:"subtitle"`
	Topic    string           `json:"topic"`
	Blocks   []CreateBlockDTO `json:"blocks" validate:"min=1,max=20"`
}

type FindNewsDTO struct {
	Authors   []FindAuthorsInFindNews      `json:"authors"`
	Title     string         `json:"title"`
	Subtitle  string         `json:"subtitle"`
	Topic     string         `json:"topic"`
	Blocks    []FindBlockDTO `json:"blocks"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Published_at *time.Time `json:"published_at"`
}