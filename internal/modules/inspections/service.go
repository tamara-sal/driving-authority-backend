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

type InspectionView struct {
	ID      string `json:"id"`
	Vehicle string `json:"vehicle"`
	Date    string `json:"date"`
	Result  string `json:"result"`
	Report  string `json:"report"`
}

func (s *Service) List(ctx context.Context, userID primitive.ObjectID, all bool) ([]InspectionView, error) {
	var (
		items []VehicleInspection
		err   error
	)
	if all {
		items, err = s.repo.FindAll(ctx)
	} else {
		items, err = s.repo.FindByUser(ctx, userID)
	}
	if err != nil {
		return nil, err
	}
	out := make([]InspectionView, 0, len(items))
	for _, i := range items {
		result := "Scheduled"
		switch i.Status {
		case domain.InspectionPassed:
			result = "Passed"
		case domain.InspectionFailed:
			if i.ReportPath == "" {
				result = "Scheduled"
			} else {
				result = "Failed"
			}
		}
		report := "Pending"
		if i.ReportPath != "" {
			report = "View report"
		}
		out = append(out, InspectionView{
			ID:      "INS-" + i.ID.Hex()[len(i.ID.Hex())-4:],
			Vehicle: i.VehicleID.Hex()[len(i.VehicleID.Hex())-4:],
			Date:    i.InspectionDate.Format("2006-01-02"),
			Result:  result,
			Report:  report,
		})
	}
	return out, nil
}
