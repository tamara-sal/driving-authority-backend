package analytics

type Overview struct {
	Users       int64 `json:"users"`
	Licenses    int64 `json:"licenses"`
	Vehicles    int64 `json:"vehicles"`
	Inspections int64 `json:"inspections"`
	Payments    int64 `json:"payments"`
}

type Revenue struct {
	TotalPaid float64 `json:"total_paid"`
	PaidCount int64   `json:"paid_count"`
}

type ExamStats struct {
	TotalAttempts int64 `json:"total_attempts"`
	Passed        int64 `json:"passed"`
	Failed        int64 `json:"failed"`
	InProgress    int64 `json:"in_progress"`
}
