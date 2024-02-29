package domain

import "time"

type Comments struct {
	CommentId    int
	PostId       int
	UserId       int
	Username     string
	Content      string
	Likes        int
	Dislikes     int
	CreationDate time.Time
}
