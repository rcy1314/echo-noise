<template>
  <div class="builtin-comments">
    <div class="waline-wrapper p-3 rounded-lg" :class="themeBg">
      <div class="text-sm mb-2" :class="themeText">评论</div>
      <div v-if="siteConfig?.commentEmailEnabled" class="text-xs mb-2" :class="themeMuted">新评论或回复会发送通知邮件</div>
      <div class="space-y-3 mb-4">
        <div class="grid grid-cols-1 md:grid-cols-3 gap-3">
          <UInput v-model="nick" placeholder="昵称" />
          <UInput v-model="mail" placeholder="邮箱" />
          <UInput v-model="link" placeholder="网址（可选）" />
        </div>
        <UTextarea v-model="content" :rows="3" :maxlength="500" placeholder="写点什么..." />
        <div class="flex justify-end">
          <UButton color="green" @click="submit">{{ replyTo ? '发布回复' : '发布评论' }}</UButton>
        </div>
      </div>

      <div v-if="rootComments.length" class="space-y-4">
        <div v-for="c in rootComments" :key="c.id" class="rounded-md p-3" :class="themeItem">
          <div class="text-xs flex items-center justify-between" :class="themeMuted">
            <span>{{ c.nick || '匿名' }}</span>
            <span>{{ formatDate(c.created_at) }}</span>
          </div>
          <div class="mt-1 text-sm" :class="themeText">{{ c.content }}</div>
          <div class="mt-2">
            <UButton size="xs" variant="ghost" @click="startReply(c.id, c.nick || '匿名')">回复</UButton>
          </div>
          <div v-if="childrenMap[c.id]?.length" class="mt-3 pl-4 border-l border-white/10 space-y-2">
            <div v-for="child in childrenMap[c.id]" :key="child.id" class="rounded-md p-3" :class="themeItem">
              <div class="text-xs flex items-center justify-between" :class="themeMuted">
                <span>{{ child.nick || '匿名' }}</span>
                <span>{{ formatDate(child.created_at) }}</span>
              </div>
              <div class="mt-1 text-sm" :class="themeText">{{ child.content }}</div>
              <div class="mt-2">
                <UButton size="xs" variant="ghost" @click="startReply(child.id, child.nick || '匿名')">回复</UButton>
              </div>
            </div>
          </div>
        </div>
      </div>
      <div v-else class="text-xs text-gray-400">暂无评论</div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, watch, computed } from 'vue'
import { useToast } from '#ui/composables/useToast'
import { getRequest, postRequest } from '~/utils/api'

const props = defineProps<{ messageId: number, siteConfig: any }>()
const comments = ref<any[]>([])
const nick = ref('')
const mail = ref('')
const link = ref('')
const content = ref('')
const replyTo = ref<number | null>(null)

const themeBg = computed(() => 'bg-[rgba(36,43,50,0.35)]')
const themeText = computed(() => 'text-gray-200')
const themeMuted = computed(() => 'text-gray-400')
const themeItem = computed(() => 'bg-[rgba(36,43,50,0.25)]')

const load = async () => {
  try {
    const res = await getRequest<any>(`messages/${props.messageId}/comments`, undefined, { credentials: 'include' })
    if (res && res.code === 1) {
      comments.value = res.data || []
    } else {
      comments.value = []
    }
  } catch (e) {
    comments.value = []
  }
}

const submit = async () => {
  try {
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
      await load()
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
watch(() => props.messageId, load)

const startReply = (id: number, nickName: string) => {
  replyTo.value = id
  if (!content.value.startsWith(`@${nickName} `)) content.value = `@${nickName} ` + content.value
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
.waline-wrapper { border: 1px solid rgba(255,255,255,0.06); }
</style>
