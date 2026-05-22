package licenses

import (
	"time"

	"driving-authority-backend/internal/domain"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ApplicationDetails struct {
	Name        string `bson:"name,omitempty" json:"name,omitempty"`
	DOB         string `bson:"dob,omitempty" json:"dob,omitempty"`
	Gender      string `bson:"gender,omitempty" json:"gender,omitempty"`
	Nationality string `bson:"nationality,omitempty" json:"nationality,omitempty"`
	Address     string `bson:"address,omitempty" json:"address,omitempty"`
	City        string `bson:"city,omitempty" json:"city,omitempty"`
	Postal      string `bson:"postal,omitempty" json:"postal,omitempty"`
}

type License struct {
	ID            primitive.ObjectID   `bson:"_id,omitempty" json:"id"`
	UserID        primitive.ObjectID   `bson:"user_id" json:"user_id"`
	ReferenceID   string               `bson:"reference_id,omitempty" json:"reference_id,omitempty"`
	LicenseNumber string               `bson:"license_number,omitempty" json:"license_number,omitempty"`
	Type          domain.LicenseType   `bson:"type" json:"type"`
	Status        domain.LicenseStatus `bson:"status" json:"status"`
	Application   *ApplicationDetails  `bson:"application,omitempty" json:"application,omitempty"`
	IssueDate     *time.Time           `bson:"issue_date,omitempty" json:"issue_date,omitempty"`
	ExpiryDate    *time.Time           `bson:"expiry_date,omitempty" json:"expiry_date,omitempty"`
	CreatedAt     time.Time            `bson:"created_at" json:"created_at"`
	UpdatedAt     time.Time            `bson:"updated_at" json:"updated_at"`
}
