<template>
  <div class="background-container" :style="backgroundStyle">
    <div class="loading" v-if="!isLoaded">
      <div class="rainbow-spinner"></div>
      <div class="loading-text">加载中...</div>
    </div>
    <div ref="contentWrapper" class="content-wrapper" :class="{ 'gpu-accelerated': true }">
      <UContainer class="container-fixed py-2 pb-4 my-4">
        <div class="moments-header">
          <div class="header-image" :style="headerImageStyle">
    <h1 class="header-title">{{ frontendConfig.siteTitle }}</h1>
    <div ref="subtitleEl" class="header-subtitle">{{ frontendConfig.subtitleText }}</div>
    
    <div class="profile-info">
        <img class="avatar" 
             @click="changeBackground" 
             :src="frontendConfig.avatarURL" 
             :alt="frontendConfig.username">
        <div class="profile-text">
            <div class="title">{{ frontendConfig.username }}</div>
            <div class="description">{{ frontendConfig.description }}</div>
        </div>
    </div>
          </div>
        </div>
       <!-- 移动热力图到这里，在头部主图下方，AddForm上方 -->
    <div v-if="showHeatmap" class="mx-auto sm:max-w-2xl mb-4">
      <HeatmapWidget />
    </div>
    
    <!-- 添加表单 -->
    <AddForm @search-result="handleSearchResult" />
    <!-- 移动标签墙到这里，并保持与 AddForm 相同宽度 -->
    <TagList 
      v-if="tags && tags.length > 0"
      :tags="tags"
      @tagClick="handleTagClick"
      @updateTags="handleTagsUpdate" 
    />
    <AnnouncementBar v-if="frontendConfig.announcementEnabled && (frontendConfig.announcementText || '').trim() !== ''" :text="frontendConfig.announcementText || '欢迎访问我的说说笔记！'" />
        <!-- 确保 MessageList 有足够的 z-index -->
      <MessageList 
      ref="messageList" 
      class="message-list-container" 
      :site-config="frontendConfig"
      :target-message-id="targetMessageId" 
      />
      <!-- 音乐播放器容器（浮动或嵌入） -->
      <div
        v-if="frontendConfig.musicEnabled"
        class="netease-mini-player"
        :data-playlist-id="frontendConfig.musicPlaylistId || ''"
        :data-song-id="frontendConfig.musicSongId || ''"
        :data-position="frontendConfig.musicPosition || 'bottom-left'"
        :data-theme="frontendConfig.musicTheme || 'auto'"
        :data-lyric="frontendConfig.musicLyric ? 'true' : 'false'"
        :data-default-minimized="frontendConfig.musicDefaultMinimized ? 'true' : 'false'"
        :data-embed="frontendConfig.musicEmbed ? 'true' : 'false'"
        :data-autoplay="frontendConfig.musicAutoplay ? 'true' : 'false'"
      />
      </UContainer>
  <Notification />
  <!-- 添加搜索模态框组件 -->
  <SearchMode v-model="showSearchModal" @search-result="handleSearchResult" />
  <div class="scroll-buttons" @mouseenter="hoverScroll = true" @mouseleave="hoverScroll = false">
    <UButton v-show="showScroll" :class="scrollButtonClass" variant="ghost" size="sm" @click="handleScrollClick">
      <UIcon :class="iconClass" :name="scrollIconName" />
    </UButton>
  </div>
  </div>
</div>
</template>

<script setup lang="ts">
import { ref, computed, inject, provide, onMounted, onUnmounted, watch } from 'vue'
import { useRoute } from '#imports'
import AddForm from '@/components/index/AddForm.vue'
import MessageList from '@/components/index/MessageList.vue'
import Notification from '~/components/widgets/Notification.vue';
import HeatmapWidget from '~/components/widgets/heatmap.vue'
import SearchMode from '~/components/index/Searchmode.vue' // 导入 SearchMode 组件
import TagList from '~/components/index/TagList.vue'
import AnnouncementBar from '~/components/widgets/AnnouncementBar.vue'
import { getRequest } from '~/utils/api'

