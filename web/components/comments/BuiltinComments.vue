<template>
  <div class="builtin-comments">
    <div class="waline-wrapper p-3 rounded-lg border" :class="[themeBg, themeBorder]">
      <div class="text-sm mb-2" :class="themeText">评论</div>
      <div v-if="siteConfig?.commentEmailEnabled" class="text-xs mb-2" :class="themeMuted">新评论或回复会发送通知邮件</div>
      <div class="space-y-3 mb-4">
        <div class="grid grid-cols-1 md:grid-cols-3 gap-3">
          <UInput v-model="nick" placeholder="昵称" :class="inputNickClass" />
          <UInput v-model="mail" placeholder="邮箱" :class="inputMailClass" />
          <UInput v-model="link" placeholder="网址（可选）" :class="inputLinkClass" />
        </div>
        <UTextarea v-model="content" :rows="3" :maxlength="500" placeholder="写点什么..." :class="textareaClass" />
        <div class="flex justify-end">
          <UButton color="green" @click="submit">{{ replyTo ? '发布回复' : '发布评论' }}</UButton>
        </div>
      </div>

      <div v-if="rootComments.length" class="space-y-4">
        <div v-for="c in rootComments" :key="c.id" class="rounded-md p-3 border" :class="[themeItem, themeBorder]">
          <div class="text-xs flex items-center justify-between" :class="themeMuted">
            <span>{{ c.nick || '匿名' }}</span>
            <span>{{ formatDate(c.created_at) }}</span>
          </div>
          <div class="mt-1 text-sm" :class="themeText">{{ c.content }}</div>
          <div class="mt-2 flex gap-2">
            <UButton size="xs" variant="ghost" @click="startReply(c.id, c.nick || '匿名')">回复</UButton>
            <UButton v-if="isAdmin" size="xs" color="red" variant="ghost" @click="confirmDelete(c.id)">删除</UButton>
          </div>
          <div v-if="childrenMap[c.id]?.length" class="mt-3 pl-4 border-l space-y-2" :class="childBorder">
            <div v-for="child in childrenMap[c.id]" :key="child.id" class="rounded-md p-3 border" :class="[themeItem, themeBorder]">
              <div class="text-xs flex items-center justify-between" :class="themeMuted">
                <span>{{ child.nick || '匿名' }}</span>
                <span>{{ formatDate(child.created_at) }}</span>
              </div>
              <div class="mt-1 text-sm" :class="themeText">{{ child.content }}</div>
              <div class="mt-2 flex gap-2">
                <UButton size="xs" variant="ghost" @click="startReply(child.id, child.nick || '匿名')">回复</UButton>
                <UButton v-if="isAdmin" size="xs" color="red" variant="ghost" @click="confirmDelete(child.id)">删除</UButton>
              </div>
            </div>
          </div>
        </div>
      </div>
      <div v-else class="text-xs" :class="themeMuted">暂无评论</div>
    </div>
    <UModal v-model="showConfirm">
      <div class="p-4">
        <div class="mb-3">确认删除该评论吗？此操作不可恢复。</div>
        <div class="flex justify-end gap-2">
          <UButton variant="ghost" @click="showConfirm=false">取消</UButton>
          <UButton color="red" @click="doDelete">删除</UButton>
        </div>
      </div>
    </UModal>
    <UNotifications />
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, watch, computed, nextTick, inject, onBeforeUnmount } from 'vue'
import { useToast } from '#ui/composables/useToast'
import { getRequest, postRequest, deleteRequest } from '~/utils/api'
import { useUserStore } from '~/store/user'

const props = defineProps<{ messageId: number, siteConfig: any }>()
const comments = ref<any[]>([])
const nick = ref('')
const mail = ref('')
const link = ref('')
const content = ref('')
const replyTo = ref<number | null>(null)
const nickError = ref(false)
const mailError = ref(false)
const showConfirm = ref(false)
const deleteId = ref<number | null>(null)
const user = useUserStore()
const isAdmin = computed(() => !!(user.user as any)?.is_admin)

// 主题注入，严格跟随页面当前模式
const injectedTheme = inject('contentTheme', ref('light')) as any
const isDark = computed(() => String(injectedTheme?.value || 'light') === 'dark' || document.documentElement.classList.contains('dark'))

const themeBg = computed(() => (isDark.value ? 'bg-[rgba(24,28,32,0.95)]' : 'bg-white'))
const themeBorder = computed(() => (isDark.value ? 'border-white/20' : 'border-black'))
const themeText = computed(() => (isDark.value ? 'text-gray-200' : 'text-[#111]'))
const themeMuted = computed(() => (isDark.value ? 'text-gray-400' : 'text-gray-500'))
const themeItem = computed(() => (isDark.value ? 'bg-[rgba(24,28,32,0.7)]' : 'bg-white'))
const childBorder = computed(() => (isDark.value ? 'border-white/20' : 'border-black'))
const inputBaseLight = 'border border-black rounded ring-0 focus:ring-0 focus:border-black'
const inputDark = 'bg-[rgba(24,28,32,0.85)] text-white border border-white/20 rounded focus:ring-0 focus:border-primary-400 placeholder:text-gray-400'
const inputNickClass = computed(() => (nickError.value ? 'ring-1 ring-red-500' : (isDark.value ? inputDark : `bg-gray-50 ${inputBaseLight}`)))
const inputMailClass = computed(() => (mailError.value ? 'ring-1 ring-red-500' : (isDark.value ? inputDark : `bg-gray-50 ${inputBaseLight}`)))
const inputLinkClass = computed(() => (isDark.value ? inputDark : `bg-white ${inputBaseLight}`))
const textareaClass = computed(() => (isDark.value ? inputDark : `bg-white ${inputBaseLight}`))

const load = async () => {
  try {
    const res = await getRequest<any>(`messages/${props.messageId}/comments`, undefined, { credentials: 'include' })
    if (res && res.code === 1) {
      const list = Array.isArray(res.data) ? res.data : []
      // 保持按时间正序或倒序，后端若未排序则按 created_at 排序
      comments.value = list.sort((a, b) => new Date(a.created_at).getTime() - new Date(b.created_at).getTime())
    } else {
      comments.value = []
    }
  } catch (e) {
    comments.value = []
  }
}

const submit = async () => {
  try {
    nickError.value = !nick.value.trim()
    mailError.value = !mail.value.trim()
    if (nickError.value || mailError.value) {
      useToast().add({ title: '缺少必填字段', description: (!nick.value ? '昵称 ' : '') + (!mail.value ? '邮箱' : ''), color: 'red' })
      return
    }
    const payload: any = { nick: nick.value.trim(), mail: mail.value.trim(), link: link.value.trim(), content: content.value.trim() }
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
      const el = document.querySelector(`.content-container[data-msg-id="${props.messageId}"]`)
      el?.scrollIntoView({ behavior: 'smooth', block: 'start' })
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
  showConfirm.value = true
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
    showConfirm.value = false
    deleteId.value = null
  }
}

const scrollToMessage = () => {
  const el = document.querySelector(`.content-container[data-msg-id="${props.messageId}"]`)
  el?.scrollIntoView({ behavior: 'smooth', block: 'start' })
}

const rootComments = computed(() => (comments.value || []).filter(c => !c.parent_id))
const childrenMap = computed<Record<number, any[]>>(() => {
  const map: Record<number, any[]> = {}
  (comments.value || []).forEach(c => {
    const pid = c.parent_id
    if (pid) {
      if (!map[pid]) map[pid] = []
      map[pid].push(c)
    }
  })
  return map
})
</script>

<style scoped>
</style>

 
