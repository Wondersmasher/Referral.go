package middleware

import (
	"fmt"
	"time"

	"github.com/Wondersmasher/Referral/env"
	"github.com/Wondersmasher/Referral/utils"
	"github.com/gin-gonic/gin"
)

func AuthMiddleware(c *gin.Context) {
	fmt.Println("AuthMiddleware")
	accessToken, err := c.Cookie("accessToken")
	if err != nil {
		c.AbortWithStatusJSON(401, utils.ApiErrorResponse("unauthenticated"))
		return
	}

	accessClaims, accessIsValid, err := utils.ValidateToken(accessToken, env.JWT_SECRET_ACCESS_KEY)
	if err != nil {
		c.AbortWithStatusJSON(401, utils.ApiErrorResponse("unauthenticated"))
		return
	}

	if !accessIsValid {
		refreshToken, err := c.Cookie("refreshToken")
		if err != nil {
			c.AbortWithStatusJSON(401, utils.ApiErrorResponse("unauthenticated"))
			return
		}

		if accessClaims.ExpiresAt != nil && time.Now().After(accessClaims.ExpiresAt.Time) {
			refreshClaims, refreshIsValid, refreshErr := utils.ValidateToken(refreshToken, env.JWT_SECRET_ACCESS_KEY)
			if refreshErr != nil {
				c.AbortWithStatusJSON(401, utils.ApiErrorResponse("unauthenticated"))
				return
			}
			if !refreshIsValid {
				c.AbortWithStatusJSON(401, utils.ApiErrorResponse("unauthenticated"))
				return
			}

			accessToken, err := utils.CreateNewToken(refreshClaims.Email, refreshClaims.UserName, refreshClaims.ReferralID, time.Now().Add(time.Minute*15), env.JWT_SECRET_ACCESS_KEY)
			if err != nil {
				c.AbortWithStatusJSON(401, utils.ApiErrorResponse("unauthenticated"))
				return
			}
			c.SetCookie("accessToken", accessToken, 60*15, "/", "localhost", false, true)
			return
		} else {
			c.AbortWithStatusJSON(401, utils.ApiErrorResponse("unauthenticated"))
			return
		}
	}

	c.Next()
}
