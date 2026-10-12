package serializer

type RegisterUserRequest struct {
	FullName string `json:"full_name"`
	Password string `json:"password"`
}

type UserResponse struct {
	ID       uint   `json:"id"`
	FullName string `json:"full_name"`
}
