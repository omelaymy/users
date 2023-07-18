package auth

type User struct {
	Username string
	Password string
	Admin    bool
}

type Credentials struct {
	Username string
	Password string
}
