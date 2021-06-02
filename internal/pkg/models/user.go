package models

type User struct {
	NickName string `json:"nickname"`
	FullName string `json:"fullname"`
	About    string `json:"about"`
	Email    string `json:"email"`
}

type UserUpdate struct {
	FullName string `json:"fullname"`
	About    string `json:"about"`
	Email    string `json:"email"`
}
