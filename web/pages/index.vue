<template>
  <div class="background-container" :style="backgroundStyle">
    <div class="loading" v-if="!isLoaded">
      <div class="rainbow-spinner"></div>
      <div class="loading-text">加载中...</div>
    </div>
    <div ref="contentWrapper" class="content-wrapper" :class="{ 'gpu-accelerated': true }">
      <UContainer class="container-fixed py-2 pb-4 my-4">
        <div :class="['layout-container', gridModeClass]">
      <ClientOnly>
      <div class="left-col" v-if="!isMobile && layoutState!=='single'">
        <UCard class="sidebar-card" :class="sidebarThemeCard">
          <div class="profile-card">
            <div class="avatar-wrap relative">
              <img class="avatar-lg" :src="profileAvatar" :alt="profileName" @click="changeBackground" @error="handleAvatarError">
              <span class="avatar-status" :class="isOnline ? 'bg-green-500' : 'bg-gray-400'"></span>
            </div>
            <div class="profile-name text-center">{{ profileName }}</div>
            <div v-if="isAdmin" class="mt-1">
              <span class="px-1.5 py-0.5 rounded bg-orange-500 text-white text-[10px]">管理员</span>
            </div>
            <div class="profile-desc">{{ profileDesc }}</div>
            <div v-if="!isOnline" class="auth-actions">
              <UTooltip text="登录">
                <UButton variant="ghost" color="indigo" class="auth-btn" @click="authMode='login'; showAuthModal=true">
                  <UIcon name="i-heroicons-arrow-right-end-on-rectangle" class="w-5 h-5" />
                </UButton>
              </UTooltip>
              <UTooltip text="注册">
                <UButton variant="ghost" color="orange" class="auth-btn" @click="switchToRegister(); showAuthModal=true">
                  <UIcon name="i-heroicons-user-plus" class="w-5 h-5" />
                </UButton>
              </UTooltip>
            </div>
          </div>
        </UCard>
        <UCard class="sidebar-card no-padding-card mt-2" :class="sidebarThemeCard">
          <div class="p-0 grid grid-cols-3 gap-2 text-center text-sm">
            <div>
              <div class="font-semibold">{{ (status?.total_messages ?? status?.TotalMessages ?? 0) }}</div>
              <div class="opacity-70">笔记数</div>
            </div>
            <div>
              <div class="font-semibold">{{ (tags?.length || 0) }}</div>
              <div class="opacity-70">标签</div>
            </div>
            <div>
              <div class="font-semibold">{{ (images?.length || 0) }}</div>
              <div class="opacity-70">图片</div>
            </div>
          </div>
        </UCard>
        <UCard class="sidebar-card no-padding-card mt-2" :class="sidebarThemeCard">
          <div class="social-list">
            <a v-for="item in (frontendConfig.socialLinks || [])" :key="item.url || item.name" class="social-item" :href="item.url" target="_blank" rel="noopener noreferrer" :data-label="item.name || item.url">
              <template v-if="item.imageURL">
                <img :src="item.imageURL" alt="icon" class="social-icon-img" />
              </template>
              <template v-else>
                <UIcon :name="getIconName(item)" class="w-7 h-7" />
              </template>
            </a>
          </div>
        </UCard>
        <UCard v-if="frontendConfig.timeEnabled" class="sidebar-card no-padding-card mt-2" :class="sidebarThemeCard">
          <div class="p-0 text-center clock-card">
            <div class="clock-display">{{ formatTime(currentTime as any) }}</div>
            <div class="clock-date">{{ formatDate(currentTime as any) }}</div>
          </div>
        </UCard>
        

        
        <UCard v-if="frontendConfig.leftAdEnabled && leftAds.length > 0" class="sidebar-card mt-2" :class="sidebarThemeCard">
          <div>
            <template v-if="leftAds.length > 0">
              <div class="relative">
                <a :href="(currentAd.linkURL || '#')" target="_blank" rel="noopener noreferrer" class="block ad-wrap group rounded-lg overflow-hidden" :style="{ '--ad-bg': `url(${imgSrc(currentAd.imageURL)})` }">
                  <img :src="imgSrc(currentAd.imageURL)" alt="ad" class="ad-image w-full object-cover transition duration-200 rounded-lg" />
                  <div class="ad-overlay">
                    <div class="ad-overlay-box transition-colors duration-200" :class="[isDark ? '' : 'group-hover:text-orange-500']">{{ (currentAd.description || '').trim() || '广告' }}</div>
                  </div>
                </a>
                <div v-if="leftAds.length > 1" class="absolute bottom-2 left-0 right-0 flex justify-center gap-2">
                  <button v-for="(ad, i) in leftAds" :key="i" @click="adIndex = i" class="w-2 h-2 rounded-full" :class="i === adIndex ? 'bg-white' : 'bg-white/40'" aria-label="switch-ad"></button>
                </div>
              </div>
            </template>
          </div>
        </UCard>
        <UCard v-if="frontendConfig.hitokotoEnabled" class="sidebar-card mt-2" :class="sidebarThemeCard">
          <div class="hidden"><span id="hitokoto">正在获取中...</span></div>
          <div class="hitokoto-container mx-auto w-full sm:max-w-2xl px-2">
            <div class="hitokoto-text" :title="hitokotoText || '正在获取中...'">
              <MarkdownRenderer :content="hitokotoText || '正在获取中...'" :enableGithubCard="false" />
            </div>
          </div>
        </UCard>
        <div v-if="layoutState==='two'" class="mt-2 space-y-3">
          <UCard v-if="frontendConfig.announcementEnabled && (frontendConfig.announcementText || '').trim() !== ''" class="sidebar-card" :class="sidebarThemeCard">
            <AnnouncementBar :text="frontendConfig.announcementText || '欢迎访问我的说说笔记！'" />
          </UCard>
          <UCard class="sidebar-card no-padding-card" :class="sidebarThemeCard">
            <div>
              <div class="text-xs opacity-70 mb-2">热门标签</div>
              <div class="scroll-tags mb-2">
                <div class="tag-grid">
                  <button v-for="t in popularTags" :key="t.name" class="px-2 py-1 rounded text-xs border opacity-80 hover:opacity-100" @click="handleTagClick(t.name)">#{{ t.name }}<span class="ml-1 opacity-60">{{ t.count }}</span></button>
                </div>
              </div>
            </div>
          </UCard>
          <UCard class="sidebar-card no-padding-card" :class="sidebarThemeCard">
            <div>
              <div class="text-xs opacity-70 mb-2">图集</div>
              <div class="scroll-images">
                <div class="recommend-grid">
                  <a v-for="img in recommendedImages" :key="img.id || img" :href="imageSrc(img)" data-fancybox="recommend-gallery" class="block">
                    <img :src="imageSrc(img)" class="recommend-image-box" loading="lazy" alt="recommend" />
                  </a>
                </div>
              </div>
            </div>
          </UCard>
          <UCard class="sidebar-card no-padding-card" :class="sidebarThemeCard">
            <HeatmapWidget />
          </UCard>
        
        <UCard class="sidebar-card no-padding-card">
          <div class="text-xs opacity-70 mb-1">最新评论</div>
          <div class="scroll-list p-0">
            <div class="flex flex-wrap gap-2">
                <div v-for="(c, i) in visibleRecentComments" :key="c.id || i" class="inline-flex items-center gap-2 px-2 py-1 rounded-full" :class="recentItemClass">
                  <img :src="recentAvatar(c)" :data-mail="c.mail || ''" :data-qq="qqNumberFromEmail(c?.mail || '')" data-try="0" class="w-5 h-5 rounded-full object-cover" alt="avatar" @error="onRecentAvatarError($event, c.nick || '')" />
                  <span class="text-xs" :class="recentTextClass">{{ shortText(c.content) }}</span>
                </div>
              </div>
            </div>
          </UCard>
          
        </div>
      </div>
      </ClientOnly>
      <div class="center-col">
        <div :class="centerContainerClass">
          <div class="moments-header">
            <div class="header-image" :style="headerImageStyle">
              <h1 class="header-title">{{ (frontendConfig.siteTitle || '说说笔记') }}</h1>
              <div class="header-subtitle" ref="subtitleEl"></div>
              <div class="hero-tabs">
                <button v-for="t in centerTabs" :key="t.key" :class="['hero-tab', activeTab===t.key ? 'active' : '']" @click="activeTab=t.key">{{ t.name }}</button>
              </div>
            </div>
          </div>
          <div v-if="activeTab==='links'" class="links-page">
            <UCard class="search-card mb-3" :ui="{ body: 'p-3 md:p-4' }">
              <div class="card-title text-center" :class="isDark ? 'text-white' : 'text-black'">{{ frontendConfig.linksTitle || '友情链接' }}</div>
              <div v-if="(frontendConfig.linksDescription || '').trim() !== ''" class="section-subtitle">{{ frontendConfig.linksDescription }}</div>
              <div class="mx-auto w-full max-w-3xl px-4 sm:px-6">
              <div class="link-grid">
                <a v-for="fl in friendLinksList" :key="fl.link || fl.title" :href="fl.link" target="_blank" rel="noopener noreferrer" class="link-card" :class="isDark ? 'link-card-dark' : 'link-card-light'">
                  <UIcon :name="getIconName({ url: fl.link, icon: fl.icon })" class="w-7 h-7" />
                  <div class="link-content">
                    <div class="link-title">{{ fl.title }}</div>
                    <div v-if="fl.description" class="text-xs opacity-75">{{ fl.description }}</div>
                  </div>
                </a>
              </div>
              </div>
            </UCard>
          </div>
          <div v-else-if="activeTab==='comment'" class="comment-page">
            <UCard class="search-card mb-3" :ui="{ body: 'p-5' }">
              <div class="card-title text-center mb-3" :class="isDark ? 'text-white' : 'text-black'">{{ frontendConfig.commentPageTitle || '留言' }}</div>
              <div v-if="(frontendConfig.commentPageDescription || '').trim() !== ''" class="section-subtitle">{{ frontendConfig.commentPageDescription }}</div>
              <div class="max-w-3xl mx-auto">
                <BuiltinComments v-if="guestbookMessageId" :message-id="guestbookMessageId" :site-config="frontendConfig" :show-input="true" />
                <div v-else class="text-sm opacity-70">正在准备留言板...</div>
              </div>
            </UCard>
          </div>
          <div v-else-if="activeTab==='about'" class="about-page">
            <UCard class="search-card mb-3" :ui="{ body: 'p-6' }">
              <div class="card-title text-center" :class="isDark ? 'text-white' : 'text-black'">{{ frontendConfig.aboutPageTitle || '关于本站' }}</div>
              <div v-if="(frontendConfig.aboutPageDescription || '').trim() !== ''" class="section-subtitle">{{ frontendConfig.aboutPageDescription }}</div>
              <div class="mx-auto w-full max-w-3xl px-4 sm:px-6">
                <MarkdownRenderer :content="(frontendConfig.aboutMarkdown || '').trim() || defaultConfig.aboutMarkdown" />
              </div>
            </UCard>
          </div>
          <template v-else>
            <AddForm @search-result="handleSearchResult" :hide-header-tools="layoutState==='three'" :wide="layoutState==='two'" />
            <div class="mx-auto w-full sm:max-w-4xl mt-3">
              <TagList 
                v-if="tags && tags.length > 0"
                :tags="tags"
                @tagClick="handleTagClick"
                @updateTags="handleTagsUpdate" 
              />
            </div>
            <MessageList 
              ref="messageList" 
              class="message-list-container" 
              :site-config="frontendConfig"
              :target-message-id="targetMessageId" 
              :wide="layoutState==='two'"
            />
          </template>
        </div>
      </div>
      <ClientOnly>
      <div class="right-col space-y-3" v-if="!isMobile && layoutState==='three'">
        <UCard v-if="frontendConfig.announcementEnabled && (frontendConfig.announcementText || '').trim() !== ''" class="sidebar-card" :class="sidebarThemeCard">
          <AnnouncementBar :text="frontendConfig.announcementText || '欢迎访问我的说说笔记！'" />
        </UCard>
        <UCard class="sidebar-card no-padding-card" :class="sidebarThemeCard">
          <div>
            <div class="text-xs opacity-70 mb-2">热门标签</div>
            <div class="scroll-tags mb-2">
              <div class="tag-grid">
                <button v-for="t in popularTags" :key="t.name" class="px-2 py-1 rounded text-xs border opacity-80 hover:opacity-100" @click="handleTagClick(t.name)">#{{ t.name }}<span class="ml-1 opacity-60">{{ t.count }}</span></button>
              </div>
            </div>
          </div>
        </UCard>
        <UCard class="sidebar-card no-padding-card" :class="sidebarThemeCard">
          <div>
            <div class="text-xs opacity-70 mb-2">图集</div>
            <div class="scroll-images">
              <div class="recommend-grid">
                <a v-for="img in recommendedImages" :key="img.id || img" :href="imageSrc(img)" data-fancybox="recommend-gallery" class="block">
                  <img :src="imageSrc(img)" class="recommend-image-box" loading="lazy" alt="recommend" />
                </a>
              </div>
            </div>
          </div>
        </UCard>
        <UCard class="sidebar-card no-padding-card" :class="sidebarThemeCard">
          <HeatmapWidget />
        </UCard>
        
        <UCard class="sidebar-card mt-2 no-padding-card" :class="sidebarThemeCard">
          <div class="text-xs opacity-70 mb-1">最新评论</div>
          <div class="scroll-list p-0">
            <div class="flex flex-wrap gap-2">
              <div v-for="(c, i) in visibleRecentComments" :key="c.id || i" class="inline-flex items-center gap-2 px-2 py-1 rounded-full" :class="recentItemClass">
                <img :src="recentAvatar(c)" :data-mail="c.mail || ''" :data-qq="qqNumberFromEmail(c?.mail || '')" data-try="0" class="w-5 h-5 rounded-full object-cover" alt="avatar" @error="onRecentAvatarError($event, c.nick || '')" />
                <span class="text-xs" :class="recentTextClass">{{ shortText(c.content) }}</span>
              </div>
            </div>
          </div>
        </UCard>
        
        
        
      </div>
      </ClientOnly>
    </div>
    <div v-if="false">
      <UCard v-if="frontendConfig.announcementEnabled && (frontendConfig.announcementText || '').trim() !== ''" class="mx-auto sm:max-w-2xl mb-3 sidebar-card">
        <AnnouncementBar :text="frontendConfig.announcementText || '欢迎访问我的说说笔记！'" />
      </UCard>
      <UCard class="mx-auto sm:max-w-2xl mb-3 sidebar-card">
        <div>
          <div class="text-xs opacity-70 mb-2">热门话题</div>
          <div class="flex flex-wrap gap-2 mb-3">
            <button v-for="t in popularTags" :key="t.name" class="px-2 py-1 rounded text-xs border opacity-80 hover:opacity-100" @click="handleTagClick(t.name)">#{{ t.name }}<span class="ml-1 opacity-60">{{ t.count }}</span></button>
          </div>
          <div class="text-xs opacity-70 mb-2">推荐图集</div>
          <div class="scroll-images">
            <div class="recommend-grid">
              <a v-for="img in recommendedImages" :key="img.id || img" :href="imageSrc(img)" data-fancybox="recommend-gallery" class="block">
                <img :src="imageSrc(img)" class="recommend-image-box" loading="lazy" alt="recommend" />
              </a>
            </div>
          </div>
        </div>
      </UCard>
      <div class="mx-auto sm:max-w-2xl mb-4">
        <HeatmapWidget />
      </div>
      <UCard v-if="frontendConfig.leftAdEnabled && leftAds.length > 0" class="mx-auto sm:max-w-2xl mb-3 sidebar-card">
        <div class="sidebar-title flex items-center gap-2">
          <UIcon name="i-heroicons-megaphone" class="w-4 h-4" />
          <span>广而告之</span>
        </div>
        <div class="mt-2">
          <template v-if="leftAds.length > 0">
            <div class="relative">
              <a :href="(currentAd.linkURL || '#')" target="_blank" rel="noopener noreferrer" class="block ad-wrap" :style="{ '--ad-bg': `url(${imgSrc(currentAd.imageURL)})` }">
                <img :src="imgSrc(currentAd.imageURL)" alt="ad" class="ad-image w-full rounded-md object-cover transition duration-200" />
                <div class="ad-overlay">
                  <div class="ad-overlay-box">{{ (currentAd.description || '').trim() || '广告' }}</div>
                </div>
              </a>
              <div v-if="leftAds.length > 1" class="absolute bottom-2 left-0 right-0 flex justify-center gap-2">
                <button v-for="(ad, i) in leftAds" :key="i" @click="adIndex = i" class="w-2 h-2 rounded-full" :class="i === adIndex ? 'bg-white' : 'bg-white/40'" aria-label="switch-ad"></button>
              </div>
            </div>
            
          </template>
        </div>
      </UCard>
      <TagList 
        v-if="tags && tags.length > 0"
        :tags="tags"
        @tagClick="handleTagClick"
        @updateTags="handleTagsUpdate" 
      />
    </div>
      <!-- 音乐播放器容器（浮动或嵌入） -->
      <div
        v-if="frontendConfig.musicEnabled"
        class="netease-mini-player"
        :data-playlist-id="frontendConfig.musicPlaylistId || ''"
        :data-song-id="frontendConfig.musicSongId || ''"
        :data-position="frontendConfig.musicPosition || 'bottom-left'"
        :data-theme="frontendConfig.musicTheme || 'auto'"
        :data-lyric="frontendConfig.musicLyric ? 'true' : 'false'"
        :data-default-minimized="frontendConfig.musicDefaultMinimized ? 'true' : 'false'"
        :data-embed="frontendConfig.musicEmbed ? 'true' : 'false'"
        :data-autoplay="frontendConfig.musicAutoplay ? 'true' : 'false'"
      />
      </UContainer>
  <Notification />
  <!-- 添加搜索模态框组件 -->
  <SearchMode v-model="showSearchModal" @search-result="handleSearchResult" />
  <FloatingToolSidebar 
    :content-theme="contentTheme"
    :layout-icon="layoutIcon"
    @search="showSearchModal = true"
    @switch-background="changeBackground"
    @toggle-theme="toggleThemeGlobal"
    @toggle-layout="cycleLayout"
    @open-admin="openAdmin"
  />
  <UModal v-model="showAuthModal" :ui="{ width: 'sm:max-w-md', container: 'items-center', base: 'backdrop-blur-sm' }">
    <UCard class="search-card">
      <template #header>
        <div class="flex justify-between items-center">
          <h3 class="text-xl font-semibold">{{ authMode === 'login' ? '欢迎回来' : '创建账户' }}</h3>
          <div class="flex items-center gap-2">
            <UButton v-if="authMode==='login'" variant="link" color="orange" @click="switchToRegister">去注册</UButton>
            <UButton v-else variant="link" color="orange" @click="switchToLogin">去登录</UButton>
            <UButton icon="i-mdi-close" variant="ghost" color="gray" @click="closeAuthModal" />
          </div>
        </div>
      </template>
      <div class="space-y-3">
        <div v-if="authMode==='login'">
          <UForm @submit.prevent="onLoginSubmit">
            <div class="space-y-3">
              <UInput v-model="loginForm.username" placeholder="用户名或邮箱" />
              <UInput v-model="loginForm.password" :type="showLoginPassword ? 'text' : 'password'" placeholder="密码" autocomplete="current-password" autocorrect="off" autocapitalize="off" spellcheck="false">
                <template #trailing>
                  <UButton icon="i-heroicons-eye" v-if="!showLoginPassword" variant="ghost" color="gray" @click="showLoginPassword = true" :ui="{ rounded: 'rounded-full' }" />
                  <UButton icon="i-heroicons-eye-slash" v-else variant="ghost" color="gray" @click="showLoginPassword = false" :ui="{ rounded: 'rounded-full' }" />
                </template>
              </UInput>
              <UButton class="w-full" :loading="loginSubmitting" :disabled="loginSubmitting" type="submit" color="primary">登录</UButton>
            </div>
          </UForm>
          <div class="flex items-center justify-between text-sm mt-2">
            <UButton variant="link" @click="switchToRegister">账号注册</UButton>
            <UButton variant="link" @click="showForgot = true">忘记密码？</UButton>
          </div>
          <div class="text-center text-sm opacity-80 mt-3" v-if="githubEnabled">其他登录方式</div>
          <div class="mt-2" v-if="githubEnabled">
            <UButton class="w-full h-10 px-3 gap-2 justify-center font-medium bg-[#24292f] hover:bg-[#1f2328] text-white ring-1 ring-black/20" @click="loginWithGithub">
              <UIcon name="i-mdi-github" class="w-5 h-5" />
              <span>GitHub 一键登录</span>
            </UButton>
          </div>
          
        </div>
        <div v-else>
          <UForm @submit.prevent="onRegisterSubmit">
            <div class="space-y-3">
              <UInput v-model="registerForm.username" placeholder="用户名" />
              <UInput v-model="registerForm.password" :type="showRegisterPassword ? 'text' : 'password'" placeholder="密码" autocomplete="new-password" autocorrect="off" autocapitalize="off" spellcheck="false">
                <template #trailing>
                  <UButton icon="i-heroicons-eye" v-if="!showRegisterPassword" variant="ghost" color="gray" @click="showRegisterPassword = true" :ui="{ rounded: 'rounded-full' }" />
                  <UButton icon="i-heroicons-eye-slash" v-else variant="ghost" color="gray" @click="showRegisterPassword = false" :ui="{ rounded: 'rounded-full' }" />
                </template>
              </UInput>
              <div class="flex items-center gap-2">
                <UInput v-model="registerForm.captcha" placeholder="验证码" />
                <img :src="captchaSrc" @click="refreshCaptcha" class="h-10 w-24 rounded border cursor-pointer" alt="captcha" />
                <UBadge :color="remaining>0 ? 'primary' : 'red'" variant="soft">{{ remaining>0 ? `有效 ${remaining}s` : '已过期' }}</UBadge>
              </div>
              <UButton class="w-full" :loading="registerSubmitting" :disabled="remaining<=0 || registerSubmitting" type="submit" color="primary">{{ remaining>0 ? '注册' : '验证码已过期' }}</UButton>
            </div>
          </UForm>
          
        </div>
      </div>
    </UCard>
  </UModal>
  <UModal v-model="showForgot" :ui="{ container: 'items-center', base: 'backdrop-blur-sm' }">
    <UCard class="search-card">
      <div class="font-semibold mb-2">找回密码</div>
      <UForm @submit.prevent="onForgot">
        <div class="space-y-3">
          <UInput v-model="forgot.account" placeholder="用户名或邮箱" />
          <div class="flex justify-end gap-2">
            <UButton variant="ghost" @click="showForgot = false">取消</UButton>
            <UButton :disabled="forgotCooldown>0 || !smtpEnabled" type="submit" color="primary">{{ smtpEnabled ? (forgotCooldown>0 ? `请稍候(${forgotCooldown}s)` : '发送重置邮件') : '邮件未开启' }}</UButton>
          </div>
        </div>
      </UForm>
    </UCard>
  </UModal>
  <div class="scroll-buttons" @mouseenter="hoverScroll = true" @mouseleave="hoverScroll = false">
    <UButton v-show="showScroll" :class="scrollButtonClass" variant="ghost" size="sm" @click="handleScrollClick">
      <UIcon :class="iconClass" :name="scrollIconName" />
    </UButton>
  </div>
  </div>
