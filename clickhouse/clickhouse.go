package clickhouse

import (
	"HezzelTask/config"
	"database/sql"
	"fmt"
	"github.com/ClickHouse/clickhouse-go"
 )

type ClickHouse struct {
	Ch *sql.DB
}

func Connect(cfg *config.Config) (*sql.DB, error) {
	err := createDefaultDatabase(cfg)
	if err != nil {
		return nil, err
	}
	connect, err := sql.Open("clickhouse", fmt.Sprintf("tcp://%s:%s?database=%s&username=default", cfg.Clickhouse.Host, cfg.Clickhouse.Port, cfg.Clickhouse.Database))
	if err != nil {
		return nil, err
	}

	if err := connect.Ping(); err != nil {
		return nil, err
	}
	//_, err = connect.Exec(`
	//	DROP TABLE IF EXISTS my_logs
	//	`)
	_, err = connect.Exec(`
		CREATE TABLE IF NOT EXISTS my_logs (
		    log_text    String,
			event_name  String,
			eventDate   DateTime
		) engine=Kafka('kafka:9092', 'addUserLogs', 'now', 'JSONEachRow')
		`)
	if err != nil {
		return nil, err
	}



	if err != nil {
		return nil, err
	}

	return connect, nil
}

func createDefaultDatabase(cfg *config.Config) error {
	connect, err := sql.Open("clickhouse", fmt.Sprintf("tcp://%s:%s?username=default", cfg.Clickhouse.Host, cfg.Clickhouse.Port))
	if err != nil {
		return err
	}

	if err = connect.Ping(); err != nil {
		if exception, ok := err.(*clickhouse.Exception); ok {
			fmt.Printf("[%d] %s \n%s\n", exception.Code, exception.Message, exception.StackTrace)
		}
		return err
	}
	_, err = connect.Exec("CREATE DATABASE IF NOT EXISTS " + "logs")
	tx, err := connect.Begin()
	if err != nil {
		return err
	}
	stmt, err := tx.Prepare("ISELECT 1")
	if err != nil {
		return err
	}
	defer stmt.Close()
	defer connect.Close()

	return nil
}
