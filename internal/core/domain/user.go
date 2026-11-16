package domain

import (
	"fmt"
	"regexp"

	core_errors "github.com/FITLPI/TODOGOLANG/internal/core/errors"
)

type User struct {
	ID      int
	Version int

	FullName    string
	PhoneNumber *string
}

func NewUser(id int, version int, fullName string, phoneNumber *string) User {
	return User{
		ID:          id,
		Version:     version,
		FullName:    fullName,
		PhoneNumber: phoneNumber,
	}
}

func NewUserUninitialized(FullName string, phoneNumber *string) User {
	return NewUser(UninitializedID, UninitializedVersion, FullName, phoneNumber)
}

type UserPatch struct {
	FullName    Nullable[string]
	PhoneNumber Nullable[string]
}

// Validate проверяет сам патч: обязательные поля нельзя выставлять в null.
func (p UserPatch) Validate() error {
	if p.FullName.Set && p.FullName.Value == nil {
		return fmt.Errorf("full name cannot be null: %w", core_errors.ErrInvalidArgument)
	}

	if p.FullName.Value != nil {
		fullNameLength := len([]rune(*p.FullName.Value))
		if fullNameLength < 3 || fullNameLength > 100 {
			return fmt.Errorf("invalid full name length: %d: %w", fullNameLength, core_errors.ErrInvalidArgument)
		}
	}

	if p.PhoneNumber.Value != nil {
		phoneNumber := *p.PhoneNumber.Value

		phoneNumberLength := len([]rune(phoneNumber))
		if phoneNumberLength < 10 || phoneNumberLength > 15 {
			return fmt.Errorf("invalid phone number length: %d: %w", phoneNumberLength, core_errors.ErrInvalidArgument)
		}

		re := regexp.MustCompile(`^\+[0-9]+$`)
		if !re.MatchString(phoneNumber) {
			return fmt.Errorf("invalid phone number format: %s: %w", phoneNumber, core_errors.ErrInvalidArgument)
		}
	}

	return nil
}

// ApplyPatch — read-modify-write «всё или ничего»: изменения применяются
// к копии, копия валидируется, и только при успехе заменяется оригинал.
// Version не трогается: инкремент делает SQL в UpdateUser.
func (u *User) ApplyPatch(p UserPatch) error {
	if err := p.Validate(); err != nil {
		return fmt.Errorf("validate user patch: %w", err)
	}

	tmp := *u

	if p.FullName.Set {
		tmp.FullName = *p.FullName.Value
	}

	if p.PhoneNumber.Set {
		tmp.PhoneNumber = p.PhoneNumber.Value
	}

	if err := tmp.Validate(); err != nil {
		return fmt.Errorf("validate patched user: %w", err)
	}

	*u = tmp

	return nil
}

func (u User) Validate() error {
	fullNameLength := len([]rune(u.FullName))
	if fullNameLength < 3 || fullNameLength > 100 {
		return fmt.Errorf("invalid full name length: %d: %w", fullNameLength, core_errors.ErrInvalidArgument)
	}

	if u.PhoneNumber != nil {
		phoneNumberLength := len([]rune(*u.PhoneNumber))
		if phoneNumberLength < 10 || phoneNumberLength > 15 {
			return fmt.Errorf("invalid phone number length: %d: %w", phoneNumberLength, core_errors.ErrInvalidArgument)
		}
	}

	re := regexp.MustCompile(`^\+[0-9]+$`)

	if u.PhoneNumber != nil && !re.MatchString(*u.PhoneNumber) {
		return fmt.Errorf("invalid phone number format: %s: %w", *u.PhoneNumber, core_errors.ErrInvalidArgument)
	}

	return nil
}
