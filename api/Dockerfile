FROM golang:1.25.12-bookworm AS build

WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download

COPY api ./api
COPY cmd ./cmd
COPY internal ./internal
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /out/nav-api ./cmd/nav
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /out/nav-healthcheck ./cmd/healthcheck
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /out/nav-bindingcheck ./cmd/bindingcheck
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /out/nav-bootstrap ./cmd/bootstrap

FROM migrate/migrate:v4.19.1 AS migration
COPY manifest/sql/migrations /migrations

FROM gcr.io/distroless/static-debian12:nonroot AS bindingcheck
COPY --from=build /out/nav-bindingcheck /usr/local/bin/nav-bindingcheck
USER nonroot:nonroot
ENTRYPOINT ["/usr/local/bin/nav-bindingcheck"]

FROM gcr.io/distroless/static-debian12:nonroot AS bootstrap
COPY --from=build /out/nav-bootstrap /usr/local/bin/nav-bootstrap
USER nonroot:nonroot
ENTRYPOINT ["/usr/local/bin/nav-bootstrap"]

FROM gcr.io/distroless/static-debian12:nonroot AS runtime
WORKDIR /app
COPY --from=build /out/nav-api /usr/local/bin/nav-api
COPY --from=build /out/nav-healthcheck /usr/local/bin/nav-healthcheck
COPY --chown=nonroot:nonroot manifest ./manifest

ENV GF_SERVER_ADDRESS=0.0.0.0:8090
ENV OTEL_TRACES_EXPORTER=none

EXPOSE 8090
USER nonroot:nonroot
ENTRYPOINT ["/usr/local/bin/nav-api"]