</div>
</template>

<script setup lang="ts">
import { ref, computed, inject, provide, onMounted, onUnmounted, watch, nextTick } from 'vue'
import { useRoute } from '#imports'
import AddForm from '@/components/index/AddForm.vue'
import MessageList from '@/components/index/MessageList.vue'
import Notification from '~/components/widgets/Notification.vue';
import HeatmapWidget from '~/components/widgets/heatmap.vue'
import SearchMode from '~/components/index/Searchmode.vue' // 导入 SearchMode 组件
import TagList from '~/components/index/TagList.vue'
import AnnouncementBar from '~/components/widgets/AnnouncementBar.vue'
import FloatingToolSidebar from '~/components/widgets/FloatingToolSidebar.vue'
import BuiltinComments from '~/components/comments/BuiltinComments.vue'
import MarkdownRenderer from '~/components/index/MarkdownRenderer.vue'
import { getRequest } from '~/utils/api'
import { useToast } from '#ui/composables/useToast'
import { useUserStore } from '~/store/user'
const router = useRouter()
const route = useRoute()
const baseApi = useRuntimeConfig().public.baseApi || '/api'
const initialLayout = ((): 'three' | 'two' | 'single' => {
  if (typeof window === 'undefined') return 'three'
  const isMobileInit = window.matchMedia('(max-width: 1024px)').matches
  const saved = localStorage.getItem(isMobileInit ? 'homeLayoutMobile' : 'homeLayoutDesktop') as any
  if (saved) return saved as any
  return isMobileInit ? 'single' : 'three'
})()
const layoutState = ref<'three' | 'two' | 'single'>(initialLayout)
const mq = typeof window !== 'undefined' ? window.matchMedia('(max-width: 1024px)') : null
const isMobile = ref<boolean>(!!mq?.matches)
const cycleLayout = () => {
  if (isMobile.value) return
  layoutState.value = layoutState.value === 'three' ? 'two' : (layoutState.value === 'two' ? 'single' : 'three')
  if (typeof window !== 'undefined') localStorage.setItem('homeLayoutDesktop', layoutState.value)
}
const gridModeClass = computed(() => (layoutState.value === 'three' ? 'grid-3' : (layoutState.value === 'two' ? 'grid-2' : 'grid-1')))
const layoutIcon = computed(() => (layoutState.value === 'three' ? 'i-mdi-view-grid' : (layoutState.value === 'two' ? 'i-mdi-view-column' : 'i-mdi-view-stream')))
const centerContainerClass = computed(() => (
  layoutState.value === 'two'
    ? 'w-full max-w-none'
    : (layoutState.value === 'single'
        ? 'mx-auto w-full max-w-[640px] sm:max-w-3xl'
        : 'mx-auto w-full sm:max-w-4xl')
))
const toggleHeatmapCard = () => { showHeatmap.value = !showHeatmap.value }
// 主题预设。统一由 ThemePresetSwitcher 控制 documentElement 类，不在容器上附加主题类
  const centerTabs = [
    { key: 'latest', name: '最新' },
    { key: 'links', name: '友链' },
    { key: 'comment', name: '留言' },
    { key: 'about', name: '关于' }
  ]
  const activeTab = ref('latest')


