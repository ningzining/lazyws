package base_handler

import (
	"log"

	"lazyws/client"
	"lazyws/message"
)

type BaseHandler struct {
}

func New() client.Handler {
	return &BaseHandler{}
}

func (b *BaseHandler) Handle(c *client.Client, msg *message.Request) {
	log.Println(msg)
	c.SendMessage("hello client")
}
