package exams

import (
	"context"
	"errors"
	"time"

	"driving-authority-backend/internal/domain"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	examQuestionCount = 30
	passThreshold     = 25
)

type Service struct {
	repo *Repo
}

func NewService(repo *Repo) *Service {
	return &Service{repo: repo}
}

func (s *Service) seedIfEmpty(ctx context.Context, adminID primitive.ObjectID) error {
	n, err := s.repo.CountQuestions(ctx)
	if err != nil {
		return err
	}
	if n > 0 {
		return nil
	}
	now := time.Now()
	seeds := []struct {
		text string
		cat  string
		opts []string
		corr int
	}{
		{"What does a red traffic light mean?", "rules", []string{"Stop", "Go", "Slow down", "Yield"}, 0},
		{"Maximum speed in urban areas unless posted?", "rules", []string{"50 km/h", "80 km/h", "100 km/h", "120 km/h"}, 0},
		{"When must you use headlights?", "safety", []string{"At night and low visibility", "Only at night", "Never in city", "Always off"}, 0},
		{"Safe following distance depends on?", "safety", []string{"Speed and conditions", "Car color", "Music volume", "Fuel level"}, 0},
		{"A yield sign means?", "rules", []string{"Give way when required", "Full stop always", "Speed up", "No entry"}, 0},
		{"Seat belts are required for?", "safety", []string{"All occupants", "Driver only", "Front only", "Children only"}, 0},
		{"Parking on a hill facing uphill with curb?", "rules", []string{"Turn wheels away from curb", "Wheels straight", "Toward curb", "No brake"}, 0},
		{"Blood alcohol limit for drivers is typically?", "rules", []string{"Low or zero", "Unlimited", "High", "Not regulated"}, 0},
		{"Before changing lanes you should?", "safety", []string{"Signal and check mirrors/blind spot", "Honk only", "Speed up blindly", "Close eyes briefly"}, 0},
		{"A pedestrian crosswalk means?", "rules", []string{"Yield to pedestrians", "Ignore pedestrians", "Speed through", "Honk always"}, 0},
	}
	// duplicate variants to reach 30+ questions
	for i := 0; len(seeds) < 35; i++ {
		base := seeds[i%len(seeds)]
		seeds = append(seeds, struct {
			text string
			cat  string
			opts []string
			corr int
		}{base.text + " (variant)", base.cat, base.opts, base.corr})
	}
	for _, sd := range seeds {
		q := Question{
			QuestionText: sd.text,
			Category:     sd.cat,
			CreatedBy:    adminID,
			CreatedAt:    now,
		}
		opts := make([]QuestionOption, len(sd.opts))
		for j, t := range sd.opts {
			opts[j] = QuestionOption{OptionText: t, IsCorrect: j == sd.corr}
		}
		if err := s.repo.InsertQuestion(ctx, q, opts); err != nil {
			return err
		}
	}
	return nil
}

func (s *Service) ListQuestions(ctx context.Context, adminID primitive.ObjectID) ([]QuestionWithOptions, error) {
	if err := s.seedIfEmpty(ctx, adminID); err != nil {
		return nil, err
	}
	return s.repo.ListQuestionsWithOptions(ctx)
}

type StartInput struct {
	LicenseType domain.LicenseType `json:"license_type" binding:"required"`
}

func (s *Service) Start(ctx context.Context, userID primitive.ObjectID, in StartInput) (StartOutput, error) {
	if err := s.seedIfEmpty(ctx, userID); err != nil {
		return StartOutput{}, err
	}
	ids, err := s.repo.RandomQuestionIDs(ctx, examQuestionCount)
	if err != nil {
		return StartOutput{}, err
	}
	if len(ids) < examQuestionCount {
		return StartOutput{}, errors.New("not enough questions in pool")
	}
	attempt := ExamAttempt{
		UserID:      userID,
		LicenseType: in.LicenseType,
		Status:      domain.ExamInProgress,
		QuestionIDs: ids,
		StartedAt:   time.Now(),
	}
	attempt, err = s.repo.InsertAttempt(ctx, attempt)
	if err != nil {
		return StartOutput{}, err
	}
	questions, err := s.repo.QuestionsWithOptionsByIDs(ctx, ids)
	if err != nil {
		return StartOutput{}, err
	}
	return StartOutput{Attempt: attempt, Questions: questions}, nil
}

func (s *Service) Submit(ctx context.Context, attemptID, userID primitive.ObjectID, in SubmitInput) (SubmitOutput, error) {
	attempt, err := s.repo.FindAttempt(ctx, attemptID)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return SubmitOutput{}, errors.New("attempt not found")
		}
		return SubmitOutput{}, err
	}
	if attempt.UserID != userID {
		return SubmitOutput{}, errors.New("forbidden")
	}
	if attempt.Status != domain.ExamInProgress {
		return SubmitOutput{}, errors.New("attempt already submitted")
	}

	score := 0
	answers := make([]ExamAnswer, 0, len(in.Answers))
	for _, a := range in.Answers {
		qid, err := primitive.ObjectIDFromHex(a.QuestionID)
		if err != nil {
			return SubmitOutput{}, errors.New("invalid question_id")
		}
		oid, err := primitive.ObjectIDFromHex(a.SelectedOptionID)
		if err != nil {
			return SubmitOutput{}, errors.New("invalid selected_option_id")
		}
		ok, qFromOpt, err := s.repo.IsOptionCorrect(ctx, oid)
		if err != nil {
			return SubmitOutput{}, err
		}
		if qFromOpt != qid {
			return SubmitOutput{}, errors.New("option does not belong to question")
		}
		if ok {
			score++
		}
		answers = append(answers, ExamAnswer{
			ExamAttemptID:    attemptID,
			QuestionID:       qid,
			SelectedOptionID: oid,
			IsCorrect:        ok,
		})
	}
	if err := s.repo.InsertAnswers(ctx, answers); err != nil {
		return SubmitOutput{}, err
	}

	pct := float64(score) / float64(examQuestionCount) * 100
	status := domain.ExamFailed
	passed := false
	if score >= passThreshold {
		status = domain.ExamPassed
		passed = true
	}
	attempt, err = s.repo.UpdateAttemptResult(ctx, attemptID, score, pct, status)
	if err != nil {
		return SubmitOutput{}, err
	}
	return SubmitOutput{Attempt: attempt, Passed: passed}, nil
}

func (s *Service) History(ctx context.Context, userID primitive.ObjectID) ([]ExamAttempt, error) {
	return s.repo.HistoryByUser(ctx, userID)
}
