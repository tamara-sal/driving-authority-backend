package exams

import (
	"context"
	"time"

	"driving-authority-backend/internal/domain"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Repo struct {
	questions *mongo.Collection
	options   *mongo.Collection
	attempts  *mongo.Collection
	answers   *mongo.Collection
}

func NewRepo(db *mongo.Database) *Repo {
	return &Repo{
		questions: db.Collection("questions"),
		options:   db.Collection("question_options"),
		attempts:  db.Collection("exam_attempts"),
		answers:   db.Collection("exam_answers"),
	}
}

func (r *Repo) CountQuestions(ctx context.Context) (int64, error) {
	return r.questions.CountDocuments(ctx, bson.M{})
}

func (r *Repo) InsertQuestion(ctx context.Context, q Question, opts []QuestionOption) error {
	res, err := r.questions.InsertOne(ctx, q)
	if err != nil {
		return err
	}
	qid := res.InsertedID.(primitive.ObjectID)
	for i := range opts {
		opts[i].QuestionID = qid
	}
	if len(opts) > 0 {
		docs := make([]any, len(opts))
		for j := range opts {
			docs[j] = opts[j]
		}
		_, err = r.options.InsertMany(ctx, docs)
	}
	return err
}

func (r *Repo) ListQuestionsWithOptions(ctx context.Context) ([]QuestionWithOptions, error) {
	cur, err := r.questions.Find(ctx, bson.M{}, options.Find().SetSort(bson.D{{Key: "created_at", Value: 1}}))
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)
	var questions []Question
	if err := cur.All(ctx, &questions); err != nil {
		return nil, err
	}
	out := make([]QuestionWithOptions, 0, len(questions))
	for _, q := range questions {
		opts, err := r.optionsForQuestion(ctx, q.ID)
		if err != nil {
			return nil, err
		}
		out = append(out, QuestionWithOptions{Question: q, Options: opts})
	}
	return out, nil
}

func (r *Repo) optionsForQuestion(ctx context.Context, qid primitive.ObjectID) ([]QuestionOption, error) {
	cur, err := r.options.Find(ctx, bson.M{"question_id": qid})
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)
	var opts []QuestionOption
	if err := cur.All(ctx, &opts); err != nil {
		return nil, err
	}
	if opts == nil {
		opts = []QuestionOption{}
	}
	return opts, nil
}

func (r *Repo) RandomQuestionIDs(ctx context.Context, n int) ([]primitive.ObjectID, error) {
	pipe := mongo.Pipeline{
		{{Key: "$sample", Value: bson.M{"size": n}}},
		{{Key: "$project", Value: bson.M{"_id": 1}}},
	}
	cur, err := r.questions.Aggregate(ctx, pipe)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)
	var ids []primitive.ObjectID
	for cur.Next(ctx) {
		var doc struct {
			ID primitive.ObjectID `bson:"_id"`
		}
		if err := cur.Decode(&doc); err != nil {
			return nil, err
		}
		ids = append(ids, doc.ID)
	}
	return ids, cur.Err()
}

func (r *Repo) QuestionsWithOptionsByIDs(ctx context.Context, ids []primitive.ObjectID) ([]QuestionWithOptions, error) {
	out := make([]QuestionWithOptions, 0, len(ids))
	for _, id := range ids {
		var q Question
		if err := r.questions.FindOne(ctx, bson.M{"_id": id}).Decode(&q); err != nil {
			return nil, err
		}
		opts, err := r.optionsForQuestion(ctx, id)
		if err != nil {
			return nil, err
		}
		// Strip is_correct for exam delivery
		safe := make([]QuestionOption, len(opts))
		for i, o := range opts {
			safe[i] = QuestionOption{ID: o.ID, QuestionID: o.QuestionID, OptionText: o.OptionText}
		}
		out = append(out, QuestionWithOptions{Question: q, Options: safe})
	}
	return out, nil
}

func (r *Repo) InsertAttempt(ctx context.Context, a ExamAttempt) (ExamAttempt, error) {
	res, err := r.attempts.InsertOne(ctx, a)
	if err != nil {
		return ExamAttempt{}, err
	}
	if oid, ok := res.InsertedID.(primitive.ObjectID); ok {
		a.ID = oid
	}
	return a, nil
}

func (r *Repo) FindAttempt(ctx context.Context, id primitive.ObjectID) (ExamAttempt, error) {
	var a ExamAttempt
	err := r.attempts.FindOne(ctx, bson.M{"_id": id}).Decode(&a)
	return a, err
}

func (r *Repo) UpdateAttemptResult(ctx context.Context, id primitive.ObjectID, score int, pct float64, status domain.ExamAttemptStatus) (ExamAttempt, error) {
	now := time.Now()
	update := bson.M{
		"$set": bson.M{
			"score":        score,
			"percentage":   pct,
			"status":       status,
			"submitted_at": now,
		},
	}
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	var out ExamAttempt
	err := r.attempts.FindOneAndUpdate(ctx, bson.M{"_id": id}, update, opts).Decode(&out)
	return out, err
}

func (r *Repo) InsertAnswers(ctx context.Context, answers []ExamAnswer) error {
	if len(answers) == 0 {
		return nil
	}
	docs := make([]any, len(answers))
	for i := range answers {
		docs[i] = answers[i]
	}
	_, err := r.answers.InsertMany(ctx, docs)
	return err
}

func (r *Repo) IsOptionCorrect(ctx context.Context, optionID primitive.ObjectID) (bool, primitive.ObjectID, error) {
	var o QuestionOption
	err := r.options.FindOne(ctx, bson.M{"_id": optionID}).Decode(&o)
	return o.IsCorrect, o.QuestionID, err
}

func (r *Repo) HistoryByUser(ctx context.Context, userID primitive.ObjectID) ([]ExamAttempt, error) {
	cur, err := r.attempts.Find(ctx, bson.M{"user_id": userID},
		options.Find().SetSort(bson.D{{Key: "started_at", Value: -1}}))
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)
	var out []ExamAttempt
	if err := cur.All(ctx, &out); err != nil {
		return nil, err
	}
	if out == nil {
		out = []ExamAttempt{}
	}
	return out, nil
}

func (r *Repo) CountAttempts(ctx context.Context, filter bson.M) (int64, error) {
	return r.attempts.CountDocuments(ctx, filter)
}
