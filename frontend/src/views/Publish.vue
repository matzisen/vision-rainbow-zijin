<template>
  <div class="publish-page">
    <div class="simple-nav">
      <div class="brand" @click="$router.push('/')">XIANQU</div>
      <div class="close-btn" @click="$router.go(-1)">
        <el-icon><Close /></el-icon> {{ isEditMode ? '取消编辑' : '取消发布' }}
      </div>
    </div>

    <div class="publish-container animate-up">
      <div class="glass-card">
        <h2 class="page-title">{{ isEditMode ? '编辑宝贝信息' : '发布闲置宝贝' }}</h2>

        <div class="publish-layout" v-loading="loading">
          <div class="upload-area">
            <el-upload
                class="image-uploader"
                action="http://127.0.0.1:8081/api/upload"
                name="file"
                :headers="uploadHeaders"
                :show-file-list="false"
                :on-success="handleUploadSuccess"
                :before-upload="beforeUpload"
                :on-error="handleUploadError"
            >
              <div class="preview-viewport" :class="{ 'has-image': publishForm.image }">
                <img v-if="publishForm.image" :src="publishForm.image" class="viewport-img" />

                <div v-else class="upload-placeholder">
                  <div class="icon-circle"><el-icon :size="40"><Plus /></el-icon></div>
                  <span class="text">上传封面图</span>
                  <span class="sub-text">支持 JPG/PNG (最大10MB)</span>
                </div>

                <div class="hover-overlay" v-if="publishForm.image">
                  <el-icon><Switch /></el-icon> 更换图片
                </div>
              </div>
            </el-upload>
          </div>

          <div class="form-area">
            <div class="form-group">
              <input v-model="publishForm.name" class="input-ghost title-input" placeholder="宝贝标题 品牌型号..." maxlength="30" />
            </div>

            <div class="form-group textarea-group">
              <textarea v-model="publishForm.description" class="input-ghost desc-input" rows="6" placeholder="描述一下宝贝的入手渠道、新旧程度和转手原因..." ></textarea>
            </div>

            <div class="form-detail-row">
              <div class="input-item price-item">
                <div class="label">价格</div>
                <div class="price-input-wrap">
                  <span class="prefix">¥</span>
                  <input v-model.number="publishForm.price" type="number" class="price-input" placeholder="0.00" />
                </div>
              </div>

              <div class="input-item">
                <div class="label">数量</div>
                <el-input-number
                    v-model="publishForm.count"
                    :min="1" :max="99"
                    controls-position="right"
                    class="minimal-number-input"
                />
              </div>

              <div class="input-item category-item">
                <div class="label">分类</div>
                <el-select v-model="publishForm.category" placeholder="选择分类" class="custom-select" :popper-append-to-body="false">
                  <el-option v-for="cat in categories" :key="cat.id" :label="cat.name" :value="cat.name" />
                </el-select>
              </div>
            </div>

            <div class="status-control" v-if="isEditMode">
              <div class="label">当前状态</div>
              <el-radio-group v-model="publishForm.status" size="large">
                <el-radio-button :label="1">上架中</el-radio-button>
                <el-radio-button :label="3">已下架</el-radio-button>
                <el-radio-button :label="2" disabled v-if="publishForm.status === 2">已售出</el-radio-button>
              </el-radio-group>
            </div>

            <div class="form-footer-bar">
              <div class="tags-group">
                <div class="tag-chip" :class="{ active: publishForm.is_free_shipping }" @click="publishForm.is_free_shipping = !publishForm.is_free_shipping">
                  <el-icon><Van /></el-icon> 包邮
                </div>
                <div class="tag-chip" :class="{ active: publishForm.is_negotiable }" @click="publishForm.is_negotiable = !publishForm.is_negotiable">
                  <el-icon><Scissor /></el-icon> 可小刀
                </div>
              </div>

              <button class="btn-submit" @click="submitPublish" :disabled="isSubmitting">
                {{ isSubmitting ? '提交中...' : (isEditMode ? '保存修改' : '确认发布') }}
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import request from '@/utils/request'
import { ElMessage } from 'element-plus'
import { Plus, Switch, Van, Scissor, Close } from '@element-plus/icons-vue'

const router = useRouter()
const route = useRoute()
const isSubmitting = ref(false)
const loading = ref(false)
const user = JSON.parse(localStorage.getItem('user') || '{}')

// 判断是否为编辑模式
const isEditMode = computed(() => !!route.query.id)

const uploadHeaders = computed(() => {
  const token = localStorage.getItem('token')
  return { Authorization: token || '' }
})

const categories = ref([
  { id: 1, name: '数码' }, { id: 2, name: '书籍' }, { id: 3, name: '生活' },
  { id: 4, name: '服饰' }, { id: 5, name: '运动' }, { id: 6, name: '美妆' },
  { id: 7, name: '乐器' }, { id: 8, name: '手游' }, { id: 9, name: '兼职' },
  { id: 10, name: '其他' }
])

