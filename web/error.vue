<script setup lang="ts">
const router = useRouter()
</script>

<template>
  <div class="fixed inset-0 bg-gradient-to-br from-slate-900 via-indigo-950 to-slate-800">
    <div class="absolute inset-0 backdrop-blur-xl"></div>
    <div class="relative z-10 flex min-h-screen items-center justify-center p-4">
      <UCard class="w-full max-w-md bg-slate-900/70 text-white border border-slate-700/40 shadow-2xl">
        <div class="flex items-center gap-3 mb-4">
          <UIcon name="i-heroicons-exclamation-triangle" class="w-7 h-7 text-orange-300" />
          <div class="text-xl font-bold">404 页面未找到</div>
        </div>
        <div class="text-slate-300 mb-4">你访问的页面不存在或已被移动。</div>
        <div class="flex gap-2">
          <UButton color="primary" class="flex-1" to="/">返回首页</UButton>
          <UButton color="gray" variant="soft" class="flex-1" @click="router.back()">上一页</UButton>
        </div>
      </UCard>
    </div>
  </div>
</template>
