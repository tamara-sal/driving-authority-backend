package payments

import (
	"context"
	"time"

	"driving-authority-backend/internal/domain"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Repo struct {
	fees     *mongo.Collection
	payments *mongo.Collection
}

func NewRepo(db *mongo.Database) *Repo {
	return &Repo{
		fees:     db.Collection("service_fees"),
		payments: db.Collection("payments"),
	}
}

func (r *Repo) CountFees(ctx context.Context) (int64, error) {
	return r.fees.CountDocuments(ctx, bson.M{})
}

func (r *Repo) InsertFees(ctx context.Context, fees []ServiceFee) error {
	docs := make([]any, len(fees))
	for i := range fees {
		docs[i] = fees[i]
	}
	_, err := r.fees.InsertMany(ctx, docs)
	return err
}

func (r *Repo) FindFeeByType(ctx context.Context, st domain.ServiceType) (ServiceFee, error) {
	var f ServiceFee
	err := r.fees.FindOne(ctx, bson.M{"service_type": st}).Decode(&f)
	return f, err
}

func (r *Repo) InsertPayment(ctx context.Context, p Payment) (Payment, error) {
	p.CreatedAt = time.Now()
	res, err := r.payments.InsertOne(ctx, p)
	if err != nil {
		return Payment{}, err
	}
	if oid, ok := res.InsertedID.(primitive.ObjectID); ok {
		p.ID = oid
	}
	return p, nil
}

func (r *Repo) MarkPaid(ctx context.Context, id primitive.ObjectID) (Payment, error) {
	now := time.Now()
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	var out Payment
	err := r.payments.FindOneAndUpdate(ctx, bson.M{"_id": id}, bson.M{
		"$set": bson.M{"status": domain.PaymentPaid, "paid_at": now},
	}, opts).Decode(&out)
	return out, err
}

func (r *Repo) HistoryByUser(ctx context.Context, userID primitive.ObjectID) ([]Payment, error) {
	cur, err := r.payments.Find(ctx, bson.M{"user_id": userID},
		options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}}))
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)
	var out []Payment
	if err := cur.All(ctx, &out); err != nil {
		return nil, err
	}
	if out == nil {
		out = []Payment{}
	}
	return out, nil
}

func (r *Repo) SumPaid(ctx context.Context) (float64, error) {
	pipe := mongo.Pipeline{
		{{Key: "$match", Value: bson.M{"status": domain.PaymentPaid}}},
		{{Key: "$group", Value: bson.M{"_id": nil, "total": bson.M{"$sum": "$amount"}}}},
	}
	cur, err := r.payments.Aggregate(ctx, pipe)
	if err != nil {
		return 0, err
	}
	defer cur.Close(ctx)
	if !cur.Next(ctx) {
		return 0, nil
	}
	var doc struct {
		Total float64 `bson:"total"`
	}
	if err := cur.Decode(&doc); err != nil {
		return 0, err
	}
	return doc.Total, nil
}

func (r *Repo) CountPayments(ctx context.Context, filter bson.M) (int64, error) {
	return r.payments.CountDocuments(ctx, filter)
}
