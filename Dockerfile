FROM golang

WORKDIR /app

RUN go mod init bookstore-api

COPY app/ ./

RUN go get github.com/labstack/echo/v4 github.com/labstack/echo/v4/middleware github.com/jackc/pgx/v5 golang.org/x/crypto/bcrypt

RUN go build -o api ./server.go

CMD ["./api"]
