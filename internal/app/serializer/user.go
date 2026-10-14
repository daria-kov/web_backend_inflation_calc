package serializer

type RegisterUserRequest struct {
	Login    string `json:"login" binding:"required"`
	FullName string `json:"full_name"`
	Password string `json:"password"`
}

type UserResponse struct {
	ID       uint   `json:"id"`
	Login    string `json:"login"`
	FullName string `json:"full_name"`
}

type AuthUserRequest struct {
	Login    string `json:"login" binding:"required"`
	Password string `json:"password" binding:"required"`
}
