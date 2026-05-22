package licenses

import (
	"context"
	"errors"
	"fmt"
	"strings"
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
	Type        domain.LicenseType    `json:"type"`
	LicenseType string                `json:"license_type"`
	Name        string                `json:"name"`
	DOB         string                `json:"dob"`
	Gender      string                `json:"gender"`
	Nationality string                `json:"nationality"`
	Address     string                `json:"address"`
	City        string                `json:"city"`
	Postal      string                `json:"postal"`
}

func (s *Service) Create(ctx context.Context, userID primitive.ObjectID, in CreateInput) (License, error) {
	licType, err := resolveLicenseType(in.Type, in.LicenseType)
	if err != nil {
		return License{}, err
	}
	ref := fmt.Sprintf("APP-%s-%s", time.Now().Format("2006"), uuid.New().String()[:8])
	lic := License{
		UserID:      userID,
		ReferenceID: ref,
		Type:        licType,
		Status:      domain.LicenseSubmitted,
	}
	if in.Name != "" || in.DOB != "" {
		lic.Application = &ApplicationDetails{
			Name: in.Name, DOB: in.DOB, Gender: in.Gender,
			Nationality: in.Nationality, Address: in.Address, City: in.City, Postal: in.Postal,
		}
	}
	return s.repo.Insert(ctx, lic)
}

func resolveLicenseType(enum domain.LicenseType, label string) (domain.LicenseType, error) {
	if enum != "" {
		switch enum {
		case domain.LicenseMotorcycle, domain.LicenseCar, domain.LicenseTruck, domain.LicenseBus:
			return enum, nil
		}
	}
	switch strings.ToLower(strings.TrimSpace(label)) {
	case "motorcycle license", "motorcycle":
		return domain.LicenseMotorcycle, nil
	case "car license", "car", "class b":
		return domain.LicenseCar, nil
	case "commercial license", "commercial", "truck":
		return domain.LicenseTruck, nil
	case "bus license", "bus":
		return domain.LicenseBus, nil
	}
	return "", errors.New("invalid license type")
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

func (s *Service) Reject(ctx context.Context, id primitive.ObjectID) (License, error) {
	lic, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return License{}, errors.New("license not found")
		}
		return License{}, err
	}
	if lic.Status != domain.LicenseSubmitted {
		return License{}, errors.New("only pending applications can be rejected")
	}
	return s.repo.Reject(ctx, id)
}
