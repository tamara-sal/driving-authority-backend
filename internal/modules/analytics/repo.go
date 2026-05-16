package analytics

import (
	"context"

	"driving-authority-backend/internal/domain"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Repo struct {
	db *mongo.Database
}

func NewRepo(db *mongo.Database) *Repo {
	return &Repo{db: db}
}

func (r *Repo) count(ctx context.Context, coll string, filter bson.M) (int64, error) {
	return r.db.Collection(coll).CountDocuments(ctx, filter)
}

func (r *Repo) Overview(ctx context.Context) (Overview, error) {
	users, err := r.count(ctx, "users", bson.M{})
	if err != nil {
		return Overview{}, err
	}
	licenses, err := r.count(ctx, "licenses", bson.M{})
	if err != nil {
		return Overview{}, err
	}
	vehicles, err := r.count(ctx, "vehicles", bson.M{})
	if err != nil {
		return Overview{}, err
	}
	inspections, err := r.count(ctx, "vehicle_inspections", bson.M{})
	if err != nil {
		return Overview{}, err
	}
	payments, err := r.count(ctx, "payments", bson.M{})
	if err != nil {
		return Overview{}, err
	}
	return Overview{
		Users:       users,
		Licenses:    licenses,
		Vehicles:    vehicles,
		Inspections: inspections,
		Payments:    payments,
	}, nil
}

func (r *Repo) Revenue(ctx context.Context) (Revenue, error) {
	paidCount, err := r.count(ctx, "payments", bson.M{"status": domain.PaymentPaid})
	if err != nil {
		return Revenue{}, err
	}
	pipe := mongo.Pipeline{
		{{Key: "$match", Value: bson.M{"status": domain.PaymentPaid}}},
		{{Key: "$group", Value: bson.M{"_id": nil, "total": bson.M{"$sum": "$amount"}}}},
	}
	cur, err := r.db.Collection("payments").Aggregate(ctx, pipe)
	if err != nil {
		return Revenue{}, err
	}
	defer cur.Close(ctx)
	total := 0.0
	if cur.Next(ctx) {
		var doc struct {
			Total float64 `bson:"total"`
		}
		_ = cur.Decode(&doc)
		total = doc.Total
	}
	return Revenue{TotalPaid: total, PaidCount: paidCount}, nil
}

func (r *Repo) ExamStats(ctx context.Context) (ExamStats, error) {
	total, err := r.count(ctx, "exam_attempts", bson.M{})
	if err != nil {
		return ExamStats{}, err
	}
	passed, err := r.count(ctx, "exam_attempts", bson.M{"status": domain.ExamPassed})
	if err != nil {
		return ExamStats{}, err
	}
	failed, err := r.count(ctx, "exam_attempts", bson.M{"status": domain.ExamFailed})
	if err != nil {
		return ExamStats{}, err
	}
	inProg, err := r.count(ctx, "exam_attempts", bson.M{"status": domain.ExamInProgress})
	if err != nil {
		return ExamStats{}, err
	}
	return ExamStats{
		TotalAttempts: total,
		Passed:        passed,
		Failed:        failed,
		InProgress:    inProg,
	}, nil
}
