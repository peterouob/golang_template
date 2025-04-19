package router

import (
	"github.com/gin-gonic/gin"
	"github.com/peterouob/golang_template/pkg/http/user"
	"github.com/peterouob/golang_template/utils"
)

func InitRouter(r *gin.Engine) {
	r.Use(utils.Cors)
	u := r.RouterGroup
	{
		u.POST("/login", user.LoginUser)
		u.POST("/register", user.RegisterUser)

	}
}