// 添加 messageList ref
const messageList = ref(null)
// 搜索模态的开关
const showSearchModal = ref(false)
const showAuthModal = ref(false)
const authMode = ref<'login'|'register'>('login')
const loginForm = reactive({ username: '', password: '' })
const registerForm = reactive({ username: '', password: '', captcha: '' })
const showLoginPassword = ref(false)
const showRegisterPassword = ref(false)
const loginSubmitting = ref(false)
const registerSubmitting = ref(false)
const captchaSrc = ref('')
const remaining = ref(0)
let captchaExpiresAt: number | null = null
let captchaTimer: any = null
const showForgot = ref(false)
const forgot = reactive({ account: '' })
const forgotCooldown = ref(0)
let forgotTimer: any = null
const smtpEnabled = ref(true)
const githubEnabled = ref(false)
const refreshCaptcha = async () => {
  try {
    captchaSrc.value = `${baseApi}/captcha?ts=${Date.now()}`
    captchaExpiresAt = Date.now() + 120000
    remaining.value = 120
    if (captchaTimer) clearInterval(captchaTimer)
    captchaTimer = setInterval(() => {
      const r = Math.max(0, Math.ceil(((captchaExpiresAt || Date.now()) - Date.now()) / 1000))
      remaining.value = r
      if (r <= 0) clearInterval(captchaTimer)
    }, 1000)
  } catch { remaining.value = 0 }
}
const switchToRegister = async () => {
  try {
    const res = await fetch(`${baseApi}/frontend/config`, { credentials: 'include' })
    const data = await res.json()
    const allowed = !!data?.data?.allowRegistration
    if (!allowed) {
      useToast().add({ title: '提示', description: '站点已关闭用户注册', color: 'orange' })
      return
    }
  } catch {}
  authMode.value = 'register'
  refreshCaptcha()
}
const switchToLogin = () => { authMode.value = 'login' }
const closeAuthModal = () => {
  showAuthModal.value = false
  // 移除 ?login=1，避免关闭后再次自动弹出
  const { login, ...rest } = route.query as any
  router.replace({ path: route.path, query: rest })
}
const onLoginSubmit = async () => {
  loginSubmitting.value = true
  try {
    const ok = await useUserStore().login({ username: loginForm.username, password: loginForm.password })
    if (ok) {
      useToast().add({ title: '登录成功', color: 'green' })
      showAuthModal.value = false
      const redirect = (route.query.redirect as string) || '/status'
      router.push(redirect)
    }
  } catch {
    useToast().add({ title: '登录失败', color: 'red' })
  } finally { loginSubmitting.value = false }
}
const onRegisterSubmit = async () => {
  registerSubmitting.value = true
  try {
    if (!registerForm.username || !registerForm.password || !registerForm.captcha) {
      throw new Error('请完整填写用户名、密码与验证码')
    }
    if ((captchaExpiresAt || 0) < Date.now()) { throw new Error('验证码已过期，请刷新后再提交') }
    const res = await fetch(`${baseApi}/register`, { method: 'POST', headers: { 'Content-Type': 'application/json' }, credentials: 'include', body: JSON.stringify(registerForm) })
    const data = await res.json().catch(() => ({}))
    if (!res.ok || data.code !== 1) throw new Error(data?.msg || '注册失败')
    useToast().add({ title: '注册成功', color: 'green' })
    authMode.value = 'login'
  } catch (e: any) {
    useToast().add({ title: '注册失败', description: e.message || '请稍后重试', color: 'red' })
    refreshCaptcha()
  } finally { registerSubmitting.value = false }
}
const openAdmin = async () => {
  const ok = await useUserStore().checkLoginStatus()
  if (ok) router.push('/status')
  else {
    authMode.value = 'login';
    showAuthModal.value = true;
    try {
      const res = await fetch(`${baseApi}/frontend/config`, { credentials: 'include' })
      const data = await res.json()
      githubEnabled.value = !!data?.data?.frontendSettings?.githubOAuthEnabled
      smtpEnabled.value = !!data?.data?.smtpEnabled
    } catch {}
  }
}
onUnmounted(() => { if (captchaTimer) clearInterval(captchaTimer) })

watch(() => route.query.login, (v) => {
  if (v) { authMode.value = 'login'; showAuthModal.value = true }
}, { immediate: true })
watch(showAuthModal, (v) => {
  if (!v && route.query.login) {
    const { login, ...rest } = route.query as any
    router.replace({ path: route.path, query: rest })
  }
})
const loginWithGithub = () => { window.location.href = `${baseApi}/oauth/github/login` }
const onForgot = async () => {
  try {
    const res = await fetch(`${baseApi}/password/forgot`, { method: 'POST', headers: { 'Content-Type': 'application/json' }, credentials: 'include', body: JSON.stringify({ account: forgot.account }) })
    const data = await res.json().catch(() => ({}))
    if (!res.ok || data.code !== 1) throw new Error(data?.msg || '发送失败')
    useToast().add({ title: data?.msg || '已发送', description: '请查收重置邮件', color: 'green' })
    forgotCooldown.value = 60
    if (forgotTimer) clearInterval(forgotTimer)
    forgotTimer = setInterval(() => { if (forgotCooldown.value > 0) forgotCooldown.value--; else clearInterval(forgotTimer) }, 1000)
  } catch (e: any) {
    useToast().add({ title: '失败', description: e.message || '发送失败', color: 'red' })
  }
}
// 修复：定义 targetMessageId，避免模板引用未定义导致列表不渲染
const targetMessageId = ref<string | null>(null)
// 添加搜索结果处理函数
const handleSearchResult = (result: any) => {
  console.log('接收到搜索结果:', result); // 添加调试日志
  if (messageList.value) {
    // 直接传递原始结果，让 MessageList 组件自己处理数据格式
    messageList.value.handleSearchResult(result);
  }
}
// 留言板目标消息ID（默认取最新公开消息）
const guestbookMessageId = ref<number | null>(null)
const loadGuestbookTarget = async () => {
  try {
    const res = await fetch((useRuntimeConfig().public.baseApi || '/api') + '/guestbook/message', {
      method: 'GET',
      headers: { 'Accept': 'application/json' },
      credentials: 'include'
    })
    const js = await res.json().catch(() => null)
    const id = js?.data?.id
    guestbookMessageId.value = typeof id === 'number' ? id : null
  } catch { guestbookMessageId.value = null }
}
watch(() => activeTab.value, (k) => { if (k === 'comment' && !guestbookMessageId.value) loadGuestbookTarget() })
onMounted(() => { loadGuestbookTarget() })


const userStore = useUserStore()
const isOnline = computed(() => !!userStore.isLogin)
const defaultAdmin = computed(() => {
  const list = ((status.value as any)?.users || ( (status.value as any)?.Users ) || []) as any[]
  const admin = Array.isArray(list) ? list.find((it: any) => !!(it?.is_admin ?? it?.IsAdmin)) : null
  return admin || (Array.isArray(list) ? list[0] : null) || null
})
const isAdmin = computed(() => {
  const u = userStore.user as any
  if (u) return !!(u?.is_admin || u?.IsAdmin)
  const a = defaultAdmin.value as any
  return !!(a && (a?.is_admin || a?.IsAdmin))
})
const profileName = computed(() => {
  const u = userStore.user as any
  const name = String(u?.username || u?.Username || '').trim()
  if (name) return name
  const wname = String((frontendConfig.value as any)?.welcomeName || '').trim()
  return wname || '匿名'
})
const fallbackAvatarURL = 'https://s2.loli.net/2025/03/24/HnSXKvibAQlosIW.png'
const profileAvatar = computed(() => {
  const u = userStore.user as any
  const raw = String(u?.avatar_url || u?.AvatarURL || '').trim()
  const base = useRuntimeConfig().public.baseApi || '/api'
  if (raw) return raw.startsWith('http') ? raw : `${base}${raw}`
  const wraw = String((frontendConfig.value as any)?.welcomeAvatarURL || '').trim()
  if (wraw) return wraw.startsWith('http') ? wraw : `${base}${wraw}`
  return fallbackAvatarURL
})
const handleAvatarError = (e: Event) => {
  const img = e.target as HTMLImageElement
  if (img && img.src !== fallbackAvatarURL) img.src = fallbackAvatarURL
}
const profileDesc = computed(() => {
  const u = userStore.user as any
  const d = String(u?.description || '').trim()
  if (d) return d
  const wd = String((frontendConfig.value as any)?.welcomeDescription || '').trim()
  return wd || '执迷不悟'
})
// 注入从AddForm组件提供的showHeatmap变量
const showHeatmap = ref(true);
const contentTheme = ref<string>(typeof window !== 'undefined' ? (localStorage.getItem('contentTheme') || 'light') : 'light')
const rightThemeIcon = computed(() => (contentTheme.value === 'dark' ? 'i-mdi-weather-night' : 'i-mdi-white-balance-sunny'))
const enableAutoScroll = ref(false)
// 提供给子组件
provide('showHeatmap', showHeatmap);
provide('contentTheme', contentTheme)
  provide('toggleContentTheme', () => {
    contentTheme.value = contentTheme.value === 'dark' ? 'light' : 'dark'
    if (typeof window !== 'undefined') {
      localStorage.setItem('contentTheme', contentTheme.value)
      document.documentElement.classList.toggle('dark', contentTheme.value === 'dark')
    }
  })

// 同步文档根节点主题类，修复浅色/深色不一致与初始化错误
watch(() => contentTheme.value, (val) => {
  if (typeof window !== 'undefined') {
    document.documentElement.classList.toggle('dark', val === 'dark')
  }
}, { immediate: true })

const toggleThemeGlobal = () => {
  contentTheme.value = contentTheme.value === 'dark' ? 'light' : 'dark'
  if (typeof window !== 'undefined') {
    localStorage.setItem('contentTheme', contentTheme.value)
    document.documentElement.classList.toggle('dark', contentTheme.value === 'dark')
  }
}

const contentWrapper = ref<HTMLElement | null>(null)
const scrollToTop = () => {
  const el = contentWrapper.value
  if (el) el.scrollTo({ top: 0, behavior: 'smooth' })
  else window.scrollTo({ top: 0, behavior: 'smooth' })
}
  const scrollToBottom = () => {
    const el = contentWrapper.value
    if (el) el.scrollTo({ top: el.scrollHeight, behavior: 'smooth' })
    else {
      const h = Math.max(document.body.scrollHeight, document.documentElement.scrollHeight)
      window.scrollTo({ top: h, behavior: 'smooth' })
    }
  }

  const autoScrollCleanups: Array<() => void> = []
  const startAutoScroll = () => {
    const lists = Array.from(document.querySelectorAll('.scroll-list, .scroll-images')) as HTMLElement[]
    lists.forEach((el) => {
      let id = 0
      let last = performance.now()
      let pauseUntil = 0
      const speed = el.classList.contains('scroll-list') ? 0.25 : 0.35
      const pauseMs = 1200
      const step = (ts: number) => {
        id = requestAnimationFrame(step)
        const dt = ts - last
        last = ts
        if (pauseUntil && ts < pauseUntil) return
        el.scrollTop += (speed * dt) / 16
        if (el.scrollTop >= el.scrollHeight - el.clientHeight - 1) {
          el.scrollTop = 0
          pauseUntil = ts + pauseMs
        }
      }
      id = requestAnimationFrame(step)
      autoScrollCleanups.push(() => cancelAnimationFrame(id))
    })
  }

// 滚动状态与按钮展示逻辑
const hoverScroll = ref(false)
const isAtTop = ref(true)
const isAtBottom = ref(false)
const updateScrollState = () => {
  const el = contentWrapper.value
  if (!el) {
    const y = window.scrollY || document.documentElement.scrollTop || 0
    const max = Math.max(document.body.scrollHeight, document.documentElement.scrollHeight)
    isAtTop.value = y <= 2
    isAtBottom.value = window.innerHeight + y >= max - 2
    return
  }
  const y = el.scrollTop
  const max = el.scrollHeight
  isAtTop.value = y <= 2
  isAtBottom.value = el.clientHeight + y >= max - 2
}
onMounted(() => {
  updateScrollState()
  contentWrapper.value?.addEventListener('scroll', updateScrollState, { passive: true })
  if (enableAutoScroll.value) nextTick(() => startAutoScroll())
})
onUnmounted(() => {
  contentWrapper.value?.removeEventListener('scroll', updateScrollState)
  autoScrollCleanups.forEach((fn) => fn())
})

