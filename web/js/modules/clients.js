// 客户端管理模块
// 处理客户端列表的获取、选择和高亮显示
// 统一使用alias作为客户端标识符

import { selectedClient, setSelectedClient, cachedClientList, setCachedClientList } from './state.js';
import { fetchScreenshots } from './screenshots.js';
import { loadClientConfig } from './commands.js';
import { showError } from '../../toast/toast.js';
import { buildUrl } from './path-utils.js';

/**
 * 比较两个客户端列表是否相同
 * @param {Array} oldList - 旧的客户端列表
 * @param {Array} newList - 新的客户端列表
 * @returns {boolean} - 是否相同
 */
function isClientListSame(oldList, newList) {
  if (oldList.length !== newList.length) {
    return false;
  }

  for (let i = 0; i < oldList.length; i++) {
    const oldClient = oldList[i];
    const newClient = newList[i];

    if (oldClient.alias !== newClient.alias || oldClient.online !== newClient.online) {
      return false;
    }
  }

  return true;
}

/**
 * 根据客户端在线状态更新功能可用性
 * @param {boolean} isOnline - 客户端是否在线
 */
export function updateClientFeatures(isOnline) {
  // 导入标签页管理模块
  import('./tabs.js').then(({ setTabDisabled }) => {
    // 设置Shell和配置标签页的可用性
    setTabDisabled('shell', !isOnline);
    setTabDisabled('config', !isOnline);
  });

  if (isOnline) {
    // 在线客户端：启用所有功能

    // 显示命令按钮
    document.getElementById('cmdButtons').style.display = 'block';
    const commandMessage = document.getElementById('commandMessage');
    if (commandMessage) {
      commandMessage.style.display = 'none';
    }        // 启用终端功能 - 使用pty-terminal模块的统一方法
    import('./pty-terminal.js').then(async ({ updatePtyTerminalState }) => {
      await updatePtyTerminalState(true);
    });

    // 显示配置区域
    const configSection = document.getElementById('configSection');
    const configMessage = document.getElementById('configMessage');
    if (configSection && configMessage) {
      configSection.style.display = 'block';
      configMessage.style.display = 'none';
    }

  } else {
    // 离线客户端：只允许查看截图

    // 隐藏命令按钮
    document.getElementById('cmdButtons').style.display = 'none';
    const commandMessage = document.getElementById('commandMessage');
    if (commandMessage) {
      commandMessage.textContent = '客户端离线，无法发送命令';
      commandMessage.style.display = 'block';
    }

    // 禁用终端功能 - 使用pty-terminal模块的统一方法
    import('./pty-terminal.js').then(async ({ updatePtyTerminalState }) => {
      await updatePtyTerminalState(false);
    });

    // 隐藏配置区域
    const configSection = document.getElementById('configSection');
    const configMessage = document.getElementById('configMessage');
    if (configSection && configMessage) {
      configSection.style.display = 'none';
      configMessage.textContent = '客户端离线，无法修改配置';
      configMessage.style.display = 'block';
    }
  }
}

/**
 * 获取客户端列表
 */
