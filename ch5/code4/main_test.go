package main

import (
	"context"
	"fmt"
	"testing"

	"github.com/go-redis/redis/v8"
)

func initClient() (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // 默认使用0号数据库
	})

	_, err := rdb.Ping(context.Background()).Result()
	if err != nil {
		return nil, err
	}
	return rdb, nil
}

func TestIsUnderMaintenace(t *testing.T) {
	conn, err := initClient()
	if err != nil {
		fmt.Println("FindCityByIp error", err)
		return
	}
	IsUnderMaintenace(conn)
}
