package domain

type Role string

const (
	RoleCitizen  Role = "citizen"
	RoleAdmin    Role = "admin"
	RoleExaminer Role = "examiner"
	RoleOfficer  Role = "officer"
)

type AccountStatus string

const (
	AccountActive      AccountStatus = "active"
	AccountSuspended   AccountStatus = "suspended"
	AccountBlacklisted AccountStatus = "blacklisted"
)
