package handler

import (
	"goalify/pkg/options"
	"strings"
)

type (
	SignupRequest struct {
		Email           string `json:"email"`
		Password        string `json:"password"`
		ConfirmPassword string `json:"confirm_password"`
	}
	LoginRequest struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	RefreshRequest struct {
		UserId       string `json:"user_id"`
		RefreshToken string `json:"refresh_token"`
	}
	UpdateRequest struct {
		Xp            options.Option[int] `json:"xp"`
		LevelId       options.Option[int] `json:"level_id"`
		CashAvailable options.Option[int] `json:"cash_available"`
	}
)

const (
	PASSWORD_MIN_LEN = 8
	DIGITS           = "0123456789"
	SYMBOLS          = "!@#$%^"
)

func ValidateEmail(problems map[string]string, email string) {
	if email == "" {
		problems["email"] = "email is required"
	} else if !strings.Contains(email, "@") {
		problems["email"] = "email is invalid"
	}
}

func ValidatePassword(problems map[string]string, password string) {
	if password == "" {
		problems["password"] = "password is required"
	} else if len(password) < PASSWORD_MIN_LEN {
		problems["password"] = "password must be at least 8 characters"
	} else if strings.Contains(password, " ") {
		problems["password"] = "password cannot contain spaces"
	} else if !strings.ContainsAny(password, DIGITS) {
		problems["password"] = "password must contain at least one digit"
	} else if !strings.ContainsAny(password, SYMBOLS) {
		problems["password"] = "password must contain at least one symbol"
	}
}

func (r SignupRequest) Valid() map[string]string {
	problems := make(map[string]string)
	r.Email = strings.TrimSpace(r.Email)

	ValidateEmail(problems, r.Email)
	ValidatePassword(problems, r.Password)

	if r.Password != r.ConfirmPassword {
		problems["confirm_password"] = "passwords do not match"
	}

	return problems
}

func (r LoginRequest) Valid() map[string]string {
	problems := make(map[string]string)

	ValidateEmail(problems, r.Email)
	ValidatePassword(problems, r.Password)

	return problems
}

func checkNonNegativeIntField(problems map[string]string, fieldName string, val options.Option[int]) {
	if val.IsPresent() && val.ValueOrZero() < 0 {
		problems[fieldName] = fieldName + " must be non-negative"
	}
}

func (r UpdateRequest) Valid() map[string]string {
	problems := make(map[string]string)
	checkNonNegativeIntField(problems, "xp", r.Xp)
	checkNonNegativeIntField(problems, "level_id", r.LevelId)
	checkNonNegativeIntField(problems, "cash_available", r.CashAvailable)

	return problems
}

func (r RefreshRequest) Valid() map[string]string {
	problems := make(map[string]string)
	if r.RefreshToken == "" {
		problems["refresh_token"] = "refresh token is required"
	}
	return problems
}
