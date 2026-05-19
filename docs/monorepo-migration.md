# Monorepo Migration Notes

Initial conversion branch: `monorepo-migration`.

Imported submodule SHAs:

- `apps/go_infra_link` -> `apps/infra_link`: `8643ca53d94d198d30d80b734a5d23c91e19899f`
- `apps/planer_link`: `11fd3f3c0848da66fe92abaa68bcce2aaedd0915`
- `apps/loka_link`: `9dd668c8f25c89d2067327c11ae0b2eecf959bf8`

The working tree now tracks these products as normal monorepo directories rather
than Git submodules. Product histories can still be recovered from their source
repositories if a later `git subtree` import is needed for long-range blame.
