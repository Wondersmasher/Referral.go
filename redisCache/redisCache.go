package redisCache

import (
	"context"
	"fmt"

	"github.com/Wondersmasher/Referral/env"
	"github.com/redis/go-redis/v9"
)

func InitRedis() {
	ctx := context.Background()

	rdb := redis.NewClient(&redis.Options{
		Addr:     env.REDIS_ADDRESS,
		Username: env.REDIS_USERNAME,
		Password: env.REDIS_PASSWORD,
		DB:       int(env.REDIS_DB),
	})

	v := rdb.Set(ctx, "foo", "bar", 60*3*60)
	fmt.Println(v, "V rdb")
	result, err := rdb.Get(ctx, "foo").Result()

	if err != nil {
		fmt.Println("Entered here ooo")
		panic(err)
	}

	fmt.Println(result) // >>> bar
}
