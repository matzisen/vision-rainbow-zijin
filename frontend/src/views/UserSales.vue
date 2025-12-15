<template>
  <div class="order-page">
    <nav class="navbar">
      <div class="container navbar-inner">
        <div class="left" @click="$router.push('/')">
          <el-icon><ArrowLeft /></el-icon>
          <span class="title">我卖出的</span> </div>
      </div>
    </nav>

    <div class="tabs-wrapper">
      <div class="container">
        <div class="custom-tabs">
          <div v-for="tab in tabs" :key="tab.value" class="tab-item" :class="{ active: currentStatus === tab.value }" @click="switchTab(tab.value)">
            {{ tab.label }}
            <div class="slider" v-if="currentStatus === tab.value"></div>
          </div>
        </div>
      </div>
    </div>

    <div class="container content-area">
      <div class="order-list" v-loading="loading">
        <div v-for="order in orderList" :key="order.id" class="order-card animate-up">
          <div class="card-header">
            <div class="shop-info">
              <el-avatar :size="24" :src="order.buyer?.avatar || defaultAvatar" />
              <span class="shop-name">买家: {{ order.buyer?.nickname || '神秘买家' }}</span>
              <el-icon><ArrowRight /></el-icon>
            </div>
            <div class="order-status" :class="getStatusClass(order.status)">
              {{ getStatusText(order.status) }}
            </div>
          </div>

          <div class="card-body" @click="$router.push(`/product/${order.product?.id}`)">
            <img :src="order.product?.image" class="thumb" />
            <div class="info">
              <div class="title">{{ order.product?.title }}</div>
              <div class="tags">
                <span class="tag">出售中</span>
              </div>
            </div>
            <div class="price-col">
              <div class="price"><small>¥</small>{{ order.amount }}</div>
              <div class="qty">x1</div>
            </div>
          </div>

          <div class="card-footer">
            <div class="summary">订单号: {{ order.id }}</div>
            <div class="actions">
              <button class="btn-outline" @click="contactBuyer(order)">联系买家</button>

              <button class="btn-primary" v-if="order.status === 1" @click="shipOrder(order)">去发货</button>
              <button class="btn-outline" v-if="order.status === 2">等待收货</button>
            </div>
          </div>
        </div>
        <el-empty v-if="!loading && orderList.length === 0" description="您还没有卖出过宝贝" />
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import request from '@/utils/request'
import { useRouter } from 'vue-router'
import { Search, ArrowLeft, ArrowRight } from '@element-plus/icons-vue'
import { ElMessage, ElMessageBox } from 'element-plus'

const router = useRouter()
const loading = ref(false)
const orderList = ref([])
const currentStatus = ref(0)
const user = ref(JSON.parse(localStorage.getItem('user') || '{}'))
const defaultAvatar = 'https://cube.elemecdn.com/3/7c/3ea6beec64369c2642b92c6726f1epng.png'

const tabs = [
  { label: '全部', value: 0 },
  { label: '待发货', value: 1 },
  { label: '已发货', value: 2 },
  { label: '已完成', value: 3 }
]

const fetchOrders = async () => {
  if (!user.value.id) return router.push('/')
  loading.value = true
  try {
    const params = {
      user_id: user.value.id,
      status: currentStatus.value,
      role: 'seller' // ★★★ 关键参数：告诉后端我是卖家 ★★★
    }
    const res = await request.get('/api/orders', { params })

    // 图片路径修复逻辑 (同 UserOrders)
    orderList.value = (res.data || []).map(order => {
      if (order.product && order.product.image) {
        let img = order.product.image
        if (!img.startsWith('http')) { img = 'http://127.0.0.1:8081' + img }
        img = img.replace('localhost', '127.0.0.1')
        order.product.image = img
      }
      return order
    })
  } catch (e) { console.error(e) } finally { loading.value = false }
}

const switchTab = (val) => { currentStatus.value = val; fetchOrders() }

const contactBuyer = (order) => {
  // order.user_id 是买家ID
  if (order.user_id) {
    router.push(`/chat/${order.user_id}`)
  }
}

// 模拟发货操作
const shipOrder = (order) => {
  ElMessageBox.confirm('确认已发货？', '发货确认').then(() => {
    // 这里需要后端支持发货接口，暂时仅前端模拟
    ElMessage.success('发货成功')
    // 实际应调用: await request.put(`/api/orders/${order.id}/ship`)
  })
}

