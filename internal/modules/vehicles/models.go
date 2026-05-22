package vehicles

import (
	"time"

	"driving-authority-backend/internal/domain"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Vehicle struct {
	ID          primitive.ObjectID  `bson:"_id,omitempty" json:"id"`
	VIN         string              `bson:"vin" json:"vin"`
	PlateNumber string              `bson:"plate_number" json:"plate_number"`
	Make        string              `bson:"make" json:"make"`
	Model       string              `bson:"model" json:"model"`
	Year        int                  `bson:"year" json:"year"`
	Color       string               `bson:"color,omitempty" json:"color,omitempty"`
	Status      domain.VehicleStatus `bson:"status" json:"status"`
	OwnerID     primitive.ObjectID  `bson:"owner_id" json:"owner_id"`
	CreatedAt   time.Time           `bson:"created_at" json:"created_at"`
}

type TransferRequest struct {
	ID         primitive.ObjectID     `bson:"_id,omitempty" json:"id"`
	VehicleID  primitive.ObjectID     `bson:"vehicle_id" json:"vehicle_id"`
	SellerID   primitive.ObjectID     `bson:"seller_id" json:"seller_id"`
	BuyerID    primitive.ObjectID     `bson:"buyer_id" json:"buyer_id"`
	Status     domain.TransferStatus  `bson:"status" json:"status"`
	ReviewedBy *primitive.ObjectID    `bson:"reviewed_by,omitempty" json:"reviewed_by,omitempty"`
	CreatedAt  time.Time              `bson:"created_at" json:"created_at"`
	ReviewedAt *time.Time             `bson:"reviewed_at,omitempty" json:"reviewed_at,omitempty"`
}
