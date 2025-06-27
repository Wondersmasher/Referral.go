package main

import (
	"github.com/Wondersmasher/Referral/env"
	mongodb "github.com/Wondersmasher/Referral/mongoDb"
	redisCache "github.com/Wondersmasher/Referral/redisCache"
	"github.com/Wondersmasher/Referral/route"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	server := gin.Default()
	gin.SetMode(gin.ReleaseMode)

	env.Env()
	mongodb.InitDb()
	redisCache.InitRedis()
	route.RegisterAllRoutes(server)
}
