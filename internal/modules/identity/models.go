package identity

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type VerificationStatus string

const (
	StatusPending  VerificationStatus = "pending"
	StatusApproved VerificationStatus = "approved"
	StatusRejected VerificationStatus = "rejected"
)

type IdentityVerification struct {
	ID                primitive.ObjectID  `bson:"_id,omitempty" json:"id"`
	UserID            primitive.ObjectID  `bson:"user_id" json:"user_id"`
	NationalIDNumber  string              `bson:"national_id_number" json:"national_id_number"`
	DocumentFrontPath string              `bson:"document_front_path" json:"document_front_path"`
	DocumentBackPath  string              `bson:"document_back_path" json:"document_back_path"`
	SelfiePath        string              `bson:"selfie_path" json:"selfie_path"`
	Status            VerificationStatus  `bson:"status" json:"status"`
	ReviewedBy        *primitive.ObjectID `bson:"reviewed_by,omitempty" json:"reviewed_by,omitempty"`
	ReviewComment     string              `bson:"review_comment,omitempty" json:"review_comment,omitempty"`
	SubmittedAt       time.Time           `bson:"submitted_at" json:"submitted_at"`
	ReviewedAt        *time.Time          `bson:"reviewed_at,omitempty" json:"reviewed_at,omitempty"`
}
