// 状态管理模块
// 全局状态变量的集中管理

export let selectedClient = null; // 存储 client_id，不再是 alias
export let currentImageIndex = 0;
export let imageUrls = [];
export let webSocket = null;
export let cachedClientList = []; // 缓存的客户端列表
export let skipNextDOMRebuild = false; // 临时标志：跳过下次DOM重建

// 常量配置
export const CLIENT_POLL_INTERVAL = 5000; // 5秒 - 客户端列表更新频率

// 状态更新函数
export function setSelectedClient(clientId) {
  selectedClient = clientId; // 现在存储的是 client_id
}

export function setSkipNextDOMRebuild(skip) {
  skipNextDOMRebuild = skip;
}

export function setCurrentImageIndex(index) {
  currentImageIndex = index;
}

export function setImageUrls(urls) {
  imageUrls = [...urls];
}

export function setWebSocket(ws) {
  webSocket = ws;
}

export function addNewImageUrl(url) {
  if (!imageUrls.includes(url)) {
    imageUrls.unshift(url);
  }
}

export function clearImageUrls() {
  imageUrls = [];
}

export function setCachedClientList(clients) {
  cachedClientList = clients;
}

/**
 * 根据 client_id 获取 alias（显示名称）
 * @param {string} clientId - 客户端ID
 * @returns {string} alias 或 client_id（如果找不到）
 */
export function getClientAlias(clientId) {
  if (!clientId) return '';
  const client = cachedClientList.find(c => c.client_id === clientId);
  return client ? client.alias : clientId;
}