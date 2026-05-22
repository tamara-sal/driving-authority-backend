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
	PlateNumber string `json:"plate_number"`
	Plate       string `json:"plate"`
	Make        string `json:"make" binding:"required"`
	Model       string `json:"model" binding:"required"`
	Year        int    `json:"year" binding:"required"`
	Color       string `json:"color"`
}

type VehicleView struct {
	Plate  string `json:"plate"`
	VIN    string `json:"vin"`
	Make   string `json:"make"`
	Model  string `json:"model"`
	Year   int    `json:"year"`
	Color  string `json:"color"`
	Status string `json:"status"`
}

func toVehicleView(v Vehicle) VehicleView {
	status := "Active"
	switch v.Status {
	case domain.VehicleSuspended:
		status = "Suspended"
	case domain.VehicleStolen:
		status = "Suspended"
	}
	return VehicleView{
		Plate: v.PlateNumber, VIN: v.VIN, Make: v.Make, Model: v.Model,
		Year: v.Year, Color: v.Color, Status: status,
	}
}

func (s *Service) Create(ctx context.Context, ownerID primitive.ObjectID, in CreateInput) (VehicleView, error) {
	plate := strings.TrimSpace(in.PlateNumber)
	if plate == "" {
		plate = strings.TrimSpace(in.Plate)
	}
	if plate == "" {
		return VehicleView{}, errors.New("plate or plate_number required")
	}
	v := Vehicle{
		VIN:         strings.TrimSpace(in.VIN),
		PlateNumber: plate,
		Make:        strings.TrimSpace(in.Make),
		Model:       strings.TrimSpace(in.Model),
		Year:        in.Year,
		Color:       strings.TrimSpace(in.Color),
		Status:      domain.VehicleActive,
		OwnerID:     ownerID,
	}
	created, err := s.repo.InsertVehicle(ctx, v)
	if err != nil {
		return VehicleView{}, err
	}
	return toVehicleView(created), nil
}

func (s *Service) MyVehicles(ctx context.Context, ownerID primitive.ObjectID) ([]VehicleView, error) {
	list, err := s.repo.FindByOwner(ctx, ownerID)
	if err != nil {
		return nil, err
	}
	out := make([]VehicleView, 0, len(list))
	for _, v := range list {
		out = append(out, toVehicleView(v))
	}
	return out, nil
}

func (s *Service) AllVehicles(ctx context.Context) ([]VehicleView, error) {
	list, err := s.repo.FindAll(ctx)
	if err != nil {
		return nil, err
	}
	out := make([]VehicleView, 0, len(list))
	for _, v := range list {
		out = append(out, toVehicleView(v))
	}
	return out, nil
}

type TransferInput struct {
	BuyerID    string `json:"buyer_id"`
	BuyerEmail string `json:"buyer_email"`
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
	buyerID, err := primitive.ObjectIDFromHex(strings.TrimSpace(in.BuyerID))
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
