FROM golang:1.26 AS builder

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

# vgrun bootstraps the vugu toolchain (installs vugugen). enumer is managed
# as a tool dependency in go.mod and resolved automatically by go generate.
RUN go install github.com/vugu/vgrun@latest \
    && vgrun -install-tools

COPY . .

RUN go generate ./...
RUN bash ./frontend/generate.bash

RUN GOARCH=wasm GOOS=js go build -o ./static/files/app.wasm ./frontend/

# Optional: pass at build time to embed version metadata.
ARG GIT_COMMIT=""
ARG GIT_TIME=""
ARG GIT_TAG=""

# CGO_ENABLED=0: production only uses PostgreSQL (lib/pq, pure Go).
# SQLite (which requires CGO) is only used in dev/test.
RUN CGO_ENABLED=0 go build -o /streepjes \
    -ldflags " \
      -X 'github.com/FallenTaters/streepjes/backend/infrastructure/router.buildCommit=${GIT_COMMIT}' \
      -X 'github.com/FallenTaters/streepjes/backend/infrastructure/router.buildTime=${GIT_TIME}' \
      -X 'github.com/FallenTaters/streepjes/backend/infrastructure/router.buildVersion=${GIT_TAG}' \
    " \
    .

FROM gcr.io/distroless/static-debian12

COPY --from=builder /streepjes /streepjes

# STREEPJES_DB_CONNECTION_STRING must be provided at runtime.
#
# To re-enable in-container TLS, set:
#   STREEPJES_DISABLE_SECURE=false
#   STREEPJES_TLS_CERT_PATH=/path/to/cert.pem
#   STREEPJES_TLS_KEY_PATH=/path/to/key.pem
ENV STREEPJES_DISABLE_SECURE=true
ENV STREEPJES_PORT=80

EXPOSE 80

USER 65534

ENTRYPOINT ["/streepjes"]
