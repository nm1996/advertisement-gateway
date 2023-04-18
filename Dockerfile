FROM golang:1.19

WORKDIR /app

COPY . .

RUN go mod tidy && go build -o app .

CMD ["./app"]