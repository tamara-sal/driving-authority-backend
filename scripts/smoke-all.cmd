@echo off
set B=https://api-production-5e10.up.railway.app/api/v1
set CIT=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1aWQiOiI2YTEwOTUzNTQxYmZhYzViYWVjZjQ0MWYiLCJlbWFpbCI6ImNpdGl6ZW5AZXhhbXBsZS5jb20iLCJyb2xlIjoiY2l0aXplbiIsImlzcyI6ImRyaXZpbmctYXV0aG9yaXR5IiwiZXhwIjoxNzc5NDg2NTA2LCJpYXQiOjE3Nzk0ODI5MDZ9.SZJBASuEZ5xRoTwC8k8c-rV-0LSXaJ-g_bzykE84I0E
set ADM=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1aWQiOiI2YTEwOTUzNTQxYmZhYzViYWVjZjQ0MjAiLCJlbWFpbCI6ImFkbWluQGV4YW1wbGUuY29tIiwicm9sZSI6ImFkbWluIiwiaXNzIjoiZHJpdmluZy1hdXRob3JpdHkiLCJleHAiOjE3Nzk0ODY1NDMsImlhdCI6MTc3OTQ4Mjk0M30.NBLvleujOLzK0LHASriqpoWyE-9I_joDE37A5uaybiQ
set EXM=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1aWQiOiI2YTEwOTUzNTQxYmZhYzViYWVjZjQ0MjEiLCJlbWFpbCI6ImV4YW1pbmVyQGV4YW1wbGUuY29tIiwicm9sZSI6ImV4YW1pbmVyIiwiaXNzIjoiZHJpdmluZy1hdXRob3JpdHkiLCJleHAiOjE3Nzk0ODY1NDYsImlhdCI6MTc3OTQ4Mjk0Nn0.DilFTtwSCFticpHcAPdo-1ESmWCoUFkuicWkyQhs_0Y
set OFF=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1aWQiOiI2YTEwOTUzNTQxYmZhYzViYWVjZjQ0MjIiLCJlbWFpbCI6Im9mZmljZXJAZXhhbXBsZS5jb20iLCJyb2xlIjoib2ZmaWNlciIsImlzcyI6ImRyaXZpbmctYXV0aG9yaXR5IiwiZXhwIjoxNzc5NDg2NTIxLCJpYXQiOjE3Nzk0ODI5MjF9.k6AD6i1dIZloLrkJPv7jmkGrKTzEwQdW2kbWHJtD1Ns
set R=%~dp0smoke-results.txt

echo API SMOKE TEST RESULTS > "%R%"
echo Base: %B% >> "%R%"
echo. >> "%R%"

curl.exe -s -o NUL -w "GET /health: %%{http_code} PASS?{}\n" "%B%/health" >> "%R%"
curl.exe -s -o NUL -w "POST /auth/seed-demo: %%{http_code}\n" -X POST "%B%/auth/seed-demo" >> "%R%"
curl.exe -s -o NUL -w "POST /auth/login citizen: %%{http_code}\n" -X POST "%B%/auth/login" -H "Content-Type: application/json" --data-binary "@%~dp0login-citizen.json" >> "%R%"

echo. >> "%R%"
echo === CITIZEN GET === >> "%R%"
curl.exe -s -o NUL -w "GET /me: %%{http_code}\n" -H "Authorization: Bearer %CIT%" "%B%/me" >> "%R%"
curl.exe -s -o NUL -w "GET /notifications: %%{http_code}\n" -H "Authorization: Bearer %CIT%" "%B%/notifications" >> "%R%"
curl.exe -s -o NUL -w "GET /activity: %%{http_code}\n" -H "Authorization: Bearer %CIT%" "%B%/activity" >> "%R%"
curl.exe -s -o NUL -w "GET /licenses/me: %%{http_code}\n" -H "Authorization: Bearer %CIT%" "%B%/licenses/me" >> "%R%"
curl.exe -s -o NUL -w "GET /vehicles/me: %%{http_code}\n" -H "Authorization: Bearer %CIT%" "%B%/vehicles/me" >> "%R%"
curl.exe -s -o NUL -w "GET /payments/history: %%{http_code}\n" -H "Authorization: Bearer %CIT%" "%B%/payments/history" >> "%R%"
curl.exe -s -o NUL -w "GET /exam/history: %%{http_code}\n" -H "Authorization: Bearer %CIT%" "%B%/exam/history" >> "%R%"
curl.exe -s -o NUL -w "GET /inspection: %%{http_code}\n" -H "Authorization: Bearer %CIT%" "%B%/inspection" >> "%R%"
curl.exe -s -o NUL -w "GET /violations: %%{http_code}\n" -H "Authorization: Bearer %CIT%" "%B%/violations" >> "%R%"
curl.exe -s -o NUL -w "GET /identity/status: %%{http_code}\n" -H "Authorization: Bearer %CIT%" "%B%/identity/status" >> "%R%"
curl.exe -s -o NUL -w "GET /centers: %%{http_code}\n" -H "Authorization: Bearer %CIT%" "%B%/centers" >> "%R%"

