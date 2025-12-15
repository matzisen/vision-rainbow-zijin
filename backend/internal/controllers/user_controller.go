package controllers

import (
	"gotest/config"
	"gotest/internal/middleware" // 确保导入了 middleware 包
	"gotest/internal/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type UserController struct{}

// Register 注册
func (u *UserController) Register(c *gin.Context) {
	var input struct {
		Username string `json:"username"`
		Password string `json:"password"`
		Phone    string `json:"phone"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "数据格式错误"})
		return
	}

	// 1. 检查用户名是否已存在
	var count int64
	config.DB.Model(&models.User{}).Where("username = ?", input.Username).Count(&count)
	if count > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "用户名已存在"})
		return
	}

	// 2. 密码加密
	hash, _ := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)

	// 3. 创建用户
	user := models.User{
		Username: input.Username,
		Password: string(hash),
		Phone:    input.Phone,
		Role:     "user",
		Avatar:   "https://cube.elemecdn.com/3/7c/3ea6beec64369c2642b92c6726f1epng.png", // 默认头像
	}

	if err := config.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "注册失败，请重试"})
		return
	}

	// 4. 生成真实 Token
	token, err := middleware.GenerateToken(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Token生成失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "注册成功", "token": token, "user": user})
}

// Login 登录
func (u *UserController) Login(c *gin.Context) {
	var input struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	// 1. 查找用户
	var user models.User
	if err := config.DB.Where("username = ?", input.Username).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "账号或密码错误"})
		return
	}

	// 2. 验证密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "账号或密码错误"})
		return
	}

	// 3. 生成真实 Token
	token, err := middleware.GenerateToken(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "系统错误: Token生成失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token, "user": user})
}

// GetMyData 获取我的发布/收藏
func (u *UserController) GetMyData(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未登录"})
		return
	}
	uid := userID.(uint)

	queryType := c.Query("type")

	if queryType == "products" {
		// 只查询当前用户发布的商品
		var products []models.Product
		config.DB.Where("user_id = ?", uid).Order("created_at desc").Find(&products)

		// 图片路径处理
		for i := range products {
			if products[i].Image == "" {
				products[i].Image = "/uploads/default_product.png"
			}
		}
		c.JSON(http.StatusOK, gin.H{"data": products})
	} else {
		// 获取我的收藏
		var favorites []models.Favorite
		config.DB.Preload("Product").Where("user_id = ?", uid).Find(&favorites)
		c.JSON(http.StatusOK, gin.H{"data": favorites})
	}
}

// UpdateProfile 更新资料
func (u *UserController) UpdateProfile(c *gin.Context) {
	userID, _ := c.Get("userID")
	var input models.User
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	var user models.User
	if err := config.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "用户不存在"})
		return
	}

	updates := map[string]interface{}{}
	if input.Nickname != "" {
		updates["nickname"] = input.Nickname
	}
	if input.Avatar != "" {
		updates["avatar"] = input.Avatar
	}
	if input.Phone != "" {
		updates["phone"] = input.Phone
	}

	config.DB.Model(&user).Updates(updates)
	c.JSON(http.StatusOK, gin.H{"message": "更新成功", "user": user})
}

// ChangePassword 修改密码
func (u *UserController) ChangePassword(c *gin.Context) {
	userID, _ := c.Get("userID")
	var input struct {
		OldPassword string `json:"old_password"`
		NewPassword string `json:"new_password"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	var user models.User
	config.DB.First(&user, userID)

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.OldPassword)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "旧密码错误"})
		return
	}

	hash, _ := bcrypt.GenerateFromPassword([]byte(input.NewPassword), bcrypt.DefaultCost)
	config.DB.Model(&user).Update("password", string(hash))

	c.JSON(http.StatusOK, gin.H{"message": "修改成功，请重新登录"})
}

// ToggleFavorite 切换收藏状态
func (u *UserController) ToggleFavorite(c *gin.Context) {
	var input struct {
		ProductID uint `json:"product_id"`
	}
	if c.ShouldBindJSON(&input) != nil {
		return
	}
	userID, _ := c.Get("userID")

	var fav models.Favorite
	if err := config.DB.Where("user_id = ? AND product_id = ?", userID, input.ProductID).First(&fav).Error; err == nil {
		config.DB.Delete(&fav)
		c.JSON(http.StatusOK, gin.H{"message": "已取消"})
	} else {
		config.DB.Create(&models.Favorite{UserID: userID.(uint), ProductID: input.ProductID})
		c.JSON(http.StatusOK, gin.H{"message": "已收藏"})
	}
}

// CheckFavorite 检查是否收藏
func (u *UserController) CheckFavorite(c *gin.Context) {
	userID, _ := c.Get("userID")
	productID := c.Query("product_id")
	var count int64
	config.DB.Model(&models.Favorite{}).Where("user_id = ? AND product_id = ?", userID, productID).Count(&count)
	c.JSON(http.StatusOK, gin.H{"is_favorite": count > 0})
}

// ★★★ 新增：GetUserInfo 获取指定用户的公开信息 ★★★
// 这就是修复“用户5”显示问题的关键接口
func (u *UserController) GetUserInfo(c *gin.Context) {
	id := c.Param("id")
	var user models.User

	// 只查询必要的公开字段，保护隐私
	result := config.DB.Select("id", "username", "nickname", "avatar").First(&user, id)

	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "用户不存在"})
		return
	}

	// 默认头像处理
	if user.Avatar == "" {
		user.Avatar = "https://cube.elemecdn.com/3/7c/3ea6beec64369c2642b92c6726f1epng.png"
	}

	c.JSON(http.StatusOK, gin.H{"data": user})
}
