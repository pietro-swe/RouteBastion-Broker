FROM golang:1.25.5-alpine

WORKDIR /app

RUN mkdir -p /home/user/.cache/go-build /go/pkg/mod && \
    chmod -R 1777 /home/user/.cache && \
    chmod -R 1777 /go/pkg

RUN go install github.com/air-verse/air@latest

EXPOSE 8080

CMD ["air"]
