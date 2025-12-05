<template>
  <div :class="containerClass">

    <div class="editor-box">
      <VditorEditor ref="vditorEditor" v-model="MessageContent" :theme="contentTheme" />
      <div class="editor-toolbar">
        <div class="toolbar-left">
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
          <button class="tb-btn" @click="triggerFileInput" title="插入图片"><UIcon name="i-fluent-image-20-regular" class="w-5 h-5" /></button>
           <!-- 新增图床上传按钮 -->
           <button class="tb-btn" @click="showImageUploader = true" title="图床上传"><UIcon name="i-mdi-cloud-upload-outline" class="w-5 h-5" /></button>
          
          <button class="tb-btn" @click="togglePrivate" :title="Private ? '设为公开' : '设为私密'">
            <UIcon :name="privateIcon" class="w-5 h-5" />
          </button>
          <button class="tb-btn" @click="toggleNotify" :title="enableNotify ? '关闭推送' : '开启推送'">
            <UIcon :name="enableNotify ? 'i-mdi-bell' : 'i-mdi-bell-off'" class="w-5 h-5" />
          </button>          
        </div>
        <div class="toolbar-right">
          
          <button class="tb-btn" @click="clearForm" title="清空"><UIcon name="i-fluent-broom-16-regular" class="w-5 h-5" /></button>
          <button class="tb-btn primary" @click="addMessage" title="发布"><UIcon name="i-fluent-add-12-filled" class="w-5 h-5" /></button>
        </div>
      </div>
    </div>
  </div>

  <!-- 内容预览区域 - 仅在有内容时显示 -->
  <div v-if="MessageContentHtml" class="mx-auto sm:max-w-2xl mt-4 preview-card">
    <div :class="[previewProseClass, 'max-w-none editor-preview']">
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
const BASE_API = useRuntimeConfig().public.baseApi || '/api';
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
const toggleTheme = () => {
  toggleContentTheme && toggleContentTheme()
  if (typeof window !== 'undefined') {
    document.documentElement.classList.toggle('dark', contentTheme.value === 'dark')
  }
}
const fileInput = ref<HTMLInputElement | null>(null);
const vditorEditor = ref<any>(null); // 需要支持 insertValue
// 预览跟随内容自动显示，无需手动开关

const privateIcon = computed(() => (Private.value ? 'i-mdi-eye-off-outline' : 'i-mdi-eye-outline'));
const previewProseClass = computed(() => contentTheme.value === 'dark' ? 'prose prose-invert' : 'prose')

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
      if (vditorEditor.value?.insertValue) {
        const origin = typeof window !== 'undefined' ? window.location.origin : ''
        const base = String(BASE_API || '/api')
        const ret = String(data.data || '')
        let full = ''
        if (ret.startsWith('http')) {
          full = ret
        } else {
          const path = ret.startsWith('/') ? ret : `/${ret}`
          if (/^https?:\/\//.test(base)) {
            full = `${base}${path}`
          } else {
            full = `${origin}${base}${path}`
          }
        }
        const imageMarkdown = `\n![](${full})\n`
        vditorEditor.value.insertValue(imageMarkdown)
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
  const raw = String(videoUrl || '')
  const baseApi = useRuntimeConfig().public.baseApi || '/api'
  let full = raw
  if (!/^https?:\/\//.test(raw)) {
    const path = raw.startsWith('/') ? raw : `/${raw}`
    if (/^https?:\/\//.test(String(baseApi))) {
      const base = String(baseApi).replace(/\/api$/, '')
      full = `${base}${path}`
    } else {
      const origin = typeof window !== 'undefined' ? window.location.origin : ''
      full = `${origin}${path}`
    }
  }
  const videoTag = `<video width="100%" height="100%" src="${full}" controls loop></video>\n`
  if (vditorEditor.value?.insertValue) {
    vditorEditor.value.insertValue(videoTag)
  }
};

const INLINE_IMAGE_REG = /!\s*(https?:\/\/[^\s!]+\.(?:png|jpe?g|gif|webp))(?:\?[^\s!]*)?/gi;
const normalizeInlineImageLinks = (md: string): string => md.replace(INLINE_IMAGE_REG, (m, url) => `![](${url})`);

