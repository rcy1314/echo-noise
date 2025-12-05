<template>
  <div 
    class="fixed z-[5000] bg-black/80 backdrop-blur-sm rounded-lg shadow-lg p-4 text-white"
    :style="getPopupPosition"
    @mousedown="startDrag"
    @touchstart="startDrag"
  >
    <div class="flex justify-between items-center mb-3">
      <h3 class="text-lg font-medium">图床上传</h3>
      <UButton 
        icon="i-heroicons-x-mark" 
        color="gray" 
        variant="ghost" 
        size="xs" 
        @click="$emit('close')" 
      />
    </div>
    
    <div class="mb-4">
      <div class="flex gap-2 mb-3">
        <UButton 
          :color="selectedHost === 'github' ? 'blue' : 'gray'" 
          variant="solid" 
          size="sm" 
          @click="selectedHost = 'github'"
        >
          GitHub
        </UButton>
      </div>
      
      <!-- GitHub 配置 -->
      <div v-if="selectedHost === 'github'" class="space-y-2 mb-3">
        <UInput 
          v-model="githubToken" 
          placeholder="GitHub Token" 
          size="sm"
          type="password"
        />
        <UInput 
          v-model="githubRepo" 
          placeholder="仓库名 (用户名/仓库)" 
          size="sm"
        />
        <UInput 
          v-model="githubBranch" 
          placeholder="分支名 (默认: main)" 
          size="sm"
        />
        <UInput 
          v-model="githubPath" 
          placeholder="存储路径 (默认: images/)" 
          size="sm"
        />
        <!-- 修改后的CDN配置 -->
        <div class="flex items-center gap-2">
          <UToggle v-model="enableCDN" />
          <span class="text-sm">启用CDN加速</span>
        </div>
        <UInput 
          v-if="enableCDN"
          v-model="cdnDomain" 
          placeholder="输入CDN域名 (如: jsd.cdn.noisework.cn)" 
          size="sm"
        />
      </div>
      
      <!-- 保存配置按钮 - 移到这里，让它在任何配置下都显示 -->
      <div class="flex justify-end mt-2 mb-3">
        <UButton 
          color="green" 
          variant="solid" 
          size="sm" 
          icon="i-heroicons-check" 
          @click="saveConfig"
        >
          保存配置
        </UButton>
      </div>
      
      <div class="flex flex-col items-center justify-center border-2 border-dashed border-gray-500 rounded-lg p-4 cursor-pointer hover:border-blue-400 transition-colors"
           @click="triggerFileInput"
           @dragover.prevent="isDragging = true"
           @dragleave.prevent="isDragging = false"
           @drop.prevent="handleFileDrop"
           :class="{ 'border-blue-400': isDragging }"
      >
        <input
          ref="fileInput"
          type="file"
          accept="image/*"
          @change="handleFileSelect"
          class="hidden"
        />
        <UIcon 
          v-if="!isUploading && !previewUrl" 
          name="i-heroicons-cloud-arrow-up" 
          class="w-10 h-10 text-gray-400 mb-2" 
        />
        <img 
          v-if="previewUrl && !isUploading" 
          :src="previewUrl" 
          class="max-h-32 max-w-full mb-2 rounded" 
        />
        <UProgress 
          v-if="isUploading" 
          :value="uploadProgress" 
          color="blue" 
          class="w-full mb-2" 
        />
        <p v-if="!isUploading && !previewUrl" class="text-sm text-gray-300">
          点击或拖拽图片到此处上传
        </p>
        <p v-if="isUploading" class="text-sm text-gray-300">
          上传中... {{ uploadProgress }}%
        </p>
      </div>
      
      <!-- 添加明确的上传按钮 -->
      <div v-if="previewUrl && !isUploading && !uploadedUrl" class="flex justify-end mt-2">
        <UButton 
          color="blue" 
          variant="solid" 
          size="sm" 
          icon="i-heroicons-cloud-arrow-up" 
          @click="startUpload"
        >
          开始上传
        </UButton>
      </div>
      
      <div v-if="errorMessage" class="mt-2 text-red-400 text-sm">
        {{ errorMessage }}
      </div>
         
      <div v-if="uploadedUrl" class="mt-2">
        <div class="flex items-center gap-2 bg-gray-800 p-2 rounded">
          <input 
            type="text" 
            :value="uploadedUrl" 
            readonly 
            class="bg-transparent flex-1 text-sm outline-none"
          />
          <UButton 
            icon="i-heroicons-clipboard" 
            color="gray" 
            variant="ghost" 
            size="xs" 
            @click="copyToClipboard(uploadedUrl)" 
          />
        </div>
        <div class="flex justify-end mt-2">
          <UButton 
            color="blue" 
            variant="solid" 
            size="sm" 
            @click="insertImage" 
          >
            插入图片
          </UButton>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch, onMounted } from 'vue';
