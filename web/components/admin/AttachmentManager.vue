<template>
  <div class="rounded-xl" :class="theme?.cardBg">
    <div class="px-4 py-3 flex items-center justify-between">
      <div class="font-semibold flex items-center gap-2" :class="theme?.text">
        <UIcon name="i-heroicons-paper-clip" class="w-5 h-5 text-indigo-300" />
        <span>附件管理</span>
      </div>
      <div class="flex items-center gap-2">
        <UButton :loading="loading" color="gray" variant="soft" class="shadow" @click="refresh">刷新</UButton>
      </div>
    </div>
    <div class="px-4">
      <div class="flex gap-2 mb-3">
        <UButton :color="activeTab==='images'?'primary':'gray'" variant="soft" @click="activeTab='images'">图片</UButton>
        <UButton :color="activeTab==='videos'?'primary':'gray'" variant="soft" @click="activeTab='videos'">视频</UButton>
      </div>
      <div v-if="activeTab==='images'">
      <div v-if="images.length===0" class="text-sm" :class="theme?.mutedText">暂无图片附件</div>
      <div v-else class="grid grid-cols-1 md:grid-cols-2 gap-3">
          <div v-for="item in imagesDisplay" :key="item.name" class="rounded-lg border p-3" :class="theme?.border">
            <div class="flex items-center justify-between">
              <div>
                <div class="text-sm" :class="theme?.text">{{ item.name }}</div>
                <div class="text-xs" :class="theme?.mutedText">{{ formatSize(item.size) }} · {{ formatDate(item.modified_at) }}</div>
              </div>
              <div class="flex items-center gap-2">
                <UButton size="xs" icon="i-heroicons-arrow-down-tray" color="gray" variant="soft" @click="downloadAttachment(item)">下载</UButton>
                <UButton size="xs" icon="i-heroicons-trash" color="red" variant="soft" @click="openDelete('image', item)">删除</UButton>
              </div>
            </div>
            <img :src="fullURL(item.url)" class="mt-2 rounded w-full max-h-40 object-contain bg-black/20" loading="lazy" />
            <div class="mt-2">
              <UButton size="xs" color="gray" variant="ghost" @click="toggleExpand(item)">{{ isExpanded(item) ? '收起关联内容' : '查看关联内容' }}</UButton>
              <div v-if="isExpanded(item)" class="mt-2 rounded p-2" :class="theme?.subtleBg">
                <div v-if="!item.belongs?.length" class="text-xs" :class="theme?.mutedText">无关联内容</div>
                <div v-else class="space-y-2">
                  <div v-for="b in item.belongs" :key="b.id" class="text-xs" :class="theme?.text">
                    <div class="flex items-center gap-2">
                      <span class="px-2 py-1 rounded text-[10px]" :class="theme?.subtleBg">#{{ b.id }}</span>
                      <span :class="theme?.mutedText">{{ formatDate(b.created_at) }}</span>
                    </div>
                    <div class="mt-1 line-clamp-2" :class="theme?.mutedText">{{ b.snippet }}</div>
                  </div>
                </div>
              </div>
            </div>
          </div>
          <div class="md:col-span-2 flex justify-center mt-2" v-if="images.length > imagesDisplay.length">
            <UButton color="gray" variant="soft" @click="loadMoreImages">加载更多</UButton>
          </div>
          <div class="md:col-span-2 flex justify-center mt-2" v-else-if="images.length > 4">
            <UButton color="gray" variant="soft" @click="collapseImages">收起</UButton>
          </div>
        </div>
      </div>
      <div v-else>
        <div v-if="videos.length===0" class="text-sm" :class="theme?.mutedText">暂无视频附件</div>
        <div v-else class="grid grid-cols-1 md:grid-cols-2 gap-3">
          <div v-for="item in videosDisplay" :key="item.name" class="rounded-lg border p-3" :class="theme?.border">
            <div class="flex items-center justify-between">
              <div>
                <div class="text-sm" :class="theme?.text">{{ item.name }}</div>
                <div class="text-xs" :class="theme?.mutedText">{{ formatSize(item.size) }} · {{ formatDate(item.modified_at) }}</div>
              </div>
              <div class="flex items-center gap-2">
                <UButton size="xs" icon="i-heroicons-arrow-down-tray" color="gray" variant="soft" @click="downloadAttachment(item)">下载</UButton>
                <UButton size="xs" icon="i-heroicons-trash" color="red" variant="soft" @click="openDelete('video', item)">删除</UButton>
              </div>
            </div>
            <video :src="fullURL(item.url)" class="mt-2 rounded w-full max-h-64 bg-black/20" controls preload="metadata"></video>
            <div class="mt-2">
              <UButton size="xs" color="gray" variant="ghost" @click="toggleExpand(item)">{{ isExpanded(item) ? '收起关联内容' : '查看关联内容' }}</UButton>
              <div v-if="isExpanded(item)" class="mt-2 rounded p-2" :class="theme?.subtleBg">
                <div v-if="!item.belongs?.length" class="text-xs" :class="theme?.mutedText">无关联内容</div>
                <div v-else class="space-y-2">
                  <div v-for="b in item.belongs" :key="b.id" class="text-xs" :class="theme?.text">
                    <div class="flex items-center gap-2">
                      <span class="px-2 py-1 rounded text-[10px]" :class="theme?.subtleBg">#{{ b.id }}</span>
                      <span :class="theme?.mutedText">{{ formatDate(b.created_at) }}</span>
                    </div>
                    <div class="mt-1 line-clamp-2" :class="theme?.mutedText">{{ b.snippet }}</div>
                  </div>
                </div>
              </div>
            </div>
          </div>
          <div class="md:col-span-2 flex justify-center mt-2" v-if="videos.length > videosDisplay.length">
            <UButton color="gray" variant="soft" @click="loadMoreVideos">加载更多</UButton>
          </div>
          <div class="md:col-span-2 flex justify-center mt-2" v-else-if="videos.length > 4">
            <UButton color="gray" variant="soft" @click="collapseVideos">收起</UButton>
          </div>
        </div>
      </div>
    </div>

    <UModal v-model="confirmOpen">
      <UCard :class="theme?.cardBg">
        <div class="text-sm" :class="theme?.text">确定删除该{{ deleteType==='image'?'图片':'视频' }}附件吗？此操作不可恢复。</div>
        <div class="flex justify-end gap-2 mt-3">
          <UButton color="gray" variant="soft" @click="confirmOpen=false">取消</UButton>
          <UButton color="red" :loading="deleting" @click="doDelete">确认删除</UButton>
        </div>
      </UCard>
    </UModal>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useToast } from '#imports'

