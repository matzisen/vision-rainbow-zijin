<template>
  <el-dialog
      v-model="visible"
      width="800px"
      :show-close="false"
      class="auth-dialog"
      destroy-on-close
      align-center
      :close-on-click-modal="false"
  >
    <div class="auth-container">
      <div class="auth-left">
        <div class="brand">
          <el-icon :size="32"><Goods /></el-icon>
          <span>闲趣</span>
        </div>
        <div class="slogan">
          <h3>发现闲置的价值</h3>
          <p>更温暖的校园交易社区</p>
        </div>
        <div class="circle c1"></div>
        <div class="circle c2"></div>
      </div>

      <div class="auth-right">
        <button type="button" class="close-btn" @click="close">
          <el-icon><Close /></el-icon>
        </button>

        <h2>{{ isLogin ? '欢迎回来' : '创建账号' }}</h2>
        <p class="subtitle">{{ isLogin ? '登录后开始交易' : '注册只需 1 分钟' }}</p>

        <el-form :model="form" class="auth-form" @submit.prevent>

          <div class="input-group">
            <el-input
                v-model="form.username"
                placeholder="用户名 (字母/数字/下划线)"
                :prefix-icon="User"
                class="custom-input"
                maxlength="15"
            />
          </div>

          <div class="input-group">
            <el-input
                v-model="form.password"
                type="password"
                placeholder="密码 (8-15位)"
                :prefix-icon="Lock"
                show-password
                class="custom-input"
                maxlength="15"
            />
          </div>

          <div class="input-group animate-height" v-if="!isLogin">
            <el-input
                v-model="form.confirmPassword"
                type="password"
                placeholder="请再次确认密码"
                :prefix-icon="Lock"
                show-password
                class="custom-input"
                maxlength="15"
            />
          </div>

          <div class="input-group animate-height" v-if="!isLogin">
            <el-input
                v-model="form.phone"
                placeholder="手机号 (11位数字)"
                :prefix-icon="Iphone"
                class="custom-input"
                maxlength="11"
                @input="handlePhoneInput"
            />
          </div>

          <button
              type="button"
              class="submit-btn"
              @click="handleSubmit"
              :disabled="loading"
          >
            <span v-if="!loading">{{ isLogin ? '立 即 登 录' : '注 册 账 号' }}</span>
            <el-icon v-else class="is-loading"><Loading /></el-icon>
          </button>

          <div class="toggle-area">
            <span v-if="isLogin">还没有账号？ <a @click="toggleMode">去注册</a></span>
            <span v-else>已有账号？ <a @click="toggleMode">去登录</a></span>
          </div>
        </el-form>
      </div>
    </div>
  </el-dialog>
</template>

<script setup>
import { ref, reactive, computed } from 'vue'
import { User, Lock, Iphone, Goods, Close, Loading } from '@element-plus/icons-vue'
import request from '@/utils/request'
import { ElMessage } from 'element-plus'

const props = defineProps(['modelValue'])
const emit = defineEmits(['update:modelValue', 'success'])

const visible = computed({
  get: () => props.modelValue,
  set: (val) => emit('update:modelValue', val)
})

const isLogin = ref(true)
const loading = ref(false)
const form = reactive({ username: '', password: '', confirmPassword: '', phone: '' })

const close = () => { visible.value = false }

const toggleMode = () => {
  isLogin.value = !isLogin.value
  Object.assign(form, { username: '', password: '', confirmPassword: '', phone: '' })
}

const handlePhoneInput = (val) => {
  form.phone = val.replace(/[^\d]/g, '').slice(0, 11)
}

