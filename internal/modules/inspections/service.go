package inspections

import (
	"context"
	"errors"
	"strings"
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

type ScheduleInput struct {
	VehicleID      string `json:"vehicle_id" binding:"required"`
	InspectionDate string `json:"inspection_date" binding:"required"`
}

func (s *Service) Schedule(ctx context.Context, userID primitive.ObjectID, in ScheduleInput) (VehicleInspection, error) {
	vid, err := primitive.ObjectIDFromHex(in.VehicleID)
	if err != nil {
		return VehicleInspection{}, errors.New("invalid vehicle_id")
	}
	dt, err := time.Parse("2006-01-02", in.InspectionDate)
	if err != nil {
		return VehicleInspection{}, errors.New("inspection_date must be YYYY-MM-DD")
	}
	return s.repo.Insert(ctx, VehicleInspection{
		VehicleID:      vid,
		RequestedBy:    userID,
		InspectionDate: dt,
		Status:         domain.InspectionFailed, // pending until report uploaded
	})
}

type UploadReportInput struct {
	ReportPath string `json:"report_path" binding:"required"`
	Status     string `json:"status"`
}

func (s *Service) UploadReport(ctx context.Context, id, userID primitive.ObjectID, in UploadReportInput) (VehicleInspection, error) {
	insp, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return VehicleInspection{}, errors.New("inspection not found")
		}
		return VehicleInspection{}, err
	}
	if insp.RequestedBy != userID {
		return VehicleInspection{}, errors.New("forbidden")
	}
	status := domain.InspectionPassed
	if in.Status != "" {
		status = domain.InspectionStatus(strings.TrimSpace(in.Status))
	}
	return s.repo.SetReport(ctx, id, strings.TrimSpace(in.ReportPath), string(status))
}
