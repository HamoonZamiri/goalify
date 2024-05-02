package service

type UserService interface {
	SignUp(email, password string) error
	Login(email, password string) (string, error)
	Refresh(email, refreshToken string) (string, error)
}
