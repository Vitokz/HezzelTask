FROM golang:1.16

COPY . /go/src/app

WORKDIR /go/src/app

#RUN #go mod download

RUN go build -o rpcserv main.go

EXPOSE 8080

CMD ["./rpcserv"]