echo. >> "%R%"
echo === ADMIN GET === >> "%R%"
curl.exe -s -o NUL -w "GET /admin/users: %%{http_code}\n" -H "Authorization: Bearer %ADM%" "%B%/admin/users" >> "%R%"
curl.exe -s -o NUL -w "GET /admin/applications: %%{http_code}\n" -H "Authorization: Bearer %ADM%" "%B%/admin/applications" >> "%R%"
curl.exe -s -o NUL -w "GET /admin/audit-logs: %%{http_code}\n" -H "Authorization: Bearer %ADM%" "%B%/admin/audit-logs" >> "%R%"
curl.exe -s -o NUL -w "GET /admin/vehicles: %%{http_code}\n" -H "Authorization: Bearer %ADM%" "%B%/admin/vehicles" >> "%R%"
curl.exe -s -o NUL -w "GET /admin/analytics/overview: %%{http_code}\n" -H "Authorization: Bearer %ADM%" "%B%/admin/analytics/overview" >> "%R%"
curl.exe -s -o NUL -w "GET /admin/analytics/revenue: %%{http_code}\n" -H "Authorization: Bearer %ADM%" "%B%/admin/analytics/revenue" >> "%R%"
curl.exe -s -o NUL -w "GET /admin/analytics/exams: %%{http_code}\n" -H "Authorization: Bearer %ADM%" "%B%/admin/analytics/exams" >> "%R%"
curl.exe -s -o NUL -w "GET /admin/analytics/trends: %%{http_code}\n" -H "Authorization: Bearer %ADM%" "%B%/admin/analytics/trends" >> "%R%"
curl.exe -s -o NUL -w "GET /exam/questions: %%{http_code}\n" -H "Authorization: Bearer %ADM%" "%B%/exam/questions" >> "%R%"
curl.exe -s -o NUL -w "GET /admin/ping: %%{http_code}\n" -H "Authorization: Bearer %ADM%" "%B%/admin/ping" >> "%R%"
curl.exe -s -o NUL -w "GET /admin/ping citizen expect 403: %%{http_code}\n" -H "Authorization: Bearer %CIT%" "%B%/admin/ping" >> "%R%"

echo. >> "%R%"
echo === OFFICER === >> "%R%"
curl.exe -s -o NUL -w "GET /violations officer: %%{http_code}\n" -H "Authorization: Bearer %OFF%" "%B%/violations" >> "%R%"

echo. >> "%R%"
echo === CITIZEN POST === >> "%R%"
curl.exe -s -o NUL -w "POST /licenses: %%{http_code}\n" -X POST -H "Authorization: Bearer %CIT%" -H "Content-Type: application/json" --data-binary "@%~dp0body-license.json" "%B%/licenses" >> "%R%"
curl.exe -s -o NUL -w "POST /identity/submit: %%{http_code}\n" -X POST -H "Authorization: Bearer %CIT%" -H "Content-Type: application/json" --data-binary "@%~dp0body-identity.json" "%B%/identity/submit" >> "%R%"
curl.exe -s -o NUL -w "POST /payments/initiate: %%{http_code}\n" -X POST -H "Authorization: Bearer %CIT%" -H "Content-Type: application/json" --data-binary "@%~dp0body-payment.json" "%B%/payments/initiate" >> "%R%"
curl.exe -s -o NUL -w "POST /exam/start: %%{http_code}\n" -X POST -H "Authorization: Bearer %CIT%" -H "Content-Type: application/json" --data-binary "@%~dp0body-exam-start.json" "%B%/exam/start" >> "%R%"
curl.exe -s -o NUL -w "POST /vehicles: %%{http_code}\n" -X POST -H "Authorization: Bearer %CIT%" -H "Content-Type: application/json" --data-binary "@%~dp0body-vehicle.json" "%B%/vehicles" >> "%R%"
curl.exe -s -o NUL -w "POST /auth/forgot-password: %%{http_code}\n" -X POST -H "Content-Type: application/json" --data-binary "@%~dp0body-forgot.json" "%B%/auth/forgot-password" >> "%R%"

echo. >> "%R%"
echo === CENTERS SLOTS === >> "%R%"
curl.exe -s "%B%/centers" -H "Authorization: Bearer %CIT%" -o "%~dp0centers.json"
curl.exe -s -o NUL -w "GET /centers saved: %%{http_code}\n" -H "Authorization: Bearer %CIT%" "%B%/centers" >> "%R%"

type "%R%"
