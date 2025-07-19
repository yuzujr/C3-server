// WebSocket 通信模块
// 处理实时通知和WebSocket连接管理

import { selectedClient, setWebSocket, webSocket } from './state.js';
import { addNewScreenshot } from './screenshots.js';
import { handlePtyShellOutput } from './pty-terminal.js';
import { fetchClients, handleClientStatusChange as handleClientOnlineStatusChange } from './clients.js';
import { buildWebSocketUrl } from './path-utils.js';
import { showSuccess, showError, showWarning } from '../../toast/toast.js';

/**
 * 初始化 WebSocket 连接
 */
export function initWebSocket() {
  const wsUrl = buildWebSocketUrl();

  const webSocket = new WebSocket(wsUrl);
  setWebSocket(webSocket);

  webSocket.onopen = () => {
    console.log('WebSocket连接已建立');
  };

  webSocket.onmessage = (event) => {
    try {
      const data = JSON.parse(event.data);
      handleWebSocketMessage(data);
    } catch (error) {
      console.error('WebSocket消息解析失败:', error);
    }
  };

  webSocket.onclose = (_event) => {
    console.log('WebSocket连接已关闭, 3秒后重连...');
    setTimeout(initWebSocket, 3000);
  };

  webSocket.onerror = (error) => {
    console.error('WebSocket连接错误:', error);
  };
}

/**
 * 通过 WebSocket 发送命令
 * @param {object} command - 命令对象
 */
export function sendCommand(command) {
  if (!selectedClient) {
    showWarning('请先选择客户端');
    return;
  }
  if (!webSocket || webSocket.readyState !== WebSocket.OPEN) {
    showError('WebSocket 未连接');
    return;
  }

  webSocket.send(JSON.stringify({
    type: 'command',
    client_id: selectedClient,
    cmd: command
  }));
}

/**
 * 处理 WebSocket 消息
 * @param {object} data - 接收到的消息数据
 */
function handleWebSocketMessage(data) {
  switch (data.type) {
    case 'new_screenshot':
      if (data.client_id === selectedClient) {
        addNewScreenshot(data.url);
      }
      break;

    case 'shell_output':
      if (data.client_id === selectedClient) {
        handlePtyShellOutput(data);
      }
      break;

    case 'client_status_change':
      handleClientOnlineStatusChange(data.client_id, data.online);
      break;

    case 'alias_updated':
      handleAliasUpdated(data.client_id, data.old_alias, data.new_alias);
      break;

    case 'command_ack':
      if (data.client_id === selectedClient) {
        if (data.success) {
          showSuccess('命令发送成功');
        } else {
          showError(`命令发送失败：${data.message}`);
        }
      }
      break;

    default:
      console.warn('未知消息类型：', data.type);
  }
}

/**
 * 处理别名更新事件
 * @param {string} clientId - 客户端ID
 * @param {string} oldAlias - 旧别名
 * @param {string} newAlias - 新别名
 */
async function handleAliasUpdated(clientId, oldAlias, newAlias) {
  console.log(`别名更新: ${oldAlias} -> ${newAlias} (clientId: ${clientId})`);
  fetchClients(); // 重新获取客户端列表以更新UI
}