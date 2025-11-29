<template>
  <StatusPanel />
  <Notification />
</template>

<script setup lang="ts">
definePageMeta({ middleware: ['auth'] })
import StatusPanel from '@/components/index/StatusPanel.vue'
import Notification from '~/components/widgets/Notification.vue'
</script>
