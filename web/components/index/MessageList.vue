<template>
  <div class="min-h-screen flex flex-col">
    <!-- 空状态显示 -->
    <div v-if="!displayMessages.length" class="text-center text-gray-500 py-8">
      <UIcon name="i-heroicons-inbox" class="w-12 h-12 mx-auto mb-4" />
      <p>暂无消息内容</p>
    </div>
    
    <div class="flex-grow mx-auto w-full sm:max-w-2xl px-2">
      <!-- 搜索模式提示 -->
      <div 
        v-if="isSearchMode" 
        class="flex justify-between items-center mb-4 p-4 rounded-lg"
      >
        <p class="text-gray-400">搜索结果 ({{ searchResults.length }} 条)</p>
        <UButton
          size="sm"
          variant="ghost"
          class="text-gray-400 hover:text-orange-500"
          icon="i-heroicons-arrow-left"
          @click="resetList"
        >
          返回完整列表
        </UButton>
      </div>
      <!-- 消息列表 -->
      <div class="my-4">
         <!-- 无搜索结果提示 -->
  <div v-if="isSearchMode && searchResults.length === 0" class="text-center text-gray-500 py-8">
    <UIcon name="i-heroicons-magnifying-glass" class="w-12 h-12 mx-auto mb-4" />
    <p>未找到相关内容</p>
  </div>
        <!-- 消息列表内容 -->
        <div v-for="msg in displayMessages" :key="msg.id" class="w-full h-auto overflow-hidden flex flex-col justify-between">
           <!-- 修改头部布局 -->
           <div class="flex justify-between items-center">
            <!-- 时间部分保持不变 -->
            <div class="flex justify-start items-center h-auto overflow-x-auto whitespace-nowrap hide-scrollbar">
  <div class="w-2 h-2 rounded-full bg-orange-600 mr-2 flex-shrink-0"></div>
  <div class="flex justify-start text-sm">
    <span class="text-orange-500">{{ formatDate(msg.created_at) }}</span>
    <span class="gradient-dot mx-2 flex-shrink-0">@</span>
    <span class="text-orange-500">{{ msg.username || '匿名用户' }}</span>
  </div>
</div>
            <!-- 优化操作按钮组样式 -->
          <div class="message-actions flex justify-end items-center space-x-2 flex-shrink-0 px-3 py-1.5 mr-[9px] -mb-[1px]">
            <!-- ... 按钮内容 ... -->
              <div v-if="msg.private" class="w-5 h-5 flex-shrink-0 transition-transform duration-200 hover:scale-110">
                <UIcon name="i-mdi-lock-outline" class="text-gray-400" />
              </div>
              <div v-if="isLogin" class="w-5 h-5 cursor-pointer flex-shrink-0 transition-all duration-200 hover:scale-110" @click="editMessage(msg)" :title="'编辑内容'">
                <UIcon name="i-mdi-pencil-outline" class="text-gray-400 hover:text-orange-500" />
              </div>
              <div class="w-5 h-5 cursor-pointer flex-shrink-0 transition-all duration-200 hover:scale-110" @click="copyContent(msg.content)" :title="'复制内容'">
                <UIcon name="i-mdi-content-copy" class="text-gray-400 hover:text-orange-500" />
              </div>
              <div class="w-5 h-5 cursor-pointer flex-shrink-0 transition-all duration-200 hover:scale-110" @click="downloadAsImage(msg.id)" :title="'下载为图片'">
                <UIcon name="i-mdi-image-outline" class="text-gray-400 hover:text-orange-500" />
              </div>
              <div class="w-5 h-5 cursor-pointer flex-shrink-0 transition-all duration-200 hover:scale-110" @click="toggleComment(msg.id)" :title="'评论'">
                <UIcon name="i-mdi-comment-outline" class="text-gray-400 hover:text-orange-500" />
              </div>
              <div v-if="isLogin" class="w-5 h-5 cursor-pointer flex-shrink-0 transition-all duration-200 hover:scale-110" @click="deleteMsg(msg.id)" :title="'删除'">
                <UIcon name="i-mdi-paper-roll-outline" class="text-gray-400 hover:text-orange-500" />
              </div>
            </div>
          </div>

          <div class="border-l-2 border-gray-300 p-4 ml-1">
            <div class="content-container" :class="listThemeClass" v-if="msg.image_url || msg.content" :data-msg-id="msg.id">
              <!-- 图片内容 -->
              <img 
  v-if="msg.image_url" 
  :src="`${BASE_API}${msg.image_url}`" 
  alt="Image" 
  class="max-w-full object-cover rounded-lg mb-4"
  loading="lazy"
  :fetchpriority="index < 3 ? 'high' : 'low'"
