<template>
  <div class="profile-page">
    <div class="profile-header-bg"></div>

    <div class="container">
      <div class="profile-nav animate-up">
        <div class="back-btn glass-effect" @click="$router.push('/')">
          <el-icon><ArrowLeft /></el-icon> 返回首页
        </div>
      </div>

      <div class="main-layout">
        <div class="left-sidebar animate-up">
          <div class="user-card">
            <div class="avatar-section">
              <el-avatar :size="100" :src="userInfo.avatar || defaultAvatar" class="user-avatar" />
              <div class="edit-avatar-btn" @click="triggerFileUpload">
                <el-icon><Camera /></el-icon>
              </div>
              <input type="file" ref="fileInput" accept="image/png, image/jpeg, image/jpg" class="hidden-file-input" @change="handleFileChange" />
            </div>

            <h2 class="nickname">{{ userInfo.nickname || userInfo.username }}</h2>
            <div class="username-badge">ID: {{ userInfo.username }}</div>

            <div class="info-list">
              <div class="info-item">
                <el-icon><Iphone /></el-icon>
                <span>{{ userInfo.phone || '暂未绑定手机' }}</span>
              </div>
              <div class="info-item">
                <el-icon><Location /></el-icon>
                <span>闲趣校园 · 实名认证</span>
              </div>
            </div>

            <div class="btn-group">
              <button class="btn-primary" @click="openEditProfileModal">编辑资料</button>
              <button class="btn-outline" @click="openPasswordModal">修改密码</button>
            </div>
          </div>

          <div class="stat-card">
            <div class="stat-item">
              <div class="num">{{ myProducts.length }}</div>
              <div class="label">发布</div>
            </div>
            <div class="divider"></div>
            <div class="stat-item">
              <div class="num">{{ myFavorites.length }}</div>
              <div class="label">收藏</div>
            </div>
          </div>
        </div>

        <div class="right-content animate-up delay-1">
          <div class="content-card">
            <el-tabs v-model="activeTab" class="custom-tabs">
              <el-tab-pane label="我发布的" name="products">
                <div class="grid-list" v-if="myProducts.length > 0">
                  <div v-for="item in myProducts" :key="item.id" class="mini-product-card" @click="goToManage(item.id)">
                    <div class="img-box">
                      <img :src="fixImageUrl(item.image)" loading="lazy" />
                      <div v-if="item.status !== 1" class="status-mask" :class="item.status === 2 ? 'sold' : 'off'">
                        <span>{{ item.status === 2 ? '已售出' : '已下架' }}</span>
                      </div>
                      <div class="manage-hint"><el-icon><Edit /></el-icon> 管理宝贝</div>
                    </div>
                    <div class="info">
                      <div class="title">{{ item.name || item.title }}</div>
                      <div class="price-row">
                        <span class="currency">¥</span><span class="amount">{{ item.price }}</span>
                      </div>
                    </div>
                  </div>
                </div>
                <el-empty v-else description="还没有发布过宝贝哦" :image-size="140" />
              </el-tab-pane>

              <el-tab-pane label="我的收藏" name="favorites">
                <div class="grid-list" v-if="myFavorites.length > 0">
                  <div v-for="item in myFavorites" :key="item.id" class="mini-product-card" @click="goToProduct(item.product?.id)">
                    <div class="img-box">
                      <img :src="fixImageUrl(item.product?.image)" loading="lazy" />
                      <div v-if="item.product?.status !== 1" class="status-mask" :class="item.product?.status === 2 ? 'sold' : 'off'">
                        <span>{{ item.product?.status === 2 ? '已售出' : '已下架' }}</span>
                      </div>
                    </div>
                    <div class="info">
                      <div class="title" :class="{'text-gray': item.product?.status !== 1}">
                        {{ item.product?.name || item.product?.title || '未知商品' }}
                      </div>
                      <div class="price-row">
                        <span class="currency">¥</span><span class="amount">{{ item.product?.price }}</span>
                      </div>
                      <div class="seller-mini">
                        <el-avatar :size="16" :src="defaultAvatar" class="mini-avatar" />
                        <span class="name">{{ getSellerName(item.product) }}</span>
                      </div>
                    </div>
                  </div>
                </div>
                <el-empty v-else description="还没有收藏宝贝" :image-size="140" />
              </el-tab-pane>
            </el-tabs>
          </div>
        </div>
      </div>
    </div>

    <el-dialog v-model="editProfileVisible" title="编辑资料" width="90%" maxWidth="480px" align-center class="custom-dialog" destroy-on-close>
      <div class="dialog-body">
        <div class="dialog-avatar-upload" @click="triggerFileUpload">
          <el-avatar :size="80" :src="previewAvatar || fixImageUrl(userInfo.avatar) || defaultAvatar" class="dialog-avatar" />
          <div class="upload-hint"><el-icon><Camera /></el-icon> 更换头像</div>
        </div>
        <el-form :model="profileForm" label-position="top" class="edit-form">
          <el-form-item label="昵称">
            <el-input v-model="profileForm.nickname" placeholder="请输入昵称" size="large" />
          </el-form-item>
          <el-form-item label="手机号">
            <el-input v-model="profileForm.phone" placeholder="请输入手机号" size="large" />
          </el-form-item>
        </el-form>
      </div>
      <template #footer>
        <div class="dialog-footer">
          <button class="btn-cancel" @click="editProfileVisible = false">取消</button>
          <button class="btn-submit" @click="submitProfileEdit" :disabled="isSubmitting">{{ isSubmitting ? '保存中...' : '保存' }}</button>
        </div>
      </template>
    </el-dialog>

    <el-dialog v-model="pwdVisible" title="修改密码" width="90%" maxWidth="400px" align-center class="custom-dialog" destroy-on-close>
      <el-form :model="pwdForm" label-width="0" class="pwd-form">
        <el-form-item><el-input v-model="pwdForm.old_password" type="password" placeholder="当前密码" show-password size="large" /></el-form-item>
        <el-form-item><el-input v-model="pwdForm.new_password" type="password" placeholder="新密码 (至少6位)" show-password size="large" /></el-form-item>
      </el-form>
      <template #footer>
        <div class="dialog-footer">
          <button class="btn-cancel" @click="pwdVisible = false">取消</button>
          <button class="btn-submit" @click="submitPwd">确认修改</button>
        </div>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import request from '@/utils/request'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { Camera, Iphone, Location, Edit, ArrowLeft } from '@element-plus/icons-vue'