// 添加 messageList ref
const messageList = ref(null)
// 修复：定义 targetMessageId，避免模板引用未定义导致列表不渲染
const targetMessageId = ref<string | null>(null)
// 添加搜索结果处理函数
const handleSearchResult = (result: any) => {
  console.log('接收到搜索结果:', result); // 添加调试日志
  if (messageList.value) {
    // 直接传递原始结果，让 MessageList 组件自己处理数据格式
    messageList.value.handleSearchResult(result);
  }
}
// 注入从AddForm组件提供的showHeatmap变量
const showHeatmap = ref(false);
const contentTheme = ref<string>(typeof window !== 'undefined' ? (localStorage.getItem('contentTheme') || 'light') : 'light')
// 提供给子组件
provide('showHeatmap', showHeatmap);
provide('contentTheme', contentTheme)
provide('toggleContentTheme', () => {
  contentTheme.value = contentTheme.value === 'dark' ? 'light' : 'dark'
  if (typeof window !== 'undefined') {
    localStorage.setItem('contentTheme', contentTheme.value)
    document.documentElement.className = contentTheme.value === 'dark' ? 'dark' : ''
  }
})

const contentWrapper = ref<HTMLElement | null>(null)
const scrollToTop = () => {
  const el = contentWrapper.value
  if (el) el.scrollTo({ top: 0, behavior: 'smooth' })
  else window.scrollTo({ top: 0, behavior: 'smooth' })
}
const scrollToBottom = () => {
  const el = contentWrapper.value
  if (el) el.scrollTo({ top: el.scrollHeight, behavior: 'smooth' })
  else {
    const h = Math.max(document.body.scrollHeight, document.documentElement.scrollHeight)
    window.scrollTo({ top: h, behavior: 'smooth' })
  }
}

// 滚动状态与按钮展示逻辑
const hoverScroll = ref(false)
const isAtTop = ref(true)
const isAtBottom = ref(false)
const updateScrollState = () => {
  const el = contentWrapper.value
  if (!el) {
    const y = window.scrollY || document.documentElement.scrollTop || 0
    const max = Math.max(document.body.scrollHeight, document.documentElement.scrollHeight)
    isAtTop.value = y <= 2
    isAtBottom.value = window.innerHeight + y >= max - 2
    return
  }
  const y = el.scrollTop
  const max = el.scrollHeight
  isAtTop.value = y <= 2
  isAtBottom.value = el.clientHeight + y >= max - 2
}
onMounted(() => {
  updateScrollState()
  contentWrapper.value?.addEventListener('scroll', updateScrollState, { passive: true })
})
onUnmounted(() => {
  contentWrapper.value?.removeEventListener('scroll', updateScrollState)
})

const showScroll = computed(() => isAtTop.value || isAtBottom.value || hoverScroll.value)
const scrollIconName = computed(() => (isAtBottom.value && !isAtTop.value) ? 'i-heroicons-arrow-up' : 'i-heroicons-arrow-down')
const handleScrollClick = () => {
  if (isAtBottom.value && !isAtTop.value) {
    scrollToTop()
  } else {
    scrollToBottom()
  }
}
const isDark = computed(() => contentTheme.value === 'dark')
const scrollButtonClass = computed(() => (
  isDark.value
    ? 'scroll-button bg-[rgba(36,43,50,0.85)] hover:bg-[rgba(36,43,50,0.95)] text-white shadow-[0_6px_16px_rgba(0,0,0,0.35)]'
    : 'scroll-button bg-white/95 hover:bg-white text-gray-700 ring-1 ring-gray-200 shadow-[0_4px_12px_rgba(0,0,0,0.12)]'
))
const iconClass = computed(() => (isDark.value ? 'text-white w-6 h-6' : 'text-gray-600 w-6 h-6'))

