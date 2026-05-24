#!/bin/sh
set -eu

corepack enable

if [ ! -d node_modules/.pnpm ]; then
	while ! mkdir node_modules/.install-lock 2>/dev/null; do
		sleep 1
	done
	trap 'rmdir node_modules/.install-lock 2>/dev/null || true' EXIT
	if [ ! -d node_modules/.pnpm ]; then
		corepack pnpm install --frozen-lockfile --config.confirmModulesPurge=false
	fi
	rmdir node_modules/.install-lock 2>/dev/null || true
fi

exec corepack pnpm "$@"