const publishForm = reactive({
  name: '',
  description: '',
  price: '',
  count: 1,
  image: '',
  category: '',
  user_id: user.id,
  status: 1, // 默认上架
  is_free_shipping: false,
  is_negotiable: false
})

// 加载待编辑商品数据
const fetchProductData = async () => {
  if (!isEditMode.value) return

  loading.value = true
  try {
    const res = await request.get(`/api/products/${route.query.id}`)
    const data = res.data

    // 权限校验
    if (data.user_id !== user.id) {
      ElMessage.error('无权编辑此商品')
      router.push('/')
      return
    }

    // 回显数据
    Object.assign(publishForm, {
      name: data.name,
      description: data.description,
      price: data.price,
      count: data.count || 1,
      category: data.category,
      image: data.image,
      status: data.status, // 回显状态
      is_free_shipping: data.is_free_shipping,
      is_negotiable: data.is_negotiable
    })

    // 修复图片路径显示
    if (publishForm.image && !publishForm.image.startsWith('http')) {
      publishForm.image = 'http://127.0.0.1:8081' + publishForm.image
    }
    publishForm.image = publishForm.image.replace('localhost', '127.0.0.1')

  } catch (e) {
    console.error(e)
    ElMessage.error('获取商品信息失败')
  } finally {
    loading.value = false
  }
}

const handleUploadSuccess = (res) => {
  if (res.url) {
    let finalUrl = res.url
    if (!finalUrl.startsWith('http')) {
      finalUrl = 'http://127.0.0.1:8081' + finalUrl
    }
    finalUrl = finalUrl.replace('localhost', '127.0.0.1')
    publishForm.image = finalUrl
    ElMessage.success('图片上传成功')
  }
}

const handleUploadError = (err) => {
  console.error('Upload failed:', err)
  ElMessage.error('上传失败，请检查网络')
}

const beforeUpload = (file) => {
  const isLt10M = file.size / 1024 / 1024 < 10
  if (!isLt10M) ElMessage.warning('图片大小不能超过 10MB')
  return isLt10M
}

const submitPublish = async () => {
  if (!publishForm.name) return ElMessage.warning('请输入标题')
  if (!publishForm.price) return ElMessage.warning('请输入价格')
  if (!publishForm.image) return ElMessage.warning('请上传图片')
  if (!publishForm.category) return ElMessage.warning('请选择分类')

  isSubmitting.value = true
  try {
    const submitData = {
      name: publishForm.name,
      description: publishForm.description,
      price: Number(publishForm.price),
      count: Number(publishForm.count),
      image: publishForm.image,
      category: publishForm.category,
      status: publishForm.status, // 提交状态
      is_free_shipping: publishForm.is_free_shipping,
      is_negotiable: publishForm.is_negotiable
    }

    if (isEditMode.value) {
      await request.put(`/api/products/${route.query.id}`, submitData)
      ElMessage.success('修改成功！')
    } else {
      await request.post('/api/products', submitData)
      ElMessage.success('发布成功！')
    }

    setTimeout(() => {
      if (isEditMode.value) {
        router.push('/profile')
      } else {
        router.push('/')
      }
    }, 1000)

  } catch (e) {
    console.error(e)
    ElMessage.error(e.response?.data?.error || (isEditMode.value ? '修改失败' : '发布失败'))
  } finally {
    isSubmitting.value = false
  }
}

onMounted(() => {
  fetchProductData()
})
</script>

<style scoped lang="scss">
$primary: #ffdf5d;
$bg-color: #f2f3f5;

.publish-page {
  min-height: 100vh;
  background: $bg-color;
  background-image: radial-gradient(circle at 10% 20%, rgba(255, 223, 93, 0.2) 0%, transparent 40%);
  display: flex;
  flex-direction: column;
  height: 100vh;
  overflow: hidden;
}