const router = useRouter()
const userInfo = ref({})
const myProducts = ref([])
const myFavorites = ref([])
const activeTab = ref('products')
const defaultAvatar = 'https://cube.elemecdn.com/3/7c/3ea6beec64369c2642b92c6726f1epng.png'
const isSubmitting = ref(false)

const editProfileVisible = ref(false)
const pwdVisible = ref(false)
const profileForm = reactive({ nickname: '', avatar: '', phone: '' })
const pwdForm = reactive({ old_password: '', new_password: '' })

// 本地上传相关
const fileInput = ref(null)
const previewAvatar = ref('')
const selectedFile = ref(null)

// 辅助函数：修复图片路径
const fixImageUrl = (url) => {
  if (!url) return ''
  let fixedUrl = url
  if (!fixedUrl.startsWith('http')) {
    fixedUrl = 'http://127.0.0.1:8081' + fixedUrl
  }
  return fixedUrl.replace('localhost', '127.0.0.1')
}

// 获取卖家名字 (收藏列表中使用)
const getSellerName = (product) => {
  if (!product) return '未知'
  const s = product.seller || product.user || {}
  return s.nickname || s.username || '闲趣用户'
}

const triggerFileUpload = () => { fileInput.value.click() }

const handleFileChange = (event) => {
  const file = event.target.files[0]
  if (!file) return
  if (!['image/jpeg', 'image/png', 'image/jpg'].includes(file.type)) { return ElMessage.error('仅支持 JPG/PNG 格式') }
  if (file.size > 5 * 1024 * 1024) { return ElMessage.error('图片不能超过 5MB') }

  selectedFile.value = file
  const reader = new FileReader()
  reader.onload = (e) => { previewAvatar.value = e.target.result }
  reader.readAsDataURL(file)

  if (!editProfileVisible.value) openEditProfileModal()
}