import { useToast, useRuntimeConfig } from '#imports';

// 修改 position 属性定义
const props = defineProps({
  position: {
    type: Object,
    required: true,
    default: () => ({ x: 0, y: 0 })
  },
  editorRef: {
    type: Object,
    required: false,
    default: null
  }
});

// 添加弹窗位置计算
const getPopupPosition = computed(() => {
  if (typeof window === 'undefined') return {};
  
  const screenWidth = window.innerWidth;
  const screenHeight = window.innerHeight;
  const popupWidth = 400;
  const popupHeight = 500;
  
  let x = props.position.x;
  let y = props.position.y;
  
  // 确保弹窗在屏幕内
  if (x + popupWidth > screenWidth) {
    x = screenWidth - popupWidth - 20;
  }
  if (x < 0) x = 20;
  
  if (y + popupHeight > screenHeight) {
    y = screenHeight - popupHeight - 20;
  }
  if (y < 0) y = 20;
  
  // 移动端/平板居中显示
  if (screenWidth < 1024) {
    x = (screenWidth - popupWidth) / 2;
    y = 60;
  }
  
  return {
    position: 'fixed',
    top: `${y}px`,
    left: `${x}px`,
    maxWidth: `${Math.min(popupWidth, screenWidth - 40)}px`,
    width: '100%',
    // 添加以下样式确保可拖动
    cursor: 'move',
    userSelect: 'none',
    touchAction: 'none'
  };
});
const isDraggingWindow = ref(false);
const dragStartPos = ref({ x: 0, y: 0 });

const startDrag = (e: MouseEvent | TouchEvent) => {
  isDraggingWindow.value = true;
  const clientX = 'touches' in e ? e.touches[0].clientX : e.clientX;
  const clientY = 'touches' in e ? e.touches[0].clientY : e.clientY;
  
  dragStartPos.value = {
    x: clientX - props.position.x,
    y: clientY - props.position.y
  };

  document.addEventListener('mousemove', handleDrag);
  document.addEventListener('touchmove', handleDrag);
  document.addEventListener('mouseup', stopDrag);
  document.addEventListener('touchend', stopDrag);
};

const handleDrag = (e: MouseEvent | TouchEvent) => {
  if (!isDraggingWindow.value) return;
  
  const clientX = 'touches' in e ? e.touches[0].clientX : e.clientX;
  const clientY = 'touches' in e ? e.touches[0].clientY : e.clientY;
  
  emit('update:position', {
    x: clientX - dragStartPos.value.x,
    y: clientY - dragStartPos.value.y
  });
};

const stopDrag = () => {
  isDraggingWindow.value = false;
  document.removeEventListener('mousemove', handleDrag);
  document.removeEventListener('touchmove', handleDrag);
  document.removeEventListener('mouseup', stopDrag);
  document.removeEventListener('touchend', stopDrag);
};

const emit = defineEmits(['close', 'upload-success']);
const toast = useToast();

// 状态变量
const selectedHost = ref('smms');
const fileInput = ref<HTMLInputElement | null>(null);
const isDragging = ref(false);
const isUploading = ref(false);
const uploadProgress = ref(0);
const previewUrl = ref('');
const uploadedUrl = ref('');
const errorMessage = ref('');

// GitHub 配置
const githubToken = ref('');
const githubRepo = ref('');
const githubBranch = ref('main');
const githubPath = ref('images/');
const enableCDN = ref(false);  // 默认不启用CDN
const cdnDomain = ref('');     // 默认没有CDN域名

// 图床配置
const smmsToken = ref('');

const freeimageToken = ref('');

// 触发文件选择
const triggerFileInput = () => {
  fileInput.value?.click();
};

// 处理文件拖放
const handleFileDrop = (event: DragEvent) => {
  isDragging.value = false;
  const files = event.dataTransfer?.files;
  if (files && files.length > 0) {
    handleFile(files[0]);
  }
};

// 处理文件选择
const handleFileSelect = (event: Event) => {
  const input = event.target as HTMLInputElement;
  const files = input.files;
  if (files && files.length > 0) {
    handleFile(files[0]);
  }
};