export async function fetchClients() {
  try {
    const res = await fetch(buildUrl('/web/clients'));

    // 检查HTTP响应状态
    if (!res.ok) {
      throw new Error(`HTTP ${res.status}: ${res.statusText}`);
    }

    const clients = await res.json();

    // 检查是否应该跳过DOM重建（例如在WebSocket别名更新后）
    const { skipNextDOMRebuild, setSkipNextDOMRebuild } = await import('./state.js');
    if (skipNextDOMRebuild) {
      setCachedClientList([...clients]);
      setSkipNextDOMRebuild(false); // 重置标志
      return;
    }

    // 检查客户端列表是否有变化，但第一次加载时要强制更新DOM
    const isFirstLoad = cachedClientList.length === 0 && document.getElementById('clientsList').textContent === '加载中...';
    if (!isFirstLoad && isClientListSame(cachedClientList, clients)) {
      // 列表没有变化，直接返回，避免不必要的DOM更新
      return;
    }

    // 更新缓存
    setCachedClientList([...clients]);

    // 更新客户端管理列表
    import('./client-management.js').then(({ updateClientManagementList }) => {
      updateClientManagementList(clients);
    });

    const clientsList = document.getElementById('clientsList');
    clientsList.innerHTML = '';

    if (clients.length === 0) {
      clientsList.textContent = '无客户端';
      setSelectedClient(null);
      document.getElementById('commands').style.display = 'none';
      document.getElementById('screenshots').textContent = '请选择客户端';

      // 更新客户端管理列表（空列表）
      import('./client-management.js').then(({ updateClientManagementList }) => {
        updateClientManagementList([]);
      });

      // 禁用终端功能 - 使用pty-terminal模块的统一方法
      import('./pty-terminal.js').then(async ({ updatePtyTerminalState }) => {
        await updatePtyTerminalState(false);
      });
      return;
    }

    // 检查是否有在线客户端
    const onlineClients = clients.filter(c => c.online);
    if (onlineClients.length === 0) {
      // 所有客户端都离线，但仍然显示客户端列表（可以查看截图）
      // 如果当前有选中的客户端，更新其功能状态
      if (selectedClient) {
        const selectedClientData = clients.find(c => c.client_id === selectedClient); // 使用 client_id 查找
        if (selectedClientData) {
          updateClientFeatures(false); // 离线状态
        }
      }
    }

    // 显示命令区域
    document.getElementById('commands').style.display = 'block';

    // 检查当前选中的客户端是否仍然在线
    if (selectedClient) {
      const selectedClientData = clients.find(c => c.client_id === selectedClient); // 使用 client_id 查找
      if (selectedClientData) {
        // 客户端存在，更新功能可用性
        updateClientFeatures(selectedClientData.online);
      } else {
        // 客户端不存在了，清除选择
        setSelectedClient(null);
        document.getElementById('cmdButtons').style.display = 'none';
        document.getElementById('screenshots').textContent = '请选择客户端';

        const commandMessage = document.getElementById('commandMessage');
        if (commandMessage) {
          commandMessage.textContent = '请选择客户端';
        }

        // 禁用终端功能
        import('./pty-terminal.js').then(async ({ updatePtyTerminalState }) => {
          await updatePtyTerminalState(false);
        });
      }
    }

    // 构建客户端列表
    clients.forEach(client => {
      const clientItem = document.createElement('div');
      clientItem.className = `client-item ${client.online ? 'online' : 'offline'}`;
      clientItem.dataset.clientId = client.client_id; // 使用 client_id
      clientItem.dataset.clientAlias = client.alias;  // 保留 alias 用于向后兼容
      clientItem.dataset.online = client.online;

      const statusText = client.online ? '在线' : '离线';
      const statusClass = client.online ? 'online' : 'offline';

      clientItem.innerHTML = `
                <span class="client-name">${client.alias}</span>
                <span class="client-status ${statusClass}">${statusText}</span>
            `;

      clientItem.addEventListener('click', () => {
        selectClient(client.client_id, client.online);
      });

      clientsList.appendChild(clientItem);
    });

    // 恢复选中状态
    if (selectedClient) {
      updateClientHighlight();
    }

  } catch (error) {
    console.error('获取客户端列表失败:', error);
    showError('获取客户端列表失败');

    // 更新DOM显示错误状态
    const clientsList = document.getElementById('clientsList');
    if (clientsList) {
      clientsList.innerHTML = '';
      clientsList.textContent = '获取客户端列表失败，请刷新页面重试';
    }

    // 清除选中状态和隐藏相关功能
    setSelectedClient(null);
    document.getElementById('commands').style.display = 'none';
    document.getElementById('screenshots').textContent = '请选择客户端';

    // 禁用终端功能
    import('./pty-terminal.js').then(async ({ updatePtyTerminalState }) => {
      await updatePtyTerminalState(false);
    });
  }
}

