package controllers

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
)

type FileController struct{}

// Upload 处理文件上传
func (fc *FileController) Upload(c *gin.Context) {
	// 1. 防崩溃保护
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("上传服务发生异常:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "服务内部错误"})
		}
	}()

	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请选择文件"})
		return
	}

	// 2. ★★★ 动态获取当前运行目录 ★★★
	// 无论你在 C盘、D盘 还是 E盘，它都能找到当前程序所在的位置
	workDir, _ := os.Getwd()
	uploadDir := filepath.Join(workDir, "uploads")

	// 3. 自动创建文件夹
	if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
		os.MkdirAll(uploadDir, 0755)
	}

	// 4. 生成文件名
	ext := filepath.Ext(file.Filename)
	// 使用纳秒时间戳，确保唯一性
	filename := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)
	dst := filepath.Join(uploadDir, filename)

	// 5. 保存文件
	if err := c.SaveUploadedFile(file, dst); err != nil {
		fmt.Println("文件保存失败:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "保存失败"})
		return
	}

	fmt.Println(">>> 图片已保存至:", dst)

	// 6. ★★★ 动态生成访问 URL ★★★
	// c.Request.Host 会自动获取当前访问的 IP:端口 (如 localhost:8081 或 192.168.1.5:8081)
	// 这样换了电脑、换了 IP 也能正常显示
	protocol := "http://"
	if c.Request.TLS != nil {
		protocol = "https://"
	}
	fullURL := fmt.Sprintf("%s%s/uploads/%s", protocol, c.Request.Host, filename)

	c.JSON(http.StatusOK, gin.H{
		"url": fullURL,
	})
}
