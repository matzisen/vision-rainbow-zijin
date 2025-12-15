<template>
  <div class="product-detail-page" v-loading="loading">

    <div class="breadcrumb-bar">
      <div class="container flex-header">
        <div class="back-home-btn" @click="router.push('/')" title="返回首页">
          <el-icon><ArrowLeft /></el-icon>
        </div>
        <el-divider direction="vertical" class="nav-divider" />
        <el-breadcrumb separator="/">
          <el-breadcrumb-item :to="{ path: '/' }">首页</el-breadcrumb-item>
          <el-breadcrumb-item>闲置宝贝</el-breadcrumb-item>
          <el-breadcrumb-item>宝贝详情</el-breadcrumb-item>
        </el-breadcrumb>
      </div>
    </div>

    <div class="container">
      <div class="detail-card">
        <div class="left-gallery">
          <div class="main-image-box">
            <el-image
                :src="product.image || defaultImg"
                fit="contain"
                class="real-img"
                :preview-src-list="[product.image || defaultImg]"
                hide-on-click-modal
                :initial-index="0"
            >
              <template #error>
                <div class="image-slot"><el-icon><Picture /></el-icon></div>
              </template>
            </el-image>
            <div class="zoom-hint"><el-icon><ZoomIn /></el-icon> 点击查看大图</div>

            <div v-if="product.status !== 1" class="status-mask">
              <span class="mask-text">{{ product.status === 2 ? '已售出' : '已下架' }}</span>
            </div>
          </div>
        </div>

        <div class="right-info">
          <div class="info-header">
            <h1 class="product-title">{{ product.name || product.title || '未知商品' }}</h1>
            <div class="share-btn" @click="handleShare"><el-icon><Share /></el-icon></div>
          </div>

          <div class="price-section">
            <div class="main-price"><span class="symbol">¥</span><span class="amount">{{ product.price }}</span></div>
            <div class="view-count"><el-icon><View /></el-icon> {{ product.view_count || 1 }}人围观</div>
          </div>

          <div class="tags-row">
            <span v-if="product.is_free_shipping" class="tag-item highlight"><el-icon><Van /></el-icon> 包邮</span>
            <span v-if="product.is_negotiable" class="tag-item highlight"><el-icon><Scissor /></el-icon> 可小刀</span>
          </div>

          <div class="description-box">
            <div class="label">宝贝描述</div>
            <p class="content">{{ product.description || '卖家很懒，没有填写详细描述...' }}</p>
          </div>

          <div class="seller-card">
            <div class="seller-left">
              <el-avatar
                  :size="48"
                  :src="getSellerAvatar()"
                  class="avatar"
              />
              <div class="seller-text">
                <div class="name">{{ getSellerName() }}</div>
                <div class="credit-row">
                  <span class="credit-badge">信用极好</span>
                  <span class="credit-badge blue">实名认证</span>
                </div>
              </div>
            </div>
            <button class="btn-contact" @click="contactSeller">
              <el-icon><ChatDotRound /></el-icon> 联系卖家
            </button>
          </div>

          <div class="action-footer">
            <button class="btn-icon-action" :class="{ active: isFavorited }" @click="toggleFavorite">
              <el-icon v-if="isFavorited"><StarFilled /></el-icon>
              <el-icon v-else><Star /></el-icon>
              <span>{{ isFavorited ? '已收藏' : '收藏' }}</span>
            </button>

            <div class="btn-group">
              <button
                  class="btn-cart"
                  @click="addToCart"
                  :disabled="product.status !== 1"
              >
                加入购物车
              </button>

              <button
                  class="btn-buy"
                  @click="handleBuy"
                  :disabled="product.status !== 1"
              >
                {{ product.status === 1 ? '立即购买' : '无法购买' }}
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>

    <PaymentModal
        v-model="payVisible"
        :order="currentOrder"
        @success="handlePaySuccess"
    />
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import request from '@/utils/request'
import { ElMessage } from 'element-plus'
import { Share, View, Van, Scissor, ChatDotRound, Star, StarFilled, Picture, ArrowLeft, ZoomIn } from '@element-plus/icons-vue'
import PaymentModal from '../components/PaymentModal.vue'

const route = useRoute()
const router = useRouter()
const loading = ref(false)
const product = ref({})
const isFavorited = ref(false)
const user = ref(JSON.parse(localStorage.getItem('user') || '{}'))
const defaultAvatar = 'https://cube.elemecdn.com/3/7c/3ea6beec64369c2642b92c6726f1epng.png'
const defaultImg = 'https://via.placeholder.com/400x400?text=No+Image'

const payVisible = ref(false)
const currentOrder = ref({})

// 获取卖家头像
const getSellerAvatar = () => {
  const s = product.value.seller || product.value.user || {}
  return s.avatar || defaultAvatar
}

// 获取卖家名字
const getSellerName = () => {
  const s = product.value.seller || product.value.user || {}
  return s.nickname || s.username || '闲趣用户'
}

