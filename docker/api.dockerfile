# syntax=docker/dockerfile:1
FROM golang:1.26-alpine AS build-stage

WORKDIR /app

COPY ./go.* /app/
RUN go mod download

COPY ./ /app

ENV GOCACHE=/root/.cache/go-build
RUN --mount=type=cache,target="/root/.cache/go-build" GOOS=linux go build -o api cmd/main.go

FROM scratch
WORKDIR /app
COPY --from=build-stage /app/api /app/api
COPY --from=build-stage /app/config.yaml /app/config.yaml

# Run
CMD ["/app/api"]