const handleSubmit = async () => {
  // ★★★ 关键点3：空值拦截，直接 Return，不发请求，也不刷新 ★★★
  if (!form.username || !form.password) {
    // 这里是你要求的提示语
    ElMessage.warning('请输入账号和密码')
    return
  }

  // 注册时的额外校验
  if (!isLogin.value) {
    if (form.password !== form.confirmPassword) {
      ElMessage.error('两次输入的密码不一致')
      return
    }
  }

  loading.value = true
  const url = isLogin.value ? '/api/login' : '/api/register'

  try {
    const submitData = {
      username: form.username,
      password: form.password,
      phone: isLogin.value ? undefined : form.phone
    }

    const res = await request.post(url, submitData)

    if (isLogin.value) {
      ElMessage.success('登录成功，欢迎回来！')
      localStorage.setItem('token', res.token)
      localStorage.setItem('user', JSON.stringify(res.user))
      emit('success', res.user)
      close()
    } else {
      ElMessage.success('注册成功，请登录')
      toggleMode()
    }
  } catch (err) {
    // ★★★ 关键点4：捕获错误，什么都不做，保持弹窗开启 ★★★
    // 具体的错误提示已由 request.js 弹出，这里只需停止 loading
    console.error('操作失败:', err)
  } finally {
    loading.value = false
  }
}
</script>

<style lang="scss">
.auth-dialog {
  border-radius: 24px !important; overflow: hidden; padding: 0 !important;
  .el-dialog__header { display: none; } .el-dialog__body { padding: 0 !important; }
}
</style>

<style scoped lang="scss">
.auth-container { display: flex; height: 540px; }

.auth-left {
  width: 320px; background: linear-gradient(135deg, #ffda44 0%, #ffc107 100%);
  position: relative; display: flex; flex-direction: column; justify-content: center; align-items: center; color: #333; overflow: hidden;
  .brand { font-size: 32px; font-weight: 800; display: flex; align-items: center; gap: 10px; z-index: 2; }
  .slogan { text-align: center; margin-top: 20px; z-index: 2; h3 { margin: 0; font-size: 20px; } p { margin: 8px 0 0; opacity: 0.8; font-size: 14px; } }
  .circle { position: absolute; border-radius: 50%; background: rgba(255, 255, 255, 0.2); }
  .c1 { width: 200px; height: 200px; top: -50px; left: -50px; }
  .c2 { width: 150px; height: 150px; bottom: -30px; right: -30px; }
}

.auth-right {
  flex: 1; padding: 40px 50px; position: relative; display: flex; flex-direction: column; justify-content: center; background: #fff;
  .close-btn { position: absolute; top: 20px; right: 20px; background: transparent; border: none; font-size: 20px; cursor: pointer; color: #999; &:hover { color: #333; } }
  h2 { margin: 0 0 8px; font-size: 28px; color: #333; }
  .subtitle { margin: 0 0 24px; color: #999; font-size: 14px; }
  .input-group { margin-bottom: 16px; }

  :deep(.custom-input .el-input__wrapper) {
    border-radius: 12px; padding: 8px 15px; box-shadow: 0 0 0 1px #eee; background: #f9f9f9; transition: all 0.3s;
    &.is-focus { box-shadow: 0 0 0 2px #ffda44 !important; background: #fff; }
  }

  .submit-btn {
    width: 100%; height: 48px; background: #333; color: #ffda44; border: none; border-radius: 12px; font-size: 16px; font-weight: bold; cursor: pointer; margin-top: 10px; transition: all 0.3s; display: flex; align-items: center; justify-content: center;
    &:hover { background: #000; transform: translateY(-2px); box-shadow: 0 8px 20px rgba(0,0,0,0.15); }
    &:disabled { opacity: 0.7; cursor: not-allowed; transform: none; }
  }

  .toggle-area {
    margin-top: 24px; text-align: center; font-size: 14px; color: #666;
    a { color: #ff9500; font-weight: bold; cursor: pointer; margin-left: 5px; &:hover { text-decoration: underline; } }
  }
}

.animate-height { animation: slideDown 0.3s ease; }
@keyframes slideDown { from { opacity: 0; transform: translateY(-10px); } to { opacity: 1; transform: translateY(0); } }
</style>