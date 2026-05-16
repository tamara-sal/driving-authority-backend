package monitoring

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Service struct {
	repo *Repo
}

func NewService(repo *Repo) *Service {
	return &Service{repo: repo}
}

func (s *Service) IngestDeviceData(ctx context.Context, in DeviceDataInput) (Trip, error) {
	vid, err := primitive.ObjectIDFromHex(in.VehicleID)
	if err != nil {
		return Trip{}, errors.New("invalid vehicle_id")
	}
	var uid primitive.ObjectID
	if in.UserID != "" {
		uid, err = primitive.ObjectIDFromHex(in.UserID)
		if err != nil {
			return Trip{}, errors.New("invalid user_id")
		}
	}
	if err := s.repo.UpsertDevice(ctx, Device{
		VehicleID:    vid,
		DeviceSerial: in.DeviceSerial,
		Status:       "active",
	}); err != nil {
		return Trip{}, err
	}
	trip := Trip{
		VehicleID:    vid,
		UserID:       uid,
		StartTime:    in.Trip.StartTime,
		EndTime:      in.Trip.EndTime,
		Distance:     in.Trip.Distance,
		AverageSpeed: in.Trip.AverageSpeed,
		SafetyScore:  in.Trip.SafetyScore,
	}
	trip, err = s.repo.InsertTrip(ctx, trip)
	if err != nil {
		return Trip{}, err
	}
	events := make([]TripEvent, 0, len(in.Events))
	for _, e := range in.Events {
		ts := e.Timestamp
		if ts.IsZero() {
			ts = time.Now()
		}
		events = append(events, TripEvent{
			TripID:    trip.ID,
			EventType: e.EventType,
			Severity:  e.Severity,
			Timestamp: ts,
		})
	}
	if err := s.repo.InsertEvents(ctx, events); err != nil {
		return Trip{}, err
	}
	return trip, nil
}

func (s *Service) TripsByVehicle(ctx context.Context, vehicleID primitive.ObjectID) ([]Trip, error) {
	return s.repo.TripsByVehicle(ctx, vehicleID)
}

func (s *Service) ScoreByUser(ctx context.Context, userID primitive.ObjectID) (ScoreOutput, error) {
	avg, count, err := s.repo.AverageScoreByUser(ctx, userID)
	if err != nil {
		return ScoreOutput{}, err
	}
	return ScoreOutput{
		UserID:       userID.Hex(),
		AverageScore: avg,
		TripCount:    count,
	}, nil
}
