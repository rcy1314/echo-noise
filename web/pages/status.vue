<template>
  <StatusPanel />
  <Notification />
</template>

<script setup lang="ts">
definePageMeta({ middleware: ['auth'] })
import StatusPanel from '@/components/index/StatusPanel.vue'
import Notification from '~/components/widgets/Notification.vue'

// 隔离首页的暗黑/明亮模式对后台页的影响：
// 进入后台页时根据管理员选择的 adminTheme 强制设置 html 上的配色类；
// 离开后台页时恢复进入前的 html 配色状态。
const prevHtmlClasses = typeof document !== 'undefined' ? Array.from(document.documentElement.classList) : []
const applyAdminColorMode = () => {
  if (typeof window === 'undefined') return
  const theme = (localStorage.getItem('adminTheme') || 'dark').toLowerCase()
  const html = document.documentElement
  // 清理影响 UI 组件的暗色类
  html.classList.remove('dark')
  // 根据 adminTheme 选择是否启用暗色基础（仅当选择 light 时不启用）
  if (theme !== 'light') {
    html.classList.add('dark')
  }
}
onMounted(applyAdminColorMode)
onBeforeUnmount(() => {
  if (typeof window === 'undefined') return
  const html = document.documentElement
  // 先清空再恢复原先的类集
  Array.from(html.classList).forEach(cls => html.classList.remove(cls))
  prevHtmlClasses.forEach(cls => html.classList.add(cls))
})
</script>