const showScroll = computed(() => isAtTop.value || isAtBottom.value || hoverScroll.value)
const scrollIconName = computed(() => (isAtBottom.value && !isAtTop.value) ? 'i-heroicons-arrow-up' : 'i-heroicons-arrow-down')
const handleScrollClick = () => {
  if (isAtBottom.value && !isAtTop.value) {
    scrollToTop()
  } else {
    scrollToBottom()
  }
}
const isDark = computed(() => contentTheme.value === 'dark')
const sidebarThemeCard = computed(() => (
  isDark.value
    ? 'bg-[rgba(36,43,50,0.95)] text-white border border-white/20'
    : 'bg-white text-black border border-black/10'
))
const scrollButtonClass = computed(() => (
  isDark.value
    ? 'scroll-button bg-[rgba(36,43,50,0.85)] hover:bg-[rgba(36,43,50,0.95)] text-white shadow-[0_6px_16px_rgba(0,0,0,0.35)]'
    : 'scroll-button bg-white/95 hover:bg-white text-gray-700 ring-1 ring-gray-200 shadow-[0_4px_12px_rgba(0,0,0,0.12)]'
))
const iconClass = computed(() => (isDark.value ? 'text-white w-6 h-6' : 'text-gray-600 w-6 h-6'))

// 添加监听，查看状态变化
watch(showHeatmap, (newVal) => {
  console.log('index.vue 中热力图状态变化:', newVal);
});
// 添加 useHead
useHead({
  meta: [
    {
      name: 'viewport',
      content: 'width=device-width, initial-scale=1.0, maximum-scale=1.0, user-scalable=no'
    }
  ]
})
// 同步路由中的消息ID到 MessageList，用于高亮或定位
watch(() => route.hash, (newHash) => {
  const id = (newHash || '').split('/messages/').pop()
  targetMessageId.value = id || null
}, { immediate: true })

// 添加前端配置的响应式对象
const frontendConfig = ref({
    siteTitle: '',
    subtitleText: '',
    avatarURL: '',
    username: '',
    description: '',
    backgrounds: [] as string[],
    cardFooterTitle: '',
    cardFooterSubtitle: '',
    pageFooterHTML: '',
    rssTitle: '',
    rssDescription: '',
  rssAuthorName: '',
  rssFaviconURL: '',
  linksTitle: '',
  linksDescription: '',
  commentPageTitle: '',
  commentPageDescription: '',
  aboutPageTitle: '',
  aboutPageDescription: '',
  aboutMarkdown: '',
  friendLinks: [] as Array<{ title?: string; link: string; icon?: string; description?: string }>,
    enableGithubCard: false,
    // PWA
    pwaEnabled: true,
    pwaTitle: '',
    pwaDescription: '',
    pwaIconURL: '',
    announcementText: '',
    announcementEnabled: true,
    hitokotoEnabled: true,
    // 评论系统
    commentEnabled: true,
    commentSystem: 'builtin',
    commentEmailEnabled: false,
    commentLoginRequired: false,
    // 音乐配置
    musicEnabled: false,
    musicPlaylistId: '2141128031',
    musicSongId: '',
    musicPosition: 'bottom-left',
    musicTheme: 'auto',
    musicLyric: true,
    musicAutoplay: false,
    musicDefaultMinimized: true,
    musicEmbed: false,
    musicCssCdnURL: '',
    musicJsCdnURL: '',
    calendarEnabled: true,
  timeEnabled: true,
  // 左栏广告
  leftAdEnabled: true,
    leftAdImageURL: `https://picsum.photos/seed/${Math.random().toString(36).slice(2)}/640/640`,
    leftAdLinkURL: 'https://note.noisework.cn',
    leftAdDescription: '示例广告（单条配置）',
    leftAds: [
      { imageURL: `https://picsum.photos/seed/${Math.random().toString(36).slice(2)}/640/640`, linkURL: 'https://note.noisework.cn', description: '写作与记录，开启灵感之旅' },
      { imageURL: `https://picsum.photos/seed/${Math.random().toString(36).slice(2)}/640/640`, linkURL: 'https://noisework.cn', description: '探索新主题与小工具' },
      { imageURL: `https://picsum.photos/seed/${Math.random().toString(36).slice(2)}/640/640`, linkURL: 'https://github.com', description: '开源项目，欢迎 Star' },
    ],
  leftAdsIntervalMs: 4000
})

const initNMP = async () => {
  await nextTick()
  const nmp = (window as any).NeteaseMiniPlayer
  if (nmp && typeof nmp.init === 'function') {
    nmp.init()
  }
}

const loadNMPAssets = async (): Promise<boolean> => {
  if (typeof window === 'undefined') return false
  // 如果已经存在并且全局对象可用，直接返回
  if ((window as any).NeteaseMiniPlayer) return true
  const head = document.head
  const body = document.body
  // 样式按需注入（本地资源，避免外链失败）
  const cssId = 'nmp-css'
  if (!document.getElementById(cssId)) {
    const link = document.createElement('link')
    link.id = cssId
    link.rel = 'stylesheet'
    const cssCdn = (frontendConfig.value as any).musicCssCdnURL || ''
    const localCss = '/NeteaseMiniPlayer/netease-mini-player-v2.css'
    link.href = (cssCdn && cssCdn.trim() !== '') ? cssCdn.trim() : localCss
    link.onerror = () => { link.href = localCss }
    head.appendChild(link)
  }
  // 脚本按需注入
  const scriptId = 'nmp-js'
  if (!document.getElementById(scriptId)) {
    return await new Promise<boolean>((resolve) => {
      const script = document.createElement('script')
      script.id = scriptId
      const jsCdn = (frontendConfig.value as any).musicJsCdnURL || ''
      const localJs = '/NeteaseMiniPlayer/netease-mini-player-v2.js'
      script.src = (jsCdn && jsCdn.trim() !== '') ? jsCdn.trim() : localJs
      script.defer = true
      script.onload = () => resolve(!!(window as any).NeteaseMiniPlayer)
      script.onerror = () => {
        if (script.src !== localJs) {
          script.src = localJs
          script.onload = () => resolve(!!(window as any).NeteaseMiniPlayer)
          script.onerror = () => resolve(false)
          return
        }
        resolve(false)
      }
      body.appendChild(script)
    })
  }
  return !!(window as any).NeteaseMiniPlayer
}

const fetchWithTimeout = async (url: string, ms = 2500): Promise<Response> => {
  const ctrl = new AbortController()
  const timer = setTimeout(() => ctrl.abort(), ms)
  try {
    const res = await fetch(url, { cache: 'no-cache', signal: ctrl.signal })
    clearTimeout(timer)
    return res
  } catch (e) {
    clearTimeout(timer)
    throw e
  }
}

const loadHitokoto = async () => {
  if (typeof window === 'undefined') return
  await nextTick()
  try {
    const r1 = await fetchWithTimeout('https://v1.hitokoto.cn/?encode=json', 2500)
    if (r1.ok) {
      const j = await r1.json()
      const t = String(j?.hitokoto || '').trim()
      if (t) { hitokotoText.value = t; return }
    }
  } catch {}
  try {
    const r2 = await fetchWithTimeout('https://v1.hitokoto.cn/?encode=text', 2500)
    if (r2.ok) {
      const t = String(await r2.text()).trim()
      if (t) { hitokotoText.value = t; return }
    }
  } catch {}
  try {
    if (!document.getElementById('hitokoto-script-left')) {
      const s = document.createElement('script')
      s.id = 'hitokoto-script-left'
      s.src = 'https://v1.hitokoto.cn/?encode=js&select=%23hitokoto'
      s.defer = true
      s.onload = () => {
        try {
          const el = document.getElementById('hitokoto')
          const txt = (el?.textContent || el?.innerText || '').trim()
          if (txt) hitokotoText.value = txt
        } catch {}
      }
      document.body.appendChild(s)
    }
  } catch {}
  const txt = String(hitokotoText.value || '').trim()
  if (!txt) hitokotoText.value = '身为冒险者，如果安静的老死在床上，那简直就是耻辱！'
}

watch(() => frontendConfig.value.musicEnabled, async (enabled) => {
  if (enabled) {
    const ok = await loadNMPAssets()
    if (ok) await initNMP()
  }
})

watch(() => frontendConfig.value.hitokotoEnabled, async (enabled) => {
  if (enabled) await loadHitokoto()
}, { immediate: true })

const backgroundStyle = computed(() => ({
    '--bg-image': `url(${currentImage.value || frontendConfig.value.backgrounds[0]})`
}))
// 添加 headerImageStyle 计算属性
const headerImageStyle = computed(() => ({
    'background-image': `url(${currentImage.value || frontendConfig.value.backgrounds[0]})`,
    'background-size': 'cover',
    'background-position': 'center'
}))
// 修改 fetchConfig 方法

// 首先添加默认配置对象
  const defaultConfig = {
    siteTitle: '说说笔记',
    subtitleText: '欢迎访问，点击头像可更换封面背景！',
    avatarURL: 'https://s2.loli.net/2025/03/24/HnSXKvibAQlosIW.png',
    username: 'Noise',
    description: '执迷不悟',
    // 系统欢迎组件（未登录时显示）
    welcomeAvatarURL: 'https://s2.loli.net/2025/03/24/HnSXKvibAQlosIW.png',
    welcomeName: 'Noise',
    welcomeDescription: '执迷不悟',
    welcomeUseAdmin: true,
    backgrounds: [
                'https://s2.loli.net/2025/03/27/KJ1trnU2ksbFEYM.jpg',
                'https://s2.loli.net/2025/03/27/MZqaLczCvwjSmW7.jpg',
                'https://s2.loli.net/2025/03/27/UMijKXwJ9yTqSeE.jpg',
                'https://s2.loli.net/2025/03/27/WJQIlkXvBg2afcR.jpg',
                'https://s2.loli.net/2025/03/27/oHNQtf4spkq2iln.jpg',
                'https://s2.loli.net/2025/03/27/PMRuX5loc6Uaimw.jpg',
                'https://s2.loli.net/2025/03/27/U2WIslbNyTLt4rD.jpg',
                'https://s2.loli.net/2025/03/27/xu1jZL5Og4pqT9d.jpg',
                'https://s2.loli.net/2025/03/27/OXqwzZ6v3PVIns9.jpg',
                'https://s2.loli.net/2025/03/27/HGuqlE6apgNywbh.jpg',
                'https://s2.loli.net/2025/03/26/d7iyuPYA8cRqD1K.jpg',
                'https://s2.loli.net/2025/03/27/wYy12qDMH6bGJOI.jpg',
                'https://s2.loli.net/2025/03/27/y67m2k5xcSdTsHN.jpg',
    ],
    cardFooterTitle: 'Noise·说说·笔记~',
    cardFooterSubtitle: 'note.noisework.cn',
    pageFooterHTML: '<div class="text-center text-xs text-gray-400 py-4">来自<a href="https://www.noisework.cn" target="_blank" rel="noopener noreferrer" class="text-orange-400 hover:text-orange-500">Noise</a> 使用<a href="https://github.com/rcy1314/echo-noise" target="_blank" rel="noopener noreferrer" class="text-orange-400 hover:text-orange-500">Ech0-Noise</a>发布</div>',
    rssTitle: 'Noise的说说笔记',
  rssDescription: '一个说说笔记~',
  rssAuthorName: 'Noise',
  rssFaviconURL: '/favicon.ico',
  aboutMarkdown: '# 关于我\n\n这里是一个默认的个人简介示例：\n\n- 喜欢记录与分享\n- 热爱开源与学习\n- 持续打磨产品体验\n\n欢迎通过友链或留言与我交流！',
  linksTitle: '友情链接',
  linksDescription: '推荐站点和朋友们的主页',
  commentPageTitle: '留言',
  commentPageDescription: '欢迎留下你的看法',
  aboutPageTitle: '关于本站',
  aboutPageDescription: '这里是站点的介绍与说明',
    enableGithubCard: false,
    // PWA 默认（为空时回退到站点设置）
    pwaEnabled: true,
    pwaTitle: '',
    pwaDescription: '',
    pwaIconURL: '',
    announcementText: '欢迎访问我的说说笔记！',
    announcementEnabled: true,
    hitokotoEnabled: true,
    // 评论系统默认值
    commentEnabled: true,
    commentSystem: 'builtin',
    commentEmailEnabled: false,
    commentLoginRequired: false,
    // 音乐默认配置
    musicEnabled: false,
    musicPlaylistId: '2141128031',
    musicSongId: '',
    musicPosition: 'bottom-left',
    musicTheme: 'auto',
    musicLyric: true,
    musicAutoplay: false,
    musicDefaultMinimized: true,
    musicEmbed: false,
    musicCssCdnURL: '',
    musicJsCdnURL: '',
  socialLinks: [
      { name: 'GitHub', url: 'https://github.com/rcy1314', icon: 'i-mdi-github' },
      { name: 'X', url: 'https://x.com/liangwenhao3', icon: 'i-mdi-twitter' },
      { name: '主页', url: 'https://www.noisework.cn/', icon: 'i-mdi-home' },
      { name: '博客', url: 'https://www.noiseblogs.top/', icon: 'i-mdi-notebook' }
  ],
  friendLinks: [
    { title: 'NoiseWork', link: 'https://www.noisework.cn/', icon: 'i-mdi-home', description: '个人主页与作品集合' },
    { title: 'NoiseBlogs', link: 'https://www.noiseblogs.top/', icon: 'i-mdi-notebook', description: '技术随笔与学习记录' }
  ],
    calendarEnabled: true,
    timeEnabled: true,
    // 左栏广告（完全后端驱动，无前端默认）
    leftAdEnabled: true,
    leftAds: [],
    leftAdsIntervalMs: 4000
  };

