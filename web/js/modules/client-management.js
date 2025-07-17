// 客户端管理模块
// 处理客户端别名编辑和删除功能

import { showSuccess, showError, showWarning } from '../../toast/toast.js';
import { buildUrl } from './path-utils.js';

let currentEditingClientId = null;
let currentEditingAlias = null;

/**
 * 初始化客户端管理功能
 */
export function initClientManagement() {
  // 绑定别名编辑相关事件
  document.getElementById('cancelAliasEdit').onclick = closeAliasModal;
  document.getElementById('saveAliasEdit').onclick = saveAliasEdit;

  // 绑定删除客户端相关事件
  document.getElementById('cancelClientDelete').onclick = closeDeleteClientModal;
  document.getElementById('confirmClientDelete').onclick = confirmDeleteClient;

  // 绑定输入框回车事件
  document.getElementById('newAliasInput').onkeypress = (e) => {
    if (e.key === 'Enter') {
      saveAliasEdit();
    }
  };

  // 点击模态框外部关闭
  document.getElementById('aliasModal').onclick = (e) => {
    if (e.target.id === 'aliasModal') {
      closeAliasModal();
    }
  };

  document.getElementById('deleteClientModal').onclick = (e) => {
    if (e.target.id === 'deleteClientModal') {
      closeDeleteClientModal();
    }
  };
}

/**
 * 更新客户端管理列表
 * @param {Array} clients - 客户端列表
 */
export function updateClientManagementList(clients) {
  const container = document.getElementById('clientManagementList');

  if (!container) return;

  if (clients.length === 0) {
    container.innerHTML = '<div style="color:#888; text-align:center; padding:20px;">暂无客户端</div>';
    return;
  } const html = clients.map(client => `
        <div class="client-manage-item" data-client-id="${client.client_id || ''}" data-client-alias="${client.alias}">
            <div class="client-info">
                <span class="client-alias">${client.alias}</span>
            </div>
            <div class="client-actions">
                <button class="edit-alias-btn" onclick="window.clientManagement.editAlias('${client.client_id || ''}', '${client.alias}')">
                    重命名
                </button>
                <button class="delete-client-btn" onclick="window.clientManagement.deleteClient('${client.client_id || ''}', '${client.alias}')">
                    删除
                </button>
            </div>
        </div>
    `).join('');

  container.innerHTML = html;
}

/**
 * 编辑客户端别名
 * @param {string} clientId - 客户端ID  
 * @param {string} alias - 当前别名
 */
export function editAlias(clientId, alias) {
  currentEditingClientId = clientId;
  currentEditingAlias = alias;
  document.getElementById('currentAliasDisplay').textContent = alias;
  document.getElementById('newAliasInput').value = '';
  document.getElementById('aliasModal').style.display = 'block';

  // 聚焦到输入框
  setTimeout(() => {
    document.getElementById('newAliasInput').focus();
  }, 100);
}

/**
 * 删除客户端
 * @param {string} clientId - 客户端ID
 * @param {string} alias - 客户端别名
 */
export function deleteClient(clientId, alias) {
  currentEditingClientId = clientId;
  currentEditingAlias = alias;
  document.getElementById('deleteClientName').textContent = alias;
  document.getElementById('deleteClientModal').style.display = 'block';
}

/**
 * 保存别名编辑
 */
async function saveAliasEdit() {
  const newAlias = document.getElementById('newAliasInput').value.trim();

  if (!newAlias) {
    showError('请输入新的别名');
    return;
  }

  if (newAlias === currentEditingAlias) {
    showWarning('新别名与当前别名相同');
    return;
  }

  // 验证别名格式
  if (!/^[a-zA-Z0-9_-]+$/.test(newAlias)) {
    showError('别名只能包含字母、数字、下划线和连字符');
    return;
  }

  try {
    const response = await fetch(buildUrl(`/web/clients/${encodeURIComponent(currentEditingClientId)}/alias`), {
      method: 'PUT',
      headers: {
        'Content-Type': 'application/json'
      },
      body: JSON.stringify({
        newAlias: newAlias
      })
    });

    const result = await response.json(); if (result.success) {
      showSuccess('别名更新成功');
      closeAliasModal();

      // 立即更新状态
      const { selectedClient, setSelectedClient, setCachedClientList } = await import('./state.js');

      // 清空缓存，强制重新获取
      setCachedClientList([]);

      // 如果当前选中的是被重命名的客户端，立即更新选中状态
      if (selectedClient === currentEditingAlias) {
        setSelectedClient(result.newAlias);
      }

      // 强制刷新客户端列表
      const { fetchClients } = await import('./clients.js');
      await fetchClients();

      // 再次检查并更新选中客户端的状态
      if (selectedClient === result.newAlias) {
        // 更新高亮显示
        const { updateClientHighlight } = await import('./clients.js');
        updateClientHighlight();

        // 找到新客户端的在线状态并更新功能状态
        const { cachedClientList } = await import('./state.js');
        const newClientData = cachedClientList.find(c => c.alias === result.newAlias);
        const isOnline = newClientData ? newClientData.online : true;

        // 更新功能状态
        const { updateClientFeatures } = await import('./clients.js');
        updateClientFeatures(isOnline);
      }
    } else {
      showError(result.message || '更新别名失败');
    }
  } catch (error) {
    console.error('更新别名失败:', error);
    showError('更新别名时发生错误');
  }
}

/**
 * 确认删除客户端
 */
async function confirmDeleteClient() {
  try {
    const response = await fetch(buildUrl(`/web/clients/${encodeURIComponent(currentEditingClientId)}`), {
      method: 'DELETE'
    });

    const result = await response.json();

    if (result.success) {
      showSuccess('客户端删除成功');
      closeDeleteClientModal();

      // 刷新客户端列表
      const { fetchClients } = await import('./clients.js');
      await fetchClients();

      // 如果当前选中的是被删除的客户端，清除选中状态
      const { selectedClient, setSelectedClient } = await import('./state.js');
      if (selectedClient === currentEditingAlias) {
        setSelectedClient(null);
        // 清空截图显示
        document.getElementById('screenshots').textContent = '请选择客户端';
        document.getElementById('commands').style.display = 'none';
      }
    } else {
      showError(result.message || '删除客户端失败');
    }
  } catch (error) {
    console.error('删除客户端失败:', error);
    showError('删除客户端时发生错误');
  }
}

/**
 * 关闭别名编辑模态框
 */
function closeAliasModal() {
  document.getElementById('aliasModal').style.display = 'none';
  currentEditingAlias = null;
}

/**
 * 关闭删除客户端模态框
 */
function closeDeleteClientModal() {
  document.getElementById('deleteClientModal').style.display = 'none';
  currentEditingAlias = null;
}

// 将函数暴露到全局，供HTML onclick使用
window.clientManagement = {
  editAlias,
  deleteClient
};
