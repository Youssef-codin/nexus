package auth

type registerRequest struct {
	Email    string `json:"email"     validate:"required,email"`
	Password string `json:"password"  validate:"required"`
	FullName string `json:"full_name" validate:"required,min=3"`
}

type registerResponse struct {
	Email        string `json:"email"`
	FullName     string `json:"full_name"`
	JwtToken     string `json:"jwt_token"`
	RefreshToken string `json:"refresh_token"`
}

type loginResponse struct {
	Email        string `json:"email"`
	FullName     string `json:"full_name"`
	JwtToken     string `json:"jwt_token"`
	RefreshToken string `json:"refresh_token"`
}

type loginRequest struct {
	Email    string `json:"email"    validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type refreshRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

type refreshResponse struct {
	JwtToken     string `json:"jwt_token"`
	RefreshToken string `json:"refresh_token"`
}
