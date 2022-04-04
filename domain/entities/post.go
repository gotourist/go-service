package entities

type Post struct {
	Id        int64  `db:"id" json:"id"`
	Title     string `db:"title" json:"title"`
	Body      string `db:"body" json:"body"`
	IsDeleted bool   `db:"is_deleted" json:"is_deleted"`
	CreatedAt string `db:"created_at" json:"created_at"`
}

type PostList struct {
	Data []Post `json:"data"`
}
