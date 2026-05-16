package vehicles

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
	vehicles  *mongo.Collection
	transfers *mongo.Collection
}

func NewRepo(db *mongo.Database) *Repo {
	v := db.Collection("vehicles")
	t := db.Collection("vehicle_transfer_requests")
	_, _ = v.Indexes().CreateMany(context.Background(), []mongo.IndexModel{
		{Keys: bson.D{{Key: "vin", Value: 1}}, Options: options.Index().SetUnique(true)},
		{Keys: bson.D{{Key: "plate_number", Value: 1}}, Options: options.Index().SetUnique(true)},
		{Keys: bson.D{{Key: "owner_id", Value: 1}}},
	})
	return &Repo{vehicles: v, transfers: t}
}

func (r *Repo) InsertVehicle(ctx context.Context, v Vehicle) (Vehicle, error) {
	v.CreatedAt = time.Now()
	res, err := r.vehicles.InsertOne(ctx, v)
	if err != nil {
		return Vehicle{}, err
	}
	if oid, ok := res.InsertedID.(primitive.ObjectID); ok {
		v.ID = oid
	}
	return v, nil
}

func (r *Repo) FindByOwner(ctx context.Context, ownerID primitive.ObjectID) ([]Vehicle, error) {
	cur, err := r.vehicles.Find(ctx, bson.M{"owner_id": ownerID})
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)
	var out []Vehicle
	if err := cur.All(ctx, &out); err != nil {
		return nil, err
	}
	if out == nil {
		out = []Vehicle{}
	}
	return out, nil
}

func (r *Repo) FindByID(ctx context.Context, id primitive.ObjectID) (Vehicle, error) {
	var v Vehicle
	err := r.vehicles.FindOne(ctx, bson.M{"_id": id}).Decode(&v)
	return v, err
}

func (r *Repo) InsertTransfer(ctx context.Context, t TransferRequest) (TransferRequest, error) {
	t.CreatedAt = time.Now()
	res, err := r.transfers.InsertOne(ctx, t)
	if err != nil {
		return TransferRequest{}, err
	}
	if oid, ok := res.InsertedID.(primitive.ObjectID); ok {
		t.ID = oid
	}
	return t, nil
}

func (r *Repo) FindTransfer(ctx context.Context, id primitive.ObjectID) (TransferRequest, error) {
	var t TransferRequest
	err := r.transfers.FindOne(ctx, bson.M{"_id": id}).Decode(&t)
	return t, err
}

func (r *Repo) ApproveTransfer(ctx context.Context, id, reviewer primitive.ObjectID) (TransferRequest, Vehicle, error) {
	tr, err := r.FindTransfer(ctx, id)
	if err != nil {
		return TransferRequest{}, Vehicle{}, err
	}
	now := time.Now()
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	err = r.transfers.FindOneAndUpdate(ctx, bson.M{"_id": id}, bson.M{
		"$set": bson.M{
			"status":      domain.TransferApproved,
			"reviewed_by": reviewer,
			"reviewed_at": now,
		},
	}, opts).Decode(&tr)
	if err != nil {
		return TransferRequest{}, Vehicle{}, err
	}
	vopts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	var v Vehicle
	err = r.vehicles.FindOneAndUpdate(ctx, bson.M{"_id": tr.VehicleID}, bson.M{
		"$set": bson.M{"owner_id": tr.BuyerID},
	}, vopts).Decode(&v)
	return tr, v, err
}

func (r *Repo) CountVehicles(ctx context.Context) (int64, error) {
	return r.vehicles.CountDocuments(ctx, bson.M{})
}

func (r *Repo) CountTransfers(ctx context.Context, filter bson.M) (int64, error) {
	return r.transfers.CountDocuments(ctx, filter)
}
