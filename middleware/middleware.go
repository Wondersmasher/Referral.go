package middleware

import "github.com/gin-gonic/gin"

func AuthMiddleware(c *gin.Context) {
	var user = map[string]any{
		"user":     "user",
		"password": "password",
		"email":    "email",
	}

	c.Set("user", user)
	c.Next()
}
