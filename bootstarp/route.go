package bootstarp

import (
	"github.com/869413421/chatgpt-web/config"
	"github.com/869413421/chatgpt-web/routes"
	"github.com/gin-gonic/gin"
	"sync"
)

var router *gin.Engine
var once sync.Once

func SetUpRoute() {
	cnf := config.LoadConfig()
	accounts := gin.Accounts{}

	for _, user := range cnf.Users {
		accounts[user.Name] = user.Password
	}

	once.Do(func() {
		router = gin.Default()

		authorized := router.Group("/", gin.BasicAuth(accounts))

		routes.RegisterWebRoutes(authorized)
	})
}
