@echo off
set B=https://api-production-5e10.up.railway.app/api/v1
set CIT=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1aWQiOiI2YTEwOTUzNTQxYmZhYzViYWVjZjQ0MWYiLCJlbWFpbCI6ImNpdGl6ZW5AZXhhbXBsZS5jb20iLCJyb2xlIjoiY2l0aXplbiIsImlzcyI6ImRyaXZpbmctYXV0aG9yaXR5IiwiZXhwIjoxNzc5NDg2NTA2LCJpYXQiOjE3Nzk0ODI5MDZ9.SZJBASuEZ5xRoTwC8k8c-rV-0LSXaJ-g_bzykE84I0E
set ADM=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1aWQiOiI2YTEwOTUzNTQxYmZhYzViYWVjZjQ0MjAiLCJlbWFpbCI6ImFkbWluQGV4YW1wbGUuY29tIiwicm9sZSI6ImFkbWluIiwiaXNzIjoiZHJpdmluZy1hdXRob3JpdHkiLCJleHAiOjE3Nzk0ODY1NDMsImlhdCI6MTc3OTQ4Mjk0M30.NBLvleujOLzK0LHASriqpoWyE-9I_joDE37A5uaybiQ
set OFF=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1aWQiOiI2YTEwOTUzNTQxYmZhYzViYWVjZjQ0MjIiLCJlbWFpbCI6Im9mZmljZXJAZXhhbXBsZS5jb20iLCJyb2xlIjoib2ZmaWNlciIsImlzcyI6ImRyaXZpbmctYXV0aG9yaXR5IiwiZXhwIjoxNzc5NDg2NTIxLCJpYXQiOjE3Nzk0ODI5MjF9.k6AD6i1dIZloLrkJPv7jmkGrKTzEwQdW2kbWHJtD1Ns
set R=%~dp0smoke-results.txt
set CID=6a0cc1436f69f9a38d354a99

echo. >> "%R%"
echo === PART 2 === >> "%R%"
curl.exe -s -o NUL -w "GET /centers/{id}/slots: %%{http_code}\n" -H "Authorization: Bearer %CIT%" "%B%/centers/%CID%/slots" >> "%R%"
curl.exe -s "%B%/centers/%CID%/slots" -H "Authorization: Bearer %CIT%" -o "%~dp0slots.json"
curl.exe -s -o NUL -w "POST /violations: %%{http_code}\n" -X POST -H "Authorization: Bearer %OFF%" -H "Content-Type: application/json" --data-binary "@%~dp0body-violation.json" "%B%/violations" >> "%R%"
curl.exe -s -o NUL -w "POST /auth/register: %%{http_code}\n" -X POST -H "Content-Type: application/json" --data-binary "@%~dp0body-register.json" "%B%/auth/register" >> "%R%"

curl.exe -s "%B%/vehicles/me" -H "Authorization: Bearer %CIT%" -o "%~dp0veh.json"
curl.exe -s "%B%/licenses/me" -H "Authorization: Bearer %CIT%" -o "%~dp0lic.json"

echo part2 done >> "%R%"
type "%R%"
