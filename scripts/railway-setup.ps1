# Railway setup: project, MongoDB, GitHub deploy, env vars, public domain
# Prerequisites: railway login  OR  set $env:RAILWAY_TOKEN = "your-token"
# Run: .\scripts\railway-setup.ps1

$ErrorActionPreference = "Stop"
Set-Location (Resolve-Path (Join-Path $PSScriptRoot ".."))

$ProjectName = "driving-authority-backend"
$GitHubRepo = "tamara-sal/driving-authority-backend"
$ApiService = "api"

Write-Host "==> Checking Railway auth..." -ForegroundColor Cyan
railway whoami 2>&1 | Out-Null
if ($LASTEXITCODE -ne 0) {
    Write-Host "Not logged in. Run: railway login" -ForegroundColor Red
    Write-Host "Or set: `$env:RAILWAY_TOKEN = 'your-project-token'" -ForegroundColor Yellow
    exit 1
}

if (-not (Test-Path ".railway")) {
    Write-Host "==> Creating Railway project '$ProjectName'..." -ForegroundColor Cyan
    railway init --name $ProjectName --json | Out-Null
} else {
    Write-Host "==> Project already linked." -ForegroundColor Green
}

Write-Host "==> Adding MongoDB..." -ForegroundColor Cyan
$dbJson = railway add --database mongo --json 2>&1
Write-Host $dbJson

Write-Host "==> Linking GitHub repo ($GitHubRepo)..." -ForegroundColor Cyan
$repoJson = railway add --repo $GitHubRepo --service $ApiService --json 2>&1
Write-Host $repoJson

$bytes = New-Object byte[] 32
[System.Security.Cryptography.RandomNumberGenerator]::Create().GetBytes($bytes)
$jwtSecret = [BitConverter]::ToString($bytes).Replace("-", "").ToLower()

Write-Host "==> Setting API environment variables..." -ForegroundColor Cyan
railway variable set MONGO_URI='${{MongoDB.MONGO_URL}}' --service $ApiService --skip-deploys
railway variable set MONGO_DB=driving_authority --service $ApiService --skip-deploys
railway variable set JWT_SECRET=$jwtSecret --service $ApiService --skip-deploys
railway variable set JWT_ISSUER=driving-authority --service $ApiService --skip-deploys
railway variable set JWT_ACCESS_TTL_MINUTES=60 --service $ApiService --skip-deploys
railway variable set APP_ENV=production --service $ApiService --skip-deploys

Write-Host "==> Generating public domain..." -ForegroundColor Cyan
$domainJson = railway domain --service $ApiService --json 2>&1
Write-Host $domainJson

Write-Host "==> Triggering deploy..." -ForegroundColor Cyan
railway up --service $ApiService

Write-Host @"

Done. Check status: railway status
Logs: railway logs --service $ApiService
Dashboard: railway open

Save JWT_SECRET (shown once): $jwtSecret
"@ -ForegroundColor Green
