package http

// ErrorResponse is the standard API error body.
type ErrorResponse struct {
	Error string `json:"error" example:"invalid credentials"`
}

// HealthResponse is returned by the health check endpoint.
type HealthResponse struct {
	OK bool `json:"ok" example:"true"`
}

// MeResponse is returned by GET /me.
type MeResponse struct {
	ID    string `json:"id" example:"507f1f77bcf86cd799439011"`
	Email string `json:"email" example:"user@example.com"`
	Role  string `json:"role" example:"citizen"`
}

// AdminPingResponse is returned by GET /admin/ping.
type AdminPingResponse struct {
	Admin bool `json:"admin" example:"true"`
}