// 添加监听，查看状态变化
watch(showHeatmap, (newVal) => {
  console.log('index.vue 中热力图状态变化:', newVal);
});
// 添加 useHead
useHead({
  meta: [
    {
      name: 'viewport',
      content: 'width=device-width, initial-scale=1.0, maximum-scale=1.0, user-scalable=no'
    }
  ]
})
// 同步路由中的消息ID到 MessageList，用于高亮或定位
const route = useRoute()
watch(() => route.hash, (newHash) => {
  const id = (newHash || '').split('/messages/').pop()
  targetMessageId.value = id || null
}, { immediate: true })

// 添加前端配置的响应式对象
const frontendConfig = ref({
    siteTitle: '',
    subtitleText: '',
    avatarURL: '',
    username: '',
    description: '',
    backgrounds: [] as string[],
    cardFooterTitle: '',
    cardFooterSubtitle: '',
    pageFooterHTML: '',
    rssTitle: '',
    rssDescription: '',
    rssAuthorName: '',
    rssFaviconURL: '',
    enableGithubCard: false,
    // PWA
    pwaEnabled: true,
    pwaTitle: '',
    pwaDescription: '',
    pwaIconURL: '',
    announcementText: '',
    announcementEnabled: true,
    // 评论系统
    commentEnabled: false,
    commentSystem: 'waline',
    commentEmailEnabled: false,
    walineServerURL: '',
    // 音乐配置
    musicEnabled: false,
    musicPlaylistId: '2141128031',
    musicSongId: '',
    musicPosition: 'bottom-left',
    musicTheme: 'auto',
    musicLyric: true,
    musicAutoplay: false,
    musicDefaultMinimized: true,
    musicEmbed: false
})

const backgroundStyle = computed(() => ({
    '--bg-image': `url(${currentImage.value || frontendConfig.value.backgrounds[0]})`
}))
// 添加 headerImageStyle 计算属性
const headerImageStyle = computed(() => ({
    'background-image': `url(${currentImage.value || frontendConfig.value.backgrounds[0]})`,
    'background-size': 'cover',
    'background-position': 'center'
}))
// 修改 fetchConfig 方法

// 首先添加默认配置对象
const defaultConfig = {
    siteTitle: 'Noise的说说笔记',
    subtitleText: '欢迎访问，点击头像可更换封面背景！',
    avatarURL: 'https://s2.loli.net/2025/03/24/HnSXKvibAQlosIW.png',
    username: 'Noise',
    description: '执迷不悟',
    backgrounds: [
                'https://s2.loli.net/2025/03/27/KJ1trnU2ksbFEYM.jpg',
                'https://s2.loli.net/2025/03/27/MZqaLczCvwjSmW7.jpg',
                'https://s2.loli.net/2025/03/27/UMijKXwJ9yTqSeE.jpg',
                'https://s2.loli.net/2025/03/27/WJQIlkXvBg2afcR.jpg',
                'https://s2.loli.net/2025/03/27/oHNQtf4spkq2iln.jpg',
                'https://s2.loli.net/2025/03/27/PMRuX5loc6Uaimw.jpg',
                'https://s2.loli.net/2025/03/27/U2WIslbNyTLt4rD.jpg',
                'https://s2.loli.net/2025/03/27/xu1jZL5Og4pqT9d.jpg',
                'https://s2.loli.net/2025/03/27/OXqwzZ6v3PVIns9.jpg',
                'https://s2.loli.net/2025/03/27/HGuqlE6apgNywbh.jpg',
                'https://s2.loli.net/2025/03/26/d7iyuPYA8cRqD1K.jpg',
                'https://s2.loli.net/2025/03/27/wYy12qDMH6bGJOI.jpg',
                'https://s2.loli.net/2025/03/27/y67m2k5xcSdTsHN.jpg',
    ],
    cardFooterTitle: 'Noise·说说·笔记~',
    cardFooterSubtitle: 'note.noisework.cn',
    pageFooterHTML: '<div class="text-center text-xs text-gray-400 py-4">来自<a href="https://www.noisework.cn" target="_blank" rel="noopener noreferrer" class="text-orange-400 hover:text-orange-500">Noise</a> 使用<a href="https://github.com/rcy1314/echo-noise" target="_blank" rel="noopener noreferrer" class="text-orange-400 hover:text-orange-500">Ech0-Noise</a>发布</div>',
    rssTitle: 'Noise的说说笔记',
    rssDescription: '一个说说笔记~',
    rssAuthorName: 'Noise',
    rssFaviconURL: '/favicon.ico',
    enableGithubCard: false,
    // PWA 默认（为空时回退到站点设置）
    pwaEnabled: true,
    pwaTitle: '',
    pwaDescription: '',
    pwaIconURL: '',
    announcementText: '欢迎访问我的说说笔记！',
    announcementEnabled: true,
    // 评论系统默认值
    commentEnabled: false,
    commentSystem: 'waline',
    commentEmailEnabled: false,
    walineServerURL: '',
    // 音乐默认配置
    musicEnabled: false,
    musicPlaylistId: '2141128031',
    musicSongId: '',
    musicPosition: 'bottom-left',
    musicTheme: 'auto',
    musicLyric: true,
    musicAutoplay: false,
    musicDefaultMinimized: true,
    musicEmbed: false
};

