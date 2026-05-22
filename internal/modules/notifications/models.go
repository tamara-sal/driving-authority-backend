package notifications

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Notification struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID      primitive.ObjectID `bson:"user_id" json:"user_id"`
	Title       string             `bson:"title" json:"title"`
	Description string             `bson:"description" json:"description"`
	Status      string             `bson:"status" json:"status"` // unread, read
	Type        string             `bson:"type" json:"type"`     // success, info, warning, danger
	CreatedAt   time.Time          `bson:"created_at" json:"created_at"`
}

type NotificationView struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Timestamp   string `json:"timestamp"`
	Status      string `json:"status"`
	Type        string `json:"type"`
}
