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

RUN go build -ldflags '-w -s' -a -o ./bin/app ./cmd/app


# Deployment environment
# ----------------------
FROM alpine
RUN apk update && apk add --no-cache bash

COPY --from=build-env /api/bin/app /api/

EXPOSE 8080
CMD ["/api/app"]