// 修改 fetchConfig 方法
const fetchConfig = async () => {
    try {
        frontendConfig.value = { ...defaultConfig };
        const res = await getRequest<any>('frontend/config', undefined, { credentials: 'include' })
        if (res && res.code === 1 && res.data && res.data.frontendSettings) {
            const settings = res.data.frontendSettings
            const booleanKeys = ['enableGithubCard', 'pwaEnabled', 'announcementEnabled', 'hitokotoEnabled', 'commentEnabled', 'commentEmailEnabled', 'commentLoginRequired', 'musicEnabled', 'musicLyric', 'musicAutoplay', 'musicDefaultMinimized', 'musicEmbed', 'calendarEnabled', 'timeEnabled', 'leftAdEnabled', 'welcomeUseAdmin']
            Object.keys(frontendConfig.value).forEach(key => {
                if (settings[key] !== null && settings[key] !== undefined) {
                    if (key === 'backgrounds' && Array.isArray(settings[key])) {
                        frontendConfig.value.backgrounds = [...settings[key]]
                    } else if (key === 'socialLinks' && Array.isArray(settings[key])) {
                        frontendConfig.value.socialLinks = [...settings[key]]
                    } else if (key === 'leftAds' && Array.isArray(settings[key])) {
                        frontendConfig.value.leftAds = [...settings[key]]
                    } else if (key === 'friendLinks' && Array.isArray(settings[key])) {
                        frontendConfig.value.friendLinks = [...settings[key]]
                    } else if (booleanKeys.includes(key)) {
                        const v = settings[key]
                        frontendConfig.value[key] = (v === true || v === 'true')
                    } else {
                        const v = settings[key]
                        frontendConfig.value[key] = typeof v === 'string' ? v.trim() : v
                    }
                }
            })
            const defaultTheme = (settings.defaultContentTheme || 'light').trim()
            if (typeof window !== 'undefined' && !localStorage.getItem('contentTheme')) {
              contentTheme.value = defaultTheme === 'light' ? 'light' : 'dark'
              document.documentElement.className = contentTheme.value === 'dark' ? 'dark' : ''
            } else if (typeof window !== 'undefined') {
              document.documentElement.className = contentTheme.value === 'dark' ? 'dark' : ''
            }
        }
        if (!frontendConfig.value.backgrounds?.length) {
            frontendConfig.value.backgrounds = [...defaultConfig.backgrounds]
        }
        if (frontendConfig.value.backgrounds.length > 0) {
            const randomIndex = Math.floor(Math.random() * frontendConfig.value.backgrounds.length)
            currentImage.value = frontendConfig.value.backgrounds[randomIndex]
        }
        isLoaded.value = true
    } catch (error) {
        console.error('获取配置失败:', error)
        frontendConfig.value = { ...defaultConfig }
        isLoaded.value = true
    }
}

const getIconName = (item: any) => {
  const icon = (item && typeof item.icon === 'string') ? item.icon.trim() : ''
  if (icon) return icon
  const url = String(item?.url || '').trim()
  if (!url) return 'i-mdi-link-variant'
  const u = url.toLowerCase()
  if (u.startsWith('mailto:')) return 'i-mdi-email'
  let host = ''
  try {
    host = new URL(u).hostname.toLowerCase()
  } catch {}
  if (!host) return 'i-mdi-link-variant'
  if (host.includes('github.com')) return 'i-mdi-github'
  if (host.includes('twitter.com') || host.includes('x.com')) return 'i-mdi-twitter'
  if (host.includes('t.me') || host.includes('telegram.me') || host.includes('telegram.org')) return 'i-mdi-telegram'
  if (host.includes('discord.com') || host.includes('discord.gg')) return 'i-mdi-discord'
  if (host.includes('youtube.com') || host.includes('youtu.be')) return 'i-mdi-youtube'
  if (host.includes('linkedin.com')) return 'i-mdi-linkedin'
  if (host.includes('facebook.com')) return 'i-mdi-facebook'
  if (host.includes('instagram.com')) return 'i-mdi-instagram'
  if (host.includes('bilibili.com')) return 'i-mdi-bilibili'
  if (host.includes('zhihu.com')) return 'i-mdi-zhihu'
  if (host.includes('medium.com')) return 'i-mdi-medium'
  if (host.includes('tiktok.com') || host.includes('douyin.com')) return 'i-mdi-tiktok'
  if (host.includes('weibo.com')) return 'i-mdi-sina-weibo'
  if (host.includes('wechat') || host.includes('weixin')) return 'i-mdi-wechat'
  if (host.includes('rss') || host.endsWith('/rss')) return 'i-mdi-rss'
  return 'i-mdi-link-variant'
}
const normalizeHost = (raw: string) => {
  const s = String(raw || '').trim()
  try {
    const u = new URL(s)
    return u.hostname.replace(/^www\./, '')
  } catch {
    return s
  }
}
const isSocialHost = (host: string) => {
  const h = String(host || '').toLowerCase()
  return (
    h.includes('github.com') ||
    h.includes('twitter.com') || h.includes('x.com') ||
    h.includes('t.me') || h.includes('telegram.me') || h.includes('telegram.org') ||
    h.includes('discord.com') || h.includes('discord.gg') ||
    h.includes('youtube.com') || h.includes('youtu.be') ||
    h.includes('linkedin.com') || h.includes('facebook.com') ||
    h.includes('instagram.com') || h.includes('bilibili.com') ||
    h.includes('zhihu.com') || h.includes('medium.com') ||
    h.includes('tiktok.com') || h.includes('douyin.com') ||
    h.includes('weibo.com') || h.includes('wechat') || h.includes('weixin')
  )
}
const simpleLinks = computed(() => {
  const raw = (frontendConfig as any).value?.socialLinks ?? (frontendConfig as any).socialLinks
  const arr = Array.isArray(raw) ? raw : []
  return arr.filter((it: any) => String(it?.url || '').trim() !== '').map((it: any) => ({ ...it, name: it?.name || it?.url }))
})
const visibleSocialLinks = computed(() => simpleLinks.value)
const friendLinksList = computed(() => {
  const raw = (frontendConfig as any).value?.friendLinks ?? (frontendConfig as any).friendLinks
  const arr = Array.isArray(raw) ? raw : []
  return arr.filter((it: any) => String(it?.link || '').trim() !== '').map((it: any) => ({
    title: String(it?.title || it?.link || '').trim(),
    link: String(it?.link || '').trim(),
    icon: String(it?.icon || '').trim(),
    description: String(it?.description || '').trim(),
  }))
})
const HITOKOTO_FALLBACKS = [
  '身为冒险者，如果安静的老死在床上，那简直就是耻辱！',
  '愿你出走半生，归来仍是少年。',
  '愿你眼里有光，心里有海。'
]
const hitokotoText = ref(HITOKOTO_FALLBACKS[0])
const currentImage = ref('')
const isLoaded = ref(false)
const imageLoading = ref(false)
// 添加图片预加载函数
const preloadImages = async (images: string[]) => {
  const loadImage = (src: string) => {
    return new Promise((resolve) => {
      const img = new Image()
      img.src = src
      img.onload = () => resolve(src)
      img.onerror = () => resolve(null)
    })
  }
  
  // 并行预加载所有图片
  await Promise.all(images.map(src => loadImage(src)))
}
// 添加配置更新事件监听
// 移除重复绑定的 frontend-config-updated 监听，避免多次拉取导致卡顿

onUnmounted(() => {
    // 移除事件监听
    window.removeEventListener('frontend-config-updated', fetchConfig);
});
// 优化背景切换函数
const changeBackground = async () => {
  if (imageLoading.value) return
  imageLoading.value = true
  
  const newIndex = Math.floor(Math.random() * frontendConfig.value.backgrounds.length)
  const newImage = frontendConfig.value.backgrounds[newIndex]
  
  if (newImage === currentImage.value) {
    imageLoading.value = false
    return
  }

  // 使用更小的缩略图
  const thumbnailImage = `${newImage}?imageView2/2/w/10/blur/1/q/10`
  currentImage.value = thumbnailImage

  // 使用 requestAnimationFrame 优化渲染
  requestAnimationFrame(() => {
    const img = new Image()
    img.src = newImage
    img.onload = () => {
      requestAnimationFrame(() => {
        currentImage.value = newImage
        imageLoading.value = false
      })
    }
    img.onerror = () => {
      imageLoading.value = false
    }
  })
}
// 定义页面元数据
definePageMeta({
  title: '说说笔记'
})

// 设置动态标题
let headUpdateTimer: any = null
const scheduleHeadUpdate = () => {
  if (headUpdateTimer) clearTimeout(headUpdateTimer)
  headUpdateTimer = setTimeout(() => updateTitle(), 100)
}
const updateTitle = () => {
  const title = (frontendConfig.value.pwaTitle || frontendConfig.value.siteTitle || '说说笔记').trim()
  const icon = (frontendConfig.value.rssFaviconURL || '/favicon.ico').trim()
  const pwaIcon = (
    frontendConfig.value.pwaIconURL && frontendConfig.value.pwaIconURL.trim() !== ''
      ? frontendConfig.value.pwaIconURL.trim()
      : (icon.toLowerCase().endsWith('.png') ? icon : '/android-chrome-192x192.png')
  )
  const description = (frontendConfig.value.pwaDescription || frontendConfig.value.description || '').trim()
  useHead({
    title,
    meta: [
      { key: 'description', name: 'description', content: description },
      { key: 'theme-color', name: 'theme-color', content: '#000000' }
    ],
    link: [
      { key: 'icon-32', rel: 'icon', type: 'image/png', href: '/favicon-32x32.png', sizes: '32x32' },
      { key: 'shortcut-icon-32', rel: 'shortcut icon', type: 'image/png', href: '/favicon-32x32.png', sizes: '32x32' },
      { key: 'icon-fallback', rel: 'icon', href: icon },
      ...(frontendConfig.value.pwaEnabled ? [
        { key: 'manifest', rel: 'manifest', href: '/manifest.webmanifest' },
        { key: 'apple-touch', rel: 'apple-touch-icon', href: pwaIcon, sizes: '180x180' },
        { key: 'pwa-192', rel: 'icon', href: pwaIcon, sizes: '192x192' },
        { key: 'pwa-512', rel: 'icon', href: (pwaIcon.toLowerCase().endsWith('.png') ? pwaIcon : '/android-chrome-512x512.png'), sizes: '512x512' }
      ] : [])
    ]
  })
}

// 监听配置变化
watch(() => [frontendConfig.value.pwaEnabled, frontendConfig.value.pwaTitle, frontendConfig.value.pwaIconURL, frontendConfig.value.pwaDescription, frontendConfig.value.siteTitle, frontendConfig.value.rssFaviconURL, frontendConfig.value.description], () => {
  scheduleHeadUpdate()
}, { immediate: true })
const subtitleEl = ref<HTMLElement | null>(null)
  const tags = ref([])
// 添加标签更新处理函数
const handleTagsUpdate = async () => {
  await fetchTags()
}
// 获取所有标签
const fetchTags = async () => {
  try {
    const res = await getRequest<any>('messages/tags')
    if (res && res.code === 1) {
      tags.value = res.data || []
    } else {
      tags.value = []
    }
  } catch (error) {
    console.error('获取标签失败:', error)
    tags.value = []
  }
}

// 图片与状态
const images = ref<any[]>([])
const status = ref<any>(null)
const fetchImages = async () => {
  try {
    const r = await getRequest<any>('messages/images')
    if (r && r.code === 1 && Array.isArray(r.data)) images.value = r.data
  } catch {}
}
const fetchStatus = async () => {
  try {
    const r = await getRequest<any>('status')
    if (r && r.code === 1) status.value = r.data
  } catch {}
}
const popularTags = computed(() => {
  const arr = Array.isArray(tags.value) ? [...tags.value] : []
  return arr.sort((a: any, b: any) => (b.count || 0) - (a.count || 0)).slice(0, 8)
})
const recommendedImages = computed(() => images.value.slice(0, 60))
const imageSrc = (img: any) => {
  const url = typeof img === 'string' ? img : (img?.image_url || img?.url)
  const base = useRuntimeConfig().public.baseApi || '/api'
  return url?.startsWith('http') ? url : `${base}${url}`
}

