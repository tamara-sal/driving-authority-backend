package monitoring

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Repo struct {
	devices *mongo.Collection
	trips   *mongo.Collection
	events  *mongo.Collection
}

func NewRepo(db *mongo.Database) *Repo {
	return &Repo{
		devices: db.Collection("devices"),
		trips:   db.Collection("trips"),
		events:  db.Collection("trip_events"),
	}
}

func (r *Repo) UpsertDevice(ctx context.Context, d Device) error {
	_, err := r.devices.UpdateOne(ctx,
		bson.M{"device_serial": d.DeviceSerial},
		bson.M{"$set": d},
		options.Update().SetUpsert(true),
	)
	return err
}

func (r *Repo) InsertTrip(ctx context.Context, t Trip) (Trip, error) {
	res, err := r.trips.InsertOne(ctx, t)
	if err != nil {
		return Trip{}, err
	}
	if oid, ok := res.InsertedID.(primitive.ObjectID); ok {
		t.ID = oid
	}
	return t, nil
}

func (r *Repo) InsertEvents(ctx context.Context, events []TripEvent) error {
	if len(events) == 0 {
		return nil
	}
	docs := make([]any, len(events))
	for i := range events {
		docs[i] = events[i]
	}
	_, err := r.events.InsertMany(ctx, docs)
	return err
}

func (r *Repo) TripsByVehicle(ctx context.Context, vehicleID primitive.ObjectID) ([]Trip, error) {
	cur, err := r.trips.Find(ctx, bson.M{"vehicle_id": vehicleID},
		options.Find().SetSort(bson.D{{Key: "start_time", Value: -1}}))
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)
	var out []Trip
	if err := cur.All(ctx, &out); err != nil {
		return nil, err
	}
	if out == nil {
		out = []Trip{}
	}
	return out, nil
}

func (r *Repo) AverageScoreByUser(ctx context.Context, userID primitive.ObjectID) (float64, int64, error) {
	pipe := mongo.Pipeline{
		{{Key: "$match", Value: bson.M{"user_id": userID}}},
		{{Key: "$group", Value: bson.M{
			"_id":   nil,
			"avg":   bson.M{"$avg": "$safety_score"},
			"count": bson.M{"$sum": 1},
		}}},
	}
	cur, err := r.trips.Aggregate(ctx, pipe)
	if err != nil {
		return 0, 0, err
	}
	defer cur.Close(ctx)
	if !cur.Next(ctx) {
		return 0, 0, nil
	}
	var doc struct {
		Avg   float64 `bson:"avg"`
		Count int64   `bson:"count"`
	}
	if err := cur.Decode(&doc); err != nil {
		return 0, 0, err
	}
	return doc.Avg, doc.Count, nil
}
