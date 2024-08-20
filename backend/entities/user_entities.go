package entities

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	RefreshTokenExpiry time.Time `db:"refresh_token_expiry" json:"refresh_token_expiry"`
	CreatedAt          time.Time `db:"created_at" json:"created_at"`
	UpdatedAt          time.Time `db:"updated_at" json:"updated_at"`
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
	CreatedAt          time.Time `db:"created_at" json:"created_at"`
	UpdatedAt          time.Time `db:"updated_at" json:"updated_at"`
	Email              string    `db:"email" json:"email"`
	AccessToken        string    `json:"access_token"`
	Xp                 int       `db:"xp" json:"xp"`
	LevelId            int       `db:"level_id" json:"level_id"`
	CashAvailable      int       `db:"cash_available" json:"cash_available"`
	Id                 uuid.UUID `db:"id" json:"id"`
	RefreshToken       uuid.UUID `db:"refresh_token" json:"refresh_token"`
}

func (u *User) ToUserDTO(accessToken string) *UserDTO {
	return &UserDTO{
		RefreshTokenExpiry: u.RefreshTokenExpiry,
		CreatedAt:          u.CreatedAt,
		UpdatedAt:          u.UpdatedAt,
		Email:              u.Email,
		AccessToken:        accessToken,
		Xp:                 u.Xp,
		LevelId:            u.LevelId,
		CashAvailable:      u.CashAvailable,
		Id:                 u.Id,
		RefreshToken:       u.RefreshToken,
	}
}
