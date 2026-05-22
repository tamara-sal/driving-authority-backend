# API handoff for frontend (fyppp)

Share this document with the React developer.

## Production base URL

```
https://api-production-5e10.up.railway.app/api/v1
```

## Swagger (interactive docs)

- UI: https://api-production-5e10.up.railway.app/swagger/index.html
- Health: https://api-production-5e10.up.railway.app/api/v1/health

## Frontend `.env`

```env
VITE_API_URL=https://api-production-5e10.up.railway.app/api/v1
```

## Authentication

All protected routes need:

```
Authorization: Bearer <access_token>
```

### Demo accounts (4 roles)

| Role | Email | Password |
|------|-------|----------|
| Citizen | `citizen@example.com` | `Password123!` |
| Admin | `admin@example.com` | `Password123!` |
| Examiner | `examiner@example.com` | `Password123!` |
| Officer | `officer@example.com` | `Password123!` |

Re-seed anytime (safe to call; upserts users):

```http
POST /api/v1/auth/seed-demo
```

### Login

```http
POST /api/v1/auth/login
Content-Type: application/json

{
  "email": "citizen@example.com",
  "password": "Password123!"
}
```

Response:

```json
{
  "access_token": "...",
  "user_id": "...",
  "email": "citizen@example.com",
  "name": "John Citizen",
  "role": "citizen",
  "redirect": "/dashboard/citizen"
}
```

Use `redirect` for React Router after login.

### Register

```http
POST /api/v1/auth/register

{
  "full_name": "John Doe",
  "email": "newuser@example.com",
  "phone": "+1234567890",
  "password": "Password123!"
}
```

## Main endpoints (by feature)

| Feature | Method | Path |
|---------|--------|------|
| Current user | GET | `/me` |
| Notifications | GET | `/notifications` |
| Activity | GET | `/activity` |
| Apply license | POST | `/licenses` |
| My licenses | GET | `/licenses/me` |
| Theory exam start | POST | `/exam/start` |
| Theory exam submit | POST | `/exam/{attemptId}/submit` |
| Exam history | GET | `/exam/history` |
| Exam questions (admin) | GET | `/exam/questions` |
| Test centers | GET | `/centers` |
| Slots | GET | `/centers/{id}/slots` |
| Book practical | POST | `/practical/book` |
| Examiner result | PUT | `/examiner/practical/{id}/result` |
| My vehicles | GET | `/vehicles/me` |
| Register vehicle | POST | `/vehicles` |
| Transfer vehicle | POST | `/vehicles/{id}/transfer` |
| Inspections list | GET | `/inspection` |
| Schedule inspection | POST | `/inspection/schedule` |
| Payments history | GET | `/payments/history` |
| Violations | GET | `/violations` |
| Update violation | PUT | `/violations/{id}/status` (use Mongo `id` from list) |
| Admin users | GET | `/admin/users` |
| License applications | GET | `/admin/applications` |
| Approve license | PUT | `/admin/licenses/{id}/approve` |
| Reject license | PUT | `/admin/licenses/{id}/reject` |
| Admin vehicles | GET | `/admin/vehicles` |
| Analytics overview | GET | `/admin/analytics/overview` |
| Analytics trends | GET | `/admin/analytics/trends` |
| Audit logs | GET | `/admin/audit-logs` |
| Identity submit | POST | `/identity/submit` |
| Identity status | GET | `/identity/status` |

Full list: [`scripts/generate-swagger-paths.md`](../scripts/generate-swagger-paths.md)

## License application body (matches Apply for License form)

```json
{
  "name": "John Citizen",
  "dob": "1998-01-15",
  "gender": "Male",
  "nationality": "Country",
  "address": "123 Main St",
  "city": "Capital",
  "postal": "10001",
  "license_type": "Car license"
}
```

`license_type` values: `Motorcycle license`, `Car license`, `Commercial license`

Response includes `reference` (e.g. `APP-2026-xxxxxxxx`).

## CORS

The API allows all origins (`*`). Browser requests from `localhost:5173` work without extra config.

## GitHub repo

https://github.com/tamara-sal/driving-authority-backend
