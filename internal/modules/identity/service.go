package identity

import (
	"context"
	"errors"
	"strings"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Service struct {
	repo *Repo
}

func NewService(repo *Repo) *Service {
	return &Service{repo: repo}
}

type SubmitInput struct {
	NationalIDNumber  string `json:"national_id_number" binding:"required"`
	DocumentFrontPath string `json:"document_front_path" binding:"required"`
	DocumentBackPath  string `json:"document_back_path" binding:"required"`
	SelfiePath        string `json:"selfie_path" binding:"required"`
}

func (s *Service) Submit(ctx context.Context, userID primitive.ObjectID, in SubmitInput) (IdentityVerification, error) {
	v := IdentityVerification{
		UserID:            userID,
		NationalIDNumber:  strings.TrimSpace(in.NationalIDNumber),
		DocumentFrontPath: strings.TrimSpace(in.DocumentFrontPath),
		DocumentBackPath:  strings.TrimSpace(in.DocumentBackPath),
		SelfiePath:        strings.TrimSpace(in.SelfiePath),
	}
	return s.repo.UpsertSubmit(ctx, v)
}

func (s *Service) MyStatus(ctx context.Context, userID primitive.ObjectID) (IdentityVerification, error) {
	v, err := s.repo.FindByUserID(ctx, userID)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return IdentityVerification{UserID: userID, Status: ""}, nil
		}
		return IdentityVerification{}, err
	}
	return v, nil
}

func (s *Service) Approve(ctx context.Context, id primitive.ObjectID, reviewer primitive.ObjectID, comment string) (IdentityVerification, error) {
	return s.repo.SetDecision(ctx, id, reviewer, StatusApproved, strings.TrimSpace(comment))
}

func (s *Service) Reject(ctx context.Context, id primitive.ObjectID, reviewer primitive.ObjectID, comment string) (IdentityVerification, error) {
	return s.repo.SetDecision(ctx, id, reviewer, StatusRejected, strings.TrimSpace(comment))
}
