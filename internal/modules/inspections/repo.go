package inspections

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
	return &Repo{coll: db.Collection("vehicle_inspections")}
}

func (r *Repo) Insert(ctx context.Context, insp VehicleInspection) (VehicleInspection, error) {
	insp.CreatedAt = time.Now()
	res, err := r.coll.InsertOne(ctx, insp)
	if err != nil {
		return VehicleInspection{}, err
	}
	if oid, ok := res.InsertedID.(primitive.ObjectID); ok {
		insp.ID = oid
	}
	return insp, nil
}

func (r *Repo) FindByID(ctx context.Context, id primitive.ObjectID) (VehicleInspection, error) {
	var insp VehicleInspection
	err := r.coll.FindOne(ctx, bson.M{"_id": id}).Decode(&insp)
	return insp, err
}

func (r *Repo) SetReport(ctx context.Context, id primitive.ObjectID, reportPath string, status string) (VehicleInspection, error) {
	expiry := time.Now().AddDate(1, 0, 0)
	update := bson.M{
		"$set": bson.M{
			"report_path": reportPath,
			"status":      status,
			"expiry_date": expiry,
		},
	}
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	var out VehicleInspection
	err := r.coll.FindOneAndUpdate(ctx, bson.M{"_id": id}, update, opts).Decode(&out)
	return out, err
}

func (r *Repo) Count(ctx context.Context) (int64, error) {
	return r.coll.CountDocuments(ctx, bson.M{})
}
