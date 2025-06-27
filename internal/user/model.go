package user

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID  `json:"id" db:"id"`
	Email        string     `json:"email" db:"email"`
	Name         string     `json:"name" db:"name"`
	PasswordHash string     `json:"-" db:"password_hash"` // Hidden from JSON responses
	Role         string     `json:"role" db:"role"`
	StartupID    *uuid.UUID `json:"startup_id,omitempty" db:"startup_id"`
	CreatedAt    time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at" db:"updated_at"`
}

func (User) TableName() string {
	return "users"
}
