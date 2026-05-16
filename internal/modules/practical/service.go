package practical

import (
	"context"
	"errors"
	"time"

	"driving-authority-backend/internal/domain"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Service struct {
	repo *Repo
}

func NewService(repo *Repo) *Service {
	return &Service{repo: repo}
}

func (s *Service) seedCenterIfEmpty(ctx context.Context) error {
	n, err := s.repo.CountCenters(ctx)
	if err != nil {
		return err
	}
	if n > 0 {
		return nil
	}
	center, err := s.repo.InsertCenter(ctx, TestCenter{
		Name:            "National Driving Test Center",
		Location:        "Capital City",
		CapacityPerSlot: 5,
	})
	if err != nil {
		return err
	}
	tomorrow := time.Now().AddDate(0, 0, 1).Format("2006-01-02")
	slots := []TimeSlot{
		{CenterID: center.ID, Date: tomorrow, StartTime: "09:00", EndTime: "10:00", Capacity: 5},
		{CenterID: center.ID, Date: tomorrow, StartTime: "10:00", EndTime: "11:00", Capacity: 5},
		{CenterID: center.ID, Date: tomorrow, StartTime: "14:00", EndTime: "15:00", Capacity: 5},
	}
	for _, sl := range slots {
		if _, err := s.repo.InsertSlot(ctx, sl); err != nil {
			return err
		}
	}
	return nil
}

func (s *Service) ListCenters(ctx context.Context) ([]TestCenter, error) {
	if err := s.seedCenterIfEmpty(ctx); err != nil {
		return nil, err
	}
	return s.repo.ListCenters(ctx)
}

func (s *Service) ListSlots(ctx context.Context, centerID primitive.ObjectID) ([]TimeSlot, error) {
	if _, err := s.repo.FindCenter(ctx, centerID); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, errors.New("center not found")
		}
		return nil, err
	}
	return s.repo.ListSlotsByCenter(ctx, centerID)
}

type BookInput struct {
	SlotID string `json:"slot_id" binding:"required"`
}

func (s *Service) Book(ctx context.Context, userID primitive.ObjectID, in BookInput) (PracticalBooking, error) {
	slotID, err := primitive.ObjectIDFromHex(in.SlotID)
	if err != nil {
		return PracticalBooking{}, errors.New("invalid slot_id")
	}
	slot, err := s.repo.FindSlot(ctx, slotID)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return PracticalBooking{}, errors.New("slot not found")
		}
		return PracticalBooking{}, err
	}
	if slot.Booked >= slot.Capacity {
		return PracticalBooking{}, errors.New("slot is full")
	}
	if err := s.repo.IncrementSlotBooked(ctx, slotID); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return PracticalBooking{}, errors.New("slot is full")
		}
		return PracticalBooking{}, err
	}
	return s.repo.InsertBooking(ctx, PracticalBooking{
		UserID:   userID,
		SlotID:   slotID,
		CenterID: slot.CenterID,
		Status:   domain.BookingBooked,
	})
}

type ResultInput struct {
	Result   domain.PracticalResult `json:"result" binding:"required"`
	Comments string                 `json:"comments"`
}

func (s *Service) RecordResult(ctx context.Context, bookingID, examinerID primitive.ObjectID, in ResultInput) (PracticalBooking, error) {
	switch in.Result {
	case domain.PracticalPass, domain.PracticalFail:
	default:
		return PracticalBooking{}, errors.New("invalid result")
	}
	_, err := s.repo.FindBooking(ctx, bookingID)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return PracticalBooking{}, errors.New("booking not found")
		}
		return PracticalBooking{}, err
	}
	return s.repo.SetResult(ctx, bookingID, examinerID, in.Result, in.Comments)
}
