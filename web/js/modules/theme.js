// 主题管理模块
// 处理暗色/亮色主题的切换

/**
 * 初始化主题切换功能
 */
export function initThemeToggle() {
    const themeToggle = document.getElementById('themeToggle');
    const html = document.documentElement;

    // 从localStorage恢复主题设置，默认为暗色模式
    const savedTheme = localStorage.getItem('theme') || 'dark';
    setTheme(savedTheme);

    // 点击切换主题
    themeToggle.onclick = () => {
        const currentTheme = html.classList.contains('dark-theme') ? 'dark' : 'light';
        const newTheme = currentTheme === 'dark' ? 'light' : 'dark';
        setTheme(newTheme);
        localStorage.setItem('theme', newTheme);
    };
}

/**
 * 设置主题
 * @param {string} theme - 主题类型 ('dark' 或 'light')
 */
export function setTheme(theme) {
    const html = document.documentElement;
    const themeToggle = document.getElementById('themeToggle');

    if (theme === 'dark') {
        html.classList.add('dark-theme');
        themeToggle.innerHTML = '☀️ 亮色模式';
    } else {
        html.classList.remove('dark-theme');
        themeToggle.innerHTML = '🌙 暗色模式';
    }
}
