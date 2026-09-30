package serializer

type RegisterUserRequest struct {
	FullName string `json:"full_name"`
	Email    string `json:"email"`
}

type UserResponse struct {
	ID       uint   `json:"id"`
	FullName string `json:"full_name"`
	Email    string `json:"email"`
}
