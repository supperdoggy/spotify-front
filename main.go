package main

import (
	"github.com/gin-gonic/gin"
	"github.com/supperdoggy/webproject/frontend/internal/handlers"
	service2 "github.com/supperdoggy/webproject/frontend/internal/service"
	"go.uber.org/zap"
)

func main() {
	logger, _ := zap.NewDevelopment()
	r := gin.Default()
	// todo to config
	port := ":8081"

	r.LoadHTMLFiles("./src/static/html/play.html")
	r.Static("src/static", "./src/static")

	service := service2.NewService(logger)
	hadlers := handlers.NewHandlers(logger, r, &service)

	hadlers.InitHandlers()

	if err := r.Run(port); err != nil {
		logger.Fatal("error running thing", zap.Error(err))
	}
}