// 添加一个变量存储当前选择的文件
const selectedFile = ref<File | null>(null);

// 处理文件
// 修改文件处理方法
const handleFile = async (file: File) => {
  // 检查文件类型
  if (!file.type.startsWith('image/')) {
    errorMessage.value = '请选择图片文件';
    return;
  }
  
  // 检查文件大小 (限制为 5MB)
  const maxSize = 5 * 1024 * 1024;
  if (file.size > maxSize) {
    errorMessage.value = '图片大小不能超过 5MB';
    return;
  }
  
  // 清除错误信息
  errorMessage.value = '';
  selectedFile.value = file;
  
  // 创建预览
  const reader = new FileReader();
  reader.onload = (e) => {
    previewUrl.value = e.target?.result as string;
  };
  reader.readAsDataURL(file);
};

// 修改插入图片方法
const insertImage = () => {
  if (!uploadedUrl.value) {
    errorMessage.value = '请先上传图片';
    return;
  }
  
  const markdownLink = `![](${uploadedUrl.value})`;
  emit('upload-success', markdownLink);
  emit('close');
  
  // 重置状态
  selectedFile.value = null;
  previewUrl.value = '';
  uploadedUrl.value = '';
  uploadProgress.value = 0;
  isUploading.value = false;
  errorMessage.value = '';
};

// 添加复制到剪贴板方法
const copyToClipboard = async (text: string) => {
  try {
    await navigator.clipboard.writeText(text);
    toast.add({
      title: '成功',
      description: '已复制到剪贴板',
      color: 'green',
      timeout: 2000
    });
  } catch (error) {
    toast.add({
      title: '错误',
      description: '复制失败',
      color: 'red',
      timeout: 2000
    });
  }
};

