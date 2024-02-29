package domain

import "time"

type Session struct {
	UserId         int
	Username       string
	SessionId      string
	CreationDate   time.Time
	ExpiritionDate time.Time
}
