// 代码清单5-16没看明白什么意思
package main

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

func main() {

}

var LASTCHECKED = int64(0)
var ISUNDERMAINTENANCE = false

// IsUnderMaintenace
func IsUnderMaintenace(conn *redis.Client) (bool, error) {
	nSend := time.Now().Unix()
	if LASTCHECKED < nSend-120 {
		// 说明要更新配置了
		LASTCHECKED = nSend
		isMatain, err := conn.Get(context.Background(), "is-under-maintenance").Result()
		if err != nil {
			fmt.Println("redis get err", err)
			return false, err
		}
		if isMatain == "true" {
			ISUNDERMAINTENANCE = true
		} else {
			ISUNDERMAINTENANCE = false
		}
	}
	fmt.Println(ISUNDERMAINTENANCE)
	return ISUNDERMAINTENANCE, nil
}
