<template>
  <UCard class="rounded-xl">
    <div class="px-3 py-2 text-xs opacity-70">主题预设</div>
    <div class="flex items-center gap-2 p-2">
      <button v-for="p in presets" :key="p.key" class="w-6 h-6 rounded-full border hover:scale-105" :style="{ background: p.color }" @click="apply(p.key)" :title="p.name"></button>
    </div>
  </UCard>
</template>

<script setup lang="ts">
type Preset = { key: string; name: string; color: string }
const presets: Preset[] = [
  { key: 'default', name: '默认', color: '#ff8c3a' },
  { key: 'mint', name: '薄荷', color: '#3bb273' },
  { key: 'rose', name: '玫瑰', color: '#e85d75' },
  { key: 'slate', name: '蓝灰', color: '#5c7cfa' },
]
const apply = (key: string) => {
  if (typeof window !== 'undefined') localStorage.setItem('themePreset', key)
  const root = document.documentElement
  root.classList.remove('theme-default','theme-mint','theme-rose','theme-slate')
  root.classList.add(`theme-${key}`)
  // 清理手动设置的 accent 以使用预设主题的全局颜色
  root.style.removeProperty('--accent')
  root.style.removeProperty('--title-color')
}
onMounted(() => {
  const saved = typeof window !== 'undefined' ? localStorage.getItem('themePreset') : null
  if (saved) {
    const root = document.documentElement
    root.classList.add(`theme-${saved}`)
  }
})
</script>

<style scoped>
.rounded-full { transition: transform .15s ease }
</style>
