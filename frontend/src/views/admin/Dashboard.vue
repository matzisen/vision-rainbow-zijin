<template>
  <div class="dashboard-container">
    <div class="welcome-banner animate-up">
      <div class="welcome-content">
        <div class="text-group">
          <h2 class="greeting">早安，管理员！🌞</h2>
          <p class="subtitle">今天是 {{ currentDate }}，闲趣平台运行平稳。准备好开始一天的工作了吗？</p>
        </div>
        <img src="https://cube.elemecdn.com/0/88/03b0d39583f48206768a7534e55bcpng.png" class="admin-avatar" alt="Admin" />
      </div>
    </div>

    <el-row :gutter="24" class="data-row">
      <el-col :xs="24" :sm="12" :md="6" v-for="(item, index) in cardData" :key="index">
        <el-card shadow="hover" class="stat-card animate-up" :class="'delay-' + (index + 1)">
          <div class="stat-content">
            <div class="stat-info">
              <div class="stat-label">{{ item.label }}</div>
              <div class="stat-num">{{ item.value }}</div>
            </div>
            <div class="stat-icon" :class="item.colorClass">
              <el-icon><component :is="item.icon" /></el-icon>
            </div>
          </div>
          <div class="bg-circle" :class="item.colorClass"></div>
        </el-card>
      </el-col>
    </el-row>

    <el-row :gutter="24" class="bottom-row animate-up delay-5">
      <el-col :span="16" :xs="24">
        <el-card shadow="never" class="action-card">
          <template #header>
            <div class="card-header">
              <span><el-icon><Operation /></el-icon> 快捷工作台</span>
            </div>
          </template>
          <div class="quick-actions">
            <div class="action-item" @click="$router.push('/admin/products')">
              <div class="icon-box blue"><el-icon><Checked /></el-icon></div>
              <span>商品审核</span>
            </div>
            <div class="action-item" @click="$router.push('/admin/users')">
              <div class="icon-box green"><el-icon><User /></el-icon></div>
              <span>用户管理</span>
            </div>
            <div class="action-item" @click="$router.push('/admin/orders')">
              <div class="icon-box orange"><el-icon><List /></el-icon></div>
              <span>订单查询</span>
            </div>
            <div class="action-item">
              <div class="icon-box purple"><el-icon><Setting /></el-icon></div>
              <span>系统设置</span>
            </div>
          </div>
        </el-card>
      </el-col>

      <el-col :span="8" :xs="24">
        <el-card shadow="never" class="info-card">
          <template #header>
            <div class="card-header">
              <span><el-icon><InfoFilled /></el-icon> 系统公告</span>
            </div>
          </template>
          <ul class="notice-list">
            <li><el-tag size="small" type="danger">紧急</el-tag> 请尽快处理 3 个违规商品举报</li>
            <li><el-tag size="small" type="warning">通知</el-tag> 本周六凌晨进行服务器维护</li>
            <li><el-tag size="small" type="success">更新</el-tag> 后台管理系统 v1.0 上线</li>
          </ul>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup>
import { ref, onMounted, computed, shallowRef } from 'vue'
import request from '@/utils/request'
// 显式引入图标组件，避免 component :is 解析失败
import {
  User, Goods, List, Wallet,
  Checked, Setting, Operation, InfoFilled
} from '@element-plus/icons-vue'

const currentDate = new Date().toLocaleDateString('zh-CN', { month: 'long', day: 'numeric', weekday: 'long' })

const stats = ref({
  user_count: 0,
  product_count: 0,
  order_count: 0,
  trade_amount: 0
})

// 使用 computed 动态生成数据，确保 stats 更新后视图自动刷新
const cardData = computed(() => [
  { label: '总用户数', value: stats.value.user_count, icon: shallowRef(User), colorClass: 'icon-blue' },
  { label: '商品总数', value: stats.value.product_count, icon: shallowRef(Goods), colorClass: 'icon-orange' },
  { label: '订单总量', value: stats.value.order_count, icon: shallowRef(List), colorClass: 'icon-red' },
  // ★★★ 修复点：后端返回的是 trade_amount ★★★
  { label: '总交易额', value: '¥ ' + (stats.value.trade_amount || 0), icon: shallowRef(Wallet), colorClass: 'icon-purple' },
])

const fetchStats = async () => {
  try {
    const token = localStorage.getItem('admin_token')
    const res = await request.get('/api/admin/stats', {
      headers: { Authorization: token } // request.js 已处理前缀，这里直接传
    })
    // 兼容空数据
    stats.value = res || { user_count: 0, product_count: 0, order_count: 0, trade_amount: 0 }
  } catch (e) {
    console.error(e)
  }
}

