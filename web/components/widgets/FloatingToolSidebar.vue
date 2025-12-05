<template>
  <div class="floating-sidebar" :class="isDark ? 'fs-dark' : 'fs-light'">
    <button class="tool-btn btn-layout" :class="isDark ? 'btn-dark' : 'btn-light'" @click="$emit('toggle-layout')" :title="'布局'">
      <UIcon :name="layoutIconProp" class="w-5 h-5" />
    </button>
    <button class="tool-btn" :class="isDark ? 'btn-dark' : 'btn-light'" @click="$emit('search')" :title="'搜索'">
      <UIcon name="i-heroicons-magnifying-glass" class="w-5 h-5" />
    </button>
    <button class="tool-btn" :class="isDark ? 'btn-dark' : 'btn-light'" @click="$emit('switch-background')" :title="'背景'">
      <UIcon name="i-mdi-image-outline" class="w-5 h-5" />
    </button>
    <button class="tool-btn" :class="isDark ? 'btn-dark' : 'btn-light'" @click="$emit('toggle-theme')" :title="'切换亮暗'">
      <UIcon :name="themeIcon" class="w-5 h-5" />
    </button>
    <a href="/rss" target="_blank" rel="noopener noreferrer" class="tool-btn" :class="isDark ? 'btn-dark' : 'btn-light'" title="RSS">
      <UIcon name="i-mdi-rss" class="w-5 h-5" />
    </a>
    <button class="tool-btn" :class="isDark ? 'btn-dark' : 'btn-light'" title="后台" @click="$emit('open-admin')">
      <UIcon name="i-mdi-server-outline" class="w-5 h-5" />
    </button>
  </div>
</template>

<script setup lang="ts">
const props = defineProps<{ contentTheme?: string; layoutIcon?: string }>()
const isDark = computed(() => props.contentTheme === 'dark')
const themeIcon = computed(() => (props.contentTheme === 'dark' ? 'i-mdi-weather-night' : 'i-mdi-white-balance-sunny'))
const layoutIconProp = computed(() => props.layoutIcon || 'i-mdi-view-grid')
</script>

<style scoped>
.floating-sidebar { position: fixed; right: 16px; top: 50%; transform: translateY(-50%); z-index: 50; display:flex; flex-direction:column; gap:10px; padding:8px; border-radius:12px; background: transparent; box-shadow: none; }
.floating-sidebar.fs-dark { background: transparent !important; }
.floating-sidebar.fs-light { background: transparent !important; box-shadow: none; }
.tool-btn { display:flex; align-items:center; justify-content:center; width:40px; height:40px; border-radius:10px; transition: all .18s ease; }
.tool-btn.btn-dark { background: rgba(36,43,50,0.85); color:#ffffff; border: 1px solid rgba(255,255,255,0.20); box-shadow: 0 6px 16px rgba(0,0,0,0.35); backdrop-filter: blur(6px); }
.tool-btn.btn-dark:hover { transform: translate3d(0,0,0) scale(1.06); background: rgba(36,43,50,0.95); border-color: rgba(255,255,255,0.28); }
.tool-btn.btn-light { background: rgba(255,255,255,0.92); color:#1f2937; border: 1px solid rgba(0,0,0,0.18); box-shadow: 0 2px 8px rgba(0,0,0,.12); }
.tool-btn.btn-light:hover { transform: translate3d(0,0,0) scale(1.06); background: #ffffff; border-color: rgba(0,0,0,0.24); box-shadow: 0 4px 12px rgba(0,0,0,.18); }
.tool-btn.btn-light:hover { transform: translate3d(0,0,0) scale(1.06); background: rgba(255,255,255,0.70); }
@media (max-width: 1024px) { .floating-sidebar { position: fixed; left: 50%; bottom: 20px; top: auto; transform: translateX(-50%); flex-direction:row; padding:8px; border-radius:14px; background: transparent; box-shadow: none; } .tool-btn { width:36px; height:36px; } .tool-btn.btn-layout { display: none; } }
</style>
