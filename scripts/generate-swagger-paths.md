# API paths (see Swagger UI for full schemas)

All under `/api/v1` unless noted.

## System
- GET `/` - API info
- GET `/health`
- GET `/swagger/index.html`

## Auth
- POST `/auth/register`, `/auth/login`, `/auth/verify-email`, `/auth/forgot-password`, `/auth/reset-password`

## Identity
- POST `/identity/submit`, GET `/identity/status`
- PUT `/admin/identity/{id}/approve`, `/admin/identity/{id}/reject`

## Licenses
- POST `/licenses`, GET `/licenses/me`, PUT `/licenses/{id}/renew`
- PUT `/admin/licenses/{id}/approve`

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

## Profile
- GET `/me`, GET `/admin/ping`
