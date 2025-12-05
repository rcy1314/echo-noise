<template>
  <div class="tags-container w-full">
    <div class="relative">
      <div class="tags-wrapper">
        <div class="tags-scroll">
          <template v-for="tag in filteredTags" :key="tag.name + timestamp">
            <span
              class="tag-item"
              @click="handleTagClick(tag.name)"
            >
              #{{ tag.name }}
              <span class="tag-count">({{ tag.count }})</span>
            </span>
          </template>
        </div>
      </div>
      <div 
        class="absolute -right-1 top-1/2 -translate-y-1/2 p-2 cursor-pointer transition-all duration-200 hover:scale-110 z-10 refresh-toggle"
        @click="refreshTags"
        title="刷新标签"
      >
        <UIcon 
          name="i-mdi-refresh" 
          class="w-5 h-5 refresh-icon"
          :class="{ 'animate-spin': isRefreshing }"
        />
      </div>
      <div class="scroll-fade"></div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'

const emit = defineEmits(['tagClick', 'updateTags'])
const isRefreshing = ref(false)
const timestamp = ref(Date.now())

const props = defineProps({
  tags: {
    type: Array,
    default: () => []
  }
})

const filteredTags = computed(() => {
  const invalidChars = /[/?=&]/;
  const isMediaLink = /^(song|video|playlist)\?id=\d+$/;
  const cache = new Map();
  
  return props.tags.reduce((acc, tag) => {
    // 使用缓存避免重复计算
    if (cache.has(tag.name)) {
      return acc;
    }
    
    if (!invalidChars.test(tag.name) && !isMediaLink.test(tag.name)) {
      cache.set(tag.name, true);
      acc.push(tag);
    }
    return acc;
  }, []);
}, { immediate: true });

const handleTagClick = (tagName: string) => {
  emit('tagClick', tagName)
}

const refreshTags = async () => {
  if (isRefreshing.value) return
  
  isRefreshing.value = true
  timestamp.value = Date.now()
  emit('updateTags')
  
  setTimeout(() => {
    isRefreshing.value = false
  }, 1000)
}
</script>

<style scoped>
.tags-container { width: 100%; margin: 0; padding: 0; position: relative; background: transparent; --title-color: #d1d5db; }
:global(html.dark) .tags-container { --title-color: #e5e7eb; }

.tags-wrapper {
  position: relative;
  overflow: hidden;
  background: transparent;
}

.tags-scroll {
  display: flex;
  flex-wrap: nowrap;
  gap: 0.4rem;
  overflow-x: auto;
  padding: 0.2rem 0;
  -webkit-overflow-scrolling: touch;
  scrollbar-width: none;
  -ms-overflow-style: none;
  background: transparent;
}

.tags-scroll::-webkit-scrollbar {
  display: none;
}

.tag-item {
  will-change: transform;
  contain: content;
  display: inline-flex;
  align-items: center;
  padding: 0.15rem 0.6rem;
  color: var(--title-color, #d1d5db);
  opacity: .9;
  cursor: pointer;
  transition: all 0.2s ease;
  font-size: 0.875rem;
  white-space: nowrap;
  flex-shrink: 0;
  
}

.tag-item:hover {
  color: #fb923c;
  opacity: 1;
  transform: translateY(-1px);
}

.tag-count {
  margin-left: 0.25rem;
  font-size: 0.75rem;
  opacity: 0.8;
}

.scroll-fade {
  position: absolute;
  right: 0;
  top: 0;
  height: 100%;
  width: 32px;
  pointer-events: none;
}
.refresh-icon { color: var(--title-color, #e5e7eb); filter: drop-shadow(0 0 1px rgba(0,0,0,0.35)); }
.refresh-toggle { opacity: 0; pointer-events: none; }
.tags-container:hover .refresh-toggle { opacity: 1; pointer-events: auto; }
</style>
