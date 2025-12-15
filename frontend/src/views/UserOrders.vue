<template>
  <div class="orders-page" v-loading="loading">

    <div class="nav-header">
      <div class="container nav-content">
        <div class="back-btn" @click="$router.push('/profile')">
          <el-icon><ArrowLeft /></el-icon> 返回个人中心
        </div>
        <div class="page-title">我的订单</div>
        <div class="placeholder"></div>
      </div>
    </div>

    <div class="tabs-bar">
      <div class="container tabs-inner">
        <div
            v-for="tab in tabs"
            :key="tab.value"
            class="tab-item"
            :class="{ active: currentTab === tab.value }"
            @click="currentTab = tab.value"
        >
          {{ tab.label }}
          <div class="active-line" v-if="currentTab === tab.value"></div>
        </div>
      </div>
    </div>

    <div class="container main-content">
      <div v-if="filteredOrders.length === 0 && !loading" class="empty-state">
        <el-empty description="暂无相关订单" :image-size="160" />
        <button class="go-home-btn" @click="$router.push('/')">去逛逛</button>
      </div>

      <div v-else class="order-list">
        <div v-for="order in filteredOrders" :key="order.id" class="order-card animate-up">

          <div class="card-header">
            <div class="shop-info" @click="contactSeller(order.seller?.id)">
              <el-avatar :size="24" :src="order.seller?.avatar || defaultAvatar" />
              <span class="shop-name">{{ order.seller?.nickname || order.seller?.username || '闲趣卖家' }}</span>
              <el-icon><ArrowRight /></el-icon>
            </div>
            <div class="status-badge" :class="getStatusClass(order.status)">
              {{ getStatusText(order.status) }}
            </div>
          </div>

          <div class="card-body" @click="$router.push(`/product/${order.product_id}`)">
            <div class="img-box">
              <el-image :src="fixUrl(order.product?.image)" fit="cover" class="prod-img" lazy>
                <template #error><div class="err-slot"><el-icon><Picture /></el-icon></div></template>
              </el-image>
            </div>

            <div class="info-box">
              <div class="prod-title">{{ order.product?.name || order.product?.title || '商品信息已失效' }}</div>
              <div class="prod-desc">{{ order.product?.description || '暂无描述...' }}</div>
              <div class="tags">
                <span class="tag">包邮</span>
                <span class="tag">担保交易</span>
              </div>
            </div>

            <div class="price-box">
              <div class="price"><span class="symbol">¥</span>{{ order.price }}</div>
              <div class="qty">x 1</div>
            </div>
          </div>

          <div class="card-footer">
            <div class="total-info">
              实付: <span class="amount">¥ {{ order.price }}</span>
            </div>

            <div class="actions">
              <button class="btn btn-contact" @click="contactSeller(order.seller?.id)">联系卖家</button>

              <button
                  v-if="order.status === 1"
                  class="btn btn-primary"
                  @click="payNow(order)"
              >
                立即支付
              </button>

              <button
                  v-if="order.status === 2 || order.status === 3"
                  class="btn btn-confirm"
                  @click="confirmReceive(order.id)"
              >
                确认收货
              </button>

              <button v-if="order.status === 4" class="btn btn-outline">评价</button>
            </div>
          </div>
        </div>
      </div>
    </div>

    <PaymentModal v-model="payVisible" :order="currentPayOrder" @success="fetchOrders" />
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import request from '@/utils/request'
import { useRouter } from 'vue-router'
import { ArrowLeft, ArrowRight, Picture } from '@element-plus/icons-vue'
import PaymentModal from '../components/PaymentModal.vue'
import { ElMessage } from 'element-plus'

const router = useRouter()
const loading = ref(false)
const orders = ref([])
// 0=全部, 1=待支付, 2=待发货/运输中, 4=已完成
const currentTab = ref(0)
const defaultAvatar = 'https://cube.elemecdn.com/3/7c/3ea6beec64369c2642b92c6726f1epng.png'
const payVisible = ref(false)
const currentPayOrder = ref({})

const tabs = [
  { label: '全部', value: 0 },
  { label: '待付款', value: 1 },
  { label: '待发货', value: 2 },
  { label: '已完成', value: 4 }
]

// 修复图片路径
const fixUrl = (url) => {
  if (!url) return ''
  if (!url.startsWith('http')) return 'http://127.0.0.1:8081' + url
  return url.replace('localhost', '127.0.0.1')
}

const fetchOrders = async () => {
  loading.value = true
  try {
    const res = await request.get('/api/orders')
    orders.value = res.data || []
  } catch (e) {
    console.error(e)
  } finally {
    loading.value = false
  }
}

// 前端过滤 (后端目前是返回全部，如果数据量大建议后端过滤)
const filteredOrders = computed(() => {
  if (currentTab.value === 0) return orders.value

  // 待发货 Tab 同时显示 status 2(待发货) 和 3(运输中)
  if (currentTab.value === 2) {
    return orders.value.filter(o => o.status === 2 || o.status === 3)
  }

  return orders.value.filter(o => o.status === currentTab.value)
})

const getStatusText = (status) => {
  const map = { 1: '待付款', 2: '待发货', 3: '运输中', 4: '交易成功', 5: '已取消' }
  return map[status] || '未知状态'
}

const getStatusClass = (status) => {
  if (status === 1) return 'wait'
  if (status === 4) return 'success'
  if (status === 5) return 'cancel'
  return 'normal'
}

const contactSeller = (sellerId) => {
  if (sellerId) router.push(`/chat/${sellerId}`)
}

