package auth

type LoginBody struct {
	Login    string `json:"login" validate:"required"`
	Password string `json:"password" validate:"required"`
}
