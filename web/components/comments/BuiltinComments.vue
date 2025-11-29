<template>
  <div class="builtin-comments">
  <div class="waline-wrapper p-3 rounded-lg border" :class="[themeBg, themeBorder]">
      <div class="text-sm mb-2" :class="themeText">评论区</div>
      <div class="text-[10px] mb-1 opacity-60" :class="themeMuted">已加载 {{ comments.length }} 条</div>
      <div v-if="sortedRootComments.length" class="space-y-4 mb-4">
        <div v-for="c in visibleRootComments" :key="c.id" class="rounded-md p-4 border" :class="[themeItem, themeBorder]">
          <div class="text-xs flex items-center justify-between" :class="themeMuted">
            <template v-if="safeURL(c.link)">
              <a :href="safeURL(c.link)" target="_blank" rel="noopener noreferrer">{{ c.nick || '匿名' }}</a>
            </template>
            <template v-else>
              <span>{{ c.nick || '匿名' }}</span>
            </template>
            <span>{{ formatDate(c.created_at) }}</span>
          </div>
          <div class="mt-1 text-sm" :class="themeText"><MarkdownRenderer :content="c.content" /></div>
          <div class="mt-2 flex gap-2">
            <button class="text-xs px-2 py-1 rounded border" @click="startReply(c.id, c.nick || '匿名')">回复</button>
            <button v-if="isAdmin" class="text-xs px-2 py-1 rounded border border-red-500 text-red-500" @click="confirmDelete(c.id)">删除</button>
          </div>
          <div v-if="childrenMap[c.id]?.length" class="mt-4 pl-4 border-l space-y-2" :class="childBorder">
            <div v-for="child in visibleChildren(c.id)" :key="child.id" class="rounded-md p-3" :class="[isDark ? 'bg-[rgba(24,28,32,0.6)]' : 'bg-gray-50']">
              <div class="text-xs flex items-center justify-between" :class="themeMuted">
                <template v-if="safeURL(child.link)">
                  <a :href="safeURL(child.link)" target="_blank" rel="noopener noreferrer">{{ child.nick || '匿名' }}</a>
                </template>
                <template v-else>
                  <span>{{ child.nick || '匿名' }}</span>
                </template>
                <span>{{ formatDate(child.created_at) }}</span>
              </div>
              <div v-if="replyNickMap[child.id]" class="text-xs mt-1" :class="themeMuted">回复 @{{ replyNickMap[child.id] }}</div>
              <div class="mt-1 text-sm" :class="themeText"><MarkdownRenderer :content="child.content" /></div>
              <div class="mt-2 flex gap-2">
                <button class="text-xs px-2 py-1 rounded border" @click="startReply(child.id, child.nick || '匿名')">回复</button>
                <button v-if="isAdmin" class="text-xs px-2 py-1 rounded border border-red-500 text-red-500" @click="confirmDelete(child.id)">删除</button>
              </div>
            </div>
            <div v-if="hasMoreReplies(c.id)" class="flex justify-start">
              <button class="text-xs px-2 py-1 rounded border" :class="themeBorder" @click="loadMoreReplies(c.id)">加载更多回复</button>
            </div>
          </div>
        </div>
        <div v-if="hasMore" class="flex justify-center">
          <button class="text-xs px-3 py-1 rounded border" :class="themeBorder" @click="loadMore">加载更多评论</button>
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
        <div class="flex flex-wrap items-center gap-2 mb-2">
          <button class="text-xs px-2 py-1 rounded border" :class="themeBorder" @click="applyFormat('bold')">加粗</button>
          <button class="text-xs px-2 py-1 rounded border" :class="themeBorder" @click="applyFormat('italic')">斜体</button>
          <button class="text-xs px-2 py-1 rounded border" :class="themeBorder" @click="applyFormat('link')">链接</button>
          <button class="text-xs px-2 py-1 rounded border" :class="themeBorder" @click="applyFormat('image')">图片</button>
          <div class="relative">
            <button class="text-xs px-2 py-1 rounded border" :class="themeBorder" @click="toggleEmoji">表情</button>
            <div v-if="showEmoji" class="absolute z-10 mt-1 p-2 rounded border bg-white shadow" :class="themeBorder">
              <div class="flex flex-wrap gap-1 w-56">
                <button v-for="e in emojis" :key="e" class="px-2 py-1 text-sm" @click="insertEmoji(e)">{{ e }}</button>
              </div>
            </div>
          </div>
        </div>
        <div class="relative">
          <textarea ref="taRef" v-model="content" :class="textareaClass" rows="4" placeholder="写下你的评论..." @input="onInput" @keydown="onKeydown" @blur="hideMention" />
          <div v-if="showMention" class="absolute left-2 top-2 mt-0 z-10 bg-white rounded border shadow max-h-40 overflow-auto" :class="[themeBorder, isDark ? 'text-black' : '']">
            <div class="px-2 py-1 text-xs">@{{ mentionQuery || '昵称' }}</div>
            <div>
              <button v-for="(n,i) in filteredNicks" :key="n" class="block w-full text-left px-2 py-1 text-sm hover:bg-gray-100" :class="i===mentionIndex? 'bg-gray-100' : ''" @mousedown.prevent="chooseNick(n)">{{ n }}</button>
            </div>
          </div>
        </div>
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
const taRef = ref<any>(null)
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
    comments.value = (list || []).sort((a, b) => new Date(b.created_at).getTime() - new Date(a.created_at).getTime())
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

