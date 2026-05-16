package practical

import (
	"time"

	"driving-authority-backend/internal/domain"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TestCenter struct {
	ID               primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name             string             `bson:"name" json:"name"`
	Location         string             `bson:"location" json:"location"`
	CapacityPerSlot  int                `bson:"capacity_per_slot" json:"capacity_per_slot"`
	CreatedAt        time.Time          `bson:"created_at" json:"created_at"`
}

type TimeSlot struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	CenterID  primitive.ObjectID `bson:"center_id" json:"center_id"`
	Date      string             `bson:"date" json:"date"`
	StartTime string             `bson:"start_time" json:"start_time"`
	EndTime   string             `bson:"end_time" json:"end_time"`
	Capacity  int                `bson:"capacity" json:"capacity"`
	Booked    int                `bson:"booked" json:"booked"`
}

type PracticalBooking struct {
	ID         primitive.ObjectID       `bson:"_id,omitempty" json:"id"`
	UserID     primitive.ObjectID       `bson:"user_id" json:"user_id"`
	SlotID     primitive.ObjectID       `bson:"slot_id" json:"slot_id"`
	CenterID   primitive.ObjectID       `bson:"center_id" json:"center_id"`
	Status     domain.BookingStatus     `bson:"status" json:"status"`
	Result     *domain.PracticalResult  `bson:"result,omitempty" json:"result,omitempty"`
	ExaminerID *primitive.ObjectID      `bson:"examiner_id,omitempty" json:"examiner_id,omitempty"`
	Comments   string                   `bson:"comments,omitempty" json:"comments,omitempty"`
	CreatedAt  time.Time                `bson:"created_at" json:"created_at"`
}