// 修改 fetchConfig 方法
const fetchConfig = async () => {
    try {
        frontendConfig.value = { ...defaultConfig };
        const res = await getRequest<any>('frontend/config', undefined, { credentials: 'include' })
        if (res && res.code === 1 && res.data && res.data.frontendSettings) {
            const settings = res.data.frontendSettings
            const booleanKeys = ['enableGithubCard', 'pwaEnabled', 'announcementEnabled', 'commentEnabled', 'commentEmailEnabled', 'musicEnabled', 'musicLyric', 'musicAutoplay', 'musicDefaultMinimized', 'musicEmbed']
            Object.keys(frontendConfig.value).forEach(key => {
                if (settings[key] !== null && settings[key] !== undefined) {
                    if (key === 'backgrounds' && Array.isArray(settings[key])) {
                        frontendConfig.value.backgrounds = [...settings[key]]
                    } else if (booleanKeys.includes(key)) {
                        const v = settings[key]
                        frontendConfig.value[key] = (v === true || v === 'true')
                    } else {
                        const v = settings[key]
                        frontendConfig.value[key] = typeof v === 'string' ? v.trim() : v
                    }
                }
            })
            const defaultTheme = (settings.defaultContentTheme || 'light').trim()
            if (typeof window !== 'undefined' && !localStorage.getItem('contentTheme')) {
              contentTheme.value = defaultTheme === 'light' ? 'light' : 'dark'
              document.documentElement.className = contentTheme.value === 'dark' ? 'dark' : ''
            } else if (typeof window !== 'undefined') {
              document.documentElement.className = contentTheme.value === 'dark' ? 'dark' : ''
            }
        }
        if (!frontendConfig.value.backgrounds?.length) {
            frontendConfig.value.backgrounds = [...defaultConfig.backgrounds]
        }
        if (frontendConfig.value.backgrounds.length > 0) {
            const randomIndex = Math.floor(Math.random() * frontendConfig.value.backgrounds.length)
            currentImage.value = frontendConfig.value.backgrounds[randomIndex]
        }
        isLoaded.value = true
    } catch (error) {
        console.error('获取配置失败:', error)
        frontendConfig.value = { ...defaultConfig }
        isLoaded.value = true
    }
}
const currentImage = ref('')
const isLoaded = ref(false)
const imageLoading = ref(false)
// 添加图片预加载函数
const preloadImages = async (images: string[]) => {
  const loadImage = (src: string) => {
    return new Promise((resolve) => {
      const img = new Image()
      img.src = src
      img.onload = () => resolve(src)
      img.onerror = () => resolve(null)
    })
  }
  
  // 并行预加载所有图片
  await Promise.all(images.map(src => loadImage(src)))
}
// 添加配置更新事件监听
// 移除重复绑定的 frontend-config-updated 监听，避免多次拉取导致卡顿

