<template>
  <div class="builtin-comments">
    <div class="waline-wrapper p-3 rounded-lg border" :class="[themeBg, themeBorder]">
      <div class="text-sm mb-2" :class="themeText">评论区</div>
      <div class="text-[10px] mb-1 opacity-60" :class="themeMuted">已加载 {{ comments.length }} 条</div>
      <div v-if="rootComments.length" class="space-y-4 mb-4">
        <div v-for="c in rootComments" :key="c.id" class="rounded-md p-3 border" :class="[themeItem, themeBorder]">
          <div class="text-xs flex items-center justify-between" :class="themeMuted">
            <span>{{ c.nick || '匿名' }}</span>
            <span>{{ formatDate(c.created_at) }}</span>
          </div>
          <div class="mt-1 text-sm" :class="themeText"><MarkdownRenderer :content="c.content" /></div>
          <div class="mt-2 flex gap-2">
            <button class="text-xs px-2 py-1 rounded border" @click="startReply(c.id, c.nick || '匿名')">回复</button>
            <button v-if="isAdmin" class="text-xs px-2 py-1 rounded border border-red-500 text-red-500" @click="confirmDelete(c.id)">删除</button>
          </div>
          <div v-if="childrenMap[c.id]?.length" class="mt-3 pl-4 border-l space-y-2" :class="childBorder">
            <div v-for="child in childrenMap[c.id]" :key="child.id" class="rounded-md p-3 border" :class="[themeItem, themeBorder]">
              <div class="text-xs flex items中心 justify之间" :class="themeMuted">
                <span>{{ child.nick || '匿名' }}</span>
                <span>{{ formatDate(child.created_at) }}</span>
              </div>
              <div class="mt-1 text-sm" :class="themeText"><MarkdownRenderer :content="child.content" /></div>
              <div class="mt-2 flex gap-2">
                <button class="text-xs px-2 py-1 rounded border" @click="startReply(child.id, child.nick || '匿名')">回复</button>
                <button v-if="isAdmin" class="text-xs px-2 py-1 rounded border border-red-500 text-red-500" @click="confirmDelete(child.id)">删除</button>
              </div>
            </div>
          </div>
        </div>
      </div>
      <div v-else class="text-xs mb-4" :class="themeMuted">暂无评论</div>

      <div v-if="(props.showInput || !!replyTo) && siteConfig?.commentEmailEnabled" class="text-xs mb-2" :class="themeMuted">新评论或回复会发送通知邮件</div>
      <div v-if="props.showInput || !!replyTo" class="space-y-3">
        <div class="grid grid-cols-1 md:grid-cols-3 gap-3">
          <input v-model="nick" placeholder="昵称" :class="inputNickClass" />
          <input v-model="mail" placeholder="邮箱" :class="inputMailClass" />
          <input v-model="link" placeholder="网址（可选）" :class="inputLinkClass" />
        </div>
        <textarea v-model="content" :class="textareaClass" rows="4" placeholder="写下你的评论..." />
        <div class="flex justify-end">
          <button class="px-3 py-1 rounded bg-green-500 text-white" @click="submit">{{ replyTo ? '发布回复' : '发布评论' }}</button>
        </div>
      </div>
      
    </div>
    
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, watch, computed, nextTick, inject, onBeforeUnmount } from 'vue'
import MarkdownRenderer from '~/components/index/MarkdownRenderer.vue'
import { useToast } from '#ui/composables/useToast'
import { getRequest, postRequest, deleteRequest } from '~/utils/api'
import { useUserStore } from '~/store/user'

const props = defineProps<{ messageId: number, siteConfig: any, showInput?: boolean }>()
const comments = ref<any[]>([])
const nick = ref('')
const mail = ref('')
const link = ref('')
const content = ref('')
const replyTo = ref<number | null>(null)
const nickError = ref(false)
const mailError = ref(false)
const deleteId = ref<number | null>(null)
const user = useUserStore()
const isAdmin = computed(() => !!(user.user as any)?.is_admin)
// 使用原始 textarea 输入框

// 主题注入，严格跟随页面当前模式
const injectedTheme = inject('contentTheme', ref('light')) as any
const isDark = computed(() => {
  const v = (injectedTheme && typeof injectedTheme.value !== 'undefined') ? injectedTheme.value : injectedTheme
  return String(v || 'light') === 'dark'
})

const themeBg = computed(() => (isDark.value ? 'bg-[rgba(24,28,32,0.95)]' : 'bg-white'))
const themeBorder = computed(() => (isDark.value ? 'border-white/20' : 'border-black'))
const themeText = computed(() => (isDark.value ? 'text-gray-200' : 'text-black'))
const themeMuted = computed(() => (isDark.value ? 'text-gray-400' : 'text-gray-500'))
const themeItem = computed(() => (isDark.value ? 'bg-[rgba(24,28,32,0.7)]' : 'bg-white'))
const childBorder = computed(() => (isDark.value ? 'border-white/20' : 'border-black'))
const inputBaseLight = 'w-full px-2 py-1 border border-black rounded ring-0 focus:ring-0 focus:border-black'
const inputDark = 'w-full px-2 py-1 bg-[rgba(24,28,32,0.85)] text-white border border-white/20 rounded focus:ring-0 focus:border-primary-400 placeholder:text-gray-400'
const inputNickClass = computed(() => (nickError.value ? 'ring-1 ring-red-500' : (isDark.value ? inputDark : `bg-gray-50 ${inputBaseLight}`)))
const inputMailClass = computed(() => (mailError.value ? 'ring-1 ring-red-500' : (isDark.value ? inputDark : `bg-gray-50 ${inputBaseLight}`)))
const inputLinkClass = computed(() => (isDark.value ? inputDark : `bg-white ${inputBaseLight}`))
const textareaClass = computed(() => (isDark.value ? inputDark : `bg-white ${inputBaseLight}`))

