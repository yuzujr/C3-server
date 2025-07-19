package websocket

import (
	"net"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

const (
	// 写控制帧超时时间
	writeWait = 10 * time.Second
	// 客户端必须在 pongWait 时间内回复 Pong，否则视为断开
	pongWait = 60 * time.Second
	// 发送心跳 Ping 的周期，应小于 pongWait
	pingPeriod = (pongWait * 9) / 10
	// TCP KeepAlive 周期
	tcpKeepAlivePeriod = 3 * time.Minute
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func CreateConn(c *gin.Context, id string, role Role) error {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return err
	}

	// 打开底层 TCP Keep-Alive
	if tcpConn, ok := conn.UnderlyingConn().(*net.TCPConn); ok {
		tcpConn.SetKeepAlive(true)
		tcpConn.SetKeepAlivePeriod(tcpKeepAlivePeriod)
	}

	client := &client{
		Conn: conn,
		Send: make(chan []byte, 256),
		ID:   id,
		Role: role,
	}

	HubInstance.register <- client

	go client.writePump()
	go client.readPump()
	return nil
}

func (c *client) readPump() {
	defer func() {
		HubInstance.unregister <- c
		c.Conn.Close()
	}()

	// 初始读超时
	c.Conn.SetReadDeadline(time.Now().Add(pongWait))
	// 收到 Pong 时重置读超时
	c.Conn.SetPongHandler(func(appData string) error {
		c.Conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			break
		}
		if c.Role == RoleAgent {
			handleAgentMessage(c.ID, message)
		}else{
			handleUserMessage(c.Send, message)
		}
	}
}

func (c *client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.Conn.Close()
	}()

	for {
		select {
		case msg, ok := <-c.Send:
			// 服务端主动关闭通道时，发个 Close 帧
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			// 普通业务消息
			if err := c.Conn.WriteMessage(websocket.TextMessage, msg); err != nil {
				return
			}

		case <-ticker.C:
			// 发送心跳 Ping
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
