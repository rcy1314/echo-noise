<template>
  <div id="site-comments-section" class="flex flex-wrap items-center rounded-lg p-3 justify-between gap-3" :class="theme?.subtleBg || subtleBg">
    <div class="flex flex-col gap-3 w-full">
      <div class="flex flex-wrap items-center gap-3 justify-between">
        <USelect v-model="local.commentSystem" :options="[{label:'内置',value:'builtin'},{label:'Waline',value:'waline'}]" />
        <UInput v-if="local.commentSystem === 'waline'" v-model="local.walineServerURL" :ui="{base: theme?.text}" placeholder="Waline 地址" class="w-full md:w-[280px]" />
        <div class="flex-1"></div>
        <UButton color="green" @click="save" class="shadow">保存</UButton>
      </div>
      <div v-if="local.commentSystem === 'builtin' && !props.config?.commentEmailEnabled" class="rounded border p-3" :class="theme?.border">
        <div :class="theme?.text || textCls">开启“邮件通知”后，可配置站点链接地址、主题前缀、发件显示名、文本/HTML 模板，并实时预览。</div>
      </div>
      <div v-if="local.commentSystem === 'builtin' && !!props.config?.commentEmailEnabled" class="grid grid-cols-1 md:grid-cols-3 gap-3">
        <UInput v-model="local.commentEmailSiteURL" :ui="{base: theme?.text}" placeholder="站点链接地址（用于邮件中的 {url} 基础） 如 https://note.noisework.cn" />
        <UInput v-model="local.commentEmailReplyName" :ui="{base: theme?.text}" placeholder="邮件发件显示名（回复通知）" />
        <UInput v-model="local.commentEmailAdminPrefix" :ui="{base: theme?.text}" placeholder="管理员通知主题前缀（可选）" />
        <UInput v-model="local.commentEmailReplyPrefix" :ui="{base: theme?.text}" placeholder="回复通知主题前缀（可选）" />
        <UTextarea v-model="local.commentEmailReplyTemplate" :ui="{base: theme?.text}" placeholder="回复通知正文模板（默认示例）支持 {site} {nick} {content} {url}" />
        <UTextarea v-model="local.commentEmailAdminTemplate" :ui="{base: theme?.text}" placeholder="管理员通知正文模板（可选）支持 {site} {nick} {mail} {link} {content} {url}" />
        <UTextarea v-model="local.commentEmailReplyTemplateHTML" :ui="{base: theme?.text}" placeholder="回复通知 HTML 模板（可选）支持 {site} {nick} {content} {url}" />
        <UTextarea v-model="local.commentEmailAdminTemplateHTML" :ui="{base: theme?.text}" placeholder="管理员通知 HTML 模板（可选）支持 {site} {nick} {mail} {link} {content} {url}" />
        <div class="md:col-span-3 rounded border p-3" :class="theme?.border">
          <div :class="theme?.text">回复通知富文本预览</div>
          <div class="mt-2 rounded p-2 bg-white text-black" v-html="previewReplyHTML"></div>
        </div>
        <div class="md:col-span-3 rounded border p-3" :class="theme?.border">
          <div :class="theme?.text">管理员通知富文本预览</div>
          <div class="mt-2 rounded p-2 bg-white text-black" v-html="previewAdminHTML"></div>
        </div>
      </div>
      
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
  walineServerURL: '',
  commentEmailReplyName: '',
  commentEmailAdminPrefix: '',
  commentEmailReplyPrefix: '',
  commentEmailReplyTemplate: '',
  commentEmailAdminTemplate: ''
  ,commentEmailReplyTemplateHTML: '',
  commentEmailAdminTemplateHTML: ''
})

watch(() => props.config, (v: any) => {
  if (!v) return
  local.commentEnabled = !!v.commentEnabled
  local.commentSystem = String(v.commentSystem || 'waline')
  local.commentEmailEnabled = !!v.commentEmailEnabled
  local.walineServerURL = String(v.walineServerURL || '')
  local.commentEmailReplyName = String(v.commentEmailReplyName || '')
  local.commentEmailAdminPrefix = String(v.commentEmailAdminPrefix || '')
  local.commentEmailReplyPrefix = String(v.commentEmailReplyPrefix || '')
  local.commentEmailReplyTemplate = String(v.commentEmailReplyTemplate || '')
  local.commentEmailAdminTemplate = String(v.commentEmailAdminTemplate || '')
  local.commentEmailSiteURL = String(v.commentEmailSiteURL || '')
  local.commentEmailReplyTemplateHTML = String(v.commentEmailReplyTemplateHTML || '')
  local.commentEmailAdminTemplateHTML = String(v.commentEmailAdminTemplateHTML || '')
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
        walineServerURL: String(local.walineServerURL || ''),
        commentEmailReplyName: String(local.commentEmailReplyName || ''),
        commentEmailAdminPrefix: String(local.commentEmailAdminPrefix || ''),
        commentEmailReplyPrefix: String(local.commentEmailReplyPrefix || ''),
        commentEmailReplyTemplate: String(local.commentEmailReplyTemplate || ''),
        commentEmailAdminTemplate: String(local.commentEmailAdminTemplate || ''),
        commentEmailSiteURL: String(local.commentEmailSiteURL || ''),
        commentEmailReplyTemplateHTML: String(local.commentEmailReplyTemplateHTML || ''),
        commentEmailAdminTemplateHTML: String(local.commentEmailAdminTemplateHTML || '')
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

const previewReplyHTML = computed(() => {
  const site = String(props.config?.siteTitle || '站点')
  const nick = 'Mike'
  const content = '这里是示例内容'
  const base = String(local.commentEmailSiteURL || '')
  const url = (base ? base.replace(/\/$/, '') : '') + '/m/123'
  const tpl = local.commentEmailReplyTemplateHTML || `
<div>
  <p>您在<strong>{site}</strong>主页上的内容有了新的评论</p>
  <p><strong>{nick}</strong> 回复说：</p>
  <p>{content}</p>
  <p>您可以点击查看回复的完整内容：<a href="{url}" target="_blank">{url}</a></p>
</div>`
  return tpl.replaceAll('{site}', site).replaceAll('{nick}', nick).replaceAll('{content}', content).replaceAll('{url}', url)
})
const previewAdminHTML = computed(() => {
  const site = String(props.config?.siteTitle || '站点')
  const nick = 'Mike'
  const mail = 'mike@example.com'
  const link = 'https://example.com'
  const content = '这里是示例内容'
  const base = String(local.commentEmailSiteURL || '')
  const url = (base ? base.replace(/\/$/, '') : '') + '/m/123'
  const tpl = local.commentEmailAdminTemplateHTML || `
<div>
  <p>站点：<strong>{site}</strong></p>
  <p>用户：{nick}（{mail}）</p>
  <p>网址：<a href="{link}" target="_blank">{link}</a></p>
  <p>内容：</p>
  <div>{content}</div>
  <p>查看：<a href="{url}" target="_blank">{url}</a></p>
</div>`
  return tpl
    .replaceAll('{site}', site)
    .replaceAll('{nick}', nick)
    .replaceAll('{mail}', mail)
    .replaceAll('{link}', link)
    .replaceAll('{content}', content)
    .replaceAll('{url}', url)
})
</script>