const leftAds = computed(() => (
  Array.isArray((frontendConfig.value as any).leftAds)
    ? (frontendConfig.value as any).leftAds.filter((ad: any) => String(ad?.imageURL || '').trim() !== '')
    : []
))
const imgSrc = (raw: string) => {
  const base = useRuntimeConfig().public.baseApi || '/api'
  const s = String(raw || '').trim()
  return s.startsWith('http') ? s : (s ? `${base}${s}` : '')
}
const adIndex = ref(0)
const currentAd = computed(() => leftAds.value[Math.max(0, Math.min(adIndex.value, leftAds.value.length - 1))] || { imageURL: '', linkURL: '', description: '' })
let adTimer: any
onMounted(() => {
  const interval = Number(frontendConfig.leftAdsIntervalMs || 5000)
  if (leftAds.value.length > 1) {
    adTimer = setInterval(() => {
      adIndex.value = (adIndex.value + 1) % leftAds.value.length
    }, Math.max(2000, interval))
  }
})
onUnmounted(() => { if (adTimer) clearInterval(adTimer) })

// 最新评论（右栏）
const recentComments = ref<Array<any>>([])
const recentIndex = ref(0)
const containsImage = (s: string) => /!\[[^\]]*\]\([^)]*\)/.test(String(s || '')) || /<img[^>]*>/i.test(String(s || ''))
const visibleRecentComments = computed(() => {
  const items = recentComments.value.filter((c) => !containsImage(c?.content || ''))
  if (items.length <= 8) return items
  const start = recentIndex.value % items.length
  const end = (start + 8) % items.length
  return start < end ? items.slice(start, end) : [...items.slice(start), ...items.slice(0, end)]
})
const shortText = (s: string) => {
  const noMdImg = String(s || '').replace(/!\[[^\]]*\]\([^)]*\)/g, '')
  const noHtmlImg = noMdImg.replace(/<img[^>]*>/gi, '')
  const t = noHtmlImg.replace(/<[^>]+>/g, '').replace(/\s+/g, ' ').trim()
  return t.length > 20 ? t.slice(0, 20) : (t || '评论')
}
const escapeHTML = (s: string) => String(s || '')
  .replace(/&/g, '&amp;')
  .replace(/</g, '&lt;')
  .replace(/>/g, '&gt;')
  .replace(/"/g, '&quot;')
  .replace(/'/g, '&#39;')
const recentHTML = (s: string) => {
  const raw = String(s || '')
  // Replace Markdown image ![alt](url)
  const replacedMd = raw.replace(/!\[[^\]]*\]\(([^)]+)\)/g, (_m, url) => `<img class=\"recent-inline-img\" src=\"${escapeHTML(url)}\" alt=\"img\" loading=\"lazy\" />`)
  // Replace HTML <img ... src="..."> with small inline
  const replacedHtmlImg = replacedMd.replace(/<img[^>]*src=[\"']([^\"']+)[\"'][^>]*>/gi, (_m, url) => `<img class=\"recent-inline-img\" src=\"${escapeHTML(url)}\" alt=\"img\" loading=\"lazy\" />`)
  // Escape other HTML
  const safe = escapeHTML(replacedHtmlImg)
  // Undo escaping for the injected <img> tags
  return safe.replace(/&lt;img class=\"recent-inline-img\" src=\"([^\"]+)\" alt=\"img\" loading=\"lazy\" \/&gt;/g, '<img class="recent-inline-img" src="$1" alt="img" loading="lazy" />')
}
const pravatar = (seed: string) => `https://i.pravatar.cc/64?u=${encodeURIComponent(seed || Math.random().toString(36).slice(2))}`
const qqNumberFromEmail = (mail?: string) => {
  const m = String(mail || '').trim()
  const match = m.match(/^(\d+)@qq\.com$/i)
  return match ? match[1] : ''
}
const qqAvatarUrl = (qq: string) => qq ? `https://q1.qlogo.cn/g?b=qq&nk=${qq}&s=100` : ''
const hashCode = (s: string) => { let h = 0; const t = String(s || ''); for (let i = 0; i < t.length; i++) { h = ((h << 5) - h) + t.charCodeAt(i); h |= 0 } return Math.abs(h) }
const initialsText = (s: string) => {
  const t = String(s || '').trim()
  if (!t) return 'U'
  if (/^[\u4e00-\u9fa5]/.test(t)) return t.slice(0, 1)
  const parts = t.split(/\s+|[_.-]+/).filter(Boolean)
  const a = (parts[0] || t)[0] || 'U'
  const b = (parts[1] || '')[0] || ''
  return (a + b).toUpperCase()
}
const initialsAvatar = (seed: string, size = 40) => {
  const txt = initialsText(seed)
  const hue = hashCode(seed) % 360
  const bg = `hsl(${hue},70%,50%)`
  const svg = `<svg xmlns="http://www.w3.org/2000/svg" width="${size}" height="${size}"><rect width="${size}" height="${size}" rx="${size/2}" fill="${bg}"/><text x="50%" y="55%" font-size="${Math.floor(size*0.5)}" text-anchor="middle" fill="#fff" font-family="-apple-system,Segoe UI,Roboto,Helvetica,Arial">${txt}</text></svg>`
  return 'data:image/svg+xml;charset=UTF-8,' + encodeURIComponent(svg)
}
const avatarPlaceholder = computed(() => {
  const raw = String((frontendConfig.value.avatarURL || '')).trim()
  const base = useRuntimeConfig().public.baseApi || '/api'
  if (raw) return raw.startsWith('http') ? raw : `${base}${raw}`
  const icon = String(frontendConfig.value.rssFaviconURL || '/favicon.svg').trim()
  return icon
})
const qqAvatarCandidates = (qq: string, size = 100) => [
  qq ? `https://q1.qlogo.cn/g?b=qq&nk=${qq}&s=${size}` : '',
  qq ? `https://q2.qlogo.cn/g?b=qq&nk=${qq}&s=${size}` : '',
  qq ? `https://q.qlogo.cn/g?b=qq&nk=${qq}&s=${size}` : ''
].filter(Boolean)
const recentAvatar = (c: any) => {
  const qq = qqNumberFromEmail(c?.mail || '')
  const name = String(c?.nick || '').trim()
  const arr = qqAvatarCandidates(qq)
  return (arr[0] || '') || initialsAvatar(name || c?.mail || 'anonymous') || avatarPlaceholder.value || pravatar(name || c?.mail || '')
}
const onRecentAvatarError = (e: Event, seed: string) => {
  const img = e.target as HTMLImageElement
  const mailAttr = (img?.dataset?.mail || '') as string
  const qq = String(img?.dataset?.qq || '')
  const tries = Number(img?.dataset?.try || 0)
  const candidates = qqAvatarCandidates(qq)
  if (qq && tries < candidates.length) {
    const nextIdx = Math.min(tries + 1, candidates.length - 1)
    const next = candidates[nextIdx]
    img.dataset.try = String(nextIdx)
    img.src = next
    return
  }
  const fallback = initialsAvatar(seed || mailAttr || 'anonymous') || avatarPlaceholder.value || pravatar(seed || mailAttr || 'anonymous')
  if (img && fallback) img.src = fallback
}
const BASE_API = useRuntimeConfig().public.baseApi || '/api'
const loadRecentComments = async () => {
  try {
    let page = 1
    let pageSize = 50
    let total = 0
    const collected: any[] = []
    for (let round = 0; round < 3; round++) {
      const resp = await fetch(`${BASE_API}/messages/page`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json', 'Accept': 'application/json' },
        credentials: 'include',
        body: JSON.stringify({ page, pageSize })
      })
      if (!resp.ok) break
      const js = await resp.json().catch(() => null)
      const items = js?.data?.items || []
      total = Number(js?.data?.total || 0)
      const ids = items.map((m: any) => m.id).filter((id: number) => id !== guestbookMessageId.value)
      const tasks = ids.map(async (id: number) => {
        try {
          const r = await fetch(`${BASE_API}/messages/${id}/comments`, { credentials: 'include', headers: { 'Accept': 'application/json' } })
          if (!r.ok) return []
          const d = await r.json()
          const arr = Array.isArray(d?.data) ? d.data : []
          return arr.map((c: any) => ({ id: c.id, nick: c.nick, mail: c.mail, content: c.content, created_at: c.created_at, message_id: id }))
        } catch { return [] }
      })
      const results = await Promise.all(tasks)
      const flat = ([] as any[]).concat(...results)
      collected.push(...flat)
      const lastPage = Math.max(1, Math.ceil(total / pageSize))
      if (collected.length >= 24 || page >= lastPage) break
      page += 1
    }
    collected.sort((a, b) => new Date(b.created_at || 0).getTime() - new Date(a.created_at || 0).getTime())
    recentComments.value = collected.filter((c: any) => c.message_id !== guestbookMessageId.value)
      .filter((c: any) => !containsImage(c?.content || ''))
      .map((c: any) => ({ id: c.id, nick: c.nick, mail: c.mail, content: c.content }))
      .slice(0, 24)
    if (recentComments.value.length > 0) return
  } catch {}
  recentComments.value = ['期待','对方的','发个','困难','路口提示','测试','你好','加油','真不错','赞同','有趣','哈哈','有用','收藏','不错','灵感','记录']
    .map((t, i) => ({ id: i + 1, nick: '匿名', mail: '', content: t }))
}
onMounted(async () => {
  await loadRecentComments()
  setInterval(() => { recentIndex.value = (recentIndex.value + 1) % Math.max(1, recentComments.value.length) }, 3000)
})

// 绑定 Fancybox 以支持推荐图集预览
onMounted(() => {
  try { (window as any).Fancybox?.bind?.('[data-fancybox]', {}) } catch {}
})

const recentItemClass = computed(() => (isDark.value ? 'bg-[rgba(36,43,50,0.75)] text-white' : 'bg-slate-100 text-slate-800'))
const recentTextClass = computed(() => (isDark.value ? 'text-white' : 'text-slate-800'))

// 监听前端配置更新事件，保存后主动刷新配置
onMounted(() => {
  const handler = () => fetchConfig()
  window.addEventListener('frontend-config-updated', handler)
  // 初始拉取
  fetchConfig()
  fetchTags()
  fetchImages()
  fetchStatus()
  onUnmounted(() => window.removeEventListener('frontend-config-updated', handler))
})
// 标签点击处理
const handleTagClick = async (tag: string) => {
  try {
    const encodedTag = encodeURIComponent(tag.trim())
    const res = await getRequest<any>(`messages/tags/${encodedTag}`, undefined, { credentials: 'include' })
    if (res && res.code === 1 && Array.isArray(res.data)) {
      if (messageList.value) {
        messageList.value.handleSearchResult(res.data)
      }
    } else {
      throw new Error(res?.msg || '获取标签内容失败')
    }
  } catch (error: any) {
    console.error('获取标签消息失败:', error)
    useToast().add({
      title: '获取标签消息失败',
      description: error.message || '服务器错误，请稍后重试',
      color: 'red',
      timeout: 3000
    })
  }
}
// 修改打字效果函数
const startTypeEffect = () => {
  if (!subtitleEl.value) return
  
  let index = 0
  let isDeleting = false
  let isWaiting = false
  
  const typeInterval = setInterval(() => {
    if (!subtitleEl.value) {
      clearInterval(typeInterval)
      return
    }
    if (isWaiting) return

    if (!isDeleting) {
      // 打字过程
      subtitleEl.value!.textContent = frontendConfig.value.subtitleText.slice(0, index + 1)
      index++
      
      if (index >= frontendConfig.value.subtitleText.length) {
        isWaiting = true
        setTimeout(() => {
          isDeleting = true
          isWaiting = false
        }, 2000)
      }
    } else {
      // 删除过程
      index--
      subtitleEl.value!.textContent = frontendConfig.value.subtitleText.slice(0, index)
      
      if (index <= 0) {
        isWaiting = true
        subtitleEl.value!.textContent = ''
        setTimeout(() => {
          isDeleting = false
          isWaiting = false
          index = 0
        }, 1000)
      }
    }
  }, 100)

  return typeInterval
}

// 修改 onMounted 钩子
onMounted(async () => {
  try {
    // 确保在任何异步操作之前设置加载状态
    isLoaded.value = false;

    // 使用 requestIdleCallback 延迟加载非关键组件
    window.requestIdleCallback = window.requestIdleCallback || ((cb) => setTimeout(cb, 1))
    
    // 关键内容优先加载
    await fetchConfig()
    if (frontendConfig.value.musicEnabled) {
      await initNMP()
    }
    
    // 非关键内容延迟加载
    requestIdleCallback(async () => {
      await fetchTags()
      await preloadImages(frontendConfig.value.backgrounds)
    })

    // 确保配置加载完成后再执行后续操作
    if (frontendConfig.value.backgrounds.length > 0) {
      const initialImage = frontendConfig.value.backgrounds[
        Math.floor(Math.random() * frontendConfig.value.backgrounds.length)
      ]
      
      // 先加载低质量版本
      const lowQualityImage = `${initialImage}?imageView2/2/w/100/q/30`
      currentImage.value = lowQualityImage

      // 后台预加载其他图片
      requestIdleCallback(async () => {
        await preloadImages(frontendConfig.value.backgrounds)
      })
      
      // 加载高质量初始图片
      const img = new Image()
      img.src = initialImage
      img.onload = () => {
        requestAnimationFrame(() => {
          currentImage.value = initialImage
          isLoaded.value = true; // 在高质量图片加载完成后设置为已加载
        })
      }
      img.onerror = () => {
        // 若加载失败，直接结束加载遮罩，避免阻塞交互
        isLoaded.value = true
      }
      // 最长等待时间到达后也结束加载遮罩，避免卡顿
      setTimeout(() => {
        if (!isLoaded.value) isLoaded.value = true
      }, 2000)
    }
    
    // 启动打字效果
    const typeInterval = startTypeEffect()
    onUnmounted(() => {
      if (typeInterval) {
        clearInterval(typeInterval)
      }
    })

    // 添加事件监听
    window.addEventListener('frontend-config-updated', async (event: CustomEvent) => {
      await fetchConfig()
      if (frontendConfig.value.backgrounds?.length > 0) {
        const randomIndex = Math.floor(Math.random() * frontendConfig.value.backgrounds.length)
        const newImage = frontendConfig.value.backgrounds[randomIndex]
        currentImage.value = newImage
      }
    })

  } catch (error) {
    console.error('初始化失败:', error)
    isLoaded.value = true
  }
})
const selectedDate = ref<Date>(new Date())
const currentTime = ref<Date>(new Date())
let timeTimer: any = null
onMounted(() => {
  timeTimer = setInterval(() => { currentTime.value = new Date() }, 1000)
})
onUnmounted(() => {
  if (timeTimer) { clearInterval(timeTimer); timeTimer = null }
})
const formatTime = (d: Date) => {
  const hh = String(d.getHours()).padStart(2, '0')
  const mm = String(d.getMinutes()).padStart(2, '0')
  const ss = String(d.getSeconds()).padStart(2, '0')
  return `${hh}:${mm}:${ss}`
}
const formatDate = (d: Date) => {
  const y = d.getFullYear()
  const m = String(d.getMonth() + 1).padStart(2, '0')
  const day = String(d.getDate()).padStart(2, '0')
  const wk = ['周日','周一','周二','周三','周四','周五','周六'][d.getDay()]
  return `${y}-${m}-${day} ${wk}`
}
</script>

<style>
html, body {
  margin: 0;
  padding: 0;
  width: 100%;
  height: 100%;
  overflow: hidden;
  overscroll-behavior: none;
}
.header-subtitle {
  position: absolute;
  top: calc(50% + 60px);  /* 增加间距 */
  left: 50%;
  transform: translate(-50%, -50%);
  color: white;
  font-size: 1rem;
  text-shadow: none;
  white-space: nowrap;
}

.layout-toggle {
  position: absolute;
  top: 14px;
  right: 14px;
  z-index: 2;
}

@media screen and (max-width: 768px) {
  .header-subtitle {
    font-size: 0.9rem;
    top: calc(50% + 45px);  /* 保持移动端的适配 */
  }
}
.background-container {
  width: 100%;
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  overflow: hidden;
  z-index: 0;
  background-color: black;
}

.background-container::before {
  content: '';
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background-image: var(--bg-image);
  background-size: cover;
  background-position: center;
  background-attachment: fixed; /* 确保背景固定 */
  filter: blur(8px);
  z-index: -1;
}

.content-wrapper {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  height: 100vh;
  overflow-y: auto;
  z-index: 1;
  pointer-events: auto;
  cursor: default;
}

.moments-header {
  margin-bottom: 20px;
}

.header-image {
  position: relative;
  width: 100%;
  height: 300px;
  background-size: cover;
  background-position: center;
  border-radius: 18px;
  overflow: hidden;
  transition: background-image 0.3s ease;
  will-change: background-image;
  transform: translateZ(0);
  margin-top: 0;
  box-shadow: 0 8px 24px rgba(0,0,0,0.12);
}

:global(html.dark) .header-image { box-shadow: 0 10px 28px rgba(0,0,0,0.35); }


.header-title {
  position: absolute;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  color: #ffffff;
  font-size: 2.5rem;
  font-weight: bold;
  text-shadow: none;
  margin: 0;
  white-space: nowrap;
  transition: font-size 0.3s ease;
}
.hero-tabs {
  position: absolute;
  left: 50%;
  bottom: 18px;
  transform: translateX(-50%);
  display: flex;
  gap: 8px;
  background: rgba(0, 0, 0, 0.3);
  backdrop-filter: blur(6px);
  padding: 6px 8px;
  border-radius: 9999px;
  flex-wrap: nowrap;
  overflow-x: auto;
  overflow-y: hidden;
  max-width: 90vw;
  white-space: nowrap;
}
/* 统一弹窗底色样式与搜索弹窗一致 */
.search-card { background: #ffffff; color: #111827; border: 1px solid #e5e7eb; border-radius: 16px; }
html.dark .search-card { background: rgba(36,43,50,0.95); color: #fff; border: 1px solid rgba(255,255,255,0.08); }
.hero-tab {
  font-size: 0.85rem;
  line-height: 1;
  padding: 8px 12px;
  border-radius: 9999px;
  color: #fff;
  opacity: 0.8;
  white-space: nowrap;
  flex-shrink: 0;
}
.hero-tab.active {
  opacity: 1;
  background: rgba(255, 255, 255, 0.18);
}
@media (max-width: 480px) {
  .hero-tabs { gap: 6px; padding: 6px 8px; }
  .hero-tab { font-size: 0.8rem; padding: 8px 12px; }
}
.hero-tabs::-webkit-scrollbar { display: none; }
.theme-default { --accent: #ff8c3a }
.theme-mint { --accent: #3bb273 }
.theme-rose { --accent: #e85d75 }
.theme-slate { --accent: #5c7cfa }
.sidebar-title { color: var(--accent, inherit) }
.clock-card { padding: 0 }
.clock-display { font-weight: 700; font-size: 1.8rem; letter-spacing: 2px }
.clock-date { margin-top: 6px; font-size: 0.85rem; opacity: 0.7 }
.calendar-card :deep(.u-calendar) { width: 100% }
.calendar-card :deep(.u-calendar) {
  width: 100%;
  display: block;
}
.calendar-card :deep(.u-calendar *) {
  color: inherit;
}
.calendar-card :deep(.u-calendar-header) {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 6px 8px;
  color: var(--accent, #ff8c3a);
}
.calendar-card :deep(table) {
  width: 100%;
}
.calendar-card :deep(th) {
  color: inherit;
  font-weight: 500;
}
.calendar-card :deep(td) {
  padding: 2px;
}
.calendar-card :deep(td > button) {
  width: 100%;
  border-radius: 8px;
  padding: 6px 0;
  transition: background-color .15s ease, transform .1s ease;
  border: 1px solid rgba(0,0,0,0.06);
}
.calendar-card :deep(td > button:hover) {
  transform: scale(1.02);
  background: rgba(0,0,0,0.06);
}
.calendar-card :deep(td > button[aria-selected="true"]) {
  background: var(--accent, #ff8c3a);
  color: #fff;
}
.calendar-card :deep(td > button[aria-current="date"]) {
  outline: 2px solid var(--accent, #ff8c3a);
}
@media screen and (max-width: 768px) {
  .content-wrapper {
   padding: 0.5rem; 
  }
  /* 优化移动端滚动性能 */
  .message-list-container {
    transform: translateZ(0);
    will-change: transform;
  }
  .container-fixed {
    width: 100%;
    margin: 0 auto;
    padding-bottom: 0.2rem; /* 底部内边距 */
    padding-left: 0.5rem; padding-right: 0.5rem; /* 移动端左右保持对称内边距 */
  }
  .background-container::before {
    filter: blur(4px); /* 减少模糊度提升性能 */
    background-attachment: scroll; /* 移动端使用普通滚动 */
    transform: scale(1.08);
  }
  
  .content-wrapper {
    overscroll-behavior: contain;
    position: absolute;
    top: 0;
    left: 0;
    right: 0;
    height: 100vh;
    overflow-y: auto;
    z-index: 1;
    pointer-events: auto;
    padding: 0.25rem; /* 收紧移动端外层边距，提升内容占比 */
  }
  
  .background-container {
    position: fixed;
  }
}
@media screen and (max-width: 768px) {
  .header-title {
    font-size: 1.8rem;
    top: 35%;
  }
  .header-subtitle {
    top: calc(35% + 50px); 
    font-size: 0.9rem;
  }
}

.profile-info {
  position: absolute;
  bottom: 20px;
  right: 20px;
  display: flex;
  flex-direction: row-reverse;  /* 改变方向，头像在右侧 */
  align-items: center;
  gap: 10px;
  max-width: 80%;  /* 限制最大宽度 */
}

.avatar {
  width: 50px;
  height: 50px;
  border-radius: 50%;
  border: 2px solid white;
  object-fit: cover;
  cursor: pointer;
  transition: transform 0.3s ease;
}

.avatar:hover {
  transform: scale(1.1);
}

.profile-text {
  text-align: left;  /* 改为左对齐 */
  min-width: 0;  /* 允许内容收缩 */
  overflow-x: auto;  /* 允许横向滚动 */
  scrollbar-width: none;  /* 隐藏滚动条 (Firefox) */
  -ms-overflow-style: none;  /* 隐藏滚动条 (IE/Edge) */
  padding: 5px 0;
}
.profile-text .title {
  font-size: 1.2rem;
  font-weight: bold;
}

.profile-text .description {
  font-size: 0.9rem;
  opacity: 0.9;
}
.profile-text::-webkit-scrollbar {
  display: none;
}

.profile-text .title {
font-size: 1.2rem;
font-weight: bold;
color: #fcfafb;  
text-shadow: none;
white-space: nowrap;  /* 防止换行 */

}

.profile-text .description {
  font-size: 0.9rem;
  color: #fcfafb; 
  text-shadow: none;
  white-space: nowrap;
  opacity: 0.95;
}
.u-container {
  backdrop-filter: blur(4px);
  border-radius: 8px;
  margin: 0 auto;
  max-width: 1440px;
  width: 100%;
  position: relative;
  z-index: 1;
  box-sizing: border-box;
  overflow-x: hidden;
  cursor: default;
}

.message-list-container { cursor: default; }

.loading {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
  background: rgba(0, 0, 0, 0.7); /* 更改背景色 */
  backdrop-filter: blur(8px);
  z-index: 1000;
  gap: 15px;
  opacity: 1;
  transition: opacity 0.3s ease;
}

.loading-text {
  font-size: 16px;
  color: #fff; /* 更改文字颜色 */
  text-shadow: none;
}

.rainbow-spinner {
  width: 50px;
  height: 50px;
  border-radius: 50%;
  border: 4px solid transparent;
  border-top: 4px solid #FF0000;
  border-right: 4px solid #00FF00;
  border-bottom: 4px solid #0000FF;
  border-left: 4px solid #FF00FF;
  animation: spin 1s linear infinite;
  will-change: transform;
}

@keyframes spin {
  0% { transform: rotate(0deg); }
  100% { transform: rotate(360deg); }
}
/* 添加新的样式 */
.message-list-container {
  position: relative;
  z-index: 10;
}

.container-fixed {
  min-height: 100vh;
  pointer-events: auto;
  max-width: 1420px;
}

/* 确保背景不会遮挡评论框 */
.background-container::before {
  z-index: -1;
}
.heatmap-container {
  background: rgba(0, 0, 0, 0.3);
  backdrop-filter: blur(4px);
  border-radius: 8px;
  padding: 1rem;
  margin-bottom: 1rem;
}
.scroll-buttons {
  position: fixed;
  right: 20px;
  bottom: 20px;
  display: flex;
  flex-direction: column;
  gap: 10px;
  z-index: 1000;
}
.scroll-button {
  width: 2.25rem;
  height: 2.25rem;
  border-radius: 9999px;
  transition: transform 0.15s ease, background-color 0.15s ease, box-shadow 0.15s ease;
  backdrop-filter: blur(6px);
  display: inline-flex;
  align-items: center;
  justify-content: center;
  opacity: 0.95;
}
.scroll-button:hover {
  transform: scale(1.06);
}

.layout-container.grid-3 {
  display: grid;
  grid-template-columns: var(--sidebar-width, 320px) 1fr var(--sidebar-width, 320px);
  gap: var(--grid-gap, 16px);
  align-items: start;
}
.layout-container.grid-2 {
  display: grid;
  grid-template-columns: var(--sidebar-width, 320px) 1fr;
  gap: var(--grid-gap, 16px);
}
.layout-container.grid-1 {
  display: block;
}
.layout-container { --sidebar-width: 320px; --grid-gap: 16px; }
.left-col, .right-col { position: sticky; top: 0; align-self: start; height: fit-content; }
.center-col { min-width: 0; }
.sidebar-card {
  border-radius: 10px;
  background: #ffffff;
  color: #111827;
  border: 1px solid #e5e7eb;
}
/* 三栏容器在浅色模式统一白色背景（不影响深色与背景图层） */
:global(html:not(.dark)) .left-col,
:global(html:not(.dark)) .center-col,
:global(html:not(.dark)) .right-col {
  background: #ffffff;
  border: 1px solid rgba(0,0,0,0.08);
  border-radius: 16px;
  padding: 8px;
  box-shadow: 0 2px 10px rgba(0,0,0,0.06);
}
:global(html:not(.dark)) .layout-container.grid-3 { gap: 18px; }
/* 统一浅色模式下所有卡片底色为纯白 */
:global(html:not(.dark)) :where(.u-card, .u-card-body, .u-card__body, .u-card-header, .u-card__header) { background-color: #ffffff !important; }
:global(html:not(.dark)) :where(.bg-gray-50, .bg-gray-100, .bg-gray-200, .bg-gray-300, .bg-slate-50, .bg-slate-100, .bg-slate-200) { background-color: #ffffff !important; }
:global(html:not(.dark)) :where(.border-gray-200, .border-gray-300, .border-slate-200) { border-color: rgba(0,0,0,0.08) !important; }
/* 去除指定卡片内部默认留白，使内容铺满容器 */
.no-padding-card :deep(.u-card-body),
.no-padding-card :deep(.u-card__body) { padding: 0 !important; }
.no-padding-card :deep(.u-card-header),
.no-padding-card :deep(.u-card__header) { padding: 8px 12px !important; }
/* 统一压缩所有侧栏卡片的默认内边距 */
.sidebar-card :deep(.u-card-body),
.sidebar-card :deep(.u-card__body) { padding: 0 !important; }
.sidebar-card :deep(.u-card-header),
.sidebar-card :deep(.u-card__header) { padding: 6px 8px !important; }
html.dark .sidebar-card { background: rgba(36,43,50,0.95); color: #fff; border: 1px solid rgba(255,255,255,0.08); }
html.dark .sidebar-card :where(.bg-white,.bg-gray-50,.bg-gray-100,.bg-gray-700,.bg-gray-800,.bg-gray-900) { background-color: rgba(36,43,50,0.95) !important; }
html.dark .sidebar-card :where(.text-black,.text-gray-900,.text-gray-800) { color: #fff !important; }
html.dark .sidebar-card :where(.border,.border-gray-200,.border-gray-300,.border-gray-600,.border-gray-700) { border-color: rgba(255,255,255,0.08) !important; }
.profile-card {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 0;
}
.auth-actions { margin-top: 6px; display: flex; align-items: center; justify-content: center; gap: 10px; }
.auth-btn { width: 36px; height: 36px; padding: 0; border-radius: 9999px; display: inline-flex; align-items: center; justify-content: center; background: transparent !important; border: none !important; box-shadow: none !important; }
.auth-btn:hover { background: transparent !important; transform: translateY(-1px); }
.avatar-lg {
  width: 72px;
  height: 72px;
  border-radius: 50%;
  object-fit: cover;
  border: 2px solid rgba(255,255,255,0.8);
  cursor: pointer;
  transition: transform .18s ease;
}
.avatar-lg:hover { transform: scale(1.05); }
/* 头像在线状态独立定位 */
.avatar-wrap { display:inline-block; position:relative; }
.avatar-status { position:absolute; right:-2px; top:-2px; width:10px; height:10px; border-radius:9999px; border:2px solid rgba(255,255,255,0.9); }
/* 桌面登录/注册圆形按钮禁用悬停效果，仅显示 Tooltip */
.plain-icon-btn { background: transparent !important; transition: none !important; box-shadow: none !important; }
.plain-icon-btn:hover { background: transparent !important; transform: none !important; box-shadow: none !important; }
/* 随机一言文本容器：无图标、自动换行显示全文 */
.hitokoto-container { margin: 0; padding: 0; }
.hitokoto-text { overflow: hidden; white-space: normal; word-break: break-word; overflow-wrap: anywhere; font-size: 14px; font-weight: 500; line-height: 1.45; }
.profile-name {
  margin-top: 8px;
  font-weight: 600;
}
.profile-desc {
  margin-top: 4px;
  font-size: 0.875rem;
  opacity: 0.8;
}
.social-list { display:flex; flex-wrap:wrap; gap:10px; padding:0; justify-content:center; align-items:center; }
.social-item { position:relative; display:inline-flex; align-items:center; justify-content:center; width:32px; height:32px; border-radius:8px; }
.social-item { color: inherit; text-decoration: none; }
.social-item:hover { transform: scale(1.06); transition: transform .12s ease; }
.social-icon-img { width:28px; height:28px; border-radius:6px; object-fit:cover; display:inline-block; }
.social-item::after { content: attr(data-label); position:absolute; bottom:calc(100% + 2px); left:50%; transform: translateX(-50%); white-space:nowrap; padding:4px 8px; font-size:12px; border-radius:6px; pointer-events:none; opacity:0; transition: opacity .12s ease; }
:global(html.dark) .social-item::after { background: rgba(36,43,50,0.95); color:#fff; border:1px solid rgba(255,255,255,0.1); }
:global(html:not(.dark)) .social-item::after { background:#fff; color:#111; border:1px solid rgba(0,0,0,0.08); }
.social-item:hover::after { opacity:1; }
.sidebar-title {
  font-weight: 600;
  padding: 8px 10px;
}

.links-page { padding-bottom: 120px; }
.links-page .section-title { font-weight: 600; font-size: 14px; margin-bottom: 12px; padding: 0; border-radius: 0; display: inline-flex; align-items: center; gap: 6px; }
.card-title { font-weight: 700; font-size: 18px; margin-bottom: 14px; padding: 0; border-radius: 0; display: block; }
.section-subtitle { text-align: center; font-size: 13px; opacity: 0.75; margin-top: -6px; margin-bottom: 12px; }
.section-title-light { color: #111827; }
.section-title-dark { color: #fff; }
.link-grid { display: grid; grid-template-columns: repeat(auto-fit, minmax(170px, 1fr)); gap: 8px; justify-items: center; }
.link-card { display: flex; align-items: center; gap: 8px; padding: 8px; border-radius: 8px; text-decoration: none; transition: background-color .16s ease, border-color .16s ease, transform .12s ease; }
.link-card { width: 100%; max-width: 300px; }
.link-card-light { background: #fff; color: #111827; border: 1px solid rgba(0,0,0,0.08); box-shadow: 0 1px 2px rgba(0,0,0,0.06); }
.link-card-dark { background: rgba(36,43,50,0.85); color: #fff; border: 1px solid rgba(255,255,255,0.12); box-shadow: 0 1px 2px rgba(255,255,255,0.06); }
.link-card:hover { background-color: rgba(0,0,0,0.02); transform: translateY(-1px); }
.link-card-dark:hover { background-color: rgba(255,255,255,0.08); }
.link-avatar { width: 42px; height: 42px; border-radius: 0; display: inline-flex; align-items: center; justify-content: center; }
.link-avatar-light { background: #eef2ff; color: #4f46e5; border: 1px solid rgba(0,0,0,0.06); }
.link-avatar-dark { background: rgba(255,255,255,0.12); color: #c7d2fe; border: 1px solid rgba(255,255,255,0.16); }
.link-content { flex: 1; min-width: 0; }
.link-title { font-weight: 600; font-size: 13px; line-height: 1.2; white-space: nowrap; overflow: hidden; text-overflow: ellipsis; }
.link-sub { display: none; }

.about-header { display: flex; align-items: center; gap: 14px; border-radius: 0; padding: 14px; }
.about-header-light { background: #fff; border: 1px solid rgba(0,0,0,0.08); box-shadow: 0 1px 2px rgba(0,0,0,0.06); }
.about-header-dark { background: rgba(36,43,50,0.85); border: 1px solid rgba(255,255,255,0.12); box-shadow: 0 1px 2px rgba(255,255,255,0.06); }
.about-avatar { width: 72px; height: 72px; border-radius: 0; object-fit: cover; }
.about-info { display: flex; flex-direction: column; gap: 4px; }
.about-title { font-size: 20px; font-weight: 700; }
.about-sub { font-size: 13px; opacity: 0.75; }
.about-desc { font-size: 14px; opacity: 0.9; }

@media screen and (max-width: 1024px) {
  .layout-container.grid-3 {
    display: block;
  }
  .layout-container.grid-2 {
    display: block;
  }
}
@media screen and (max-width: 768px) {
  .layout-container.grid-3,
  .layout-container.grid-2 {
    display: block !important;
  }
  .left-col,
  .right-col {
    display: none !important;
  }
}
@media screen and (max-width: 1280px) {
  .layout-container { --sidebar-width: 280px; --grid-gap: 16px; }
}
</style>
<style>
/* 推荐图集图片 box 效果与悬停动画 */
.recommend-image-box {
  width: 100%;
  height: 100%;
  border-radius: 12px;
  object-fit: cover;
  display: block;
  transition: transform .18s ease, box-shadow .18s ease, filter .18s ease;
  box-shadow: 0 1px 2px rgba(0,0,0,0.10);
}
.recommend-grid { display: grid; grid-template-columns: repeat(3, 1fr); grid-auto-rows: 88px; gap: 6px; }
.recommend-grid a { display:block; height:100%; }
.recommend-image-box:hover {
  transform: translate3d(0,0,0) scale(1.03);
  box-shadow: 0 6px 18px rgba(0,0,0,0.28);
  filter: saturate(1.06) contrast(1.02);
}
@media (prefers-color-scheme: dark) {
  .recommend-image-box { box-shadow: 0 1px 2px rgba(255,255,255,0.06); }
  .recommend-image-box:hover { box-shadow: 0 8px 22px rgba(255,255,255,0.12); }
}
.scroll-list { height: 64px; overflow-y: auto; -webkit-overflow-scrolling: touch; padding: 2px 2px; scroll-behavior: smooth; }
.recent-inline-img { display:inline-block; width:18px; height:18px; object-fit:cover; border-radius:4px; vertical-align:middle; margin: -2px 2px 0 2px; }
.ad-wrap { position: relative; aspect-ratio: var(--ad-aspect, 1 / 1); }
.ad-image { width: 100%; height: 100%; object-fit: contain; transition: filter .12s ease, transform .12s ease; }
.ad-wrap::before { content: ""; position: absolute; inset: 0; background-image: var(--ad-bg); background-size: cover; background-position: center; filter: blur(12px) brightness(0.95); transform: scale(1.05); }
.ad-wrap > .ad-image { position: relative; z-index: 1; }
.ad-overlay { position:absolute; inset:0; display:flex; align-items:center; justify-content:center; opacity:0; transition: opacity .12s ease; pointer-events:none; z-index: 2; }
.ad-overlay-box { max-width: 90%; max-height: 70%; overflow-y: auto; padding: 8px 10px; border-radius: 10px; font-size: 14px; line-height: 1.5; word-break: break-word; overflow-wrap: anywhere; }
:global(html.dark) .ad-overlay-box { background: rgba(36,43,50,0.90); color:#f59e0b !important; border:1px solid rgba(255,255,255,0.12); box-shadow: 0 6px 18px rgba(0,0,0,0.28); }
:global(html.dark) .ad-overlay-box a { color:#f59e0b !important; text-decoration:none; }
:global(html:not(.dark)) .ad-overlay-box { background:#ffffff; color:#f59e0b !important; border:1px solid rgba(0,0,0,0.08); box-shadow: 0 6px 18px rgba(0,0,0,0.12); }
:global(html:not(.dark)) .ad-overlay-box a { color:#f59e0b !important; text-decoration:none; }
.ad-wrap:hover .ad-overlay { opacity:1; }
.ad-wrap:hover .ad-image { filter: contrast(0.95) brightness(0.9); }
.scroll-images { height: 240px; overflow-y: auto; -webkit-overflow-scrolling: touch; padding-right: 2px; }
/* 标签三栏栅格与滚动容器 */
.scroll-tags { max-height: 160px; overflow-y: auto; -webkit-overflow-scrolling: touch; padding-right: 2px; min-height: 0; overscroll-behavior: contain; }
.tag-grid { display: grid; grid-template-columns: repeat(3, 1fr); grid-auto-rows: minmax(28px, auto); gap: 6px; }
@media screen and (max-width: 1024px) { .tag-grid { grid-template-columns: repeat(2, 1fr); } .scroll-tags { max-height: 120px; } }
@media screen and (max-width: 1024px) {
  .scroll-images { height: 180px; }
  .recommend-grid { grid-auto-rows: 56px; gap: 4px; }
}
</style>
  mq?.addEventListener?.('change', (e: MediaQueryListEvent) => {
    isMobile.value = e.matches
    if (isMobile.value) {
      layoutState.value = 'single'
      localStorage.setItem('homeLayoutMobile', 'single')
    } else {
      const saved = localStorage.getItem('homeLayoutDesktop') as any
      layoutState.value = (saved as any) || 'three'
    }
  })
/* 移除未使用的 header-inner 样式 */
