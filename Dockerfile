# Stage 1

# build stage: compile the Go binary
FROM golang:1.25-alpine AS builder

# work inside /app
WORKDIR /app

# tools needed during build (e.g. for go modules)
RUN apk add --no-cache git

# copy module files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# copy the rest of the code and build the binary from cmd/
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /bin/pismo-account ./cmd

# Stage 2

# runtime stage: minimal image that just runs the binary
FROM alpine:3.19

# workdir for the running app
WORKDIR /app

# basic runtime dependencies
RUN apk add --no-cache ca-certificates tzdata

# copy the built binary from the build stage
COPY --from=builder /bin/pismo-account /app/pismo-account
COPY --from=builder /app/config/config.json /config/config.json

# default location of the config file (overridable)
ENV CONFIG_PATH=/config/config.json

# app listens on 8080
EXPOSE 8080

# start the service
CMD ["/app/pismo-account"]
