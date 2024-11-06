package client

import (
	"lazyws/message"
)

type Handler interface {
	Handle(c *Client, msg *message.Request)
}

func HandleTopic(manager *Manager, client *Client, msg *message.Request) {
	if handler := manager.GetHandler(msg.Topic); handler != nil {
		handler.Handle(client, msg)
	}
}
