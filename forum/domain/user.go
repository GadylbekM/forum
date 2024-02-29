package domain

import "time"

type User struct {
	Username         string
	UserId           int
	Password         string
	Email            string
	RegistrationDate time.Time
}