onUnmounted(() => {
    // 移除事件监听
    window.removeEventListener('frontend-config-updated', fetchConfig);
});
// 优化背景切换函数
const changeBackground = async () => {
  if (imageLoading.value) return
  imageLoading.value = true
  
  const newIndex = Math.floor(Math.random() * frontendConfig.value.backgrounds.length)
  const newImage = frontendConfig.value.backgrounds[newIndex]
  
  if (newImage === currentImage.value) {
    imageLoading.value = false
    return
  }

  // 使用更小的缩略图
  const thumbnailImage = `${newImage}?imageView2/2/w/10/blur/1/q/10`
  currentImage.value = thumbnailImage

  // 使用 requestAnimationFrame 优化渲染
  requestAnimationFrame(() => {
    const img = new Image()
    img.src = newImage
    img.onload = () => {
      requestAnimationFrame(() => {
        currentImage.value = newImage
        imageLoading.value = false
      })
    }
    img.onerror = () => {
      imageLoading.value = false
    }
  })
}
// 定义页面元数据
definePageMeta({
  title: '说说笔记'
})

// 设置动态标题
let headUpdateTimer: any = null
const scheduleHeadUpdate = () => {
  if (headUpdateTimer) clearTimeout(headUpdateTimer)
  headUpdateTimer = setTimeout(() => updateTitle(), 100)
}
const updateTitle = () => {
  const title = (frontendConfig.value.pwaTitle || frontendConfig.value.siteTitle || '说说笔记').trim()
  const icon = (frontendConfig.value.rssFaviconURL || '/favicon.ico').trim()
  const pwaIcon = (
    frontendConfig.value.pwaIconURL && frontendConfig.value.pwaIconURL.trim() !== ''
      ? frontendConfig.value.pwaIconURL.trim()
      : (icon.toLowerCase().endsWith('.png') ? icon : '/android-chrome-192x192.png')
  )
  const description = (frontendConfig.value.pwaDescription || frontendConfig.value.description || '').trim()
  useHead({
    title,
    meta: [
      { key: 'description', name: 'description', content: description },
      { key: 'theme-color', name: 'theme-color', content: '#000000' }
    ],
    link: [
      { key: 'icon-32', rel: 'icon', type: 'image/png', href: '/favicon-32x32.png', sizes: '32x32' },
      { key: 'shortcut-icon-32', rel: 'shortcut icon', type: 'image/png', href: '/favicon-32x32.png', sizes: '32x32' },
      { key: 'icon-fallback', rel: 'icon', href: icon },
      ...(frontendConfig.value.pwaEnabled ? [
        { key: 'manifest', rel: 'manifest', href: '/manifest.webmanifest' },
        { key: 'apple-touch', rel: 'apple-touch-icon', href: pwaIcon, sizes: '180x180' },
        { key: 'pwa-192', rel: 'icon', href: pwaIcon, sizes: '192x192' },
        { key: 'pwa-512', rel: 'icon', href: (pwaIcon.toLowerCase().endsWith('.png') ? pwaIcon : '/android-chrome-512x512.png'), sizes: '512x512' }
      ] : [])
    ]
  })
}

// 监听配置变化
watch(() => [frontendConfig.value.pwaEnabled, frontendConfig.value.pwaTitle, frontendConfig.value.pwaIconURL, frontendConfig.value.pwaDescription, frontendConfig.value.siteTitle, frontendConfig.value.rssFaviconURL, frontendConfig.value.description], () => {
  scheduleHeadUpdate()
}, { immediate: true })
const subtitleEl = ref<HTMLElement | null>(null)
  const tags = ref([])
