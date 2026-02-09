package handler

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/yuzujr/C3/internal/config"
)

// ServeHTMLWithConfig 提供HTML文件并注入配置
func ServeHTMLWithConfig(htmlFile string) gin.HandlerFunc {
	return func(c *gin.Context) {
		cwd, err := os.Getwd()
		if err != nil {
			c.String(http.StatusInternalServerError, "Internal server error")
			return
		}

		htmlPath := filepath.Join(cwd, "web", htmlFile)
		content, err := os.ReadFile(htmlPath)
		if err != nil {
			c.String(http.StatusNotFound, "File not found")
			return
		}

		// 获取配置
		cfg := config.Get()
		basePath := cfg.Server.BasePath
		
		// 构建base标签（确保以/结尾）
		baseHref := basePath
		if baseHref != "" && !strings.HasSuffix(baseHref, "/") {
			baseHref += "/"
		}

		// 构建配置注入内容
		var injections strings.Builder
		
		// 1. base标签 - 让所有相对路径自动加上BASE_PATH前缀
		if baseHref != "" {
			injections.WriteString(`<base href="` + baseHref + `">`)
			injections.WriteString("\n  ")
		}
		
		// 2. 配置脚本 - 提供给JavaScript使用
		injections.WriteString(`<script>
    window.APP_CONFIG = {
      BASE_PATH: "` + basePath + `"
    };
  </script>`)

		// 在<head>后注入（放在最前面）
		htmlContent := string(content)
		htmlContent = strings.Replace(htmlContent, "<head>", "<head>\n  "+injections.String(), 1)

		c.Header("Content-Type", "text/html; charset=utf-8")
		c.String(http.StatusOK, htmlContent)
	}
}
