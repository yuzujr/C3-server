// 高级PTY终端模块，基于xterm.js，利用后端PTY能力，支持完整的终端体验

import { selectedClient } from './state.js';
import { showWarning, showError } from '../../toast/toast.js';
import { sendCommand } from './websocket.js';

// xterm.js 相关
let terminal = null;
let fitAddon = null;
let webLinksAddon = null;
let terminalContainer = null;
let isXtermLoaded = false;
let welcomeMessageShown = false; // 标志，防止重复显示欢迎信息
let sessionTerminating = false; // 会话正在终止的标志

/**
 * 动态加载 xterm.js 及其插件
 */
async function loadXtermLibrary() {
  if (isXtermLoaded) return true;

  try {
    // 动态创建 script 标签加载 xterm.js
    await loadScript('https://cdn.jsdelivr.net/npm/xterm@5.3.0/lib/xterm.js');
    await loadScript('https://cdn.jsdelivr.net/npm/xterm-addon-fit@0.8.0/lib/xterm-addon-fit.js');
    await loadScript('https://cdn.jsdelivr.net/npm/xterm-addon-web-links@0.9.0/lib/xterm-addon-web-links.js');

    // 动态创建 CSS 链接
    const link = document.createElement('link');
    link.rel = 'stylesheet';
    link.href = 'https://cdn.jsdelivr.net/npm/xterm@5.3.0/css/xterm.css';
    document.head.appendChild(link);

    isXtermLoaded = true;
    return true;
  } catch (error) {
    console.error('Failed to load xterm.js:', error);
    return false;
  }
}

/**
 * 加载外部脚本的辅助函数
 */
function loadScript(src) {
  return new Promise((resolve, reject) => {
    const script = document.createElement('script');
    script.src = src;
    script.onload = resolve;
    script.onerror = reject;
    document.head.appendChild(script);
  });
}

/**
 * 初始化PTY终端
 */
export async function initPtyTerminal() {
  // 加载 xterm.js 库
  const loaded = await loadXtermLibrary();
  if (!loaded) {
    console.error('Failed to load xterm.js');
    return;
  }

  // 等待 xterm.js 全局对象可用
  if (typeof Terminal === 'undefined') {
    console.error('xterm.js not loaded properly');
    return;
  }

  await initXtermTerminal();
}

/**
 * 初始化 xterm.js 终端
 */
async function initXtermTerminal() {
  // 创建终端实例
  terminal = new Terminal({
    cursorBlink: true,
    cursorStyle: 'block',
    fontSize: 14,
    fontFamily: 'Consolas, "Liberation Mono", Menlo, Courier, monospace',
    theme: {
      background: '#1e1e1e',
      foreground: '#d4d4d4',
      cursor: '#ffffff',
      selection: '#264f78',
      black: '#000000',
      red: '#cd3131',
      green: '#0dbc79',
      yellow: '#e5e510',
      blue: '#2472c8',
      magenta: '#bc3fbc',
      cyan: '#11a8cd',
      white: '#e5e5e5',
      brightBlack: '#666666',
      brightRed: '#f14c4c',
      brightGreen: '#23d18b',
      brightYellow: '#f5f543',
      brightBlue: '#3b8eea',
      brightMagenta: '#d670d6',
      brightCyan: '#29b8db',
      brightWhite: '#e5e5e5'
    },
    allowTransparency: false,
    convertEol: true,
    scrollback: 10000,
    tabStopWidth: 4
  });

  // 创建插件实例
  fitAddon = new FitAddon.FitAddon();
  webLinksAddon = new WebLinksAddon.WebLinksAddon();

  // 加载插件
  terminal.loadAddon(fitAddon);
  terminal.loadAddon(webLinksAddon);

  // 获取或创建终端容器
  terminalContainer = document.getElementById('terminalOutput');
  if (!terminalContainer) {
    console.error('Terminal container not found');
    return;
  }

  // 清空容器并打开终端
  terminalContainer.innerHTML = '';
  terminal.open(terminalContainer);

  // 适配大小
  fitAddon.fit();

  // 监听窗口大小变化
  window.addEventListener('resize', () => {
    if (fitAddon && terminal) {
      fitAddon.fit();
      // 通知后端PTY调整大小
      resizePtySession();
    }
  });    // 处理用户输入 - 直接发送给客户端处理
  terminal.onData(data => {
    if (!selectedClient) return;

    // 直接将用户输入发送给客户端，让客户端处理所有逻辑
    // 包括回显、命令执行、提示符显示等
    sendPtyInput(data);
  });    // 绑定传统按钮（兼容性）
  bindTraditionalButtons();

  // 不在初始化时显示欢迎信息，等选择客户端后再显示
}

