FROM node:22-alpine AS admin-builder

WORKDIR /app

RUN corepack enable

COPY admin/package.json admin/pnpm-lock.yaml admin/pnpm-workspace.yaml ./
RUN pnpm install --frozen-lockfile

COPY admin/. ./
COPY shared /shared

ARG VITE_APP_BASE=/admin/
ARG VITE_APP_NAME=Grtblog Admin
ARG VITE_APP_TITLE=管理后台
ARG VITE_WATERMARK_CONTENT=
ARG VITE_API_BASE_URL=/api/v2

ENV VITE_APP_BASE=${VITE_APP_BASE} \
    VITE_APP_NAME=${VITE_APP_NAME} \
    VITE_APP_TITLE=${VITE_APP_TITLE} \
    VITE_WATERMARK_CONTENT=${VITE_WATERMARK_CONTENT} \
    VITE_API_BASE_URL=${VITE_API_BASE_URL}

RUN pnpm build

FROM golang:1.25-alpine AS builder

WORKDIR /src/server

RUN apk add --no-cache ca-certificates git

ARG GOOSE_VERSION=v3.26.0
RUN GOBIN=/out go install github.com/pressly/goose/v3/cmd/goose@${GOOSE_VERSION}

COPY server/go.mod server/go.sum ./
RUN go mod download

COPY server/. .

ARG APP_VERSION=dev
ARG BUILD_COMMIT=unknown

RUN CGO_ENABLED=0 GOOS=linux \
  go build -trimpath -ldflags="-s -w \
  -X github.com/grtsinry43/grtblog-v2/server/internal/buildinfo.BuildVersion=${APP_VERSION} \
  -X github.com/grtsinry43/grtblog-v2/server/internal/buildinfo.BuildCommit=${BUILD_COMMIT}" \
  -o /out/grtblog-server ./cmd/api

RUN CGO_ENABLED=0 GOOS=linux \
  go build -trimpath -ldflags="-s -w" \
  -o /out/grtblog-restore ./cmd/restore

FROM alpine:3.21 AS runtime

RUN apk add --no-cache ca-certificates tzdata su-exec postgresql17-client \
  && addgroup -g 10001 -S app \
  && adduser -u 10001 -S app -G app

WORKDIR /app

COPY --from=builder /out/grtblog-server /app/grtblog-server
COPY --from=builder /out/grtblog-restore /app/grtblog-restore
COPY --from=builder /out/goose /usr/local/bin/goose
COPY --from=builder /src/server/docs /app/docs
COPY --from=builder /src/server/migrations /app/migrations
COPY --from=admin-builder /app/dist /app/admin
COPY deploy/docker/server-entrypoint.sh /usr/local/bin/server-entrypoint.sh

RUN mkdir -p /app/storage/html /app/storage/uploads /app/storage/backups /app/storage/geoip \
  && chown -R app:app /app \
  && chmod +x /usr/local/bin/server-entrypoint.sh

EXPOSE 8080

ENTRYPOINT ["/usr/local/bin/server-entrypoint.sh"]
