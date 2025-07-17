// 认证模块
// 处理前端登录、登出、状态检查等所有认证相关功能

import { buildUrl } from './path-utils.js';

/**
 * 检查登录状态
 * @returns {Promise<boolean>} 是否已登录
 */
export async function checkAuthStatus() {
  try {
    const response = await fetch(buildUrl('/auth/session'), {
      credentials: 'include'
    });
    const sessionInfo = await response.json();

    if (!sessionInfo.authEnabled) {
      // 认证被禁用
      return true;
    }

    if (!sessionInfo.authenticated) {
      // 未登录，重定向到登录页
      window.location.href = buildUrl('/login');
      return false;
    }

    return true;
  } catch (error) {
    console.error('检查登录状态失败:', error);
    // 发生错误时重定向到登录页
    window.location.href = buildUrl('/login');
    return false;
  }
}

/**
 * 登出
 */
export async function logout() {
  try {
    const response = await fetch(buildUrl('/auth/logout'), {
      method: 'POST',
      credentials: 'include'
    });

    if (response.ok) {
      // 登出成功，重定向到登录页
      window.location.href = buildUrl('/login');
    } else {
      console.error('登出失败');
    }
  } catch (error) {
    console.error('登出时出错:', error);
    // 即使出错也重定向到登录页
    window.location.href = buildUrl('/login');
  }
}

/**
 * 初始化认证相关功能（主应用使用）
 */
export function initAuth() {
  // 检查登录状态
  checkAuthStatus();

  // 设置登出按钮事件
  const logoutBtn = document.getElementById('logoutBtn');
  if (logoutBtn) {
    logoutBtn.addEventListener('click', logout);
  }
}

// ===== 登录页面相关功能 =====

/**
 * 显示错误信息
 */
function showError(message) {
  const errorMessage = document.getElementById('errorMessage');
  errorMessage.textContent = message;
  errorMessage.style.display = 'block';
}

/**
 * 隐藏错误信息
 */
function hideError() {
  const errorMessage = document.getElementById('errorMessage');
  errorMessage.style.display = 'none';
}

/**
 * 设置加载状态
 */
function setLoading(isLoading) {
  const loginBtn = document.getElementById('loginBtn');
  const loading = document.getElementById('loading');

  loginBtn.disabled = isLoading;
  loading.style.display = isLoading ? 'block' : 'none';
  loginBtn.textContent = isLoading ? '登录中...' : '登录';
}

/**
 * 处理登录表单提交
 */
async function handleLogin(event) {
  event.preventDefault();
  hideError();
  setLoading(true);

  const formData = new FormData(event.target);
  const loginData = {
    username: formData.get('username'),
    password: formData.get('password')
  };

  try {
    const response = await fetch(buildUrl('/auth/login'), {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(loginData),
      credentials: 'include'
    });

    if (response.ok) {
      // 登录成功，重定向到主页
      window.location.href = buildUrl('/');
    } else {
      const error = await response.json();
      showError(error.error || '登录失败');
    }
  } catch {
    showError('网络错误，请重试');
  } finally {
    setLoading(false);
  }
}

/**
 * 检查当前登录状态（登录页面使用）
 */
async function checkLoginStatusForLoginPage() {
  try {
    const response = await fetch(buildUrl('/auth/session'), {
      credentials: 'include'
    });
    const sessionInfo = await response.json();

    if (sessionInfo.authenticated) {
      // 已经登录，重定向到主页
      window.location.href = buildUrl('/');
    }
  } catch {
    // 忽略错误，继续显示登录页面
  }
}

/**
 * 初始化登录页面
 */
export function initLogin() {
  // 绑定表单提交事件
  const loginForm = document.getElementById('loginForm');
  if (loginForm) {
    loginForm.addEventListener('submit', handleLogin);
  }

  // 检查登录状态
  checkLoginStatusForLoginPage();
}

// 页面加载完成后自动初始化（仅在登录页面）
document.addEventListener('DOMContentLoaded', () => {
  // 检查是否在登录页面（通过是否存在登录表单来判断）
  if (document.getElementById('loginForm')) {
    initLogin();
  }
});
