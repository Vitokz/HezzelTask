package main

import (
	"HezzelTask/clickhouse"
	"HezzelTask/config"
	"HezzelTask/handler"
	"HezzelTask/kafka"
	"HezzelTask/proto"
	"HezzelTask/repository"
	"google.golang.org/grpc"
	"log"
	"net"
	"time"
)

const (
	port = "8080"
)

func main() {
	time.Sleep(10*time.Second)
	cfg := config.Get()

	//Создаю соединение с PostgreSQL
	db, err := repository.Connect(cfg.Pg.PgUrl)
	defer db.Close()
    err = repository.RunPgMigrations(cfg)
    if err != nil {
    	panic(err)
	}

	//Создаю соединение с брокером Kafka
	kf, wr, err := kafka.Connect(cfg)
	if err != nil {
		panic(err)
	}
	defer kf.Close()
	defer wr.Close()
	time.Sleep(5*time.Second)
	//Создаю соединение с Clickhouse( только для того чтобы настроить базу )
	ch,err := clickhouse.Connect(cfg)
	if err != nil {
		panic(err)
	}
	defer ch.Close()

	hdlr := &handler.Handler{
		Db:    &repository.Pg{Db: db},
		Ch:    &clickhouse.ClickHouse{Ch: ch},
		Kafka: &kafka.Kafka{Kf: kf, Writer: wr},
	}

	s := grpc.NewServer()
	srv := proto.GRPCServer{Handler: hdlr}
	proto.RegisterHezzelUsersServer(s, &srv)

	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatal(err)
	}
	if err = s.Serve(lis); err != err {
		log.Fatal(err)
	}
}
