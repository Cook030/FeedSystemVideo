# syntax=docker/dockerfile:1.7
#
# Multi-stage build for backend API and worker.
# Uses vendor directory for offline builds (no network needed during build).
#
# Build:
#   docker build -f backend/Dockerfile -t feedsystem-backend:api --target api .
#   docker build -f backend/Dockerfile -t feedsystem-backend:worker --target worker .
#
# Prerequisite: run `go mod vendor` in backend/ before building.

ARG GO_VERSION=1.24.5

FROM golang:${GO_VERSION} AS build

WORKDIR /src/backend

COPY backend/go.mod backend/go.sum ./
COPY backend/vendor/ ./vendor/

ENV CGO_ENABLED=0
COPY backend/ ./
RUN go build -mod=vendor -trimpath -ldflags="-s -w" -o /out/api ./cmd
RUN go build -mod=vendor -trimpath -ldflags="-s -w" -o /out/worker ./cmd/worker

FROM alpine:3.21 AS base
RUN apk add --no-cache ca-certificates tzdata && adduser -D -H -s /sbin/nologin app
WORKDIR /app
COPY --from=build /src/backend/configs ./configs
RUN mkdir -p ./.run/uploads && chown -R app:app /app
USER app

FROM base AS api
COPY --from=build /out/api /app/api
EXPOSE 8080
ENTRYPOINT ["/app/api"]

FROM base AS worker
COPY --from=build /out/worker /app/worker
ENTRYPOINT ["/app/worker"]