// 核心数据获取
const fetchUserData = async () => {
  try {
    const localUser = JSON.parse(localStorage.getItem('user') || '{}')
    if (!localUser.id) { router.push('/'); return }
    userInfo.value = localUser

    // 获取发布的商品
    const resProducts = await request.get('/api/products', { params: { user_id: localUser.id } })
    myProducts.value = (resProducts.list || []).filter(item => Number(item.user_id) === Number(localUser.id))

    // 获取收藏列表
    const resFav = await request.get('/api/user/data?type=favorites')
    myFavorites.value = (resFav.data || []).filter(item => item.product && item.product.id)

  } catch (e) { console.error(e) }
}

const goToManage = (id) => { router.push(`/publish?id=${id}`) }
const goToProduct = (id) => { if (id) router.push(`/product/${id}`) }

const openEditProfileModal = () => {
  Object.assign(profileForm, { nickname: userInfo.value.nickname, phone: userInfo.value.phone })
  if (!selectedFile.value) { previewAvatar.value = fixImageUrl(userInfo.value.avatar) }
  editProfileVisible.value = true
}

const submitProfileEdit = async () => {
  if (isSubmitting.value) return
  isSubmitting.value = true
  try {
    let newAvatarUrl = userInfo.value.avatar
    if (selectedFile.value) {
      const formData = new FormData()
      formData.append('file', selectedFile.value)
      const uploadRes = await request.post('/api/upload', formData, { headers: { 'Content-Type': 'multipart/form-data' } })

      let uploadedUrl = uploadRes.url
      if (uploadedUrl && !uploadedUrl.startsWith('http')) { uploadedUrl = 'http://127.0.0.1:8081' + uploadedUrl }
      uploadedUrl = uploadedUrl.replace('localhost', '127.0.0.1')
      newAvatarUrl = uploadedUrl
    }

    const updateData = { ...profileForm, avatar: newAvatarUrl }
    await request.put('/api/user/profile', updateData)

    ElMessage.success('资料更新成功')
    const newUser = { ...userInfo.value, ...updateData }
    localStorage.setItem('user', JSON.stringify(newUser))
    userInfo.value = newUser
    selectedFile.value = null
    editProfileVisible.value = false

  } catch (e) { ElMessage.error(e.response?.data?.error || '更新失败') }
  finally { isSubmitting.value = false }
}

const openPasswordModal = () => { pwdForm.old_password = ''; pwdForm.new_password = ''; pwdVisible.value = true }
const submitPwd = async () => { if (!pwdForm.old_password || !pwdForm.new_password) return ElMessage.warning('请填写完整'); try { await request.put('/api/user/password', pwdForm); ElMessage.success('密码修改成功，请重新登录'); localStorage.clear(); window.location.href = '/' } catch (e) { ElMessage.error(e.response?.data?.error || '修改失败') } }

onMounted(fetchUserData)
</script>

<style scoped lang="scss">
$primary: #ffdf5d;
$text-main: #333;
$text-light: #999;
$bg-page: #f6f7f9;

/* 全局布局 */
.profile-page { min-height: 100vh; background: $bg-page; padding-bottom: 60px; }
.profile-header-bg { height: 200px; background: linear-gradient(180deg, $primary 0%, rgba(255, 223, 93, 0.05) 100%); margin-bottom: -100px; }
.container { max-width: 1100px; margin: 0 auto; padding: 0 20px; box-sizing: border-box; }
.hidden-file-input { display: none; }

