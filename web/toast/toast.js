// toast.js - 现代化Toast提示系统
// 支持多个toast同时显示、队列管理、不同类型的提示

let toastCounter = 0;
const activeToasts = new Map();

/**
 * 显示Toast提示框
 * @param {string} message - 提示消息
 * @param {string} type - 提示类型: 'success', 'error', 'warning', 'info'
 * @param {number} duration - 显示时长(毫秒)，默认3000
 * @param {boolean} closable - 是否可手动关闭，默认true
 */
export function showToast(message, type = 'info', duration = 3000, closable = true) {
  // 生成唯一ID
  const toastId = `toast-${++toastCounter}`;

  // 创建toast元素
  const toastElement = createToastElement(toastId, message, type, closable);

  // 获取或创建toast容器
  const container = getToastContainer();

  // 添加到容器
  container.appendChild(toastElement);

  // 记录活动toast
  activeToasts.set(toastId, toastElement);

  // 触发入场动画
  window.requestAnimationFrame(() => {
    toastElement.classList.add('show');
  });

  // 自动关闭
  const autoCloseTimer = setTimeout(() => {
    closeToast(toastId);
  }, duration);

  // 如果可关闭，添加点击事件
  if (closable) {
    const closeBtn = toastElement.querySelector('.toast-close');
    if (closeBtn) {
      closeBtn.addEventListener('click', () => {
        clearTimeout(autoCloseTimer);
        closeToast(toastId);
      });
    }
  }

  // 整个toast点击也可关闭
  toastElement.addEventListener('click', (e) => {
    if (e.target !== toastElement.querySelector('.toast-close')) {
      clearTimeout(autoCloseTimer);
      closeToast(toastId);
    }
  });

  return toastId;
}

/**
 * 创建toast元素
 * @param {string} id - toast ID
 * @param {string} message - 消息内容
 * @param {string} type - 类型
 * @param {boolean} closable - 是否可关闭
 * @returns {HTMLElement} toast元素
 */
function createToastElement(id, message, type, closable) {
  const toast = document.createElement('div');
  toast.id = id;
  toast.className = `toast toast-${type}`;

  // 获取图标
  const icon = getTypeIcon(type);

  // 构建HTML结构
  toast.innerHTML = `
        <div class="toast-content">
            <div class="toast-icon">${icon}</div>
            <div class="toast-message">${escapeHtml(message)}</div>
            ${closable ? '<div class="toast-close">×</div>' : ''}
        </div>
    `;

  return toast;
}

/**
 * 获取或创建toast容器
 * @returns {HTMLElement} 容器元素
 */
function getToastContainer() {
  let container = document.getElementById('toast-container');
  if (!container) {
    container = document.createElement('div');
    container.id = 'toast-container';
    container.className = 'toast-container';
    document.body.appendChild(container);
  }
  return container;
}

/**
 * 关闭toast
 * @param {string} toastId - toast ID
 */
function closeToast(toastId) {
  const toastElement = activeToasts.get(toastId);
  if (!toastElement) return;

  // 添加退出动画
  toastElement.classList.add('hiding');

  // 动画结束后移除元素
  setTimeout(() => {
    if (toastElement.parentNode) {
      toastElement.parentNode.removeChild(toastElement);
    }
    activeToasts.delete(toastId);
  }, 300);
}

/**
 * 获取类型对应的图标
 * @param {string} type - 提示类型
 * @returns {string} 图标HTML
 */
function getTypeIcon(type) {
  const icons = {
    'success': '✓',
    'error': '✕',
    'warning': '⚠',
    'info': 'ℹ'
  };
  return icons[type] || icons.info;
}

/**
 * 转义HTML特殊字符
 * @param {string} text - 要转义的文本
 * @returns {string} 转义后的文本
 */
function escapeHtml(text) {
  const div = document.createElement('div');
  div.textContent = text;
  return div.innerHTML;
}

/**
 * 关闭所有toast
 */
export function closeAllToasts() {
  for (const [toastId] of activeToasts) {
    closeToast(toastId);
  }
}

/**
 * 便捷方法：显示成功提示
 * @param {string} message - 提示消息
 * @param {number} duration - 显示时长
 */
export function showSuccess(message, duration = 3000) {
  return showToast(message, 'success', duration);
}

/**
 * 便捷方法：显示错误提示
 * @param {string} message - 提示消息
 * @param {number} duration - 显示时长
 */
export function showError(message, duration = 4000) {
  return showToast(message, 'error', duration);
}

/**
 * 便捷方法：显示警告提示
 * @param {string} message - 提示消息
 * @param {number} duration - 显示时长
 */
export function showWarning(message, duration = 3500) {
  return showToast(message, 'warning', duration);
}

/**
 * 便捷方法：显示信息提示
 * @param {string} message - 提示消息
 * @param {number} duration - 显示时长
 */
export function showInfo(message, duration = 3000) {
  return showToast(message, 'info', duration);
}
