package exams

import (
	"time"

	"driving-authority-backend/internal/domain"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Question struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	QuestionText string             `bson:"question_text" json:"question_text"`
	Category     string             `bson:"category" json:"category"`
	CreatedBy    primitive.ObjectID `bson:"created_by,omitempty" json:"created_by,omitempty"`
	CreatedAt    time.Time          `bson:"created_at" json:"created_at"`
}

type QuestionOption struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	QuestionID primitive.ObjectID `bson:"question_id" json:"question_id"`
	OptionText string             `bson:"option_text" json:"option_text"`
	IsCorrect  bool               `bson:"is_correct" json:"is_correct"`
}

type QuestionWithOptions struct {
	Question Question         `json:"question"`
	Options  []QuestionOption `json:"options"`
}

type ExamAttempt struct {
	ID          primitive.ObjectID    `bson:"_id,omitempty" json:"id"`
	UserID      primitive.ObjectID    `bson:"user_id" json:"user_id"`
	LicenseType domain.LicenseType    `bson:"license_type" json:"license_type"`
	Score       int                   `bson:"score" json:"score"`
	Percentage  float64               `bson:"percentage" json:"percentage"`
	Status      domain.ExamAttemptStatus `bson:"status" json:"status"`
	QuestionIDs []primitive.ObjectID  `bson:"question_ids" json:"question_ids"`
	StartedAt   time.Time             `bson:"started_at" json:"started_at"`
	SubmittedAt *time.Time            `bson:"submitted_at,omitempty" json:"submitted_at,omitempty"`
}

type ExamAnswer struct {
	ID               primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	ExamAttemptID    primitive.ObjectID `bson:"exam_attempt_id" json:"exam_attempt_id"`
	QuestionID       primitive.ObjectID `bson:"question_id" json:"question_id"`
	SelectedOptionID primitive.ObjectID `bson:"selected_option_id" json:"selected_option_id"`
	IsCorrect        bool               `bson:"is_correct" json:"is_correct"`
}

type StartOutput struct {
	Attempt   ExamAttempt           `json:"attempt"`
	Questions []QuestionWithOptions `json:"questions"`
}

type SubmitInput struct {
	Answers []AnswerInput `json:"answers" binding:"required"`
}

type AnswerInput struct {
	QuestionID       string `json:"question_id" binding:"required"`
	SelectedOptionID string `json:"selected_option_id" binding:"required"`
}

type SubmitOutput struct {
	Attempt ExamAttempt `json:"attempt"`
	Passed  bool        `json:"passed"`
}
