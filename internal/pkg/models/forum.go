package models

type Forum struct {
	Title string `json:"title"`
	AuthorNickName string `json:"user"`
	Slug string `json:"slug"`
	Posts uint64 `json:"posts"`
	Threads uint64 `json:"threads"`
}

type ForumCreate struct {
	Title string `json:"title"`
	AuthorNickName string `json:"user"`
	Slug string `json:"slug"`
}
