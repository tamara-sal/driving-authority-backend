# Push to GitHub and deploy on Railway
# Run from repo root: .\scripts\deploy.ps1

$ErrorActionPreference = "Stop"
Set-Location (Resolve-Path (Join-Path $PSScriptRoot ".."))

Write-Host "==> Checking GitHub CLI..." -ForegroundColor Cyan
gh auth status 2>$null
if ($LASTEXITCODE -ne 0) {
    Write-Host "Log in to GitHub (browser will open):" -ForegroundColor Yellow
    gh auth login --hostname github.com --git-protocol https --web
}

$repoName = "driving-authority-backend"
$ghUser = (gh api user -q .login)
$originUrl = "https://github.com/$ghUser/$repoName.git"

Write-Host "==> Creating GitHub repo $ghUser/$repoName (if needed)..." -ForegroundColor Cyan
gh repo view "$ghUser/$repoName" 2>$null | Out-Null
if ($LASTEXITCODE -ne 0) {
    gh repo create $repoName --public --source=. --remote=origin --push
} else {
    $remotes = @(git remote)
    if ($remotes -contains "origin") {
        git remote set-url origin $originUrl
    } else {
        git remote add origin $originUrl
    }
    git push -u origin main
}

Write-Host "==> GitHub: https://github.com/$ghUser/$repoName" -ForegroundColor Green

Write-Host "`n==> Railway deploy..." -ForegroundColor Cyan
railway whoami 2>$null | Out-Null
if ($LASTEXITCODE -ne 0) {
    Write-Host "Log in to Railway:" -ForegroundColor Yellow
    railway login
}

if (-not (Test-Path ".railway")) {
    railway init
}

Write-Host @"

Next steps in Railway dashboard (https://railway.com):
1. Add a MongoDB service to your project
2. On the API service, set variables:
   MONGO_URI     = (reference from MongoDB service, e.g. `${{MongoDB.MONGO_URL}})
   MONGO_DB      = driving_authority
   JWT_SECRET    = (long random string)
   JWT_ISSUER    = driving-authority
   APP_ENV       = production
3. Generate a public domain for the API service

Or run: railway up
"@ -ForegroundColor Yellow

$deploy = Read-Host "Deploy now with 'railway up'? (y/N)"
if ($deploy -eq "y") { railway up }