// 判断是否是自己
const isMe = (userId) => {
  return Number(userId) === Number(user.value.id)
}

const fetchProduct = async () => {
  loading.value = true
  try {
    const res = await request.get(`/api/products/${route.params.id}`)
    let data = res.data || {}

    // 图片路径修复
    if (data.image && !data.image.startsWith('http')) {
      data.image = 'http://127.0.0.1:8081' + data.image
      data.image = data.image.replace('localhost', '127.0.0.1')
    }

    product.value = data
    if (user.value.id) checkFavoriteStatus()
  } catch (e) {
    ElMessage.error('商品加载失败')
  } finally {
    loading.value = false
  }
}

const checkFavoriteStatus = async () => {
  try {
    const res = await request.get('/api/user/favorite/check', {
      params: { user_id: user.value.id, product_id: route.params.id }
    })
    isFavorited.value = !!res.is_favorite
  } catch (e) {}
}

const contactSeller = () => {
  if (!user.value.id) return ElMessage.warning('请先登录')
  const sellerId = product.value.seller?.id || product.value.user_id
  if (!sellerId) return ElMessage.error('卖家信息获取失败')
  if (isMe(sellerId)) return ElMessage.warning('这是您自己的商品')
  router.push(`/chat/${sellerId}`)
}

const handleShare = () => {
  navigator.clipboard.writeText(window.location.href).then(() => ElMessage.success('链接已复制'))
}

const toggleFavorite = async () => {
  if (!user.value.id) return ElMessage.warning('请先登录')
  const prev = isFavorited.value
  isFavorited.value = !prev
  try {
    await request.post('/api/user/favorite', { user_id: user.value.id, product_id: product.value.id })
    ElMessage.success(isFavorited.value ? '已收藏' : '取消收藏')
  } catch (e) { isFavorited.value = prev }
}

const addToCart = async () => {
  if (!user.value.id) return ElMessage.warning('请先登录')

  const sellerId = product.value.seller?.id || product.value.user_id
  if (isMe(sellerId)) {
    return ElMessage.warning('不能将自己的商品加入购物车')
  }

  try {
    // 默认数量 1
    await request.post('/api/cart', { user_id: user.value.id, product_id: product.value.id, count: 1 })
    ElMessage.success('已加入购物车')
  } catch (e) {
    // request.js 统一处理错误
  }
}

const handleBuy = async () => {
  if (!user.value.id) return ElMessage.warning('请先登录')

  const sellerId = product.value.seller?.id || product.value.user_id
  if (isMe(sellerId)) {
    return ElMessage.warning('不能购买自己发布的商品')
  }

  loading.value = true
  try {
    const res = await request.post('/api/orders', { product_id: product.value.id })
    const orderData = res.data

    currentOrder.value = {
      id: orderData.id,
      amount: orderData.price, // 映射 price -> amount
      order_no: orderData.order_no
    }

    payVisible.value = true
  } catch (e) {
    ElMessage.error(e.response?.data?.error || '下单失败')
  } finally {
    loading.value = false
  }
}

const handlePaySuccess = () => {
  fetchProduct() // 刷新商品状态
}

onMounted(fetchProduct)
</script>

<style scoped lang="scss">
$primary: #ffdf5d;
$bg-page: #f6f7f9;

.product-detail-page { min-height: 100vh; background: $bg-page; padding-bottom: 60px; }
.container { max-width: 1100px; margin: 0 auto; padding: 0 20px; box-sizing: border-box; }

