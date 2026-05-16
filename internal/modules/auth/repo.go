package auth

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type UserRepo struct {
	coll *mongo.Collection
}

func NewUserRepo(db *mongo.Database) *UserRepo {
	c := db.Collection("users")

	// Fire-and-forget index creation (safe to call on startup).
	_, _ = c.Indexes().CreateMany(context.Background(), []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "email", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys:    bson.D{{Key: "uuid", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
	})

	return &UserRepo{coll: c}
}

func (r *UserRepo) Insert(ctx context.Context, u User) (User, error) {
	now := time.Now()
	u.CreatedAt = now
	u.UpdatedAt = now
	res, err := r.coll.InsertOne(ctx, u)
	if err != nil {
		return User{}, err
	}
	if oid, ok := res.InsertedID.(primitive.ObjectID); ok {
		u.ID = oid
	}
	return u, nil
}

func (r *UserRepo) FindByEmail(ctx context.Context, email string) (User, error) {
	var u User
	err := r.coll.FindOne(ctx, bson.M{"email": email}).Decode(&u)
	return u, err
}

func (r *UserRepo) SetEmailVerified(ctx context.Context, id primitive.ObjectID, verified bool) error {
	_, err := r.coll.UpdateOne(ctx, bson.M{"_id": id}, bson.M{
		"$set": bson.M{"email_verified": verified, "updated_at": time.Now()},
	})
	return err
}

func (r *UserRepo) UpdatePasswordHash(ctx context.Context, id primitive.ObjectID, hash string) error {
	_, err := r.coll.UpdateOne(ctx, bson.M{"_id": id}, bson.M{
		"$set": bson.M{"password_hash": hash, "updated_at": time.Now()},
	})
	return err
}

func (r *UserRepo) Count(ctx context.Context) (int64, error) {
	return r.coll.CountDocuments(ctx, bson.M{})
}
