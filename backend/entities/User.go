package entities

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	RefreshTokenExpiry time.Time `db:"refresh_token_expiry" json:"refresh_token_expiry"`
	Email              string    `db:"email" json:"email"`
	Password           string    `db:"password"`
	Xp                 int       `db:"xp" json:"xp"`
	LevelId            int       `db:"level_id" json:"level_id"`
	CashAvailable      int       `db:"cash_available" json:"cash_available"`
	Id                 uuid.UUID `db:"id" json:"id"`
	RefreshToken       uuid.UUID `db:"refresh_token" json:"refresh_token"`
}

type UserDTO struct {
	RefreshTokenExpiry time.Time `db:"refresh_token_expiry" json:"refresh_token_expiry"`
	Email              string    `db:"email" json:"email"`
	AccessToken        string    `json:"access_token"`
	Xp                 int       `db:"xp" json:"xp"`
	LevelId            int       `db:"level_id" json:"level_id"`
	CashAvailable      int       `db:"cash_available" json:"cash_available"`
	Id                 uuid.UUID `db:"id" json:"id"`
	RefreshToken       uuid.UUID `db:"refresh_token" json:"refresh_token"`
}
