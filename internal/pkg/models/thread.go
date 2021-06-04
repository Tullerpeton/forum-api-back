package models

import "time"

type ThreadCreate struct {
	Title          string    `json:"title"`
	AuthorNickName string    `json:"user"`
	Message        string    `json:"message"`
	DateCreated    time.Time `json:"created"`
	Slug           string    `json:"slug"`
}

type Thread struct {
	Id          uint64    `json:"id"`
	Title       string    `json:"title"`
	Author      string    `json:"author"`
	Forum       string    `json:"forum"`
	Message     string    `json:"message"`
	Votes       string    `json:"votes"`
	Slug        string    `json:"slug"`
	DateCreated time.Time `json:"created"`
}

type ThreadUpdate struct {
	Title   string `json:"title"`
	Message string `json:"message"`
}

type ThreadVote struct {
	NickName string `json:"nickname"`
	Voice    int    `json:"voice"`
}

type ThreadPaginator struct {
	Limit     uint64    `json:"limit"`
	Since     time.Time `json:"since"`
	SortOrder bool      `json:"desc"`
}
