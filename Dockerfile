FROM golang:1.26.1-alpine

RUN apk add --no-cache gcc musl-dev

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY  . .

RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -o /todo-app_v1 main.go

ENV TODO_PORT=7540 \
    TODO_DBFILE=/app/db/scheduler.db \
    TODO_PASSWORD=""

EXPOSE 7540

CMD ["/todo-app_v1"]