package payments

import (
	"context"
	"errors"

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

func (s *Service) seedFeesIfEmpty(ctx context.Context) error {
	n, err := s.repo.CountFees(ctx)
	if err != nil {
		return err
	}
	if n > 0 {
		return nil
	}
	return s.repo.InsertFees(ctx, []ServiceFee{
		{ServiceType: domain.ServiceExamFee, Amount: 25, Currency: "USD"},
		{ServiceType: domain.ServiceLicenseFee, Amount: 50, Currency: "USD"},
		{ServiceType: domain.ServiceInspection, Amount: 40, Currency: "USD"},
		{ServiceType: domain.ServiceTransfer, Amount: 30, Currency: "USD"},
	})
}

type InitiateInput struct {
	ServiceType domain.ServiceType `json:"service_type" binding:"required"`
}

func (s *Service) Initiate(ctx context.Context, userID primitive.ObjectID, in InitiateInput) (Payment, error) {
	if err := s.seedFeesIfEmpty(ctx); err != nil {
		return Payment{}, err
	}
	fee, err := s.repo.FindFeeByType(ctx, in.ServiceType)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return Payment{}, errors.New("unknown service type")
		}
		return Payment{}, err
	}
	return s.repo.InsertPayment(ctx, Payment{
		UserID:        userID,
		ServiceType:   fee.ServiceType,
		Amount:        fee.Amount,
		Currency:      fee.Currency,
		Status:        domain.PaymentPending,
		TransactionID: uuid.NewString(),
	})
}

func (s *Service) MarkPaid(ctx context.Context, id primitive.ObjectID) (Payment, error) {
	return s.repo.MarkPaid(ctx, id)
}

func (s *Service) History(ctx context.Context, userID primitive.ObjectID) ([]Payment, error) {
	return s.repo.HistoryByUser(ctx, userID)
}

func (s *Service) TotalRevenue(ctx context.Context) (float64, error) {
	return s.repo.SumPaid(ctx)
}
