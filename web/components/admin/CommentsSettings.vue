<template>
  <div id="site-comments-section" class="flex flex-wrap items-center rounded-lg p-3 justify-between gap-3" :class="theme?.subtleBg || subtleBg">
    <div class="flex flex-wrap items-center gap-3 w-full justify-end">
      <USelect v-model="local.commentSystem" :options="[{label:'内置',value:'builtin'},{label:'Waline',value:'waline'}]" />
      <UInput v-if="local.commentSystem === 'waline'" v-model="local.walineServerURL" :ui="{base: theme?.text}" placeholder="Waline 地址" class="w-full md:w-[280px]" />
      <UButton color="green" @click="save" class="shadow">保存</UButton>
    </div>
  </div>
</template>

<script setup lang="ts">
import { reactive, watch, computed } from 'vue'
import { useToast } from '#imports'

const props = defineProps<{ config: any, theme?: Record<string, string> }>()
const emit = defineEmits<{ (e: 'update:config', v: any): void }>()

const local = reactive({
  commentEnabled: false,
  commentSystem: 'waline',
  commentEmailEnabled: false,
  walineServerURL: ''
})

watch(() => props.config, (v: any) => {
  if (!v) return
  local.commentEnabled = !!v.commentEnabled
  local.commentSystem = String(v.commentSystem || 'waline')
  local.commentEmailEnabled = !!v.commentEmailEnabled
  local.walineServerURL = String(v.walineServerURL || '')
}, { immediate: true, deep: true })

const subtleBg = computed(() => 'bg-gray-800')
const mutedText = computed(() => 'text-slate-400')
const textCls = computed(() => 'text-white')

  const save = async () => {
    try {
      const payload = {
        frontendSettings: {
        commentEnabled: !!props.config?.commentEnabled,
        commentSystem: String(local.commentSystem || 'waline'),
        commentEmailEnabled: !!props.config?.commentEmailEnabled,
        walineServerURL: String(local.walineServerURL || '')
      }
    }
    const response = await fetch('/api/settings', {
      method: 'PUT',
      headers: { 'Content-Type': 'application/json' },
      credentials: 'include',
      body: JSON.stringify(payload)
    })
    const data = await response.json()
    if (response.ok && data.code === 1) {
      const next = { ...props.config, ...payload.frontendSettings }
      emit('update:config', next)
      window.dispatchEvent(new Event('frontend-config-updated'))
      useToast().add({ title: '成功', description: '评论设置已更新', color: 'green' })
    } else {
      throw new Error(data.msg || '保存失败')
    }
  } catch (e: any) {
    useToast().add({ title: '错误', description: e?.message || '保存失败', color: 'red' })
  }
}
</script>
