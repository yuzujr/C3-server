package eventbus

// 主要用于避免websocket和service包之间的循环依赖
type Messenger interface {
	Broadcast(msg any)                        //向所有用户广播消息
	SendCommand(clientID string, cmd Command) //向指定客户端发送命令
}

// 全局的广播器实例（在 websocket 包里初始化）
var Global Messenger
