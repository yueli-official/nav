FROM node:24.22.1-bookworm-slim AS build

RUN corepack enable && corepack prepare pnpm@10.28.2 --activate
WORKDIR /app/web
COPY web/package.json web/pnpm-lock.yaml ./
COPY web/vendor ./vendor
RUN pnpm install --frozen-lockfile --ignore-scripts

COPY web/ ./
COPY api/contracts /app/api/contracts
ENV NODE_OPTIONS=--max-old-space-size=4096
RUN pnpm build

FROM node:24.22.1-bookworm-slim

WORKDIR /app
COPY --from=build /app/web/.output ./.output

ENV HOST=0.0.0.0
ENV PORT=3006

EXPOSE 3006
CMD ["node", ".output/server/index.mjs"]
