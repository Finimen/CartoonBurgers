package models

import "errors"

type User struct {
	Id       int    `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Bonus    int    `json:"bonus"`
}

func (u *User) Validate() error {
	if u.Username == "" {
		return errors.New("username cannot be empty")
	}
	if len(u.Password) < 6 {
		return errors.New("password must be at least 6 characters")
	}
	return nil
}
