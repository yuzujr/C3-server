package websocket

import (
	"encoding/json"

	"github.com/yuzujr/C3/internal/logger"
)

func sendCommandToAgent(clientID string, message Command) {
	// 将命令转换为JSON格式
	messageBytes, err := json.Marshal(message)
	if err != nil {
		logger.Errorf("Failed to marshal message for client %s: %v", clientID, err)
		return
	}

	// 查找客户端并发送消息
	client, ok := HubInstance.agents[clientID]
	if !ok {
		logger.Errorf("Client %s not found", clientID)
		return
	}
	select {
	case client.Send <- messageBytes:
	default:
		logger.Errorf("Failed to send message to client %s: channel full", clientID)
	}
}

// Agent 发来的包，只关心 type 和 data.output
type agentMsg struct {
	Type string          `json:"type"`
	Data json.RawMessage `json:"data"`
}

// data 中仅包含 output
type shellData struct {
	Output string `json:"output"`
}

func handleAgentMessage(clientID string, message []byte) {
	logger.Infof("收到 Agent %s 消息: %s", clientID, message)

	// 解析最外层，只取 type 和 data
	var am agentMsg
	if err := json.Unmarshal(message, &am); err != nil {
		logger.Errorf("解析 agentMsg 失败: %v", err)
		return
	}

	// 只处理 shell_output
	if am.Type != "shell_output" {
		return
	}

	// 解析 data.output
	var sd shellData
	if err := json.Unmarshal(am.Data, &sd); err != nil {
		logger.Errorf("解析 shellData 失败: %v", err)
		return
	}

	// 构造前端需要的最简消息
	outMsg := struct {
		Type   string `json:"type"`
		Output string `json:"output"`
	}{
		Type:   am.Type,
		Output: sd.Output,
	}

	// 序列化并广播
	bs, err := json.Marshal(outMsg)
	if err != nil {
		logger.Errorf("序列化 outMsg 失败: %v", err)
		return
	}
	broadcastMsgToUsers(bs)
}
