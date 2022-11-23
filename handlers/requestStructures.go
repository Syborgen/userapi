package handlers

import (
	"errors"
	"net/http"
	"strings"
)

type CreateUserRequest struct {
	DisplayName string `json:"display_name"`
	Email       string `json:"email"`
}

func (c *CreateUserRequest) Bind(r *http.Request) error {
	if c.DisplayName == "" {
		return errors.New("display_name cannot be empty")
	}

	if !strings.Contains(c.Email, "@") {
		return errors.New("email must conatin symbol @")
	}

	return nil
}

type UpdateUserRequest struct {
	DisplayName string `json:"display_name"`
	Email       string `json:"email"`
}

func (c *UpdateUserRequest) Bind(r *http.Request) error {
	if c.DisplayName == "" && c.Email == "" {
		return errors.New("display_name or email must be set")
	}

	if !strings.Contains(c.Email, "@") {
		return errors.New("email must conatin symbol @")
	}

	return nil
}
