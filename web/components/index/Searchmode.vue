<template>
  <!-- 搜索内容显示 -->
  <UModal v-model="showModal">
    <div class="bg-gray-800 p-6 rounded-lg">
      <h3 class="text-xl font-semibold text-white mb-4">搜索内容</h3>
      <div class="space-y-4">
        <UInput
          v-model="searchQuery"
          placeholder="输入搜索关键词"
          class="w-full"
          @keyup.enter="handleSearch"
        />
        <div class="flex justify-end gap-2">
          <UButton
            variant="solid"
            color="gray"
            @click="closeModal"
            class="bg-gray-700 hover:bg-gray-700 text-white hover:text-white border border-gray-600"
          >
            取消
          </UButton>
          <UButton
            color="primary"
            @click="handleSearch"
            class="text-white hover:text-white"
          >
            搜索
          </UButton>
        </div>
      </div>
    </div>
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
.modal-content {
  background-color: rgba(36, 43, 50, 0.95);
  border-radius: 16px;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.15);
}
</style>
