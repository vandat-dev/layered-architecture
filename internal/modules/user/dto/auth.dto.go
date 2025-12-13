package dto

// LoginRequestDto represents the login request structure
type LoginRequestDto struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// RegisterRequestDto represents the registration request structure
type RegisterRequestDto struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	Role     string `json:"role" binding:"required,oneof=ADMIN USER SUPER_ADMIN"`
}

// AuthResponseDto represents the authentication response
type AuthResponseDto struct {
	Token        string `json:"token"`
	RefreshToken string `json:"refresh_token"`
}
