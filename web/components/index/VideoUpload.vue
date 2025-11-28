<template>
  <div>
    <input
      ref="videoInput"
      type="file"
      accept="video/*"
      class="hidden"
      @change="handleVideoChange"
    />
    <UButton
      color="gray"
      variant="solid"
      size="sm"
      icon="i-mdi-video"
      class="cursor-pointer"
      @click="triggerVideoInput"
      :title="'上传视频'"
    />
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useToast } from '#imports'

const emit = defineEmits(['video-uploaded', 'upload-progress'])
const videoInput = ref<HTMLInputElement | null>(null)
const toast = useToast()
const BASE_API = useRuntimeConfig().public.baseApi || '/api'

const triggerVideoInput = () => {
  videoInput.value?.click()
}

const handleVideoChange = async (event: Event) => {
  const input = event.target as HTMLInputElement
  const file = input.files ? input.files[0] : null

  if (!file) {
    toast.add({ title: '错误', description: '未选择视频', color: 'red' })
    return
  }

  const maxSize = 200 * 1024 * 1024 // 200MB
  if (file.size > maxSize) {
    toast.add({ title: '错误', description: '视频不能超过200MB', color: 'red' })
    return
  }

  const formData = new FormData()
  formData.append('video', file)

  // 使用 XMLHttpRequest 以支持进度
  const xhr = new XMLHttpRequest()
  xhr.open('POST', `${BASE_API}/videos/upload`, true)
  xhr.withCredentials = true

  xhr.upload.onprogress = (event) => {
    if (event.lengthComputable) {
      const percent = Math.round((event.loaded / event.total) * 100)
      emit('upload-progress', percent)
    }
  }

  xhr.onload = () => {
    if (xhr.status === 200) {
      try {
        const data = JSON.parse(xhr.responseText)
        if (data.code === 1 && data.data) {
          emit('video-uploaded', data.data)
          toast.add({ title: '成功', description: '视频上传成功', color: 'green' })
        } else {
          throw new Error(data.msg || '视频上传失败')
        }
      } catch (error: any) {
        toast.add({ title: '错误', description: error.message || '视频上传失败', color: 'red' })
      }
    } else {
      toast.add({ title: '错误', description: '视频上传失败', color: 'red' })
    }
    emit('upload-progress', 0) // 上传结束后重置
    if (videoInput.value) videoInput.value.value = ''
  }

  xhr.onerror = () => {
    toast.add({ title: '错误', description: '视频上传失败', color: 'red' })
    emit('upload-progress', 0)
    if (videoInput.value) videoInput.value.value = ''
  }

  xhr.send(formData)
}
</script>
