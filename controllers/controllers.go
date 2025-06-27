package controllers

import (
	// "fmt"
	"fmt"
	"time"

	"github.com/Wondersmasher/Referral/env"
	"github.com/Wondersmasher/Referral/model"
	"github.com/Wondersmasher/Referral/utils"
	"github.com/gin-gonic/gin"
)

type SignInStruct struct {
	Email    string `json:"email" bson:"email" validate:"required,email"`
	Password string `json:"password" bson:"password" validate:"required"`
}

func SignIn(c *gin.Context) {
	u := SignInStruct{}

	err := c.ShouldBindJSON(&u)
	if err != nil {
		c.JSON(400, utils.ApiErrorResponse("invalid credentials"))
		return
	}

	foundUser, err := model.GetUserByEmail(u.Email, u.Password)
	if err != nil {
		c.JSON(400, utils.ApiErrorResponse("invalid credentials"))
		return
	}

	accessToken, err := utils.CreateNewToken(foundUser.Email, foundUser.Username, time.Now().Add(time.Minute*15), env.JWT_SECRET_ACCESS_KEY)
	if err != nil {
		c.JSON(400, utils.ApiErrorResponse("could'nt generate accessToken"))
		return
	}

	refreshToken, err := utils.CreateNewToken(foundUser.Email, foundUser.Username, time.Now().Add(time.Hour*24*3), env.JWT_SECRET_REFRESH_KEY)
	if err != nil {
		c.JSON(400, utils.ApiErrorResponse("could'nt generate refreshToken"))
		return
	}
	c.SetCookie("accessToken", accessToken, int(time.Minute)*15, "/", "localhost", false, true)
	c.SetCookie("refreshToken", refreshToken, int(time.Hour)*24*3, "/", "localhost", false, true)

	// fmt.Println(accessToken, refreshToken)
	c.JSON(200, utils.ApiSuccessResponse(foundUser.TrimUser(), "success"))
}

func SignOut(c *gin.Context) {
	cookie, err := c.Cookie("accessToken")
	fmt.Println(cookie, err)
	c.SetCookie("accessToken", "", -1, "/", "localhost", false, true)
	c.SetCookie("refreshToken", "", -1, "/", "localhost", false, true)
	c.JSON(200, utils.ApiSuccessResponse(nil, "success"))
}

func SignUp(c *gin.Context) {

}

func GetReferrals(c *gin.Context) {

}
