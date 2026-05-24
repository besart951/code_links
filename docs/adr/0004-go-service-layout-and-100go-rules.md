# Go Service Layout And 100go Rules

Accepted.

CodeLinks Go service backends follow the 100 Go Mistakes guidance for project organization, interfaces, HTTP servers, SQL resource handling, and tests: keep command entrypoints thin, put service implementation under `internal/`, avoid producer-side interface pollution, avoid utility dumping grounds, configure HTTP server timeouts, check SQL row errors, and run race-enabled tests when touching concurrent code. We keep Auth Service domain language in `CONTEXT.md`; implementation package seams live in code and ADRs.

Sources: https://100go.co/, especially #5 interface pollution, #6 producer-side interfaces, #12 project/package organization, #13 utility packages, #78 SQL mistakes, #81 default HTTP client/server, #83 race flag.
