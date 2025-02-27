package views

import (
	"errors"
	"net/mail"
	"strings"

	"github.com/gorilla/schema"
)

// decoder caches meta-data and can be shared safely.
var userDecoder = schema.NewDecoder()

// newUser is a user provided over the web
type newUser struct {
	Username string `schema:"username,required"`
	Email    string `schema:"email,required"`
	Password string `schema:"password,required"`
}

// newUserValidate valides a newUser
func newUserValidate(input map[string][]string) (*newUser, error) {
	var user newUser
	err := userDecoder.Decode(&user, input)
	if err != nil {
		return nil, err
	}
	if len(user.Username) < 6 || len(user.Username) > 16 {
		return nil, errors.New("username must be between 6 and 16 characters in length")
	}
	if _, err := mail.ParseAddress(user.Email); err != nil {
		return nil, errors.New("email not valid")
	}
	if len(user.Password) < 8 || len(user.Password) > 32 {
		return nil, errors.New("password should be between 8 and 32 characters in length")
	}
	alpha, Alpha, num, special := 0, 0, 0, 0
	// owasp specials https://owasp.org/www-community/password-special-characters
	specials := `!"#$%&'()*+,-./:;<=>?@[\]^_{|}~` + "`"
	for i := 0; i < len(user.Password); i++ {
		switch up := string(user.Password[i]); { // true switch
		case strings.ContainsAny(up, "abcdefghijklmnopqrstuvwxyz"):
			alpha++
		case strings.ContainsAny(up, "ABCDEFGHIJKLMNOPQRSTUVWXYZ"):
			Alpha++
		case strings.ContainsAny(up, "0123456789"):
			num++
		case strings.ContainsAny(up, specials):
			special++
		}
	}
	if !(alpha > 0 && Alpha > 0 && num > 0 && special > 0) {
		return nil, errors.New("password requires numbers, upper, lowercase and special characters")
	}
	return &user, nil
}
