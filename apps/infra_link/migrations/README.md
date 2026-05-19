# infra_link migrations

`infra_link` owns the `codelinks_infra_link` database.

The current Go backend still keeps its forward-only migration implementation in
`backend/internal/db`. New product migrations should stay product-local and must
not read or write the platform database.
