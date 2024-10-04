package utils

type NotFoundErrorResponse struct {
	Status  int    `json:"status" example:"404"`
	Message string `json:"message" example:"URL not found"`
}

type InternalServerErrorResponse struct {
	Status  int    `json:"status" example:"500"`
	Message string `json:"message" example:"Internal server error"`
}

type ConflictErrorResponse struct {
	Status  int    `json:"status" example:"409"`
	Message string `json:"message" example:"Conflict"`
}

type BadRequestErrorResponse struct {
	Status  int    `json:"status" example:"400"`
	Message string `json:"message" example:"Bad request"`
}
