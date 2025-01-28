package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/yonchando/chirpy/internal/database"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (u *User) Map(user database.User) *User {

	u.ID = user.ID
	u.Email = user.Email
	u.CreatedAt = user.CreatedAt
	u.UpdatedAt = user.UpdatedAt

	return u
}
