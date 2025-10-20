package dto


type CreateUserFromJsonFileDTO struct {
	Username string `json:"username" validate:"required,max=50"`
	Role string `json:"role" validate:"required,role"`
}

type FindAuthorsInFindNews struct {
	Username string `json:"username"`
}