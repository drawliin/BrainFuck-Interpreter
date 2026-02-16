# syntax=docker/dockerfile:1

FROM golang:1.25-alpine AS build
WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /out/bf .

FROM alpine:3.21
WORKDIR /app

RUN addgroup -S app && adduser -S -G app app

COPY --from=build /out/bf /app/bf
COPY static /app/static

EXPOSE 8080
USER app
ENTRYPOINT ["/app/bf"]
