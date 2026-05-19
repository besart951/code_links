#!/usr/bin/env sh
set -eu

if [ "${CONFIRM_RESET:-}" != "codelinks" ]; then
  echo "Refusing to reset database volumes. Re-run with CONFIRM_RESET=codelinks."
  exit 1
fi

ROOT_DIR="$(CDPATH= cd -- "$(dirname -- "$0")/../.." && pwd)"
cd "$ROOT_DIR"

docker compose -f deploy/docker/docker-compose.dev.yml down -v
docker compose -f deploy/docker/docker-compose.dev.yml up --build postgres platform-migrate infra-migrate
