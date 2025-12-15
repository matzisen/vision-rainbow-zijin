package controllers

import (
	"gotest/config"
	"gotest/internal/models"
	"net/http"
	"strconv" // ★★★ 新增：用于订单号拼接
	"time"

	"github.com/gin-gonic/gin"
)

type OrderController struct {
	// 直接使用 config.DB
}

// Create 创建订单 (单商品直接购买)
func (o *OrderController) Create(c *gin.Context) {
	var input struct {
		ProductID uint `json:"product_id"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	userID, _ := c.Get("userID")
	uid := userID.(uint)

	// 1. 开启事务
	tx := config.DB.Begin()

	// 2. 查询商品
	var product models.Product
	// 加锁查询防止超卖
	if err := tx.Set("gorm:query_option", "FOR UPDATE").First(&product, input.ProductID).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusNotFound, gin.H{"error": "商品不存在"})
		return
	}

	// 3. 校验状态
	if product.Status != 1 {
		tx.Rollback()
		c.JSON(http.StatusBadRequest, gin.H{"error": "商品已售出或下架"})
		return
	}
	if product.UserID == uid {
		tx.Rollback()
		c.JSON(http.StatusBadRequest, gin.H{"error": "不能购买自己的商品"})
		return
	}

	// 4. 创建订单
	order := models.Order{
		OrderNo:   generateOrderNo(),
		UserID:    uid,
		SellerID:  product.UserID,
		ProductID: product.ID,
		Price:     product.Price,
		Status:    1, // 1: 待支付
	}

	if err := tx.Create(&order).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建订单失败"})
		return
	}

	// 5. 更新商品状态为已售出 (2)
	if err := tx.Model(&product).Update("status", 2).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新商品状态失败"})
		return
	}

	tx.Commit()
	c.JSON(http.StatusOK, gin.H{"message": "下单成功", "data": order})
}

func (o *OrderController) List(c *gin.Context) {
	userID, _ := c.Get("userID")
	role := c.Query("role")

	var orders []models.Order

	// ★★★ 核心修复：必须 Preload 三层数据 ★★★
	// 1. Product: 获取商品名、图片
	// 2. Seller: 获取卖家头像、昵称
	// 3. User: 获取买家信息
	db := config.DB.Preload("Product").Preload("User").Preload("Seller")

	if role == "seller" {
		db = db.Where("seller_id = ?", userID)
	} else {
		db = db.Where("user_id = ?", userID)
	}

	// 按时间倒序
	if err := db.Order("created_at desc").Find(&orders).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取订单失败"})
		return
	}

	// 图片路径处理 (防止前端裂图)
	for i := range orders {
		if orders[i].Product.Image == "" {
			orders[i].Product.Image = "/uploads/default_product.png"
		}
	}

	c.JSON(http.StatusOK, gin.H{"data": orders})
}

// Pay 模拟支付
func (o *OrderController) Pay(c *gin.Context) {
	id := c.Param("id")
	userID, _ := c.Get("userID") // 获取当前登录用户 ID

	var order models.Order
	if err := config.DB.First(&order, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "订单不存在"})
		return
	}

	// ★★★ 关键检查：只能支付自己的订单 ★★★
	if order.UserID != userID.(uint) {
		// 如果你是用管理员账号测试普通用户的订单，这里会报错
		c.JSON(http.StatusForbidden, gin.H{"error": "无权操作此订单"})
		return
	}

	// ★★★ 关键检查：防止重复支付 ★★★
	if order.Status != 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "订单状态不正确，可能已支付"})
		return
	}

	// 更新为待发货 (2)
	config.DB.Model(&order).Update("status", 2)
	c.JSON(http.StatusOK, gin.H{"message": "支付成功"})
}

// BatchCreate 购物车批量结算 (★★★ 核心修复实现 ★★★)
func (o *OrderController) BatchCreate(c *gin.Context) {
	// 1. 定义接收格式
	var input struct {
		CartIDs []uint `json:"cart_ids"` // 前端传来的购物车ID数组
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	if len(input.CartIDs) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请选择要结算的商品"})
		return
	}

	userID, _ := c.Get("userID")
	uid := userID.(uint)

	// 2. 开启事务 (Transaction)
	tx := config.DB.Begin()

	// 准备一个数组存放生成的订单，用于返回给前端
	var createdOrders []models.Order

	for _, cartID := range input.CartIDs {
		// A. 查找购物车记录 (确保是自己的)
		var cartItem models.Cart
		if err := tx.Preload("Product").Where("id = ? AND user_id = ?", cartID, uid).First(&cartItem).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusNotFound, gin.H{"error": "购物车记录不存在或已删除"})
			return
		}

		// B. 检查商品状态
		if cartItem.Product.Status != 1 {
			tx.Rollback()
			c.JSON(http.StatusBadRequest, gin.H{"error": "商品 [" + cartItem.Product.Name + "] 已失效，无法购买"})
			return
		}

		// C. 防自己买自己
		if cartItem.Product.UserID == uid {
			tx.Rollback()
			c.JSON(http.StatusBadRequest, gin.H{"error": "不能购买自己的商品"})
			return
		}

		// D. 创建订单
		order := models.Order{
			OrderNo:   generateOrderNo() + strconv.Itoa(int(cartID)), // 防止高并发下订单号冲突
			UserID:    uid,
			SellerID:  cartItem.Product.UserID,
			ProductID: cartItem.Product.ID,
			Price:     cartItem.Product.Price,
			Status:    1, // 待支付
		}

		if err := tx.Create(&order).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "创建订单失败"})
			return
		}
		createdOrders = append(createdOrders, order)

		// E. 标记商品为已售出 (2)
		if err := tx.Model(&cartItem.Product).Update("status", 2).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "锁定商品失败"})
			return
		}

		// F. 从购物车移除该条目
		if err := tx.Delete(&cartItem).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "移除购物车失败"})
			return
		}
	}

	// 3. 提交事务
	tx.Commit()

	c.JSON(http.StatusOK, gin.H{"message": "结算成功", "data": createdOrders})
}

func generateOrderNo() string {
	return time.Now().Format("20060102150405") + "001"
}