// 添加标签更新处理函数
const handleTagsUpdate = async () => {
  await fetchTags()
}
// 获取所有标签
const fetchTags = async () => {
  try {
    const res = await getRequest<any>('messages/tags')
    if (res && res.code === 1) {
      tags.value = res.data || []
    } else {
      tags.value = []
    }
  } catch (error) {
    console.error('获取标签失败:', error)
    tags.value = []
  }
}

// 监听前端配置更新事件，保存后主动刷新配置
onMounted(() => {
  const handler = () => fetchConfig()
  window.addEventListener('frontend-config-updated', handler)
  // 初始拉取
  fetchConfig()
  onUnmounted(() => window.removeEventListener('frontend-config-updated', handler))
})
// 标签点击处理
const handleTagClick = async (tag: string) => {
  try {
    const encodedTag = encodeURIComponent(tag.trim())
    const res = await getRequest<any>(`messages/tags/${encodedTag}`, undefined, { credentials: 'include' })
    if (res && res.code === 1 && Array.isArray(res.data)) {
      if (messageList.value) {
        messageList.value.handleSearchResult(res.data)
      }
    } else {
      throw new Error(res?.msg || '获取标签内容失败')
    }
  } catch (error: any) {
    console.error('获取标签消息失败:', error)
    useToast().add({
      title: '获取标签消息失败',
      description: error.message || '服务器错误，请稍后重试',
      color: 'red',
      timeout: 3000
    })
  }
}
// 修改打字效果函数
const startTypeEffect = () => {
  if (!subtitleEl.value) return
  
  let index = 0
  let isDeleting = false
  let isWaiting = false
  
  const typeInterval = setInterval(() => {
    if (!subtitleEl.value) {
      clearInterval(typeInterval)
      return
    }
    if (isWaiting) return

    if (!isDeleting) {
      // 打字过程
      subtitleEl.value!.textContent = frontendConfig.value.subtitleText.slice(0, index + 1)
      index++
      
      if (index >= frontendConfig.value.subtitleText.length) {
        isWaiting = true
        setTimeout(() => {
          isDeleting = true
          isWaiting = false
        }, 2000)
      }
    } else {
      // 删除过程
      index--
      subtitleEl.value!.textContent = frontendConfig.value.subtitleText.slice(0, index)
      
      if (index <= 0) {
        isWaiting = true
        subtitleEl.value!.textContent = ''
        setTimeout(() => {
          isDeleting = false
          isWaiting = false
          index = 0
        }, 1000)
      }
    }
  }, 100)

  return typeInterval
}

// 修改 onMounted 钩子
onMounted(async () => {
  try {
    // 确保在任何异步操作之前设置加载状态
    isLoaded.value = false;

    // 使用 requestIdleCallback 延迟加载非关键组件
    window.requestIdleCallback = window.requestIdleCallback || ((cb) => setTimeout(cb, 1))
    
    // 关键内容优先加载
    await fetchConfig()
    
    // 非关键内容延迟加载
    requestIdleCallback(async () => {
      await fetchTags()
      await preloadImages(frontendConfig.value.backgrounds)
    })

    // 确保配置加载完成后再执行后续操作
    if (frontendConfig.value.backgrounds.length > 0) {
      const initialImage = frontendConfig.value.backgrounds[
        Math.floor(Math.random() * frontendConfig.value.backgrounds.length)
      ]
      
      // 先加载低质量版本
      const lowQualityImage = `${initialImage}?imageView2/2/w/100/q/30`
      currentImage.value = lowQualityImage

      // 后台预加载其他图片
      requestIdleCallback(async () => {
        await preloadImages(frontendConfig.value.backgrounds)
      })
      
      // 加载高质量初始图片
      const img = new Image()
      img.src = initialImage
      img.onload = () => {
        requestAnimationFrame(() => {
          currentImage.value = initialImage
          isLoaded.value = true; // 在高质量图片加载完成后设置为已加载
        })
      }
      img.onerror = () => {
        // 若加载失败，直接结束加载遮罩，避免阻塞交互
        isLoaded.value = true
      }
      // 最长等待时间到达后也结束加载遮罩，避免卡顿
      setTimeout(() => {
        if (!isLoaded.value) isLoaded.value = true
      }, 2000)
    }
    
    // 启动打字效果
    const typeInterval = startTypeEffect()
    onUnmounted(() => {
      if (typeInterval) {
        clearInterval(typeInterval)
      }
    })

    // 添加事件监听
    window.addEventListener('frontend-config-updated', async (event: CustomEvent) => {
      await fetchConfig()
      if (frontendConfig.value.backgrounds?.length > 0) {
        const randomIndex = Math.floor(Math.random() * frontendConfig.value.backgrounds.length)
        const newImage = frontendConfig.value.backgrounds[randomIndex]
        currentImage.value = newImage
      }
    })

  } catch (error) {
    console.error('初始化失败:', error)
    isLoaded.value = true
  }
})
</script>

