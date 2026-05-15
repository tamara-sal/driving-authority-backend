package auth

import (
	"time"

	"driving-authority-backend/internal/domain"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID            primitive.ObjectID   `bson:"_id,omitempty" json:"id"`
	UUID          string               `bson:"uuid" json:"uuid"`
	FirstName     string               `bson:"first_name" json:"first_name"`
	LastName      string               `bson:"last_name" json:"last_name"`
	Email         string               `bson:"email" json:"email"`
	PasswordHash  string               `bson:"password_hash" json:"-"`
	Phone         string               `bson:"phone,omitempty" json:"phone,omitempty"`
	Role          domain.Role          `bson:"role" json:"role"`
	Status        domain.AccountStatus `bson:"status" json:"status"`
	EmailVerified bool                 `bson:"email_verified" json:"email_verified"`
	CreatedAt     time.Time            `bson:"created_at" json:"created_at"`
	UpdatedAt     time.Time            `bson:"updated_at" json:"updated_at"`
}
