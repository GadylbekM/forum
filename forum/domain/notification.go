package domain

import "time"

type Notification struct {
	Id        int
	UserId    int
	Type      string
	PostId    int
	OwnerId   int
	Username  string
	Timestamp time.Time
}

type Notification_comments struct {
	Id        int
	UserId    int
	Type      string
	CommentId int
	OwnerId   int
	Username  string
	PostId    int
	Timestamp time.Time
}