<style>
html, body {
  margin: 0;
  padding: 0;
  width: 100%;
  height: 100%;
  overflow: hidden;
  overscroll-behavior: none;
}
.header-subtitle {
  position: absolute;
  top: calc(50% + 60px);  /* 增加间距 */
  left: 50%;
  transform: translate(-50%, -50%);
  color: white;
  font-size: 1rem;
  text-shadow: none;
  white-space: nowrap;
}

@media screen and (max-width: 768px) {
  .header-subtitle {
    font-size: 0.9rem;
    top: calc(50% + 45px);  /* 保持移动端的适配 */
  }
}
.background-container {
  width: 100%;
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  overflow: hidden;
  z-index: 0;
  background-color: black;
}

.background-container::before {
  content: '';
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background-image: var(--bg-image);
  background-size: cover;
  background-position: center;
  background-attachment: fixed; /* 确保背景固定 */
  filter: blur(8px);
  z-index: -1;
}

.content-wrapper {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  height: 100vh;
  overflow-y: auto;
  z-index: 1;
  pointer-events: auto;
  cursor: default;
}

.moments-header {
  margin-bottom: 20px;
}

.header-image {
  position: relative;
  width: 100%;
  height: 300px;
  background-size: cover;
  background-position: center;
  border-radius: 18px;
  overflow: hidden;
  transition: background-image 0.3s ease;
  will-change: background-image;
  transform: translateZ(0);
  margin-top: 1px; /* 调整这个值来控制主图下移距离 */
}

.header-title {
  position: absolute;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  color: white;
  font-size: 2.5rem;
  font-weight: bold;
  text-shadow: none;
  margin: 0;
  white-space: nowrap;
  transition: font-size 0.3s ease;
}
@media screen and (max-width: 768px) {
  .content-wrapper {
   padding: 0.5rem; 
  }
  /* 优化移动端滚动性能 */
  .message-list-container {
    transform: translateZ(0);
    will-change: transform;
  }
  .container-fixed {
    width: 100%;
    margin: 0 auto;
    padding-bottom: 0.2rem; /* 底部内边距 */
  }
  .background-container::before {
    filter: blur(4px); /* 减少模糊度提升性能 */
    background-attachment: scroll; /* 移动端使用普通滚动 */
    transform: scale(1.08);
  }
  
  .content-wrapper {
    overscroll-behavior: contain;
    position: absolute;
    top: 0;
    left: 0;
    right: 0;
    height: 100vh;
    overflow-y: auto;
    z-index: 1;
    pointer-events: auto;
  }
  
  .background-container {
    position: fixed;
  }
}
@media screen and (max-width: 768px) {
  .header-title {
    font-size: 1.8rem;
    top: 35%;
  }
  .header-subtitle {
    top: calc(35% + 50px); 
    font-size: 0.9rem;
  }
}

