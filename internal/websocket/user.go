package websocket

import (
	"encoding/json"

	"github.com/yuzujr/C3/internal/logger"
)

func broadcastMsgToUsers(message []byte) {
	for client := range HubInstance.users {
		select {
		case client.Send <- message:
		default:
			logger.Errorf("Failed to broadcast message: channel full")
		}
	}
}

// handleUserMessage 解析前端发来的消息，只处理 type="command"
// cSend 即 client.Send，用于向用户写 ACK
func handleUserMessage(cSend chan []byte, raw []byte) {
	type userMsg struct {
		Type     string          `json:"type"`
		ClientID string          `json:"client_id"`
		Cmd      json.RawMessage `json:"cmd"`
	}

	// 构造 ACK 的基础字段
	ack := map[string]interface{}{
		"type":      "command_ack",
		"client_id": "",
		"success":   false,
		"message":   "",
	}

	// 1. 解析消息结构
	var msg userMsg
	if err := json.Unmarshal(raw, &msg); err != nil {
		ack["message"] = "无效的消息格式: " + err.Error()
		writeAck(cSend, ack)
		return
	}
	// 填充 client_id 到 ACK
	ack["client_id"] = msg.ClientID

	// 2. 只关心 command 类型
	if msg.Type != "command" {
		ack["message"] = "不支持的消息类型: " + msg.Type
		writeAck(cSend, ack)
		return
	}

	// 3. 解析具体命令
	var cmd Command
	if err := json.Unmarshal(msg.Cmd, &cmd); err != nil {
		ack["message"] = "解析命令失败: " + err.Error()
		writeAck(cSend, ack)
		return
	}

	// 4. 查找对应的 agent
	agent := HubInstance.agents[msg.ClientID]
	if agent == nil {
		ack["message"] = "目标客户端不存在或未连接"
		writeAck(cSend, ack)
		return
	}

	// 5. 转发给 agent
	sendCommandToAgent(agent, cmd)

	// 6. 高频命令不发送ACK
	if cmd.Type == "pty_input" || cmd.Type == "pty_resize" {
		return
	}

	// 7. 所有流程无误，标记成功并回 ACK
	ack["success"] = true
	ack["message"] = "命令已成功发送"
	writeAck(cSend, ack)
}

// writeAck 把 ack map 序列化后通过 cSend 通道发回前端
func writeAck(cSend chan []byte, ack map[string]interface{}) {
	if b, err := json.Marshal(ack); err == nil {
		cSend <- b
	}
}
