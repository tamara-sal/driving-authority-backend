package notifications

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Repo struct {
	coll *mongo.Collection
}

func NewRepo(db *mongo.Database) *Repo {
	return &Repo{coll: db.Collection("notifications")}
}

func (r *Repo) Insert(ctx context.Context, n Notification) (Notification, error) {
	now := time.Now()
	n.CreatedAt = now
	res, err := r.coll.InsertOne(ctx, n)
	if err != nil {
		return Notification{}, err
	}
	if oid, ok := res.InsertedID.(primitive.ObjectID); ok {
		n.ID = oid
	}
	return n, nil
}

func (r *Repo) ListByUser(ctx context.Context, userID primitive.ObjectID) ([]Notification, error) {
	cur, err := r.coll.Find(ctx, bson.M{"user_id": userID}, options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}}))
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)
	var out []Notification
	if err := cur.All(ctx, &out); err != nil {
		return nil, err
	}
	if out == nil {
		out = []Notification{}
	}
	return out, nil
}

func (r *Repo) MarkRead(ctx context.Context, id, userID primitive.ObjectID) (Notification, error) {
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	var out Notification
	err := r.coll.FindOneAndUpdate(ctx, bson.M{"_id": id, "user_id": userID}, bson.M{
		"$set": bson.M{"status": "read"},
	}, opts).Decode(&out)
	return out, err
}
