# =============================================================================
# Builder stage
# =============================================================================
# Uses the official Go image (Debian-based), which provides bash (needed by
# the vugugen script) and a full Go toolchain.
#
# Alternative: golang:1.26-alpine would yield a smaller builder image, but
# the alpine shell lacks bash and `shopt` which frontend/generate.bash uses.
# You'd need to `apk add bash` or rewrite the script for POSIX sh.
FROM golang:1.26 AS builder

WORKDIR /src

# -- Dependency caching -------------------------------------------------------
# Copy only the module files first so that `go mod download` is cached as long
# as dependencies don't change, even when source code does.
COPY go.mod go.sum ./
RUN go mod download

# -- Code generation tools ----------------------------------------------------
# vgrun bootstraps the vugu toolchain (installs vugugen); enumer generates
# enum helper methods used by domain types.
RUN go install github.com/vugu/vgrun@latest \
    && vgrun -install-tools \
    && go install github.com/dmarkham/enumer@latest

# -- Source code --------------------------------------------------------------
COPY . .

# -- Code generation ----------------------------------------------------------
# 1. enumer: generates *_enumer.go files for Club, Status, Role types
# 2. vugugen: generates *_vgen.go component code from .vugu templates
RUN go generate ./...
RUN bash ./frontend/generate.bash

# -- WASM frontend ------------------------------------------------------------
# Compiles the Vugu frontend to WebAssembly. The resulting app.wasm is placed
# in static/files/ where it gets embedded into the backend binary via //go:embed.
RUN GOARCH=wasm GOOS=js go build -o ./static/files/app.wasm ./frontend/

# -- Build arguments for version info -----------------------------------------
# These are optional. Pass them at build time to embed version metadata:
#   docker build \
#     --build-arg GIT_COMMIT=$(git rev-parse --short HEAD) \
#     --build-arg GIT_TIME="$(git show -s --format=%ci HEAD)" \
#     --build-arg GIT_TAG=$(git describe --exact-match --tags HEAD 2>/dev/null || echo "") \
#     .
ARG GIT_COMMIT=""
ARG GIT_TIME=""
ARG GIT_TAG=""

# -- Backend binary -----------------------------------------------------------
# CGO_ENABLED=0: the production binary only uses PostgreSQL (lib/pq, pure Go).
#   SQLite (which requires CGO) is only used in dev/test and is not imported by
#   the production main.go. This lets us produce a fully static binary.
#
# NOTE: the Makefile's PACKAGE path (src/infrastructure/router) appears stale;
#   the actual package is backend/infrastructure/router. We use the correct
#   path here so that -ldflags -X actually injects the version variables.
RUN CGO_ENABLED=0 go build -o /streepjes \
    -ldflags " \
      -X 'github.com/FallenTaters/streepjes/backend/infrastructure/router.buildCommit=${GIT_COMMIT}' \
      -X 'github.com/FallenTaters/streepjes/backend/infrastructure/router.buildTime=${GIT_TIME}' \
      -X 'github.com/FallenTaters/streepjes/backend/infrastructure/router.buildVersion=${GIT_TAG}' \
    " \
    .

# =============================================================================
# Runtime stage
# =============================================================================
# gcr.io/distroless/static contains only CA certificates and tzdata -- no
# shell, no package manager. This minimises the attack surface and image size.
#
# Alternatives:
#   - `scratch`: even smaller, but has no CA certs or tzdata. You'd need to
#     COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
#     and import the `time/tzdata` package in Go if timezone support is needed.
#   - `alpine:3`: ~7 MB, includes a shell for debugging. Useful during early
#     development when you want to `kubectl exec` into the pod.
#   - `gcr.io/distroless/static-debian12:debug`: distroless with a busybox
#     shell, handy for one-off debugging without switching to alpine.
FROM gcr.io/distroless/static-debian12

COPY --from=builder /streepjes /streepjes

# -- Default configuration ----------------------------------------------------
# STREEPJES_DISABLE_SECURE=true: assumes TLS termination at the Kubernetes
#   ingress controller (e.g. nginx-ingress, Traefik, or a cloud LB). The
#   container serves plain HTTP only.
#
# STREEPJES_PORT=80: the default HTTP port. Override via k8s env or configmap.
#
# STREEPJES_DB_CONNECTION_STRING: must be provided at runtime, e.g. via a k8s
#   Secret. There is no sensible default for production.
#
# To re-enable in-container TLS (not typical for k8s), set:
#   STREEPJES_DISABLE_SECURE=false
#   STREEPJES_TLS_CERT_PATH=/path/to/cert.pem
#   STREEPJES_TLS_KEY_PATH=/path/to/key.pem
#   and mount the certificate files into the container.
ENV STREEPJES_DISABLE_SECURE=true
ENV STREEPJES_PORT=80

EXPOSE 80
# Optional: if you re-enable TLS, also expose 443:
# EXPOSE 443

# Run as non-root. 65534 is the "nobody" user in distroless.
USER 65534

ENTRYPOINT ["/streepjes"]
