package controllers

import (
	"gotest/config"
	"gotest/internal/models"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ProductController struct {
	// 移除 Service 依赖，直接使用 config.DB
}

// List 获取商品列表
func (p *ProductController) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	category := c.Query("category")
	search := c.Query("search")
	isRandom := c.Query("is_random") == "true"

	var products []models.Product
	var total int64

	db := config.DB.Model(&models.Product{})

	// 1. 过滤状态：只显示在售商品
	db = db.Where("status = ?", 1)

	// 2. 搜索逻辑
	if search != "" {
		db = db.Where("name LIKE ? OR description LIKE ?", "%"+search+"%", "%"+search+"%")
	}

	// 3. 分类筛选
	if category != "" && category != "全部" {
		db = db.Where("category = ?", category)
	}

	// 4. 核心逻辑
	if isRandom {
		// ★★★ 修复点：随机推荐也要 Preload("User")，否则首页显示不出卖家头像 ★★★
		// SQLite 使用 RANDOM()，MySQL 使用 RAND()
		if err := db.Order("RANDOM()").Limit(pageSize).Preload("User").Find(&products).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "获取推荐失败"})
			return
		}
		total = int64(len(products))
	} else {
		// 普通分页列表
		db.Count(&total)
		offset := (page - 1) * pageSize
		// 这里原代码已经加了 Preload("User")，保持不变
		if err := db.Order("created_at desc").Offset(offset).Limit(pageSize).Preload("User").Find(&products).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "获取列表失败"})
			return
		}
	}

	// 图片路径处理 (防止前端图片裂开)
	for i := range products {
		if products[i].Image == "" {
			products[i].Image = "/uploads/default_product.png"
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"list":  products,
		"total": total,
		"page":  page,
	})
}

// GetDetail 获取商品详情
func (p *ProductController) GetDetail(c *gin.Context) {
	id := c.Param("id")
	var product models.Product

	// 增加浏览量 (使用 gorm.Expr 原子更新)
	config.DB.Model(&models.Product{}).Where("id = ?", id).UpdateColumn("view_count", gorm.Expr("view_count + ?", 1))

	// ★★★ 核心修复：必须 Preload("User") 才能获取卖家头像和昵称 ★★★
	if err := config.DB.Preload("User").First(&product, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "商品不存在"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": product})
}

// Categories 获取分类列表
func (p *ProductController) Categories(c *gin.Context) {
	categories := []string{"数码", "书籍", "生活", "服饰", "运动", "美妆", "乐器", "手游", "兼职", "其他"}
	c.JSON(http.StatusOK, gin.H{"data": categories})
}

// Create 发布商品
func (p *ProductController) Create(c *gin.Context) {
	var input models.Product
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未登录"})
		return
	}

	// 绑定当前登录用户 ID
	input.UserID = userID.(uint)
	input.Status = 1
	input.CreatedAt = time.Now()
	input.UpdatedAt = time.Now()

	if err := config.DB.Create(&input).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "发布失败: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "发布成功", "data": input})
}

// Update 更新商品
func (p *ProductController) Update(c *gin.Context) {
	id := c.Param("id")
	userID, _ := c.Get("userID")

	var product models.Product
	if err := config.DB.First(&product, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "商品不存在"})
		return
	}

	// 权限校验
	if product.UserID != userID.(uint) {
		c.JSON(http.StatusForbidden, gin.H{"error": "无权修改此商品"})
		return
	}

	var input models.Product
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := config.DB.Model(&product).Updates(input).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "更新成功", "data": product})
}
