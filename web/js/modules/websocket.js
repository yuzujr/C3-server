// WebSocket 通信模块
// 处理实时通知和WebSocket连接管理

import { selectedClient, setWebSocket, webSocket } from './state.js';
import { addNewScreenshot } from './screenshots.js';
import { handlePtyShellOutput } from './pty-terminal.js';
import { handleClientStatusChange } from './clients.js';
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
        addNewScreenshot(data.screenshot_url);
      }
      break;

    case 'shell_output':
      if (data.client_id === selectedClient) {
        handlePtyShellOutput(data);
      }
      break;

    case 'client_status_change':
      handleClientStatusChange(data.client_id, data.online);
      break;

    case 'alias_updated':
      handleAliasUpdated(data.client_id, data.oldAlias, data.newAlias);
      break;

    case 'client_deleted':
      handleClientDeleted(data.client_id);
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
  // 导入需要的模块
  const { selectedClient, cachedClientList, setCachedClientList, setSkipNextDOMRebuild } = await import('./state.js');

  console.log(`别名更新: ${oldAlias} -> ${newAlias} (clientId: ${clientId})`);

  // 设置标志，让接下来的定时刷新跳过DOM重建
  setSkipNextDOMRebuild(true);

  // 立即更新缓存中的客户端别名
  const updatedClients = cachedClientList.map(client => {
    if (client.client_id === clientId) {
      return { ...client, alias: newAlias };
    }
    return client;
  });
  setCachedClientList(updatedClients);

  // 立即更新DOM中的别名显示
  const clientElement = document.querySelector(`[data-client-id="${clientId}"]`);
  if (clientElement) {
    const nameElement = clientElement.querySelector('.client-name');
    if (nameElement) {
      nameElement.textContent = newAlias;
    }
  }

  // 如果当前选中的客户端被重命名了，更新选中状态显示
  if (selectedClient === clientId) {
    // 更新客户端管理列表中的别名显示
    import('./client-management.js').then(({ updateClientManagementList }) => {
      updateClientManagementList(updatedClients);
    });
  }
}

/**
 * 处理客户端删除事件
 * @param {string} clientId - 被删除的客户端ID
 */
async function handleClientDeleted(clientId) {

  // 导入需要的模块
  const { selectedClient, cachedClientList, setCachedClientList, setSelectedClient } = await import('./state.js');
  const { fetchClients } = await import('./clients.js');

  // 从缓存中移除该客户端
  const updatedClients = cachedClientList.filter(client => client.client_id !== clientId);
  setCachedClientList(updatedClients);

  // 从DOM中移除客户端元素
  const clientElement = document.querySelector(`[data-client-id="${clientId}"]`);
  if (clientElement) {
    clientElement.remove();
  }

  // 如果删除的是当前选中的客户端，清除选中状态
  if (selectedClient === clientId) {
    setSelectedClient(null);
    // 清除右侧面板
    const rightPanel = document.querySelector('.right-panel');
    if (rightPanel) {
      rightPanel.innerHTML = '<div class="no-client-selected">请选择一个客户端</div>';
    }
  }

  // 显示删除成功的提示

  // 立即刷新客户端列表以确保同步
  setTimeout(() => {
    fetchClients();
  }, 100);
}