/**
 * 发送用户输入到PTY会话（逐字符发送）
 */
function sendPtyInput(input) {
  if (!selectedClient) {
    showWarning('请先选择一个客户端');
    return;
  }

  sendCommand({
    type: 'pty_input',
    data: {
      input: input,
      session_id: selectedClient
    }
  });
}

/**
 * 发送新建终端命令（创建新的PTY会话）
 */
function sendNewTerminal() {
  if (!selectedClient) {
    showWarning('请先选择一个客户端');
    return;
  }

  // 获取当前终端大小
  const cols = terminal ? terminal.cols : 80;
  const rows = terminal ? terminal.rows : 24;

  // 发送创建PTY会话命令
  sendCommand({
    type: 'pty_create_session',
    data: {
      session_id: selectedClient,
      cols: cols,
      rows: rows
    }
  });
}

/**
 * 绑定传统按钮
 */
function bindTraditionalButtons() {
  // 绑定新建终端按钮
  document.getElementById('cmdNewTerminal')?.addEventListener('click', () => sendNewTerminal());
  // 绑定终止会话按钮
  document.getElementById('cmdKillSession')?.addEventListener('click', killSession);// 绑定自定义命令输入框（如果存在）
  const executeBtn = document.getElementById('executeCustom');
  const customInput = document.getElementById('customCommand');

  if (executeBtn && customInput) {
    executeBtn.addEventListener('click', () => {
      if (customInput.value.trim()) {
        sendCustomCommand(customInput.value.trim());
        customInput.value = '';
      }
    });

    customInput.addEventListener('keypress', (e) => {
      if (e.key === 'Enter' && customInput.value.trim()) {
        sendCustomCommand(customInput.value.trim());
        customInput.value = '';
      }
    });
  }
}

/**
 * 发送自定义命令
 */
function sendCustomCommand(command) {
  if (!selectedClient) {
    showWarning('请先选择一个客户端');
    return;
  }

  sendPtyInput(command + '\r');
}

/**
 * 显示欢迎信息
 */
async function showWelcomeMessage() {
  if (!terminal) return;

  // 防止重复显示欢迎信息
  if (welcomeMessageShown) {
    console.debug('Welcome message already shown, skipping');
    return;
  }

  // 获取客户端别名用于显示
  const { getClientAlias } = await import('./state.js');
  const clientAlias = selectedClient ? getClientAlias(selectedClient) : 'No client selected';

  // 计算动态填充空格，确保边框对齐
  const boxWidth = 41; // 总宽度
  const connectedToText = 'Connected to: ';
  const usedWidth = connectedToText.length + clientAlias.length + 2; // 2 for │ characters
  const paddingSpaces = Math.max(0, boxWidth - usedWidth);
  const padding = ' '.repeat(paddingSpaces);

  const welcomeMsg = `\r\n\x1b[32m┌─ PTY Terminal ─────────────────────────┐\x1b[0m\r\n`;
  const welcomeMsg2 = `\x1b[32m│\x1b[0m ${connectedToText}\x1b[36m${clientAlias}\x1b[0m${padding}\x1b[32m│\x1b[0m\r\n`;
  const welcomeMsg3 = `\x1b[32m│\x1b[0m Full PTY support with xterm.js         \x1b[32m│\x1b[0m\r\n`;
  const welcomeMsg4 = `\x1b[32m└────────────────────────────────────────┘\x1b[0m\r\n\r\n`;

  terminal.write(welcomeMsg + welcomeMsg2 + welcomeMsg3 + welcomeMsg4);
  welcomeMessageShown = true; // 设置标志，防止重复显示
}