const props = defineProps<{ theme?: Record<string, string> }>()

const activeTab = ref<'images'|'videos'>('images')
const images = ref<any[]>([])
const videos = ref<any[]>([])
const loading = ref(false)
const expanded = ref<Record<string, boolean>>({})
const confirmOpen = ref(false)
const deleting = ref(false)
const deleteType = ref<'image'|'video'>('image')
const deleteItem = ref<any>(null)
const imagesVisible = ref(4)
const videosVisible = ref(4)

const baseApi = useRuntimeConfig().public.baseApi || '/api'

const fullURL = (u: string) => {
  const origin = typeof window !== 'undefined' ? window.location.origin : ''
  if (u.startsWith('http')) return u
  if (u.startsWith('/api')) return origin + u
  return origin + u
}

const formatSize = (n: number) => {
  if (n < 1024) return `${n} B`
  if (n < 1024*1024) return `${(n/1024).toFixed(1)} KB`
  if (n < 1024*1024*1024) return `${(n/1024/1024).toFixed(1)} MB`
  return `${(n/1024/1024/1024).toFixed(1)} GB`
}
const formatDate = (s: string | Date) => {
  const d = new Date(s)
  return d.toLocaleString()
}

const isExpanded = (item: any) => !!expanded.value[item.name]
const toggleExpand = (item: any) => { expanded.value[item.name] = !expanded.value[item.name] }

const fetchImages = async () => {
  const resp = await fetch(`${baseApi}/attachments/images`, { credentials: 'include' })
  const js = await resp.json().catch(() => null)
  const arr = (js && js.code === 1 && Array.isArray(js.data)) ? js.data : []
  images.value = arr.filter((it: any) => /\.(png|jpe?g|gif|webp)$/i.test(String(it.name || '')))
  imagesVisible.value = 4
}
const fetchVideos = async () => {
  const resp = await fetch(`${baseApi}/attachments/video`, { credentials: 'include' })
  const js = await resp.json().catch(() => null)
  const arr = (js && js.code === 1 && Array.isArray(js.data)) ? js.data : []
  videos.value = arr.filter((it: any) => /\.(mp4|webm|mov|avi)$/i.test(String(it.name || '')))
  videosVisible.value = 4
}
const refresh = async () => {
  try {
    loading.value = true
    await Promise.all([fetchImages(), fetchVideos()])
  } finally {
    loading.value = false
  }
}
onMounted(refresh)

const openDelete = (type: 'image'|'video', item: any) => {
  deleteType.value = type
  deleteItem.value = item
  confirmOpen.value = true
}
const doDelete = async () => {
  if (!deleteItem.value) return
  try {
    deleting.value = true
    const url = `${baseApi}/attachments/${deleteType.value === 'image' ? 'images' : 'video'}/${encodeURIComponent(deleteItem.value.name)}`
    const resp = await fetch(url, { method: 'DELETE', credentials: 'include' })
    const js = await resp.json().catch(() => null)
    if (!resp.ok || !js || js.code !== 1) throw new Error(js?.msg || '删除失败')
    useToast().add({ title: '已删除', color: 'green' })
    confirmOpen.value = false
    deleteItem.value = null
    await refresh()
  } catch (e: any) {
    useToast().add({ title: '删除失败', description: e?.message, color: 'red' })
  } finally {
    deleting.value = false
  }
}

const imagesDisplay = computed(() => images.value.slice(0, imagesVisible.value))
const videosDisplay = computed(() => videos.value.slice(0, videosVisible.value))
const loadMoreImages = () => { imagesVisible.value += 4 }
const loadMoreVideos = () => { videosVisible.value += 4 }
const collapseImages = () => { imagesVisible.value = 4 }
const collapseVideos = () => { videosVisible.value = 4 }
const downloadAttachment = (item: any) => {
  const a = document.createElement('a')
  a.href = fullURL(item.url)
  a.download = item.name || ''
  document.body.appendChild(a)
  a.click()
  document.body.removeChild(a)
}

const theme = computed(() => props.theme || {})
</script>

<style scoped>
</style>
