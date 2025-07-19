package websocket

import "github.com/gorilla/websocket"

type Role string

const (
	RoleAgent Role = "agent" // C++ 客户端
	RoleUser  Role = "user"  // 前端页面
)

type client struct {
	Conn *websocket.Conn
	Send chan []byte
	ID   string
	Role Role
}

type Command struct {
	Type      string         `json:"type"`
	Data      map[string]any `json:"data"`
}