.profile-info {
  position: absolute;
  bottom: 20px;
  right: 20px;
  display: flex;
  flex-direction: row-reverse;  /* 改变方向，头像在右侧 */
  align-items: center;
  gap: 10px;
  max-width: 80%;  /* 限制最大宽度 */
}

.avatar {
  width: 50px;
  height: 50px;
  border-radius: 50%;
  border: 2px solid white;
  object-fit: cover;
  cursor: pointer;
  transition: transform 0.3s ease;
}

.avatar:hover {
  transform: scale(1.1);
}

.profile-text {
  text-align: left;  /* 改为左对齐 */
  min-width: 0;  /* 允许内容收缩 */
  overflow-x: auto;  /* 允许横向滚动 */
  scrollbar-width: none;  /* 隐藏滚动条 (Firefox) */
  -ms-overflow-style: none;  /* 隐藏滚动条 (IE/Edge) */
  padding: 5px 0;
}
.profile-text .title {
  font-size: 1.2rem;
  font-weight: bold;
}

.profile-text .description {
  font-size: 0.9rem;
  opacity: 0.9;
}
.profile-text::-webkit-scrollbar {
  display: none;
}

.profile-text .title {
font-size: 1.2rem;
font-weight: bold;
color: #fcfafb;  
text-shadow: none;
white-space: nowrap;  /* 防止换行 */

}

.profile-text .description {
  font-size: 0.9rem;
  color: #fcfafb; 
  text-shadow: none;
  white-space: nowrap;
  opacity: 0.95;
}
.u-container {
  backdrop-filter: blur(4px);
  border-radius: 8px;
  margin: 0 auto;
  max-width: 1200px;
  width: 100%;
  position: relative;
  z-index: 1;
  box-sizing: border-box;
  overflow-x: hidden;
  cursor: default;
}

.message-list-container { cursor: default; }

.loading {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
  background: rgba(0, 0, 0, 0.7); /* 更改背景色 */
  backdrop-filter: blur(8px);
  z-index: 1000;
  gap: 15px;
  opacity: 1;
  transition: opacity 0.3s ease;
}

.loading-text {
  font-size: 16px;
  color: #fff; /* 更改文字颜色 */
  text-shadow: none;
}

.rainbow-spinner {
  width: 50px;
  height: 50px;
  border-radius: 50%;
  border: 4px solid transparent;
  border-top: 4px solid #FF0000;
  border-right: 4px solid #00FF00;
  border-bottom: 4px solid #0000FF;
  border-left: 4px solid #FF00FF;
  animation: spin 1s linear infinite;
  will-change: transform;
}

@keyframes spin {
  0% { transform: rotate(0deg); }
  100% { transform: rotate(360deg); }
}
/* 添加新的样式 */
.message-list-container {
  position: relative;
  z-index: 10;
}

.container-fixed {
  min-height: 100vh;
  pointer-events: auto;
}

/* 确保背景不会遮挡评论框 */
.background-container::before {
  z-index: -1;
}
.heatmap-container {
  background: rgba(0, 0, 0, 0.3);
  backdrop-filter: blur(4px);
  border-radius: 8px;
  padding: 1rem;
  margin-bottom: 1rem;
}
.scroll-buttons {
  position: fixed;
  right: 20px;
  bottom: 20px;
  display: flex;
  flex-direction: column;
  gap: 10px;
  z-index: 1000;
}
.scroll-button {
  width: 2.25rem;
  height: 2.25rem;
  border-radius: 9999px;
  transition: transform 0.15s ease, background-color 0.15s ease, box-shadow 0.15s ease;
  backdrop-filter: blur(6px);
  display: inline-flex;
  align-items: center;
  justify-content: center;
  opacity: 0.95;
}
.scroll-button:hover {
  transform: scale(1.06);
}
</style>
