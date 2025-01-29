package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/yonchando/chirpy/internal/database"
)

type Chirp struct {
	ID        uuid.UUID `json:"id"`
	Body      string    `json:"body"`
	UserID    uuid.UUID `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (c *Chirp) Map(chirp database.Chirp) {
	c.ID = chirp.ID
	c.Body = chirp.Body
	c.UserID = chirp.UserID
	c.CreatedAt = chirp.CreatedAt
	c.UpdatedAt = chirp.UpdatedAt
}
