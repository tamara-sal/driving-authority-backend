package payments

import (
	"time"

	"driving-authority-backend/internal/domain"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ServiceFee struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	ServiceType domain.ServiceType `bson:"service_type" json:"service_type"`
	Amount      float64            `bson:"amount" json:"amount"`
	Currency    string             `bson:"currency" json:"currency"`
}

type Payment struct {
	ID            primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID        primitive.ObjectID `bson:"user_id" json:"user_id"`
	ServiceType   domain.ServiceType `bson:"service_type" json:"service_type"`
	Amount        float64            `bson:"amount" json:"amount"`
	Currency      string             `bson:"currency" json:"currency"`
	Status        domain.PaymentStatus `bson:"status" json:"status"`
	TransactionID string             `bson:"transaction_id" json:"transaction_id"`
	CreatedAt     time.Time          `bson:"created_at" json:"created_at"`
	PaidAt        *time.Time         `bson:"paid_at,omitempty" json:"paid_at,omitempty"`
}
