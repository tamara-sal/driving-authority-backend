package platform

import (
	"context"
	"fmt"
	"strings"
	"time"

	"driving-authority-backend/internal/audit"
	"driving-authority-backend/internal/domain"
	"driving-authority-backend/internal/modules/auth"
	"driving-authority-backend/internal/modules/inspections"
	"driving-authority-backend/internal/modules/licenses"
	"driving-authority-backend/internal/modules/vehicles"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Service struct {
	users       *auth.UserRepo
	licenses    *licenses.Repo
	vehicles    *vehicles.Repo
	inspections *inspections.Repo
	audit       *audit.Logger
}

func NewService(users *auth.UserRepo, lic *licenses.Repo, veh *vehicles.Repo, insp *inspections.Repo, auditLog *audit.Logger) *Service {
	return &Service{users: users, licenses: lic, vehicles: veh, inspections: insp, audit: auditLog}
}

type UserView struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	Email      string `json:"email"`
	Role       string `json:"role"`
	Status     string `json:"status"`
	Registered string `json:"registered"`
}

type ApplicationView struct {
	ID        string `json:"id"`
	Applicant string `json:"applicant"`
	Type      string `json:"type"`
	Status    string `json:"status"`
	Submitted string `json:"submitted"`
}

type ActivityView struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Time        string `json:"time"`
}

type AuditView struct {
	ID         string         `json:"id"`
	UserID     string         `json:"user_id"`
	ActionType string         `json:"action_type"`
	Metadata   map[string]any `json:"metadata,omitempty"`
	CreatedAt  string         `json:"created_at"`
}

func (s *Service) ListUsers(ctx context.Context) ([]UserView, error) {
	users, err := s.users.List(ctx)
	if err != nil {
		return nil, err
	}
	out := make([]UserView, 0, len(users))
	for _, u := range users {
		out = append(out, UserView{
			ID:         "USR-" + u.ID.Hex()[len(u.ID.Hex())-4:],
			Name:       strings.TrimSpace(u.FirstName + " " + u.LastName),
			Email:      u.Email,
			Role:       titleRole(u.Role),
			Status:     titleStatus(u.Status),
			Registered: u.CreatedAt.Format("2006-01-02"),
		})
	}
	return out, nil
}

func (s *Service) ListApplications(ctx context.Context) ([]ApplicationView, error) {
	licList, err := s.licenses.FindAll(ctx)
	if err != nil {
		return nil, err
	}
	users, _ := s.users.List(ctx)
	nameByID := map[string]string{}
	for _, u := range users {
		nameByID[u.ID.Hex()] = strings.TrimSpace(u.FirstName + " " + u.LastName)
	}
	out := make([]ApplicationView, 0, len(licList))
	for _, lic := range licList {
		ref := lic.ReferenceID
		if ref == "" {
			ref = "APP-" + lic.ID.Hex()[len(lic.ID.Hex())-8:]
		}
		out = append(out, ApplicationView{
			ID:        ref,
			Applicant: nameByID[lic.UserID.Hex()],
			Type:      licenseTypeLabel(lic.Type),
			Status:    licenseStatusLabel(lic.Status),
			Submitted: lic.CreatedAt.Format("2006-01-02"),
		})
	}
	return out, nil
}

func (s *Service) ListActivity(ctx context.Context, userID primitive.ObjectID) ([]ActivityView, error) {
	logs, err := s.audit.ListByUser(ctx, userID, 20)
	if err != nil {
		return nil, err
	}
	out := make([]ActivityView, 0, len(logs))
	for _, l := range logs {
		out = append(out, ActivityView{
			ID:          l.ID.Hex(),
			Title:       l.ActionType,
			Description: fmt.Sprintf("%v", l.Metadata),
			Time:        l.CreatedAt.Format("2006-01-02 15:04"),
		})
	}
	return out, nil
}

func (s *Service) ListAuditLogs(ctx context.Context) ([]AuditView, error) {
	logs, err := s.audit.List(ctx, 100)
	if err != nil {
		return nil, err
	}
	out := make([]AuditView, 0, len(logs))
	for _, l := range logs {
		out = append(out, AuditView{
			ID:         l.ID.Hex(),
			UserID:     l.UserID.Hex(),
			ActionType: l.ActionType,
			Metadata:   l.Metadata,
			CreatedAt:  l.CreatedAt.Format(time.RFC3339),
		})
	}
	return out, nil
}

func titleRole(r domain.Role) string {
	switch r {
	case domain.RoleCitizen:
		return "Citizen"
	case domain.RoleAdmin:
		return "Admin"
	case domain.RoleExaminer:
		return "Examiner"
	case domain.RoleOfficer:
		return "Officer"
	default:
		return string(r)
	}
}

func titleStatus(s domain.AccountStatus) string {
	switch s {
	case domain.AccountActive:
		return "Active"
	case domain.AccountSuspended:
		return "Suspended"
	default:
		return string(s)
	}
}

func licenseTypeLabel(t domain.LicenseType) string {
	switch t {
	case domain.LicenseMotorcycle:
		return "Motorcycle License"
	case domain.LicenseCar:
		return "Car License"
	case domain.LicenseTruck:
		return "Commercial License"
	case domain.LicenseBus:
		return "Bus License"
	default:
		return string(t)
	}
}

func licenseStatusLabel(s domain.LicenseStatus) string {
	switch s {
	case domain.LicenseSubmitted:
		return "Pending"
	case domain.LicenseApproved, domain.LicenseIssued:
		return "Approved"
	case domain.LicenseRejected:
		return "Rejected"
	default:
		return string(s)
	}
}
