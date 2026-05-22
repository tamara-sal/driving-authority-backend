@echo off
set B=https://api-production-5e10.up.railway.app/api/v1
set C=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1aWQiOiI2YTEwOTUzNTQxYmZhYzViYWVjZjQ0MWYiLCJlbWFpbCI6ImNpdGl6ZW5AZXhhbXBsZS5jb20iLCJyb2xlIjoiY2l0aXplbiIsImlzcyI6ImRyaXZpbmctYXV0aG9yaXR5IiwiZXhwIjoxNzc5NDg2NTA2LCJpYXQiOjE3Nzk0ODI5MDZ9.SZJBASuEZ5xRoTwC8k8c-rV-0LSXaJ-g_bzykE84I0E
set O=%~dp0smoke-results.txt
echo CITIZEN GETS >> "%O%"
curl.exe -s -o NUL -w "GET /me: %%{http_code}\n" -H "Authorization: Bearer %C%" "%B%/me"
curl.exe -s -o NUL -w "GET /notifications: %%{http_code}\n" -H "Authorization: Bearer %C%" "%B%/notifications"
curl.exe -s -o NUL -w "GET /activity: %%{http_code}\n" -H "Authorization: Bearer %C%" "%B%/activity"
curl.exe -s -o NUL -w "GET /licenses/me: %%{http_code}\n" -H "Authorization: Bearer %C%" "%B%/licenses/me"
curl.exe -s -o NUL -w "GET /vehicles/me: %%{http_code}\n" -H "Authorization: Bearer %C%" "%B%/vehicles/me"
curl.exe -s -o NUL -w "GET /payments/history: %%{http_code}\n" -H "Authorization: Bearer %C%" "%B%/payments/history"
curl.exe -s -o NUL -w "GET /exam/history: %%{http_code}\n" -H "Authorization: Bearer %C%" "%B%/exam/history"
curl.exe -s -o NUL -w "GET /inspection: %%{http_code}\n" -H "Authorization: Bearer %C%" "%B%/inspection"
curl.exe -s -o NUL -w "GET /violations: %%{http_code}\n" -H "Authorization: Bearer %C%" "%B%/violations"
curl.exe -s -o NUL -w "GET /identity/status: %%{http_code}\n" -H "Authorization: Bearer %C%" "%B%/identity/status"
curl.exe -s -o NUL -w "GET /centers: %%{http_code}\n" -H "Authorization: Bearer %C%" "%B%/centers"