/>
              <!-- 分隔线 -->
              <div v-if="msg.image_url && msg.content" class="border-t border-gray-600 my-4"></div>
              <!-- 文本内容区域 -->
              <div class="overflow-y-hidden relative" :class="[{ 'max-h-[700px]': !isExpanded[msg.id] }, listThemeTextClass]">
                <MarkdownRenderer :content="msg.content" :enableGithubCard="siteConfig?.enableGithubCard === true" @tagClick="handleTagClick" link-target="_blank"/>
                <div v-if="shouldShowExpandButton[msg.id] && !isExpanded[msg.id]"
    :class="['absolute bottom-0 left-0 right-0 h-32 bg-gradient-to-t', gradientClass]">
  </div>
              </div>
              <!-- 展开/折叠按钮 -->
              <div v-if="shouldShowExpandButton[msg.id]" class="text-center mt-2 relative" style="z-index: 9999;">
                <button @click="toggleExpand(msg.id)"
                  class="flex items-center justify-center space-x-1 mx-auto px-4 py-2 text-orange-500 hover:text-orange-600 focus:outline-none transition-colors duration-200">
                  <span>{{ isExpanded[msg.id] ? "收起内容" : "展开全文" }}</span>
                  <UIcon :name="isExpanded[msg.id] ? 'i-mdi-chevron-up' : 'i-mdi-chevron-down'" class="w-5 h-5" />
                </button>
              </div>
            </div>
            <!-- 评论区域 -->
            <div v-show="activeCommentId === msg.id" class="mt-4">
              <div :id="`waline-${msg.id}`"></div>
            </div>
          </div>
        </div>
      </div>
      <!-- 分页控制区域 -->
      <div v-if="!isSearchMode" class="flex justify-center items-center space-x-4 w-full my-4 flex-wrap md:flex-nowrap">
  <div class="flex justify-center items-center space-x-4 w-full md:w-auto">
    <UButton 
      v-if="message.page > 1"
      color="gray" 
      variant="outline" 
      size="xs" 
      class="rounded-full px-4 py-1.5 bg-[rgba(36,43,50,0.95)] text-white hover:text-white border-none shadow-lg hover:shadow-xl transition-all duration-300 backdrop-blur-sm"
      @click="loadPreviousPage"
      :disabled="isPageLoading"
    >
      <UIcon name="i-heroicons-arrow-left" class="mr-1 w-4 h-4" /> 
      上一页
    </UButton>

    <UButton 
      v-if="message.hasMore"
      color="gray" 
      variant="outline" 
      size="xs" 
      class="rounded-full px-4 py-1.5 bg-[rgba(36,43,50,0.95)] text-white hover:text-white border-none shadow-lg hover:shadow-xl transition-all duration-300 backdrop-blur-sm"
      @click="loadNextPage"
      :disabled="isPageLoading"
    >
      下一页
      <UIcon name="i-heroicons-arrow-right" class="ml-1 w-4 h-4" />
    </UButton>
    <span v-if="isPageLoading" class="ml-2 text-orange-400">加载中...</span>
  </div>

  <!-- 页码显示和跳转 -->
  <div class="flex items-center justify-center space-x-2 w-full md:w-auto mt-3 md:mt-0">
    <span class="text-gray-500 text-shadow-sm text-sm">第 {{ message.page }} 页</span> 
    <UInput
      v-model="targetPage"
      type="number"
      min="1"
      :max="totalPages"
      class="w-12 text-center text-sm" 
      placeholder="#"
      @keyup.enter="jumpToPage"
    />
    <UButton
      size="xs" 
      color="gray"
      variant="ghost"
      class="text-gray-400 hover:text-orange-500 text-sm"  
      @click="jumpToPage"
    >
      跳转
    </UButton>
  </div>
</div>
      <!-- 加载完毕提示 -->
      <div v-if="!isSearchMode && message.messages.length > 0 && !message.hasMore" class="text-center text-gray-500 mt-4">
        <UIcon name="i-fluent-emoji-flat-confetti-ball" size="lg" />
        加载完毕~
      </div>
    </div>
    <!-- 来源信息 - 固定在底部 -->
    <div v-if="!siteConfig.pageFooterHTML" class="text-center text-xs text-gray-400 py-4">
    来自<a href="https://www.noisework.cn" target="_blank" rel="noopener noreferrer" class="text-orange-400 hover:text-orange-500">Noise</a> 
    使用<a href="https://github.com/lin-snow/Ech0" target="_blank" rel="noopener noreferrer" class="text-orange-400 hover:text-orange-500">Ech0</a>发布
  </div>
  <div v-else v-html="siteConfig.pageFooterHTML"></div>
</div>
  <!-- 编辑对话框 -->
  <UModal v-model="showEditModal" :ui="{ width: 'sm:max-w-2xl' }">
    <UCard>
      <template #header>
        <div class="flex justify-between items-center">
          <h3 class="text-lg font-medium">编辑内容</h3>
          <UButton color="gray" variant="ghost" icon="i-mdi-close" class="-my-1" @click="showEditModal = false" />
        </div>
      </template>
      <div class="flex flex-col space-y-4">
        <UTextarea
          v-model="editingContent"
          placeholder="编辑内容..."
          rows="10"
          class="font-mono text-sm"
        />
        <div class="border-t border-gray-200 my-2 pt-2">
          <div class="text-sm text-gray-500 mb-2">预览：</div>
          <!-- 修改预览区域样式 -->
          <div class="p-4 rounded-lg overflow-auto max-h-[300px] bg-[rgba(36,43,50,0.95)]">
            <div class="text-white">
              <MarkdownRenderer :content="editingContent" :enableGithubCard="siteConfig?.enableGithubCard === true" />
            </div>
          </div>
        </div>
      </div>
      <template #footer>
        <div class="flex justify-end space-x-2">
          <UButton color="gray" variant="outline" @click="showEditModal = false" class="text-white">
            取消
          </UButton>
          <UButton color="orange" @click="saveEditedMessage" :loading="isSaving" class="text-white">
            保存
          </UButton>
        </div>
      </template>
    </UCard>
  </UModal>
</template>

<script setup lang="ts">
import { useMessageStore } from "~/store/message";
import { useUserStore } from "~/store/user";
import MarkdownRenderer from "~/components/index/MarkdownRenderer.vue";
const contentTheme = inject('contentTheme', ref<string>(typeof window !== 'undefined' ? (localStorage.getItem('contentTheme') || 'dark') : 'dark'))
const listThemeClass = computed(() => contentTheme.value === 'dark' ? 'bg-[rgba(36,43,50,0.95)] text-white' : 'bg-white text-black')
const listThemeTextClass = computed(() => contentTheme.value === 'dark' ? 'text-white' : 'text-black')
const gradientClass = computed(() => contentTheme.value === 'dark' ? 'from-[rgba(36,43,50,1)] via-[rgba(36,43,50,0.8)] to-transparent' : 'from-[rgba(255,255,255,1)] via-[rgba(255,255,255,0.8)] to-transparent')