const applyImageGridHTML = (html: string) => {
  const parser = new DOMParser();
  const doc = parser.parseFromString(html, 'text/html');
  const isPureImageParagraph = (p: Element) => {
    let ok = true;
    const children = Array.from(p.childNodes);
    if (children.length === 0) return false;
    for (const node of children) {
      if (node.nodeType === Node.ELEMENT_NODE) {
        const el = node as Element;
        const tag = el.tagName.toLowerCase();
        if (tag === 'img') continue;
        if (tag === 'a' && el.childElementCount === 1 && el.querySelector('img')) continue;
        if (tag === 'br') { ok = false; break; }
        ok = false; break;
      } else if (node.nodeType === Node.TEXT_NODE) {
        if ((node.textContent || '').trim() !== '') { ok = false; break; }
      }
    }
    return ok;
  };

  const paras = Array.from(doc.body.querySelectorAll('p'));
  const runs: Element[][] = [];
  let current: Element[] = [];
  for (const p of paras) {
    if (isPureImageParagraph(p)) {
      const last = current[current.length - 1];
      if (!last || last.nextElementSibling === p) {
        current.push(p);
      } else {
        if (current.length >= 2) runs.push(current);
        current = [p];
      }
    } else {
      if (current.length >= 2) runs.push(current);
      current = [];
    }
  }
  if (current.length >= 2) runs.push(current);

  for (const run of runs) {
    const grid = doc.createElement('div');
    const count = run.length;
    const cols = count === 2 || count === 4 ? 2 : Math.min(3, count);
    grid.className = `image-grid cols-${cols}`;
    const group = `grid-${Math.random().toString(36).slice(2)}`;
    for (const p of run) {
      const img = p.querySelector('img') as HTMLImageElement | null;
      const a = p.querySelector('a') as HTMLAnchorElement | null;
      if (!img && !a) continue;
      const item = doc.createElement('div');
      item.className = 'image-grid-item';
      let anchor: HTMLAnchorElement;
      if (a && a.querySelector('img')) {
        anchor = a;
        anchor.setAttribute('data-fancybox', group);
        if (!anchor.getAttribute('href')) {
          const innerImg = a.querySelector('img') as HTMLImageElement;
          anchor.setAttribute('href', innerImg.src);
        }
      } else if (img) {
        anchor = doc.createElement('a');
        anchor.setAttribute('href', img.src);
        anchor.setAttribute('data-fancybox', group);
        anchor.appendChild(img);
      } else {
        continue;
      }
      item.appendChild(anchor);
      grid.appendChild(item);
    }
    const first = run[0];
    first.replaceWith(grid);
    for (let i = 1; i < run.length; i++) run[i].remove();
  }
  return doc.body.innerHTML;
};

