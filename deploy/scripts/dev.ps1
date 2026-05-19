Set-StrictMode -Version Latest
$ErrorActionPreference = "Stop"

$composeFile = Join-Path $PSScriptRoot "..\docker\docker-compose.yml"
docker compose -f $composeFile up --build
