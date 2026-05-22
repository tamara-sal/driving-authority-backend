package auth

import (
	"context"
	"errors"
	"strings"
	"time"

	"driving-authority-backend/internal/domain"
	"driving-authority-backend/internal/http/middleware"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	FullName  string `json:"full_name"`
	Email     string `json:"email" binding:"required,email"`
	Password  string `json:"password" binding:"required,min=8"`
	Phone     string `json:"phone"`
}

type AuthOutput struct {
	AccessToken       string `json:"access_token"`
	UserID            string `json:"user_id"`
	Email             string `json:"email"`
	Name              string `json:"name"`
	Role              string `json:"role"`
	Redirect          string `json:"redirect,omitempty"`
	VerificationToken string `json:"verification_token,omitempty"`
}

func normalizeRegisterInput(in *RegisterInput) error {
	if strings.TrimSpace(in.FirstName) == "" && strings.TrimSpace(in.LastName) == "" {
		full := strings.TrimSpace(in.FullName)
		if full == "" {
			return errors.New("full_name or first_name and last_name required")
		}
		parts := strings.Fields(full)
		in.FirstName = parts[0]
		if len(parts) > 1 {
			in.LastName = strings.Join(parts[1:], " ")
		}
	}
	if strings.TrimSpace(in.FirstName) == "" {
		return errors.New("first_name required")
	}
	return nil
}

func userDisplayName(u User) string {
	return strings.TrimSpace(u.FirstName + " " + u.LastName)
}

func roleRedirect(role domain.Role) string {
	switch role {
	case domain.RoleAdmin:
		return "/dashboard/admin"
	case domain.RoleExaminer:
		return "/dashboard/examiner"
	case domain.RoleOfficer:
		return "/dashboard/officer"
	default:
		return "/dashboard/citizen"
	}
}

func authOutput(u User, tok, vTok string) AuthOutput {
	return AuthOutput{
		AccessToken:       tok,
		UserID:            u.ID.Hex(),
		Email:             u.Email,
		Name:              userDisplayName(u),
		Role:              string(u.Role),
		Redirect:          roleRedirect(u.Role),
		VerificationToken: vTok,
	}
}

func (s *Service) Register(ctx context.Context, in RegisterInput) (AuthOutput, error) {
	if err := normalizeRegisterInput(&in); err != nil {
		return AuthOutput{}, err
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

	return authOutput(u, tok, vTok), nil
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

	return authOutput(u, tok, ""), nil
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

	return authOutput(u, tok, ""), nil
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

func (s *Service) GetUser(ctx context.Context, id primitive.ObjectID) (User, error) {
	return s.repo.FindByID(ctx, id)
}

// DemoPassword is the shared password for seeded demo accounts (matches frontend README).
const DemoPassword = "Password123!"

var demoAccounts = []User{
	{FirstName: "John", LastName: "Citizen", Email: "citizen@example.com", Role: domain.RoleCitizen},
	{FirstName: "Alice", LastName: "Admin", Email: "admin@example.com", Role: domain.RoleAdmin},
	{FirstName: "Bob", LastName: "Examiner", Email: "examiner@example.com", Role: domain.RoleExaminer},
	{FirstName: "Officer", LastName: "Davis", Email: "officer@example.com", Role: domain.RoleOfficer},
}

type SeedDemoOutput struct {
	Message  string   `json:"message"`
	Password string   `json:"password"`
	Accounts []string `json:"accounts"`
}

func (s *Service) SeedDemoUsers(ctx context.Context) (SeedDemoOutput, error) {
	pwHash, err := bcrypt.GenerateFromPassword([]byte(DemoPassword), bcrypt.DefaultCost)
	if err != nil {
		return SeedDemoOutput{}, err
	}
	emails := make([]string, 0, len(demoAccounts))
	for _, d := range demoAccounts {
		u := User{
			UUID:          uuid.NewString(),
			FirstName:     d.FirstName,
			LastName:      d.LastName,
			Email:         strings.ToLower(d.Email),
			PasswordHash:  string(pwHash),
			Phone:         "+10000000000",
			Role:          d.Role,
			Status:        domain.AccountActive,
			EmailVerified: true,
		}
		if err := s.repo.UpsertDemo(ctx, u); err != nil {
			return SeedDemoOutput{}, err
		}
		emails = append(emails, u.Email)
	}
	return SeedDemoOutput{
		Message:  "4 demo users ready",
		Password: DemoPassword,
		Accounts: emails,
	}, nil
}
