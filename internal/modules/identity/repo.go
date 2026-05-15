package identity

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
	c := db.Collection("identity_verifications")

	_, _ = c.Indexes().CreateMany(context.Background(), []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "user_id", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys: bson.D{{Key: "status", Value: 1}, {Key: "submitted_at", Value: -1}},
		},
	})

	return &Repo{coll: c}
}

func (r *Repo) UpsertSubmit(ctx context.Context, v IdentityVerification) (IdentityVerification, error) {
	v.Status = StatusPending
	v.SubmittedAt = time.Now()
	v.ReviewedAt = nil
	v.ReviewedBy = nil
	v.ReviewComment = ""

	filter := bson.M{"user_id": v.UserID}
	update := bson.M{"$set": v}
	opts := options.FindOneAndUpdate().SetUpsert(true).SetReturnDocument(options.After)

	var out IdentityVerification
	err := r.coll.FindOneAndUpdate(ctx, filter, update, opts).Decode(&out)
	return out, err
}

func (r *Repo) FindByUserID(ctx context.Context, userID primitive.ObjectID) (IdentityVerification, error) {
	var out IdentityVerification
	err := r.coll.FindOne(ctx, bson.M{"user_id": userID}).Decode(&out)
	return out, err
}

func (r *Repo) FindByID(ctx context.Context, id primitive.ObjectID) (IdentityVerification, error) {
	var out IdentityVerification
	err := r.coll.FindOne(ctx, bson.M{"_id": id}).Decode(&out)
	return out, err
}

func (r *Repo) SetDecision(ctx context.Context, id primitive.ObjectID, reviewer primitive.ObjectID, status VerificationStatus, comment string) (IdentityVerification, error) {
	now := time.Now()
	update := bson.M{
		"$set": bson.M{
			"status":         status,
			"review_comment": comment,
			"reviewed_by":    reviewer,
			"reviewed_at":    now,
		},
	}
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)

	var out IdentityVerification
	err := r.coll.FindOneAndUpdate(ctx, bson.M{"_id": id}, update, opts).Decode(&out)
	return out, err
}
