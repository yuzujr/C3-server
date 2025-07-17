// 标签页管理模块
// 处理标签页的切换和状态管理

/**
 * 初始化标签页功能
 */
export function initTabs() {
  const tabButtons = document.querySelectorAll('.tab-button');

  tabButtons.forEach(button => {
    button.addEventListener('click', () => {
      // 检查按钮是否被禁用
      if (button.disabled) {
        return;
      }

      const targetTab = button.dataset.tab;
      switchTab(targetTab);
    });
  });
}

/**
 * 切换到指定标签页
 * @param {string} tabName - 标签页名称
 */
export function switchTab(tabName) {
  // 移除所有激活状态
  document.querySelectorAll('.tab-button').forEach(btn => {
    btn.classList.remove('active');
  });

  document.querySelectorAll('.tab-panel').forEach(panel => {
    panel.classList.remove('active');
  });

  // 激活指定标签页
  const targetButton = document.querySelector(`[data-tab="${tabName}"]`);
  const targetPanel = document.getElementById(`tab-${tabName}`);

  if (targetButton && targetPanel) {
    targetButton.classList.add('active');
    targetPanel.classList.add('active');
  }
}

/**
 * 禁用/启用指定标签页
 * @param {string} tabName - 标签页名称
 * @param {boolean} disabled - 是否禁用
 */
export function setTabDisabled(tabName, disabled) {
  const tabButton = document.querySelector(`[data-tab="${tabName}"]`);
  if (tabButton) {
    if (disabled) {
      tabButton.disabled = true;
      tabButton.style.opacity = '0.5';
      tabButton.style.cursor = 'not-allowed';
    } else {
      tabButton.disabled = false;
      tabButton.style.opacity = '1';
      tabButton.style.cursor = 'pointer';
    }
  }
}

/**
 * 获取当前激活的标签页
 * @returns {string} 当前激活的标签页名称
 */
export function getCurrentTab() {
  const activeButton = document.querySelector('.tab-button.active');
  return activeButton ? activeButton.dataset.tab : 'overview';
}