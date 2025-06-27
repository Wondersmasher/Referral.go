package controllers

import (
	"github.com/Wondersmasher/Referral/model"
	"github.com/Wondersmasher/Referral/utils"
	"github.com/gin-gonic/gin"
)

type LoginStruct struct {
	Email    string `json:"email" bson:"email" validate:"required,email"`
	Password string `json:"password" bson:"password" validate:"required"`
}

func Test(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "Hello, World!",
	})
}

func Login(c *gin.Context) {
	u := LoginStruct{}

	err := c.ShouldBindJSON(&u)
	if err != nil {
		c.JSON(400, utils.ApiErrorResponse("invalid credentials"))
		return
	}

	_, err = model.GetUserByEmail(u.Email, u.Password)
	if err != nil {
		c.JSON(400, utils.ApiErrorResponse("invalid credentials"))
		return
	}

	c.JSON(200, gin.H{
		"message": "Hello, World!",
	})
}
