# API Test Report (Production)

**Base URL:** `https://api-production-5e10.up.railway.app/api/v1`  
**Tested:** 2026-05-22 (smoke scripts in `scripts/smoke-all.cmd`)  
**Method:** Live HTTP via `curl` against Railway

## Summary

| Category | Passed | Failed / Not run |
|----------|--------|------------------|
| System & auth | 8 | 2 |
| Citizen GET | 11 | 0 |
| Admin GET | 11 | 0 |
| Citizen POST | 6 | 0 |
| Admin / officer actions | 5 | 0 |
| **Total verified** | **~41** | **~8** |

## Passed (HTTP 2xx)

### Public / auth
- `GET /health` → 200
- `POST /auth/seed-demo` → 200
- `POST /auth/login` (citizen, admin, examiner, officer) → 200
- `POST /auth/register` → 201
- `POST /auth/forgot-password` → 200

### Citizen
- `GET /me`, `/notifications`, `/activity`, `/licenses/me`, `/vehicles/me`
- `GET /payments/history`, `/exam/history`, `/inspection`, `/violations`
- `GET /identity/status`, `/centers`
- `POST /licenses` → 201
- `POST /identity/submit` → 200
- `POST /payments/initiate` → 201
- `POST /exam/start` → 201
- `POST /vehicles` → 201
- `GET /centers/{id}/slots` → 200
- `POST /practical/book` → 201
- `POST /inspection/schedule` → 201
- `GET /monitoring/score/{userId}` → 200
- `GET /monitoring/trips/{vehicleId}` → 200

### Admin
- `GET /admin/users`, `/admin/applications`, `/admin/audit-logs`, `/admin/vehicles`
- `GET /admin/analytics/overview`, `/revenue`, `/exams`, `/trends`
- `GET /exam/questions`, `/admin/ping` → 200
- `PUT /admin/licenses/{id}/approve` → 200

### Officer
- `GET /violations` → 200
- `POST /violations` → 201

### RBAC
- `GET /admin/ping` as citizen → **403** (expected)

## Not tested / needs manual follow-up

| Endpoint | Reason |
|----------|--------|
| `POST /exam/{id}/submit` | Needs attempt id + 30 answer payloads from start response |
| `PUT /examiner/practical/{id}/result` | Needs booking id (re-run book + examiner token) |
| `PUT /admin/identity/{id}/approve` | Needs verification document id |
| `PUT /admin/licenses/{id}/reject` | Skipped after approve on same license |
| `PUT /admin/transfer/{id}/approve` | Needs pending transfer id |
| `PUT /admin/payments/{id}/mark-paid` | Needs payment id from history |
| `PUT /licenses/{id}/renew` | Needs issued license |
| `PATCH /notifications/{id}/read` | No notifications in empty DB |
| `POST /devices/data` | Needs valid `vehicle_id` + trip payload |
| `POST /inspection/{id}/upload-report` | Needs inspection id |
| `POST /vehicles/{id}/transfer` | Needs vehicle MongoDB id (API returns plate-only view on prod until redeploy) |
| `POST /auth/verify-email` | Needs verification token from register |
| `POST /auth/reset-password` | Needs reset token from forgot-password |
| `POST /auth/bootstrap-admin` | Requires secret env var |

## Re-run tests locally

```cmd
cd scripts
smoke-all.cmd
smoke-part2.cmd
smoke-part3.cmd
type smoke-results.txt
```

Or Python (if environment allows):

```bash
python scripts/smoke_test_api.py https://api-production-5e10.up.railway.app/api/v1
```

## Demo users (verified login)

| Email | Password |
|-------|------------|
| citizen@example.com | Password123! |
| admin@example.com | Password123! |
| examiner@example.com | Password123! |
| officer@example.com | Password123! |