watch(MessageContent, (val) => {
  const raw = Vditor.md2html(normalizeInlineImageLinks(val || ""));
  MessageContentHtml.value = applyImageGridHTML(raw);
  nextTick(() => {
    const roots = document.querySelectorAll('.editor-preview');
    roots.forEach((root) => {
      root.querySelectorAll('.image-grid-item img').forEach((imgEl) => {
        const img = imgEl as HTMLImageElement;
        const parent = img.parentElement as HTMLElement;
        const setAR = () => {
          const w = img.naturalWidth;
          const h = img.naturalHeight;
          parent.classList.remove('ar-169','ar-34','ar-11');
          if (w > h) parent.classList.add('ar-169');
          else if (h > w) parent.classList.add('ar-34');
          else parent.classList.add('ar-11');
        };
        if (img.complete && img.naturalWidth && img.naturalHeight) setAR();
        else img.addEventListener('load', setAR, { once: true });
      });
    });
  });
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

<style scoped>
.editor-box { background: #ffffff; border: 1px solid #e5e7eb; border-radius: 12px; box-shadow: 0 10px 24px rgba(0,0,0,.08); padding: 8px; color:#111827; }
.editor-toolbar { display:flex; align-items:center; justify-content:space-between; gap:8px; margin-top:6px; padding:6px; border-radius:12px; background: rgba(255,255,255,0.85); flex-wrap: wrap; overflow:hidden; }
.toolbar-left, .toolbar-right { display:flex; align-items:center; gap:8px; flex-wrap: wrap; }
.tb-btn { display:flex; align-items:center; justify-content:center; width:36px; height:36px; border-radius:12px; background: rgba(0,0,0,0.06); color:#374151; transition: all .18s ease; border:none; }
.tb-btn:hover { transform: translate3d(0,0,0) scale(1.06); background: rgba(0,0,0,0.12); }
.tb-btn.primary { background: linear-gradient(135deg, rgba(251,146,60,.95), rgba(234,88,12,.95)); color: #fff; }
.tb-sep { width:1px; height:24px; background: rgba(0,0,0,0.12); margin: 0 2px; }
.preview-card { backdrop-filter: blur(8px); background: #ffffff; border: 1px solid #e5e7eb; border-radius: 12px; padding: 8px; color:#111827; }
html.dark .editor-box { background: rgba(36,43,50,0.95); border: 1px solid rgba(255,255,255,0.08); color:#fff; }
html.dark .editor-toolbar { background: rgba(36,43,50,0.6); }
html.dark .tb-btn { background: rgba(255,255,255,0.06); color:#cbd5e1; border:none; }
html.dark .tb-btn:hover { background: rgba(255,255,255,0.12); }
html.dark .tb-sep { background: rgba(255,255,255,0.12); }
html.dark .preview-card { background: rgba(36,43,50,0.6); border: 1px solid rgba(255,255,255,0.12); color:#fff; }
.editor-toolbar :deep(.u-button) { border:none !important; box-shadow:none !important; background: transparent !important; color:#374151 !important; }
html.dark .editor-toolbar :deep(.u-button) { border:none !important; box-shadow:none !important; background: rgba(255,255,255,0.06) !important; color:#cbd5e1 !important; }
.editor-preview p { margin: 0.5rem 0; }
.editor-preview img { margin: 0.4rem 0; }
.image-grid {
  display: grid;
  gap: 6px;
  margin: 0;
  width: 100%;
  grid-auto-flow: dense;
  align-items: stretch;
  justify-items: stretch;
}
.image-grid.cols-2 { grid-template-columns: repeat(2, 1fr); }
.image-grid.cols-3 { grid-template-columns: repeat(3, 1fr); }
.image-grid-item {
  position: relative;
  aspect-ratio: 1 / 1;
  overflow: hidden;
  border-radius: 10px;
}
.image-grid-item > a { display: block; width: 100%; height: 100%; }
.image-grid-item > a > img { width: 100%; height: 100%; object-fit: cover; object-position: center; display: block; }
.image-grid-item.ar-169 { aspect-ratio: 16 / 9; }
.image-grid-item.ar-34 { aspect-ratio: 3 / 4; }
.image-grid-item.ar-11 { aspect-ratio: 1 / 1; }
.image-grid-item img {
  width: 100%;
  height: 100%;
  object-fit: cover;
  display: block;
  margin: 0;
}
</style>

const applyFormat = (type: string) => {
  const ins = (val: string) => vditorEditor.value?.insertValue && vditorEditor.value.insertValue(val)
  switch (type) {
    case 'bold': ins('**加粗**'); break
    case 'italic': ins('*斜体*'); break
    case 'strike': ins('~~删除线~~'); break
    case 'inlineCode': ins('`code`'); break
    case 'codeBlock': ins('\n```\n// code\n```\n'); break
    case 'quote': ins('\n> 引用\n'); break
    case 'ul': ins('\n* 列表项\n'); break
    case 'ol': ins('\n1. 列表项\n'); break
    case 'task': ins('\n- [ ] 待办\n'); break
    case 'hr': ins('\n---\n'); break
    case 'link': ins('[链接文本](https://)'); break
  }
}
const props = defineProps<{ wide?: boolean }>()
const containerClass = computed(() => (props.wide ? 'w-full max-w-none' : 'mx-auto sm:max-w-2xl'))
