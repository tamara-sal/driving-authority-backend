# Frontend (fyppp zip) ↔ Backend API alignment

The React app in `fyppp (2).zip` is a **UI prototype** that originally used `mockData.js` and `fakeApi.js`. The Go backend exposes real REST endpoints under `/api/v1` that cover every screen in the app.

## Demo users (required for login)

| Role | Email | Password |
|------|-------|----------|
| Citizen | citizen@example.com | Password123! |
| Admin | admin@example.com | Password123! |
| Examiner | examiner@example.com | Password123! |
| Officer | officer@example.com | Password123! |

Seeded on startup (`SEED_DEMO_USERS=true`) or via `POST /api/v1/auth/seed-demo`.

## Screen → API mapping

| Frontend route / page | Backend endpoint(s) | Status |
|----------------------|---------------------|--------|
| Register | `POST /auth/register` (`full_name`, email, password, phone) | ✅ |
| Login | `POST /auth/login` → JWT + `name`, `redirect` | ✅ |
| Verify email | `POST /auth/verify-email` | ✅ |
| Forgot password | `POST /auth/forgot-password`, `POST /auth/reset-password` | ✅ |
| Profile / me | `GET /me` | ✅ |
| Apply for license | `POST /licenses` (application fields + `license_type`) | ✅ |
| License status | `GET /licenses/me` | ✅ |
| Theory exam | `POST /exam/start`, `POST /exam/:id/submit`, `GET /exam/history` | ✅ |
| Exam questions (admin) | `GET /exam/questions` | ✅ |
| Practical booking | `GET /centers`, `GET /centers/:id/slots`, `POST /practical/book` | ✅ |
| Submit practical result | `PUT /examiner/practical/:id/result` | ✅ |
| Vehicles list / register | `GET /vehicles/me`, `POST /vehicles` (`plate` or `plate_number`) | ✅ |
| Transfer vehicle | `POST /vehicles/:id/transfer` (`buyer_email` or `buyer_id`) | ✅ |
| Inspections / history | `POST /inspection/schedule`, `GET /inspection` | ✅ |
| Payments | `POST /payments/initiate`, `GET /payments/history` | ✅ |
| Monitoring | `GET /monitoring/trips/:vehicleId`, `GET /monitoring/score/:userId` | ✅ |
| Notifications | `GET /notifications`, `PATCH /notifications/:id/read` | ✅ |
| Activity feed | `GET /activity` | ✅ |
| Admin users | `GET /admin/users` | ✅ |
| License applications | `GET /admin/applications`, `PUT /admin/licenses/:id/approve`, `PUT .../reject` | ✅ |
| Identity verify | `POST /identity/submit`, `GET /identity/status`, approve/reject by id | ✅ |
| Violations | `GET /violations`, `POST /violations`, `PUT /violations/:id/status` | ✅ |
| Admin vehicles | `GET /admin/vehicles` | ✅ |
| Analytics | `GET /admin/analytics/overview`, `/revenue`, `/exams`, `/trends` | ✅ |
| Audit logs | `GET /admin/audit-logs` | ✅ |

## Frontend integration

Point the Vite app at the API:

```env
VITE_API_URL=http://localhost:8080/api/v1
```

Or use the dev proxy in `vite.config.js` (`/api` → `http://localhost:8080`).

Reference client: `frontend-extracted/fyppp/src/lib/api.js` (auth, tables, license apply).

## UI-only (no separate API)

These are simulated in the UI only; backend has related data via other endpoints:

- File upload previews (documents stored as paths on `POST /identity/submit`)
- System settings form (no settings collection yet)
- Hard-coded dashboard stat cards (use `/admin/analytics/overview` for live counts)
