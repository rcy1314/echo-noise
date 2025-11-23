<template>
  <UCard class="mx-auto sm:max-w-2xl hover:shadow-md backdrop-blur-sm bg-black/40 shadow-lg text-white">
    <div class="flex justify-between mb-3">
      <div class="flex justify-start items-center gap-2">
        <UIcon 
          name="i-heroicons-pencil-square" 
          class="w-6 h-6 transition-all duration-300 hover:scale-110 hover:text-blue-400 animate-pulse" 
        />
        <h2 class="text-lg font-bold italic text-white">说说·笔记~</h2>
      </div>
      <div class="flex gap-2">
        <ClientOnly>
          <div @click="showSearchModal=true" class="cursor-pointer flex">
            <UIcon name="i-heroicons-magnifying-glass" class="w-5 h-5 text-gray-200" />
          </div>
          <button @click="toggleHeatmap">
            <UIcon name="i-mdi-calendar-month" class="w-5 h-5 text-gray-200" />
          </button>
          <button @click="toggleContentTheme && toggleContentTheme()" title="切换暗黑/白天模式">
            <UIcon 
              :name="(contentTheme === 'dark' ? 'i-mdi-weather-night' : 'i-mdi-white-balance-sunny')"
              class="w-5 h-5 text-gray-200"
            />
          </button>
          <a href="/rss" target="_blank">
            <UIcon name="i-mdi-rss" class="w-5 h-5 text-gray-200" />
          </a>
        </ClientOnly>
        <NuxtLink to="/status">
          <UIcon name="i-mdi-server-outline" class="w-5 h-5 text-gray-200" />
        </NuxtLink>
      </div>
    </div>

    <div>
      <VditorEditor ref="vditorEditor" v-model="MessageContent" />
      <div class="flex justify-between items-center">
        <div class="flex items-center justify-start gap-2">
          <input
            id="file-input"
            ref="fileInput"
            type="file"
            accept="image/*"
            @change="addImage"
            class="hidden"
            placeholder="选择图片"
          />
          <!-- 视频上传按钮 -->
          <VideoUpload
    @video-uploaded="handleVideoUploaded"
    @before-upload="checkVideoLogin"
    @upload-progress="handleVideoUploadProgress"
  />
  <div v-if="videoUploadProgress > 0 && videoUploadProgress < 100" class="w-full mt-2">
    <div class="bg-gray-700 rounded h-2">
      <div class="bg-blue-500 h-2 rounded" :style="{ width: videoUploadProgress + '%' }"></div>
    </div>
    <div class="text-xs text-gray-300 mt-1 text-right">{{ videoUploadProgress }}%</div>
  </div>
          <UButton
            color="gray"
            variant="solid"
            class="cursor-pointer"
            size="sm"
            icon="i-fluent-image-20-regular"
            @click="triggerFileInput"
          />
           <!-- 新增图床上传按钮 -->
           <UButton
            color="gray"
            variant="solid"
            class="cursor-pointer"
            size="sm"
            icon="i-mdi-cloud-upload-outline"
            @click="showImageUploader = true"
            title="图床上传"
          />
          <UButton
            color="gray"
            variant="solid"
            size="sm"
            @click="togglePrivate"
            :icon="privateIcon"
            :title="Private ? '设为公开' : '设为私密'"
            :ui="{ tooltip: { text: Private ? '设为公开' : '设为私密' } }"
          />
          <UButton
            color="gray"
            variant="solid"
            size="sm"
            @click="toggleNotify"
            :icon="enableNotify ? 'i-mdi-bell' : 'i-mdi-bell-off'"
            :title="enableNotify ? '关闭推送' : '开启推送'"
            :ui="{ tooltip: { text: enableNotify ? '关闭推送' : '开启推送' } }"
            class="text-gray-600" 
          />          
        </div>
        <div class="flex gap-2">
          <UButton
            icon="i-fluent-broom-16-regular"
            variant="solid"
            color="gray"
            size="sm"
            @click="clearForm"
          />
          <UButton
            icon="i-fluent-add-12-filled"
            variant="solid"
            color="gray"
            size="sm"
            @click="addMessage"
          />
        </div>
      </div>
    </div>
  </UCard>

  <!-- 内容预览区域 - 仅在有内容时显示 -->
  <div
    v-if="MessageContentHtml"
    class="mx-auto sm:max-w-2xl mt-5 backdrop-blur-sm bg-black/40 p-4 rounded-md"
  >
    <div class="prose prose-invert max-w-none">
      <div v-html="MessageContentHtml"></div>
    </div>
    <hr class="border-gray-600 my-4">
    <div v-if="MessageContentHtml" class="prose prose-invert max-w-none">
      <div v-html="MessageContentHtml"></div>
    </div>
  </div>

  <SearchMode 
    v-model="showSearchModal" 
    @search-result="handleSearchResult" 
  />
  <ImageHostingUploader
  v-if="showImageUploader"
  :position="imageUploaderPosition"
  @close="showImageUploader = false"
  @upload-success="handleImageHostingSuccess"
  @update:position="handlePositionUpdate"
