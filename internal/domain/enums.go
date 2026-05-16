package domain

type LicenseType string

const (
	LicenseMotorcycle LicenseType = "motorcycle"
	LicenseCar        LicenseType = "car"
	LicenseTruck      LicenseType = "truck"
	LicenseBus        LicenseType = "bus"
)

type LicenseStatus string

const (
	LicenseDraft     LicenseStatus = "draft"
	LicenseSubmitted LicenseStatus = "submitted"
	LicenseApproved  LicenseStatus = "approved"
	LicenseIssued    LicenseStatus = "issued"
	LicenseSuspended LicenseStatus = "suspended"
	LicenseExpired   LicenseStatus = "expired"
	LicenseRevoked   LicenseStatus = "revoked"
)

type ExamAttemptStatus string

const (
	ExamInProgress ExamAttemptStatus = "in_progress"
	ExamPassed     ExamAttemptStatus = "passed"
	ExamFailed     ExamAttemptStatus = "failed"
)

type BookingStatus string

const (
	BookingBooked     BookingStatus = "booked"
	BookingCompleted  BookingStatus = "completed"
	BookingCancelled  BookingStatus = "cancelled"
)

type PracticalResult string

const (
	PracticalPass PracticalResult = "pass"
	PracticalFail PracticalResult = "fail"
)

type VehicleStatus string

const (
	VehicleActive    VehicleStatus = "active"
	VehicleSuspended VehicleStatus = "suspended"
	VehicleStolen    VehicleStatus = "stolen"
)

type TransferStatus string

const (
	TransferPending  TransferStatus = "pending"
	TransferApproved TransferStatus = "approved"
	TransferRejected TransferStatus = "rejected"
)

type InspectionStatus string

const (
	InspectionPassed  InspectionStatus = "passed"
	InspectionFailed  InspectionStatus = "failed"
	InspectionExpired InspectionStatus = "expired"
)

type PaymentStatus string

const (
	PaymentPending PaymentStatus = "pending"
	PaymentPaid    PaymentStatus = "paid"
	PaymentFailed  PaymentStatus = "failed"
)

type ViolationStatus string

const (
	ViolationPending  ViolationStatus = "pending"
	ViolationPaid     ViolationStatus = "paid"
	ViolationDisputed ViolationStatus = "disputed"
)

type ServiceType string

const (
	ServiceExamFee      ServiceType = "theory_exam"
	ServiceLicenseFee   ServiceType = "license"
	ServiceInspection   ServiceType = "inspection"
	ServiceTransfer     ServiceType = "transfer"
)
