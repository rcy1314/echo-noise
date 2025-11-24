<template>
  <div class="background-container" :style="backgroundStyle">
    <div class="loading" v-if="!isLoaded">
      <div class="rainbow-spinner"></div>
      <div class="loading-text">加载中...</div>
    </div>
    <div class="content-wrapper" :class="{ 'gpu-accelerated': true }">
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
        <!-- 确保 MessageList 有足够的 z-index -->
      <MessageList 
      ref="messageList" 
      class="message-list-container" 
      :site-config="frontendConfig"
      :target-message-id="targetMessageId" 
      />
      </UContainer>
      <Notification />
      <!-- 添加搜索模态框组件 -->
      <SearchMode v-model="showSearchModal" @search-result="handleSearchResult" />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, inject, provide, onMounted, onUnmounted, watch, useRoute } from 'vue' // 添加 provide
import AddForm from '@/components/index/AddForm.vue'
import MessageList from '@/components/index/MessageList.vue'
import Notification from '~/components/widgets/Notification.vue';
import HeatmapWidget from '~/components/widgets/heatmap.vue'
import SearchMode from '~/components/index/Searchmode.vue' // 导入 SearchMode 组件
import TagList from '~/components/index/TagList.vue'

// 添加 messageList ref
const messageList = ref(null)
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
const contentTheme = ref<string>(typeof window !== 'undefined' ? (localStorage.getItem('contentTheme') || 'dark') : 'dark')
// 提供给子组件
provide('showHeatmap', showHeatmap);
provide('contentTheme', contentTheme)
provide('toggleContentTheme', () => {
  contentTheme.value = contentTheme.value === 'dark' ? 'light' : 'dark'
  if (typeof window !== 'undefined') localStorage.setItem('contentTheme', contentTheme.value)
})

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
    pwaIconURL: ''
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
    pageFooterHTML: '<div class="text-center text-xs text-gray-400 py-4">来自<a href="https://www.noisework.cn" target="_blank" rel="noopener noreferrer" class="text-orange-400 hover:text-orange-500">Noise</a> 使用<a href="https://github.com/lin-snow/Ech0" target="_blank" rel="noopener noreferrer" class="text-orange-400 hover:text-orange-500">Ech0</a>发布</div>',
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
    walineServerURL: 'https://s9cewscb.lc-cn-n1-shared.com'
};

// 修改 fetchConfig 方法
const fetchConfig = async () => {
    try {
        // 先设置默认值
        frontendConfig.value = { ...defaultConfig };
        
        const response = await fetch('/api/frontend/config', {
            credentials: 'include',
            headers: {
                'Cache-Control': 'no-cache',
                'Pragma': 'no-cache'
            }
        });
        
        if (!response.ok) {
            throw new Error('Network response was not ok');
        }
        
        const data = await response.json();
        console.log('获取到的配置数据:', data);
        
        if (data?.data?.frontendSettings) {
            const settings = data.data.frontendSettings;
            // 更新配置
            const booleanKeys = ['enableGithubCard', 'pwaEnabled']
            Object.keys(frontendConfig.value).forEach(key => {
                if (settings[key] !== null && settings[key] !== undefined) {
                    if (key === 'backgrounds' && Array.isArray(settings[key])) {
                        frontendConfig.value.backgrounds = [...settings[key]];
                    } else if (booleanKeys.includes(key)) {
                        const v = settings[key]
                        frontendConfig.value[key] = (v === true || v === 'true')
                    } else {
                        const v = settings[key]
                        frontendConfig.value[key] = typeof v === 'string' ? v.trim() : v
                    }
                }
            });
        }

        // 确保背景图片数组存在且有效
        if (!frontendConfig.value.backgrounds?.length) {
            frontendConfig.value.backgrounds = [...defaultConfig.backgrounds];
        }

        // 设置初始背景图
        if (frontendConfig.value.backgrounds.length > 0) {
            const randomIndex = Math.floor(Math.random() * frontendConfig.value.backgrounds.length);
            currentImage.value = frontendConfig.value.backgrounds[randomIndex];
        }
        
        isLoaded.value = true;
    } catch (error) {
        console.error('获取配置失败:', error);
        // 发生错误时保持默认配置
        frontendConfig.value = { ...defaultConfig };
        isLoaded.value = true;
    }
};
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
onMounted(() => {
    window.addEventListener('frontend-config-updated', async (event: CustomEvent) => {
        // 获取最新配置
        await fetchConfig();
        
        // 如果背景图片列表已更新，随机选择一张新图片
        if (frontendConfig.value.backgrounds?.length > 0) {
            const randomIndex = Math.floor(Math.random() * frontendConfig.value.backgrounds.length);
            const newImage = frontendConfig.value.backgrounds[randomIndex];
            
            // 更新当前显示的图片
            currentImage.value = newImage;
        }
    });
});

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
const updateTitle = () => {
  const title = (frontendConfig.value.pwaTitle || frontendConfig.value.siteTitle || '说说笔记').trim()
  const icon = '/favicon.ico'
  const description = (frontendConfig.value.pwaDescription || frontendConfig.value.description || '').trim()
  useHead({
    title,
    meta: [
      { name: 'description', content: description },
      { name: 'theme-color', content: '#000000' }
    ],
    link: [
      { rel: 'icon', href: icon },
      ...(frontendConfig.value.pwaEnabled ? [{ rel: 'manifest', href: '/site.webmanifest' }] : [])
    ]
  })
}

