package models

import "time"

type PostUpdate struct {
	Message string `json:"message"`
}

type Post struct {
	Id          uint64    `json:"id"`
	Parent      uint64    `json:"parent"`
	Author      string    `json:"author"`
	Message     string    `json:"message"`
	IsEdited    bool      `json:"isEdited"`
	Forum       string    `json:"forum"`
	Thread      uint64    `json:"thread"`
	DateCreated time.Time `json:"created"`
}

type PostCreate struct {
	Parent  uint64 `json:"parent"`
	Author  string `json:"author"`
	Message string `json:"message"`
}

type PostDetails struct {
	Post   *Post   `json:"post"`
	Author *User   `json:"author"`
	Thread *Thread `json:"thread"`
	Forum  *Forum  `json:"forum"`
}

type PostPaginator struct {
	Limit     uint64 `json:"limit"`
	Since     uint64 `json:"since"`
	Sort      string `json:"sort"`
	SortOrder bool   `json:"desc"`
}