.simple-nav {
  height: 60px;
  flex-shrink: 0;
  display: flex; justify-content: space-between; align-items: center;
  padding: 0 40px;
  .brand { font-weight: 900; font-size: 20px; cursor: pointer; }
  .close-btn {
    display: flex; align-items: center; gap: 6px; cursor: pointer; color: #666; font-weight: bold; transition: 0.2s;
    &:hover { color: #000; }
  }
}

.publish-container {
  flex: 1;
  display: flex;
  justify-content: center;
  align-items: center;
  padding: 20px;
  overflow: hidden;
}

/* Acrylic Glass Card */
.glass-card {
  width: 900px;
  max-height: calc(100vh - 100px);
  background: rgba(255, 255, 255, 0.85);
  backdrop-filter: blur(20px);
  border-radius: 32px;
  box-shadow: 0 20px 60px rgba(0,0,0,0.08), inset 0 0 0 1px rgba(255,255,255,0.8);
  padding: 40px;
  display: flex;
  flex-direction: column;

  .page-title { margin: 0 0 30px 0; font-size: 24px; font-weight: 800; color: #333; text-align: center; flex-shrink: 0; }
}

.publish-layout {
  display: flex;
  gap: 40px;
  overflow-y: auto;
  padding-right: 10px;

  &::-webkit-scrollbar { width: 6px; }
  &::-webkit-scrollbar-thumb { background: rgba(0,0,0,0.1); border-radius: 3px; }

  .upload-area {
    flex: 0 0 350px;
    min-height: 350px;
    background: #f9f9f9;
    border-radius: 24px;
    display: flex; align-items: center; justify-content: center;
    border: 2px dashed #eee;
    transition: 0.3s;

    &:hover { border-color: $primary; background: #fffcf0; }

    .image-uploader { width: 100%; height: 100%; display: flex; align-items: center; justify-content: center; }

    .preview-viewport {
      width: 100%; height: 350px;
      display: flex; flex-direction: column; align-items: center; justify-content: center;
      position: relative;

      .viewport-img { width: 100%; height: 100%; object-fit: contain; border-radius: 20px; }

      .upload-placeholder { text-align: center; .icon-circle { width: 70px; height: 70px; background: #fff; border-radius: 50%; display: flex; align-items: center; justify-content: center; margin: 0 auto 16px; box-shadow: 0 4px 12px rgba(0,0,0,0.05); } .text { font-weight: bold; color: #666; display: block; } .sub-text { font-size: 12px; color: #ccc; margin-top: 4px; } }

      .hover-overlay { position: absolute; inset: 0; background: rgba(0,0,0,0.4); border-radius: 20px; display: flex; align-items: center; justify-content: center; color: #fff; font-weight: bold; gap: 8px; opacity: 0; transition: 0.2s; cursor: pointer; }
      &:hover .hover-overlay { opacity: 1; }
    }
  }

  .form-area {
    flex: 1;
    display: flex; flex-direction: column;

    .input-ghost { width: 100%; background: transparent; border: none; border-bottom: 2px solid #f0f0f0; padding: 12px 0; font-size: 16px; outline: none; transition: 0.3s; &:focus { border-color: $primary; } }
    .title-input { font-size: 24px; font-weight: 800; margin-bottom: 20px; }
    .desc-input { resize: none; height: 120px; margin-bottom: 30px; font-size: 15px; line-height: 1.6; }

    .form-detail-row {
      display: grid; grid-template-columns: 1.5fr 1fr 1.2fr; gap: 20px; align-items: end; margin-bottom: 30px;
      .label { font-size: 12px; color: #999; margin-bottom: 6px; font-weight: bold; }
      .price-input-wrap { display: flex; align-items: baseline; border-bottom: 2px solid #333; padding-bottom: 4px; .prefix { font-size: 24px; font-weight: bold; color: #ff5000; margin-right: 4px; } .price-input { width: 100%; border: none; background: transparent; font-size: 32px; font-weight: 900; color: #ff5000; outline: none; } }
    }

    /* 状态管理样式 */
    .status-control {
      margin-bottom: 30px;
      .label { font-size: 12px; color: #999; margin-bottom: 8px; font-weight: bold; }
      /* 修改 radio button 选中颜色 */
      :deep(.el-radio-button__original-radio:checked + .el-radio-button__inner) {
        background-color: $primary;
        border-color: $primary;
        color: #000;
        box-shadow: none;
      }
    }

    .form-footer-bar {
      margin-top: auto; display: flex; justify-content: space-between; align-items: center;
      .tags-group { display: flex; gap: 12px; .tag-chip { padding: 8px 18px; background: #fff; border: 1px solid #eee; border-radius: 99px; cursor: pointer; font-size: 13px; font-weight: bold; color: #666; transition: 0.2s; &:hover { border-color: #ccc; } &.active { background: #ffdf5d; color: #000; border-color: #ffdf5d; } } }
      .btn-submit { padding: 12px 40px; background: #1a1a1a; color: #ffdf5d; border: none; border-radius: 99px; font-weight: bold; font-size: 16px; cursor: pointer; transition: 0.2s; box-shadow: 0 4px 16px rgba(0,0,0,0.15); &:hover { transform: translateY(-2px); box-shadow: 0 8px 24px rgba(0,0,0,0.2); } &:disabled { opacity: 0.6; cursor: not-allowed; } }
    }
  }
}

.minimal-number-input { width: 100%; :deep(.el-input__wrapper) { box-shadow: none !important; background: #f5f5f5; border-radius: 8px; } :deep(.el-input__inner) { font-weight: bold; background: transparent; text-align: center; } }
:deep(.custom-select) { width: 100%; .el-input__wrapper { box-shadow: none !important; background: #f5f5f5; border-radius: 8px; padding: 0 12px; } .el-input__inner { font-weight: bold; height: 40px; } }

.animate-up { animation: fadeInUp 0.6s cubic-bezier(0.2, 0.8, 0.2, 1); }
@keyframes fadeInUp { from { opacity: 0; transform: translateY(40px); } to { opacity: 1; transform: translateY(0); } }
</style>