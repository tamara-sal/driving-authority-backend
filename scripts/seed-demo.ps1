# Seed the 4 demo users (citizen, admin, examiner, officer) into the running API.
# Usage: .\scripts\seed-demo.ps1
#        .\scripts\seed-demo.ps1 -BaseUrl "https://your-api.up.railway.app"

param(
    [string]$BaseUrl = "http://localhost:8080"
)

$ErrorActionPreference = "Stop"
$url = "$BaseUrl/api/v1/auth/seed-demo"

Write-Host "Seeding demo users at $url ..." -ForegroundColor Cyan
$response = Invoke-RestMethod -Method POST -Uri $url -ContentType "application/json"
$response | ConvertTo-Json -Depth 5
Write-Host "`nLogin with any account above and password: $($response.password)" -ForegroundColor Green
