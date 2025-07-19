// 事件监听器模块
// 集中管理各种DOM事件监听器的设置

import { updateClientConfig, sendOfflineCommand } from './commands.js';
import { sendCommand } from './websocket.js';

/**
 * 初始化命令按钮事件监听器
 */
export function initCommandListeners() {
  // 暂停命令按钮
  document.getElementById('pauseBtn')?.addEventListener('click', () => {
    sendCommand({ type: 'pause_screenshots' });
  });

  // 继续命令按钮
  document.getElementById('resumeBtn')?.addEventListener('click', () => {
    sendCommand({ type: 'resume_screenshots' });
  });

  // 下线命令按钮
  document.getElementById('offlineBtn')?.addEventListener('click', () => {
    if (confirm('确定要让客户端下线吗？这将断开客户端连接。')) {
      sendOfflineCommand();
    }
  });

  // 立即截图按钮
  document.getElementById('screenshotBtn')?.addEventListener('click', () => {
    sendCommand({ type: 'take_screenshot' });
  });

  // 更新配置按钮事件监听
  document.getElementById('updateConfigBtn').onclick = updateClientConfig;
}

/**
 * 初始化所有事件监听器
 */
export function initAllEventListeners() {
  initCommandListeners();
}
