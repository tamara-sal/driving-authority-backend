# API paths (see Swagger UI for full schemas)

All under `/api/v1` unless noted.

## System
- GET `/` - API info
- GET `/health`
- GET `/swagger/index.html`

## Auth
- POST `/auth/register`, `/auth/login`, `/auth/verify-email`, `/auth/forgot-password`, `/auth/reset-password`
- POST `/auth/seed-demo` — upsert 4 demo users (citizen/admin/examiner/officer @example.com)

## Identity
- POST `/identity/submit`, GET `/identity/status`
- PUT `/admin/identity/{id}/approve`, `/admin/identity/{id}/reject`

## Licenses
- POST `/licenses`, GET `/licenses/me`, PUT `/licenses/{id}/renew`
- PUT `/admin/licenses/{id}/approve`, `/admin/licenses/{id}/reject`

## Theory exam
- GET `/exam/questions`, POST `/exam/start`, POST `/exam/{attemptId}/submit`, GET `/exam/history`

## Practical test
- GET `/centers`, GET `/centers/{id}/slots`, POST `/practical/book`
- PUT `/examiner/practical/{id}/result`

## Vehicles
- POST `/vehicles`, GET `/vehicles/me`, POST `/vehicles/{id}/transfer`
- PUT `/admin/transfer/{id}/approve`

## Inspection
- POST `/inspection/schedule`, POST `/inspection/{id}/upload-report`

## Monitoring
- POST `/devices/data`, GET `/monitoring/trips/{vehicleId}`, GET `/monitoring/score/{userId}`

## Payments
- POST `/payments/initiate`, GET `/payments/history`
- PUT `/admin/payments/{id}/mark-paid`

## Analytics (admin)
- GET `/admin/analytics/overview`, `/admin/analytics/revenue`, `/admin/analytics/exams`

## Notifications
- GET `/notifications`, PATCH `/notifications/{id}/read`

## Violations
- GET `/violations`, POST `/violations` (officer/admin), PUT `/violations/{id}/status`

## Activity & admin lists
- GET `/activity`
- GET `/admin/users`, `/admin/applications`, `/admin/audit-logs`
- PUT `/admin/licenses/{id}/reject`
- GET `/admin/vehicles`, `/admin/analytics/trends`

## Profile
- GET `/me`, GET `/admin/ping`
