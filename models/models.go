package models

const (
	RoleAdmin = "admin"
	RoleUser  = "user"
	JwtKey    = "SecretKey"
)

type User struct {
	UserName string `json:"username"`
	Password string `json:"password"`
	Role     string `json:"role"`
}
