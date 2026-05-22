@echo off
set B=https://api-production-5e10.up.railway.app/api/v1
set CIT=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1aWQiOiI2YTEwOTUzNTQxYmZhYzViYWVjZjQ0MWYiLCJlbWFpbCI6ImNpdGl6ZW5AZXhhbXBsZS5jb20iLCJyb2xlIjoiY2l0aXplbiIsImlzcyI6ImRyaXZpbmctYXV0aG9yaXR5IiwiZXhwIjoxNzc5NDg2NTA2LCJpYXQiOjE3Nzk0ODI5MDZ9.SZJBASuEZ5xRoTwC8k8c-rV-0LSXaJ-g_bzykE84I0E
set ADM=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1aWQiOiI2YTEwOTUzNTQxYmZhYzViYWVjZjQ0MjAiLCJlbWFpbCI6ImFkbWluQGV4YW1wbGUuY29tIiwicm9sZSI6ImFkbWluIiwiaXNzIjoiZHJpdmluZy1hdXRob3JpdHkiLCJleHAiOjE3Nzk0ODY1NDMsImlhdCI6MTc3OTQ4Mjk0M30.NBLvleujOLzK0LHASriqpoWyE-9I_joDE37A5uaybiQ
set EXM=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1aWQiOiI2YTEwOTUzNTQxYmZhYzViYWVjZjQ0MjEiLCJlbWFpbCI6ImV4YW1pbmVyQGV4YW1wbGUuY29tIiwicm9sZSI6ImV4YW1pbmVyIiwiaXNzIjoiZHJpdmluZy1hdXRob3JpdHkiLCJleHAiOjE3Nzk0ODY1NDYsImlhdCI6MTc3OTQ4Mjk0Nn0.DilFTtwSCFticpHcAPdo-1ESmWCoUFkuicWkyQhs_0Y
set R=%~dp0smoke-results.txt
set LIC=6a10c17ef4b4ac8dfa2d4021
set SLOT=6a0cc1436f69f9a38d354a9b
set UID=6a10953541bfac5baecf441f

echo. >> "%R%"
echo === PART 3 === >> "%R%"
curl.exe -s -o NUL -w "POST /practical/book: %%{http_code}\n" -X POST -H "Authorization: Bearer %CIT%" -H "Content-Type: application/json" -d "{\"slot_id\":\"%SLOT%\"}" "%B%/practical/book" >> "%R%"
curl.exe -s "%B%/practical/book" -X POST -H "Authorization: Bearer %CIT%" -H "Content-Type: application/json" -d "{\"slot_id\":\"%SLOT%\"}" -o "%~dp0booking.json"
curl.exe -s -o NUL -w "PUT /admin/licenses/approve: %%{http_code}\n" -X PUT -H "Authorization: Bearer %ADM%" "%B%/admin/licenses/%LIC%/approve" >> "%R%"
curl.exe -s -o NUL -w "GET /monitoring/score: %%{http_code}\n" -H "Authorization: Bearer %CIT%" "%B%/monitoring/score/%UID%" >> "%R%"
curl.exe -s -o NUL -w "POST /inspection/schedule: %%{http_code}\n" -X POST -H "Authorization: Bearer %CIT%" -H "Content-Type: application/json" -d "{\"vehicle_id\":\"000000000000000000000001\",\"inspection_date\":\"2026-06-15\"}" "%B%/inspection/schedule" >> "%R%"
curl.exe -s -o NUL -w "GET /monitoring/trips: %%{http_code}\n" -H "Authorization: Bearer %CIT%" "%B%/monitoring/trips/000000000000000000000001" >> "%R%"

for /f "tokens=*" %%i in ('type "%~dp0booking.json" ^| findstr /i "id"') do echo booking raw >> "%R%"
echo part3 done >> "%R%"
type "%R%"