// 修改startUpload方法，添加全局错误处理
const startUpload = async () => {
  if (!selectedFile.value) {
    errorMessage.value = '请先选择图片';
    return;
  }

  isUploading.value = true;
  uploadProgress.value = 0;
  errorMessage.value = '';
  
  try {
    switch(selectedHost.value) {
      case 'github':
        await uploadToGitHub(selectedFile.value);
        break;
      default:
        throw new Error('未知的图床类型');
    }
  } catch (error: any) {
    console.error('上传错误:', error);
    errorMessage.value = error.message || '上传失败';
    isUploading.value = false;
    
    // 添加详细错误诊断
    if (error.name === 'AbortError') {
      errorMessage.value = '请求超时，请检查网络连接';
    } else if (error.message.includes('Failed to fetch')) {
      errorMessage.value = '网络请求失败，可能是跨域问题或API不可用';
    }
  }
};
// 添加 GitHub 上传方法
const uploadToGitHub = async (file: File) => {
  if (!githubToken.value || !githubRepo.value) {
    errorMessage.value = '请先配置 GitHub Token 和仓库名';
    return;
  }

  try {
    uploadProgress.value = 30;
    
    // 生成文件名和路径
    const timestamp = Date.now();
    const fileExt = file.name.split('.').pop();
    const fileName = `${timestamp}.${fileExt}`;
    const filePath = `${githubPath.value}${fileName}`.replace(/\/\//g, '/');
    
    // 读取文件内容
    const fileContent = await readFileAsBase64(file);
    
    // 准备请求数据
    const requestData = {
      message: `Upload image ${fileName}`,
      content: fileContent,
      branch: githubBranch.value || 'main'
    };
    
    uploadProgress.value = 50;
    
    // 发送请求
    const response = await fetch(`https://api.github.com/repos/${githubRepo.value}/contents/${filePath}`, {
      method: 'PUT',
      headers: {
        'Authorization': `token ${githubToken.value}`,
        'Content-Type': 'application/json'
      },
      body: JSON.stringify(requestData)
    });
    
    uploadProgress.value = 80;
    
    if (!response.ok) {
      throw new Error(`GitHub 上传失败: ${response.status}`);
    }
    
    const data = await response.json();
    // 只有当启用CDN且填写了CDN域名时才使用CDN URL
    uploadedUrl.value = enableCDN.value && cdnDomain.value
      ? `https://${cdnDomain.value}/gh/${githubRepo.value}@${githubBranch.value || 'main'}/${filePath}`
      : data.content.download_url;
    
    uploadProgress.value = 100;
    isUploading.value = false;
    
    toast.add({
      title: '成功',
      description: `图片已上传到 GitHub${enableCDN.value && cdnDomain.value ? ' (CDN加速)' : ''}`,
      color: 'green',
      timeout: 2000
    });
  } catch (error: any) {
    console.error('GitHub 上传错误:', error);
    errorMessage.value = error.message || 'GitHub 上传失败';
    isUploading.value = false;
  }
};

// 添加文件转 Base64 的辅助方法
const readFileAsBase64 = (file: File): Promise<string> => {
  return new Promise((resolve, reject) => {
    const reader = new FileReader();
    reader.onload = () => resolve((reader.result as string).split(',')[1]);
    reader.onerror = reject;
    reader.readAsDataURL(file);
  });
};



const uploadToFreeImage = async (file: File) => {
  const formData = new FormData();
  formData.append('source', file);
  formData.append('key', '6d207e02198a847aa98d0a2a901485a5');
  formData.append('action', 'upload');
  formData.append('format', 'json');

  try {
    uploadProgress.value = 30;
    
    const controller = new AbortController();
    const timeoutId = setTimeout(() => controller.abort(), 30000); // 延长超时时间
    
    const response = await fetchWithProxy('https://freeimage.host/api/1/upload', {
      method: 'POST',
      body: formData,
      signal: controller.signal,
      headers: {
        'Accept': 'application/json'
      }
    });
    clearTimeout(timeoutId);
    
    uploadProgress.value = 70;
    
    if (!response.ok) {
      const errorData = await response.json().catch(() => ({}));
      throw new Error(errorData.message || `FreeImage上传失败: ${response.status}`);
    }
    
    const data = await response.json();
    if (data.status_code === 200) {
      uploadedUrl.value = data.image.url;
      uploadProgress.value = 100;
      isUploading.value = false;
      
      toast.add({
        title: '成功',
        description: '图片已上传到FreeImage',
        color: 'green',
        timeout: 2000
      });
    } else {
      throw new Error(data.message || 'FreeImage上传失败');
    }
  } catch (error: any) {
    console.error('FreeImage上传错误:', error);
    errorMessage.value = error.name === 'AbortError'
      ? '请求超时，请检查网络连接'
      : error.message || 'FreeImage上传失败，请稍后再试';
    isUploading.value = false;
  }
};
// 保存配置到 localStorage
const saveConfig = () => {
  try {
    const config = {
      selectedHost: selectedHost.value,
      github: {
        token: githubToken.value,
        repo: githubRepo.value,
        branch: githubBranch.value,
        path: githubPath.value,
        enableCDN: enableCDN.value,
        cdnDomain: cdnDomain.value
      },
      smms: {
        token: smmsToken.value
      },
      freeimage: {
        token: freeimageToken.value
      }
    };
    
    localStorage.setItem('imageHostingConfig', JSON.stringify(config));
    
    toast.add({
      title: '成功',
      description: '配置已保存',
      color: 'green',
      timeout: 2000
    });
  } catch (error) {
    console.error('保存配置失败:', error);
    toast.add({
      title: '错误',
      description: '保存配置失败',
      color: 'red',
      timeout: 2000
    });
  }
};

// 从 localStorage 加载配置
const loadConfig = () => {
  const configStr = localStorage.getItem('imageHostingConfig');
  if (configStr) {
    try {
      const config = JSON.parse(configStr);
      selectedHost.value = config.selectedHost || 'smms';
      
      if (config.github) {
        githubToken.value = config.github.token || '';
        githubRepo.value = config.github.repo || '';
        githubBranch.value = config.github.branch || 'main';
        githubPath.value = config.github.path || 'images/';
        enableCDN.value = config.github.enableCDN || false;
        cdnDomain.value = config.github.cdnDomain || ''; // 加载时不设置默认值
      }
      
      if (config.smms) {
        smmsToken.value = config.smms.token || '';
      }
      if (config.imgtp) {
    imgtpToken.value = config.imgtp.token || '';
  }
  if (config.freeimage) {
    freeimageToken.value = config.freeimage.token || '';
  }
    } catch (error) {
      console.error('加载配置失败:', error);
    }
  }
};

// 监听配置变化并自动保存
watch([selectedHost, githubToken, githubRepo, githubBranch, githubPath, smmsToken], 
  () => {
    saveConfig();
  },
  { deep: true }
);

// 组件挂载时加载配置
onMounted(() => {
  loadConfig();
});
</script>

<style scoped>
.fixed {
  position: fixed !important;
}

@media (max-width: 768px) {
  .fixed {
    max-height: 90vh;
    overflow-y: auto;
  }
}
</style>
