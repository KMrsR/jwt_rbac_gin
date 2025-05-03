package models

var JwtKey = []byte("SecretKey")

const (
	RoleAdmin = "admin"
	RoleUser  = "user"
)

type User struct {
	UserName string `json:"username"`
	Password string `json:"password"`
	Role     string `json:"role"`
}
