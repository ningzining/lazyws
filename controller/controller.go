package controller

import (
	"log"
	"net/http"

	"lazyws/client"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{} // use default options

func Test(manager *client.Manager, c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("upgrade:", err)
		c.String(http.StatusInternalServerError, "无法升级到WebSocket连接: %v", err)
		return
	}
	defer conn.Close()

	cl := client.New(conn)
	if err := manager.AddClient(cl); err != nil {
		return
	}

	cl.Run(manager)
}
