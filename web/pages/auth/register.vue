<template>
  <div class="fixed inset-0 bg-gradient-to-br from-slate-900 via-indigo-950 to-slate-800">
    <div class="absolute inset-0 backdrop-blur-xl"></div>
    <div class="relative z-10 flex min-h-screen items-center justify-center p-4">
      <UCard class="w-full max-w-md bg-slate-900/70 text-white border border-slate-700/40 shadow-2xl">
        <div class="flex items-center justify-between mb-4">
          <div class="flex items-center gap-2">
            <UIcon name="i-heroicons-user-plus" class="w-6 h-6 text-indigo-300" />
            <h1 class="text-lg font-semibold">注册</h1>
          </div>
          <NuxtLink to="/auth/login" class="text-sm text-indigo-300 hover:text-indigo-200">去登录</NuxtLink>
        </div>
        <UForm @submit.prevent="onSubmit">
          <UFormGroup label="用户名" class="mb-3">
            <UInput v-model="form.username" placeholder="请输入用户名" />
          </UFormGroup>
          <UFormGroup label="密码" class="mb-3">
            <UInput v-model="form.password" type="password" placeholder="请输入密码" />
          </UFormGroup>
          <UFormGroup label="验证码" class="mb-2">
            <div class="flex items-center gap-2">
              <UInput v-model="form.captcha" placeholder="请输入验证码" />
              <img :src="captchaSrc" @click="refreshCaptcha" class="h-10 w-24 rounded border border-slate-700/40 cursor-pointer" alt="captcha" />
              <UBadge :color="remaining>0 ? 'primary' : 'red'" variant="soft">{{ remaining>0 ? `有效 ${remaining}s` : '已过期' }}</UBadge>
            </div>
          </UFormGroup>
          <div class="flex justify-end">
            <UButton :loading="submitting" :disabled="remaining<=0 || submitting" type="submit" color="primary">{{ remaining>0 ? '注册' : '验证码已过期' }}</UButton>
          </div>
        </UForm>
        <div class="mt-2">
          <UButton class="w-full h-10 px-3 gap-2 justify-center font-medium bg-[#24292f] hover:bg-[#1f2328] text-white ring-1 ring-black/20" @click="loginWithGithub">
            <UIcon name="i-mdi-github" class="w-5 h-5" />
            <span>GitHub 一键注册/登录</span>
          </UButton>
        </div>
      </UCard>
    </div>
  </div>
  <UNotifications />
</template>

<script setup lang="ts">
definePageMeta({ layout: false })
import { useUserStore } from '~/store/user'
import { useToast } from '#imports'
const user = useUserStore()
const toast = useToast()
const router = useRouter()
const baseApi = useRuntimeConfig().public.baseApi || '/api'

const form = reactive({ username: '', password: '', captcha: '' })
const submitting = ref(false)
const captchaSrc = ref('')
const captchaExpiresAt = ref<number | null>(null)
const remaining = ref(0)
let timer: any = null

const refreshCaptcha = async () => {
  try {
    captchaSrc.value = `${baseApi}/captcha?ts=${Date.now()}`
    captchaExpiresAt.value = Date.now() + 120000
    remaining.value = 120
    if (timer) clearInterval(timer)
    timer = setInterval(() => {
      const r = Math.max(0, Math.ceil(((captchaExpiresAt.value || Date.now()) - Date.now()) / 1000))
      remaining.value = r
      if (r <= 0) clearInterval(timer)
    }, 1000)
  } catch (e) {
    remaining.value = 0
  }
}

const onSubmit = async () => {
  try {
    submitting.value = true
    if (!form.username || !form.password || !form.captcha) {
      throw new Error('请完整填写用户名、密码与验证码')
    }
    if ((captchaExpiresAt.value || 0) < Date.now()) {
      throw new Error('验证码已过期，请刷新后再提交')
    }
    const res = await fetch(`${baseApi}/register`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      credentials: 'include',
      body: JSON.stringify({ username: form.username, password: form.password, captcha: form.captcha })
    })
    const data = await res.json().catch(() => ({}))
    if (!res.ok || data.code !== 1) throw new Error(data?.msg || '注册失败')
    toast.add({ title: '注册成功', color: 'green' })
    router.push('/auth/login')
  } catch (e: any) {
    toast.add({ title: '注册失败', description: e.message || '请稍后重试', color: 'red' })
    refreshCaptcha()
  }
  finally {
    submitting.value = false
  }
}

const loginWithGithub = () => {
  window.location.href = `${baseApi}/oauth/github/login`
}
onMounted(() => {
  refreshCaptcha()
  // 进入注册页时校验是否允许注册
  ;(async () => {
    try {
      const res = await fetch(`${baseApi}/frontend/config`, { credentials: 'include' })
      const data = await res.json()
      const allowed = !!data?.data?.allowRegistration
      if (!allowed) {
        useToast().add({ title: '提示', description: '站点已关闭用户注册', color: 'orange' })
        useRouter().push('/auth/login')
      }
    } catch {}
  })()
})
onUnmounted(() => {
  if (timer) clearInterval(timer)
})
</script>
