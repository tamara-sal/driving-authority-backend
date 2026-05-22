package violations

import (
	"context"
	"errors"
	"strings"

	"driving-authority-backend/internal/domain"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Service struct {
	repo *Repo
}

func NewService(repo *Repo) *Service {
	return &Service{repo: repo}
}

func (s *Service) ListAll(ctx context.Context) ([]ViolationView, error) {
	return s.list(ctx, bson.M{})
}

func (s *Service) ListByDriver(ctx context.Context, driverID primitive.ObjectID) ([]ViolationView, error) {
	return s.list(ctx, bson.M{"driver_id": driverID})
}

func (s *Service) list(ctx context.Context, filter bson.M) ([]ViolationView, error) {
	items, err := s.repo.List(ctx, filter)
	if err != nil {
		return nil, err
	}
	out := make([]ViolationView, 0, len(items))
	for _, v := range items {
		out = append(out, toView(v))
	}
	return out, nil
}

func (s *Service) Create(ctx context.Context, in CreateInput) (ViolationView, error) {
	var driverID primitive.ObjectID
	var err error
	if in.DriverID != "" {
		driverID, err = primitive.ObjectIDFromHex(in.DriverID)
		if err != nil {
			return ViolationView{}, errors.New("invalid driver_id")
		}
	}
	status := domain.ViolationPending
	v, err := s.repo.Insert(ctx, Violation{
		DriverID: driverID,
		Driver:   strings.TrimSpace(in.Driver),
		Type:     strings.TrimSpace(in.Type),
		Severity: strings.TrimSpace(in.Severity),
		Status:   status,
	})
	if err != nil {
		return ViolationView{}, err
	}
	return toView(v), nil
}

func (s *Service) UpdateStatus(ctx context.Context, id primitive.ObjectID, status domain.ViolationStatus) (ViolationView, error) {
	switch status {
	case domain.ViolationPending, domain.ViolationPaid, domain.ViolationDisputed:
	default:
		return ViolationView{}, errors.New("invalid status")
	}
	v, err := s.repo.UpdateStatus(ctx, id, string(status))
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return ViolationView{}, errors.New("violation not found")
		}
		return ViolationView{}, err
	}
	return toView(v), nil
}

func toView(v Violation) ViolationView {
	status := string(v.Status)
	switch status {
	case "pending":
		status = "Open"
	case "paid":
		status = "Paid"
	case "disputed":
		status = "Appealed"
	}
	return ViolationView{
		ID:       v.ID.Hex(),
		Driver:   v.Driver,
		Type:     v.Type,
		Severity: v.Severity,
		Date:     v.CreatedAt.Format("2006-01-02"),
		Status:   status,
	}
}
