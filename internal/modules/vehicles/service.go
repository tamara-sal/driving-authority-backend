package vehicles

import (
	"context"
	"errors"
	"strings"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"driving-authority-backend/internal/domain"
)

type Service struct {
	repo *Repo
}

func NewService(repo *Repo) *Service {
	return &Service{repo: repo}
}

type CreateInput struct {
	VIN         string `json:"vin" binding:"required"`
	PlateNumber string `json:"plate_number" binding:"required"`
	Make        string `json:"make" binding:"required"`
	Model       string `json:"model" binding:"required"`
	Year        int    `json:"year" binding:"required"`
}

func (s *Service) Create(ctx context.Context, ownerID primitive.ObjectID, in CreateInput) (Vehicle, error) {
	v := Vehicle{
		VIN:         strings.TrimSpace(in.VIN),
		PlateNumber: strings.TrimSpace(in.PlateNumber),
		Make:        strings.TrimSpace(in.Make),
		Model:       strings.TrimSpace(in.Model),
		Year:        in.Year,
		Status:      domain.VehicleActive,
		OwnerID:     ownerID,
	}
	return s.repo.InsertVehicle(ctx, v)
}

func (s *Service) MyVehicles(ctx context.Context, ownerID primitive.ObjectID) ([]Vehicle, error) {
	return s.repo.FindByOwner(ctx, ownerID)
}

type TransferInput struct {
	BuyerID string `json:"buyer_id" binding:"required"`
}

func (s *Service) RequestTransfer(ctx context.Context, vehicleID, sellerID primitive.ObjectID, in TransferInput) (TransferRequest, error) {
	v, err := s.repo.FindByID(ctx, vehicleID)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return TransferRequest{}, errors.New("vehicle not found")
		}
		return TransferRequest{}, err
	}
	if v.OwnerID != sellerID {
		return TransferRequest{}, errors.New("forbidden")
	}
	buyerID, err := primitive.ObjectIDFromHex(in.BuyerID)
	if err != nil {
		return TransferRequest{}, errors.New("invalid buyer_id")
	}
	return s.repo.InsertTransfer(ctx, TransferRequest{
		VehicleID: vehicleID,
		SellerID:  sellerID,
		BuyerID:   buyerID,
		Status:    domain.TransferPending,
	})
}

func (s *Service) ApproveTransfer(ctx context.Context, transferID, adminID primitive.ObjectID) (TransferRequest, error) {
	tr, err := s.repo.FindTransfer(ctx, transferID)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return TransferRequest{}, errors.New("transfer not found")
		}
		return TransferRequest{}, err
	}
	if tr.Status != domain.TransferPending {
		return TransferRequest{}, errors.New("transfer is not pending")
	}
	tr, _, err = s.repo.ApproveTransfer(ctx, transferID, adminID)
	return tr, err
}
