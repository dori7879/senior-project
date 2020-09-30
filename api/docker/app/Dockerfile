# # Multi-Stage Build
#
# Build environment
# -----------------
FROM golang:1.15.2-alpine as build-env
WORKDIR /api

RUN apk update && apk add --no-cache gcc musl-dev git

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -ldflags '-w -s' -a -o ./bin/app ./cmd/app \
    && go build -ldflags '-w -s' -a -o ./bin/migrate ./cmd/migrate


# Deployment environment
# ----------------------
FROM alpine
RUN apk update && apk add --no-cache bash mysql-client

COPY --from=build-env /api/bin/app /api/
COPY --from=build-env /api/bin/migrate /api/
COPY --from=build-env /api/migrations /api/migrations

COPY --from=build-env /api/docker/app/bin /usr/local/bin/myapi/
RUN chmod +x /usr/local/bin/myapi/*

EXPOSE 8080
CMD ["/api/app"]