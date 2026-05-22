package violations

import (
	"time"

	"driving-authority-backend/internal/domain"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Violation struct {
	ID        primitive.ObjectID    `bson:"_id,omitempty" json:"id"`
	DriverID  primitive.ObjectID    `bson:"driver_id" json:"driver_id"`
	Driver    string                `bson:"driver_name" json:"driver"`
	Type      string                `bson:"type" json:"type"`
	Severity  string                `bson:"severity" json:"severity"`
	Status    domain.ViolationStatus `bson:"status" json:"status"`
	CreatedAt time.Time             `bson:"created_at" json:"created_at"`
}

type ViolationView struct {
	ID       string `json:"id"`
	Driver   string `json:"driver"`
	Type     string `json:"type"`
	Severity string `json:"severity"`
	Date     string `json:"date"`
	Status   string `json:"status"`
}

type CreateInput struct {
	DriverID string `json:"driver_id"`
	Driver   string `json:"driver"`
	Type     string `json:"type" binding:"required"`
	Severity string `json:"severity" binding:"required"`
}

type UpdateStatusInput struct {
	Status domain.ViolationStatus `json:"status" binding:"required"`
}
