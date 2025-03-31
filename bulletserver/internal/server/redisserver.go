package server

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/peterhuang621/CelebratoryBullets/bulletserver/configs"
	"github.com/peterhuang621/CelebratoryBullets/bulletserver/pkg"
	goredis "github.com/redis/go-redis/v9"
)

func initRedis() {
	redisClient = goredis.NewClient(&goredis.Options{
		Addr:     configs.Redis_Addr,
		Password: "",
		DB:       0,
	})

	ctx := context.Background()
	getCount, err := redisClient.Get(ctx, configs.Redis_Key).Result()
	if err == goredis.Nil {
		log.Printf("Redis server: [INFO] Key %s not found, set to 0\n", configs.Redis_Key)

		err = redisClient.Set(ctx, configs.Redis_Key, getCount, configs.Redis_Key_DefaultExpr).Err()

		pkg.FailOnError(err, "Redis server: [ERROR] Failed to set key!")
	} else if err != nil {
		pkg.FailOnError(err, "Redis server: [ERROR] Failed to get key!")
		return
	} else {
		num, err := strconv.Atoi(getCount)
		pkg.FailOnError(err, "Redis server: [ERROR] Key is not an vaild integar!")

		log.Printf("Redis server: [INFO] Old %s=%d will be Set to 0\n", configs.Redis_Key, num)
		err = redisClient.Set(ctx, configs.Redis_Key, 0, configs.Redis_Key_DefaultExpr).Err()
		pkg.FailOnError(err, "Redis server: [ERROR] Failed to set key!")

		log.Printf("Redis server: [INFO] Key %s not found, set to 0\n", configs.Redis_Key)
	}
}

func getKey() (int, error) {
	ctx := context.Background()
	getCount, err := redisClient.Get(ctx, configs.Redis_Key).Result()
	if err != nil {
		return -1, fmt.Errorf("Redis server: [ERROR] Get error! %v\n", err)
	}
	currCount, err := strconv.Atoi(getCount)
	if err != nil {
		return -1, fmt.Errorf("Redis server: [ERROR] Failed to convert the getCount into a number, which hardly occurs! %v\n", err)
	}
	return currCount, nil
}

func isAllowed(count int, window time.Duration) error {
	ctx := context.Background()
	currCount, err := getKey()
	if err != nil {
		return err
	}

	currCount += count
	if currCount > configs.DrawingSpeed_Max {
		return fmt.Errorf("Redis server: [ERROR] Failed to handle your bullets num=%d, as the server has already %d bullet(s) in waiting! %v\n", count, currCount-count, err)
	} else if currCount > configs.DrawingSpeed_Heavy {
		log.Printf("Redis server: [WARN] PRESSURE-Heavy now!\n")
	} else if currCount > configs.DrawingSpeed_Light {
		log.Printf("Redis server: [WARN] PRESSURE-Light now!\n")
	}

	_, err = redisClient.IncrBy(ctx, configs.Redis_Key, int64(count)).Result()
	if err != nil {
		return fmt.Errorf("Redis server: [ERROR] Incrby %d error! %v\n", count, err)
	}
	if count == 0 {
		redisClient.Expire(ctx, configs.Redis_Key, window)
		log.Printf("Redis server: [INFO] First time to require service at Redis server, so setting expiration time to %d\n", window)
	}
	return nil
}

func consumeKey(count int) error {
	ctx := context.Background()
	_, err = redisClient.DecrBy(ctx, configs.Redis_Key, int64(count)).Result()
	if err != nil {
		return fmt.Errorf("Redis server: [ERROR] Decrby %d error! %v\n", count, err)
	}
	return nil
}