const payNow = (order) => {
  currentPayOrder.value = {
    id: order.id,           // 单个订单 ID
    amount: order.price,    // 订单金额
    order_no: order.order_no, // 订单号(用于生成二维码)
    isBatch: false          // 明确标记为非批量
  }
  payVisible.value = true
}

const confirmReceive = async (id) => {
  // 预留接口，后端暂未实现 confirm 接口，这里仅做 UI 演示
  // 实际应调用 await request.put(`/api/orders/${id}/confirm`)
  ElMessage.success('收货成功！交易完成')
}

onMounted(fetchOrders)
</script>

<style scoped lang="scss">
$primary: #ffdf5d;
$bg-page: #f6f7f9;

.orders-page { min-height: 100vh; background: $bg-page; padding-bottom: 40px; }

/* 顶部导航 */
.nav-header {
  height: 60px; background: #fff; position: sticky; top: 0; z-index: 100; box-shadow: 0 2px 8px rgba(0,0,0,0.02);
  .nav-content { height: 100%; display: flex; align-items: center; justify-content: space-between; }
  .back-btn { cursor: pointer; display: flex; align-items: center; gap: 4px; font-size: 14px; font-weight: 500; color: #666; &:hover { color: #333; } }
  .page-title { font-weight: 800; font-size: 18px; color: #333; }
  .placeholder { width: 80px; }
}

/* Tabs */
.tabs-bar {
  background: #fff; margin-bottom: 20px; border-top: 1px solid #f9f9f9;
  .tabs-inner { display: flex; gap: 30px; }
  .tab-item {
    padding: 14px 0; font-size: 15px; color: #666; cursor: pointer; position: relative; transition: 0.2s;
    &:hover { color: #333; }
    &.active { font-weight: 800; color: #333; font-size: 16px; }
    .active-line { position: absolute; bottom: 0; left: 50%; transform: translateX(-50%); width: 24px; height: 3px; background: $primary; border-radius: 2px; }
  }
}

.container { max-width: 800px; margin: 0 auto; padding: 0 16px; }

/* 订单卡片 */
.order-card {
  background: #fff; border-radius: 16px; padding: 20px; margin-bottom: 16px; box-shadow: 0 4px 12px rgba(0,0,0,0.02); transition: 0.2s;
  &:hover { box-shadow: 0 8px 20px rgba(0,0,0,0.05); transform: translateY(-2px); }

  .card-header {
    display: flex; justify-content: space-between; align-items: center; padding-bottom: 14px; border-bottom: 1px dashed #eee; margin-bottom: 14px;
    .shop-info { display: flex; align-items: center; gap: 8px; cursor: pointer; font-size: 14px; font-weight: bold; color: #333; &:hover { opacity: 0.8; } }
    .status-badge { font-size: 13px; font-weight: bold; &.wait { color: #ff5000; } &.success { color: #52c41a; } &.cancel { color: #999; } &.normal { color: #1677ff; } }
  }

  .card-body {
    display: flex; gap: 16px; cursor: pointer;
    .img-box { width: 90px; height: 90px; border-radius: 12px; overflow: hidden; background: #f9f9f9; .prod-img { width: 100%; height: 100%; transition: transform 0.3s; &:hover { transform: scale(1.05); } } }
    .info-box { flex: 1; .prod-title { font-size: 15px; font-weight: bold; color: #333; margin-bottom: 6px; line-height: 1.4; max-height: 42px; overflow: hidden; } .prod-desc { font-size: 13px; color: #999; white-space: nowrap; overflow: hidden; text-overflow: ellipsis; margin-bottom: 8px; } .tags .tag { font-size: 10px; background: #fffbe6; color: #d48806; padding: 2px 6px; border-radius: 4px; margin-right: 6px; } }
    .price-box { text-align: right; .price { font-weight: 800; font-size: 16px; color: #333; } .qty { font-size: 12px; color: #999; margin-top: 4px; } }
  }

  .card-footer {
    display: flex; justify-content: space-between; align-items: center; margin-top: 20px; padding-top: 10px;
    .total-info { font-size: 13px; color: #666; .amount { font-size: 18px; font-weight: 800; color: #ff5000; margin-left: 4px; } }
    .actions {
      display: flex; gap: 10px;
      .btn { padding: 8px 20px; border-radius: 99px; font-size: 13px; cursor: pointer; border: 1px solid #ddd; background: #fff; font-weight: 600; transition: 0.2s; }
      .btn-contact { color: #666; &:hover { border-color: #333; color: #333; } }
      .btn-primary { border: none; background: #1a1a1a; color: $primary; &:hover { background: #000; transform: translateY(-1px); box-shadow: 0 4px 10px rgba(0,0,0,0.15); } }
      .btn-confirm { border-color: $primary; color: #d48806; background: #fffbe6; &:hover { background: #fff1b8; } }
      .btn-outline { &:hover { background: #f5f5f5; } }
    }
  }
}

.empty-state { display: flex; flex-direction: column; align-items: center; padding-top: 60px; .go-home-btn { margin-top: 20px; padding: 10px 32px; background: #333; color: $primary; border: none; border-radius: 99px; font-weight: bold; cursor: pointer; transition: 0.2s; &:hover { transform: translateY(-2px); box-shadow: 0 4px 12px rgba(0,0,0,0.15); } } }

.animate-up { animation: fadeInUp 0.4s ease; } @keyframes fadeInUp { from { opacity: 0; transform: translateY(10px); } to { opacity: 1; transform: translateY(0); } }
</style>