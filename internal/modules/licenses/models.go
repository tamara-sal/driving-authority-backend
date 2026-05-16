package licenses

import (
	"time"

	"driving-authority-backend/internal/domain"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type License struct {
	ID            primitive.ObjectID  `bson:"_id,omitempty" json:"id"`
	UserID        primitive.ObjectID  `bson:"user_id" json:"user_id"`
	LicenseNumber string              `bson:"license_number,omitempty" json:"license_number,omitempty"`
	Type          domain.LicenseType  `bson:"type" json:"type"`
	Status        domain.LicenseStatus `bson:"status" json:"status"`
	IssueDate     *time.Time          `bson:"issue_date,omitempty" json:"issue_date,omitempty"`
	ExpiryDate    *time.Time          `bson:"expiry_date,omitempty" json:"expiry_date,omitempty"`
	CreatedAt     time.Time           `bson:"created_at" json:"created_at"`
	UpdatedAt     time.Time           `bson:"updated_at" json:"updated_at"`
}
