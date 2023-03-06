package bootstarp

import (
	"github.com/869413421/chatgpt-web/config"
	"github.com/869413421/chatgpt-web/pkg/logger"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func StartWebServer() {

	gin.SetMode(gin.ReleaseMode)

	// 注册启动所需各类参数
	SetUpRoute()
	initTemplateDir()
	initStaticServer()

	// 启动服务
	port := config.LoadConfig().Port
	portString := strconv.Itoa(port)
	// err := router.Run(":" + portString)
	// openssl req -newkey rsa:4096 -nodes -keyout server.key -out server.csr
	// openssl x509 -signkey server.key -in server.csr -req -days 365 -out server.crt
	err := router.RunTLS(":"+portString, "./server.crt", "./server.key")
	if err != nil {
		logger.Danger("run webserver error %s", err)
		return
	}
}

// initTemplate 初始化HTML模板加载路径
func initTemplateDir() {
	router.LoadHTMLGlob("resources/view/*")
}

// initStaticServer 初始化静态文件处理
func initStaticServer() {
	router.StaticFS("/static", http.Dir("static"))
	router.StaticFile("logo192.png", "static/logo192.png")
	router.StaticFile("logo512.png", "static/logo512.png")
}
