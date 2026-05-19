FROM node:24-alpine AS builder
RUN corepack enable && corepack prepare pnpm@10.29.1 --activate
WORKDIR /repo
COPY package.json pnpm-workspace.yaml turbo.json tsconfig.base.json pnpm-lock.yaml ./
COPY packages ./packages
COPY apps/infra_link/frontend ./apps/infra_link/frontend
RUN pnpm install --filter @codelinks/infra-link-frontend... --frozen-lockfile
RUN pnpm --filter @codelinks/infra-link-frontend build

FROM caddy:2.11.2-alpine AS runtime
COPY apps/infra_link/frontend/Caddyfile /etc/caddy/Caddyfile
COPY --from=builder /repo/apps/infra_link/frontend/build /usr/share/caddy
EXPOSE 80
