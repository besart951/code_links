FROM node:24-alpine AS builder
RUN corepack enable && corepack prepare pnpm@10.29.1 --activate
WORKDIR /repo
COPY package.json pnpm-workspace.yaml turbo.json tsconfig.base.json pnpm-lock.yaml* ./
COPY packages ./packages
COPY apps/startseite ./apps/startseite
RUN pnpm install --filter @codelinks/startseite... --frozen-lockfile=false
RUN pnpm --filter @codelinks/startseite build

FROM node:24-alpine AS runtime
WORKDIR /app
ENV NODE_ENV=production
ENV PORT=3000
COPY --from=builder /repo/apps/startseite/build ./build
COPY --from=builder /repo/apps/startseite/package.json ./package.json
COPY --from=builder /repo/node_modules ./node_modules
COPY --from=builder /repo/apps/startseite/node_modules ./node_modules
EXPOSE 3000
CMD ["node", "build"]
