<template>
  <div class="cart-page">
    <div class="bg-orb orb-1"></div>
    <div class="bg-orb orb-2"></div>

    <div class="container">
      <div class="page-header glass-panel">
        <div class="left" @click="$router.push('/')">
          <el-icon><ArrowLeft /></el-icon> <span>返回首页</span>
        </div>
        <div class="center">购物车 <span class="count">({{ cartList.length }})</span></div>
        <div class="right" @click="isEditMode = !isEditMode">{{ isEditMode ? '完成' : '管理' }}</div>
      </div>

      <div class="cart-content" v-loading="loading">
        <div v-if="cartList.length > 0" class="cart-list">
          <transition-group name="list-anim">
            <div
                v-for="item in cartList"
                :key="item.id"
                class="cart-item-card"
                :class="{ 'is-selected': selectedIds.includes(item.id) }"
                @click="toggleSelect(item.id)"
            >
              <div class="checkbox-area" @click.stop="toggleSelect(item.id)">
                <div class="custom-checkbox" :class="{ checked: selectedIds.includes(item.id) }">
                  <el-icon v-if="selectedIds.includes(item.id)"><Check /></el-icon>
                </div>
              </div>

              <div class="image-wrapper" @click.stop="$router.push(`/product/${item.product?.id}`)">
                <el-image :src="item.product?.image || defaultImg" fit="cover" class="thumb" loading="lazy" />
                <div v-if="item.product?.status !== 1" class="invalid-mask">失效</div>
              </div>

              <div class="info-wrapper">
                <div class="title-row">
                  <h3 class="product-title">{{ item.product?.name || item.product?.title || '未知商品' }}</h3>
                  <el-button
                      v-if="isEditMode"
                      type="danger"
                      circle
                      size="small"
                      class="delete-btn"
                      @click.stop="deleteItem(item.id)"
                  >
                    <el-icon><Delete /></el-icon>
                  </el-button>
                </div>
                <div class="meta-row">
                  <div class="seller-tag">
                    <el-avatar :size="16" :src="item.product?.seller?.avatar || defaultAvatar" />
                    <span>{{ getSellerName(item.product) }}</span>
                  </div>
                </div>
                <div class="price-row">
                  <div class="price"><span class="symbol">¥</span><span class="num">{{ item.product?.price }}</span></div>
                  <span class="qty-tag">x{{ item.count || 1 }}</span>
                </div>
              </div>
            </div>
          </transition-group>
        </div>

        <div v-else-if="!loading" class="empty-cart glass-panel">
          <img src="https://cube.elemecdn.com/e/fd/0fc7d20532fdaf769a25683617711png.png" class="empty-img" />
          <p>购物车空空如也，去逛逛吧~</p>
          <el-button round color="#333" class="go-home-btn" @click="$router.push('/')">去首页</el-button>
        </div>
      </div>
    </div>

    <div class="checkout-bar-wrapper" v-if="cartList.length > 0">
      <div class="checkout-bar glass-panel">
        <div class="select-all" @click="toggleSelectAll">
          <div class="custom-checkbox" :class="{ checked: isAllSelected }">
            <el-icon v-if="isAllSelected"><Check /></el-icon>
          </div>
          <span>全选</span>
        </div>

        <div class="right-section" v-if="!isEditMode">
          <div class="total-section">
            <span class="label">合计:</span>
            <span class="total-price"><span class="symbol">¥</span>{{ totalPrice }}</span>
          </div>
          <button
              class="btn-checkout"
              :disabled="selectedIds.length === 0"
              @click="handleCheckout"
          >
            立即结算 ({{ selectedIds.length }})
          </button>
        </div>

        <div class="right-section" v-else>
          <button
              class="btn-delete-all"
              :disabled="selectedIds.length === 0"
              @click="handleBatchDelete"
          >
            删除 ({{ selectedIds.length }})
          </button>
        </div>
      </div>
    </div>

    <PaymentModal
        v-model="payVisible"
        :order="payOrderInfo"
        @success="onPaySuccess"
    />
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import request from '@/utils/request'
import { useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { ArrowLeft, Delete, Check } from '@element-plus/icons-vue'
import PaymentModal from '../components/PaymentModal.vue'

const router = useRouter()
const loading = ref(false)
const cartList = ref([])
const selectedIds = ref([])
const isEditMode = ref(false)
const user = ref(JSON.parse(localStorage.getItem('user') || '{}'))
const defaultAvatar = 'https://cube.elemecdn.com/3/7c/3ea6beec64369c2642b92c6726f1epng.png'
const defaultImg = 'https://via.placeholder.com/150?text=No+Image'

const payVisible = ref(false)
const payOrderInfo = ref({})

// 获取卖家名字
const getSellerName = (p) => {
  if (!p) return '未知卖家'
  return p.seller?.nickname || p.seller?.username || p.user?.nickname || '闲趣用户'
}

const fetchCart = async () => {
  if (!user.value.id) return router.push('/')
  loading.value = true
  try {
    const res = await request.get('/api/cart', { params: { user_id: user.value.id } })

    // 图片路径修复逻辑
    cartList.value = (res.data || []).map(item => {
      if (item.product?.image && !item.product.image.startsWith('http')) {
        item.product.image = 'http://127.0.0.1:8081' + item.product.image
      }
      if (item.product?.image) {
        item.product.image = item.product.image.replace('localhost', '127.0.0.1')
      }
      return item
    })
  } catch (e) {
    console.error(e)
  } finally {
    loading.value = false
  }
}

// ★★★ 修复点3：计算总价时乘以数量 ★★★
const totalPrice = computed(() => {
  let sum = 0
  cartList.value.forEach(item => {
    if (selectedIds.value.includes(item.id) && item.product?.status === 1) {
      sum += Number(item.product.price) * (item.count || 1)
    }
  })
  return sum.toFixed(2)
})

const isAllSelected = computed(() => {
  // 只统计有效商品
  const validItems = cartList.value.filter(i => isEditMode.value || i.product?.status === 1)
  if (validItems.length === 0) return false
  return validItems.every(i => selectedIds.value.includes(i.id))
})

const toggleSelect = (id) => {
  const item = cartList.value.find(i => i.id === id)
  // 非管理模式下，失效商品不能选
  if (!isEditMode.value && item?.product?.status !== 1) return

  const index = selectedIds.value.indexOf(id)
  index > -1 ? selectedIds.value.splice(index, 1) : selectedIds.value.push(id)
}

const toggleSelectAll = () => {
  if (isAllSelected.value) {
    selectedIds.value = []
  } else {
    selectedIds.value = cartList.value
        .filter(i => isEditMode.value || i.product?.status === 1)
        .map(i => i.id)
  }
}

const deleteItem = (id) => {
  ElMessageBox.confirm('确定要把这件宝贝移出购物车吗？', '提示', {
    confirmButtonText: '移除', cancelButtonText: '取消', type: 'warning', center: true, customClass: 'warm-theme-box'
  }).then(async () => {
    await request.delete(`/api/cart/${id}`)
    ElMessage.success('已移除')
    cartList.value = cartList.value.filter(i => i.id !== id)
    selectedIds.value = selectedIds.value.filter(sid => sid !== id)
  }).catch(() => {})
}

// 批量删除
const handleBatchDelete = async () => {
  if (selectedIds.value.length === 0) return
  try {
    await ElMessageBox.confirm(`确定删除选中的 ${selectedIds.value.length} 件商品吗？`, '提示', { type: 'warning' })
    for (const id of selectedIds.value) {
      await request.delete(`/api/cart/${id}`)
    }
    ElMessage.success('已删除')
    selectedIds.value = []
    fetchCart()
  } catch (e) {}
}

// ★★★ 核心修复：批量下单逻辑 ★★★
const handleCheckout = async () => {
  if (selectedIds.value.length === 0) return

  loading.value = true
  try {
    // 1. 调用后端批量下单接口
    const res = await request.post('/api/orders/batch', {
      cart_ids: selectedIds.value,
      address: "默认收货地址"
    })

    // 2. 获取生成的订单数据
    const orders = res.data || []
    if (!orders || orders.length === 0) throw new Error('下单异常，未返回订单信息')

    // 提取 ID 数组，供 PaymentModal 批量支付使用
    const orderIds = orders.map(o => o.id)
    const totalAmount = orders.reduce((sum, o) => sum + o.price, 0)

    // 3. 打开支付弹窗
    payOrderInfo.value = {
      ids: orderIds,      // 传入数组，支持批量支付
      amount: totalAmount,
      isBatch: true,
      order_no: orders[0].order_no // 取第一个订单号用于展示二维码
    }
    payVisible.value = true

  } catch (e) {
    console.error(e)
    ElMessage.error(e.response?.data?.error || '结算失败，请稍后重试')
  } finally {
    loading.value = false
  }
}

// 支付成功回调
const onPaySuccess = () => {
  selectedIds.value = []
  fetchCart() // 重新拉取购物车(已结算的会被后端移除)
  // 延迟跳转，提升体验
  setTimeout(() => {
    router.push('/orders')
  }, 500)
}

onMounted(fetchCart)
</script>

<style scoped lang="scss">
$primary: #ffdf5d;
$bg-page: #f2f4f7;

.cart-page { min-height: 100vh; background-color: $bg-page; padding-bottom: 100px; position: relative; overflow-x: hidden; }
.bg-orb { position: fixed; border-radius: 50%; filter: blur(80px); opacity: 0.5; z-index: 0; pointer-events: none;
  &.orb-1 { width: 300px; height: 300px; background: #ffdf5d; top: -100px; left: -50px; }
  &.orb-2 { width: 250px; height: 250px; background: #a0cfff; bottom: 100px; right: -80px; }
}
.container { max-width: 800px; margin: 0 auto; padding: 20px; position: relative; z-index: 1; }

.glass-panel { background: rgba(255, 255, 255, 0.75); backdrop-filter: blur(20px); -webkit-backdrop-filter: blur(20px); border: 1px solid rgba(255, 255, 255, 0.6); box-shadow: 0 8px 32px rgba(0, 0, 0, 0.03); }

.page-header { display: flex; justify-content: space-between; align-items: center; padding: 0 24px; height: 60px; margin-bottom: 24px; border-radius: 99px;
  .left { display: flex; align-items: center; gap: 4px; cursor: pointer; font-size: 14px; color: #666; transition: 0.2s; &:hover { color: #000; } }
  .center { font-size: 18px; font-weight: 800; color: #333; .count { font-size: 14px; font-weight: normal; color: #999; } }
  .right { font-size: 14px; font-weight: bold; cursor: pointer; color: #333; &:hover { color: $primary; } }
}

.cart-list { display: flex; flex-direction: column; gap: 16px; }
.cart-item-card { display: flex; align-items: center; padding: 16px; background: #fff; border-radius: 20px; transition: all 0.3s cubic-bezier(0.25, 0.8, 0.25, 1); cursor: pointer; border: 1px solid transparent;
  &:hover { transform: translateY(-2px); box-shadow: 0 8px 20px rgba(0,0,0,0.04); }
  &.is-selected { border-color: $primary; background: #fffdf0; }
  .checkbox-area { margin-right: 16px; .custom-checkbox { width: 22px; height: 22px; border-radius: 50%; border: 2px solid #ddd; display: flex; align-items: center; justify-content: center; color: #fff; transition: 0.2s; &.checked { background: #000; border-color: #000; } } }
  .image-wrapper { width: 100px; height: 100px; border-radius: 12px; overflow: hidden; background: #f9f9f9; margin-right: 16px; position: relative; .thumb { width: 100%; height: 100%; object-fit: cover; } .invalid-mask { position: absolute; inset: 0; background: rgba(0,0,0,0.6); color: #fff; display: flex; align-items: center; justify-content: center; font-size: 12px; } }
  .info-wrapper { flex: 1; display: flex; flex-direction: column; justify-content: space-between; height: 100px; padding: 2px 0;
    .title-row { display: flex; justify-content: space-between; align-items: flex-start; .product-title { margin: 0; font-size: 16px; font-weight: bold; color: #333; line-height: 1.4; display: -webkit-box; -webkit-line-clamp: 2; -webkit-box-orient: vertical; overflow: hidden; max-width: 300px; } .delete-btn { color: #ff4d4f; background: #fff0f0; border: none; &:hover { background: #ff4d4f; color: #fff; } } }
    .meta-row { .seller-tag { display: inline-flex; align-items: center; gap: 6px; background: #f5f5f5; padding: 2px 8px; border-radius: 6px; font-size: 12px; color: #666; } }
    .price-row { display: flex; justify-content: space-between; align-items: flex-end; .price { color: #ff5000; font-weight: 900; .symbol { font-size: 14px; margin-right: 2px; } .num { font-size: 20px; } } .qty-tag { color: #999; font-size: 12px; background: #f5f5f5; padding: 2px 6px; border-radius: 4px; } }
  }
}

.empty-cart { text-align: center; padding: 60px 0; border-radius: 24px; .empty-img { width: 160px; margin-bottom: 20px; opacity: 0.8; } p { color: #999; margin-bottom: 24px; font-size: 15px; } .go-home-btn { padding: 12px 32px; font-weight: bold; } }

.checkout-bar-wrapper { position: fixed; bottom: 30px; left: 0; right: 0; z-index: 99; display: flex; justify-content: center; pointer-events: none; }
.checkout-bar { pointer-events: auto; width: 100%; max-width: 800px; margin: 0 20px; height: 72px; border-radius: 99px; display: flex; align-items: center; justify-content: space-between; padding: 0 10px 0 24px;
  .select-all { display: flex; align-items: center; gap: 10px; cursor: pointer; font-size: 14px; font-weight: 600; .custom-checkbox { width: 22px; height: 22px; border-radius: 50%; border: 2px solid #999; display: flex; align-items: center; justify-content: center; color: #fff; &.checked { background: #000; border-color: #000; } } }

  .right-section {
    display: flex; align-items: center; gap: 16px;
    .total-section { display: flex; align-items: baseline; gap: 8px; .label { font-size: 14px; color: #666; } .total-price { color: #ff5000; font-weight: 900; font-size: 24px; font-family: 'DIN Alternate', sans-serif; .symbol { font-size: 16px; margin-right: 2px; } } }
    .btn-checkout { height: 52px; padding: 0 40px; background: #1a1a1a; color: $primary; border: none; border-radius: 99px; font-size: 16px; font-weight: bold; cursor: pointer; transition: 0.3s; box-shadow: 0 8px 20px rgba(0,0,0,0.2); &:hover { transform: translateY(-2px); background: #000; box-shadow: 0 12px 24px rgba(0,0,0,0.25); } &:disabled { background: #ccc; color: #fff; cursor: not-allowed; transform: none; box-shadow: none; } }
    .btn-delete-all { height: 44px; padding: 0 24px; background: #ff4d4f; color: #fff; border-radius: 22px; border: none; font-weight: bold; cursor: pointer; transition: 0.2s; &:hover { background: #ff7875; } &:disabled { opacity: 0.5; cursor: not-allowed; } }
  }
}

.list-anim-enter-active, .list-anim-leave-active { transition: all 0.4s ease; }
.list-anim-enter-from { opacity: 0; transform: translateY(20px); }
.list-anim-leave-to { opacity: 0; transform: translateX(-100%); }
</style>