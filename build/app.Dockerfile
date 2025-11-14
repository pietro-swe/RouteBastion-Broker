FROM golang:1.25.1-alpine AS build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o broker ./cmd

FROM alpine:3.20

WORKDIR /app
COPY --from=build /app/broker .
COPY --from=build /app/app.env ./app.env

EXPOSE 8080
CMD ["./broker"]
