# Lightweight API smoke test (one request at a time, low memory)
param([string]$Base = "https://api-production-5e10.up.railway.app/api/v1")
$out = Join-Path $PSScriptRoot "smoke-results.txt"
"" | Set-Content $out -Encoding utf8

function Log($line) { Add-Content $out $line -Encoding utf8; Write-Host $line }

function Api($Method, $Path, $Token, $JsonBody) {
  $args = @("-s", "-w", "`n%{http_code}", "-X", $Method, "$Base$Path", "-H", "Content-Type: application/json")
  if ($Token) { $args += @("-H", "Authorization: Bearer $Token") }
  $tmp = [System.IO.Path]::GetTempFileName()
  if ($JsonBody) { $JsonBody | Set-Content $tmp -Encoding utf8 -NoNewline; $args += @("-d", "@$tmp") }
  $raw = & curl.exe @args 2>$null
  if ($JsonBody) { Remove-Item $tmp -Force -ErrorAction SilentlyContinue }
  $parts = $raw -split "`n"
  $code = [int]($parts[-1] -replace '\D','')
  $body = ($parts[0..($parts.Length-2)] -join "`n")
  return @{ Code = $code; Body = $body }
}

Log "=== API Smoke Test $(Get-Date -Format o) ==="
Log "Base: $Base`n"

$r = Api GET "/health" $null $null
Log ("GET /health -> {0} {1}" -f $r.Code, $(if($r.Code -eq 200){'PASS'}else{'FAIL'}))

$r = Api POST "/auth/seed-demo" $null $null
Log ("POST /auth/seed-demo -> {0} {1}" -f $r.Code, $(if($r.Code -eq 200){'PASS'}else{'FAIL'}))

function Get-Token($email) {
  $j = "{`"email`":`"$email`",`"password`":`"Password123!`"}"
  $r = Api POST "/auth/login" $null $j
  Log ("POST /auth/login ($email) -> {0}" -f $r.Code)
  if ($r.Code -eq 200) { return ($r.Body | ConvertFrom-Json).access_token }
  return $null
}

$citizen = Get-Token "citizen@example.com"
$admin = Get-Token "admin@example.com"
$examiner = Get-Token "examiner@example.com"
$officer = Get-Token "officer@example.com"

$reads = @(
  @("GET","/me",$citizen),
  @("GET","/notifications",$citizen),
  @("GET","/activity",$citizen),
  @("GET","/licenses/me",$citizen),
  @("GET","/vehicles/me",$citizen),
  @("GET","/payments/history",$citizen),
  @("GET","/exam/history",$citizen),
  @("GET","/inspection",$citizen),
  @("GET","/violations",$citizen),
  @("GET","/identity/status",$citizen),
  @("GET","/centers",$citizen),
  @("GET","/admin/users",$admin),
  @("GET","/admin/applications",$admin),
  @("GET","/admin/audit-logs",$admin),
  @("GET","/admin/vehicles",$admin),
  @("GET","/admin/analytics/overview",$admin),
  @("GET","/admin/analytics/revenue",$admin),
  @("GET","/admin/analytics/exams",$admin),
  @("GET","/admin/analytics/trends",$admin),
  @("GET","/exam/questions",$admin),
  @("GET","/admin/ping",$admin),
  @("GET","/violations",$officer)
)

foreach ($x in $reads) {
  $r = Api $x[0] $x[1] $x[2] $null
  $ok = if ($r.Code -ge 200 -and $r.Code -lt 300) { 'PASS' } else { 'FAIL' }
  Log ("{0} {1} -> {2} {3}" -f $x[0], $x[1], $r.Code, $ok)
}

# RBAC: citizen on admin ping should be 403
$r = Api GET "/admin/ping" $citizen $null
Log ("GET /admin/ping (citizen) -> {0} {1} (expect 403)" -f $r.Code, $(if($r.Code -eq 403){'PASS'}else{'FAIL'}))

# Writes
$licJson = '{"name":"Smoke Test","dob":"1998-01-01","gender":"Male","nationality":"X","address":"1 St","city":"C","postal":"1","license_type":"Car license"}'
$r = Api POST "/licenses" $citizen $licJson
Log ("POST /licenses -> {0} {1}" -f $r.Code, $(if($r.Code -in 200,201){'PASS'}else{'FAIL'}))

$idJson = '{"national_id_number":"SMK001","document_front_path":"/f.pdf","document_back_path":"/b.pdf","selfie_path":"/s.jpg"}'
$r = Api POST "/identity/submit" $citizen $idJson
Log ("POST /identity/submit -> {0} {1}" -f $r.Code, $(if($r.Code -eq 200){'PASS'}else{'FAIL'}))

$payJson = '{"service_type":"license"}'
$r = Api POST "/payments/initiate" $citizen $payJson
Log ("POST /payments/initiate -> {0} {1}" -f $r.Code, $(if($r.Code -in 200,201){'PASS'}else{'FAIL'}))

$examJson = '{"license_type":"car"}'
$r = Api POST "/exam/start" $citizen $examJson
Log ("POST /exam/start -> {0} {1}" -f $r.Code, $(if($r.Code -in 200,201){'PASS'}else{'FAIL'}))

$vin = "SMK" + (Get-Random -Maximum 99999999)
$vehJson = "{`"vin`":`"$vin`",`"plate`":`"P-$vin`",`"make`":`"Toyota`",`"model`":`"Yaris`",`"year`":2021,`"color`":`"Red`"}"
$r = Api POST "/vehicles" $citizen $vehJson
Log ("POST /vehicles -> {0} {1}" -f $r.Code, $(if($r.Code -in 200,201){'PASS'}else{'FAIL'}))

$fpJson = '{"email":"citizen@example.com"}'
$r = Api POST "/auth/forgot-password" $null $fpJson
Log ("POST /auth/forgot-password -> {0} {1}" -f $r.Code, $(if($r.Code -eq 200){'PASS'}else{'FAIL'}))

Log "`nDone. Full log: $out"
