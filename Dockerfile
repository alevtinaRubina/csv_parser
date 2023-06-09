FROM golang:1.17

WORKDIR /app

COPY . .

RUN go mod init localhost/promotions
RUN go get github.com/gorilla/mux
RUN go build -o main .

COPY promotions.csv .

EXPOSE 1321

CMD ["./main"]
