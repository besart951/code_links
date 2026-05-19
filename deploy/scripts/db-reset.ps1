Set-StrictMode -Version Latest
$ErrorActionPreference = "Stop"

if ($env:CONFIRM_RESET -ne "codelinks") {
    Write-Error "Refusing to reset database volumes. Re-run with CONFIRM_RESET=codelinks."
}

$repoRoot = Resolve-Path (Join-Path $PSScriptRoot "..\..")
$composeFile = Join-Path $repoRoot "deploy\docker\docker-compose.dev.yml"

Push-Location $repoRoot
try {
    docker compose -f $composeFile down -v
    docker compose -f $composeFile up --build postgres platform-migrate infra-migrate
}
finally {
    Pop-Location
}
