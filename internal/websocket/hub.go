package websocket

import (
	"encoding/json"

	"github.com/yuzujr/C3/internal/eventbus"
	"github.com/yuzujr/C3/internal/logger"
	"github.com/yuzujr/C3/internal/models"
	"github.com/yuzujr/C3/internal/service"
)

type hub struct {
	agents     map[string]*client // 客户端连接
	users      map[*client]bool   // 前端用户连接
	register   chan *client
	unregister chan *client
}

func (h *hub) Broadcast(msg any) {
	data, err := json.Marshal(msg)
	if err != nil {
		logger.Errorf("Failed to marshal broadcast message: %v", err)
		return
	}
	for user := range HubInstance.users {
		select {
		case user.Send <- data:
		default:
			logger.Errorf("Failed to send message to user %s: channel full", user.ID)
		}
	}
}

func (h *hub) SendCommand(id string, cmd eventbus.Command) {
	c := h.agents[id]
	// 检查客户端角色是否为 Agent
	if c.Role != RoleAgent {
		logger.Errorf("Client %s is not an agent, cannot send command", c.ID)
		return
	}

	// 查找客户端
	client, ok := h.agents[c.ID]
	if !ok {
		logger.Errorf("Client %s not found", c.ID)
		return
	}

	// 序列化命令
	messageBytes, err := json.Marshal(cmd)
	if err != nil {
		logger.Errorf("Failed to marshal message for client %s: %v", c.ID, err)
		return
	}

	// 发送消息
	select {
	case client.Send <- messageBytes:
	default:
		logger.Errorf("Failed to send message to client %s: channel full", c.ID)
	}
}

var HubInstance = newHub()

func newHub() *hub {
	return &hub{
		agents:     make(map[string]*client),
		users:      make(map[*client]bool),
		register:   make(chan *client),
		unregister: make(chan *client),
	}
}

func init() {
	eventbus.Global = HubInstance
}

func (h *hub) Run() {
	logger.Infof("WebSocket hub started")
	for {
		select {
		case client := <-h.register:
			switch client.Role {
			case RoleAgent:
				h.handleAgentRegister(client)
			case RoleUser:
				h.handleUserRegister(client)
			}
		case client := <-h.unregister:
			switch client.Role {
			case RoleAgent:
				h.handleAgentUnregister(client)
			case RoleUser:
				h.handleUserUnregister(client)
			}
		}
	}
}

func (h *hub) handleAgentRegister(client *client) {
	// 加入到 agents 列表
	h.agents[client.ID] = client
	// 更新数据库状态
	service.UpsertClient(&models.Client{
		ClientID: client.ID,
		Online:   true,
	})
	// 广播客户端上线消息给所有用户
	msg := eventbus.StatusChangeMsg{
		Type:     "client_status_change",
		ClientID: client.ID,
		Online:   true,
	}
	eventbus.Global.Broadcast(msg)
	logger.Infof("Client %s Websocket connected", client.ID)
}

func (h *hub) handleAgentUnregister(client *client) {
	// 从 agents 列表中删除
	delete(h.agents, client.ID)
	// 更新数据库状态
	service.UpdateClient(&models.Client{
		ClientID: client.ID,
		Online:   false,
	})
	// 广播客户端下线消息给所有用户
	msg := eventbus.StatusChangeMsg{
		Type:     "client_status_change",
		ClientID: client.ID,
		Online:   false,
	}
	eventbus.Global.Broadcast(msg)
	close(client.Send)
	logger.Infof("Client %s Websocket disconnected", client.ID)
}

func (h *hub) handleUserRegister(client *client) {
	h.users[client] = true
	logger.Infof("User Websocket connected")
}

func (h *hub) handleUserUnregister(client *client) {
	delete(h.users, client)
	close(client.Send)
	logger.Infof("User Websocket disconnected")
}
