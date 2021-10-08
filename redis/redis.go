package redis

import (
	"HezzelTask/config"
	"HezzelTask/models"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis"
	"time"
)

type Redis struct {
	Redis *redis.Client
}

func Connect(cfg *config.Config) (*redis.Client,error){
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s",cfg.Redis.Host,cfg.Redis.Port),
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	err := rdb.Set("key", "value", 1*time.Minute).Err()
	if err != nil {
		return nil,err
	}
	return rdb, nil
}

func (r *Redis)AddUsers(users *[]models.User) error {
	for _,v := range *users {
		user, err := json.Marshal(v)
		if err != nil {
			return err
		}
		err = r.Redis.Set(v.Email,user,time.Minute*1).Err()
		if err != nil {
			return err
		}
	}
	return nil
}