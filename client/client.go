package client

import (
	"encoding/json"
	"sync"

	"lazyws/message"

	"github.com/gorilla/websocket"
)

// Client 客户端结构
type Client struct {
	conn   *websocket.Conn     // 客户端的连接
	topics map[string]struct{} // 客户端订阅的主题列表

	mutex     sync.Mutex
	writeChan chan any // 写通道
}

func New(conn *websocket.Conn) *Client {
	return &Client{
		conn:      conn,
		topics:    make(map[string]struct{}),
		writeChan: make(chan any),
	}
}

// SubscribeTopic 订阅主题
func (c *Client) SubscribeTopic(topic string) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.topics[topic] = struct{}{}
}

// UnSubscribeTopic 取消订阅主题
func (c *Client) UnSubscribeTopic(topic string) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	delete(c.topics, topic)
}

func (c *Client) SendMessage(message any) {
	c.writeChan <- message
}

// Run 运行客户端
func (c *Client) Run(manager *Manager) {
	go c.readLoop(manager)
	go c.writeLoop()

	// 阻塞，防止客户端退出
	select {}
}

// readLoop 读循环
func (c *Client) readLoop(manager *Manager) {
	for {
		// 读取消息
		messageType, bytes, err := c.conn.ReadMessage()
		if err != nil {
			break
		}
		switch messageType {
		case websocket.PingMessage: // 心跳包
			if err := c.conn.WriteMessage(websocket.PongMessage, []byte{}); err != nil {
				continue
			}
		case websocket.TextMessage: // 文本消息
			req := new(message.Request)
			if err := json.Unmarshal(bytes, req); err != nil {
				continue
			}
			switch req.TopicType {
			case message.SubscribeTopicType: // 订阅主题
				manager.AddTopic(c, req.Topic)
				// 处理主题
				HandleTopic(manager, c, req)
			case message.UnsubscribeTopicType: // 取消订阅主题
				manager.RemoveTopic(c, req.Topic)
			}
		}
	}
	// 连接断开，移除连接，释放资源
	c.close(manager)
}

// writeLoop 写循环
func (c *Client) writeLoop() {
	for {
		select {
		case msg := <-c.writeChan:
			if err := c.conn.WriteJSON(msg); err != nil {
				return
			}
		}
	}
}

// close 关闭连接
func (c *Client) close(manager *Manager) {
	close(c.writeChan)
	manager.RemoveClient(c)
}
