package domain

type User struct {
	Id       int64  `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Nickname string `json:"nickname"`
	Birthday string `json:"birthday"`
	AboutMe  string `json:"aboutMe"`
}
