Set-StrictMode -Version Latest
$ErrorActionPreference = "Stop"

$repoRoot = Resolve-Path (Join-Path $PSScriptRoot "..\..")
$composeFile = Join-Path $repoRoot "deploy\docker\docker-compose.dev.yml"

Push-Location $repoRoot
try {
    docker compose -f $composeFile up --build @args
}
finally {
    Pop-Location
}