const getStatusText = (status) => {
  switch (status) {
    case 0: return '买家未付款';
    case 1: return '待发货';
    case 2: return '已发货';
    case 3: return '交易成功';
    default: return ''
  }
}
const getStatusClass = (status) => {
  switch (status) {
    case 3: return 'success';
    case 2: return 'warning';
    case 1: return 'danger'; // 待发货对卖家来说是重点
    default: return ''
  }
}

onMounted(fetchOrders)
</script>

<style scoped lang="scss">
/* 复用 UserOrders.vue 的样式，完全一致 */
$primary: #ffdf5d;
$bg: #f6f7f9;

.order-page { min-height: 100vh; background: $bg; padding-top: 110px; }
.container { max-width: 800px; margin: 0 auto; padding: 0 20px; }
.navbar {
  height: 60px; background: #fff; position: fixed; top: 0; left: 0; right: 0; z-index: 100; border-bottom: 1px solid #f0f0f0;
  .navbar-inner { height: 100%; display: flex; align-items: center; justify-content: space-between; }
  .left { display: flex; align-items: center; gap: 8px; cursor: pointer; font-weight: bold; font-size: 16px; &:hover { opacity: 0.7; } }
}
.tabs-wrapper { position: fixed; top: 60px; left: 0; right: 0; z-index: 99; background: #fff; box-shadow: 0 4px 12px rgba(0,0,0,0.03); }
.custom-tabs {
  display: flex; gap: 30px; padding: 0 10px;
  .tab-item {
    position: relative; padding: 14px 0; cursor: pointer; font-size: 15px; color: #666; transition: 0.3s;
    &:hover { color: #333; }
    &.active { font-weight: bold; color: #000; font-size: 16px; }
    .slider { position: absolute; bottom: 0; left: 50%; transform: translateX(-50%); width: 20px; height: 3px; background: $primary; border-radius: 3px; }
  }
}
.content-area { padding-top: 20px; padding-bottom: 40px; }
.order-card {
  background: #fff; border-radius: 16px; padding: 16px; margin-bottom: 16px; box-shadow: 0 2px 8px rgba(0,0,0,0.02); transition: 0.2s;
  &:hover { transform: translateY(-2px); box-shadow: 0 8px 16px rgba(0,0,0,0.05); }
}
.card-header {
  display: flex; justify-content: space-between; align-items: center; margin-bottom: 12px; padding-bottom: 12px; border-bottom: 1px solid #f9f9f9;
  .shop-info { display: flex; align-items: center; gap: 8px; font-size: 14px; font-weight: bold; }
  .order-status { font-size: 13px; color: #666; font-weight: 500;
    &.success { color: #00b578; }
    &.warning { color: #ff9500; }
    &.danger { color: #ff4d4f; }
  }
}
.card-body {
  display: flex; gap: 12px; cursor: pointer;
  .thumb { width: 80px; height: 80px; border-radius: 8px; object-fit: cover; background: #f9f9f9; }
  .info { flex: 1; .title { font-size: 14px; font-weight: bold; margin-bottom: 8px; line-height: 1.4; height: 40px; overflow: hidden; } .tags { display: flex; gap: 4px; .tag { font-size: 10px; background: #fef0f0; color: #ff4d4f; padding: 2px 6px; border-radius: 4px; } } }
  .price-col { text-align: right; .price { font-size: 16px; font-weight: bold; } .qty { font-size: 12px; color: #999; margin-top: 4px; } }
}
.card-footer {
  margin-top: 16px; display: flex; justify-content: space-between; align-items: center;
  .summary { font-size: 12px; color: #999; }
  .actions {
    display: flex; gap: 10px;
    button { padding: 6px 16px; border-radius: 99px; font-size: 13px; cursor: pointer; font-weight: 500; transition: 0.2s; background: #fff; border: 1px solid #ddd; }
    .btn-primary { background: #fff; border: 1px solid $primary; color: #d48806; font-weight: bold; &:hover { background: #fffcf0; } }
  }
}
.animate-up { animation: fadeInUp 0.5s ease backwards; }
@keyframes fadeInUp { from { opacity: 0; transform: translateY(20px); } to { opacity: 1; transform: translateY(0); } }
</style>