package views

import (
	"fmt"
	"testing"
)

func TestNewUser(t *testing.T) {

	tests := []struct {
		tName    string
		username string
		email    string
		password string
		isErr    bool
	}{
		{
			tName:    "empty",
			username: "",
			email:    "",
			password: "",
			isErr:    true,
		},
		{
			tName:    "ok",
			username: "okfine",
			email:    "email@google.com",
			password: "abcdefG1'",
			isErr:    false,
		},
		{
			tName:    "username too short",
			username: "okfin", // too short
			email:    "email@google.com",
			password: "abcdefG1'",
			isErr:    true,
		},
		{
			tName:    "invalid email",
			username: "okfine",
			email:    "@google.com", // invalidd
			password: "abcdefG1'",
			isErr:    true,
		},
		{
			tName:    "no special password char",
			username: "okfine",
			email:    "email@google.com",
			password: "abcdefG1", // no special
			isErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("test_%s", tt.tName), func(t *testing.T) {
			_, err := newUserValidate(map[string][]string{
				"username": []string{tt.username},
				"email":    []string{tt.email},
				"password": []string{tt.password},
			})
			if (err != nil) != tt.isErr {
				t.Errorf("got unexpected error %v", err)
				return
			}
		})

	}
}
