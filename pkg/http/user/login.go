package user

import (
	"github.com/gin-gonic/gin"
	mdb "github.com/peterouob/golang_template/pkg/db/mysql"
	"github.com/peterouob/golang_template/pkg/grpc/client"
	"github.com/peterouob/golang_template/utils"
	"net/http"
)

func LoginUser(c *gin.Context) {
	var model mdb.UserModel
	if err := c.BindJSON(&model); err != nil {
		utils.Error("error in bind json for msg", err)
		//TODO:發生錯誤時不結束重試機制
		return
	}
	resp := client.LoginUserGrpc("8082", c, model)
	if resp == nil {
		c.JSON(http.StatusBadGateway, gin.H{
			"msg": "error in grpc client",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"resp": resp,
	})
}

func RegisterUser(c *gin.Context) {
	var model mdb.UserModel
	if err := c.BindJSON(&model); err != nil {
		utils.Error("error in bind json for msg", err)
		//TODO:發生錯誤時不結束重試機制
		return
	}
	resp := client.RegisterUser("8085", c, model)
	if resp == nil {
		c.JSON(http.StatusBadGateway, gin.H{
			"msg": "error in grpc client",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"resp":  resp,
		"model": model,
	})
}
