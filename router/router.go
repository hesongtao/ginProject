package router

import (
	"ginProject/app/controller/other"
	"ginProject/app/controller/user"
	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	router := gin.Default()
	//router.GET("/", websocket.Web WebsocketManager.WsClient)

	router.GET("/user/login", user.Login)

	router.GET("/user/info", user.Info)

	router.GET("/other/wash", other.Wash)

	router.GET("/friend/smile", user.Smile)

	return router
}
