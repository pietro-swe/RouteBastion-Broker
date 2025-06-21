FROM golang:1.24.4-alpine

WORKDIR /app

RUN go install github.com/air-verse/air@latest

COPY go.mod go.sum ./
RUN go mod download

ADD cmd /cmd
ADD scripts /scripts
ADD internal /internal
ADD sql /sql

COPY .air.toml .
COPY sqlc.yml .
COPY app.env .

EXPOSE 8080

CMD ["air"]
