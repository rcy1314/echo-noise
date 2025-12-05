<template>
  <div class="announcement-container mx-auto w-full sm:max-w-2xl px-2">
    <div class="announcement-wrapper">
      <UIcon name="i-mdi-bell-outline" class="w-4 h-4 mr-2 flex-shrink-0 text-orange-500" />
      <div class="announcement-text" ref="textContainer" :title="text">
        <div class="marquee animate" ref="marquee" :style="{ '--distance': distance + 'px', animationDuration: duration + 's' }">
          <div class="marquee-item" ref="firstItem">
            <MarkdownRenderer :content="text" :enableGithubCard="false" />
          </div>
          <div class="marquee-item">
            <MarkdownRenderer :content="text" :enableGithubCard="false" />
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import MarkdownRenderer from "~/components/index/MarkdownRenderer.vue"
import { ref, onMounted, onUnmounted, nextTick, watch } from 'vue'
const props = defineProps<{ text: string }>()
const textContainer = ref<HTMLElement | null>(null)
const marquee = ref<HTMLElement | null>(null)
const firstItem = ref<HTMLElement | null>(null)
const distance = ref(0)
const duration = ref(18)
const baseSpeed = 60 // px/s
const measure = () => {
  const item = firstItem.value
  if (!item) return
  const contentWidth = item.scrollWidth
  const gap = 40
  distance.value = contentWidth + gap
  duration.value = Math.max(Math.round(distance.value / baseSpeed), 12)
}
onMounted(() => { nextTick(() => { measure(); window.addEventListener('resize', measure) }) })
onUnmounted(() => { window.removeEventListener('resize', measure) })
watch(() => props.text, () => nextTick(measure))
</script>

<style scoped>
.announcement-container { margin: 0; padding: 0; }
.announcement-wrapper { display: flex; align-items: center; overflow: hidden; padding: 2px 6px; margin: 0; }
.announcement-text { overflow: hidden; white-space: nowrap; font-size: 14px; font-weight: 500; line-height: 1.25; width: 100%; }
.marquee { display: inline-flex; align-items: center; gap: 24px; will-change: transform; }
.marquee-item { display: inline-block; }
.animate { animation-name: marquee; animation-timing-function: linear; animation-iteration-count: infinite; }
@keyframes marquee { 0% { transform: translateX(0); } 100% { transform: translateX(calc(-1 * var(--distance))); } }
:global(html.dark) .announcement-text { color: var(--title-color, #f3f4f6); }
:global(html:not(.dark)) .announcement-text { color: var(--title-color, #111827); }
/* 强制公告文本随标题色变化 */
:deep(.markdown-preview) { display: inline; white-space: nowrap; color: var(--title-color) !important; }
:deep(.markdown-preview p) { display: inline; margin: 0; color: var(--title-color) !important; }
:deep(.markdown-preview h1),
:deep(.markdown-preview h2),
:deep(.markdown-preview h3),
:deep(.markdown-preview h4),
:deep(.markdown-preview h5),
:deep(.markdown-preview h6) { display: inline; margin: 0 0.5rem 0 0; font-size: 1rem; color: var(--title-color) !important; }
:deep(.markdown-preview h1),
:deep(.markdown-preview h2),
:deep(.markdown-preview h3),
:deep(.markdown-preview h4),
:deep(.markdown-preview h5),
:deep(.markdown-preview h6) { display: inline; margin: 0 0.5rem 0 0; font-size: 1rem; }
:deep(.markdown-preview a) { color: var(--title-color, #fb923c); text-decoration: none; }
:deep(.markdown-preview a:hover) { text-decoration: underline; }
</style>
