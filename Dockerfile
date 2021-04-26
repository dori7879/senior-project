# Build the Go app
FROM golang:1.16-buster AS builder
ADD . /app
WORKDIR /app/api
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-w" -a -o ./bin/app ./cmd/apid

# Build the React app
FROM node:12-alpine AS node_builder
COPY --from=builder /app/webapp ./
RUN npm install
RUN npm i react-router-dom
RUN npm run build

# Final stage build, this will be the container
# that we will deploy to production
FROM alpine:latest
RUN apk update && apk --no-cache add ca-certificates bash
COPY --from=builder /app/api/bin ./
COPY --from=node_builder /build ./web

EXPOSE 8080

RUN chmod +x ./app
ENTRYPOINT ["./app"]