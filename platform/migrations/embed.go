package migrations

import "embed"

// FS contains forward-only platform database migrations.
//
//go:embed *.sql
var FS embed.FS