const targetPage = ref('');
const totalPages = computed(() => Math.ceil(message.total / 15));
const jumpToPage = async () => {
  const page = parseInt(targetPage.value);
  if (!page || page < 1 || page > totalPages.value || message.loading) {
    useToast().add({
      title: '页码无效',
      description: `请输入 1-${totalPages.value} 之间的数字`,
      color: 'orange',
      timeout: 2000
    });
    return;
  }

  try {
    const result = await message.getMessages({
      page,
      pageSize: 15,
    });
    
    if (!result) {
      throw new Error('跳转页面失败');
    }
    
    // 更新消息列表
    message.messages = result.items;
    
    // 滚动到顶部
    window.scrollTo({ 
      top: 0,
      behavior: 'instant'
    });
    
    targetPage.value = ''; // 清空输入框
  } catch (error) {
    console.error('跳转页面失败:', error);
    useToast().add({
      title: '跳转失败',
      color: 'red',
      timeout: 2000
    });
  }
};
// 添加 props 定义
const props = defineProps({
  siteConfig: {
    type: Object,
    required: true
  },
  targetMessageId: {  // 添加新的 prop
    type: String,
    default: null
  }
});
// 添加监听器
watch(() => props.targetMessageId, async (newId) => {
  if (!newId) return;
  
  await nextTick();
  const targetElement = document.querySelector(`.content-container[data-msg-id="${newId}"]`);
  if (targetElement) {
    targetElement.scrollIntoView({ behavior: 'smooth', block: 'center' });
    // 添加高亮效果
    targetElement.classList.add('highlight-message');
    setTimeout(() => {
      targetElement.classList.remove('highlight-message');
    }, 2000);
  }
}, { immediate: true });

const BASE_API = useRuntimeConfig().public.baseApi;
const { deleteMessage } = useMessage();
const message = useMessageStore();

