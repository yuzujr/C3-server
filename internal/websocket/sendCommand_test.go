package websocket

import (
	"encoding/json"
	"log"
	"os"
	"testing"
	"time"

	"github.com/joho/godotenv"
	"github.com/yuzujr/C3/internal/database"
	"github.com/yuzujr/C3/internal/eventbus"
	"github.com/yuzujr/C3/internal/models"
	"github.com/yuzujr/C3/internal/service"
)

func TestHubCommandSend(t *testing.T) {
	err := godotenv.Load("../../.env")
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
	os.Setenv("ENV", "test")
	database.InitDatabase()

	service.UpsertClient(&models.Client{
		ClientID:  "test-client",
		Alias:     "test-client",
		IPAddress: "",
		Online:    true,
		LastSeen:  time.Now(),
	})

	h := HubInstance
	go h.Run()

	c := &client{
		ID:   "test-client",
		Send: make(chan []byte, 10),
		Role: RoleAgent,
	}

	// 注册 client
	h.register <- c
	time.Sleep(10 * time.Millisecond)

	// 依次发送命令
	cmds := []eventbus.Command{
		{Type: "pause_screenshots"},
		{Type: "resume_screenshots"},
		{Type: "take_screenshot"},
		{
			Type: "pty_create_session",
			Data: map[string]any{
				"session_id": c.ID,
				"cols":       80,
				"rows":       24,
			},
		},
		{
			Type: "pty_resize",
			Data: map[string]any{
				"session_id": c.ID,
				"cols":       100,
				"rows":       40,
			},
		},
		{
			Type: "pty_input",
			Data: map[string]any{
				"session_id": c.ID,
				"input":      "echo Hello World\n",
			},
		},
		{
			Type: "pty_kill_session",
			Data: map[string]any{
				"session_id": c.ID,
			},
		},
	}

	go func() {
		for _, cmd := range cmds {
			eventbus.Global.SendCommand(c.ID, cmd)
			time.Sleep(10 * time.Millisecond) // 确保消息被处理
		}
	}()

	// 检查是否收到所有命令
	for i, want := range cmds {
		select {
		case msg := <-c.Send:
			var got eventbus.Command
			if err := json.Unmarshal(msg, &got); err != nil {
				t.Fatalf("unmarshal failed: %v", err)
			}
			if got.Type != want.Type {
				t.Errorf("command %d: want type %q, got %q", i, want.Type, got.Type)
			}
		case <-time.After(100 * time.Millisecond):
			t.Fatalf("timeout waiting for command %d", i)
		}
	}
}
