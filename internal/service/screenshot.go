package service

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"time"

	"github.com/yuzujr/C3/internal/config"
)

type shot struct {
	url   string
	mtime time.Time
}

// 返回截图url列表，要求前端可访问
func GetScreenshotsSince(clientID string, sinceMs int64) ([]string, error) {
	cfg := config.Get()
	// 由于已经把实际存储路径挂载到了 uploads 路径，这里可直接使用 uploads
	baseDir := filepath.Join("uploads", clientID)

	var shots []shot
	err := filepath.WalkDir(baseDir, func(pathStr string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return nil
		}
		info, err := d.Info()
		if err != nil || info.ModTime().UnixMilli() <= sinceMs {
			return nil
		}

		// 1. 先算出相对于 UPLOAD_DIR 目录的相对路径
		relPath, err := filepath.Rel("uploads", pathStr)
		if err != nil {
			return nil
		}
		// 2. 拼接成完整的 URL
		url := fmt.Sprintf("%s/%s/%s",
			cfg.Server.BasePath, //反向代理的基础路径
			"uploads",
			filepath.ToSlash(relPath), // e.g. "687657f2-86cf-9999/2025-07-19_15/xxx.jpg"
		)

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