/>
</template>

<script setup lang="ts">
import { ref, computed, inject, onMounted, onBeforeUnmount, watch } from 'vue'
import type { MessageToSave } from "~/types/models";
import { UButton } from "#components";
import { useMessage } from "~/composables/useMessage";
import { useUserStore } from '~/store/user'
import { Fancybox } from '@fancyapps/ui'
import '@fancyapps/ui/dist/fancybox/fancybox.css'
import Vditor from 'vditor'
import 'vditor/dist/index.css'
import VditorEditor from './VditorEditor.vue'
import SearchMode from './Searchmode.vue'
import { useMessageStore } from '~/store/message'
import { useNotifyStore } from '~/store/notify'
import VideoUpload from './VideoUpload.vue'
import ImageHostingUploader from '~/components/widgets/ImageHostingUploader.vue'
const showImageUploader = ref(false)
const imageUploaderPosition = ref({ x: 400, y: 320 }) // 可根据实际调整
// 处理图床上传成功，插入编辑器
const handleImageHostingSuccess = (markdown: string) => {
  if (vditorEditor.value?.insertValue) {
    vditorEditor.value.insertValue(markdown)
  }
  showImageUploader.value = false
}
const handlePositionUpdate = (newPosition: { x: number; y: number }) => {
  imageUploaderPosition.value = newPosition;
};
const videoUploadProgress = ref(0); // 新增进度变量
const handleVideoUploadProgress = (percent: number) => {
  videoUploadProgress.value = percent;
};
const showSearchModal = ref(false);
const emit = defineEmits(['search-result','video-uploaded', 'before-upload', 'upload-progress']);
const handleSearchResult = (result: any) => {
  emit('search-result', result);
};
const uploadVideo = (file: File) => {
  const formData = new FormData();
  formData.append('video', file);

  const xhr = new XMLHttpRequest();
  xhr.open('POST', `${BASE_API}/video/upload`, true);
  xhr.withCredentials = true;

  xhr.upload.onprogress = (event) => {
    if (event.lengthComputable) {
      const percent = Math.round((event.loaded / event.total) * 100);
      emit('upload-progress', percent);
    }
  };

  xhr.onload = () => {
    if (xhr.status === 200) {
      const res = JSON.parse(xhr.responseText);
      if (res.code === 1 && res.data) {
        emit('video-uploaded', res.data);
      } else {
        // 错误处理
      }
    } else {
      // 错误处理
    }
    emit('upload-progress', 0); // 上传结束后重置
  };

  xhr.onerror = () => {
    // 错误处理
    emit('upload-progress', 0);
  };

  xhr.send(formData);
};
const toast = useToast()
const BASE_API = useRuntimeConfig().public.baseApi;
const { save } = useMessage();

const showHeatmap = inject('showHeatmap') as Ref<boolean>;
provide('showHeatmap', showHeatmap);

const toggleHeatmap = () => {
  showHeatmap.value = !showHeatmap.value;
};

const Username = ref("");
const MessageContent = ref("");
const MessageContentHtml = ref("");
const Private = ref<boolean>(typeof window !== 'undefined' && localStorage.getItem('postPrivate') === 'true');
const contentTheme = inject('contentTheme') as Ref<string>
const toggleContentTheme = inject('toggleContentTheme') as (() => void) | undefined
const fileInput = ref<HTMLInputElement | null>(null);
const vditorEditor = ref<any>(null); // 需要支持 insertValue

const privateIcon = computed(() => (Private.value ? 'i-mdi-eye-off-outline' : 'i-mdi-eye-outline'));

const notifyStore = useNotifyStore()
const enableNotify = ref(localStorage.getItem('enableNotify') === 'true')

const clearForm = () => {
  Username.value = "";
  MessageContent.value = "";
  MessageContentHtml.value = "";
  
  if (vditorEditor.value) {
    vditorEditor.value.clear();
  }
};

const userStore = useUserStore();

const checkLogin = () => {
  if (!userStore.isLogin) {
    toast.add({
      title: '提示',
      description: '请先登录',
      color: 'orange',
      timeout: 2000
    });
    return false;
  }
  return true;
};

const triggerFileInput = () => {
  const input = document.getElementById("file-input");
  if (input) {
    input.click();
  }
};

