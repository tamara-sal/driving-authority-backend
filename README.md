## Digital Driving Authority System (Backend)

Go + MongoDB backend starter for:
- Auth (register/login) + JWT
- Role-based access control (Citizen/Admin/Examiner/Officer)
- Identity verification workflow (submit + admin approve/reject)

### Run locally

1) Start MongoDB:

```bash
docker compose up -d
```

2) Create `.env` from example:

```bash
copy .env.example .env
```

3) Run API:

```bash
go mod tidy
go run .\cmd\api
```

API base: `http://localhost:8080/api/v1`

### Deploy to Railway

1. Push this repo to GitHub (see below).
2. In [Railway](https://railway.com), create a project → **Deploy from GitHub repo** → select this repository.
3. Add a **MongoDB** service (Railway plugin) and link it to the API service.
4. Set these variables on the API service:

| Variable | Notes |
|----------|--------|
| `MONGO_URI` | From the MongoDB service (`MONGO_URL` or connection string variable) |
| `MONGO_DB` | `driving_authority` |
| `JWT_SECRET` | Long random secret (required) |
| `JWT_ISSUER` | `driving-authority` |
| `JWT_ACCESS_TTL_MINUTES` | `60` |
| `APP_ENV` | `production` |
| `BOOTSTRAP_ADMIN_SECRET` | Optional; for one-time admin bootstrap |

Railway sets `PORT` automatically. Swagger UI: `https://<your-app>.up.railway.app/swagger/index.html`

### Swagger (API docs for frontend)

The OpenAPI spec lives in [`docs/swagger.json`](docs/swagger.json). Send that file (or the URLs below) to your frontend developer.

Run the API and open Swagger UI:

- UI: `http://localhost:8080/swagger/index.html`
- OpenAPI JSON: `http://localhost:8080/swagger/doc.json`

In Swagger UI, click **Authorize** and enter `Bearer <access_token>` (from login/register).

### API modules (from system spec)

Full endpoint list: [`scripts/generate-swagger-paths.md`](scripts/generate-swagger-paths.md)

| Module | Highlights |
|--------|------------|
| Auth | register, login, verify-email, forgot/reset password |
| Identity | submit, status, admin approve/reject |
| Licenses | apply, approve, renew |
| Theory exam | 30 MCQs, start/submit, history |
| Practical | centers, slots, book, examiner result |
| Vehicles | register, transfer, admin approve |
| Inspection | schedule, upload report path |
| Monitoring | device data, trips, safety score |
| Payments | initiate, history, admin mark-paid |
| Analytics | admin overview, revenue, exams |

Roles: `citizen`, `admin`, `examiner`, `officer` (RBAC on routes).

