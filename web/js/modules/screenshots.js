// 截图处理模块
// 处理截图的获取、显示和新截图的添加

import { imageUrls, setImageUrls, addNewImageUrl } from './state.js';
import { openImageModal } from './modal.js';
import { buildUrl } from './path-utils.js';

/**
 * 获取指定客户端的截图列表
 * @param {string} clientId - 客户端ID
 * @returns {Promise<number>} 返回当前时间戳用于下次请求
 */
export async function fetchScreenshots(clientId) {
  try {
    const res = await fetch(buildUrl(`/web/screenshots/${clientId}`));

    if (!res.ok) {
      throw new Error(`HTTP error! status: ${res.status}`);
    }

    const screenshots = await res.json();

    const container = document.getElementById('screenshots');
    container.innerHTML = '';

    if (screenshots.length === 0) {
      container.textContent = '该客户端暂无截图';
      setImageUrls([]);
      return;
    }

    setImageUrls(screenshots);

    screenshots.forEach((url, index) => {
      const img = document.createElement('img');
      img.src = url;
      img.onclick = () => openImageModal(url, index);
      container.appendChild(img);
    });

    return Date.now(); // 返回当前时间戳用于下次请求
  } catch (error) {
    console.error('获取截图列表失败:', error);
    const container = document.getElementById('screenshots');
    if (container) {
      container.textContent = `获取截图失败: ${error.message}`;
    }
    throw error;
  }
}

/**
 * 添加新截图到界面
 * @param {string} screenshotUrl - 新截图的URL
 */
export function addNewScreenshot(screenshotUrl) {
  // 检查是否已存在相同的截图
  if (imageUrls.includes(screenshotUrl)) {
    return;
  }

  // 将新图片添加到URL列表的开头
  addNewImageUrl(screenshotUrl);

  const container = document.getElementById('screenshots');

  // 如果容器中只有文字内容（没有截图），清空容器
  if (container.children.length === 0 && container.textContent.trim() !== '') {
    container.innerHTML = '';
  }

  // 创建新的图片元素
  const img = document.createElement('img');
  img.src = screenshotUrl;
  img.onclick = () => openImageModal(screenshotUrl, 0); // 新图片索引为0

  // 将新图片插入到容器的开头
  container.insertBefore(img, container.firstChild);

  // 更新所有现有图片的点击事件索引
  const images = container.querySelectorAll('img');
  images.forEach((imgElement, index) => {
    imgElement.onclick = () => openImageModal(imageUrls[index], index);
  });
}
