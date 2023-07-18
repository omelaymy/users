package auth

type Usecase interface {
	Authentication(username, password string) bool
	AdminAuthorization(username, password string) bool
}
