package licenses

import (
	"context"
	"errors"
	"fmt"
	"time"

	"driving-authority-backend/internal/domain"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Service struct {
	repo *Repo
}

func NewService(repo *Repo) *Service {
	return &Service{repo: repo}
}

type CreateInput struct {
	Type domain.LicenseType `json:"type" binding:"required"`
}

func (s *Service) Create(ctx context.Context, userID primitive.ObjectID, in CreateInput) (License, error) {
	switch in.Type {
	case domain.LicenseMotorcycle, domain.LicenseCar, domain.LicenseTruck, domain.LicenseBus:
	default:
		return License{}, errors.New("invalid license type")
	}
	lic := License{
		UserID: userID,
		Type:   in.Type,
		Status: domain.LicenseSubmitted,
	}
	return s.repo.Insert(ctx, lic)
}

func (s *Service) MyLicenses(ctx context.Context, userID primitive.ObjectID) ([]License, error) {
	return s.repo.FindByUserID(ctx, userID)
}

func (s *Service) Approve(ctx context.Context, id primitive.ObjectID) (License, error) {
	lic, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return License{}, errors.New("license not found")
		}
		return License{}, err
	}
	if lic.Status != domain.LicenseSubmitted && lic.Status != domain.LicenseApproved {
		return License{}, errors.New("license cannot be approved in current status")
	}
	issueDate := time.Now()
	expiryDate := issueDate.AddDate(5, 0, 0)
	num := fmt.Sprintf("DL-%s-%s", issueDate.Format("2006"), uuid.New().String()[:8])
	return s.repo.Approve(ctx, id, num, issueDate, expiryDate)
}

func (s *Service) Renew(ctx context.Context, id, userID primitive.ObjectID) (License, error) {
	lic, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return License{}, errors.New("license not found")
		}
		return License{}, err
	}
	if lic.UserID != userID {
		return License{}, errors.New("forbidden")
	}
	if lic.Status != domain.LicenseIssued && lic.Status != domain.LicenseExpired {
		return License{}, errors.New("license cannot be renewed in current status")
	}
	base := time.Now()
	if lic.ExpiryDate != nil && lic.ExpiryDate.After(base) {
		base = *lic.ExpiryDate
	}
	expiry := base.AddDate(5, 0, 0)
	return s.repo.Renew(ctx, id, userID, expiry)
}
