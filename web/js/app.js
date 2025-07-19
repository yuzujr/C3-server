// 主应用程序模块
// 应用程序的入口点和初始化逻辑

import { CLIENT_POLL_INTERVAL } from './modules/state.js';
import { fetchClients } from './modules/clients.js';
import { initImageModal } from './modules/modal.js';
import { initDeleteScreenshots } from './modules/delete.js';
import { initThemeToggle } from './modules/theme.js';
import { initWebSocket } from './modules/websocket.js';
import { initAllEventListeners } from './modules/events.js';
import { initPtyTerminal } from './modules/pty-terminal.js';
import { initTabs } from './modules/tabs.js';

/**
 * 应用程序初始化函数
 */
async function initApp() {
  try {
    // 首先初始化认证，等待认证完成
    const { checkAuthStatus } = await import('./modules/auth.js');
    const isAuthenticated = await checkAuthStatus();

    if (!isAuthenticated) {
      // 如果未认证，保持加载遮罩显示直到跳转完成
      return;
    }

    // 认证成功后初始化其他功能
    const { initAuth } = await import('./modules/auth.js');
    initAuth();

    // 初始化 WebSocket 连接
    initWebSocket();

    // 获取客户端列表
    fetchClients();

    // 初始隐藏命令按钮
    document.getElementById('cmdButtons').style.display = 'none';        // 初始化各种功能模块
    initTabs();
    initPtyTerminal();
    initImageModal();
    initDeleteScreenshots();
    initThemeToggle();
    initAllEventListeners();

    // 初始化客户端管理
    const { initClientManagement } = await import('./modules/client-management.js');
    initClientManagement();

    // 所有初始化完成后，隐藏加载遮罩并显示主内容
    showMainContent();
  } catch (error) {
    console.error('应用初始化失败:', error);
    // 发生错误时也要显示主内容，避免永远卡在加载页面
    showMainContent();
  }
}

/**
 * 显示主内容并隐藏加载遮罩
 */
function showMainContent() {
  const loadingOverlay = document.getElementById('loadingOverlay');
  const mainContent = document.getElementById('mainContent');

  // 移除body的loading类，显示固定按钮
  document.body.classList.remove('loading');

  if (loadingOverlay) {
    loadingOverlay.style.display = 'none';
  }

  if (mainContent) {
    mainContent.style.display = 'block';
  }
}

// 当DOM加载完成时启动应用程序
window.addEventListener('load', initApp);
