package licenses

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
	coll *mongo.Collection
}

func NewRepo(db *mongo.Database) *Repo {
	c := db.Collection("licenses")
	_, _ = c.Indexes().CreateMany(context.Background(), []mongo.IndexModel{
		{Keys: bson.D{{Key: "user_id", Value: 1}}},
		{Keys: bson.D{{Key: "license_number", Value: 1}}, Options: options.Index().SetUnique(true).SetSparse(true)},
	})
	return &Repo{coll: c}
}

func (r *Repo) Insert(ctx context.Context, lic License) (License, error) {
	now := time.Now()
	lic.CreatedAt = now
	lic.UpdatedAt = now
	res, err := r.coll.InsertOne(ctx, lic)
	if err != nil {
		return License{}, err
	}
	if oid, ok := res.InsertedID.(primitive.ObjectID); ok {
		lic.ID = oid
	}
	return lic, nil
}

func (r *Repo) FindByUserID(ctx context.Context, userID primitive.ObjectID) ([]License, error) {
	cur, err := r.coll.Find(ctx, bson.M{"user_id": userID}, options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}}))
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)
	var out []License
	if err := cur.All(ctx, &out); err != nil {
		return nil, err
	}
	if out == nil {
		out = []License{}
	}
	return out, nil
}

func (r *Repo) FindByID(ctx context.Context, id primitive.ObjectID) (License, error) {
	var lic License
	err := r.coll.FindOne(ctx, bson.M{"_id": id}).Decode(&lic)
	return lic, err
}

func (r *Repo) Approve(ctx context.Context, id primitive.ObjectID, licenseNumber string, issueDate, expiryDate time.Time) (License, error) {
	now := time.Now()
	update := bson.M{
		"$set": bson.M{
			"status":         domain.LicenseIssued,
			"license_number": licenseNumber,
			"issue_date":     issueDate,
			"expiry_date":    expiryDate,
			"updated_at":     now,
		},
	}
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	var out License
	err := r.coll.FindOneAndUpdate(ctx, bson.M{"_id": id}, update, opts).Decode(&out)
	return out, err
}

func (r *Repo) Renew(ctx context.Context, id, userID primitive.ObjectID, expiryDate time.Time) (License, error) {
	now := time.Now()
	update := bson.M{
		"$set": bson.M{
			"status":      domain.LicenseIssued,
			"expiry_date": expiryDate,
			"updated_at":  now,
		},
	}
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	var out License
	err := r.coll.FindOneAndUpdate(ctx, bson.M{"_id": id, "user_id": userID}, update, opts).Decode(&out)
	return out, err
}

func (r *Repo) Count(ctx context.Context, filter bson.M) (int64, error) {
	return r.coll.CountDocuments(ctx, filter)
}
