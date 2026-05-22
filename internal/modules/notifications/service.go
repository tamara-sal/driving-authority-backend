package notifications

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Service struct {
	repo *Repo
}

func NewService(repo *Repo) *Service {
	return &Service{repo: repo}
}

func (s *Service) List(ctx context.Context, userID primitive.ObjectID) ([]NotificationView, error) {
	items, err := s.repo.ListByUser(ctx, userID)
	if err != nil {
		return nil, err
	}
	out := make([]NotificationView, 0, len(items))
	for _, n := range items {
		out = append(out, toView(n))
	}
	return out, nil
}

func (s *Service) MarkRead(ctx context.Context, id, userID primitive.ObjectID) (NotificationView, error) {
	n, err := s.repo.MarkRead(ctx, id, userID)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return NotificationView{}, fmt.Errorf("notification not found")
		}
		return NotificationView{}, err
	}
	return toView(n), nil
}

func (s *Service) Create(ctx context.Context, userID primitive.ObjectID, title, description, nType string) error {
	_, err := s.repo.Insert(ctx, Notification{
		UserID:      userID,
		Title:       title,
		Description: description,
		Status:      "unread",
		Type:        nType,
	})
	return err
}

func toView(n Notification) NotificationView {
	return NotificationView{
		ID:          n.ID.Hex(),
		Title:       n.Title,
		Description: n.Description,
		Timestamp:   formatTimestamp(n.CreatedAt),
		Status:      n.Status,
		Type:        n.Type,
	}
}

func formatTimestamp(t time.Time) string {
	d := time.Since(t)
	switch {
	case d < time.Hour:
		return fmt.Sprintf("%d min ago", int(d.Minutes()))
	case d < 48*time.Hour:
		return "Yesterday"
	default:
		return t.Format("2006-01-02")
	}
}