const BASE_API = useRuntimeConfig().public.baseApi || '/api'
const load = async () => {
  try {
    const tryFetch = async (url: string) => {
      const resp = await fetch(url, { credentials: 'include', headers: { 'Accept': 'application/json' } })
      if (!resp || !resp.ok) return null
      const js = await resp.json()
      if (!js || js.code !== 1 || !Array.isArray(js.data)) return []
      return js.data
    }
    const origin = typeof window !== 'undefined' ? window.location.origin : ''
    const urls = [
      `${BASE_API}/messages/${props.messageId}/comments`,
      `${origin}/api/messages/${props.messageId}/comments`,
      `http://localhost:1315/api/messages/${props.messageId}/comments`,
      `http://127.0.0.1:1315/api/messages/${props.messageId}/comments`
    ]
    let list: any[] = []
    for (const u of urls) {
      const data = await tryFetch(u)
      if (data && data.length >= 0) {
        list = data
        if (list.length > 0) break
      }
    }
    comments.value = (list || []).sort((a, b) => new Date(a.created_at).getTime() - new Date(b.created_at).getTime())
  } catch (e) {
    comments.value = []
  }
  await nextTick()
  const container = document.querySelector(`.content-container[data-msg-id="${props.messageId}"] .builtin-comments`)
  const items = container?.querySelectorAll('.rounded-md')
  const target = items && items.length ? (items[items.length - 1] as HTMLElement) : null
  target?.scrollIntoView({ behavior: 'smooth', block: 'center' })
}

const submit = async () => {
  try {
    nickError.value = !nick.value.trim()
    mailError.value = !mail.value.trim()
    if (nickError.value || mailError.value) {
      useToast().add({ title: '缺少必填字段', description: (!nick.value ? '昵称 ' : '') + (!mail.value ? '邮箱' : ''), color: 'red' })
      return
    }
    const md = content.value.trim()
    const payload: any = { nick: nick.value.trim(), mail: mail.value.trim(), link: link.value.trim(), content: md }
    if (!payload.content) {
      useToast().add({ title: '内容不能为空', color: 'red' })
      return
    }
    if (payload.mail && !/^\S+@\S+\.\S+$/.test(payload.mail)) {
      useToast().add({ title: '邮箱格式不正确', color: 'red' })
      return
    }
    if (replyTo.value) payload.parent_id = replyTo.value
    const res = await postRequest<any>(`messages/${props.messageId}/comments`, payload, { credentials: 'include' })
    if (res && res.code === 1) {
      content.value = ''
      replyTo.value = null
      comments.value = [...comments.value, res.data]
      await load()
      await nextTick()
      const container = document.querySelector(`.content-container[data-msg-id="${props.messageId}"] .builtin-comments`)
      const items = container?.querySelectorAll('.rounded-md')
      const target = items && items.length ? (items[items.length - 1] as HTMLElement) : null
      target?.scrollIntoView({ behavior: 'smooth', block: 'center' })
      useToast().add({ title: '已发布', color: 'green' })
    } else {
      useToast().add({ title: '发布失败', description: res?.msg, color: 'red' })
    }
  } catch (e: any) {
    useToast().add({ title: '发布失败', color: 'red' })
  }
}

const formatDate = (s: string) => {
  const d = new Date(s)
  return d.toLocaleString()
}

onMounted(load)
// 保持与父组件的显示控制，但不再初始化富文本编辑器
// 监听来自父级的刷新事件（每次展开评论时确保重新拉取）
const handler = () => load()
onMounted(() => {
  window.addEventListener(`refresh-comments-${props.messageId}`, handler)
})
onBeforeUnmount(() => {
  window.removeEventListener(`refresh-comments-${props.messageId}`, handler)
})
watch(() => props.messageId, load)

const startReply = (id: number, nickName: string) => {
  replyTo.value = id
  if (!content.value.startsWith(`@${nickName} `)) content.value = `@${nickName} ` + content.value
}

const confirmDelete = (id: number) => {
  deleteId.value = id
  if (confirm('确认删除该评论吗？此操作不可恢复。')) {
    doDelete()
  } else {
    deleteId.value = null
  }
}

const doDelete = async () => {
  if (!deleteId.value) return
  try {
    const res = await deleteRequest<any>(`messages/${props.messageId}/comments/${deleteId.value}`, undefined, { credentials: 'include' })
    if (res && res.code === 1) {
      comments.value = comments.value.filter(c => c.id !== deleteId.value)
      useToast().add({ title: '已删除', color: 'green' })
      scrollToMessage()
    } else {
      useToast().add({ title: '删除失败', description: res?.msg, color: 'red' })
    }
  } catch (e: any) {
    useToast().add({ title: '删除失败', color: 'red' })
  } finally {
    deleteId.value = null
  }
}

const scrollToMessage = () => {
  const el = document.querySelector(`.content-container[data-msg-id="${props.messageId}"]`)
  el?.scrollIntoView({ behavior: 'smooth', block: 'start' })
}

const rootComments = computed(() => {
  const list = Array.isArray(comments.value) ? comments.value : []
  const roots = list.filter((c: any) => c && (c.parent_id === null || Number(c.parent_id || 0) === 0))
  return roots.length ? roots : list
})
const childrenMap = computed(() => {
  const map: Record<number, any[]> = {}
  const list = Array.isArray(comments.value) ? comments.value : []
  list.forEach((c: any) => {
    const pid = Number(c?.parent_id || 0)
    if (pid > 0) {
      if (!map[pid]) map[pid] = []
      map[pid].push(c)
    }
  })
  return map
})
</script>

<style scoped>
</style>

 