/**
 * 选择客户端
 * @param {string} clientId - 客户端ID（现在使用 client_id 而不是 alias）
 * @param {boolean} isOnline - 客户端是否在线
 */
export async function selectClient(clientId, isOnline = true) {
  // 如果选择的是同一个客户端，且在线状态没有变化，则不需要重新处理
  const currentClientData = cachedClientList.find(c => c.client_id === selectedClient);
  if (selectedClient === clientId && currentClientData && currentClientData.online === isOnline) {
    // 只更新高亮显示即可
    updateClientHighlight();
    return;
  }

  // 记录是否是切换到不同的客户端
  const isDifferentClient = selectedClient !== clientId;

  setSelectedClient(clientId);
  updateClientHighlight();

  document.getElementById('screenshots').textContent = '加载截图中...';

  // 根据在线状态启用/禁用功能
  updateClientFeatures(isOnline);    // 重置终端工作目录（只在切换到不同客户端时）
  if (isDifferentClient) {
    import('./pty-terminal.js').then(async ({ resetPtyTerminal }) => {
      await resetPtyTerminal();
    });
  }

  // 始终加载截图（在线和离线都可以查看）
  try {
    await fetchScreenshots(clientId);
  } catch (error) {
    console.error(`加载截图失败:`, error);
    const container = document.getElementById('screenshots');
    if (container) {
      container.textContent = `加载截图失败: ${error.message}`;
    }
  }

  // 只有在线客户端才加载配置
  if (isOnline) {
    try {
      await loadClientConfig(clientId);
    } catch (error) {
      console.error(`加载配置失败:`, error);
    }
  }
}

/**
 * 更新客户端高亮显示
 */
export function updateClientHighlight() {
  document.querySelectorAll('.client-item').forEach(div => {
    if (div.dataset.clientId === selectedClient) { // 使用 client_id 比较
      div.classList.add('selected');
    } else {
      div.classList.remove('selected');
    }
  });
}

/**
 * 处理WebSocket客户端状态变化通知
 * @param {string} clientId - 客户端ID
 * @param {boolean} isOnline - 是否在线
 */
export function handleClientStatusChange(clientId, isOnline) {
  // 更新缓存中客户端的状态
  const clientIndex = cachedClientList.findIndex(c => c.client_id === clientId);
  if (clientIndex !== -1) {
    cachedClientList[clientIndex].online = isOnline;

    // 只有当状态变化的客户端是当前选中的客户端时，才更新功能可用性
    if (clientId === selectedClient) { // 使用 client_id 比较
      updateClientFeatures(isOnline);

      // 更新标签页状态
      import('./tabs.js').then(({ setTabDisabled }) => {
        setTabDisabled('shell', !isOnline);
        setTabDisabled('config', !isOnline);
      });
    }

    // 更新DOM显示 - 使用 client_id
    updateSingleClientDisplay(clientId, isOnline);
  } else {
    // 如果是新的客户端，直接刷新客户端列表
    fetchClients();
  }
}

/**
 * 更新单个客户端在DOM中的显示状态
 * @param {string} clientId - 客户端ID
 * @param {boolean} isOnline - 是否在线
 */
function updateSingleClientDisplay(clientId, isOnline) {
  const clientElement = document.querySelector(`[data-client-id="${clientId}"]`);
  if (clientElement) {
    const statusText = isOnline ? '在线' : '离线';
    const statusClass = isOnline ? 'online' : 'offline';

    // 更新状态显示，保持原有的HTML结构
    const statusElement = clientElement.querySelector('.client-status');
    if (statusElement) {
      statusElement.textContent = statusText;
      statusElement.className = `client-status ${statusClass}`;
    }

    // 更新CSS类，保持选中状态
    clientElement.className = `client-item ${statusClass}`;
    if (clientId === selectedClient) {
      clientElement.classList.add('selected');
    }

    // 更新数据属性
    clientElement.dataset.online = isOnline;
  }
}
