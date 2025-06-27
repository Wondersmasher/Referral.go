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

	if u.Email == "" || u.Password == "" {
		c.JSON(400, utils.ApiErrorResponse("email and password are required"))
		return
	}

	foundUser, err := model.GetUserByEmail(u.Email, u.Password)
	if err != nil {
		c.JSON(400, utils.ApiErrorResponse("invalid credentials"))
		return
	}

	accessToken, err := utils.CreateNewToken(foundUser.Email, foundUser.Username, foundUser.ReferralID, time.Now().Add(time.Minute*15), env.JWT_SECRET_ACCESS_KEY)
	if err != nil {
		c.JSON(400, utils.ApiErrorResponse("could'nt generate accessToken"))
		return
	}

	refreshToken, err := utils.CreateNewToken(foundUser.Email, foundUser.Username, foundUser.ReferralID, time.Now().Add(time.Hour*24*3), env.JWT_SECRET_REFRESH_KEY)
	if err != nil {
		c.JSON(400, utils.ApiErrorResponse("could'nt generate refreshToken"))
		return
	}
	c.SetCookie("accessToken", accessToken, int(time.Minute)*15, "/", "localhost", false, true)
	c.SetCookie("refreshToken", refreshToken, int(time.Hour)*24*3, "/", "localhost", false, true)

	// fmt.Println(accessToken, refreshToken)
	c.JSON(200, utils.ApiSuccessResponse(foundUser.TrimUser(false), "success"))
}

func SignOut(c *gin.Context) {
	cookie, err := c.Cookie("accessToken")
	fmt.Println(cookie, err)
	c.SetCookie("accessToken", "", -1, "/", "localhost", false, true)
	c.SetCookie("refreshToken", "", -1, "/", "localhost", false, true)
	c.JSON(200, utils.ApiSuccessResponse(nil, "success"))
}

func SignUp(c *gin.Context) {
	user := model.User{}

	err := c.ShouldBindJSON(&user)
	if err != nil {
		c.JSON(400, utils.ApiErrorResponse(err.Error()))
		return
	}
	if user, err = user.CreateUser(user.ReferredBy); err != nil {
		c.JSON(400, utils.ApiErrorResponse(err.Error()))
		return
	}

	c.JSON(200, utils.ApiSuccessResponse(user.TrimUser(false), "success"))
}

func GetReferrals(c *gin.Context) {
	referredBy := c.Param("referredBy")
	accessToken, err := c.Cookie("accessToken")

	if err != nil {
		c.JSON(400, utils.ApiErrorResponse(err.Error()))
		return
	}

	claims, isValid, err := utils.ValidateToken(accessToken, env.JWT_SECRET_ACCESS_KEY)
	if !isValid {
		c.JSON(400, utils.ApiErrorResponse("invalid token"))
		return
	}
	if err != nil {
		c.JSON(400, utils.ApiErrorResponse(err.Error()))
		return
	}

	referrals, err := model.GetReferrals(referredBy)
	if err != nil {
		c.JSON(400, utils.ApiErrorResponse(err.Error()))
		return
	}
	fmt.Println(claims, isValid, err)
	c.JSON(200, utils.ApiSuccessResponse(referrals, "success"))

}
