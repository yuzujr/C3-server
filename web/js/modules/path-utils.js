// 路径工具模块
// 统一处理前端基础路径相关逻辑

/**
 * 获取当前应用的基础路径
 * @returns {string} 基础路径
 */
export function getBasePath() {
  // 优先使用服务器注入的配置
  if (window.APP_CONFIG && window.APP_CONFIG.BASE_PATH) {
    return window.APP_CONFIG.BASE_PATH;
  }

  // 降级方案：从当前URL智能分析
  const pathname = window.location.pathname;
  const pathParts = pathname.split('/').filter(part => part);

  // 如果第一个路径段不是常见的页面名，可能是base path
  if (pathParts.length > 0 && !['login', 'index.html', ''].includes(pathParts[0])) {
    return '/' + pathParts[0];
  }

  return '';
}

/**
 * 构建完整的URL路径
 * @param {string} path - 相对路径
 * @returns {string} 完整路径
 */
export function buildUrl(path) {
  return getBasePath() + path;
}

/**
 * 构建WebSocket连接URL
 * @returns {string} WebSocket URL
 */
export function buildWebSocketUrl() {
  const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
  const host = window.location.host;
  const basePath = getBasePath();

  // 确保basePath以/结尾，以匹配Nginx的location规则
  const normalizedBasePath = basePath.endsWith('/') ? basePath : basePath + '/';

  return `${protocol}//${host}${normalizedBasePath}ws?type=web`;
}
