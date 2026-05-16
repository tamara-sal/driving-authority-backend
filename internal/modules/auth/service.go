package auth

import (
	"context"
	"errors"
	"strings"
	"time"

	"driving-authority-backend/internal/domain"
	"driving-authority-backend/internal/http/middleware"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

var ErrInvalidCredentials = errors.New("invalid credentials")

type Service struct {
	repo                 *UserRepo
	tokens               *TokenRepo
	jwt                  *middleware.JWT
	bootstrapAdminSecret string
}

func NewService(repo *UserRepo, tokens *TokenRepo, jwt *middleware.JWT, bootstrapAdminSecret string) *Service {
	return &Service{repo: repo, tokens: tokens, jwt: jwt, bootstrapAdminSecret: bootstrapAdminSecret}
}

type RegisterInput struct {
	FirstName string `json:"first_name" binding:"required"`
	LastName  string `json:"last_name" binding:"required"`
	Email     string `json:"email" binding:"required,email"`
	Password  string `json:"password" binding:"required,min=8"`
	Phone     string `json:"phone"`
}

type AuthOutput struct {
	AccessToken       string `json:"access_token"`
	UserID            string `json:"user_id"`
	Email             string `json:"email"`
	Role              string `json:"role"`
	VerificationToken string `json:"verification_token,omitempty"`
}

func (s *Service) Register(ctx context.Context, in RegisterInput) (AuthOutput, error) {
	email := strings.ToLower(strings.TrimSpace(in.Email))
	pwHash, err := bcrypt.GenerateFromPassword([]byte(in.Password), bcrypt.DefaultCost)
	if err != nil {
		return AuthOutput{}, err
	}

	u := User{
		UUID:          uuid.NewString(),
		FirstName:     strings.TrimSpace(in.FirstName),
		LastName:      strings.TrimSpace(in.LastName),
		Email:         email,
		PasswordHash:  string(pwHash),
		Phone:         strings.TrimSpace(in.Phone),
		Role:          domain.RoleCitizen,
		Status:        domain.AccountActive,
		EmailVerified: false,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	u, err = s.repo.Insert(ctx, u)
	if err != nil {
		return AuthOutput{}, err
	}

	vTok, _ := s.tokens.CreateEmailVerification(ctx, u.ID)
	tok, err := s.jwt.SignAccessToken(u.ID, u.Email, u.Role)
	if err != nil {
		return AuthOutput{}, err
	}

	return AuthOutput{
		AccessToken:       tok,
		UserID:            u.ID.Hex(),
		Email:             u.Email,
		Role:              string(u.Role),
		VerificationToken: vTok,
	}, nil
}

type BootstrapAdminInput struct {
	Secret string `json:"secret" binding:"required"`
	RegisterInput
}

func (s *Service) BootstrapAdmin(ctx context.Context, in BootstrapAdminInput) (AuthOutput, error) {
	if s.bootstrapAdminSecret == "" || in.Secret != s.bootstrapAdminSecret {
		return AuthOutput{}, errors.New("bootstrap not allowed")
	}

	email := strings.ToLower(strings.TrimSpace(in.Email))
	pwHash, err := bcrypt.GenerateFromPassword([]byte(in.Password), bcrypt.DefaultCost)
	if err != nil {
		return AuthOutput{}, err
	}

	u := User{
		UUID:          uuid.NewString(),
		FirstName:     strings.TrimSpace(in.FirstName),
		LastName:      strings.TrimSpace(in.LastName),
		Email:         email,
		PasswordHash:  string(pwHash),
		Phone:         strings.TrimSpace(in.Phone),
		Role:          domain.RoleAdmin,
		Status:        domain.AccountActive,
		EmailVerified: true,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	u, err = s.repo.Insert(ctx, u)
	if err != nil {
		return AuthOutput{}, err
	}

	tok, err := s.jwt.SignAccessToken(u.ID, u.Email, u.Role)
	if err != nil {
		return AuthOutput{}, err
	}

	return AuthOutput{
		AccessToken: tok,
		UserID:      u.ID.Hex(),
		Email:       u.Email,
		Role:        string(u.Role),
	}, nil
}

type LoginInput struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

func (s *Service) Login(ctx context.Context, in LoginInput) (AuthOutput, error) {
	email := strings.ToLower(strings.TrimSpace(in.Email))
	u, err := s.repo.FindByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return AuthOutput{}, ErrInvalidCredentials
		}
		return AuthOutput{}, err
	}

	if u.Status != domain.AccountActive {
		return AuthOutput{}, errors.New("account is not active")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(in.Password)); err != nil {
		return AuthOutput{}, ErrInvalidCredentials
	}

	tok, err := s.jwt.SignAccessToken(u.ID, u.Email, u.Role)
	if err != nil {
		return AuthOutput{}, err
	}

	return AuthOutput{
		AccessToken: tok,
		UserID:      u.ID.Hex(),
		Email:       u.Email,
		Role:        string(u.Role),
	}, nil
}

type VerifyEmailInput struct {
	Token string `json:"token" binding:"required"`
}

func (s *Service) VerifyEmail(ctx context.Context, in VerifyEmailInput) error {
	return s.tokens.VerifyEmail(ctx, in.Token)
}

type ForgotPasswordInput struct {
	Email string `json:"email" binding:"required,email"`
}

type ForgotPasswordOutput struct {
	ResetToken string `json:"reset_token,omitempty"`
	Message    string `json:"message"`
}

func (s *Service) ForgotPassword(ctx context.Context, in ForgotPasswordInput) (ForgotPasswordOutput, error) {
	tok, err := s.tokens.CreatePasswordReset(ctx, strings.ToLower(strings.TrimSpace(in.Email)))
	if err != nil {
		return ForgotPasswordOutput{}, err
	}
	out := ForgotPasswordOutput{Message: "if the email exists, a reset link was sent"}
	if tok != "" {
		out.ResetToken = tok // dev convenience; remove when email service is wired
	}
	return out, nil
}

type ResetPasswordInput struct {
	Token       string `json:"token" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=8"`
}

func (s *Service) ResetPassword(ctx context.Context, in ResetPasswordInput) error {
	return s.tokens.ResetPassword(ctx, in.Token, in.NewPassword)
}
