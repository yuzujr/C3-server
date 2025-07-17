// 图片模态框模块
// 处理图片的弹窗显示、导航和关闭

import { imageUrls, currentImageIndex, setCurrentImageIndex } from './state.js';

/**
 * 解析图片名中的日期
 * @param {string} imageUrl - 图片URL
 * @returns {string} 格式化的日期字符串或文件名
 */
function parseImageDate(imageUrl) {
  const filename = imageUrl.split('/').pop();

  // 尝试匹配常见的日期时间格式
  // 格式如: screenshot_20240610_143025.jpg, 20240610-143025.png 等
  const patterns = [
    /(\d{4})(\d{2})(\d{2})[_-](\d{2})(\d{2})(\d{2})/,  // YYYYMMDD_HHMMSS 或 YYYYMMDD-HHMMSS
    /(\d{4})[_-](\d{2})[_-](\d{2})[_-](\d{2})[_-](\d{2})[_-](\d{2})/, // YYYY_MM_DD_HH_MM_SS
    /(\d{4})\.(\d{2})\.(\d{2})\.(\d{2})\.(\d{2})\.(\d{2})/, // YYYY.MM.DD.HH.MM.SS
  ];

  for (const pattern of patterns) {
    const match = filename.match(pattern);
    if (match) {
      const [, year, month, day, hour, minute, second] = match;
      const date = new Date(year, month - 1, day, hour, minute, second);

      // 格式化显示
      const options = {
        year: 'numeric',
        month: '2-digit',
        day: '2-digit',
        hour: '2-digit',
        minute: '2-digit',
        second: '2-digit',
        hour12: false
      };

      return date.toLocaleString('zh-CN', options);
    }
  }

  // 如果没有匹配到日期格式，返回文件名
  return filename;
}

/**
 * 打开图片模态框
 * @param {string} imageUrl - 图片URL
 * @param {number} index - 图片在列表中的索引
 */
export function openImageModal(imageUrl, index = 0) {
  const modal = document.getElementById('imageModal');
  const modalImg = document.getElementById('modalImage');
  const caption = document.getElementById('caption');

  setCurrentImageIndex(index);
  modal.style.display = 'block';
  modalImg.src = imageUrl;
  caption.innerHTML = parseImageDate(imageUrl);

  updateNavigationButtons();
}

/**
 * 关闭图片模态框
 */
export function closeImageModal() {
  const modal = document.getElementById('imageModal');
  modal.style.display = 'none';
}

/**
 * 更新当前显示的图片
 * @param {number} newIndex - 新的图片索引
 */
function updateCurrentImage(newIndex) {
  if (imageUrls.length === 0 || newIndex < 0 || newIndex >= imageUrls.length) return;

  setCurrentImageIndex(newIndex);
  const imageUrl = imageUrls[currentImageIndex];

  const modalImg = document.getElementById('modalImage');
  const caption = document.getElementById('caption');

  modalImg.src = imageUrl;
  caption.innerHTML = parseImageDate(imageUrl);

  updateNavigationButtons();
}

/**
 * 显示上一张图片
 */
export function showPreviousImage() {
  if (imageUrls.length === 0 || currentImageIndex <= 0) return;
  updateCurrentImage(currentImageIndex - 1);
}

/**
 * 显示下一张图片
 */
export function showNextImage() {
  if (imageUrls.length === 0 || currentImageIndex >= imageUrls.length - 1) return;
  updateCurrentImage(currentImageIndex + 1);
}

/**
 * 更新导航按钮状态
 */
function updateNavigationButtons() {
  const prevBtn = document.getElementById('prevBtn');
  const nextBtn = document.getElementById('nextBtn');

  // 如果只有一张图片，隐藏导航按钮
  if (imageUrls.length <= 1) {
    prevBtn.style.display = 'none';
    nextBtn.style.display = 'none';
    return;
  }

  // 显示导航按钮
  prevBtn.style.display = 'block';
  nextBtn.style.display = 'block';

  // 设置按钮状态的通用函数
  function setButtonState(button, enabled) {
    if (enabled) {
      button.style.opacity = '1';
      button.style.cursor = 'pointer';
      button.style.pointerEvents = 'auto';
    } else {
      button.style.opacity = '0.3';
      button.style.cursor = 'not-allowed';
      button.style.pointerEvents = 'none';
    }
  }

  // 设置按钮启用/禁用状态
  setButtonState(prevBtn, currentImageIndex > 0);
  setButtonState(nextBtn, currentImageIndex < imageUrls.length - 1);
}

/**
 * 初始化图片模态框事件监听器
 */
export function initImageModal() {
  const modal = document.getElementById('imageModal');
  const closeBtn = document.querySelector('.close');
  const prevBtn = document.getElementById('prevBtn');
  const nextBtn = document.getElementById('nextBtn');

  // 点击关闭按钮关闭模态框
  closeBtn.onclick = closeImageModal;

  // 点击导航箭头
  prevBtn.onclick = showPreviousImage;
  nextBtn.onclick = showNextImage;

  // 点击模态框背景关闭模态框
  modal.onclick = function (event) {
    if (event.target === modal) {
      closeImageModal();
    }
  };

  // 按键导航
  document.addEventListener('keydown', function (event) {
    if (modal.style.display === 'block') {
      switch (event.key) {
        case 'Escape':
          closeImageModal();
          break;
        case 'ArrowLeft':
          showPreviousImage();
          break;
        case 'ArrowRight':
          showNextImage();
          break;
      }
    }
  });
}
