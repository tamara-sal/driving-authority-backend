package practical

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
	centers  *mongo.Collection
	slots    *mongo.Collection
	bookings *mongo.Collection
}

func NewRepo(db *mongo.Database) *Repo {
	return &Repo{
		centers:  db.Collection("test_centers"),
		slots:    db.Collection("time_slots"),
		bookings: db.Collection("practical_test_bookings"),
	}
}

func (r *Repo) CountCenters(ctx context.Context) (int64, error) {
	return r.centers.CountDocuments(ctx, bson.M{})
}

func (r *Repo) InsertCenter(ctx context.Context, c TestCenter) (TestCenter, error) {
	c.CreatedAt = time.Now()
	res, err := r.centers.InsertOne(ctx, c)
	if err != nil {
		return TestCenter{}, err
	}
	if oid, ok := res.InsertedID.(primitive.ObjectID); ok {
		c.ID = oid
	}
	return c, nil
}

func (r *Repo) ListCenters(ctx context.Context) ([]TestCenter, error) {
	cur, err := r.centers.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)
	var out []TestCenter
	if err := cur.All(ctx, &out); err != nil {
		return nil, err
	}
	if out == nil {
		out = []TestCenter{}
	}
	return out, nil
}

func (r *Repo) FindCenter(ctx context.Context, id primitive.ObjectID) (TestCenter, error) {
	var c TestCenter
	err := r.centers.FindOne(ctx, bson.M{"_id": id}).Decode(&c)
	return c, err
}

func (r *Repo) ListSlotsByCenter(ctx context.Context, centerID primitive.ObjectID) ([]TimeSlot, error) {
	cur, err := r.slots.Find(ctx, bson.M{"center_id": centerID},
		options.Find().SetSort(bson.D{{Key: "date", Value: 1}, {Key: "start_time", Value: 1}}))
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)
	var out []TimeSlot
	if err := cur.All(ctx, &out); err != nil {
		return nil, err
	}
	if out == nil {
		out = []TimeSlot{}
	}
	return out, nil
}

func (r *Repo) InsertSlot(ctx context.Context, s TimeSlot) (TimeSlot, error) {
	res, err := r.slots.InsertOne(ctx, s)
	if err != nil {
		return TimeSlot{}, err
	}
	if oid, ok := res.InsertedID.(primitive.ObjectID); ok {
		s.ID = oid
	}
	return s, nil
}

func (r *Repo) FindSlot(ctx context.Context, id primitive.ObjectID) (TimeSlot, error) {
	var s TimeSlot
	err := r.slots.FindOne(ctx, bson.M{"_id": id}).Decode(&s)
	return s, err
}

func (r *Repo) IncrementSlotBooked(ctx context.Context, slotID primitive.ObjectID) error {
	res, err := r.slots.UpdateOne(ctx,
		bson.M{"_id": slotID, "$expr": bson.M{"$lt": []any{"$booked", "$capacity"}}},
		bson.M{"$inc": bson.M{"booked": 1}},
	)
	if err != nil {
		return err
	}
	if res.MatchedCount == 0 {
		return mongo.ErrNoDocuments
	}
	return nil
}

func (r *Repo) InsertBooking(ctx context.Context, b PracticalBooking) (PracticalBooking, error) {
	b.CreatedAt = time.Now()
	res, err := r.bookings.InsertOne(ctx, b)
	if err != nil {
		return PracticalBooking{}, err
	}
	if oid, ok := res.InsertedID.(primitive.ObjectID); ok {
		b.ID = oid
	}
	return b, nil
}

func (r *Repo) FindBooking(ctx context.Context, id primitive.ObjectID) (PracticalBooking, error) {
	var b PracticalBooking
	err := r.bookings.FindOne(ctx, bson.M{"_id": id}).Decode(&b)
	return b, err
}

func (r *Repo) SetResult(ctx context.Context, id primitive.ObjectID, examinerID primitive.ObjectID, result domain.PracticalResult, comments string) (PracticalBooking, error) {
	update := bson.M{
		"$set": bson.M{
			"status":      domain.BookingCompleted,
			"result":      result,
			"examiner_id": examinerID,
			"comments":    comments,
		},
	}
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	var out PracticalBooking
	err := r.bookings.FindOneAndUpdate(ctx, bson.M{"_id": id}, update, opts).Decode(&out)
	return out, err
}
