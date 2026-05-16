package audit

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Log struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID     primitive.ObjectID `bson:"user_id" json:"user_id"`
	ActionType string             `bson:"action_type" json:"action_type"`
	Metadata   map[string]any     `bson:"metadata,omitempty" json:"metadata,omitempty"`
	IPAddress  string             `bson:"ip_address,omitempty" json:"ip_address,omitempty"`
	CreatedAt  time.Time          `bson:"created_at" json:"created_at"`
}

type Logger struct {
	coll *mongo.Collection
}

func NewLogger(db *mongo.Database) *Logger {
	return &Logger{coll: db.Collection("account_activity_logs")}
}

func (l *Logger) Record(ctx context.Context, userID primitive.ObjectID, action, ip string, metadata map[string]any) {
	if metadata == nil {
		metadata = map[string]any{}
	}
	_, _ = l.coll.InsertOne(ctx, Log{
		UserID:     userID,
		ActionType: action,
		Metadata:   metadata,
		IPAddress:  ip,
		CreatedAt:  time.Now(),
	})
}
