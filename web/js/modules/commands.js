// 命令系统模块
// 处理向客户端发送命令和配置管理

import { selectedClient } from './state.js';
import { showError, showWarning } from '../../toast/toast.js';
import { buildUrl } from './path-utils.js';
import { sendCommand } from './websocket.js';


/**
 * 加载客户端配置
 * @param {string} clientId - 客户端ID
 */
export async function loadClientConfig(clientId) {
  try {
    const res = await fetch(buildUrl(`/web/config/${clientId}`));
    if (!res.ok) {
      console.warn('获取配置失败');
      console.warn(await res.text());
      return;
    }
    const config = await res.json();

    if (config.api) {
      document.getElementById('hostname').value = config.api.hostname || '';
      document.getElementById('port').value = config.api.port || '';
      document.getElementById('intervalSeconds').value = config.api.interval_seconds ?? '';
      document.getElementById('maxRetries').value = config.api.max_retries ?? '';
      document.getElementById('retryDelayMs').value = config.api.retry_delay_ms ?? '';
      document.getElementById('addToStartup').checked = !!config.api.add_to_startup;
    }
  } catch (err) {
    console.error('加载客户端配置出错:', err);
  }
}

/**
 * 更新客户端配置
 */
export async function updateClientConfig() {
  if (!selectedClient) {
    showWarning('请先选择客户端');
    return;
  }

  const hostname = document.getElementById('hostname').value.trim();
  const port = parseInt(document.getElementById('port').value);
  const intervalSeconds = parseInt(document.getElementById('intervalSeconds').value);
  const maxRetries = parseInt(document.getElementById('maxRetries').value);
  const retryDelayMs = parseInt(document.getElementById('retryDelayMs').value);
  const addToStartup = document.getElementById('addToStartup').checked;

  // 验证必填字段
  if (!hostname || isNaN(port) || isNaN(intervalSeconds) || isNaN(maxRetries) || isNaN(retryDelayMs)) {
    showError('请填写正确的配置参数');
    return;
  }

  // 验证端口范围
  if (port < 1 || port > 65535) {
    showError('端口号必须在1-65535之间');
    return;
  }

  // 验证间隔时间
  if (intervalSeconds < 1) {
    showError('截图间隔必须大于0秒');
    return;
  }

  const newConfig = {
    api: {
      hostname: hostname,
      port: port,
      interval_seconds: intervalSeconds,
      max_retries: maxRetries,
      retry_delay_ms: retryDelayMs,
      add_to_startup: addToStartup,
    }
  };

  const cmd = {
    type: "update_config",
    data: newConfig
  };

  sendCommand(cmd);
}

/**
 * 发送下线命令给客户端
 */
export async function sendOfflineCommand() {
  if (!selectedClient) {
    showWarning('请先选择客户端');
    return false;
  }

  sendCommand({
    type: 'offline',
    data: {
      reason: 'User requested offline',
      timestamp: new Date().toISOString()
    }
  }, false);
}