/* ★★★ 新增：返回按钮区域 ★★★ */
.profile-nav {
  margin-bottom: 24px;
  position: relative; z-index: 10;

  .back-btn {
    display: inline-flex; align-items: center; gap: 6px;
    padding: 10px 20px; border-radius: 99px;
    background: rgba(255, 255, 255, 0.85); backdrop-filter: blur(10px);
    font-weight: bold; color: #333; cursor: pointer;
    box-shadow: 0 4px 12px rgba(0,0,0,0.05); transition: 0.3s;
    border: 1px solid rgba(255,255,255,0.6);

    &:hover { transform: translateX(-4px); background: #fff; box-shadow: 0 6px 16px rgba(0,0,0,0.1); }
  }
}

.main-layout { display: grid; grid-template-columns: 340px 1fr; gap: 24px; align-items: start; }

/* 左侧：用户信息卡片 */
.user-card {
  background: #fff; border-radius: 24px; padding: 36px 30px; text-align: center; box-shadow: 0 8px 30px rgba(0,0,0,0.04);

  .avatar-section {
    position: relative; display: inline-block; margin-bottom: 20px; cursor: pointer;
    .user-avatar { border: 4px solid #fff; box-shadow: 0 4px 16px rgba(0,0,0,0.08); transition: transform 0.3s cubic-bezier(0.34, 1.56, 0.64, 1); }
    .edit-avatar-btn {
      position: absolute; bottom: 0; right: 0; width: 40px; height: 40px; background: #fff; color: #333; border-radius: 50%; display: flex; align-items: center; justify-content: center; cursor: pointer; transition: 0.2s; border: 1px solid #eee; box-shadow: 0 4px 10px rgba(0,0,0,0.1); font-size: 20px;
      &:hover { background: $primary; border-color: $primary; color: #000; transform: scale(1.1); }
    }
    &:hover .user-avatar { transform: scale(1.05) rotate(-2deg); }
  }

  .nickname { margin: 0 0 8px 0; font-size: 26px; font-weight: 800; color: $text-main; }
  .username-badge { color: $text-light; font-size: 13px; margin-bottom: 30px; background: #f7f7f7; padding: 4px 12px; border-radius: 99px; display: inline-block; font-weight: 600; }

  .info-list {
    text-align: left; padding: 0 10px; margin-bottom: 30px;
    .info-item { display: flex; align-items: center; gap: 14px; margin-bottom: 16px; color: #555; font-size: 15px; font-weight: 500; .el-icon { font-size: 20px; color: #ddd; } }
  }

  .btn-group {
    display: flex; flex-direction: column; gap: 14px;
    button { width: 100%; height: 50px; border-radius: 16px; font-weight: 700; cursor: pointer; transition: 0.2s; font-size: 16px; display: flex; align-items: center; justify-content: center; }
    .btn-primary { background: #1a1a1a; border: none; color: $primary; &:hover { background: #000; transform: translateY(-2px); box-shadow: 0 6px 16px rgba(0,0,0,0.15); } }
    .btn-outline { background: #fff; border: 2px solid #eee; color: #666; &:hover { border-color: #333; color: #333; background: #fafafa; } }
  }
}

/* 左侧：统计卡片 */
.stat-card {
  background: #fff; border-radius: 20px; padding: 24px; display: flex; justify-content: space-around; align-items: center; box-shadow: 0 8px 20px rgba(0,0,0,0.03); margin-top: 20px;
  .stat-item { text-align: center; cursor: pointer; transition: 0.2s; &:hover { transform: translateY(-3px); } }
  .num { font-size: 24px; font-weight: 900; color: $text-main; font-family: 'Arial Black', sans-serif; }
  .label { font-size: 13px; color: $text-light; margin-top: 6px; font-weight: 600; }
  .divider { width: 1px; height: 36px; background: #eee; }
}

/* 右侧：内容区域 */
.content-card { background: #fff; border-radius: 24px; padding: 20px 30px; min-height: 500px; box-shadow: 0 8px 30px rgba(0,0,0,0.03); }

/* Tab 美化 */
:deep(.custom-tabs .el-tabs__item) { font-size: 16px; height: 56px; color: #999; transition: 0.3s; &.is-active { font-weight: 800; color: #333; font-size: 18px; } &:hover { color: #666; } }
:deep(.custom-tabs .el-tabs__active-bar) { background-color: $primary; height: 4px; border-radius: 4px; bottom: 6px; width: 40px !important; margin-left: 20px; } /* 短横线风格 */
:deep(.custom-tabs .el-tabs__nav-wrap::after) { height: 1px; background-color: #f2f2f2; }

/* 迷你商品列表 */
.grid-list {
  display: grid; grid-template-columns: repeat(auto-fill, minmax(200px, 1fr)); gap: 20px; padding-top: 24px;
}

.mini-product-card {
  border-radius: 16px; overflow: hidden; background: #fff; border: 1px solid #f9f9f9; cursor: pointer; transition: all 0.3s cubic-bezier(0.25, 0.8, 0.25, 1); position: relative;
  &:hover { transform: translateY(-6px); box-shadow: 0 12px 30px rgba(0,0,0,0.06); .manage-hint { opacity: 1; transform: translateY(0); } }

  .img-box {
    width: 100%; aspect-ratio: 1; background: #f9f9f9; position: relative; overflow: hidden;
    img { width: 100%; height: 100%; object-fit: cover; transition: transform 0.5s; }
    &:hover img { transform: scale(1.05); }

    .status-mask {
      position: absolute; inset: 0; display: flex; align-items: center; justify-content: center; backdrop-filter: grayscale(100%) blur(2px);
      span { border: 3px solid #fff; color: #fff; padding: 6px 16px; font-size: 14px; font-weight: 900; transform: rotate(-10deg); border-radius: 8px; letter-spacing: 2px; }
      &.sold { background: rgba(0,0,0,0.6); }
      &.off { background: rgba(100,100,100,0.7); }
    }

    .manage-hint {
      position: absolute; bottom: 0; left: 0; right: 0; background: rgba(255, 223, 93, 0.96);
      color: #1a1a1a; font-size: 13px; padding: 10px; text-align: center; font-weight: bold;
      opacity: 0; transform: translateY(100%); transition: all 0.3s;
      display: flex; align-items: center; justify-content: center; gap: 6px;
    }
  }

  .info {
    padding: 14px;
    .title { font-size: 15px; color: $text-main; margin-bottom: 8px; white-space: nowrap; overflow: hidden; text-overflow: ellipsis; font-weight: 500; &.text-gray { color: #aaa; text-decoration: line-through; } }
    .price-row { display: flex; align-items: baseline; gap: 2px; .currency { font-size: 12px; font-weight: bold; color: #ff5000; } .amount { font-size: 18px; font-weight: 800; color: #ff5000; } }
    .seller-mini { display: flex; align-items: center; gap: 6px; margin-top: 8px; font-size: 12px; color: #bbb; .mini-avatar { border: 1px solid #f0f0f0; } }
  }
}

/* 弹窗通用 */
.custom-dialog { border-radius: 24px; overflow: hidden; .el-dialog__header { margin-right: 0; padding: 24px; font-weight: 800; } }
.dialog-body { text-align: center; }
.dialog-avatar-upload {
  display: inline-block; position: relative; margin-bottom: 24px; cursor: pointer;
  .dialog-avatar { border: 4px solid #f5f5f5; }
  .upload-hint { position: absolute; bottom: -6px; left: 50%; transform: translateX(-50%); background: #333; color: #ffdf5d; font-size: 12px; padding: 4px 12px; border-radius: 99px; white-space: nowrap; display: flex; align-items: center; gap: 4px; box-shadow: 0 4px 10px rgba(0,0,0,0.1); font-weight: bold; }
}
.dialog-footer {
  display: flex; justify-content: flex-end; gap: 12px;
  button { height: 44px; padding: 0 28px; border-radius: 12px; font-weight: 700; cursor: pointer; transition: 0.2s; font-size: 15px; }
  .btn-cancel { background: #f5f5f5; border: 1px solid #eee; color: #666; &:hover { background: #eee; color: #333; } }
  .btn-submit { background: $primary; border: none; color: #1a1a1a; &:hover { background: #f0cf4b; transform: translateY(-2px); box-shadow: 0 4px 12px rgba(255, 223, 93, 0.3); } &:disabled { opacity: 0.6; cursor: not-allowed; transform: none; } }
}

/* 动画 */
.animate-up { animation: fadeInUp 0.6s cubic-bezier(0.25, 0.8, 0.25, 1) backwards; }
.delay-1 { animation-delay: 0.15s; }
@keyframes fadeInUp { from { opacity: 0; transform: translateY(30px); } to { opacity: 1; transform: translateY(0); } }

@media (max-width: 768px) {
  .main-layout { grid-template-columns: 1fr; }
  .profile-header-bg { height: 160px; margin-bottom: -80px; }
  .user-card { padding: 24px 20px; }
  .grid-list { grid-template-columns: repeat(2, 1fr); gap: 12px; }
}
</style>