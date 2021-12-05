package main

import (
	"bufio"
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/go-redis/redis/v8"
)

type CityDetailed struct {
	CityID   string
	Country  string
	Region   string
	CityName string
}

func main() {

}

//
func FindCityByIp(conn *redis.Client, ipAddress string) error {
	score := strconv.Itoa(IpToScore(ipAddress))
	res, err := conn.ZRevRangeByScore(context.Background(), "ip2cityid:", &redis.ZRangeBy{
		Max: score,
	}).Result()

	if err != nil {
		fmt.Println("ZRevRangeByScore error", err)
		return err
	}
	if len(res) < 1 {
		fmt.Println("not find")
		return nil
	}

	data, err := conn.HGet(context.Background(), "cityId2city", res[0]).Result()
	if err != nil {
		fmt.Println("hget err", err)
		return err
	}
	cityDetail := CityDetailed{}
	err = json.Unmarshal([]byte(data), &cityDetail)
	if err != nil {
		fmt.Println("Unmarshal error", err)
	}
	fmt.Println("cityDetail is ", cityDetail)
	return nil
}

// ImportIpsToRedis 建立cityid-ip的映射
func ImportIpsToRedis(conn *redis.Client, fileName string) {
	csvFile, err := os.Open(fileName)
	if err != nil {
		fmt.Println("open file err", err)
		return
	}
	count := 0
	reader := csv.NewReader(bufio.NewReader(csvFile))
	for {
		line, err := reader.Read()
		if err == io.EOF {
			break
		}
		count++
		if line[0] == "startIpNum" {
			continue
		}
		score := IpToScore(line[1])

		cityID := line[2]

		conn.ZAdd(context.Background(), "ip2cityid:", &redis.Z{
			Member: cityID,
			Score:  float64(score),
		})
	}
}

// ImportCitiesToRedis 建立cityId-详细信息的映射
func ImportCitiesToRedis(conn *redis.Client, fileName string) error {
	csvFile, err := os.Open(fileName)
	if err != nil {
		fmt.Println("open file err", err)
		return err
	}
	reader := csv.NewReader(bufio.NewReader(csvFile))

	for {
		line, err := reader.Read()
		if err == io.EOF {
			break
		}
		if line[0] == "cityid" {
			continue
		}
		cityDetail, err := json.Marshal(CityDetailed{
			CityID:   line[0],
			CityName: line[3],
			Region:   line[2],
		})
		if err != nil {
			return err
		}
		err = conn.HSet(context.Background(), "cityId2city", line[0], string(cityDetail)).Err()
		if err != nil {
			return err
		}
	}

	return nil

}

// IpToScore 将ip地址转化为整数值
func IpToScore(ipAddress string) int {
	score := 0
	for _, v := range strings.Split(ipAddress, ".") {
		vi, err := strconv.Atoi(v)
		if err != nil {

		}
		score = score*256 + vi
	}

	return score
}
