package service

import (
	"fmt"
	"io/fs"
	"mime/multipart"
	"os"
	"path/filepath"
	"sort"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/yuzujr/C3/internal/config"
	"github.com/yuzujr/C3/internal/eventbus"
	"github.com/yuzujr/C3/internal/logger"
	"github.com/yuzujr/C3/internal/models"
	"github.com/yuzujr/C3/internal/repository"
)

type shot struct {
	url   string
	mtime time.Time
}

// 保存客户端截图文件并记录日志
func SaveScreenshot(c *gin.Context, file *multipart.FileHeader, clientID string) error {
	// 生成带随机数的文件名，防止同一秒文件被覆盖
	ext := filepath.Ext(file.Filename)
	name := file.Filename[:len(file.Filename)-len(ext)]
	randomSuffix := fmt.Sprintf("_%d", time.Now().UnixNano()%1e6)
	newFilename := fmt.Sprintf("%s%s%s", name, randomSuffix, ext)

	dst := filepath.Join(config.Get().Upload.Directory, clientID, time.Now().Format("2006-01-02_15"), newFilename)
	if err := c.SaveUploadedFile(file, dst); err != nil {
		logger.Errorf("Failed to save file: %v", err)
		return err
	}

	// 记录截图日志
	// 如果记录失败，则删除已保存的文件
	client, err := repository.FindClientByID(clientID)
	if err != nil {
		logger.Errorf("Failed to find client: %v", err)
		_ = os.Remove(dst)
		return err
	}

	screenshot := models.ScreenshotLog{
		ClientID: client.ClientID,
		Filename: newFilename,
		FilePath: dst,
		FileSize: int(file.Size),
	}

	if err := repository.CreateScreenshotLog(&screenshot); err != nil {
		logger.Errorf("Failed to log screenshot: %v", err)
		_ = os.Remove(dst)
		return err
	}

	// 广播新截图消息
	msg := eventbus.NewScreenshotMsg{
		Type:     "new_screenshot",
		ClientID: client.ClientID,
		URL:      buildScreenshotURL(dst),
	}
	eventbus.Global.Broadcast(msg)
	return nil
}

// 返回截图url列表
func GetScreenshotsSince(clientID string, sinceMs int64) ([]string, error) {
	cfg := config.Get()
	baseDir := filepath.Join(cfg.Upload.Directory, clientID)

	var shots []shot
	err := filepath.WalkDir(baseDir, func(pathStr string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return nil
		}
		info, err := d.Info()
		if err != nil || info.ModTime().UnixMilli() <= sinceMs {
			return nil
		}

		// 构造URL
		url := buildScreenshotURL(pathStr)

		shots = append(shots, shot{url: url, mtime: info.ModTime()})
		return nil
	})
	if err != nil {
		return nil, err
	}

	// 排序并返回 URL 列表
	sort.Slice(shots, func(i, j int) bool {
		return shots[i].mtime.After(shots[j].mtime)
	})
	urls := make([]string, len(shots))
	for i, s := range shots {
		urls[i] = s.url
	}
	return urls, nil
}

// DeleteScreenshotsAfterHours 删除 clientID 目录下，
// 在 hours 小时之前的所有截图文件，返回删除总数
func DeleteScreenshotsAfterHours(clientID string, hours int) (int, error) {
	cfg := config.Get()
	dir := filepath.Join(cfg.Upload.Directory, clientID)
	fi, err := os.Stat(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return 0, fmt.Errorf("client not found: %s", clientID)
		}
		return 0, err
	}
	if !fi.IsDir() {
		return 0, fmt.Errorf("invalid client directory: %s", dir)
	}

	cutoff := time.Now().Add(-time.Duration(hours) * time.Hour)
	deleted := 0

	entries, _ := os.ReadDir(dir)
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		subdir := filepath.Join(dir, e.Name())
		files, _ := os.ReadDir(subdir)
		for _, f := range files {
			full := filepath.Join(subdir, f.Name())
			st, err := os.Stat(full)
			if err != nil {
				continue
			}
			if st.ModTime().Before(cutoff) {
				if err := os.Remove(full); err == nil {
					deleted++
				}
			}
		}
		// 若该子目录空了，则删掉子目录
		remaining, _ := os.ReadDir(subdir)
		if len(remaining) == 0 {
			os.Remove(subdir)
		}
	}

	return deleted, nil
}

// DeleteAllScreenshots 删除 clientID 目录下的所有截图文件和日期子目录，返回删除总数
func DeleteAllScreenshots(clientID string) (int, error) {
	cfg := config.Get()
	dir := filepath.Join(cfg.Upload.Directory, clientID)
	fi, err := os.Stat(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return 0, fmt.Errorf("client not found: %s", clientID)
		}
		return 0, err
	}
	if !fi.IsDir() {
		return 0, fmt.Errorf("invalid client directory: %s", dir)
	}

	deleted := 0
	entries, _ := os.ReadDir(dir)
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		subdir := filepath.Join(dir, e.Name())
		files, _ := os.ReadDir(subdir)
		for _, f := range files {
			full := filepath.Join(subdir, f.Name())
			if err := os.Remove(full); err == nil {
				deleted++
			}
		}
		os.Remove(subdir)
	}
	return deleted, nil
}

// 传入实际存储路径，返回前端可访问的 URL
func buildScreenshotURL(realPath string) string {
	cfg := config.Get()
	relPath, err := filepath.Rel(cfg.Upload.Directory, realPath)
	if err != nil {
		return ""
	}
	return filepath.Join(
		cfg.Server.BasePath,       // 反向代理的基础路径
		"uploads",                 // 实际的存储位置挂载在"uploads" 路径下
		filepath.ToSlash(relPath), // e.g. "687657f2-86cf-9999/2025-07-19_15/xxx.jpg
	)
}
