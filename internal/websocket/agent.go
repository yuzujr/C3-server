package websocket

import (
	"encoding/json"

	"github.com/yuzujr/C3/internal/eventbus"
	"github.com/yuzujr/C3/internal/logger"
)

func handleAgentMessage(clientID string, message []byte) {
	// Agent 发来的包
	type agentMsg struct {
		Type      string         `json:"type"`
		SessionID string         `json:"session_id"`
		Data      map[string]any `json:"data"`
	}

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

	// 广播客户端终端输出消息给所有用户
	msg := eventbus.ShellOutputMsg{
		Type:     am.Type,
		Output:   output,
		ClientID: clientID,
	}
	eventbus.Global.Broadcast(msg)
}
