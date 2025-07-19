package websocket

import (
	"encoding/json"

	"github.com/yuzujr/C3/internal/logger"
)

func sendCommandToAgent(c *client, cmd Command) {
	// 检查客户端角色是否为 Agent
	if c.Role != RoleAgent {
		logger.Errorf("Client %s is not an agent, cannot send command", c.ID)
		return
	}

	// 查找客户端
	client, ok := HubInstance.agents[c.ID]
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

func handleAgentMessage(clientID string, message []byte) {
	// Agent 发来的包
	type agentMsg struct {
		Type      string         `json:"type"`
		SessionID string         `json:"session_id"`
		Data      map[string]any `json:"data"`
	}

	// 发给 user 的包
	type userMsg struct {
		Type     string `json:"type"`
		Output   string `json:"output"`
		ClientID string `json:"client_id"`
	}

	logger.Debugf("收到 Agent %s 消息: %s", clientID, message)

	// 解析 agentMsg
	var am agentMsg
	if err := json.Unmarshal(message, &am); err != nil {
		logger.Errorf("解析 agentMsg 失败: %v", err)
		return
	}

	output, ok := am.Data["output"].(string)
	if !ok {
		// 不是pty发来的信息，忽略
		return
	}

	// 构造前端需要的消息
	outMsg := userMsg{
		Type:     am.Type,
		Output:   output,
		ClientID: clientID,
	}

	// 序列化并广播
	bs, err := json.Marshal(outMsg)
	if err != nil {
		logger.Errorf("序列化 outMsg 失败: %v", err)
		return
	}
	logger.Debugf("广播 Agent %s 消息: %s", clientID, bs)
	broadcastMsgToUsers(bs)
}
