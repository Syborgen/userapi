package store

import "time"

type User struct {
	CreatedAt   time.Time `json:"created_at"`
	DisplayName string    `json:"display_name"`
	Email       string    `json:"email"`
}

func (u *User) Update(newUserData User) {
	if newUserData.DisplayName != "" {
		u.DisplayName = newUserData.DisplayName
	}

	if newUserData.Email != "" {
		u.Email = newUserData.Email
	}
}
