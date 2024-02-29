package domain

import "time"

type Posts struct {
	PostId       int
	UserId       int
	Username     string
	Category     string
	Title        string
	ImageField   string
	Content      string
	CategoryId   int
	Likes        int
	Dislikes     int
	Comments     []Comments
	CreationDate time.Time
}
