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

type UserPaginator struct {
	Limit     uint64 `json:"limit"`
	Since     string `json:"since"`
	SortOrder bool   `json:"desc"`
}
