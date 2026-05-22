@echo off
set BASE=https://api-production-5e10.up.railway.app/api/v1
set OUT=%~dp0smoke-results.txt
echo === Smoke Test %date% %time% === > "%OUT%"
echo Base: %BASE% >> "%OUT%"
echo. >> "%OUT%"

call :test GET /health
curl.exe -s -X POST "%BASE%/auth/seed-demo" >> "%OUT%" 2>&1
echo POST /auth/seed-demo >> "%OUT%"

curl.exe -s -X POST "%BASE%/auth/login" -H "Content-Type: application/json" -d "{\"email\":\"citizen@example.com\",\"password\":\"Password123!\"}" -o "%TEMP%\cit.json"
curl.exe -s -X POST "%BASE%/auth/login" -H "Content-Type: application/json" -d "{\"email\":\"admin@example.com\",\"password\":\"Password123!\"}" -o "%TEMP%\adm.json"
curl.exe -s -X POST "%BASE%/auth/login" -H "Content-Type: application/json" -d "{\"email\":\"examiner@example.com\",\"password\":\"Password123!\"}" -o "%TEMP%\exm.json"
curl.exe -s -X POST "%BASE%/auth/login" -H "Content-Type: application/json" -d "{\"email\":\"officer@example.com\",\"password\":\"Password123!\"}" -o "%TEMP%\off.json"

for %%f in (cit adm exm off) do (
  for /f "delims=" %%t in ('powershell -NoProfile -Command "(Get-Content $env:TEMP\%%f.json | ConvertFrom-Json).access_token" 2^>nul') do set TOKEN_%%f=%%t
)

echo. >> "%OUT%"
echo --- Citizen reads --- >> "%OUT%"
call :auth GET /me TOKEN_cit
call :auth GET /notifications TOKEN_cit
call :auth GET /activity TOKEN_cit
call :auth GET /licenses/me TOKEN_cit
call :auth GET /vehicles/me TOKEN_cit
call :auth GET /payments/history TOKEN_cit
call :auth GET /exam/history TOKEN_cit
call :auth GET /inspection TOKEN_cit
call :auth GET /violations TOKEN_cit
call :auth GET /identity/status TOKEN_cit
call :auth GET /centers TOKEN_cit

echo. >> "%OUT%"
echo --- Admin reads --- >> "%OUT%"
call :auth GET /admin/users TOKEN_adm
call :auth GET /admin/applications TOKEN_adm
call :auth GET /admin/audit-logs TOKEN_adm
call :auth GET /admin/vehicles TOKEN_adm
call :auth GET /admin/analytics/overview TOKEN_adm
call :auth GET /admin/analytics/revenue TOKEN_adm
call :auth GET /admin/analytics/exams TOKEN_adm
call :auth GET /admin/analytics/trends TOKEN_adm
call :auth GET /exam/questions TOKEN_adm
call :auth GET /admin/ping TOKEN_adm

echo. >> "%OUT%"
echo --- Officer --- >> "%OUT%"
call :auth GET /violations TOKEN_off

echo. >> "%OUT%"
echo --- Citizen writes --- >> "%OUT%"
call :auth POST /licenses TOKEN_cit "{\"name\":\"T\",\"dob\":\"1998-01-01\",\"gender\":\"M\",\"nationality\":\"X\",\"address\":\"1\",\"city\":\"C\",\"postal\":\"1\",\"license_type\":\"Car license\"}"
call :auth POST /identity/submit TOKEN_cit "{\"national_id_number\":\"X1\",\"document_front_path\":\"/a\",\"document_back_path\":\"/b\",\"selfie_path\":\"/c\"}"
call :auth POST /payments/initiate TOKEN_cit "{\"service_type\":\"license\"}"
call :auth POST /exam/start TOKEN_cit "{\"license_type\":\"car\"}"
call :auth POST /vehicles TOKEN_cit "{\"vin\":\"BAT12345678901\",\"plate\":\"BAT-99\",\"make\":\"Toyota\",\"model\":\"X\",\"year\":2020,\"color\":\"Blue\"}"

echo. >> "%OUT%"
echo --- Public --- >> "%OUT%"
call :test POST /auth/forgot-password "{\"email\":\"citizen@example.com\"}"

echo. >> "%OUT%"
echo DONE >> "%OUT%"
type "%OUT%"
exit /b 0

:test
curl.exe -s -o NUL -w "%%{http_code}" -X %1 "%BASE%%2" %3 %4 %5 %6
echo %1 %2 HTTP:%%~3 >> "%OUT%"
exit /b 0

:auth
set TOK=!%3!
curl.exe -s -w " HTTP:%%{http_code}" -X %1 "%BASE%%2" -H "Authorization: Bearer !TOK!" -H "Content-Type: application/json" %4 %5 %6 %7 %8 >> "%OUT%" 2>&1
echo %1 %2 >> "%OUT%"
exit /b 0
