<template>
    <client-only>
        <StatusPanel @restore-success="handleRestoreSuccess" />
        <Notification />
    </client-only>
</template>

<script setup lang="ts">
import StatusPanel from '@/components/index/StatusPanel.vue'
import Notification from '~/components/widgets/Notification.vue';
import { useToast } from '#imports'
definePageMeta({ middleware: ['auth'] })

// 添加恢复成功处理函数
const handleRestoreSuccess = () => {
    const toast = useToast()
    toast.add({
        title: '恢复成功',
        description: '数据已恢复，页面即将刷新',
        color: 'green',
        timeout: 2000
    })
    
    // 2秒后刷新页面
    setTimeout(() => {
        window.location.reload()
    }, 2000)
}
</script>

<style scoped>
</style>
