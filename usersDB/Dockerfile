# Arguments
ARG GOLANG_VERSION
ARG ALPINE_VERSION

# Build part
FROM golang:${GOLANG_VERSION}-alpine${ALPINE_VERSION} AS builder

RUN adduser -D -h /tmp/build build
RUN apk --no-cache add \
        make \
        curl \
        tar \
        git

USER build
WORKDIR /tmp/build

COPY --chown=build Makefile Makefile
COPY --chown=build go.mod go.mod
COPY --chown=build go.sum go.sum

RUN make download

ARG VERSION
ARG NAME
ARG BUILD_TIME
ARG LAST_COMMIT_USER
ARG LAST_COMMIT_HASH
ARG LAST_COMMIT_TIME

COPY --chown=build rest.go rest.go
COPY --chown=build http.go http.go
COPY --chown=build config.go config.go
COPY --chown=build main.go main.go

RUN make build

# Exec part
FROM gcr.io/distroless/base:nonroot

ARG VERSION
ARG NAME

# Appication Configuration
ENV DEBUG="false"
ENV HTTP_PORT="8080"
ENV DB_HOST="mysql.default.svc.cluster.local:3306"
ENV DB_USER="root"
ENV DB_PASS="password"
ENV DB_NAME="users"
ENV DB_MIGRATIONS_PATH="/opt/app/migrations/"

# Copy from builder
COPY --chown=nonroot:nonroot migrations ${DB_MIGRATIONS_PATH}
COPY --from=builder /tmp/build/${NAME}-${VERSION} /usr/local/bin/app

# Exec
CMD ["app"]
