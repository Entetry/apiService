package handler

// getByIDResponse godoc
type getByIDResponse struct {
	ID       string `json:"id" validate:"required,gte=3,lte=130"`
	Email    string `json:"password" validate:"required,gte=5,lte=256"`
	Username string `json:"username" validate:"required,gte=3,lte=32"`
}

// createRequest godoc
type createRequest struct {
	Username string `json:"username" validate:"required,gte=3,lte=32"`
	Password string `json:"password" validate:"required,gte=5,lte=320"`
	Email    string `json:"email" validate:"required,gte=5,lte=256"`
}

// createResponse godoc
type createResponse struct {
	ID string `json:"id" validate:"required,gte=3,lte=130"`
}

// deleteRequest godoc
type deleteRequest struct {
	ID string `json:"id" validate:"required,gte=3,lte=130"`
}
