package router

import (
	"github.com/gin-gonic/gin"
	"github.com/peterouob/golang_template/pkg/http/user"
)

func InitRouter(r *gin.Engine) {
	u := r.RouterGroup
	{
		u.POST("/login", user.LoginUser)
		u.POST("/register", user.RegisterUser)
	}
}