// 监听配置变化
watch(() => [frontendConfig.value.pwaEnabled, frontendConfig.value.pwaTitle, frontendConfig.value.pwaIconURL, frontendConfig.value.pwaDescription, frontendConfig.value.siteTitle, frontendConfig.value.rssFaviconURL, frontendConfig.value.description], () => {
  updateTitle()
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
    const response = await fetch(`${useRuntimeConfig().public.baseApi}/messages/tags`, {
      headers: {
        'Cache-Control': 'no-cache',
        'Pragma': 'no-cache'
      }
    })
    const data = await response.json()
    if (data.code === 1) {
      tags.value = data.data || []
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
    const encodedTag = encodeURIComponent(tag.trim());
    const response = await fetch(`${useRuntimeConfig().public.baseApi}/messages/tags/${encodedTag}`, {
      credentials: 'include',
      headers: {
        'Accept': 'application/json'
      }
    });
    
    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`);
    }
    
    const data = await response.json();
    
    if (data.code === 1 && Array.isArray(data.data)) {
      // 直接调用 messageList 的处理方法
      if (messageList.value) {
        messageList.value.handleSearchResult(data.data);
      }
    } else {
      throw new Error(data.msg || '获取标签内容失败');
    }
  } catch (error: any) {
    console.error('获取标签消息失败:', error);
    useToast().add({
      title: '获取标签消息失败',
      description: error.message || '服务器错误，请稍后重试',
      color: 'red',
      timeout: 3000
    });
  }
};
// 修改打字效果函数
const startTypeEffect = () => {
  if (!subtitleEl.value) return
  
  let index = 0
  let isDeleting = false
  let isWaiting = false
  
  const typeInterval = setInterval(() => {
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
  overflow: hidden; /* 防止滚动条出现 */
  overscroll-behavior: none; /* 防止橡皮筋效果 */
}
.header-subtitle {
  position: absolute;
  top: calc(50% + 60px);  /* 增加间距 */
  left: 50%;
  transform: translate(-50%, -50%);
  color: white;
  font-size: 1rem;
  text-shadow: 1px 1px 2px rgba(0, 0, 0, 0.5);
  white-space: nowrap;
}

@media screen and (max-width: 768px) {
  .header-subtitle {
    font-size: 0.9rem;
    top: calc(50% + 45px);  /* 保持移动端的适配 */
  }
}
.background-container {
  min-height: 100vh;
  width: 100%;  /* 移除 vw 单位 */
  position: fixed;  /* 改为 fixed */
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  overflow-y: auto;  /* 允许内容滚动 */
  overflow: hidden; /* 禁止背景容器滚动 */
  background-color: black; 
  z-index: -1; /* 确保背景在最底层 */
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
  position: relative;
  top: 0;
  left: 0;
  right: 0;
  height: auto;
  min-height: 100vh;
  overflow-y: auto; /* 允许内容区域滚动 */
  overscroll-behavior: contain; /* 防止滚动穿透 */
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
  text-shadow: 2px 2px 4px rgba(0, 0, 0, 0.5);
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
  }
  
  .content-wrapper {
    overscroll-behavior: contain;
    position: relative; /* 移动端改为相对定位 */
    height: auto;
    min-height: 100vh;
  }
  
  .background-container {
    position: absolute; /* 移动端改为绝对定位 */
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
text-shadow: 2px 2px 4px rgba(0, 0, 0, 0.5);
white-space: nowrap;  /* 防止换行 */

}

.profile-text .description {
  font-size: 0.9rem;
  color: #fcfafb; 
  text-shadow: 1px 1px 3px rgba(0, 0, 0, 0.5);
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
}

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
  text-shadow: 0 2px 4px rgba(0, 0, 0, 0.5);
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
  backdrop-filter: blur(4px);
  border-radius: 8px;
  margin: 0 auto;
  max-width: 1200px;
  width: 100%;
  position: relative;
  /* 降低 container 的 z-index，确保不会遮挡评论框 */
  z-index: 1;
  box-sizing: border-box;
  overflow: visible; /* 修改为 visible，允许评论框溢出 */
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
.background-container {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  overflow: hidden; /* 禁止背景滚动 */
  z-index: 0; /* 调整为0 */
}

.content-wrapper {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  height: 100vh;
  overflow-y: auto;
  z-index: 1; /* 确保内容在背景之上 */
  pointer-events: auto; /* 确保内容可点击 */
}

.container-fixed {
  min-height: 100vh; /* 确保容器足够高 */
  pointer-events: auto; /* 确保内容可点击 */
}
</style>
