<template>
  <div class="fixed inset-0 bg-gradient-to-br from-slate-900 via-indigo-950 to-slate-800">
    <div class="absolute inset-0 backdrop-blur-xl"></div>
    <div class="relative z-10 flex min-h-screen items-center justify-center p-4">
      <UCard class="w-full max-w-md bg-slate-900/70 text-white border border-slate-700/40 shadow-2xl">
        <div class="flex items-center justify-between mb-4">
          <div class="flex items-center gap-2">
            <UIcon name="i-heroicons-lock-closed" class="w-6 h-6 text-indigo-300" />
            <h1 class="text-lg font-semibold">登录</h1>
          </div>
          <UButton variant="link" color="indigo" class="text-sm" @click="goRegister">去注册</UButton>
        </div>
        <UForm @submit.prevent="onSubmit">
          <UFormGroup label="用户名/已绑定邮箱" class="mb-3">
            <UInput v-model="form.username" placeholder="请输入用户名或已绑定邮箱" />
          </UFormGroup>
          <UFormGroup label="密码" class="mb-2">
            <UInput v-model="form.password" type="password" placeholder="请输入密码" autocomplete="current-password" autocorrect="off" autocapitalize="off" spellcheck="false" />
          </UFormGroup>
          <div class="flex justify-between items-center mb-3">
            <UButton variant="ghost" size="sm" @click="showForgot = true">忘记密码</UButton>
            <UButton :loading="submitting" :disabled="submitting" type="submit" color="primary">登录</UButton>
          </div>
        </UForm>
        <div class="mt-2" v-if="githubEnabled">
          <UButton class="w-full h-10 px-3 gap-2 justify-center font-medium bg-[#24292f] hover:bg-[#1f2328] text-white ring-1 ring-black/20" @click="loginWithGithub">
            <UIcon name="i-mdi-github" class="w-5 h-5" />
            <span>GitHub 一键登录</span>
          </UButton>
        </div>
      </UCard>
    </div>

    <UModal v-model="showForgot">
      <UCard class="bg-slate-900/80 text-white border border-slate-700/40">
        <div class="font-semibold mb-2">找回密码</div>
        <UForm @submit="onForgot">
          <UFormGroup label="用户名或邮箱" class="mb-3">
            <UInput v-model="forgot.account" placeholder="请输入用户名或邮箱" />
          </UFormGroup>
          <div class="flex justify-end gap-2">
            <UButton variant="ghost" @click="showForgot = false">取消</UButton>
            <UButton :disabled="forgotCooldown>0 || !smtpEnabled" type="submit" color="primary">
              {{ smtpEnabled ? (forgotCooldown>0 ? `请稍候(${forgotCooldown}s)` : '发送重置邮件') : '邮件未开启' }}
            </UButton>
          </div>
        </UForm>
      </UCard>
    </UModal>
  </div>
  <UNotifications />
</template>

<script setup lang="ts">
definePageMeta({ layout: false })
import { useUserStore } from '~/store/user'
import { useToast } from '#imports'
const user = useUserStore()
const toast = useToast()
const route = useRoute()
const router = useRouter()
const baseApi = useRuntimeConfig().public.baseApi || '/api'

const form = reactive({ username: '', password: '' })
const submitting = ref(false)
const githubEnabled = ref(true)
const showForgot = ref(false)
const forgot = reactive({ account: '' })
const forgotCooldown = ref(0)
let forgotTimer: any = null
const smtpEnabled = ref(true)
const allowRegistration = ref(true)

const onSubmit = async () => {
  submitting.value = true
  const controller = new AbortController()
  const timeout = setTimeout(() => {
    controller.abort()
    useToast().add({ title: '登录失败', description: '请求超时或服务器不可用', color: 'red' })
    submitting.value = false
  }, 8000)
  try {
    const ok = await user.login({ username: form.username, password: form.password })
    if (ok) {
      useToast().add({ title: '登录成功', color: 'green' })
      const redirect = (route.query.redirect as string) || '/status'
      router.push(redirect)
    }
  } catch (e) {
    useToast().add({ title: '登录失败', description: '请检查账号密码与后端服务', color: 'red' })
  } finally {
    clearTimeout(timeout)
    submitting.value = false
  }
}

const loginWithGithub = () => {
  window.location.href = `${baseApi}/oauth/github/login`
}

const goRegister = async () => {
  try {
    const res = await fetch(`${baseApi}/frontend/config`, { credentials: 'include' })
    const data = await res.json()
    const allowed = !!data?.data?.allowRegistration
    if (!allowed) {
      useToast().add({ title: '提示', description: '站点已关闭用户注册', color: 'orange' })
      return
    }
    useRouter().push('/auth/register')
  } catch {
    useRouter().push('/auth/register')
  }
}

const onForgot = async () => {
  try {
    const res = await fetch(`${baseApi}/password/forgot`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      credentials: 'include',
      body: JSON.stringify({ account: forgot.account })
    })
    const data = await res.json().catch(() => ({}))
    if (!res.ok || data.code !== 1) throw new Error(data?.msg || '发送失败')
    toast.add({ title: data?.msg || '已发送', description: '请查收重置邮件', color: 'green' })
    forgotCooldown.value = 60
    if (forgotTimer) clearInterval(forgotTimer)
    forgotTimer = setInterval(() => {
      if (forgotCooldown.value > 0) forgotCooldown.value--
      else clearInterval(forgotTimer)
    }, 1000)
  } catch (e: any) {
    toast.add({ title: '失败', description: e.message || '发送失败', color: 'red' })
  }
}

onMounted(async () => {
  const ok = await user.checkLoginStatus()
  if (ok) router.push('/status')
  try {
    const res = await fetch(`${baseApi}/frontend/config`, { credentials: 'include' })
    const data = await res.json()
    githubEnabled.value = !!data?.data?.frontendSettings?.githubOAuthEnabled
    smtpEnabled.value = !!data?.data?.smtpEnabled
    allowRegistration.value = !!data?.data?.allowRegistration
  } catch {}
})
</script>
