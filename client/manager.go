package client

import (
	"errors"
	"log"
	"sort"
	"sync"

	"github.com/robfig/cron/v3"
)

type Manager struct {
	mutex sync.Mutex // 互斥锁

	clients  map[*Client]map[string]struct{} // 保存所有客户端集合
	handlers map[string]Handler              // 保存所有的topic的控制处理器
	cron     *cron.Cron                      // 定时任务
}

func (m *Manager) AddHandler(topic string, h Handler) error {
	if _, ok := m.handlers[topic]; ok {
		return errors.New("handler already exist,don't repeat register")
	}
	m.handlers[topic] = h
	return nil
}

func NewManager() *Manager {
	return &Manager{
		clients:  make(map[*Client]map[string]struct{}),
		cron:     cron.New(cron.WithSeconds()),
		handlers: make(map[string]Handler),
	}
}

func NewDefaultManager() (*Manager, error) {
	manager := NewManager()
	if err := manager.AddCronFunc("0/5 * * * * *", manager.Logger); err != nil {
		return nil, err
	}
	return manager, nil
}

// AddCronFunc 添加定时任务
func (m *Manager) AddCronFunc(spec string, cmds ...func()) error {
	for _, cmd := range cmds {
		if _, err := m.cron.AddFunc(spec, cmd); err != nil {
			return err
		}
	}

	return nil
}

func (m *Manager) Logger() {
	defer func() {
		if err := recover(); err != nil {
			log.Println("panic error: ", err)
			return
		}
	}()

	log.Printf("当前存活的客户端数量：%d\n", len(m.clients))
	for client, topics := range m.clients {
		var ts []string
		for topic := range topics {
			ts = append(ts, topic)
		}
		sort.Strings(ts)
		log.Println(client.conn.RemoteAddr(), "topic:", ts)
	}
}

// StartCron 启动定时任务
func (m *Manager) StartCron() {
	m.cron.Start()
}

// AddClient 添加客户端
func (m *Manager) AddClient(c *Client) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if _, ok := m.clients[c]; ok {
		return errors.New("目标客户端连接已存在")
	}

	m.clients[c] = make(map[string]struct{})

	return nil
}

// RemoveClient 移除客户端
func (m *Manager) RemoveClient(c *Client) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	// 移除客户端
	delete(m.clients, c)

	log.Println("客户端已断开", c.conn.RemoteAddr())
}

// AddTopic 为客户端添加订阅的主题
func (m *Manager) AddTopic(c *Client, topic string) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	// 为客户端添加主题
	if client, ok := m.clients[c]; ok {
		client[topic] = struct{}{}
	}

	c.SubscribeTopic(topic)
}

// RemoveTopic 为客户端移除订阅的主题
func (m *Manager) RemoveTopic(c *Client, topic string) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	// 为客户端移除主题
	delete(m.clients[c], topic)

	c.UnSubscribeTopic(topic)
}

func (m *Manager) Broadcast(topic string, message any) {
	// 如果topic为空，则表示全局广播，反之则只广播给订阅了的该topic的客户端
	if topic == "" {
		for client := range m.clients {
			client.SendMessage(message)
		}
	}
	for client, topics := range m.clients {
		for item := range topics {
			if item == topic {
				client.SendMessage(message)
			}
		}
	}
}

func (m *Manager) GetHandler(topic string) Handler {
	return m.handlers[topic]
}
