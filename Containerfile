FROM golang:1.26 AS builder

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go generate ./...

ARG GIT_COMMIT=""
ARG GIT_TIME=""
ARG GIT_TAG=""

RUN CGO_ENABLED=0 go build -o /streepjes \
    -ldflags " \
      -X 'github.com/FallenTaters/streepjes/backend/infrastructure/router.buildCommit=${GIT_COMMIT}' \
      -X 'github.com/FallenTaters/streepjes/backend/infrastructure/router.buildTime=${GIT_TIME}' \
      -X 'github.com/FallenTaters/streepjes/backend/infrastructure/router.buildVersion=${GIT_TAG}' \
    " \
    .

FROM gcr.io/distroless/static-debian12

COPY --from=builder /streepjes /streepjes

ENV STREEPJES_DISABLE_SECURE=true
ENV STREEPJES_PORT=80

EXPOSE 80

USER 65534

ENTRYPOINT ["/streepjes"]