/**
 * 强制显示欢迎信息（忽略防重复标志）
 */
async function forceShowWelcomeMessage() {
  if (!terminal) return;

  welcomeMessageShown = false; // 重置标志
  await showWelcomeMessage();
}

/**
 * 调整PTY会话大小
 */
function resizePtySession() {
  if (!selectedClient || !terminal || !fitAddon) return;

  const cols = terminal.cols;
  const rows = terminal.rows;

  sendCommand({
    type: 'pty_resize',
    data: {
      session_id: selectedClient,
      cols,
      rows
    }
  });
}

/**
 * 处理从WebSocket接收到的PTY输出（完整保留原始数据）
 */
export function handlePtyShellOutput(data) {
  if (!data) return;

  // 如果会话正在终止，忽略输出
  if (sessionTerminating) {
    return;
  }

  // 使用原始输出
  const output = data.output || '';

  // 使用 xterm.js 直接写入原始数据
  if (terminal) {
    // 保持原始数据完整性，包括所有控制字符和ANSI序列
    // 客户端会返回所有内容，包括命令回显、输出和新的提示符
    terminal.write(output);
  } else {
    console.warn('Terminal not initialized, ignoring output');
  }
}

/**
 * 强制终止PTY会话
 */
async function killSession() {
  if (!selectedClient) {
    showWarning('请先选择一个客户端');
    return;
  }

  if (!confirm('确定要强制终止当前会话吗？这将关闭所有正在运行的程序。')) {
    return;
  }

  try {
    // 设置会话终止标志
    sessionTerminating = true;

    sendCommand({
      type: 'pty_kill_session',
      data: { session_id: selectedClient }
    });

    // 终止会话后清屏并显示连接信息
    if (terminal) {
      terminal.clear();
      forceShowWelcomeMessage();
      // 清除终止状态标志
      sessionTerminating = false;
    }
  } catch (error) {
    console.error('强制终止会话失败:', error);
    showError('强制终止会话失败');
    // 出错时也要清除终止状态标志
    sessionTerminating = false;
  }
}

/**
 * 更新终端状态
 */
export async function updatePtyTerminalState(enabled) {
  const buttons = ['cmdNewTerminal', 'executeCustom', 'cmdKillSession'];
  const customInput = document.getElementById('customCommand');

  buttons.forEach(id => {
    const button = document.getElementById(id);
    if (button) {
      button.disabled = !enabled;
    }
  });

  if (customInput) {
    customInput.disabled = !enabled;
  }

  // 只有在首次启用且有选择客户端时才显示欢迎信息
  // 并且只有在终端为空时才显示（避免重复显示）
  if (terminal && enabled && selectedClient) {
    try {
      // 检查终端是否为空（通过缓冲区行数判断）
      const isEmpty = terminal.buffer.active.length <= 1;

      if (isEmpty) {
        // 清屏后显示欢迎信息
        terminal.clear();
        await showWelcomeMessage();
      }
    } catch (error) {
      console.debug('Terminal buffer check failed, showing welcome message:', error);
      // 如果检查失败，默认显示欢迎信息
      terminal.clear();
      await showWelcomeMessage();
    }
    // 让客户端自己显示提示符
  }
}

/**
 * 重置终端（切换客户端时调用）
 */
export async function resetPtyTerminal() {
  if (terminal) {
    terminal.clear();
    welcomeMessageShown = false; // 重置标志，允许重新显示欢迎信息
    if (selectedClient) {
      await showWelcomeMessage();
      // 让客户端自己显示提示符
    }
  } else {
    console.warn('Terminal not initialized');
  }
}