const safeURL = (s: string) => {
  const url = String(s || '').trim()
  if (!url) return ''
  if (/^https?:\/\//i.test(url)) return url
  return ''
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

const allNicks = computed(() => {
  const list = Array.isArray(comments.value) ? comments.value : []
  const set = new Set<string>()
  list.forEach((c: any) => { const n = String(c?.nick || '').trim(); if (n) set.add(n) })
  return Array.from(set)
})
const showMention = ref(false)
const mentionQuery = ref('')
const mentionIndex = ref(0)
const filteredNicks = computed(() => {
  const q = mentionQuery.value.toLowerCase()
  const arr = allNicks.value.filter(n => n.toLowerCase().startsWith(q))
  return arr.slice(0, 20)
})
const hideMention = () => { showMention.value = false; mentionIndex.value = 0; mentionQuery.value = '' }
const openMention = () => { showMention.value = true; mentionIndex.value = 0 }
const getCaret = () => {
  const el = taRef.value as HTMLTextAreaElement
  if (!el) return { start: 0, end: 0 }
  return { start: el.selectionStart || 0, end: el.selectionEnd || 0 }
}
const replaceRange = (text: string, start: number, end: number, insert: string) => {
  const before = text.slice(0, start)
  const after = text.slice(end)
  return before + insert + after
}
const computeMention = () => {
  const el = taRef.value as HTMLTextAreaElement
  if (!el) return
  const pos = el.selectionStart || 0
  const s = content.value
  let i = pos - 1
  while (i >= 0 && s[i] !== '\n' && s[i] !== ' ') i--
  const start = i + 1
  if (s[start] !== '@') { hideMention(); return }
  const end = pos
  const q = s.slice(start + 1, end)
  mentionQuery.value = q
  openMention()
}
const onInput = () => { computeMention() }
const onKeydown = (e: KeyboardEvent) => {
  if (e.key === '@') { nextTick(computeMention); return }
  if (!showMention.value) return
  if (e.key === 'ArrowDown') { e.preventDefault(); mentionIndex.value = Math.min(mentionIndex.value + 1, filteredNicks.value.length - 1) }
  else if (e.key === 'ArrowUp') { e.preventDefault(); mentionIndex.value = Math.max(mentionIndex.value - 1, 0) }
  else if (e.key === 'Enter') {
    e.preventDefault()
    const n = filteredNicks.value[mentionIndex.value]
    if (n) chooseNick(n)
  } else if (e.key === 'Escape') { hideMention() }
}
const chooseNick = (nick: string) => {
  const el = taRef.value as HTMLTextAreaElement
  if (!el) return
  const pos = el.selectionStart || 0
  const s = content.value
  let i = pos - 1
  while (i >= 0 && s[i] !== '\n' && s[i] !== ' ') i--
  const start = i + 1
  const end = pos
  content.value = replaceRange(s, start, end, `@${nick} `)
  hideMention()
  nextTick(() => { const p = start + nick.length + 2; el.setSelectionRange(p, p); el.focus() })
}

const showEmoji = ref(false)
const emojis = ['😀','😄','😁','😆','😊','😍','🤔','👍','🔥','🎉','❤️','🥳','✨','🌟','🍀']
const toggleEmoji = () => { showEmoji.value = !showEmoji.value }
const insertAtCaret = (text: string) => {
  const el = taRef.value as HTMLTextAreaElement
  if (!el) { content.value += text; return }
  const { start, end } = getCaret()
  content.value = replaceRange(content.value, start, end, text)
  const p = start + text.length
  nextTick(() => { el.setSelectionRange(p, p); el.focus() })
}
const insertEmoji = (e: string) => { insertAtCaret(e) ; showEmoji.value = false }
const applyFormat = (type: string) => {
  const el = taRef.value as HTMLTextAreaElement
  const { start, end } = getCaret()
  const sel = content.value.slice(start, end)
  if (type === 'bold') insertAtCaret(sel ? `**${sel}**` : `**加粗**`)
  else if (type === 'italic') insertAtCaret(sel ? `*${sel}*` : `*斜体*`)
  else if (type === 'link') {
    const url = window.prompt('请输入链接地址', 'https://') || ''
    if (/^https?:\/\//i.test(url)) insertAtCaret(sel ? `[${sel}](${url})` : `[链接文本](${url})`)
  } else if (type === 'image') {
    const url = window.prompt('请输入图片地址', 'https://') || ''
    if (/^https?:\/\//i.test(url)) insertAtCaret(`![图片](${url})`)
  }
}

const rootComments = computed(() => {
  const list = Array.isArray(comments.value) ? comments.value : []
  const roots = list.filter((c: any) => c && (c.parent_id === null || Number(c.parent_id || 0) === 0))
  return roots
})
const byId = computed(() => {
  const m: Record<number, any> = {}
  const list = Array.isArray(comments.value) ? comments.value : []
  list.forEach((c: any) => { m[Number(c.id)] = c })
  return m
})
const childrenWithTarget = computed(() => {
  const map: Record<number, any[]> = {}
  const targetMap: Record<number, string> = {}
  const list = Array.isArray(comments.value) ? comments.value : []
  const mentionReg = /^@([^\s:：]+)/
  list.forEach((c: any) => {
    const pid = Number(c?.parent_id || 0)
    if (pid > 0) {
      const parent = byId.value[pid]
      if (parent) {
        targetMap[Number(c.id)] = String(parent.nick || '')
        let rootNode: any = parent
        while (Number(rootNode?.parent_id || 0) > 0) {
          const next = byId.value[Number(rootNode.parent_id)]
          if (!next) break
          rootNode = next
        }
        const key = Number(rootNode.id)
        if (!map[key]) map[key] = []
        map[key].push(c)
        return
      }
      return
    }
    const text = String(c?.content || '')
    const mm = text.match(mentionReg)
    if (mm && mm[1]) {
      const nick = mm[1]
      const ct = new Date(c.created_at).getTime()
      const candidates = list.filter((r: any) => String(r?.nick || '').trim() === nick)
      let parentCandidate: any = null
      const earlier = candidates.filter((r: any) => new Date(r.created_at).getTime() <= ct)
      if (earlier.length) {
        parentCandidate = earlier.sort((a: any, b: any) => new Date(b.created_at).getTime() - new Date(a.created_at).getTime())[0]
      } else if (candidates.length) {
        parentCandidate = candidates.sort((a: any, b: any) => new Date(b.created_at).getTime() - new Date(a.created_at).getTime())[0]
      }
      if (parentCandidate && parentCandidate.id) {
        let rootNode: any = parentCandidate
        while (Number(rootNode?.parent_id || 0) > 0) {
          const next = byId.value[Number(rootNode.parent_id)]
          if (!next) break
          rootNode = next
        }
        const key = Number(rootNode.id)
        if (!map[key]) map[key] = []
        if (!map[key].some((x) => x.id === c.id)) map[key].push(c)
        targetMap[Number(c.id)] = String(parentCandidate.nick || '')
      }
    }
  })
  Object.keys(map).forEach((k) => {
    map[Number(k)] = (map[Number(k)] || []).sort((a: any, b: any) => new Date(b.created_at).getTime() - new Date(a.created_at).getTime())
  })
  return { map, targetMap }
})
const childrenMap = computed(() => childrenWithTarget.value.map)
const replyNickMap = computed(() => childrenWithTarget.value.targetMap)
const sortedRootComments = computed(() => {
  const mentionReg = /^@([^\s:：]+)/
  const roots = (Array.isArray(rootComments.value) ? rootComments.value : []).filter((c: any) => {
    const text = String(c?.content || '')
    return !mentionReg.test(text)
  })
  return roots.slice().sort((a: any, b: any) => new Date(b.created_at).getTime() - new Date(a.created_at).getTime())
})
const visibleCount = ref(3)
const visibleRootComments = computed(() => sortedRootComments.value.slice(0, visibleCount.value))
const hasMore = computed(() => sortedRootComments.value.length > visibleCount.value)
const loadMore = () => { visibleCount.value += 3 }
watch(() => props.messageId, () => { visibleCount.value = 3 })

const visibleChildrenCount = ref<Record<number, number>>({})
const visibleChildren = (rootId: number) => {
  const n = visibleChildrenCount.value[rootId] ?? 3
  return (childrenMap.value[rootId] || []).slice(0, n)
}
const hasMoreReplies = (rootId: number) => {
  const total = (childrenMap.value[rootId] || []).length
  const n = visibleChildrenCount.value[rootId] ?? 3
  return total > n
}
const loadMoreReplies = (rootId: number) => {
  const cur = visibleChildrenCount.value[rootId] ?? 3
  visibleChildrenCount.value[rootId] = cur + 3
}
watch(childrenMap, (m) => {
  const next: Record<number, number> = { ...visibleChildrenCount.value }
  Object.keys(m || {}).forEach((k) => {
    const id = Number(k)
    if (!next[id]) next[id] = 3
  })
  visibleChildrenCount.value = next
})
watch(() => props.messageId, () => { visibleChildrenCount.value = {} })
</script>

<style scoped>
</style>

 
