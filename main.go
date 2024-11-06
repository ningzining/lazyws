package main

import (
	"log"

	"lazyws/client"
	"lazyws/client/base_handler"
	"lazyws/controller"
	"lazyws/server"

	"github.com/gin-gonic/gin"
)

func main() {
	s := server.New(8080)
	manager, err := client.NewDefaultManager()
	if err != nil {
		log.Println(err)
		return
	}
	if err := manager.AddHandler("aaa", base_handler.New()); err != nil {
		log.Println(err)
		return
	}
	manager.StartCron()

	s.AddHandler("/ws", func(c *gin.Context) {
		controller.Test(manager, c)
	})

	if err := s.Start(); err != nil {
		return
	}
}
