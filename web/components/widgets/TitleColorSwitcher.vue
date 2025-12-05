<template>
  <div class="px-3 py-2 text-xs opacity-70">标题颜色</div>
  <div class="flex items-center gap-2 p-2">
    <button v-for="c in colors" :key="c" class="w-6 h-6 rounded-full hover:scale-105" :style="{ background: c }" @click="apply(c)"></button>
  </div>
</template>

<script setup lang="ts">
const colors = ['#ff8c3a', '#f59e0b', '#22c55e', '#3b82f6', '#ef4444', '#e85d75', '#5c7cfa']
const apply = (c: string) => {
  document.documentElement.style.setProperty('--accent', c)
  document.documentElement.style.setProperty('--title-color', c)
  if (typeof window !== 'undefined') localStorage.setItem('accentColor', c)
}
onMounted(() => {
  const saved = typeof window !== 'undefined' ? (localStorage.getItem('accentColor') || localStorage.getItem('titleColor')) : null
  if (saved) {
    document.documentElement.style.setProperty('--accent', saved)
    document.documentElement.style.setProperty('--title-color', saved)
  }
})
</script>

<style scoped>
.rounded-full { transition: transform .15s ease }
</style>