const addImage = async (event: Event) => {
  if (!checkLogin()) return;
  const input = event.target as HTMLInputElement;
  const file = input.files ? input.files[0] : null;

  if (!file) {
    toast.add({
      title: '错误',
      description: '没有选择文件',
      color: 'red',
      timeout: 2000
    });
    return;
  }

  const allowedTypes = ['image/jpeg', 'image/png', 'image/gif', 'image/webp'];
  const allowedExtensions = ['.jpg', '.jpeg', '.png', '.gif', '.webp'];
  const fileExtension = file.name.toLowerCase().substring(file.name.lastIndexOf('.'));
  if (!allowedTypes.includes(file.type) || !allowedExtensions.includes(fileExtension)) {
    toast.add({
      title: '错误',
      description: '仅支持 JPG、PNG、GIF、WEBP 格式的图片',
      color: 'red',
      timeout: 2000
    });
    return;
  }
  const maxSize = 5 * 1024 * 1024; // 5MB
  if (file.size > maxSize) {
    toast.add({
      title: '错误',
      description: '图片大小不能超过 5MB',
      color: 'red',
      timeout: 2000
    });
    return;
  }

  try {
    const formData = new FormData();
    formData.append('image', file);

    const response = await fetch(`${BASE_API}/images/upload`, {
      method: 'POST',
      body: formData,
      credentials: 'include'
    });

    if (!response.ok) {
      const errorData = await response.json();
      throw new Error(errorData.msg || '图片上传失败');
    }

    const data = await response.json();
    if (data.code === 1 && data.data) {
      // 直接插入编辑器内容
      if (vditorEditor.value?.insertValue) {
        // 拼接完整图片链接
        const origin = typeof window !== 'undefined' ? window.location.origin : '';
        const imageMarkdown = `\n![](${origin}${BASE_API}${data.data})\n`;
        vditorEditor.value.insertValue(imageMarkdown);
      }
      toast.add({
        title: '成功',
        description: '图片上传成功',
        color: 'green',
        timeout: 2000
      });
    } else {
      throw new Error(data.msg || '图片上传失败');
    }
  } catch (error: any) {
    console.error('上传错误:', error);
    toast.add({
      title: '错误',
      description: error.message || '图片上传失败',
      color: 'red',
      timeout: 2000
    });
  } finally {
    if (fileInput.value) {
      fileInput.value.value = '';
    }
  }
};

const handleVideoUploaded = (videoUrl: string) => {
  const origin = typeof window !== 'undefined' ? window.location.origin : '';
  const url = videoUrl.startsWith('/') ? videoUrl : `/${videoUrl}`;
  const videoTag = `<video width="100%" height="100%" src="${origin}${url}" controls loop></video>\n`;
  if (vditorEditor.value?.insertValue) {
    vditorEditor.value.insertValue(videoTag);
  }
};

watch(MessageContent, (val) => {
  MessageContentHtml.value = Vditor.md2html(val || "");
});

watch(() => userStore.isLogin, (newLoginState) => {
  if (newLoginState) {
    enableNotify.value = localStorage.getItem('enableNotify') === 'true';
  }
}, { immediate: true });

onMounted(async () => {
  Fancybox.bind("[data-fancybox]", {});
  if (!userStore.isLogin) {
    const token = localStorage.getItem('token');
    if (token && userStore.fetchUserInfo) {
      await userStore.fetchUserInfo();
    }
  }
  Private.value = localStorage.getItem('postPrivate') === 'true'
  contentTheme.value = localStorage.getItem('contentTheme') || contentTheme.value
});

onBeforeUnmount(() => {
  Fancybox.destroy();
});
const toggleNotify = () => {
  enableNotify.value = !enableNotify.value;
  localStorage.setItem('enableNotify', enableNotify.value.toString());
};

const togglePrivate = () => {
  Private.value = !Private.value
  localStorage.setItem('postPrivate', Private.value ? 'true' : 'false')
}


const checkVideoLogin = (e: Event) => {
  if (!userStore.isLogin) {
    toast.add({
      title: '提示',
      description: '请登录后操作',
      color: 'orange',
      timeout: 2000
    });
    e.preventDefault && e.preventDefault();
    return false;
  }
  return true;
};

const addMessage = async () => {
  if (!checkLogin()) return;

  if (!MessageContent.value.trim()) {
    toast.add({
      title: '错误',
      description: '请输入内容或上传图片/视频',
      color: 'red',
      timeout: 2000
    });
    return;
  }

  const message: MessageToSave = {
    username: Username.value,
    content: MessageContent.value,
    private: Private.value,
    notify: enableNotify.value,
  };

  try {
    const response = await save(message);
    if (response) {
      clearForm();
    }
  } catch (error: any) {
    console.error('发布错误:', error);
    toast.add({
      title: '错误',
      description: error.message || '发布失败',
      color: 'red',
      timeout: 2000
    });
  }
};
</script>
