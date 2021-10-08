package handler

import (
	"HezzelTask/clickhouse"
	"HezzelTask/models"
	"context"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"log"
	"time"
)

type Handler struct {
	Ch    *clickhouse.ClickHouse
	Kafka Kafka
	Db    Db
}

type Db interface {
	AddUser(ctx context.Context, user *models.User) (*models.User, error)
	DeleteUser(ctx context.Context, email string) error
	UserList(ctx context.Context) (*[]models.User, error)
}

type Kafka interface {
	WriteMessage(key []byte, value []byte) error
}

func (h *Handler) AddUser(ctx context.Context, name, email, phone string) error {
	_, err := h.Db.AddUser(ctx, &models.User{
		Name:  name,
		Email: email,
		Phone: phone,
	})
	if err != nil {
		return errors.Wrap(err, "Failed creating a new user: ")
	}
	log.Println("User add in DB")
	myLog, err := json.Marshal(models.Log{
		EventName: "User create",
		LogText:   fmt.Sprintf("Add user with parametrs Phone:%s Name:%s Email:%s", phone, name, email),
		EventDate: time.Now(),
	})
	if err != nil {
		return err
	}
	err = h.Kafka.WriteMessage([]byte(email), myLog)
	if err != nil {
		return err
	}
	log.Println("User add log in DB")
	return nil
}

func (h *Handler) DeleteUser(ctx context.Context, email string) error {
	err := h.Db.DeleteUser(ctx, email)
	if err != nil {
		return errors.Wrap(err, "Failed deleting user: ")
	}

	return nil
}

func (h *Handler) UserList(ctx context.Context) (*[]models.User, error) {
	users, err := h.Db.UserList(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to take user list: ")
	}
	//Добавить редиску
	return users, nil
}