onMounted(() => {
  fetchStats()
})
</script>

<style scoped lang="scss">
.dashboard-container {
  padding: 24px;
  background-color: #f6f8f9;
  min-height: calc(100vh - 60px);
}

/* 1. 欢迎横幅 */
.welcome-banner {
  background: linear-gradient(135deg, #ffffff 0%, #eef1f5 100%);
  border-radius: 16px;
  padding: 30px 40px;
  margin-bottom: 30px;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.03);
  display: flex; align-items: center; position: relative; overflow: hidden;

  .welcome-content {
    display: flex; justify-content: space-between; align-items: center; width: 100%; position: relative; z-index: 2;
  }
  .text-group {
    .greeting { margin: 0 0 10px 0; font-size: 28px; color: #333; font-weight: 800; }
    .subtitle { margin: 0; color: #666; font-size: 15px; }
  }
  .admin-avatar {
    width: 100px; height: 100px; filter: drop-shadow(0 4px 10px rgba(0,0,0,0.1));
    transition: transform 0.3s; &:hover { transform: rotate(10deg) scale(1.1); }
  }
}

/* 2. 数据卡片 */
.data-row { margin-bottom: 30px; }
.stat-card {
  border: none; border-radius: 16px; position: relative; overflow: hidden; transition: all 0.3s ease; cursor: default;
  &:hover { transform: translateY(-5px); box-shadow: 0 12px 30px rgba(0, 0, 0, 0.08); }

  .stat-content {
    display: flex; justify-content: space-between; align-items: center; padding: 10px; position: relative; z-index: 2;
  }
  .stat-info {
    .stat-label { font-size: 14px; color: #909399; margin-bottom: 8px; font-weight: 500; }
    .stat-num { font-size: 28px; font-weight: 900; color: #303133; letter-spacing: 1px; font-family: 'Arial', sans-serif; }
  }
  .stat-icon {
    width: 56px; height: 56px; border-radius: 14px; display: flex; align-items: center; justify-content: center; font-size: 28px;
    &.icon-blue { background: #ecf5ff; color: #409eff; }
    &.icon-orange { background: #fdf6ec; color: #e6a23c; }
    &.icon-red { background: #fef0f0; color: #f56c6c; }
    &.icon-purple { background: #f4f4f5; color: #6f42c1; }
  }
  /* 装饰背景圆 */
  .bg-circle {
    position: absolute; top: -20px; right: -20px; width: 100px; height: 100px; border-radius: 50%; opacity: 0.1; z-index: 1;
    &.icon-blue { background: #409eff; }
    &.icon-orange { background: #e6a23c; }
    &.icon-red { background: #f56c6c; }
    &.icon-purple { background: #6f42c1; }
  }
}

/* 3. 底部区域 */
.bottom-row {
  .el-card { border: none; border-radius: 16px; margin-bottom: 20px; }
  .card-header { font-weight: bold; font-size: 16px; display: flex; align-items: center; gap: 8px; }
}

.quick-actions {
  display: flex; gap: 40px; padding: 20px 10px;
  .action-item {
    display: flex; flex-direction: column; align-items: center; gap: 10px; cursor: pointer; transition: 0.3s;
    &:hover { transform: translateY(-3px); span { color: #409eff; } }
    .icon-box {
      width: 60px; height: 60px; border-radius: 20px; display: flex; align-items: center; justify-content: center; font-size: 28px; color: #fff; box-shadow: 0 8px 16px rgba(0,0,0,0.1);
      &.blue { background: linear-gradient(135deg, #409eff, #79bbff); }
      &.green { background: linear-gradient(135deg, #67c23a, #95d475); }
      &.orange { background: linear-gradient(135deg, #e6a23c, #f3d19e); }
      &.purple { background: linear-gradient(135deg, #909399, #c8c9cc); }
    }
    span { font-size: 14px; color: #606266; font-weight: 500; transition: 0.3s; }
  }
}

.notice-list {
  list-style: none; padding: 0; margin: 0;
  li {
    display: flex; align-items: center; gap: 10px; padding: 12px 0; border-bottom: 1px solid #f0f2f5; font-size: 14px; color: #606266;
    &:last-child { border-bottom: none; }
  }
}

/* 简单入场动画 */
.animate-up { animation: slideUp 0.6s ease-out backwards; }
.delay-1 { animation-delay: 0.1s; } .delay-2 { animation-delay: 0.2s; }
.delay-3 { animation-delay: 0.3s; } .delay-4 { animation-delay: 0.4s; } .delay-5 { animation-delay: 0.5s; }

@keyframes slideUp {
  from { opacity: 0; transform: translateY(20px); }
  to { opacity: 1; transform: translateY(0); }
}
</style>