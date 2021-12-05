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

func TestIpToScore(t *testing.T) {
	ip := "192.168.255.10"
	score := IpToScore(ip)
	fmt.Println(score)
}

func TestImportIpsToRedis(t *testing.T) {
	conn, err := initClient()
	if err != nil {
		fmt.Println("init redis client err", err)
		return
	}

	ImportIpsToRedis(conn, "ip-city.csv")
}

func TestImportCitiesToRedis(t *testing.T) {
	conn, err := initClient()
	if err != nil {
		fmt.Println("init redis client err", err)
		return
	}

	err = ImportCitiesToRedis(conn, "citydetail.csv")
	if err != nil {
		fmt.Println(err)
	}
}

func TestFindCityByIp(t *testing.T) {
	conn, err := initClient()
	if err != nil {
		fmt.Println("FindCityByIp error", err)
		return
	}
	ip := "192.168.4.101"
	err = FindCityByIp(conn, ip)
}
