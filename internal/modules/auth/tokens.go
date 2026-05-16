package auth

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

var ErrInvalidToken = errors.New("invalid or expired token")

type emailVerification struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	UserID    primitive.ObjectID `bson:"user_id"`
	Token     string             `bson:"token"`
	ExpiresAt time.Time          `bson:"expires_at"`
	Used      bool               `bson:"used"`
}

type passwordReset struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	UserID    primitive.ObjectID `bson:"user_id"`
	Token     string             `bson:"token"`
	ExpiresAt time.Time          `bson:"expires_at"`
	Used      bool               `bson:"used"`
}

type TokenRepo struct {
	emailColl *mongo.Collection
	resetColl *mongo.Collection
	users     *UserRepo
}

func NewTokenRepo(db *mongo.Database, users *UserRepo) *TokenRepo {
	return &TokenRepo{
		emailColl: db.Collection("email_verifications"),
		resetColl: db.Collection("password_resets"),
		users:     users,
	}
}

func newToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

func (r *TokenRepo) CreateEmailVerification(ctx context.Context, userID primitive.ObjectID) (string, error) {
	tok, err := newToken()
	if err != nil {
		return "", err
	}
	_, err = r.emailColl.InsertOne(ctx, emailVerification{
		UserID:    userID,
		Token:     tok,
		ExpiresAt: time.Now().Add(24 * time.Hour),
	})
	return tok, err
}

func (r *TokenRepo) VerifyEmail(ctx context.Context, token string) error {
	var ev emailVerification
	err := r.emailColl.FindOne(ctx, bson.M{"token": token, "used": false}).Decode(&ev)
	if err != nil {
		return ErrInvalidToken
	}
	if time.Now().After(ev.ExpiresAt) {
		return ErrInvalidToken
	}
	_, _ = r.emailColl.UpdateOne(ctx, bson.M{"_id": ev.ID}, bson.M{"$set": bson.M{"used": true}})
	return r.users.SetEmailVerified(ctx, ev.UserID, true)
}

func (r *TokenRepo) CreatePasswordReset(ctx context.Context, email string) (string, error) {
	u, err := r.users.FindByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return "", nil
		}
		return "", err
	}
	tok, err := newToken()
	if err != nil {
		return "", err
	}
	_, err = r.resetColl.InsertOne(ctx, passwordReset{
		UserID:    u.ID,
		Token:     tok,
		ExpiresAt: time.Now().Add(time.Hour),
	})
	return tok, err
}

func (r *TokenRepo) ResetPassword(ctx context.Context, token, newPassword string) error {
	var pr passwordReset
	err := r.resetColl.FindOne(ctx, bson.M{"token": token, "used": false}).Decode(&pr)
	if err != nil {
		return ErrInvalidToken
	}
	if time.Now().After(pr.ExpiresAt) {
		return ErrInvalidToken
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	if err := r.users.UpdatePasswordHash(ctx, pr.UserID, string(hash)); err != nil {
		return err
	}
	_, _ = r.resetColl.UpdateOne(ctx, bson.M{"_id": pr.ID}, bson.M{"$set": bson.M{"used": true}})
	return nil
}
