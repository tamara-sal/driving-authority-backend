package inspections

import (
	"time"

	"driving-authority-backend/internal/domain"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type VehicleInspection struct {
	ID              primitive.ObjectID      `bson:"_id,omitempty" json:"id"`
	VehicleID       primitive.ObjectID      `bson:"vehicle_id" json:"vehicle_id"`
	RequestedBy     primitive.ObjectID      `bson:"requested_by" json:"requested_by"`
	InspectionDate  time.Time               `bson:"inspection_date" json:"inspection_date"`
	ExpiryDate      *time.Time              `bson:"expiry_date,omitempty" json:"expiry_date,omitempty"`
	Status          domain.InspectionStatus `bson:"status" json:"status"`
	ReportPath      string                  `bson:"report_path,omitempty" json:"report_path,omitempty"`
	InspectorID     *primitive.ObjectID     `bson:"inspector_id,omitempty" json:"inspector_id,omitempty"`
	CreatedAt       time.Time               `bson:"created_at" json:"created_at"`
}