const activeCommentId = ref<number | null>(null);
const userStore = useUserStore();
const isLogin = computed(() => userStore.isLogin);
const openInNewTab = (url: string) => {
  window.open(url, '_blank', 'noopener,noreferrer');
};
// 修改标签点击处理函数
const handleTagClick = async (tag: string) => {
  try {
    // 确保 tag 被正确编码
    const encodedTag = encodeURIComponent(tag.trim());
    const response = await fetch(`${BASE_API}/messages/tags/${encodedTag}`, {
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
      isSearchMode.value = true;
      searchResults.value = data.data;
      
      await nextTick();
      checkContentHeight();
      initFancybox();
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
// 修改重置搜索函数名称，使其更通用
// 修改 resetList 函数
const resetList = async () => {
  searchResults.value = [];
  isSearchMode.value = false;
  
  // 重新获取完整消息列表
  await message.getMessages({
    page: 1,
    pageSize: 15,
  });
  
  await nextTick();
  checkContentHeight();
  initFancybox();
};

const deleteMsg = async (id: number) => {
  const confirmDelete = confirm("确定要删除这条消息吗？");
  if (confirmDelete) {
    try {
      await message.deleteMessage(id); // 使用 store 中的方法
      message.messages = message.messages.filter(msg => msg.id !== id);
      useToast().add({
        title: '删除成功',
        color: 'green',
        timeout: 2000
      });
    } catch (error) {
      console.error('删除失败:', error);
      useToast().add({
        title: '删除失败',
        color: 'red',
        timeout: 2000
      });
    }
  }
};

const initFancybox = () => {
  if (window.Fancybox) {
    window.Fancybox.destroy();
    const fancyboxOptions = {
      Carousel: {
        infinite: false,
      },
      Toolbar: {
        display: [
          { id: "prev", position: "center" },
          { id: "counter", position: "center" },
          { id: "next", position: "center" },
          "zoom",
          "slideshow",
          "fullscreen",
          "close",
        ],
      },
      Image: {
        zoom: true,
        click: true,
        wheel: "slide",
      },
    };

    // 只为远程图片添加灯箱，且每张图片只包裹一次
    const mdImages = document.querySelectorAll(".markdown-preview img");
    mdImages.forEach((img) => {
      const src = img.getAttribute("src") || "";
      const isRemote = /^https?:\/\//i.test(src);
      const parent = img.parentElement;
      // 如果已经被包裹且包裹正确，跳过
      if (parent && parent.tagName === "A" && parent.hasAttribute("data-fancybox")) {
        // 如果不是远程图片，移除包裹
        if (!isRemote) {
          parent.replaceWith(img);
        }
        return;
      }
      // 只包裹远程图片
      if (isRemote) {
        const wrapper = document.createElement("a");
        wrapper.href = src;
        wrapper.setAttribute("data-fancybox", "uploaded-image");
        wrapper.style.display = "block";
        img.parentNode.insertBefore(wrapper, img);
        wrapper.appendChild(img);
      } else if (parent && parent.tagName === "A" && parent.hasAttribute("data-fancybox")) {
        // 非远程图片且被包裹，移除包裹
        parent.replaceWith(img);
      }
    });

    window.Fancybox.bind("[data-fancybox]", fancyboxOptions);
  }
};

const toggleComment = async (msgId: number) => {
  if (activeCommentId.value === msgId) {
    activeCommentId.value = null;
  } else {
    activeCommentId.value = msgId;
    await nextTick();
    if (window.Waline) {
      const el = document.querySelector(`#waline-${msgId}`);
      if (el) {
        window.Waline.init({
          el: `#waline-${msgId}`,
          serverURL: props.siteConfig.walineServerURL,
          path: 'messages/${msgId}',
          reaction: false,
          meta: ["nick", "mail", "link"],
          requiredMeta: ["mail", "nick"],
          pageview: true,
          search: false,
          wordLimit: 200,
          pageSize: 5,
          avatar: "monsterid",
          emoji: ["https://unpkg.com/@waline/emojis@1.2.0/tieba"],
          imageUploader: false,
          copyright: false,
          dark: 'html[class="dark"]',
        });
      } else {
        console.error(`评论容器 #waline-${msgId} 未找到`);
      }
    } else {
      console.error("Waline 未加载");
    }
  }
};

const formatDate = (dateString: string) => {
  const date = new Date(dateString);
  const now = new Date();
  const diff = now.getTime() - date.getTime();
  const diffInDays = Math.floor(diff / (1000 * 60 * 60 * 24));
  const diffInHours = Math.floor(diff / (1000 * 60 * 60));
  const diffInMinutes = Math.floor(diff / (1000 * 60));

  const diffInSeconds = Math.floor(diff / 1000);
  if (diffInSeconds < 60) {
    return "刚刚";
  } else if (diffInMinutes < 60) {
    return `${diffInMinutes}分钟前`;
  } else if (diffInHours < 24) {
    return `${diffInHours}小时前`;
  } else if (diffInDays < 3) {
    return `${diffInDays}天前`;
  } else {
    return date.toLocaleString();
  }
};
// 添加展开状态管理
const isExpanded = ref<{ [key: number]: boolean }>({});
const shouldShowExpandButton = ref<{ [key: number]: boolean }>({});

// 添加展开/折叠切换函数
const toggleExpand = (msgId: number) => {
  isExpanded.value[msgId] = !isExpanded.value[msgId];
};

// 修改检查内容高度的函数
// 修改检查内容高度的函数
const checkContentHeight = () => {
  nextTick(() => {
    // 获取当前显示的消息列表（可能是普通列表或搜索结果）
    const currentMessages = isSearchMode.value ? searchResults.value : message.messages;
    
    // 检查每条消息的内容高度
    currentMessages.forEach((msg) => {
      const contentEl = document.querySelector(
        `.content-container[data-msg-id="${msg.id}"] .overflow-y-hidden`
      );
      if (contentEl && contentEl.scrollHeight > 700) {
        shouldShowExpandButton.value[msg.id] = true;
        if (isExpanded.value[msg.id] === undefined) {
          isExpanded.value[msg.id] = false;
        }
      }
    });
  });
};

// 确保在内容变化时重新检查高度
watch(() => message.messages, () => {
  // 如果是单条消息查看模式，不执行滚动
  if (route.hash.includes('/messages/')) {
    return;
  }
  nextTick(() => {
    requestAnimationFrame(() => {
      checkContentHeight();
      initFancybox();
    });
  });
}, { deep: true });
// 添加路由相关
const route = useRoute();
onMounted(async () => {
  try {
    // 获取路由中的消息ID
    const messageId = route.hash.split('/messages/').pop();
    
    // 加载 Waline
    if (!window.Waline) {
      const link = document.createElement("link");
      link.rel = "stylesheet";
      link.href = "https://unpkg.com/@waline/client@v2/dist/waline.css";
      document.head.appendChild(link);

      await new Promise((resolve, reject) => {
        const script = document.createElement("script");
        script.src = "https://unpkg.com/@waline/client@v2/dist/waline.js";
        script.crossOrigin = "anonymous";
        script.onload = resolve;
        script.onerror = reject;
        document.head.appendChild(script);
      });
    }

    // 根据是否有消息ID来决定加载方式
    if (messageId) {
      const response = await fetch(`${BASE_API}/messages/${messageId}`, {
        credentials: 'include',
        headers: {
          'Accept': 'application/json'
        }
      });
      
      if (!response.ok) throw new Error('消息加载失败');
      
      const data = await response.json();
      if (data.code === 1 && data.data) {
        // 设置单条消息模式
        message.messages = [data.data];
        message.hasMore = false;
        message.page = 1;
        
        await nextTick();
        const targetElement = document.querySelector(`.content-container[data-msg-id="${messageId}"]`);
        if (targetElement) {
          targetElement.scrollIntoView({ behavior: 'instant', block: 'start' });
        }
      } else {
        throw new Error('消息不存在');
      }
    } else {
      // 只有在非消息详情页时才加载列表
      if (!route.hash.includes('/messages/')) {
        await message.getMessages({
          page: 1,
          pageSize: 15,
        });
      }
    }

    // 初始化视图
    await nextTick();
    checkContentHeight();
    initFancybox();
    
  } catch (error) {
    console.error('初始化失败:', error);
    if (error instanceof Error) {
      useToast().add({
        title: '加载失败',
        description: error.message || '请刷新重试',
        color: 'red',
        timeout: 2000
      });
    }
  }
});

// 修改路由监听
watch(() => route.hash, async (newHash) => {
  const messageId = newHash.split('/messages/').pop();
  
  // 如果没有消息ID且不是从消息详情页返回，则加载列表
  if (!messageId) {
    if (!route.hash.includes('/messages/')) {
      await message.getMessages({
        page: 1,
        pageSize: 15,
      });
    }
    return;
  }
  
  try {
    const response = await fetch(`${BASE_API}/messages/${messageId}`, {
      credentials: 'include',
      headers: {
        'Accept': 'application/json'
      }
    });
    
    if (!response.ok) throw new Error('消息加载失败');
    
    const data = await response.json();
    if (data.code === 1 && data.data) {
      message.messages = [data.data];
      message.hasMore = false;
      message.page = 1;
      
      await nextTick();
      const targetElement = document.querySelector(`.content-container[data-msg-id="${messageId}"]`);
      if (targetElement) {
        targetElement.scrollIntoView({ 
          behavior: 'instant',
          block: 'start'
        });
      }
    }
  } catch (error) {
    console.error('加载消息失败:', error);
    useToast().add({
      title: '加载失败',
      color: 'red',
      timeout: 2000
    });
  }
}, { immediate: true });

// 修改 loadMore 为 loadNextPage
const isPageLoading = ref(false);

const loadPreviousPage = async () => {
  if (isPageLoading.value || message.page <= 1) return;
  isPageLoading.value = true;
  try {
    const targetPage = message.page - 1;
    const result = await message.getMessages({
      page: targetPage,
      pageSize: 15,
    });
    if (result && Array.isArray(result.items)) {
      message.messages = result.items;
      message.page = result.page || targetPage;
    } else {
      message.page = targetPage;
    }
    window.scrollTo({ top: 0, behavior: 'instant' });
  } catch (error) {
    useToast().add({
      title: '加载失败',
      color: 'red',
      timeout: 2000
    });
  } finally {
    isPageLoading.value = false;
  }
};

const loadNextPage = async () => {
  if (isPageLoading.value || !message.hasMore) return;
  isPageLoading.value = true;
  try {
    const targetPage = message.page + 1;
    const result = await message.getMessages({
      page: targetPage,
      pageSize: 15,
    });
    if (result && Array.isArray(result.items)) {
      message.messages = result.items;
      message.page = result.page || targetPage;
    } else {
      message.page = targetPage;
    }
    window.scrollTo({ top: 0, behavior: 'instant' });
  } catch (error) {
    useToast().add({
      title: '加载失败',
      color: 'red',
      timeout: 2000
    });
  } finally {
    isPageLoading.value = false;
  }
};
// 添加登录状态变化监听
watch(
  () => userStore.isLogin,
  (newVal) => {
    if (newVal) {
      // 用户登录后的处理
      message.getMessages({
        page: 1,
        pageSize: 15,
      });
    }
  }
);

// 监听消息变化
watch(
  () => message.messages,
  async () => {
    try {
      await nextTick();
      checkContentHeight();
      initFancybox();
    } catch (error) {
      console.error('更新视图失败:', error);
    }
  },
  { deep: true }
);
// 组件卸载时清理
onBeforeUnmount(() => {
  if (window.Fancybox) {
    window.Fancybox.destroy();
  }
});
// 添加复制功能
const copyContent = async (content: string) => {
  try {
    await navigator.clipboard.writeText(content);
    // 可以使用 Nuxt 的 toast 提示复制成功
    useToast().add({
      title: '复制成功',
      color: 'green',
      timeout: 2000
    });
  } catch (err) {
    console.error('复制失败:', err);
    useToast().add({
      title: '复制失败',
      color: 'red',
      timeout: 2000
    });
  }
};
// 添加编辑功能
const showEditModal = ref(false);
const editingContent = ref('');
const editingMessageId = ref<number | null>(null);
const isSaving = ref(false);

const editMessage = (msg: any) => {
  editingMessageId.value = msg.id;
  
  // 保存原始内容，不包含附件图片
  editingContent.value = msg.content;
  
  // 如果存在附件图片，添加到编辑器中以便用户可以看到和编辑
  if (msg.image_url) {
    const imageMarkdown = `\n\n<!-- 附件图片(编辑时可删除) -->\n![附件图片](${BASE_API}${msg.image_url})\n<!-- 附件图片结束 -->`;
    editingContent.value += imageMarkdown;
  }
  
  showEditModal.value = true;
};

const saveEditedMessage = async () => {
  if (!editingMessageId.value) return;
  
  isSaving.value = true;
  try {
    // 获取当前编辑的消息
    const currentMsg = message.messages.find(msg => msg.id === editingMessageId.value);
    if (!currentMsg) return;

    // 处理编辑内容，移除附件图片的 Markdown 标记
    let processedContent = editingContent.value;
    
    // 移除附件图片的 Markdown 标记
    processedContent = processedContent.replace(/\n*<!-- 附件图片\(编辑时可删除\) -->\n!\[附件图片\]\(.*?\)\n<!-- 附件图片结束 -->\n*/g, '');
    
    // 检查内容是否有修改
    if (processedContent === currentMsg.content) {
      useToast().add({
        title: '内容未修改',
        description: '请修改内容后再保存',
        color: 'orange',
        timeout: 2000
      });
      isSaving.value = false;
      return;
    }
    // 直接使用编辑器中的内容，不做任何修改
    const response = await fetch(`${BASE_API}/messages/${editingMessageId.value}`, {
      method: 'PUT',
      headers: {
        'Content-Type': 'application/json',
        'Accept': 'application/json'
      },
      credentials: 'include',
      body: JSON.stringify({
        content: processedContent,
        image_url: currentMsg.image_url
      })
    });

    if (!response.ok) throw new Error(`HTTP error! status: ${response.status}`);

    const data = await response.json();
    if (data.code === 1) {
      const index = message.messages.findIndex(msg => msg.id === editingMessageId.value);
      if (index !== -1) {
        message.messages[index] = {
          ...message.messages[index],
          content: editingContent.value,  // 修正：使用 editingContent.value 替代 pureContent
          image_url: currentMsg.image_url  // 修正：使用 currentMsg.image_url 替代 imageUrl
        };
      }
      showEditModal.value = false;
      useToast().add({
        title: '更新成功',
        color: 'green',
        timeout: 2000
      });
    } else {
      throw new Error(data.msg || '保存失败');
    }
  } catch (error) {
    console.error('更新消息失败:', error);
    useToast().add({
      title: '更新失败',
      color: 'red',
      timeout: 2000
    });
  } finally {
    isSaving.value = false;
  }
};
const downloadAsImage = async (msgId: number) => {
  try {
    const element = document.querySelector(`.content-container[data-msg-id="${msgId}"]`);
    if (!element) return;

    // 检查内容类型
    const hasText = element.querySelector('.markdown-preview')?.textContent?.trim();
    const hasImage = element.querySelector('img');
    const hasVideo = element.querySelector('video');
    const hasAudio = element.querySelector('audio');

    // 纯视频或纯音频内容不生成卡片
    if ((!hasText && !hasImage && hasVideo) || (!hasText && !hasImage && hasAudio)) {
      useToast().add({
        title: '此内容不可生成卡片',
        color: 'orange',
        timeout: 2000
      });
      return;
    }

    // 设置超时检测
    const timeout = setTimeout(() => {
      useToast().add({
        title: '生成超时',
        description: '卡片生成时间过长，请稍后重试',
        color: 'red',
        timeout: 3000
      });
    }, 10000);

    // 1. 临时展开内容
    const originalExpanded = isExpanded.value[msgId];
    isExpanded.value[msgId] = true;
    await nextTick();

    // 2. 创建临时容器
    const tempContainer = document.createElement('div');
   tempContainer.style.cssText = `
  padding: 16px;
  background: transparent;
  border-radius: 12px;
  width: ${hasImage ? '640px' : '480px'};
  position: absolute;
  left: -9999px;
  top: 0;
  z-index: -1;
  overflow: visible;
  min-height: fit-content;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.1);
  border: none;
`;
    document.body.appendChild(tempContainer);
    
    // 3. 复制并处理内容
    const contentClone = element.cloneNode(true) as HTMLElement;
    const parentStyles = window.getComputedStyle(element as HTMLElement);
    const bgColor = parentStyles.backgroundColor && parentStyles.backgroundColor !== 'rgba(0, 0, 0, 0)'
      ? parentStyles.backgroundColor
      : '#ffffff';
    const textColor = parentStyles.color || '#333';
    
    // 移除所有控制元素和限制
    contentClone.querySelectorAll('.text-center.mt-2, .bg-gradient-to-t').forEach(el => el.remove());
    contentClone.style.cssText = `
      max-height: none;
      overflow: visible;
      padding: 0;
      margin: 0;
      background: ${bgColor};
      border-radius: 12px;
    `;
    
    // 处理内容区域
    const contentArea = contentClone.querySelector('.overflow-y-hidden');
    if (contentArea) {
      contentArea.className = '';
      contentArea.style.cssText = `
        overflow: visible;
        max-height: none !important;
        height: auto !important;
        padding: 12px;
        line-height: 1.6;
        margin-bottom: 0;
        white-space: pre-wrap;
        background: ${bgColor};
        border-radius: 12px;
        font-size: 14px;
        color: ${textColor};
      `;
    }

    // 处理媒体元素
    const mediaElements = contentClone.querySelectorAll('video, audio, iframe');
    mediaElements.forEach(media => {
      const placeholder = document.createElement('div');
      placeholder.style.cssText = `
        padding: 15px;
        background: rgba(251, 146, 60, 0.1);
        border: 1px solid rgba(251, 146, 60, 0.3);
        border-radius: 8px;
        color: #fb923c;
        margin: 15px 0;
        word-break: break-all;
      `;
      
      if (media instanceof HTMLVideoElement) {
        placeholder.innerHTML = `🎬 视频链接：${media.src || '未知链接'}`;
      } else if (media instanceof HTMLAudioElement) {
        placeholder.innerHTML = `🎵 音频链接：${media.src || '未知链接'}`;
      } else if (media instanceof HTMLIFrameElement) {
        placeholder.innerHTML = `🔗 嵌入内容链接：${media.src || '未知链接'}`;
      }
      
      media.parentNode?.replaceChild(placeholder, media);
    });

   // 处理图片
   const processImages = async () => {
  const images = contentClone.querySelectorAll('img');
  await Promise.all(Array.from(images).map(async (img) => {
    return new Promise<void>((resolve) => {
      const originalSrc = img.src;
      img.crossOrigin = 'anonymous';
      
      // 处理图片路径并添加 credentials
      if (originalSrc.startsWith('/')) {
        img.src = `${BASE_API}${originalSrc}`;
        // 为图片请求添加 credentials
        fetch(img.src, { credentials: 'include' })
          .then(response => response.blob())
          .then(blob => {
            img.src = URL.createObjectURL(blob);
            resolve();
          })
          .catch(() => {
            console.error('图片加载失败:', originalSrc);
            img.parentElement?.removeChild(img);
            resolve();
          });
      } else {
        if (img.complete) {
          resolve();
        } else {
          img.onload = () => resolve();
          img.onerror = () => {
            console.error('图片加载失败:', originalSrc);
            img.parentElement?.removeChild(img);
            resolve();
          };
        }
      }
    });
  }));
};

await processImages();

    tempContainer.appendChild(contentClone);

    // 添加 footer
    const footer = document.createElement('div');
  footer.style.cssText = `
    margin-top: 12px;
    padding-top: 12px;
    text-align: center;
    font-family: -apple-system, BlinkMacSystemFont, sans-serif;
    background: transparent;
  `;
  footer.innerHTML = `
    <div style="color: #fb923c; font-size: 13px; margin-bottom: 4px; font-weight: 500;">
      ${props.siteConfig.cardFooterTitle}
    </div>
    <a href="https://note.noisework.cn" 
       target="_blank" 
       rel="noopener noreferrer" 
       style="color: rgba(255,255,255,0.5); text-decoration: none;">
      ${props.siteConfig.cardFooterSubtitle}
    </a>
  `;
    tempContainer.appendChild(footer);

    // 生成图片
    await nextTick();
    const canvas = await html2canvas(tempContainer, {
      backgroundColor: null,
      scale: 2,
      useCORS: true,
      allowTaint: true,
      logging: false,
      width: tempContainer.offsetWidth,
      height: tempContainer.scrollHeight,
      borderRadius: '16px',
      onclone: (clonedDoc) => {
        const clonedElement = clonedDoc.querySelector('.content-container');
        if (clonedElement) {
          clonedElement.style.cssText = `
            overflow: visible !important;
            max-height: none !important;
            height: auto !important;
            padding: 0;
            min-height: ${contentArea?.scrollHeight || 0}px;
            background: ${bgColor};
            border-radius: 12px;
          `;
        }
      }
    });

    // 清除超时检测
    clearTimeout(timeout);
    // 下载图片
    const link = document.createElement('a');
    link.download = `message-${msgId}.png`;
    link.href = canvas.toDataURL('image/png');
    link.click();

    // 清理临时元素
    document.body.removeChild(tempContainer);
    
    // 恢复原始展开状态
    isExpanded.value[msgId] = originalExpanded;

    useToast().add({
      title: '下载成功',
      color: 'green',
      timeout: 2000
    });
  } catch (error) {
    console.error('下载失败:', error);
    useToast().add({
      title: '下载失败',
      color: 'red',
      timeout: 2000
    });
  }
};

// 添加搜索相关变量
const isSearchMode = ref(false);
const searchResults = ref([]);

// 添加搜索结果处理函数
const handleSearchResult = async (results: any) => {
  try {
    // 如果当前不是搜索模式，记录滚动位置
    const scrollPosition = !isSearchMode.value ? window.scrollY : null;
    
    console.debug('API返回的原始数据:', results);
    
    if (!results) {
      throw new Error('API返回数据为空');
    }
    
    let items = [];
    let total = 0;
    
    // 统一数据处理逻辑
    if (results.code === 1) {
      if (Array.isArray(results.data)) {
        items = results.data;
      } else if (results.data?.items) {
        items = results.data.items;
      }
    } else if (Array.isArray(results)) {
      items = results;
    }
    
    if (!Array.isArray(items)) {
      throw new Error('无效的数据格式');
    }
    
    total = items.length;
    
    // 更新搜索状态和结果
    isSearchMode.value = true;
    searchResults.value = items;
    
    // 显示结果提示
    if (total === 0) {
      useToast().add({
        title: '未找到相关内容',
        color: 'orange',
        timeout: 2000
      });
    } else {
      useToast().add({
        title: `找到 ${total} 条结果`,
        color: 'green',
        timeout: 2000
      });
    }
    
    // 如果是从非搜索模式切换来的，滚动到顶部
    if (scrollPosition !== null) {
      window.scrollTo({ top: 0, behavior: 'smooth' });
    }
    
    await nextTick();
    checkContentHeight();
    initFancybox();
    
  } catch (error: any) {
    console.error('处理搜索结果时出错:', error);
    useToast().add({
      title: '搜索失败',
      description: error.message || '处理搜索结果时发生错误',
      color: 'red',
      timeout: 2000
    });
    resetSearch();
  }
};
// 添加重置搜索函数
const resetSearch = () => {
  // 先清空结果数组
  searchResults.value = [];
  // 再关闭搜索模式
  isSearchMode.value = false;
  
  console.log('重置搜索 - searchResults:', searchResults.value);
  console.log('重置搜索 - isSearchMode:', isSearchMode.value);
  
  // 重置后更新UI
  nextTick(() => {
    checkContentHeight();
    initFancybox();
  });
};

// 修改displayMessages计算属性以支持搜索模式
const displayMessages = computed(() => {
  if (isSearchMode.value && Array.isArray(searchResults.value)) {
    return searchResults.value;
  }
  return message.messages || []; // 确保返回数组，即使是空数组
});

// 添加事件监听
defineExpose({
  handleSearchResult
});

// 添加watch监听searchResults变化
watch(searchResults, (newVal) => {
  console.log('searchResults变化:', newVal);
  // 强制更新内容高度检查
  nextTick(() => {
    checkContentHeight();
    initFancybox();
  });
}, { deep: true, immediate: true });

// 添加watch监听isSearchMode变化
watch(isSearchMode, (newVal) => {
  console.log('isSearchMode变化:', newVal);
  // 强制更新内容高度检查
  nextTick(() => {
    checkContentHeight();
    initFancybox();
  });
});
// 优化图片加载
const optimizeImage = (url: string) => {
  if (!url) return url;
  // 添加图片压缩参数
  return `${url}?imageView2/2/w/800/q/85`;
}
// 确保在模板中使用正确的配置数据
const footerConfig = computed(() => ({
  cardFooterTitle: props.siteConfig.cardFooterTitle,
  cardFooterSubtitle: props.siteConfig.cardFooterSubtitle,
  pageFooterHTML: props.siteConfig.pageFooterHTML,
  walineServerURL: props.siteConfig.walineServerURL
}));

</script>

<style scoped>
/* 修改内容卡片样式 */
.content-container {
  padding: 12px;
  border-radius: 12px;
  transition: all 0.3s ease;
  margin: 4px 0 0.2rem 0; /* 调整外边距 */
  width: 100%;
  box-sizing: border-box;
  position: relative;
  overflow: hidden;
}
/* 优化图片渲染 */
.content-container img {
  width: 100%;
  height: auto;
  min-height: 150px;
  object-fit: cover;
  border-radius: 12px;
  box-shadow: none;  /* 移除阴影 */
  transform: translate3d(0, 0, 0);  /* 启用硬件加速 */
}
/* 简化过渡动画 */
.overflow-y-hidden {
  transition: max-height 0.2s ease;  /* 缩短动画时间 */
}
/* 优化移动端滚动 */
@media screen and (max-width: 768px) {
  html, body {
    -webkit-overflow-scrolling: touch;
    overflow-scrolling: touch;
  }
}
/* 添加移动端适配 */
@media screen and (max-width: 768px) {
  .content-container {
    margin: 2px 0;
    padding: 8px;
    box-shadow: none;
    backdrop-filter: none;
    -webkit-backdrop-filter: none;
  }
  
  /* 调整内容区域的内边距 */
  .border-l-2 {
    padding: 0.8rem !important;
  }
  /* 优化移动端滚动 */
  .message-list-container {
    transform: translate3d(0, 0, 0);
    -webkit-overflow-scrolling: touch;
  }
  .content-container img {
    min-height: 100px;
  }
  /* 移除移动端动画效果 */
  .message-actions > div {
    transition: none;
  }
}
.content-container::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  z-index: -1;
  border-radius: inherit;
}

/* 添加展开/折叠按钮样式 */
button {
  background: rgba(36, 43, 50, 0.95);
  border: 1px solid rgba(251, 146, 60, 0.3);
  border-radius: 20px;
  position: relative;
  z-index: 9999;
}

button:hover {
  background: rgba(46, 53, 60, 0.95);
  border-color: rgba(251, 146, 60, 0.5);
  cursor: pointer;
}

/* 确保内容区域的层级正确 */
.overflow-y-hidden {
  transition: max-height 0.3s ease-in-out;
  position: relative;
  z-index: 1;
}
/* 添加内容过渡动画 */
.overflow-y-hidden {
  transition: max-height 0.3s ease-in-out;
}

/* 修正展开状态下的最大高度限制 */
.content-container .overflow-y-hidden:not(.max-h-\[700px\]) {
  max-height: none;
}
/* 添加页脚固定样式 */
:deep(.text-center.text-xs.text-gray-400.py-4) {
  margin-top: auto;
  padding-top: 2rem;
}
/* 修改评论区样式 */
:deep(.wl-comment) {
  background: rgba(36, 43, 50, 0.95) !important;
  border-radius: 8px;
  padding: 12px !important;
  margin-bottom: 12px !important;
}
/* 修改输入框文本颜色为黑色 */
:deep(.wl-input) {
  color: #fff !important;  /* 输入文字颜色保持白色 */
  background-color: rgba(36, 43, 50, 0.95) !important; /* 输入框背景改为深色 */
  border-color: rgba(251, 146, 60, 0.3) !important; /* 边框颜色调整 */
}
:deep(.wl-input::placeholder) {
  color: rgba(255, 255, 255, 0.5) !important; /* placeholder文字颜色 */
}

:deep(.wl-editor) {
  background: rgba(36, 43, 50, 0.95) !important;
  color: #fff !important;
}

:deep(.wl-comment .wl-content) {
  color: #fff !important;
  background: transparent !important;
}

/* 确保评论内容为白色 */
:deep(.wl-content),
:deep(.wl-content p),
:deep(.wl-content *) {
  color: #fff !important;
}
/* 调整编辑器区域样式 */
:deep(.wl-comment .wl-meta .wl-like),
:deep(.wl-comment .wl-meta .wl-reply) {
  color: #999 !important;
}


/* 调整输入框边框 */
:deep(.wl-input-row) {
  border-color: rgba(0, 0, 0, 0.1) !important;
}
:deep(.wl-comment .wl-meta .wl-like:hover),
:deep(.wl-comment .wl-meta .wl-reply:hover) {
  color: #fff !important;
}

/* 确保所有评论相关文本为白色 */
:deep(.wl-comment *) {
  color: #fff !important;
}
/* 调整按钮样式 */
:deep(.wl-btn) {
  background-color: rgba(251, 146, 60, 0.8) !important;
  color: #fff !important;
}

:deep(.wl-action) {
  color: #fff !important;
}

:deep(.wl-header) {
  border-bottom: 1px solid rgba(14, 14, 14, 0.2) !important;
}

:deep(.wl-card) {
  background: rgba(36, 43, 50, 0.95) !important;
  border: 1px solid rgba(14, 14, 14, 0.2) !important;
}
/* 添加评论框样式 */
:deep(.wl-panel),
:deep(.wl-card) {
  position: relative;
  z-index: 100;
  background: rgba(36, 43, 50, 0.95) !important;
  border: 1px solid rgba(14, 14, 14, 0.2) !important;
}

/* 确保评论区域不会被遮挡 */
.content-container {
  position: relative;
  z-index: 1;
}
/* 添加评论内容文本颜色 */
:deep(.wl-comment .wl-content) {
  color: #fff !important;
}

:deep(.wl-comment .wl-meta) {
  color: #fff !important;
}

:deep(.wl-comment .wl-meta > span),
:deep(.wl-comment .wl-meta > a) {
  color: #fff !important;
}
/* 移除 markdown 图片的 hover 效果 */
:deep(.markdown-preview img) {
  cursor: pointer;
  transform: none !important; /* 移除 hover 时的缩放效果 */
  transition: none !important; /* 移除过渡效果 */
}

:deep(.markdown-preview img:hover) {
  transform: none !important;
}

/* 确保灯箱层级最高 */
:deep(.fancybox__container) {
  --fancybox-bg: rgba(0, 0, 0, 0.9);
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  z-index: 9999 !important;
}

:deep(.fancybox__backdrop) {
  z-index: 9998 !important;
}
/* 按钮组样式 */
.message-actions {
  position: relative;
  z-index: 1;
}

/* 按钮悬停效果 */
.message-actions > div {
  position: relative;
  transition: all 0.3s ease;
}

.message-actions > div:hover {
  transform: translateY(-2px);
}

.message-actions > div:hover .text-gray-400 {
  color: #fb923c;
  filter: drop-shadow(0 0 2px rgba(251, 146, 60, 0.3));
}
.gradient-dot {
  /* 添加明亮色彩的动态渐变动画 */
  background: linear-gradient(
    45deg,
    #ff6b6b,
    #ffd93d,
    #ff9a9e,
    #cd4e67,
    #ffb347,
    #ff7eb3,
    #ffa07a
  );
  background-size: 400% 400%;
  animation: rainbow 10s ease infinite;
  -webkit-background-clip: text;
  background-clip: text;
  color: transparent;
  font-weight: bold;
}

@keyframes rainbow {
  0% {
    background-position: 0% 50%;
  }
  50% {
    background-position: 100% 50%;
  }
  100% {
    background-position: 0% 50%;
  }
}

/* 隐藏滚动条但保持功能 */
.hide-scrollbar {
  -ms-overflow-style: none;
  scrollbar-width: none;
}
.hide-scrollbar::-webkit-scrollbar {
  display: none;
}
/* ... 跳转页文本 ... */
.text-shadow-sm {
  text-shadow: 0 1px 2px rgba(0, 0, 0, 0.1),
               0 2px 4px rgba(0, 0, 0, 0.1);
  font-weight: 500;
  letter-spacing: 0.5px;
}
/* 添加移动端分页按钮适配 */
@media screen and (max-width: 768px) {
  .UButton {
    font-size: 0.875rem;
    padding: 0.375rem 0.75rem;
  }
  
  .UInput {
    height: 2rem;
    font-size: 0.875rem;
  }
  
  /* 调整按钮间距 */
  .space-x-4 > * + * {
    margin-left: 0.5rem;
  }
  
  /* 优化移动端分页布局 */
  .flex-wrap {
    flex-wrap: wrap;
  }
  
  .mt-3 {
    margin-top: 0.75rem;
  }
}
/* 添加高亮动画样式 */
@keyframes highlight {
  0% { background: rgba(251, 146, 60, 0.3); }
  100% { background: rgba(36, 43, 50, 0.95); }
}

.highlight-message {
  animation: highlight 2s ease-out;
}

/* 轻模式覆盖 Markdown 颜色 */
.content-container.text-black :deep(.markdown-preview h1),
.content-container.text-black :deep(.markdown-preview h2),
.content-container.text-black :deep(.markdown-preview h3),
.content-container.text-black :deep(.markdown-preview h4),
.content-container.text-black :deep(.markdown-preview h5),
.content-container.text-black :deep(.markdown-preview h6) {
  color: #111 !important;
}

.content-container.text-black :deep(.markdown-preview p),
.content-container.text-black :deep(.markdown-preview li),
.content-container.text-black :deep(.markdown-preview span:not(.clickable-tag)) {
  color: #333 !important;
}

/* 白天模式下链接颜色与悬停颜色 */
.content-container.text-black :deep(.markdown-preview a) {
  color: #0366d6 !important;
  text-decoration: none;
}
.content-container.text-black :deep(.markdown-preview a:hover) {
  color: #1d4ed8 !important;
  text-decoration: underline;
}

.content-container.text-black :deep(pre) {
  background-color: #f5f5f5 !important;
  border: 1px solid #e5e7eb !important;
  color: #1f2937 !important;
}

.content-container.text-black :deep(.hljs) {
  color: #1f2937 !important;
}
</style>
