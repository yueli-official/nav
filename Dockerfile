# syntax=docker/dockerfile:1.4

# ----------------------
# 前端构建阶段 (Vue)
# ----------------------
FROM node:20-alpine AS frontend-builder
WORKDIR /frontend

# 先复制依赖清单文件
COPY frontend/package.json frontend/pnpm-lock.yaml ./

# 使用 BuildKit Secret 传递 GH_TOKEN（不会进入镜像层）
RUN --mount=type=secret,id=GH_TOKEN \
  GH_TOKEN=$(cat /run/secrets/GH_TOKEN) && \
  npm install -g pnpm && \
  echo "@yuelioi:registry=https://npm.pkg.github.com" > .npmrc && \
  echo "//npm.pkg.github.com/:_authToken=${GH_TOKEN}" >> .npmrc && \
  corepack enable && corepack prepare pnpm@latest --activate && \
  pnpm install --frozen-lockfile && \
  pnpm store prune && \
  rm -f .npmrc

# 然后再复制完整源码
COPY frontend/ ./

# 构建
RUN pnpm build

# 清理开发文件，只保留 dist
RUN rm -rf src node_modules public package.json pnpm-lock.yaml vite.config.ts tsconfig*.json eslint.config.ts env.d.ts

# ----------------------
# 后端构建阶段 (Go)
# ----------------------
FROM golang:1.25.1-alpine AS backend-builder
WORKDIR /app

ENV DATABASE_URL=/app/index.db
ENV APP_MODE=release

RUN apk add --no-cache git build-base

COPY go.mod go.sum ./
RUN go mod download

# 复制后端源码
COPY . .

# 复制前端构建产物
COPY --from=frontend-builder /frontend/dist ./frontend/dist

RUN go build -o server .


# ----------------------
# 运行阶段
# ----------------------
FROM alpine:latest
WORKDIR /app

COPY --from=backend-builder /app/server ./server
COPY --from=backend-builder /app/frontend/dist ./frontend/dist

ENV DATABASE_URL=/app/index.db

EXPOSE 9000
CMD ["./server"]
