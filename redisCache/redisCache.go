package redisCache

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

func InitRedis() {
	ctx := context.Background()

	rdb := redis.NewClient(&redis.Options{
		Addr:     "redis-14747.c57.us-east-1-4.ec2.redns.redis-cloud.com:14747",
		Username: "default",
		Password: "itx9LvQ1LYkKupGWBhrwEDdQAezwMRtF",
		DB:       0,
	})

	rdb.Set(ctx, "foo", "bar", 60*3*60)
	result, err := rdb.Get(ctx, "foo").Result()

	if err != nil {
		panic(err)
	}

	fmt.Println(result) // >>> bar
}
