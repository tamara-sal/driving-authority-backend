package violations

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
	return &Repo{coll: db.Collection("violations")}
}

func (r *Repo) Insert(ctx context.Context, v Violation) (Violation, error) {
	now := time.Now()
	v.CreatedAt = now
	res, err := r.coll.InsertOne(ctx, v)
	if err != nil {
		return Violation{}, err
	}
	if oid, ok := res.InsertedID.(primitive.ObjectID); ok {
		v.ID = oid
	}
	return v, nil
}

func (r *Repo) List(ctx context.Context, filter bson.M) ([]Violation, error) {
	cur, err := r.coll.Find(ctx, filter, options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}}))
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)
	var out []Violation
	if err := cur.All(ctx, &out); err != nil {
		return nil, err
	}
	if out == nil {
		out = []Violation{}
	}
	return out, nil
}

func (r *Repo) UpdateStatus(ctx context.Context, id primitive.ObjectID, status string) (Violation, error) {
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	var out Violation
	err := r.coll.FindOneAndUpdate(ctx, bson.M{"_id": id}, bson.M{"$set": bson.M{"status": status}}, opts).Decode(&out)
	return out, err
}
