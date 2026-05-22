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
pnpm check
pnpm test
pnpm build
```

Run the landing page locally:

```sh
pnpm dev
```

Run the Docker/Caddy dev stack:

```sh
docker compose up --build
```

Local Caddy hosts:

- `http://code-links.codelinks.localhost`
- `http://auth.codelinks.localhost`
- `http://infra-link.codelinks.localhost`
- `http://planer-link.codelinks.localhost`
- `http://loka-link.codelinks.localhost`

Demo credentials:

- Email: `demo@codelinks.dev`
- Password: `password`
