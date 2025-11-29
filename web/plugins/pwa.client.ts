export default defineNuxtPlugin(() => {
  const applyHeadFromConfig = (fs: any) => {
    const enabled = typeof fs.pwaEnabled === 'boolean' ? fs.pwaEnabled : true
    const title = (fs.pwaTitle || fs.siteTitle || '说说笔记').trim()
    const icon = (fs.pwaIconURL || '/favicon.svg')
    const description = (fs.pwaDescription || fs.description || '').trim()

  try {
    // 在 Nuxt 插件中使用 useHead 更新全局 Head
    // 当禁用 PWA 时，不注入 manifest，并移除已存在的 DOM 节点
    if (enabled && useHead) {
      useHead({
        title,
        meta: [
          { name: 'description', content: description },
          { name: 'theme-color', content: '#000000' }
        ],
        link: [
          { rel: 'icon', href: icon },
          { rel: 'apple-touch-icon', href: icon },
          { rel: 'manifest', href: '/manifest.webmanifest' }
        ]
      })
    } else {
      // 移除可能存在的 manifest 链接，避免首页仍显示 PWA
      const manifestEl = document.querySelector('link[rel="manifest"]')
      if (manifestEl) manifestEl.parentElement?.removeChild(manifestEl)
    }
    } catch {}

    // 同步 Service Worker 状态
    if ('serviceWorker' in navigator) {
      if (enabled) {
        navigator.serviceWorker.register('/sw.js')
      } else {
        navigator.serviceWorker.getRegistrations().then(async regs => {
          for (const r of regs) await r.unregister()
          const keys = await caches.keys()
          await Promise.all(keys.map(k => caches.delete(k)))
        })
      }
    }
  }

  const loadAndApply = async () => {
    try {
      const res = await fetch('/api/frontend/config', { credentials: 'include' })
      const data = await res.json()
      const fs = data?.data?.frontendSettings || {}
      applyHeadFromConfig(fs)
      try {
        const last = localStorage.getItem('pwaEnabledLast')
        const prev = last === null ? null : (last === 'true')
        const curr = typeof fs.pwaEnabled === 'boolean' ? fs.pwaEnabled : true
        if (prev !== curr) {
          localStorage.setItem('pwaEnabledLast', curr ? 'true' : 'false')
          if (!curr && 'serviceWorker' in navigator) {
            const regs = await navigator.serviceWorker.getRegistrations()
            for (const r of regs) await r.unregister()
            const keys = await caches.keys()
            await Promise.all(keys.map(k => caches.delete(k)))
          }
        }
      } catch {}
    } catch {}
  }

  // 初次加载
  window.addEventListener('load', loadAndApply)
  // 监听后台面板触发的配置更新事件
  window.addEventListener('frontend-config-updated', loadAndApply)
})
import { useHead } from '#imports'
