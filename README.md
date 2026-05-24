# CodeLinks Platform

pnpm monorepo scaffold for the CodeLinks SaaS platform.

## Stack

- SvelteKit / Svelte 5 frontends with runes
- Go product and auth backends
- PostgreSQL for the central auth service
- RS256 access JWTs with JWKS-based local validation in product backends
- Tailwind CSS, shadcn-svelte, Paraglide i18n

## Local Commands

```sh
pnpm install
pnpm generate:contracts
pnpm check
pnpm test
pnpm build
```

`pnpm test` starts a throwaway Postgres container for the Auth Service store contract. Use `pnpm go:test` for the non-Docker Go path only.

Run the landing page locally:

```sh
pnpm dev
```

Run the Admin Console locally with mock auth/data:

```sh
pnpm dev:admin
```

Open `http://localhost:5178/admin`.

Run the Docker/Caddy dev stack:

```sh
pnpm compose:dev
```

The Docker/Caddy stack needs port 80. If a `*.codelinks.localhost` URL shows an Apache default page, stop the local Apache service that is already bound to port 80, then restart the stack.

Local Caddy hosts:

- `http://code-links.codelinks.localhost`
- `http://auth.codelinks.localhost`
- `http://admin-link.codelinks.localhost`
- `http://infra-link.codelinks.localhost`
- `http://planer-link.codelinks.localhost`
- `http://loka-link.codelinks.localhost`

Demo credentials:

- Email: `demo@codelinks.dev`
- Password: `password`
