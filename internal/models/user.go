package models

import (
	"errors"
	"fmt"
	"time"

	"github.com/asaskevich/govalidator"
)

var (
	// ErrInvalidEmail defines the error when an invalid email address is encountered
	ErrInvalidEmail = errors.New("invalid email")
)

// User describes the structure of the user JSON data
type User struct {
	Name       string    `json:"name"`
	Email      string    `json:"email"`
	Alias      string    `json:"alias"`
	GpgKeyID   string    `json:"gpg_key_id"`
	AddedAt    time.Time `json:"added_at"`
	ModifiedAt time.Time `json:"modified_at"`
}

// NewUser returns a new user
func NewUser(name, email, alias, gpgKeyID string) *User {
	return &User{
		Name:     name,
		Email:    email,
		Alias:    alias,
		GpgKeyID: gpgKeyID,
	}
}

// Modify updates fields if there are changes. It also updated the 'ModifiedAt' field accordingly
func (u *User) Modify(name, email, alias, gpgKeyID string) {
	var modified = 0
	if name != "" {
		u.Name = name
		modified++
	}

	if email != "" {
		u.Email = email
		modified++
	}

	if alias != "" {
		u.Alias = alias
		modified++
	}

	if gpgKeyID != "" {
		u.GpgKeyID = gpgKeyID
		modified++
	}

	if modified > 0 {
		u.ModifiedAt = time.Now()
	}
}

// Format formats user profile data as a string
func (u *User) Format(padding int) string {
	if u.Alias != "" {
		return fmt.Sprintf("%-*s %s <%s>", padding, fmt.Sprintf("[%s]", u.Alias), u.Name, u.Email)
	}
	return fmt.Sprintf("%-*s %s <%s>", padding, "", u.Name, u.Email)
}

// ValidateEmail validates the provided email address. If is was modified this function also accepts empty addresses
// (= no change)
func ValidateEmail(email string, modified bool) error {
	if modified && email == "" {
		return nil
	}

	if !govalidator.IsExistingEmail(email) {
		return ErrInvalidEmail
	}
	return nil
}