/* 顶部导航 */
.breadcrumb-bar {
  height: 60px; background: #fff; border-bottom: 1px solid #eee; margin-bottom: 30px; box-shadow: 0 2px 8px rgba(0,0,0,0.02);
  .flex-header { height: 100%; display: flex; align-items: center; gap: 16px; }
  .back-home-btn {
    width: 36px; height: 36px; border-radius: 50%; background: #fff; border: 1px solid #eee;
    display: flex; align-items: center; justify-content: center; cursor: pointer; color: #333; transition: 0.3s;
    &:hover { background: $primary; border-color: $primary; transform: scale(1.1); }
  }
  .nav-divider { height: 16px; border-color: #ddd; }
}

/* 主卡片 */
.detail-card {
  display: grid; grid-template-columns: 480px 1fr; gap: 40px; background: #fff;
  border-radius: 24px; padding: 40px; box-shadow: 0 8px 30px rgba(0,0,0,0.04); align-items: start;
}

/* 左侧图片 */
.left-gallery .main-image-box {
  width: 100%; height: 480px; background: #f9f9f9; border-radius: 20px; border: 1px solid #f0f0f0;
  position: relative; overflow: hidden; display: flex; align-items: center; justify-content: center;
  .real-img { width: 100%; height: 100%; cursor: zoom-in; }
  .zoom-hint {
    position: absolute; bottom: 20px; left: 50%; transform: translateX(-50%);
    background: rgba(0,0,0,0.6); color: #fff; padding: 6px 16px; border-radius: 99px; font-size: 12px;
    opacity: 0; transition: 0.3s; display: flex; align-items: center; gap: 4px; pointer-events: none;
  }
  &:hover .zoom-hint { opacity: 1; bottom: 30px; }
  .status-mask {
    position: absolute; inset: 0; background: rgba(255,255,255,0.7); backdrop-filter: grayscale(100%);
    display: flex; align-items: center; justify-content: center; z-index: 10;
    .mask-text { border: 4px solid #333; color: #333; font-size: 32px; font-weight: 900; padding: 10px 30px; transform: rotate(-15deg); border-radius: 12px; }
  }
}

/* 右侧信息 */
.right-info {
  display: flex; flex-direction: column; height: 100%;
  .info-header {
    display: flex; justify-content: space-between; align-items: flex-start; margin-bottom: 16px;
    .product-title { font-size: 26px; font-weight: 800; color: #1a1a1a; line-height: 1.3; margin: 0; flex: 1; }
    .share-btn {
      width: 36px; height: 36px; border-radius: 50%; background: #f5f5f5; display: flex; align-items: center; justify-content: center;
      color: #666; cursor: pointer; margin-left: 16px; &:hover { background: #eee; }
    }
  }
  .price-section {
    display: flex; align-items: baseline; gap: 16px; margin-bottom: 20px; background: #fffbf0; padding: 16px 20px; border-radius: 12px;
    .main-price { color: #ff5000; font-weight: 900; .symbol { font-size: 20px; } .amount { font-size: 36px; } }
    .view-count { font-size: 13px; color: #d48806; margin-left: auto; display: flex; align-items: center; gap: 4px; }
  }
  .tags-row {
    display: flex; gap: 12px; margin-bottom: 30px; border-bottom: 1px dashed #eee; padding-bottom: 20px;
    .tag-item {
      display: flex; align-items: center; gap: 6px; font-size: 13px; padding: 6px 12px; border-radius: 6px; font-weight: bold;
      &.highlight { background: #fff8d6; color: #d48806; }
      &.gray { background: #f5f5f5; color: #999; font-weight: normal; }
    }
  }
  .description-box {
    flex: 1; margin-bottom: 30px;
    .label { font-size: 14px; font-weight: bold; color: #333; margin-bottom: 8px; }
    .content { font-size: 15px; color: #555; line-height: 1.6; white-space: pre-wrap; }
  }
  .seller-card {
    display: flex; justify-content: space-between; align-items: center; padding: 16px; background: #fff; border: 1px solid #eee; border-radius: 16px; margin-bottom: 30px;
    .seller-left {
      display: flex; align-items: center; gap: 12px;
      .avatar { border: 2px solid #f5f5f5; }
      .seller-text {
        display: flex; flex-direction: column; gap: 2px;
        .name { font-weight: bold; color: #333; font-size: 16px; }
        .credit-row { display: flex; gap: 6px; .credit-badge { font-size: 10px; background: #e6fdfb; color: #00b578; padding: 1px 6px; border-radius: 4px; } .blue { background: #e6f7ff; color: #1890ff; } }
      }
    }
    .btn-contact {
      border: 1px solid #ddd; background: #fff; padding: 8px 18px; border-radius: 99px; font-size: 13px; font-weight: bold; color: #333; cursor: pointer; display: flex; align-items: center; gap: 6px; transition: 0.2s;
      &:hover { background: #333; color: #fff; border-color: #333; }
    }
  }
  .action-footer {
    display: flex; align-items: center; gap: 20px; margin-top: auto;
    .btn-icon-action {
      display: flex; flex-direction: column; align-items: center; gap: 4px; border: none; background: transparent; cursor: pointer; color: #999; font-size: 12px;
      .el-icon { font-size: 24px; } &.active { color: #ff5000; } &:hover { color: #333; }
    }
    .btn-group {
      flex: 1; display: flex; gap: 12px; height: 48px;
      button { flex: 1; border-radius: 99px; border: none; font-size: 16px; font-weight: bold; cursor: pointer; transition: 0.2s; }
      .btn-cart { background: #333; color: #fff; &:hover { background: #000; } &:disabled { background: #eee; color: #999; cursor: not-allowed; } }
      .btn-buy { background: $primary; color: #000; &:hover { background: #f0cf4b; } &:disabled { background: #f5f5f5; color: #ccc; cursor: not-allowed; } }
    }
  }
}
@media (max-width: 768px) {
  .detail-card { grid-template-columns: 1fr; padding: 20px; }
  .action-footer { position: fixed; bottom: 0; left: 0; right: 0; background: #fff; padding: 10px 20px; z-index: 100; box-shadow: 0 -4px 10px rgba(0,0,0,0.05); }
}
</style>