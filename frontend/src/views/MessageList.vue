<template>
  <div class="message-list-page" v-loading="loading">
    <div class="nav-header">
      <div class="nav-content container">
        <div class="back-btn" @click="$router.push('/')">
          <el-icon><ArrowLeft /></el-icon> 首页
        </div>
        <div class="page-title">消息中心</div>
        <div class="placeholder"></div>
      </div>
    </div>

    <div class="container main-content">
      <div class="chat-list-card animate-up">
        <div v-if="contacts.length === 0 && !loading" class="empty-state">
          <el-empty description="暂无消息记录，去看看宝贝吧~" :image-size="140" />
          <button class="go-home-btn" @click="$router.push('/')">去逛逛</button>
        </div>

        <div
            v-for="item in contacts"
            :key="item.id"
            class="contact-item"
            @click="toChat(item.id)"
        >
          <div class="avatar-box">
            <el-avatar :size="54" :src="item.avatar || defaultAvatar" class="avatar-img" />
          </div>

          <div class="info-box">
            <div class="row-top">
              <span class="nickname">{{ item.nickname || item.username || '闲趣用户' }}</span>
              <span class="time">{{ formatTime(item.time) }}</span>
            </div>
            <div class="row-bottom">
              <p class="last-msg">{{ item.last_msg || '暂无消息内容...' }}</p>
            </div>
          </div>

          <div class="action-icon">
            <el-icon><ArrowRight /></el-icon>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import request from '@/utils/request'
import { ArrowLeft, ArrowRight } from '@element-plus/icons-vue'

const router = useRouter()
const loading = ref(false)
const contacts = ref([])
const defaultAvatar = 'https://cube.elemecdn.com/3/7c/3ea6beec64369c2642b92c6726f1epng.png'

// 安全获取当前用户
const getCurrentUser = () => {
  try {
    const u = localStorage.getItem('user')
    return u ? JSON.parse(u) : null
  } catch (e) {
    return null
  }
}

const user = ref(getCurrentUser())

const fetchContacts = async () => {
  if (!user.value || !user.value.id) {
    router.push('/')
    return
  }

  loading.value = true
  try {
    const res = await request.get('/api/chat/contacts')
    contacts.value = res.data || []
  } catch (e) {
    console.error('获取消息列表失败:', e)
  } finally {
    loading.value = false
  }
}

const toChat = (targetId) => {
  router.push(`/chat/${targetId}`)
}

const formatTime = (timeStr) => {
  if (!timeStr) return ''
  const date = new Date(timeStr)
  const now = new Date()

  if (date.toDateString() === now.toDateString()) {
    // 当天显示时分
    return date.getHours().toString().padStart(2, '0') + ':' + date.getMinutes().toString().padStart(2, '0')
  } else if (date.getFullYear() === now.getFullYear()) {
    // 今年显示月日
    return (date.getMonth() + 1) + '-' + date.getDate()
  }
  // 跨年显示年月日
  return date.getFullYear() + '-' + (date.getMonth() + 1) + '-' + date.getDate()
}

onMounted(() => {
  fetchContacts()
})
</script>

<style scoped lang="scss">
$primary: #ffdf5d;
$bg-color: #f6f7f9;

.message-list-page {
  min-height: 100vh;
  background: $bg-color;
  padding-bottom: 40px;
}

/* 顶部导航 */
.nav-header {
  height: 60px;
  background: rgba(255, 255, 255, 0.9);
  backdrop-filter: blur(10px);
  position: sticky;
  top: 0;
  z-index: 100;
  box-shadow: 0 2px 10px rgba(0,0,0,0.03);

  .nav-content {
    height: 100%;
    display: flex;
    align-items: center;
    justify-content: space-between;
  }

  .back-btn {
    cursor: pointer;
    font-size: 15px;
    font-weight: 500;
    color: #333;
    display: flex;
    align-items: center;
    gap: 4px;
    padding: 6px 12px;
    border-radius: 20px;
    transition: 0.2s;
    &:hover { background: #f0f0f0; }
  }

  .page-title {
    font-size: 18px;
    font-weight: 800;
    color: #333;
  }

  .placeholder { width: 60px; }
}

.container {
  max-width: 800px;
  margin: 0 auto;
  padding: 0 20px;
  box-sizing: border-box;
}

.main-content {
  margin-top: 24px;
}

/* 列表卡片容器 */
.chat-list-card {
  background: #fff;
  border-radius: 20px;
  box-shadow: 0 8px 30px rgba(0,0,0,0.04);
  overflow: hidden;
  min-height: 400px;
}

/* 单个联系人项 */
.contact-item {
  display: flex;
  align-items: center;
  padding: 20px;
  border-bottom: 1px solid #f7f7f7;
  cursor: pointer;
  transition: all 0.25s ease;
  position: relative;

  &:last-child { border-bottom: none; }

  &:hover {
    background: #fffcf0; /* 悬停时淡淡的黄色背景 */
    padding-left: 26px; /* 微妙的位移 */
    .action-icon { opacity: 1; transform: translateX(0); }
  }

  /* 头像区 */
  .avatar-box {
    margin-right: 18px;
    position: relative;
    .avatar-img { border: 1px solid #f0f0f0; }
    .red-dot {
      position: absolute; top: 0; right: 0;
      width: 10px; height: 10px; background: #ff4d4f;
      border-radius: 50%; border: 2px solid #fff;
    }
  }

  /* 内容区 */
  .info-box {
    flex: 1;
    overflow: hidden;
    display: flex;
    flex-direction: column;
    gap: 4px;

    .row-top {
      display: flex;
      justify-content: space-between;
      align-items: center;

      .nickname {
        font-size: 16px;
        font-weight: 700;
        color: #333;
      }
      .time {
        font-size: 12px;
        color: #999;
        font-weight: 500;
      }
    }

    .row-bottom {
      .last-msg {
        font-size: 14px;
        color: #666;
        white-space: nowrap;
        overflow: hidden;
        text-overflow: ellipsis;
        margin: 0;
      }
    }
  }

  /* 右侧箭头 (默认隐藏或淡化，hover时显现) */
  .action-icon {
    margin-left: 16px;
    color: #ccc;
    font-size: 18px;
    opacity: 0.5;
    transform: translateX(-5px);
    transition: all 0.3s;
  }
}

/* 空状态美化 */
.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 80px 0;

  .go-home-btn {
    margin-top: 20px;
    padding: 10px 28px;
    background: #333;
    color: $primary;
    border: none;
    border-radius: 99px;
    font-weight: 700;
    cursor: pointer;
    transition: 0.2s;
    &:hover { transform: translateY(-3px); box-shadow: 0 4px 12px rgba(0,0,0,0.15); }
  }
}

/* 简单的入场动画 */
.animate-up {
  animation: slideUp 0.5s ease-out;
}

@keyframes slideUp {
  from { opacity: 0; transform: translateY(20px); }
  to { opacity: 1; transform: translateY(0); }
}

@media (max-width: 600px) {
  .container { padding: 0 12px; }
  .chat-list-card { border-radius: 12px; }
  .contact-item { padding: 16px; }
}
</style>