package monitoring

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Device struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	VehicleID    primitive.ObjectID `bson:"vehicle_id" json:"vehicle_id"`
	DeviceSerial string             `bson:"device_serial" json:"device_serial"`
	Status       string             `bson:"status" json:"status"`
}

type Trip struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	VehicleID    primitive.ObjectID `bson:"vehicle_id" json:"vehicle_id"`
	UserID       primitive.ObjectID `bson:"user_id,omitempty" json:"user_id,omitempty"`
	StartTime    time.Time          `bson:"start_time" json:"start_time"`
	EndTime      time.Time          `bson:"end_time" json:"end_time"`
	Distance     float64            `bson:"distance" json:"distance"`
	AverageSpeed float64            `bson:"average_speed" json:"average_speed"`
	SafetyScore  float64            `bson:"safety_score" json:"safety_score"`
}

type TripEvent struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	TripID    primitive.ObjectID `bson:"trip_id" json:"trip_id"`
	EventType string             `bson:"event_type" json:"event_type"`
	Severity  string             `bson:"severity" json:"severity"`
	Timestamp time.Time          `bson:"timestamp" json:"timestamp"`
}

type DeviceDataInput struct {
	DeviceSerial string       `json:"device_serial" binding:"required"`
	VehicleID    string       `json:"vehicle_id" binding:"required"`
	UserID       string       `json:"user_id"`
	Trip         TripPayload  `json:"trip" binding:"required"`
	Events       []EventInput `json:"events"`
}

type TripPayload struct {
	StartTime    time.Time `json:"start_time" binding:"required"`
	EndTime      time.Time `json:"end_time" binding:"required"`
	Distance     float64   `json:"distance"`
	AverageSpeed float64   `json:"average_speed"`
	SafetyScore  float64   `json:"safety_score"`
}

type EventInput struct {
	EventType string    `json:"event_type" binding:"required"`
	Severity  string    `json:"severity"`
	Timestamp time.Time `json:"timestamp"`
}

type ScoreOutput struct {
	UserID       string  `json:"user_id"`
	AverageScore float64 `json:"average_score"`
	TripCount    int64   `json:"trip_count"`
}
