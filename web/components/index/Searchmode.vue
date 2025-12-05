<template>
  <!-- 搜索内容显示 -->
  <UModal v-model="showModal" :ui="{ container: 'items-center', base: 'backdrop-blur-sm' }">
    <UCard class="search-card">
      <template #header>
        <div class="flex items-center gap-2">
          <UIcon name="i-heroicons-magnifying-glass" class="w-5 h-5" />
          <h3 class="text-base font-semibold">站内搜索</h3>
        </div>
      </template>
      <div class="space-y-4">
        <div class="relative">
          <UIcon name="i-heroicons-magnifying-glass" class="w-5 h-5 absolute left-3 top-1/2 -translate-y-1/2 opacity-70" />
          <UInput
            v-model="searchQuery"
            placeholder="请输入关键词"
            class="w-full pl-10"
            :ui="{ base: 'rounded-xl' }"
            @keyup.enter="handleSearch"
            autofocus
          />
        </div>
        <div class="flex justify-between items-center">
          <div class="text-xs opacity-70">按 Enter 搜索，Esc 关闭</div>
          <div class="flex justify-end gap-2">
            <UButton variant="ghost" color="gray" @click="closeModal">取消</UButton>
            <UButton color="orange" @click="handleSearch">搜索</UButton>
          </div>
        </div>
      </div>
    </UCard>
  </UModal>  
</template>

<script setup lang="ts">
import { ref, computed } from 'vue';

const toast = useToast();
const config = useRuntimeConfig();
const BASE_API = config.public.baseApi || '/api';

// 添加props和emits以支持v-model
const props = defineProps({
  modelValue: {
    type: Boolean,
    default: false
  }
});

const emit = defineEmits(['update:modelValue', 'search-result']);

// 使用计算属性处理v-model
const showModal = computed({
  get: () => props.modelValue,
  set: (value) => emit('update:modelValue', value)
});

// 关闭模态框
const closeModal = () => {
  emit('update:modelValue', false);
};

const searchQuery = ref('');

// 搜索处理函数
const handleSearch = async () => {
  if (!searchQuery.value.trim()) {
    toast.add({
      title: '提示',
      description: '请输入搜索关键词',
      color: 'yellow'
    });
    return;
  }
  
  try {
    const { data: response, error } = await useFetch('/messages/search', {
      method: 'GET',
      baseURL: BASE_API,
      params: {
        keyword: searchQuery.value,
        page: 1,
        pageSize: 10
      }
    });

    if (error.value) {
      throw new Error(error.value?.message || '网络请求失败');
    }

    if (!response.value) {
      throw new Error('未收到服务器响应');
    }

    if (response.value.code === 1) {
      // 确保发送正确的数据结构
      emit('search-result', response.value);
      emit('update:modelValue', false);
      searchQuery.value = '';
      toast.add({
        title: '成功',
        description: '搜索完成',
        color: 'green'
      });
    } else {
      throw new Error(response.value?.msg || '搜索失败');
    }
  } catch (error) {
    console.error('Search error:', error);
    toast.add({
      title: '错误',
      description: error instanceof Error ? error.message : '搜索失败，请稍后重试',
      color: 'red'
    });
  }
};

// 暴露方法和属性给父组件
defineExpose({
  handleSearch
});
</script>

<style scoped>
.search-card { background: #ffffff; color: #111827; border: 1px solid #e5e7eb; border-radius: 16px; }
html.dark .search-card { background: rgba(36,43,50,0.95); color: #fff; border: 1px solid rgba(255,255,255,0.08); }
</style>
