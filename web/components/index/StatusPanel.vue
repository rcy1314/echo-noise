<template>
 
    <div class="min-h-screen w-screen" :class="theme.pageBg">
        <div class="min-h-screen w-full">
      <aside
        class="w-72 h-screen overflow-y-auto backdrop-blur-md flex flex-col fixed left-0 top-0 z-40 transition-transform duration-300 border-r"
        :class="[{ 'translate-x-0': sidebarOpen, '-translate-x-full md:translate-x-0': !sidebarOpen }, theme.sidebarBg, theme.border, theme.sidebarText]"
      >
        <div class="px-4 py-4 border-b border-slate-700/40 flex flex-col items-center gap-2">
          <img :src="avatarSrc" class="w-14 h-14 rounded-full ring-2 ring-indigo-400/60 shadow-lg object-cover" alt="avatar" />
          <div class="w-full text-center">
            <div class="font-semibold truncate">{{ userStore.user?.username || frontendConfig.username || '未登录' }}</div>
            <div class="text-xs" :class="theme.mutedText">总笔记 {{ userStore?.status?.total_messages || 0 }}</div>
          </div>
        </div>
        <nav class="flex-1 overflow-y-auto px-2 py-3 space-y-2">
          <button class="w-full flex justify-center items-center gap-2 px-3 py-2 rounded-lg transition shadow" :class="[theme.navBtnBg, theme.navBtnHoverBg]" @click="setActive('system')">
            <UIcon name="i-heroicons-cpu-chip" class="w-5 h-5 text-indigo-300" />
            <span class="text-sm text-center">系统信息</span>
          </button>
          <button class="w-full flex justify-center items-center gap-2 px-3 py-2 rounded-lg transition shadow" :class="[theme.navBtnBg, theme.navBtnHoverBg]" @click="setActive('user')">
            <UIcon name="i-heroicons-user-circle" class="w-5 h-5 text-indigo-300" />
            <span class="text-sm text-center">用户信息</span>
          </button>
          <button class="w-full flex justify-center items-center gap-2 px-3 py-2 rounded-lg transition shadow" :class="[theme.navBtnBg, theme.navBtnHoverBg]" @click="setActive('site')">
            <UIcon name="i-heroicons-wrench-screwdriver" class="w-5 h-5 text-indigo-300" />
            <span class="text-sm text-center">网站配置</span>
          </button>
          <div v-if="isAdmin" class="space-y-2">
            <button class="w-full flex justify-center items-center gap-2 px-3 py-2 rounded-lg transition shadow" :class="[theme.navBtnBg, theme.navBtnHoverBg]" @click="setActive('site-register')">
              <UIcon name="i-heroicons-user-plus" class="w-5 h-5 text-indigo-300" />
              <span class="text-sm text-center">注册开关</span>
            </button>
            <button class="w-full flex justify-center items-center gap-2 px-3 py-2 rounded-lg transition shadow" :class="[theme.navBtnBg, theme.navBtnHoverBg]" @click="setActive('site-pwa')">
              <UIcon name="i-heroicons-rocket-launch" class="w-5 h-5 text-indigo-300" />
              <span class="text-sm text-center">PWA 模式</span>
            </button>
            <button class="w-full flex justify-center items-center gap-2 px-3 py-2 rounded-lg transition shadow" :class="[theme.navBtnBg, theme.navBtnHoverBg]" @click="setActive('site-github-card')">
              <UIcon name="i-mdi-github" class="w-5 h-5 text-indigo-300" />
              <span class="text-sm text-center">GitHub 卡片</span>
            </button>
            <button class="w-full flex justify-center items-center gap-2 px-3 py-2 rounded-lg transition shadow" :class="[theme.navBtnBg, theme.navBtnHoverBg]" @click="setActive('site-github-login')">
              <UIcon name="i-mdi-github" class="w-5 h-5 text-indigo-300" />
              <span class="text-sm text-center">GitHub 登录</span>
            </button>
            <button class="w-full flex justify-center items-center gap-2 px-3 py-2 rounded-lg transition shadow" :class="[theme.navBtnBg, theme.navBtnHoverBg]" @click="setActive('site-announcement')">
              <UIcon name="i-heroicons-megaphone" class="w-5 h-5 text-indigo-300" />
              <span class="text-sm text-center">公告栏</span>
            </button>
            <button class="w-full flex justify-center items-center gap-2 px-3 py-2 rounded-lg transition shadow" :class="[theme.navBtnBg, theme.navBtnHoverBg]" @click="setActive('site-music')">
              <UIcon name="i-heroicons-musical-note" class="w-5 h-5 text-indigo-300" />
              <span class="text-sm text-center">音乐配置</span>
            </button>
            <button class="w-full flex justify-center items-center gap-2 px-3 py-2 rounded-lg transition shadow" :class="[theme.navBtnBg, theme.navBtnHoverBg]" @click="setActive('site-default-theme')">
              <UIcon name="i-heroicons-swatch" class="w-5 h-5 text-indigo-300" />
              <span class="text-sm text-center">默认主题</span>
            </button>
            <button class="w-full flex justify-center items-center gap-2 px-3 py-2 rounded-lg transition shadow" :class="[theme.navBtnBg, theme.navBtnHoverBg]" @click="setActive('site-configs')">
              <UIcon name="i-heroicons-cog-6-tooth" class="w-5 h-5 text-indigo-300" />
              <span class="text-sm text-center">站点文案</span>
            </button>
            <button class="w-full flex justify-center items-center gap-2 px-3 py-2 rounded-lg transition shadow" :class="[theme.navBtnBg, theme.navBtnHoverBg]" @click="setActive('comments')">
              <UIcon name="i-heroicons-chat-bubble-left-right" class="w-5 h-5 text-indigo-300" />
              <span class="text-sm text-center">评论系统</span>
            </button>
            <button class="w-full flex justify-center items-center gap-2 px-3 py-2 rounded-lg transition shadow" :class="[theme.navBtnBg, theme.navBtnHoverBg]" @click="setActive('email')">
              <UIcon name="i-heroicons-envelope" class="w-5 h-5 text-indigo-300" />
              <span class="text-sm text-center">邮件设置</span>
            </button>
            <button class="w-full flex justify-center items-center gap-2 px-3 py-2 rounded-lg transition shadow" :class="[theme.navBtnBg, theme.navBtnHoverBg]" @click="setActive('admin-users')">
              <UIcon name="i-heroicons-shield-check" class="w-5 h-5 text-indigo-300" />
              <span class="text-sm text-center">管理员用户</span>
            </button>
          </div>
          <button class="w-full flex justify-center items-center gap-2 px-3 py-2 rounded-lg transition shadow" :class="[theme.navBtnBg, theme.navBtnHoverBg]" @click="setActive('notify')">
            <UIcon name="i-heroicons-bell-alert" class="w-5 h-5 text-indigo-300" />
            <span class="text-sm text-center">推送配置</span>
          </button>
          <button class="w-full flex justify-center items-center gap-2 px-3 py-2 rounded-lg transition shadow" :class="[theme.navBtnBg, theme.navBtnHoverBg]" @click="setActive('db')">
            <UIcon name="i-heroicons-database" class="w-5 h-5 text-indigo-300" />
            <span class="text-sm text-center">数据库管理</span>
          </button>
        </nav>
        <div class="px-4 py-3 border-t border-slate-700/40">
          <div class="text-xs text-slate-400">当前版本: latest</div>
          <div class="mt-2 flex items-center gap-2">
            <UButton size="xs" color="gray" variant="soft" :loading="versionInfo.checking" class="shadow-md" @click="checkVersion">{{ versionInfo.checking ? '检测中...' : '检查版本发布时间' }}</UButton>
          </div>
          <div v-if="versionInfo.hasUpdate" class="mt-2 text-orange-400 flex items-center gap-2">
            <UIcon name="i-heroicons-arrow-up-circle" class="w-4 h-4" />
            <span>最近更新于 {{ versionInfo.latestVersion }}</span>
          </div>
        </div>
      </aside>
      <main class="flex flex-col flex-1 h-full md:ml-72 min-h-screen overflow-y-auto" :class="theme.text">
        <div class="md:hidden flex items-center justify-between px-4 py-3 border-b" :class="[theme.headerBg, theme.border, theme.text]">
          <div class="flex items-center gap-2">
            <button class="p-2 rounded-lg bg-slate-800/70 shadow" @click="sidebarOpen = !sidebarOpen"><UIcon name="i-heroicons-bars-3" class="w-5 h-5" /></button>
            <span class="font-semibold">系统管理面板</span>
          </div>
          <div class="flex items-center gap-2">
            <UButton variant="soft" :color="panelTheme === 'light' ? 'gray' : 'white'" class="shadow" @click="$router.push('/')">返回首页</UButton>
            <UButton v-if="isLogin" icon="i-heroicons-power" color="red" variant="solid" @click="handleLogout">退出登录</UButton>
          </div>
        </div>
        <div v-if="sidebarOpen" class="fixed inset-0 bg-black/40 md:hidden" @click="sidebarOpen=false"></div>
        <div class="flex-1 p-4 pb-24 flex flex-col w-full" :class="isAdmin ? 'gap-4' : 'gap-0'">
          <div class="col-span-12">
            <h1 class="text-2xl md:text-3xl font-bold text-center" :class="theme.text">系统管理面板</h1>
          </div>
          <div class="col-span-12">
            <div class="rounded-xl border" :class="[theme.cardBg, theme.border]">
              <div class="px-4 py-3 flex items-center justify-between">
                <div class="flex items-center gap-4">
                  <span :class="theme.text">配色</span>
                  <div class="flex items-center gap-3">
                    <div class="flex items-center"><URadio v-model="panelTheme" value="dark" class="mr-2" /><span :class="panelTheme === 'dark' ? theme.text : 'text-slate-400'">暗黑</span></div>
                    <div class="flex items-center"><URadio v-model="panelTheme" value="midnight" class="mr-2" /><span :class="panelTheme === 'midnight' ? theme.text : 'text-slate-400'">深蓝</span></div>
                    <div class="flex items-center"><URadio v-model="panelTheme" value="slate" class="mr-2" /><span :class="panelTheme === 'slate' ? theme.text : 'text-slate-400'">石板</span></div>
                    <div class="flex items-center"><URadio v-model="panelTheme" value="light" class="mr-2" /><span :class="panelTheme === 'light' ? theme.text : 'text-slate-400'">明亮</span></div>
                  </div>
                </div>
                
                <div class="flex items-center gap-2">
                  <UButton size="sm" color="green" class="shadow" @click="saveAdminTheme">保存</UButton>
                </div>
              </div>
            </div>
          </div>

          <div id="comments-section" class="col-span-12">
            <div class="rounded-xl border shadow-xl" :class="[theme.cardBg, theme.border]">
              <div class="flex items-center justify-between px-4 py-3">
                <div class="font-semibold flex items-center gap-2" :class="theme.text">
                  <UIcon name="i-heroicons-chat-bubble-left-right" class="w-5 h-5" />
                  <span>评论系统</span>
                </div>
              </div>
              <div class="px-4 pb-4">
                <CommentsSettings v-model:config="frontendConfig" />
              </div>
            </div>
          </div>
          <div id="system-section" class="col-span-12">
            <div class="rounded-xl border shadow-xl" :class="[theme.cardBg, theme.border]">
              <div class="px-4 py-3 flex flex-wrap items-center gap-6">
                <div>
                  <span :class="theme.text">系统管理员</span>
                  <span class="inline-flex items-center font-medium rounded-md px-2 py-1 gap-1 ml-2 text-sm bg-primary-50 dark:bg-primary-400 dark:bg-opacity-10 text-slate-800 dark:text-slate-200">{{ userStore?.status?.username }}</span>
                </div>
                <div>
                  <span :class="theme.text">当前用户</span>
                  <span class="inline-flex items-center font-medium rounded-md px-2 py-1 gap-1 ml-2 text-sm bg-primary-50 dark:bg-primary-400 dark:bg-opacity-10 text-slate-800 dark:text-slate-200">{{ isLogin ? (userStore.user?.username || '未登录') : '未登录' }}</span>
                </div>
                <div>
                  <span :class="theme.text">笔记总数</span>
                  <span class="inline-flex items-center font-medium rounded-md px-2 py-1 gap-1 ml-2 text-sm bg-primary-50 dark:bg-primary-400 dark:bg-opacity-10 text-slate-800 dark:text-slate-200">{{ userStore?.status?.total_messages }} 条</span>
                </div>
              </div>
            </div>
          </div>
          
          <div id="site-music-section" class="col-span-12">
            <div class="rounded-xl border shadow-xl" :class="[theme.cardBg, theme.border]">
              <div class="flex items-center justify-between px-4 py-3">
                <div class="font-semibold flex items-center gap-2" :class="theme.text">
                  <UIcon name="i-heroicons-musical-note" class="w-5 h-5" />
                  <span>音乐配置</span>
                </div>
              </div>
              <div class="px-4 pb-4">
                <div class="rounded-lg p-4 space-y-4" :class="theme.subtleBg">
                  <div class="flex items-center justify-between">
                    <span class="font-medium" :class="theme.text">启用音乐播放器</span>
                    <USwitch v-model="frontendConfig.musicEnabled" />
                  </div>
                  <div class="grid grid-cols-1 md:grid-cols-2 gap-3">
                    <div>
                      <label class="text-sm mb-1 block" :class="theme.mutedText">歌单 ID</label>
                      <UInput v-model="frontendConfig.musicPlaylistId" placeholder="网易云歌单ID（可选）" />
                    </div>
                    <div>
                      <label class="text-sm mb-1 block" :class="theme.mutedText">歌曲 ID</label>
                      <UInput v-model="frontendConfig.musicSongId" placeholder="网易云歌曲ID（可选）" />
                    </div>
                    <div>
                      <label class="text-sm mb-1 block" :class="theme.mutedText">位置</label>
                      <USelect v-model="frontendConfig.musicPosition" :options="[{label:'左下',value:'bottom-left'},{label:'右下',value:'bottom-right'}]" />
                    </div>
                    <div>
                      <label class="text-sm mb-1 block" :class="theme.mutedText">主题</label>
                      <USelect v-model="frontendConfig.musicTheme" :options="[{label:'自动',value:'auto'},{label:'深色',value:'dark'},{label:'浅色',value:'light'}]" />
                    </div>
                    <div class="flex items-center gap-2">
                      <span class="text-sm" :class="theme.mutedText">显示歌词</span>
                      <USwitch v-model="frontendConfig.musicLyric" />
                    </div>
                    <div class="flex items-center gap-2">
                      <span class="text-sm" :class="theme.mutedText">自动播放</span>
                      <USwitch v-model="frontendConfig.musicAutoplay" />
                    </div>
                    <div class="flex items-center gap-2">
                      <span class="text-sm" :class="theme.mutedText">默认最小化</span>
                      <USwitch v-model="frontendConfig.musicDefaultMinimized" />
                    </div>
                    <div class="flex items-center gap-2 md:col-span-2">
                      <span class="text-sm" :class="theme.mutedText">展示模式</span>
                      <USelect v-model="musicEmbedMode" :options="[{label:'嵌入',value:'embed'},{label:'浮动',value:'float'}]" />
                    </div>
                  </div>
                  <div class="flex justify-end gap-2">
                    <UButton variant="soft" color="gray" @click="resetMusicConfig">重置</UButton>
                    <UButton color="green" @click="saveMusicConfig">保存</UButton>
                  </div>
                </div>
              </div>
            </div>
          </div>

          <div id="user-section" class="col-span-12" v-if="isLogin">
            <div class="rounded-xl border shadow-xl" :class="[theme.cardBg, theme.border]">
              <div class="flex items-center justify-between px-4 py-3" :class="theme.text">
                <div class="font-semibold">用户信息配置</div>
              </div>
              <div class="px-4 pb-4 grid grid-cols-1 md:grid-cols-3 gap-y-2 md:gap-x-4">
                <div class="rounded-lg p-3 h-full" :class="theme.subtleBg">
                  <div class="flex justify-between items-center mb-2">
                    <span :class="theme.mutedText">用户名</span>
                    <UButton size="sm" @click="editUserInfo.username = !editUserInfo.username" :color="editUserInfo.username ? 'gray' : 'green'" :variant="editUserInfo.username ? 'soft' : 'solid'" class="shadow">{{ editUserInfo.username ? '取消' : '编辑' }}</UButton>
                  </div>
                  <div v-if="editUserInfo.username">
                    <UInput v-model="userForm.username" placeholder="新用户名" class="w-full mb-2" />
                    <div class="flex justify-end gap-2">
                      <UButton @click="updateUsername" color="primary" class="shadow">保存</UButton>
                    </div>
                  </div>
                    <div v-else>
                      <UBadge color="primary" variant="soft" class="ml-0 text-xs px-2 py-1 rounded-md !text-slate-800 dark:!text-slate-200">{{ userStore.user?.username }}</UBadge>
                    </div>
                </div>
                <div class="rounded-lg p-3 h-full md:col-span-2" :class="theme.subtleBg">
                  <div class="flex justify-between items-center mb-2">
                    <div class="flex items-center gap-2">
                      <span :class="theme.mutedText">API Token</span>
                      <UBadge color="primary" variant="subtle" class="text-xs px-2 py-1 rounded-md !text-primary-600 dark:!text-primary-300">{{ userToken ? '已生成' : '未生成' }}</UBadge>
                      <UButton size="xs" :loading="regeneratingToken" @click="regenerateToken" color="gray" variant="soft" class="shadow text-xs px-2 py-1 rounded-md text-slate-600 dark:text-slate-200" title="重新生成将使旧 Token 失效">重新生成</UButton>
                    </div>
                  </div>
                  <div v-if="userToken" class="mb-2">
                    <div class="flex items-center gap-2 w-full flex-nowrap">
                      <UInput v-model="userToken" :type="showToken ? 'text' : 'password'" readonly class="font-mono text-sm flex-1 min-w-0" />
                      <UButton :icon="showToken ? 'i-heroicons-eye-slash' : 'i-heroicons-eye'" color="gray" variant="ghost" @click="showToken = !showToken" :title="showToken ? '隐藏' : '显示'" />
                      <UButton icon="i-heroicons-clipboard" color="gray" variant="ghost" @click="copyToken" title="复制 Token" />
                    </div>
                    <p class="text-xs mt-1" :class="theme.mutedText">请妥善保管此 Token</p>
                  </div>
                  <div v-else>
                    <p :class="theme.mutedText">暂无 Token</p>
                  </div>
                </div>
                <div class="rounded-lg p-3 md:col-span-3" :class="theme.subtleBg">
                  <div class="flex justify-between items-center mb-2">
                    <span :class="theme.mutedText">修改密码</span>
                    <UButton size="sm" @click="editUserInfo.password = !editUserInfo.password" :color="editUserInfo.password ? 'gray' : 'green'" :variant="editUserInfo.password ? 'soft' : 'solid'" class="shadow">{{ editUserInfo.password ? '取消' : '编辑' }}</UButton>
                  </div>
                  <div v-if="editUserInfo.password">
                    <div class="w-full mb-2 flex items-center gap-2">
                      <UInput v-model="userForm.oldPassword" :type="showOldPassword ? 'text' : 'password'" placeholder="当前密码" class="flex-1" />
                      <UButton :icon="showOldPassword ? 'i-heroicons-eye-slash' : 'i-heroicons-eye'" color="gray" variant="ghost" @click="showOldPassword = !showOldPassword" />
                    </div>
                    <div class="w-full mb-2 flex items-center gap-2">
                      <UInput v-model="userForm.newPassword" :type="showNewPassword ? 'text' : 'password'" placeholder="新密码" class="flex-1" />
                      <UBadge :color="passwordStrengthColor" variant="soft">{{ passwordStrengthLabel }}</UBadge>
                      <UButton :icon="showNewPassword ? 'i-heroicons-eye-slash' : 'i-heroicons-eye'" color="gray" variant="ghost" @click="showNewPassword = !showNewPassword" />
                    </div>
                    <div class="w-full mb-2 flex items-center gap-2">
                      <UInput v-model="userForm.confirmPassword" :type="showConfirmPassword ? 'text' : 'password'" placeholder="确认新密码" class="flex-1" />
                      <UButton :icon="showConfirmPassword ? 'i-heroicons-eye-slash' : 'i-heroicons-eye'" color="gray" variant="ghost" @click="showConfirmPassword = !showConfirmPassword" />
                    </div>
                    <div class="flex justify-end gap-2">
                      <UButton @click="updatePassword" :disabled="!canSavePassword" color="primary" class="shadow">保存</UButton>
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </div>

          <div id="site-section" class="col-span-12">
            <div class="rounded-xl border shadow-xl" :class="[theme.cardBg, theme.border]">
              <div class="flex items-center justify-between px-4 py-3">
                <div class="font-semibold flex items-center gap-2" :class="theme.text">
                  <UIcon name="i-heroicons-cog-6-tooth" class="w-5 h-5" />
                  <span>网站配置</span>
                </div>
              </div>
              <div class="px-4 pb-4 space-y-4">
                <div id="site-register-section" class="flex items-center rounded-lg p-3 justify-between" :class="theme.subtleBg">
                  <div class="flex items-center gap-2" :class="theme.text"><UIcon name="i-heroicons-user-plus" class="w-4 h-4" /> <span>新用户注册</span></div>
                  <div class="flex items-center gap-4">
                    <div class="flex items-center">
                      <URadio v-model="registerEnabled" :value="true" class="mr-2" />
                      <span :class="registerEnabled ? theme.text : 'text-slate-400'">允许</span>
                    </div>
                    <div class="flex items-center">
                      <URadio v-model="registerEnabled" :value="false" class="mr-2" />
                      <span :class="!registerEnabled ? theme.text : 'text-slate-400'">不允许</span>
                    </div>
                    <UButton color="green" @click="saveRegisterConfig" class="shadow">保存</UButton>
                  </div>
                </div>
                <div id="site-pwa-section" class="rounded-lg p-4" :class="theme.subtleBg">
                  <div class="flex justify-between items-center mb-3">
                    <div class="flex items-center gap-2" :class="theme.text"><UIcon name="i-heroicons-rocket-launch" class="w-4 h-4" /> <span>PWA 模式</span></div>
                    <div class="flex items-center gap-4">
                      <UToggle v-model="frontendConfig.pwaEnabled" />
                      <UButton color="green" @click="savePWAConfig" class="shadow">保存</UButton>
                    </div>
                  </div>
  <div class="grid grid-cols-1 md:grid-cols-2 gap-3">
                    <div>
                      <label :class="[theme.mutedText, 'text-sm mb-1 block']">PWA 标题</label>
                      <UInput v-model="frontendConfig.pwaTitle" :placeholder="frontendConfig.siteTitle || '说说笔记'" />
                    </div>
                    <div>
                      <label :class="[theme.mutedText, 'text-sm mb-1 block']">PWA 图标</label>
                      <UInput v-model="frontendConfig.pwaIconURL" :placeholder="'/favicon.ico'" />
                    </div>
                    <div class="md:col-span-2">
                      <label :class="[theme.mutedText, 'text-sm mb-1 block']">PWA 描述</label>
                      <UTextarea v-model="frontendConfig.pwaDescription" :placeholder="frontendConfig.description || ''" />
  </div>
  </div>
                </div>
                <div id="site-github-card-section" class="flex items-center rounded-lg p-3 justify-between" :class="theme.subtleBg">
                  <div class="flex items-center gap-2" :class="theme.text"><UIcon name="i-mdi-github" class="w-4 h-4" /> <span>GitHub 链接卡片解析</span></div>
                  <div class="flex items-center gap-4">
                    <div class="flex items-center">
                      <URadio v-model="githubCardEnabled" :value="true" class="mr-2" />
                      <span :class="githubCardEnabled ? theme.text : 'text-slate-400'">开启</span>
                    </div>
                    <div class="flex items-center">
                      <URadio v-model="githubCardEnabled" :value="false" class="mr-2" />
                      <span :class="!githubCardEnabled ? theme.text : 'text-slate-400'">关闭</span>
                    </div>
                    <UButton color="green" @click="saveGithubCardConfig" class="shadow">保存</UButton>
                  </div>
                </div>
                <div id="site-announcement-section" class="flex items-center rounded-lg p-3 justify-between" :class="theme.subtleBg">
                  <div class="flex items-center gap-2" :class="theme.text"><UIcon name="i-heroicons-megaphone" class="w-4 h-4" /> <span>公告栏开关</span></div>
                  <div class="flex items-center gap-4">
                    <UToggle v-model="frontendConfig.announcementEnabled" />
                    <UButton color="green" @click="saveConfigItem('announcementEnabled')" class="shadow">保存</UButton>
                  </div>
                </div>
                <div class="rounded-lg p-3 mt-3" :class="theme.subtleBg">
                  <div class="text-sm mb-2" :class="theme.mutedText">公告栏文本</div>
                  <UTextarea v-model="frontendConfig.announcementText" placeholder="请输入公告内容" class="w-full mb-2" />
                  <div class="flex justify-end">
                    <UButton color="primary" class="shadow" @click="saveConfigItem('announcementText')">保存公告文本</UButton>
                  </div>
                </div>
                <div id="site-default-theme-section" class="flex items-center rounded-lg p-3 justify-between" :class="theme.subtleBg">
                  <div class="flex items-center gap-2" :class="theme.text"><UIcon name="i-heroicons-swatch" class="w-4 h-4" /> <span>默认主题色</span></div>
                  <div class="flex items-center gap-4">
                    <div class="flex items-center">
                      <URadio v-model="frontendConfig.defaultContentTheme" value="dark" class="mr-2" />
                      <span :class="frontendConfig.defaultContentTheme === 'dark' ? theme.text : 'text-slate-400'">暗黑</span>
                    </div>
                    <div class="flex items-center">
                      <URadio v-model="frontendConfig.defaultContentTheme" value="light" class="mr-2" />
                      <span :class="frontendConfig.defaultContentTheme === 'light' ? theme.text : 'text-slate-400'">白天</span>
                    </div>
                    <UButton color="green" @click="saveConfigItem('defaultContentTheme')" class="shadow">保存</UButton>
                  </div>
                </div>
                <div id="site-configs-section" class="space-y-4">
                <div v-for="(label, key) in configLabels" :key="key" class="rounded-lg p-3" :class="theme.subtleBg">
                    <div class="flex justify-between items-center mb-2">
                      <span :class="theme.mutedText">{{ label }}</span>
                      <UButton size="sm" @click="editItem[key] = !editItem[key]" :color="editItem[key] ? 'gray' : 'green'" :variant="editItem[key] ? 'soft' : 'solid'" class="shadow">{{ editItem[key] ? '取消' : '编辑' }}</UButton>
                    </div>
                    <div v-if="editItem[key]">
                      <template v-if="key === 'backgrounds'">
                        <div class="space-y-2">
                          <div v-for="(bg, index) in frontendConfig.backgrounds" :key="index" class="flex gap-2">
                            <UInput v-model="frontendConfig.backgrounds[index]" placeholder="输入图片URL" class="flex-1" />
                            <UButton @click="removeBackground(index)" icon="i-heroicons-trash" color="red" variant="ghost" />
                          </div>
                          <div class="flex gap-2">
                            <UButton @click="addBackground" icon="i-heroicons-plus" variant="ghost" class="mr-2">添加链接</UButton>
                            <UButton @click="triggerFileInput" icon="i-heroicons-cloud-arrow-up" variant="ghost">上传图片</UButton>
                          </div>
                        </div>
                      </template>
                      <template v-else-if="key === 'subtitleText'">
                        <UTextarea v-model="frontendConfig[key]" :placeholder="`输入${label}`" class="w-full mb-2" />
                      </template>
                      <template v-else>
                        <UInput v-model="frontendConfig[key]" :placeholder="`输入${label}`" class="w-full mb-2" />
                      </template>
                      <div class="flex justify-end gap-2">
                        <UButton @click="resetConfigItem(key)" variant="ghost" color="gray">重置</UButton>
                        <UButton @click="saveConfigItem(key)" color="primary" class="shadow">保存</UButton>
                      </div>
                    </div>
                    <div v-else>
                      <template v-if="key === 'backgrounds'">
                        <div class="grid grid-cols-3 gap-2">
                          <img v-for="(bg, index) in frontendConfig.backgrounds" :key="index" :src="bg" class="w-full h-20 object-cover rounded cursor-pointer" @click="previewImage(bg)" />
                        </div>
                      </template>
                      <template v-else>
                        <div :class="[theme.subtleBg, theme.border, 'rounded-md p-3 border']">
                          <p :class="[theme.text, 'break-words']">{{ frontendConfig[key] }}</p>
                        </div>
                      </template>
                    </div>
                  </div>
                </div>
                <div v-if="editMode" class="flex justify-end gap-2">
                  <UButton @click="resetConfig" variant="ghost" color="gray">重置</UButton>
                  <UButton @click="saveConfig" color="primary" class="shadow">保存所有更改</UButton>
                </div>
              </div>
            </div>
          </div>

          <div id="email-section" class="col-span-12">
            <div class="rounded-xl border shadow-xl" :class="[theme.cardBg, theme.border]">
              <div class="flex items-center justify-between px-4 py-3">
                <div class="font-semibold flex items-center gap-2" :class="theme.text">
                  <UIcon name="i-heroicons-envelope" class="w-5 h-5" />
                  <span>邮件设置（SMTP）</span>
                </div>
              </div>
              <div class="px-4 pb-4">
                <div class="rounded-lg p-4 space-y-4" :class="theme.subtleBg">
                  <div class="flex items-center justify-between">
                    <span class="font-medium" :class="theme.text">启用邮件</span>
                    <USwitch v-model="smtp.enabled" />
                  </div>
                  <div>
                    <div class="text-sm font-medium mb-2" :class="theme.text">地址</div>
                    <UInput v-model="smtp.from" placeholder="发件地址，如 name@example.com" />
                  </div>
                  <div>
                    <div class="text-sm font-medium mb-2" :class="theme.text">驱动</div>
                    <USelect v-model="smtp.driver" :options="['smtp']" />
                  </div>
                  <div class="text-sm font-semibold mt-1 mb-2" :class="theme.text">SMTP 设置</div>
                  <div class="grid grid-cols-1 md:grid-cols-2 gap-3">
                    <div>
                      <div class="text-sm mb-2" :class="theme.text">主机</div>
                      <UInput v-model="smtp.host" placeholder="smtp.example.com" />
                    </div>
                    <div>
                      <div class="text-sm mb-2" :class="theme.text">端口</div>
                      <UInput v-model="smtp.port" placeholder="465 或 587" />
                    </div>
                    <div>
                      <div class="text-sm mb-2" :class="theme.text">加密协议（小写 ssl 或 tls）</div>
                      <USelect v-model="smtp.encryption" :options="['ssl','tls']" />
                    </div>
                    <div class="md:col-span-1"></div>
                    <div>
                      <div class="text-sm mb-2" :class="theme.text">用户名</div>
                      <UInput v-model="smtp.user" placeholder="通常与发件地址一致" />
                    </div>
                    <div>
                      <div class="text-sm mb-2" :class="theme.text">密码</div>
                      <UInput v-model="smtp.pass" :type="showSmtpPass ? 'text' : 'password'" placeholder="邮箱或应用专用密码" />
                    </div>
                  </div>
                  <div class="flex items-center justify-between mt-2" :class="theme.mutedText">
                    <span class="text-xs">使用上述设置发送测试邮件到：{{ smtp.from || smtp.user || '请先填写地址' }}</span>
                    <UButton :disabled="!(smtp.from || smtp.user)" :loading="testingSmtp" color="primary" @click="testSmtp">发送测试邮件</UButton>
                  </div>
                  <div class="flex justify-end gap-2 mt-3">
                    <UButton variant="soft" color="gray" @click="loadSmtp">刷新</UButton>
                    <UButton color="green" @click="saveSmtp">保存</UButton>
                  </div>
                </div>
              </div>
            </div>
          </div>

          <div id="admin-users-section" class="col-span-12">
            <div class="rounded-xl border shadow-xl" :class="[theme.cardBg, theme.border]">
              <div class="flex items-center justify-between px-4 py-3">
                <div class="font-semibold flex items-center gap-2" :class="theme.text">
                  <UIcon name="i-heroicons-shield-check" class="w-5 h-5" />
                  <span>用户管理</span>
                </div>
              </div>
              <div class="px-4 pb-4">
                <div class="rounded-lg p-3 mb-3" :class="theme.subtleBg">
                  <div class="flex items-center gap-2">
                    <UInput v-model="userSearch" placeholder="搜索用户名或ID" class="flex-1" />
                    <UButton color="primary" variant="soft" @click="refreshUsers">搜索</UButton>
                    <UButton variant="soft" color="gray" @click="refreshUsers">刷新</UButton>
                    <UButton variant="soft" :color="showUsers ? 'gray' : 'indigo'" @click="showUsers=!showUsers">{{ showUsers ? '折叠' : '展开' }}</UButton>
                  </div>
                </div>
                <div v-if="showUsers" class="rounded-lg p-3" :class="theme.subtleBg">
                  <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-2">
                    <div v-for="u in filteredUsers" :key="(u.id ?? u.ID)" class="flex items-center justify-between rounded border px-3 py-2" :class="theme.border">
                      <div class="flex items-center gap-2 min-w-0">
                        <UBadge color="gray" variant="soft">#{{ u.id ?? u.ID }}</UBadge>
                        <span class="truncate" :class="theme.text">{{ u.username ?? u.Username }}</span>
                        <UBadge :color="(u.is_admin ?? u.IsAdmin) ? 'primary' : 'gray'" variant="subtle">{{ (u.is_admin ?? u.IsAdmin) ? '管理员' : '普通' }}</UBadge>
                      </div>
                      <div class="flex items-center gap-2">
                        <UButton :color="(u.is_admin ?? u.IsAdmin) ? 'orange' : 'green'" :variant="(u.is_admin ?? u.IsAdmin) ? 'soft' : 'solid'" class="shadow" @click="confirmToggleAdmin(u)">{{ (u.is_admin ?? u.IsAdmin) ? '取消管理员' : '设为管理员' }}</UButton>
                        <UButton color="red" variant="soft" class="shadow" @click="confirmDeleteUser(u)">删除</UButton>
                      </div>
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </div>
          <UModal v-model="showAdminResetModal">
            <UCard>
              <div class="font-semibold mb-2">重置管理员密码</div>
              <div class="space-y-3">
                <div class="flex items-center gap-2">
                  <UInput v-model="adminReset.newPass" :type="showAdminPassword ? 'text' : 'password'" placeholder="新密码" class="flex-1" />
                  <UBadge :color="adminResetStrengthColor" variant="soft">{{ adminResetStrengthLabel }}</UBadge>
                </div>
                <div class="flex items-center gap-2">
                  <UInput v-model="adminReset.confirmPass" :type="showAdminPassword ? 'text' : 'password'" placeholder="确认新密码" class="flex-1" />
                </div>
                <div class="flex justify-end gap-2">
                  <UButton variant="ghost" color="gray" @click="showAdminResetModal = false">取消</UButton>
                  <UButton :disabled="!canSaveAdminReset" color="primary" @click="resetAdminPassword">保存</UButton>
                </div>
              </div>
            </UCard>
          </UModal>

          <div id="notify-section" class="col-span-12">
            <div class="rounded-xl border shadow-xl" :class="[theme.cardBg, theme.border]">
              <div class="flex items-center justify-between px-4 py-3">
                <div class="font-semibold flex items-center gap-2" :class="theme.text">
                  <UIcon name="i-heroicons-bell-alert" class="w-5 h-5" />
                  <span>推送配置</span>
                </div>
              </div>
              <div class="px-4 pb-4">
                <NotifyPanel v-model:config="notifyConfig" :immediate="true" :subtleBg="theme.subtleBg" :text="theme.text" :mutedText="theme.mutedText" />
              </div>
            </div>
          </div>
          

          <div id="site-github-login-section" class="col-span-12">
            <div class="rounded-xl border shadow-xl" :class="[theme.cardBg, theme.border]">
              <div class="flex items-center justify-between px-4 py-3">
                <div class="font-semibold flex items-center gap-2" :class="theme.text">
                  <UIcon name="i-mdi-github" class="w-5 h-5" />
                  <span>GitHub 登录</span>
                </div>
              </div>
              <div class="px-4 pb-4 space-y-3">
                <div class="grid grid-cols-1 md:grid-cols-3 gap-3">
                  <div class="md:col-span-1 flex items-center gap-2">
                    <span class="text-sm" :class="theme.text">启用</span>
                    <USwitch v-model="frontendConfig.githubOAuthEnabled" />
                  </div>
                  <div>
                    <div class="text-sm mb-2" :class="theme.text">Client ID</div>
                    <UInput v-model="frontendConfig.githubClientId" placeholder="GitHub OAuth App Client ID" />
                  </div>
                  <div>
                    <div class="text-sm mb-2" :class="theme.text">Client Secret</div>
                    <UInput v-model="frontendConfig.githubClientSecret" type="password" placeholder="GitHub OAuth App Client Secret" />
                  </div>
                  <div class="md:col-span-3">
                    <div class="text-sm mb-2" :class="theme.text">回调地址</div>
                    <UInput v-model="frontendConfig.githubCallbackURL" placeholder="例如 https://your.domain.com/oauth/github/callback" />
                  </div>
                </div>
                <div class="flex justify-end gap-2">
                  <UButton variant="soft" color="gray" @click="fetchConfig">刷新</UButton>
                  <UButton color="green" @click="saveGithubOAuthConfig">保存</UButton>
                  <UButton color="primary" @click="testGithubOAuth">测试</UButton>
                </div>
                <div class="text-xs" :class="theme.mutedText">默认不开启，开启后登录页显示“GitHub 一键登录”按钮</div>
              </div>
            </div>
          </div>

          <div id="db-section" class="col-span-12">
            <div class="rounded-xl border shadow-xl" :class="[theme.cardBg, theme.border]">
              <div class="flex items-center justify-between px-4 py-3">
                <div class="font-semibold flex items-center gap-2" :class="theme.text">
                  <UIcon name="i-heroicons-database" class="w-5 h-5" />
                  <span>数据库管理</span>
                </div>
              </div>
              <div class="px-4 pb-4 space-y-4">
                <div class="flex gap-4 flex-wrap">
                  <UButton color="primary" icon="i-heroicons-arrow-down-tray" class="shadow ring-1 ring-inset ring-primary-400 text-white" @click="downloadBackup">下载备份</UButton>
                  <UButton color="warning" variant="solid" icon="i-heroicons-arrow-up-tray" class="shadow ring-1 ring-inset ring-warning-400 text-white transition-colors duration-200 hover:opacity-90" @click="triggerDatabaseUpload">恢复数据库</UButton>
                </div>
                <div class="text-yellow-400 text-sm max-h-16 overflow-y-auto rounded p-2" :class="theme.subtleBg">🔔：SQLite一键备份恢复，云端数据库请在服务端操作</div>
                <input type="file" ref="databaseFileInput" accept=".zip" class="hidden" @change="handleDatabaseUpload" />

                <div class="rounded-lg p-3" :class="theme.subtleBg">
                  <div class="font-semibold mb-2" :class="theme.text">云存储接入（R2/S3）</div>
                  <div class="grid grid-cols-1 md:grid-cols-2 gap-3">
                    <div class="flex items-center gap-2">
                      <span class="text-sm" :class="theme.mutedText">启用</span>
                      <USwitch v-model="storageEnabled" />
                    </div>
                    <div>
                      <label class="text-sm mb-1 block" :class="theme.mutedText">提供方</label>
                      <USelect v-model="storageConfig.provider" :options="[{label:'S3',value:'s3'},{label:'R2',value:'r2'}]" />
                    </div>
                    <div>
                      <label class="text-sm mb-1 block" :class="theme.mutedText">Endpoint</label>
                      <UInput v-model="storageConfig.endpoint" placeholder="https://..." />
                    </div>
                    <div>
                      <label class="text-sm mb-1 block" :class="theme.mutedText">Region</label>
                      <UInput v-model="storageConfig.region" placeholder="auto 或区域名" />
                    </div>
                    <div>
                      <label class="text-sm mb-1 block" :class="theme.mutedText">Bucket</label>
                      <UInput v-model="storageConfig.bucket" placeholder="bucket 名称" />
                    </div>
                    <div>
                      <label class="text-sm mb-1 block" :class="theme.mutedText">Access Key</label>
                      <UInput v-model="storageConfig.accessKey" />
                    </div>
                    <div>
                      <label class="text-sm mb-1 block" :class="theme.mutedText">Secret Key</label>
                      <UInput v-model="storageConfig.secretKey" type="password" />
                    </div>
                    <div class="flex items-center gap-2">
                      <span class="text-sm" :class="theme.mutedText">PathStyle</span>
                      <USwitch v-model="storageConfig.usePathStyle" />
                    </div>
                    <div class="md:col-span-2">
                      <label class="text-sm mb-1 block" :class="theme.mutedText">公共访问前缀</label>
                      <UInput v-model="storageConfig.publicBaseURL" placeholder="https://bucket.example.com/" />
                    </div>
                  </div>
                  <div class="flex justify-end gap-2 mt-2">
                    <UButton variant="soft" color="gray" @click="loadStorageConfig">刷新</UButton>
                    <UButton color="green" @click="saveStorageConfig">保存</UButton>
                  </div>
                  <div class="grid grid-cols-1 md:grid-cols-2 gap-3 mt-3">
                    <div>
                      <label class="text-sm mb-1 block" :class="theme.mutedText">上传URL（预签名）</label>
                      <UInput v-model="uploadURL" placeholder="粘贴R2/S3预签名上传URL" />
                    </div>
                    <div>
                      <label class="text-sm mb-1 block" :class="theme.mutedText">下载URL（预签名）</label>
                      <UInput v-model="downloadURL" placeholder="粘贴R2/S3预签名下载URL" />
                    </div>
                  </div>
                  <div class="flex justify-end gap-2 mt-2">
                    <UButton color="primary" @click="uploadCloudBackup">上传备份到云</UButton>
                    <UButton color="warning" @click="restoreCloudBackup">从云恢复备份</UButton>
                  </div>
                </div>
              </div>
            </div>
          </div>

        </div>
        
      </main>
      <div class="hidden md:flex fixed bottom-0 left-0 right-0 md:left-72 z-40 border-t px-3 py-3 justify-between items-center" :class="[theme.bottomBg, theme.border]">
        <UButton
          icon="i-heroicons-arrow-left"
          :color="panelTheme === 'light' ? 'gray' : 'white'"
          variant="soft"
          @click="$router.push('/')"
        >
          返回首页
        </UButton>
        <div v-if="isLogin">
          <UButton
            icon="i-heroicons-power"
            color="red"
            variant="solid"
            @click="handleLogout"
          >
            退出登录
          </UButton>
        </div>
        <div v-else class="flex gap-2">
          <UButton color="primary" @click="showLoginModal = true; authmode = true">登录</UButton>
          <UButton color="secondary" @click="showLoginModal = true; authmode = false">注册</UButton>
        </div>
      </div>
      <div v-if="false">
            
                    <div id="section-version" class="mb-6">
                        <div class="text-center mb-6 flex items-center justify-center gap-2">
                            <span class="text-gray-300">当前版本: latest</span>
                            <UButton
                                size="xs"
                                color="gray"
                                variant="ghost"
                                :loading="versionInfo.checking"
                                @click="checkVersion"
                            >
                                {{ versionInfo.checking ? '检测中...' : '检查版本发布时间' }}
                            </UButton>
                        </div>
                        <div v-if="versionInfo.hasUpdate" class="text-center mb-6">
                            <div class="flex items-center justify-center gap-2 text-orange-400">
                                <UIcon name="i-heroicons-arrow-up-circle" class="w-5 h-5" />
                                <span>发现版本最近更新（于 {{ versionInfo.latestVersion }}）</span>
                                <a 
                                    href="https://hub.docker.com/r/noise233/echo-noise/tags" 
                                    target="_blank"
                                    class="text-blue-400 hover:text-blue-300 ml-2"
                                >
                                    查看详情
                                </a>
                            </div>
                        </div>
                    </div>
                    <div id="section-status" class="bg-gray-700 rounded-lg p-4 mb-6">
                        <h2 class="text-xl font-semibold text-white mb-4">系统状态</h2>
                        <div class="grid gap-4">
                            <div class="flex justify-between items-center">
                                <span class="text-gray-300">系统管理员</span>
                                <span class="text-white font-medium">{{ userStore?.status?.username }}</span>
                            </div>
                            <div class="flex justify-between items-center">
                                <span class="text-gray-300">当前用户</span>
                                <span class="text-white font-medium">
                                    {{ isLogin ? userStore.user?.username : "未登录" }}
                                </span>
                            </div>
                            <div class="flex justify-between items-center">
                                <span class="text-gray-300">笔记总数</span>
                                <span class="text-white font-medium">{{ userStore?.status?.total_messages }} 条</span>
                            </div>
                        </div>
                    </div>
                <!-- 添加版本信息和检测按钮 -->
                
                  <!-- 用户信息配置面板 -->
 
                <div v-if="isLogin" class="rounded-lg p-4 mb-6" :class="theme.cardBg">
                    <h2 class="text-xl font-semibold mb-4" :class="theme.text">用户信息配置</h2>
 
                <div id="section-user" v-if="isLogin" class="bg-gray-700 rounded-lg p-4 mb-6">
                    <h2 class="text-xl font-semibold text-white mb-4">用户信息配置</h2>
 
                    <div class="space-y-4">
                        <!-- 用户名修改 -->
                        <div class="bg-gray-800 rounded p-3">
                            <div class="flex justify-between items-center mb-2">
                                <span class="text-gray-300">用户名</span>
                                <UButton
                                    size="sm"
                                    @click="editUserInfo.username = !editUserInfo.username"
                                    :color="editUserInfo.username ? 'gray' : 'green'"
                                    :variant="editUserInfo.username ? 'soft' : 'solid'"
                                >
                                    {{ editUserInfo.username ? '取消' : '编辑' }}
                                </UButton>
                            </div>
                            <div v-if="editUserInfo.username">
                                <UInput
                                    v-model="userForm.username"
                                    placeholder="新用户名"
                                    class="w-full mb-2"
                                />
                                <div class="flex justify-end gap-2">
                                    <UButton @click="updateUsername" color="primary">
                                        保存
                                    </UButton>
                                </div>
                            </div>
                            <div v-else>
                                <p :class="theme.text">{{ userStore.user?.username }}</p>
                            </div>
                        </div>
                         <!-- 在用户信息配置面板中添加 -->
<div class="bg-gray-800 rounded p-3">
    <div class="flex justify-between items-center mb-2">
        <span class="text-gray-300">API Token</span>
        <UButton
            size="sm"
            @click="regenerateToken"
            color="primary"
            variant="soft"
        >
            重新生成
        </UButton>
    </div>
    <div v-if="userToken" class="mb-2">
        <div class="flex items-center gap-2">
            <UInput
                v-model="userToken"
                readonly
                class="font-mono text-sm"
            />
            <UButton
                icon="i-heroicons-clipboard"
                color="gray"
                variant="ghost"
                @click="copyToken"
            />
        </div>
        <p class="text-xs text-gray-400 mt-1">请妥善保管此 Token，它用于 API 访问认证</p>
    </div>
    <div v-else>
        <p class="text-gray-400">暂无 Token</p>
    </div>
</div>
                        <!-- 密码修改 -->
                        <div class="bg-gray-800 rounded p-3 mt-4">
                            <div class="flex justify-between items-center mb-2">
                                <span class="text-gray-300">修改密码</span>
                                <UButton
                                    size="sm"
                                    @click="editUserInfo.password = !editUserInfo.password"
                                    :color="editUserInfo.password ? 'gray' : 'green'"
                                    :variant="editUserInfo.password ? 'soft' : 'solid'"
                                >
                                    {{ editUserInfo.password ? '取消' : '编辑' }}
                                </UButton>
                            </div>
                            <div v-if="editUserInfo.password">
                                <UInput
                                    v-model="userForm.oldPassword"
                                    type="password"
                                    placeholder="当前密码"
                                    class="w-full mb-2"
                                />
                                <UInput
                                    v-model="userForm.newPassword"
                                    type="password"
                                    placeholder="新密码"
                                    class="w-full mb-2"
                                />
                                <div class="flex justify-end gap-2">
                                    <UButton @click="updatePassword" color="primary">
                                        保存
                                    </UButton>
                                </div>
                            </div>
                        </div>

                        

                    </div>
                </div>
                               <!-- 网站配置区域 -->
                <div id="section-site" v-if="isAdmin" class="bg-gray-700 rounded-lg p-4 mb-6">
                    <div class="flex justify-between items-center mb-4">
                        <h2 class="text-xl font-semibold" :class="theme.text">网站配置</h2>
                    </div>
                    <!-- 新用户注册开关 -->
                    <div class="flex items-center bg-gray-800 rounded p-3 mb-4 justify-between">
                        <span :class="theme.text">新用户注册</span>
                        <div class="flex items-center gap-4">
                            <div class="flex items-center">
                                <URadio v-model="registerEnabled" :value="true" class="mr-2" />
                                <span :class="registerEnabled ? theme.text : 'text-gray-400'">允许</span>
                            </div>
                            <div class="flex items-center">
                                <URadio v-model="registerEnabled" :value="false" class="mr-2" />
                                <span :class="!registerEnabled ? theme.text : 'text-gray-400'">不允许</span>
                            </div>
                            <UButton color="green" @click="saveRegisterConfig">保存</UButton>
                        </div>
                    </div>

                    <!-- PWA 配置区域 -->
                    <div class="bg-gray-800 rounded p-4 mb-4">
                        <div class="flex justify-between items-center mb-3">
                            <span :class="theme.text">PWA 模式</span>
                        <div class="flex items-center gap-4">
                                <UToggle v-model="frontendConfig.pwaEnabled" />
                                <UButton color="green" @click="savePWAConfig">保存</UButton>
                            </div>
                        </div>
                        <div class="grid grid-cols-1 md:grid-cols-2 gap-3">
                            <div>
                                <label :class="[theme.mutedText, 'text-sm mb-1 block']">PWA 标题</label>
                                <UInput v-model="frontendConfig.pwaTitle" :placeholder="frontendConfig.siteTitle || '说说笔记'" />
                            </div>
                            <div>
                                <label :class="[theme.mutedText, 'text-sm mb-1 block']">PWA 图标</label>
                                <UInput v-model="frontendConfig.pwaIconURL" :placeholder="'/favicon.ico'" />
                            </div>
                            <div class="md:col-span-2">
                                <label :class="[theme.mutedText, 'text-sm mb-1 block']">PWA 描述</label>
                      <UTextarea v-model="frontendConfig.pwaDescription" :placeholder="frontendConfig.description || ''" />
                    </div>
                  </div>
                </div>
                    </div>

                    <!-- GitHub 链接卡片解析（独立设置） -->
                    <div class="flex items-center bg-gray-800 rounded p-3 mb-4 justify-between">
                        <span :class="theme.text">GitHub 链接卡片解析</span>
                        <div class="flex items-center gap-4">
                            <div class="flex items-center">
                                <URadio v-model="githubCardEnabled" :value="true" class="mr-2" />
                                <span :class="githubCardEnabled ? theme.text : 'text-gray-400'">开启</span>
                            </div>
                            <div class="flex items-center">
                                <URadio v-model="githubCardEnabled" :value="false" class="mr-2" />
                                <span :class="!githubCardEnabled ? theme.text : 'text-gray-400'">关闭</span>
                            </div>
                            <UButton color="green" @click="saveGithubCardConfig">保存</UButton>
                        </div>
                    </div>
                    <!-- 公告栏开关（独立设置） -->
                    <div class="flex items-center bg-gray-800 rounded p-3 mb-4 justify-between">
                        <span :class="theme.text">公告栏开关</span>
                        <div class="flex items-center gap-4">
                            <UToggle v-model="frontendConfig.announcementEnabled" />
                            <UButton color="green" @click="saveConfigItem('announcementEnabled')">保存</UButton>
                        </div>
                    </div>

                <!-- 默认主题色设置 -->
                <div class="flex items-center bg-gray-800 rounded p-3 mb-4 justify-between">
                    <span :class="theme.text">默认主题色</span>
                    <div class="flex items-center gap-4">
                        <div class="flex items-center">
                            <URadio v-model="frontendConfig.defaultContentTheme" value="dark" class="mr-2" />
                            <span :class="frontendConfig.defaultContentTheme === 'dark' ? theme.text : 'text-gray-400'">暗黑</span>
                        </div>
                        <div class="flex items-center">
                            <URadio v-model="frontendConfig.defaultContentTheme" value="light" class="mr-2" />
                            <span :class="frontendConfig.defaultContentTheme === 'light' ? theme.text : 'text-gray-400'">白天</span>
                        </div>
                        <UButton color="green" @click="saveConfigItem('defaultContentTheme')">保存</UButton>
                    </div>
                </div>

                <!-- 音乐配置 -->
                <div id="site-music-section" class="rounded p-4 mb-4" :class="[theme.cardBg]">
                  <div class="flex justify-between items-center mb-3">
                    <span :class="theme.text">音乐播放器</span>
                    <div class="flex items-center gap-2">
                      <UButton :color="frontendConfig.musicEnabled ? 'gray' : 'green'" variant="soft" @click="toggleMusic(true)">开启</UButton>
                      <UButton color="red" variant="soft" @click="toggleMusic(false)">关闭</UButton>
                      <UToggle v-model="frontendConfig.musicEnabled" />
                      <UButton color="green" @click="saveMusicConfig">保存</UButton>
                      <UButton variant="soft" color="gray" @click="resetMusicConfig">重置</UButton>
                    </div>
                  </div>
                    <div class="grid grid-cols-1 md:grid-cols-2 gap-3">
                    <div>
                      <label :class="[theme.mutedText, 'text-sm mb-1 block']">歌单 ID</label>
                      <UInput v-model="frontendConfig.musicPlaylistId" placeholder="如 14273792576" />
                    </div>
                    <div>
                      <label :class="[theme.mutedText, 'text-sm mb-1 block']">单曲 ID</label>
                      <UInput v-model="frontendConfig.musicSongId" placeholder="可选，优先歌单" />
                    </div>
                    <div>
                      <label :class="[theme.mutedText, 'text-sm mb-1 block']">显示位置</label>
                      <USelect v-model="frontendConfig.musicPosition" :options="[
                        {label:'左下角',value:'bottom-left'},
                        {label:'右下角',value:'bottom-right'},
                        {label:'左上角',value:'top-left'},
                        {label:'右上角',value:'top-right'}
                      ]" />
                    </div>
                    <div>
                      <label :class="[theme.mutedText, 'text-sm mb-1 block']">主题风格</label>
                      <USelect v-model="frontendConfig.musicTheme" :options="[
                        {label:'自动',value:'auto'},
                        {label:'浅色',value:'light'},
                        {label:'深色',value:'dark'}
                      ]" />
                    </div>
                    <div class="flex items-center gap-3">
                      <span class="text-sm" :class="theme.mutedText">显示歌词</span>
                      <USwitch v-model="frontendConfig.musicLyric" />
                    </div>
                    <div class="flex items-center gap-3">
                      <span class="text-sm" :class="theme.mutedText">自动播放</span>
                      <USwitch v-model="frontendConfig.musicAutoplay" />
                    </div>
                    <div class="flex items-center gap-3">
                      <span class="text-sm" :class="theme.mutedText">默认最小化</span>
                      <USwitch v-model="frontendConfig.musicDefaultMinimized" />
                    </div>
                    <div class="flex items-center gap-3">
                      <span class="text-sm" :class="theme.mutedText">展示模式</span>
                      <USelect v-model="musicEmbedMode" :options="[{label:'嵌入',value:'embed'},{label:'浮动',value:'float'}]" />
                    </div>
                  </div>
                  <div class="text-xs mt-2" :class="theme.mutedText">保存后首页自动刷新显示播放器；歌单与单曲任选其一</div>
                </div>

                    <!-- 配置展示/编辑表单 -->
                    <div class="space-y-4">
                        <div v-for="(label, key) in configLabels" :key="key" class="bg-gray-800 rounded p-3">
                            <div class="flex justify-between items-center mb-2">
                                <span :class="theme.mutedText">{{ label }}</span>
                                <UButton
                                    size="sm"
                                    @click="editItem[key] = !editItem[key]"
                                    :color="editItem[key] ? 'gray' : 'green'"
                                    :variant="editItem[key] ? 'soft' : 'solid'"
                                >
                                    {{ editItem[key] ? '取消' : '编辑' }}
                                </UButton>
                            </div>
                            
                            <div v-if="editItem[key]">
                        <template v-if="key === 'backgrounds'">
                            <div class="space-y-2">
                                <div v-for="(bg, index) in frontendConfig.backgrounds" 
                                     :key="index" 
                                     class="flex gap-2">
                                    <UInput v-model="frontendConfig.backgrounds[index]" 
                                           placeholder="输入图片URL" 
                                           class="flex-1" />
                                    <UButton @click="removeBackground(index)" 
                                            icon="i-heroicons-trash" 
                                            color="red" 
                                            variant="ghost" />
                                </div>
                                <div class="flex gap-2">
                                    <UButton @click="addBackground" 
                                            icon="i-heroicons-plus" 
                                            variant="ghost" 
                                            class="mr-2">
                                        添加链接
                                            </UButton>
                                            <UButton @click="triggerFileInput" 
                                                    icon="i-heroicons-cloud-arrow-up" 
                                                    variant="ghost">
                                                上传图片
                                            </UButton>
                                        </div>
                                    </div>
                                </template>
                                <template v-else-if="key === 'subtitleText'">
                                    <UTextarea
                                        v-model="frontendConfig[key]"
                                        :placeholder="`输入${label}`"
                                        class="w-full mb-2"
                                    />
                                </template>
                                <template v-else>
                                    <UInput
                                        v-model="frontendConfig[key]"
                                        :placeholder="`输入${label}`"
                                        class="w-full mb-2"
                                    />
                                </template>
                                <div class="flex justify-end gap-2">
                                    <UButton @click="resetConfigItem(key)" variant="ghost" color="gray">
                                        重置
                                    </UButton>
                                    <UButton @click="saveConfigItem(key)" color="primary">
                                        保存
                                    </UButton>
                                </div>
                            </div>
                            <div v-else>
                        <template v-if="key === 'backgrounds'">
                            <div class="grid grid-cols-3 gap-2">
                                <img v-for="(bg, index) in frontendConfig.backgrounds"
                                     :key="index"
                                     :src="bg"
                                     class="w-full h-20 object-cover rounded cursor-pointer"
                                     @click="previewImage(bg)" />
                            </div>
                        </template>
                        <template v-else>
                        <p :class="[theme.text, 'break-words']">{{ frontendConfig[key] }}</p>
                        </template>
                    </div>
                        </div>
                    </div>
                </div>

                <!-- 保存按钮 -->
                <div v-if="editMode" class="flex justify-end gap-2 mb-6">
                    <UButton @click="resetConfig" variant="ghost" color="gray">
                        重置
                    </UButton>
                    <UButton @click="saveConfig" color="primary">
                        保存所有更改
                    </UButton>
                </div>
 
                
                
<!-- 底部操作栏 -->
                
            </div>
 
<!-- 数据库管理面板 -->
<div id="section-db" v-if="isAdmin" class="bg-gray-700 rounded-lg p-4 mb-6">
    <h2 class="text-xl font-semibold text-white mb-4">数据库管理</h2>
    <div class="space-y-4">
        <div class="flex gap-4">
            <UButton
                color="primary"
                icon="i-heroicons-arrow-down-tray"
                @click="downloadBackup"
            >
                下载备份
            </UButton>
            <UButton
                color="warning"
                variant="solid"
                icon="i-heroicons-arrow-up-tray"
                @click="triggerDatabaseUpload"
            >
                恢复数据库
            </UButton>
        </div>
        <div class="text-yellow-400 text-sm max-h-16 overflow-y-auto bg-gray-800/50 rounded p-2">
            🔔：SQLite一键备份恢复，因兼容问题，不支持云端的PostgreSQL/MySQL数据库，如有使用云端数据库，请前往云服务端来备份和恢复
        </div>
        <input
            type="file"
            ref="databaseFileInput"
            accept=".zip"
            class="hidden"
            @change="handleDatabaseUpload"
        />
    </div>
                </div>
            
 
        
        

      </div>
    <!-- 登录模态框 -->
    <UModal v-model="showLoginModal">
        <div class="bg-gray-800 p-6 rounded-lg">
            <h3 class="text-xl font-semibold mb-4" :class="theme.text">
                {{ authmode ? '用户登录' : '用户注册' }}
            </h3>
                <UForm :state="authForm" class="space-y-4">
                    <UFormGroup>
                        <UInput
                            v-model="authForm.username"
                            placeholder="用户名"
                            class="w-full"
                        />
                    </UFormGroup>
                    <UFormGroup>
                        <UInput
                            v-model="authForm.password"
                            type="password"
                            placeholder="密码"
                            class="w-full"
                        />
                    </UFormGroup>
                    <div class="flex justify-between items-center">
                        <UButton
                            variant="ghost"
                            @click="authmode = !authmode"
                        >
                            {{ authmode ? '去注册' : '去登录' }}
                        </UButton>
                        <UButton
                            color="primary"
                            @click="authmode ? login(authForm) : register(authForm)"
                        >
                            {{ authmode ? '登录' : '注册' }}
                        </UButton>
                    </div>
                </UForm>
            </div> 
        </UModal>
        <input
            type="file"
            ref="fileInput"
            accept="image/*"
            multiple
            class="hidden"
            @change="handleFileUpload"
        />
  
    </div>
</template>

<script setup lang="ts">
import { ref, reactive, watch, computed, onMounted, nextTick } from 'vue'
import type { UserToLogin, UserToRegister } from '~/types/models'
import { useUser } from '~/composables/useUser'
import { useUserStore } from '~/store/user'
import { useToast } from '#ui/composables/useToast'
import NotifyPanel from './NotifyPanel.vue'
 
import CommentsSettings from '~/components/admin/CommentsSettings.vue'
import { getRequest, putRequest, postRequest, deleteRequest } from '~/utils/api'
 
const scrollTo = (id: string) => {
    const el = document.getElementById(id)
    if (el) el.scrollIntoView({ behavior: 'smooth', block: 'start' })
}
 

// 新用户注册开关相关
const registerEnabled = ref(true);
const sidebarOpen = ref(true)
const panelTheme = ref<'dark' | 'midnight' | 'slate' | 'light'>(
  (typeof window !== 'undefined' && (localStorage.getItem('adminTheme') as any)) || 'dark'
)
const baseApi = useRuntimeConfig().public.baseApi || '/api'
const avatarSrc = computed(() => {
  const siteAvatar = (frontendConfig.avatarURL || '').trim()
  const username = (userStore.user as any)?.username || (userStore.user as any)?.Username || frontendConfig.username
  if (siteAvatar) {
    if (siteAvatar.startsWith('http')) return siteAvatar
    if (siteAvatar.startsWith('/api')) return `${baseApi}${siteAvatar}`
    return `${baseApi}${siteAvatar}`
  }
  if (username) {
    return `https://i.pravatar.cc/100?u=${encodeURIComponent(username)}`
  }
  return '/favicon.ico'
})

const setActive = async (name: 'system' | 'user' | 'site' | 'notify' | 'db' | 'site-register' | 'site-pwa' | 'site-github-card' | 'site-github-login' | 'site-announcement' | 'site-music' | 'site-default-theme' | 'site-configs' | 'comments' | 'email') => {
  await nextTick()
  const el = document.getElementById(`${name}-section`)
  el?.scrollIntoView({ behavior: 'smooth', block: 'start' })
  if (window.innerWidth < 768) sidebarOpen.value = false
}

onMounted(() => {
  if (typeof window !== 'undefined' && window.innerWidth < 768) sidebarOpen.value = false
  loadStorageConfig()
})

const theme = computed(() => {
  if (panelTheme.value === 'light') {
    return {
      sidebarBg: 'bg-white/90',
      headerBg: 'bg-white/80',
      bottomBg: 'bg-white/80',
      cardBg: 'bg-white/80',
      subtleBg: 'bg-gray-100',
      border: 'border-gray-200',
      text: 'text-slate-900',
      sidebarText: 'text-slate-900',
      mutedText: 'text-slate-700',
      pageBg: 'bg-gray-50',
      navBtnBg: 'bg-gray-100',
      navBtnHoverBg: 'hover:bg-gray-200'
    }
  }
  if (panelTheme.value === 'dark') {
    return {
      sidebarBg: 'bg-gray-950/80',
      headerBg: 'bg-gray-950/70',
      bottomBg: 'bg-gray-950/70',
      cardBg: 'bg-gray-900/70',
      subtleBg: 'bg-gray-900/50',
      border: 'border-gray-800/60',
      text: 'text-white',
      sidebarText: 'text-white',
      mutedText: 'text-gray-300',
      pageBg: 'bg-gradient-to-br from-slate-900 via-indigo-950 to-slate-800',
      navBtnBg: 'bg-slate-800/70',
      navBtnHoverBg: 'hover:bg-slate-700/70'
    }
  }
  if (panelTheme.value === 'midnight') {
    return {
      sidebarBg: 'bg-indigo-950/60',
      headerBg: 'bg-indigo-950/50',
      bottomBg: 'bg-indigo-950/50',
      cardBg: 'bg-indigo-900/40',
      subtleBg: 'bg-indigo-950/30',
      border: 'border-indigo-800/40',
      text: 'text-white',
      sidebarText: 'text-white',
      mutedText: 'text-indigo-200',
      pageBg: 'bg-gradient-to-br from-indigo-950 via-indigo-900 to-slate-900',
      navBtnBg: 'bg-indigo-950/40',
      navBtnHoverBg: 'hover:bg-indigo-900/40'
    }
  }
  if (panelTheme.value === 'slate') {
    return {
      sidebarBg: 'bg-slate-900/70',
      headerBg: 'bg-slate-900/60',
      bottomBg: 'bg-slate-900/60',
      cardBg: 'bg-slate-800/70',
      subtleBg: 'bg-slate-900/40',
      border: 'border-slate-700/40',
      text: 'text-white',
      sidebarText: 'text-white',
      mutedText: 'text-slate-200',
      pageBg: 'bg-gradient-to-br from-slate-900 via-slate-800 to-slate-900',
      navBtnBg: 'bg-slate-800/70',
      navBtnHoverBg: 'hover:bg-slate-700/70'
    }
  }
  return {
    sidebarBg: 'bg-slate-900/70',
    headerBg: 'bg-slate-900/60',
    bottomBg: 'bg-slate-900/60',
    cardBg: 'bg-slate-800/70',
    subtleBg: 'bg-slate-900/40',
      border: 'border-slate-700/40',
    text: 'text-white',
    sidebarText: 'text-white',
    mutedText: 'text-slate-200',
    pageBg: 'bg-gradient-to-br from-slate-900 via-indigo-950 to-slate-800',
    navBtnBg: 'bg-slate-800/70',
    navBtnHoverBg: 'hover:bg-slate-700/70'
  }
})

const saveAdminTheme = async () => {
  localStorage.setItem('adminTheme', panelTheme.value)
  try {
    const resConfig = await fetch('/api/frontend/config', { credentials: 'include' })
    const dataConfig = await resConfig.json()
    let payload: any = {}
    if (dataConfig.code === 1) {
      payload = { ...dataConfig.data, adminTheme: panelTheme.value }
    } else {
      payload = { adminTheme: panelTheme.value }
    }
    const res = await fetch('/api/settings', {
      method: 'PUT',
      headers: { 'Content-Type': 'application/json' },
      credentials: 'include',
      body: JSON.stringify(payload)
    })
    const data = await res.json()
    if (data.code === 1) {
      useToast().add({ title: res?.msg || '已保存', color: 'green' })
    }
  } catch {}
}

// 页面加载时获取配置
const fetchRegisterConfig = async () => {
    try {
        const res = await fetch('/api/frontend/config', { credentials: 'include' });
        const data = await res.json();
        if (data.code === 1 && typeof data.data.allowRegistration === 'boolean') {
            registerEnabled.value = data.data.allowRegistration;
        }
    } catch (e) {
        useToast().add({ title: '获取注册配置失败', color: 'red' });
    }
};
onMounted(fetchRegisterConfig);

// 保存配置
const saveRegisterConfig = async () => {
    try {
        // 先获取完整配置
        const resConfig = await fetch('/api/frontend/config', { credentials: 'include' });
        const dataConfig = await resConfig.json();
        let payload = {};
        if (dataConfig.code === 1) {
            payload = {
                ...dataConfig.data,
                allowRegistration: registerEnabled.value
            };
        } else {
            // 如果获取失败，只发 allowRegistration（兼容旧接口）
            payload = { allowRegistration: registerEnabled.value };
        }

        const res = await fetch('/api/settings', {
            method: 'PUT',
            headers: { 'Content-Type': 'application/json' },
            credentials: 'include',
            body: JSON.stringify(payload)
        });
        const data = await res.json();
        if (data.code === 1) {
            useToast().add({ title: '保存成功', color: 'green' });
        } else {
            throw new Error(data.msg || '保存失败');
        }
    } catch (e) {
        useToast().add({ title: '保存失败', color: 'red' });
    }
};

const userStore = useUserStore()
const { login, register, logout } = useUser()
const router = useRouter()
const userToken = ref('')
const versionInfo = reactive({
    checking: false,
    hasUpdate: false,
    latestVersion: ''
})
// 推送配置
let notifyConfig = reactive({
    webhookEnabled: false,
    webhookURL: '',
    telegramEnabled: false,
    telegramToken: '',
    telegramChatID: '',
    weworkEnabled: false,
    weworkKey: '',
    feishuEnabled: false,
    feishuWebhook: '',
    feishuSecret: ''
})

// 获取推送配置
const fetchNotifyConfig = async () => {
    try {
        const response = await fetch('/api/notify/config', {
            credentials: 'include'
        })
        const data = await response.json()
        if (data.code === 1) {
            Object.assign(notifyConfig, data.data)
        }
    } catch (error) {
        console.error('获取推送配置失败:', error)
    }
}
const smtp = reactive({ enabled: false, driver: 'smtp', host: '', port: '', user: '', pass: '', from: '', tls: false, encryption: 'tls' })
const showSmtpPass = ref(false)
const loadSmtp = async () => {
  try {
    const res = await getRequest<any>('frontend/config', undefined, { credentials: 'include' })
    if (res && res.code === 1) {
      const cfg = res.data || {}
      smtp.enabled = !!cfg.smtpEnabled
      smtp.driver = cfg.smtpDriver || 'smtp'
      smtp.host = cfg.smtpHost || ''
      smtp.port = (cfg.smtpPort ?? '').toString()
      smtp.user = cfg.smtpUser || ''
      smtp.pass = cfg.smtpPass || ''
      smtp.from = cfg.smtpFrom || ''
      smtp.tls = !!cfg.smtpTLS
      smtp.encryption = (cfg.smtpEncryption || (smtp.tls ? 'tls' : 'ssl'))
    }
  } catch {}
}
onMounted(loadSmtp)
const saveSmtp = async () => {
  try {
    const resCfg = await getRequest<any>('frontend/config', undefined, { credentials: 'include' })
    const payload: any = resCfg?.code === 1 ? { ...resCfg.data } : {}
    payload.smtpEnabled = smtp.enabled
    payload.smtpDriver = smtp.driver
    payload.smtpHost = smtp.host
    payload.smtpPort = parseInt(smtp.port || '0') || 0
    payload.smtpUser = smtp.user
    payload.smtpPass = smtp.pass
    payload.smtpFrom = smtp.from
    payload.smtpEncryption = smtp.encryption
    payload.smtpTLS = smtp.encryption === 'tls'
    const res = await putRequest<any>('settings', payload, { credentials: 'include' })
    if (res && res.code === 1) {
      useToast().add({ title: res?.msg || '已保存', color: 'green' })
    } else {
      throw new Error(res?.msg || '保存失败')
    }
  } catch (e: any) {
    useToast().add({ title: '保存失败', description: e.message, color: 'red' })
  }
}

const adminUsers = ref<string[]>([])
const newAdmin = ref('')
const adminPasswordMasked = ref('')
const showAdminPassword = ref(false)
const showAdminResetModal = ref(false)
const adminReset = reactive({ newPass: '', confirmPass: '' })
const adminResetStrength = computed(() => {
  const v = adminReset.newPass || ''
  let score = 0
  if (v.length >= 8) score++
  if (/[A-Z]/.test(v) && /[a-z]/.test(v)) score++
  if (/\d/.test(v) && /[^A-Za-z0-9]/.test(v)) score++
  return Math.min(score, 3)
})
const adminResetStrengthLabel = computed(() => {
  if (adminResetStrength.value <= 1) return '弱'
  if (adminResetStrength.value === 2) return '中'
  return '强'
})
const adminResetStrengthColor = computed(() => {
  if (adminResetStrength.value <= 1) return 'red'
  if (adminResetStrength.value === 2) return 'orange'
  return 'green'
})
const canSaveAdminReset = computed(() => {
  if (!adminReset.newPass || !adminReset.confirmPass) return false
  if (adminReset.newPass !== adminReset.confirmPass) return false
  return adminResetStrength.value >= 2
})
const loadAdmins = async () => {
  try {
    const res = await getRequest<any>('frontend/config', undefined, { credentials: 'include' })
    if (res && res.code === 1) {
      const cfg = res.data || {}
      adminUsers.value = Array.isArray(cfg.adminUsers) ? cfg.adminUsers : []
      adminPasswordMasked.value = cfg.adminPasswordMasked || ''
    }
  } catch {}
}
onMounted(loadAdmins)
const addAdmin = () => {
  const name = (newAdmin.value || '').trim()
  if (!name) return
  if (!adminUsers.value.includes(name)) adminUsers.value.push(name)
  newAdmin.value = ''
}
const removeAdmin = (name: string) => {
  adminUsers.value = adminUsers.value.filter(n => n !== name)
}
const saveAdmins = async () => {
  try {
    const resCfg = await getRequest<any>('frontend/config', undefined, { credentials: 'include' })
    const payload: any = resCfg?.code === 1 ? { ...resCfg.data } : {}
    payload.adminUsers = [...adminUsers.value]
    const res = await putRequest<any>('settings', payload, { credentials: 'include' })
    if (res && res.code === 1) {
      useToast().add({ title: res?.msg || '已保存', color: 'green' })
      await userStore.getStatus()
    } else {
      throw new Error(res?.msg || '保存失败')
    }
  } catch (e: any) {
    useToast().add({ title: '保存失败', description: e.message, color: 'red' })
  }
}
const resetAdminPassword = async () => {
  try {
    if (!canSaveAdminReset.value) throw new Error('请填写符合强度的新密码并确认一致')
    const resCfg = await getRequest<any>('frontend/config', undefined, { credentials: 'include' })
    const payload: any = resCfg?.code === 1 ? { ...resCfg.data } : {}
    payload.adminPasswordReset = adminReset.newPass
    const res = await putRequest<any>('settings', payload, { credentials: 'include' })
    if (res && res.code === 1) {
      useToast().add({ title: res?.msg || '管理员密码已重置', color: 'green' })
      await loadAdmins()
      showAdminPassword.value = false
      showAdminResetModal.value = false
      adminReset.newPass = ''
      adminReset.confirmPass = ''
    } else {
      throw new Error(res?.msg || '重置失败')
    }
  } catch (e: any) {
    useToast().add({ title: '重置失败', description: e.message, color: 'red' })
  }
}
// 管理员用户列表与搜索
const userSearch = ref('')
const allUsers = computed<any[]>(() => {
  const s: any = userStore.status || {}
  const list = s.users || s.Users || []
  return Array.isArray(list) ? list : []
})
const filteredUsers = computed<any[]>(() => {
  const q = (userSearch.value || '').trim().toLowerCase()
  if (!q) return allUsers.value
  return allUsers.value.filter((u: any) => {
    const id = String(u.id ?? u.ID ?? u.user_id ?? '')
    const name = String(u.username ?? u.Username ?? '').toLowerCase()
    return id.includes(q) || name.includes(q)
  })
})
const refreshUsers = async () => {
  await userStore.getStatus()
}
const showUsers = ref(true)
const confirmToggleAdmin = async (u: any) => {
  try {
    const name = u.username ?? u.Username
    if (!window.confirm(`确定要切换用户“${name}”的管理员权限吗？`)) return
    if (!window.confirm('该操作存在风险，是否继续？')) return
    const id = u.id ?? u.ID ?? u.user_id
    const res = await putRequest<any>(`user/admin?id=${id}`, {}, { credentials: 'include' })
    if (res && res.code === 1) {
      useToast().add({ title: res?.msg || '已更新管理员状态', color: 'green' })
      await userStore.getStatus()
    } else {
      throw new Error(res?.msg || '更新失败')
    }
  } catch (e: any) {
    useToast().add({ title: '更新失败', description: e.message, color: 'red' })
  }
}
const confirmDeleteUser = async (u: any) => {
  try {
    const name = u.username ?? u.Username
    if (!window.confirm(`确定要删除用户“${name}”吗？删除后不可恢复。`)) return
    if (!window.confirm('该操作存在风险，是否继续？')) return
    const id = u.id ?? u.ID ?? u.user_id
    const res = await deleteRequest<any>('user', { id }, { credentials: 'include' })
    if (res && res.code === 1) {
      useToast().add({ title: res?.msg || '已删除用户', color: 'green' })
      await userStore.getStatus()
    } else {
      throw new Error(res?.msg || '删除失败')
    }
  } catch (e: any) {
    useToast().add({ title: '删除失败', description: e.message, color: 'red' })
  }
}
const toggleAdmin = async (u: any) => {
  try {
    const id = u.id ?? u.ID ?? u.user_id
    const res = await putRequest<any>(`user/admin?id=${id}`, {}, { credentials: 'include' })
    if (res && res.code === 1) {
      useToast().add({ title: res?.msg || '已更新管理员状态', color: 'green' })
      await userStore.getStatus()
    } else {
      throw new Error(res?.msg || '更新失败')
    }
  } catch (e: any) {
    useToast().add({ title: '更新失败', description: e.message, color: 'red' })
  }
}
const testingSmtp = ref(false)
const testSmtp = async () => {
  try {
    const to = (smtp.from || smtp.user || '').trim()
    if (!to || !smtp.host || !smtp.port || !smtp.user || !smtp.pass || !smtp.encryption) {
      throw new Error('请完整填写地址、主机、端口、加密协议、用户名和密码')
    }
    testingSmtp.value = true
    if (!smtp.enabled) {
      smtp.enabled = true
      await saveSmtp()
    }
    // 优先使用现有通知测试接口
    let res = await postRequest<any>('notify/test', { type: 'email', to }, { credentials: 'include' })
    if (!res || res.code !== 1) {
      // 回退到专用邮箱测试接口（部分后端可能未提供）
      res = await postRequest<any>('email/test', { to }, { credentials: 'include' })
    }
    if (res && res.code === 1) {
      useToast().add({ title: res?.msg || '测试邮件已发送', color: 'green' })
    } else {
      throw new Error(res?.msg || '发送失败或接口不存在')
    }
  } catch (e: any) {
    useToast().add({ title: '失败', description: e.message, color: 'red' })
  } finally {
    testingSmtp.value = false
  }
}
const testGithubOAuth = () => {
  try {
    if (!frontendConfig.githubOAuthEnabled) throw new Error('请先开启 GitHub 登录')
    if (!frontendConfig.githubClientId || !frontendConfig.githubCallbackURL) throw new Error('请先填写 Client ID 与回调地址')
    const BASE_API = useRuntimeConfig().public.baseApi || '/api'
    window.open(`${BASE_API}/oauth/github/login`, '_blank')
  } catch (e: any) {
    useToast().add({ title: '无法测试', description: e.message, color: 'red' })
  }
}

// 检查版本更新
const checkVersion = async () => {
    versionInfo.checking = true;
    try {
        const response = await fetch('/api/version/check', {
            credentials: 'include',
            headers: {
                'Cache-Control': 'no-cache',
                'Pragma': 'no-cache'
            }
        });
        
        const data = await response.json();
        if (data.code === 1) {
            const { hasUpdate, lastUpdateTime } = data.data;
            versionInfo.hasUpdate = hasUpdate;
            versionInfo.latestVersion = new Date(lastUpdateTime).toLocaleString('zh-CN', {
                year: 'numeric',
                month: '2-digit',
                day: '2-digit',
                hour: '2-digit',
                minute: '2-digit'
            });

            if (hasUpdate) {
                useToast().add({
                    title: '发现版本',
                    description: `最新版本发布于 ${versionInfo.latestVersion}`,
                    color: 'orange'
                });
            } else {
                useToast().add({
                    title: '已是最新版本',
                    description: '当前使用的是最新版本',
                    color: 'green'
                });
            }
        } else {
            throw new Error(data.msg || '检查更新失败');
        }
    } catch (error) {
        console.error('检查版本更新失败:', error);
        useToast().add({
            title: '检查更新失败',
            description: '请科学上网后重试',
            color: 'red'
        });
    } finally {
        versionInfo.checking = false;
    }
};
// 重新生成 Token
// 修改 regenerateToken 函数
const regenerateToken = async () => {
    if (!userStore.isLogin) {
        useToast().add({
            title: '错误',
            description: '请先登录',
            color: 'red'
        });
        return;
    }

    try {
        if (typeof window !== 'undefined') {
            const ok = window.confirm('重新生成将使旧 Token 失效，确认继续？')
            if (!ok) return
        }
        regeneratingToken.value = true
        const response = await fetch('/api/user/token/regenerate', {
            method: 'POST',
            credentials: 'include',
            headers: {
                'Content-Type': 'application/json'
            }
        });

        const data = await response.json();
        if (!response.ok) {
            throw new Error(data.msg || 'Token生成请求失败');
        }

        if (data.code === 1 && data.data?.token) {
            userToken.value = data.data.token;
            showToken.value = false
            useToast().add({
                title: '成功',
                description: data?.msg || 'Token 已更新',
                color: 'green'
            });
        } else {
            throw new Error(data.msg || 'Token 生成失败');
        }
    } catch (error: any) {
        console.error('Token生成错误:', error);
        useToast().add({
            title: '错误',
            description: error.message || 'Token 生成失败',
            color: 'red'
        });
    } finally {
        regeneratingToken.value = false
    }
};

// 复制 Token
const copyToken = async () => {
    try {
        await navigator.clipboard.writeText(userToken.value)
        useToast().add({
            title: '成功',
            description: 'Token 已复制到剪贴板',
            color: 'green'
        })
    } catch (error) {
        useToast().add({
            title: '错误',
            description: '复制失败',
            color: 'red'
        })
    }
}
// 添加退出登录处理函数
const handleLogout = async () => {
    try {
        const response = await fetch('/api/user/logout', {
            method: 'POST',
            credentials: 'include'
        })
        const data = await response.json().catch(() => ({}))
        if (!response.ok || data.code !== 1) {
            throw new Error(data?.msg || '退出失败')
        }
        userStore.clearUserStatus()
        useToast().add({ title: '成功', description: '已退出登录', color: 'green' })
        router.push('/')
    } catch (error) {
        userStore.clearUserStatus()
        useToast().add({ title: '成功', description: '已退出登录', color: 'green' })
        router.push('/')
    }
}
// 状态变量
const isLogin = computed(() => userStore?.isLogin ?? false)
const isAdmin = computed(() => {
    return userStore.user?.is_admin ?? false
})
const authmode = ref(true)
const showLoginModal = ref(false)
const editMode = ref(false)
const fileInput = ref<HTMLInputElement | null>(null) 
    const userForm = reactive({
    username: '',
    oldPassword: '',
    newPassword: '',
    confirmPassword: ''
})
const editUserInfo = reactive({
    username: false,
    password: false
})
const showToken = ref(false)
const regeneratingToken = ref(false)
const showOldPassword = ref(false)
const showNewPassword = ref(false)
const showConfirmPassword = ref(false)
const passwordStrength = computed(() => {
    const v = userForm.newPassword || ''
    let score = 0
    if (v.length >= 8) score++
    if (/[A-Z]/.test(v) && /[a-z]/.test(v)) score++
    if (/\d/.test(v) && /[^A-Za-z0-9]/.test(v)) score++
    return Math.min(score, 3)
})
const passwordStrengthLabel = computed(() => {
    if (passwordStrength.value <= 1) return '弱'
    if (passwordStrength.value === 2) return '中'
    return '强'
})
const passwordStrengthColor = computed(() => {
    if (passwordStrength.value <= 1) return 'red'
    if (passwordStrength.value === 2) return 'orange'
    return 'green'
})
const canSavePassword = computed(() => {
    if (!userForm.oldPassword || !userForm.newPassword || !userForm.confirmPassword) return false
    if (userForm.newPassword === userForm.oldPassword) return false
    if (userForm.newPassword !== userForm.confirmPassword) return false
    return passwordStrength.value >= 2
})
const updateUsername = async () => {
    try {
        if (!userForm.username.trim()) {
            throw new Error('用户名不能为空')
        }
        const res = await putRequest<any>('user/update', { username: userForm.username, type: 'username' }, { credentials: 'include' })
        if (res && res.code === 1) {
            await userStore.getUser()
            editUserInfo.username = false
            userForm.username = ''
            useToast().add({
                title: '成功',
                description: res?.msg || '用户名已更新',
                color: 'green'
            })
        } else {
            throw new Error(res?.msg)
        }
    } catch (error) {
        useToast().add({
            title: '错误',
            description: error.message || '更新失败',
            color: 'red'
        })
    }
}

const updatePassword = async () => {
    try {
        if (!userForm.newPassword || !userForm.oldPassword || !userForm.confirmPassword) {
            throw new Error('密码不能为空')
        }
        if (userForm.newPassword === userForm.oldPassword) {
            throw new Error('新密码不能与当前密码相同')
        }
        if (userForm.newPassword !== userForm.confirmPassword) {
            throw new Error('两次输入不一致')
        }
        if (passwordStrength.value < 2) {
            throw new Error('密码强度不足')
        }
        const res = await putRequest<any>('user/change_password', { password: userForm.newPassword, oldPassword: userForm.oldPassword }, { credentials: 'include' })
        if (res && res.code === 1) {
            editUserInfo.password = false
            userForm.oldPassword = ''
            userForm.newPassword = ''
            userForm.confirmPassword = ''
            useToast().add({
                title: '成功',
                description: res?.msg || '密码已更新',
                color: 'green'
            })
        } else {
            throw new Error(res?.msg)
        }
    } catch (error) {
        useToast().add({
            title: '错误',
            description: error.message || '更新失败',
            color: 'red'
        })
    }
}


// 配置相关
const configLabels = {
    siteTitle: '站点标题',
    subtitleText: '欢迎语',
    avatarURL: '头像链接',
    username: '用户名',
    description: '个人描述',
    backgrounds: '背景图片',
    cardFooterTitle: '卡片页脚标题',
    cardFooterLink: '卡片页脚链接',
    pageFooterHTML: '页面底部HTML',
    rssTitle: 'RSS 标题',
    rssDescription: 'RSS 描述',
    rssAuthorName: 'RSS 作者',
    rssFaviconURL: 'RSS 图标链接',
    walineServerURL: 'Waline 评论服务器地址',
}

const frontendConfig = reactive({
    siteTitle: '',
    subtitleText: '',
    avatarURL: '',
    username: '',
    description: '',
    backgrounds: [] as string[],
    cardFooterTitle: '',
    cardFooterLink: '',
    pageFooterHTML: '',
    rssTitle: '',
    rssDescription: '',
    rssAuthorName: '',
    rssFaviconURL: '',
  walineServerURL: '',
  commentEnabled: false,
  commentSystem: 'waline',
  commentEmailEnabled: false,
  githubOAuthEnabled: false,
  githubClientId: '',
  githubClientSecret: '',
  githubCallbackURL: '',
    enableGithubCard: false,
    // PWA 设置
    pwaEnabled: true,
    pwaTitle: '',
    pwaDescription: '',
    pwaIconURL: '',
    defaultContentTheme: 'dark',
    announcementText: '',
    announcementEnabled: true,
    // 音乐播放器
    musicEnabled: false,
    musicPlaylistId: '',
    musicSongId: '',
    musicPosition: 'bottom-left',
    musicTheme: 'auto',
    musicLyric: true,
    musicAutoplay: false,
    musicDefaultMinimized: true,
    musicEmbed: false,
})

// GitHub 链接卡片解析开关的双向绑定（与 frontendConfig.enableGithubCard 同步）
const githubCardEnabled = computed({
    get: () => frontendConfig.enableGithubCard === true,
    set: (val: any) => {
        const b = (val === true || val === 'true' || val === 1 || val === '1')
        ;(frontendConfig as any).enableGithubCard = b
    }
})

const authForm = reactive<UserToLogin | UserToRegister>({
    username: '',
    password: ''
})

const editItem = reactive({
    siteTitle: false,
    subtitleText: false,
    avatarURL: false,
    username: false,
    description: false,
    backgrounds: false,
    cardFooterTitle: false,
    cardFooterLink: false, 
    pageFooterHTML: false,
    rssTitle: false,
    rssDescription: false,
    rssAuthorName: false,
    rssFaviconURL: false,
    walineServerURL: false,
})

// 更新默认配置
const defaultConfig = {
    siteTitle: 'Noise的说说笔记',
    subtitleText: '欢迎访问，点击头像可更换封面背景！',
    avatarURL: 'https://s2.loli.net/2025/03/24/HnSXKvibAQlosIW.png',
    username: 'Noise',
    description: '执迷不悟',
    backgrounds: [
        "https://s2.loli.net/2025/03/27/KJ1trnU2ksbFEYM.jpg",
        "https://s2.loli.net/2025/03/27/MZqaLczCvwjSmW7.jpg",
        "https://s2.loli.net/2025/03/27/UMijKXwJ9yTqSeE.jpg",
        "https://s2.loli.net/2025/03/27/WJQIlkXvBg2afcR.jpg",
        "https://s2.loli.net/2025/03/27/oHNQtf4spkq2iln.jpg",
        "https://s2.loli.net/2025/03/27/PMRuX5loc6Uaimw.jpg",
        "https://s2.loli.net/2025/03/27/U2WIslbNyTLt4rD.jpg",
        "https://s2.loli.net/2025/03/27/xu1jZL5Og4pqT9d.jpg",
        "https://s2.loli.net/2025/03/27/OXqwzZ6v3PVIns9.jpg",
        "https://s2.loli.net/2025/03/27/HGuqlE6apgNywbh.jpg",
        "https://s2.loli.net/2025/03/26/d7iyuPYA8cRqD1K.jpg",
        "https://s2.loli.net/2025/03/27/wYy12qDMH6bGJOI.jpg",
        "https://s2.loli.net/2025/03/27/y67m2k5xcSdTsHN.jpg",
        ],
        cardFooterTitle: "Noise·说说·笔记~",
        cardFooterLink: "note.noisework.cn",
    pageFooterHTML: `<div class="text-center text-xs text-gray-400 py-4">来自<a href="https://www.noisework.cn" target="_blank" rel="noopener noreferrer" class="text-orange-400 hover:text-orange-500">Noise</a> 使用<a href="https://github.com/rcy1314/echo-noise" target="_blank" rel="noopener noreferrer" class="text-orange-400 hover:text-orange-500">Ech0-Noise</a>发布</div>`,
    rssTitle: 'Noise的说说笔记',
    rssDescription: '一个说说笔记~',
    rssAuthorName: 'Noise',
    rssFaviconURL: '/favicon.ico',
    walineServerURL: '请前往waline官网https://waline.js.org查看部署配置',
    githubOAuthEnabled: false,
    githubClientId: '',
    githubClientSecret: '',
    
    githubCallbackURL: '',
    // PWA 设置默认值（为空时回退到站点设置）
    pwaEnabled: true,
    pwaTitle: '',
    pwaDescription: '',
    pwaIconURL: ''
    ,defaultContentTheme: 'dark'
    ,announcementText: '欢迎访问我的说说笔记！'
    ,announcementEnabled: true
    ,musicEnabled: false
    ,musicPlaylistId: '2141128031'
    ,musicSongId: ''
    ,musicPosition: 'bottom-left'
    ,musicTheme: 'auto'
    ,musicLyric: true
    ,musicAutoplay: false
    ,musicDefaultMinimized: true
    ,musicEmbed: false
}
// 添加单个配置项保存方法

// 添加单个配置项重置方法
const resetConfigItem = (key: string) => {
    frontendConfig[key] = defaultConfig[key]
    editItem[key] = false
}
// 修改 fetchConfig 方法// ... existing code ...

const fetchConfig = async () => {
    try {
        const response = await fetch('/api/frontend/config', {
            credentials: 'include',
            headers: {
                'Cache-Control': 'no-cache',
                'Pragma': 'no-cache'
            }
        });
        
        const data = await response.json();
        
        if (data?.data?.frontendSettings) {
            const settings = data.data.frontendSettings;
            
            // 遍历配置项进行更新（布尔型键需强制转换）
            const booleanKeys = ['enableGithubCard', 'pwaEnabled', 'announcementEnabled', 'musicEnabled', 'musicLyric', 'musicAutoplay', 'musicDefaultMinimized', 'musicEmbed', 'commentEnabled', 'commentEmailEnabled', 'githubOAuthEnabled']
            Object.keys(frontendConfig).forEach(key => {
                if (key === 'backgrounds') {
                    const serverBackgrounds = settings[key];
                    if (Array.isArray(serverBackgrounds)) {
                        frontendConfig[key] = [...serverBackgrounds];
                    }
                } else if (booleanKeys.includes(key)) {
                    const v = settings[key] ?? defaultConfig[key]
                    ;(frontendConfig as any)[key] = (v === true || v === 'true')
                } else {
                    const v = settings[key] ?? defaultConfig[key]
                    ;(frontendConfig as any)[key] = typeof v === 'string' ? v.trim() : v
                }
            });

            // 独立处理布尔型未包含在 frontendConfig 键中的字段
            if (settings.enableGithubCard !== undefined) {
                const v = settings.enableGithubCard
                // @ts-ignore 动态添加字段
                frontendConfig.enableGithubCard = (v === true || v === 'true')
            }

            // 自动应用到页面 Head（标题、描述、图标）
            const title = (frontendConfig.pwaTitle || frontendConfig.siteTitle || '说说笔记').trim()
            const icon = (frontendConfig.rssFaviconURL || '/favicon.ico').trim()
            const description = (frontendConfig.pwaDescription || frontendConfig.description || '').trim()
            useHead({
                title,
                meta: [
                    { name: 'description', content: description },
                    { name: 'theme-color', content: '#000000' }
                ],
                link: [
                    { rel: 'manifest', href: '/manifest.webmanifest' },
                    { rel: 'icon', href: icon }
                ]
            })
        }
    } catch (error) {
        console.error('获取配置失败:', error);
    }
};
const saveConfigItem = async (key: string) => {
    try {
        // 特殊处理背景图片数组
        if (key === 'backgrounds') {
            const validBackgrounds = frontendConfig.backgrounds.filter(url => url && url.trim() !== '');
            frontendConfig.backgrounds = validBackgrounds;
        }

        const settingsToSave = {
            frontendSettings: frontendConfig  // 直接发送整个配置对象
        };

        const response = await fetch('/api/settings', {
            method: 'PUT',
            headers: {
                'Content-Type': 'application/json'
            },
            credentials: 'include',
            body: JSON.stringify(settingsToSave)
        });
        
        if (!response.ok) {
            const errorData = await response.json();
            throw new Error(errorData.msg || '请求失败');
        }
        
        const data = await response.json();
        if (data.code === 1) {
            editItem[key] = false;
            // 重新获取配置
            await fetchConfig();
            // 特殊提示：公告开关
            if (key === 'announcementEnabled') {
                const enabled = !!frontendConfig.announcementEnabled
                useToast().add({
                    title: '成功',
                    description: enabled ? '已开启公告' : '已关闭公告',
                    color: enabled ? 'green' : 'gray'
                })
            } else {
                const label = key === 'defaultContentTheme' ? '默认主题色' : (configLabels[key] || (key === 'pwa' ? 'PWA 设置' : key))
                useToast().add({
                    title: '成功',
                    description: `${label}已更新`,
                    color: 'green'
                })
            }
            if (key === 'defaultContentTheme') {
                const theme = (frontendConfig.defaultContentTheme || 'dark').trim();
                // 不触发任何前端切换，仅在后续首次加载时生效
            }
        } else {
            throw new Error(data.msg || '保存失败');
        }
    } catch (error: any) {
        console.error('保存配置失败:', error);
        const label = key === 'defaultContentTheme' ? '默认主题色' : (configLabels[key] || key)
        useToast().add({
            title: '失败',
            description: `${label}保存失败`,
            color: 'red'
        });
    }
};

const savePWAConfig = async () => {
    try {
        const settingsToSave = {
            frontendSettings: frontendConfig
        }
        const response = await fetch('/api/settings', {
            method: 'PUT',
            headers: { 'Content-Type': 'application/json' },
            credentials: 'include',
            body: JSON.stringify(settingsToSave)
        })
        const data = await response.json()
        if (response.ok && data.code === 1) {
            await fetchConfig()
            useToast().add({ title: '成功', description: 'PWA 设置已更新', color: 'green' })

            // 立即切换 Service Worker 状态
            if ('serviceWorker' in navigator) {
                const regs = await navigator.serviceWorker.getRegistrations()
                if (frontendConfig.pwaEnabled) {
                    await navigator.serviceWorker.register('/sw.js')
                } else {
                    for (const r of regs) await r.unregister()
                    const keys = await caches.keys()
                    await Promise.all(keys.map(k => caches.delete(k)))
                }
            }

            // 通知全局插件重新应用 Head 与 SW 状态
            window.dispatchEvent(new Event('frontend-config-updated'))
        } else {
            throw new Error(data.msg || '保存失败')
        }
    } catch (error: any) {
        useToast().add({ title: '错误', description: error.message || '保存失败', color: 'red' })
    }
}

const saveCommentConfig = async () => {
  try {
    const payload = {
      frontendSettings: {
        commentEnabled: !!frontendConfig.commentEnabled,
        commentSystem: String(frontendConfig.commentSystem || 'waline'),
        commentEmailEnabled: !!frontendConfig.commentEmailEnabled,
        walineServerURL: String(frontendConfig.walineServerURL || '')
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
      await fetchConfig()
      window.dispatchEvent(new Event('frontend-config-updated'))
      useToast().add({ title: '成功', description: '评论设置已更新', color: 'green' })
    } else {
      throw new Error(data.msg || '保存失败')
    }
  } catch (error: any) {
    useToast().add({ title: '错误', description: error.message || '保存失败', color: 'red' })
  }
}

// 保存 GitHub 卡片解析配置（独立项）
const saveGithubCardConfig = async () => {
    try {
        const payload = {
            frontendSettings: {
                enableGithubCard: !!githubCardEnabled.value
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
            // 同步本地状态
            // @ts-ignore
            frontendConfig.enableGithubCard = !!githubCardEnabled.value
            await fetchConfig()
            window.dispatchEvent(new Event('frontend-config-updated'))
            
            useToast().add({ title: '成功', description: 'GitHub 解析设置已保存', color: 'green' })
        } else {
            throw new Error(data.msg || '保存失败')
        }
    } catch (error: any) {
        useToast().add({ title: '错误', description: error.message || '保存失败', color: 'red' })
    }
}

// 音乐配置保存与重置
const saveMusicConfig = async () => {
  try {
    const payload = {
      frontendSettings: {
        musicEnabled: !!frontendConfig.musicEnabled,
        musicPlaylistId: String(frontendConfig.musicPlaylistId || ''),
        musicSongId: String(frontendConfig.musicSongId || ''),
        musicPosition: String(frontendConfig.musicPosition || 'bottom-left'),
        musicTheme: String(frontendConfig.musicTheme || 'auto'),
        musicLyric: !!frontendConfig.musicLyric,
        musicAutoplay: !!frontendConfig.musicAutoplay,
        musicDefaultMinimized: !!frontendConfig.musicDefaultMinimized,
        musicEmbed: !!frontendConfig.musicEmbed
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
      await fetchConfig()
      window.dispatchEvent(new Event('frontend-config-updated'))
      useToast().add({ title: '成功', description: '音乐配置已更新', color: 'green' })
    } else {
      throw new Error(data.msg || '保存失败')
    }
  } catch (error: any) {
    useToast().add({ title: '错误', description: error.message || '保存失败', color: 'red' })
  }
}

const resetMusicConfig = () => {
  (frontendConfig as any).musicEnabled = defaultConfig.musicEnabled
  ;(frontendConfig as any).musicPlaylistId = defaultConfig.musicPlaylistId
  ;(frontendConfig as any).musicSongId = defaultConfig.musicSongId
  ;(frontendConfig as any).musicPosition = defaultConfig.musicPosition
  ;(frontendConfig as any).musicTheme = defaultConfig.musicTheme
  ;(frontendConfig as any).musicLyric = defaultConfig.musicLyric
  ;(frontendConfig as any).musicAutoplay = defaultConfig.musicAutoplay
  ;(frontendConfig as any).musicDefaultMinimized = defaultConfig.musicDefaultMinimized
  ;(frontendConfig as any).musicEmbed = defaultConfig.musicEmbed
}

const toggleMusic = async (enabled: boolean) => {
  ;(frontendConfig as any).musicEnabled = enabled
  if (enabled) {
    if (!(frontendConfig as any).musicPlaylistId) {
      ;(frontendConfig as any).musicPlaylistId = '2141128031'
    }
    ;(frontendConfig as any).musicPosition = 'bottom-left'
    ;(frontendConfig as any).musicDefaultMinimized = true
    ;(frontendConfig as any).musicAutoplay = false
    ;(frontendConfig as any).musicTheme = 'auto'
  }
  await saveMusicConfig()
}
const musicEmbedMode = computed({
  get: () => (frontendConfig.musicEmbed ? 'embed' : 'float'),
  set: (v: string) => { (frontendConfig as any).musicEmbed = (v === 'embed') }
})

const saveGithubOAuthConfig = async () => {
  try {
    await saveConfigItem('githubOAuthEnabled')
    await saveConfigItem('githubClientId')
    await saveConfigItem('githubClientSecret')
    await saveConfigItem('githubCallbackURL')
    useToast().add({ title: '成功', description: 'GitHub 登录配置已保存', color: 'green' })
  } catch (error: any) {
    useToast().add({ title: '错误', description: error?.message || '保存失败', color: 'red' })
  }
}

const applyPWAConfig = () => {
    const title = (frontendConfig.pwaTitle || frontendConfig.siteTitle || '说说笔记')
    const icon = (frontendConfig.rssFaviconURL || '/favicon.ico')
    const description = (frontendConfig.pwaDescription || frontendConfig.description || '')

    useHead({
        title,
        meta: [
            { name: 'description', content: description },
            { name: 'theme-color', content: '#000000' }
        ],
        link: [
            { rel: 'manifest', href: '/manifest.webmanifest' },
            { rel: 'icon', href: icon }
        ]
    })
}

// 修改文件上传处理
const handleFileUpload = async (event: Event) => {
    const files = (event.target as HTMLInputElement).files
    if (!files) return

    const allowedTypes = ['image/jpeg', 'image/png', 'image/webp']
    
    for (const file of Array.from(files)) {
        try {
            if (!allowedTypes.includes(file.type)) {
                throw new Error('仅支持 JPG/PNG/WEBP 格式')
            }

            const formData = new FormData()
            formData.append('image', file)

            const response = await fetch('/api/images/upload', {
                method: 'POST',
                credentials: 'include',
                body: formData
            })

            const data = await response.json()
            
            if (!response.ok || data.code !== 1) {
                throw new Error(data.msg || '上传失败')
            }

            if (data.code === 1 && data.data) {
                const imageUrl = data.data.startsWith('http') 
                    ? data.data 
                    : `/api${data.data}`
                
                // 更新背景图片列表并保存
                const newBackgrounds = [...frontendConfig.backgrounds, imageUrl];
                frontendConfig.backgrounds = newBackgrounds;
                await saveConfigItem('backgrounds');

                useToast().add({
                    title: '上传成功',
                    description: `${file.name} 已添加到背景图片列表`,
                    color: 'green'
                })
            }
        } catch (error: any) {
            useToast().add({
                title: '上传失败',
                description: error.message || '文件上传失败',
                color: 'red'
            })
        }
    }

    // 清空文件输入框
    if (fileInput.value) {
        fileInput.value.value = ''
    }
}

// 添加配置更新监听器
onMounted(() => {
    window.addEventListener('frontend-config-updated', (event: any) => {
        const { key, value } = event.detail;
        if (key && value !== undefined) {
            frontendConfig[key] = value;
        }
    });
});
// ... existing code ...
const resetConfig = () => {
    fetchConfig()
    editMode.value = false
}

const addBackground = async () => {
    frontendConfig.backgrounds.push(''); 
}

const removeBackground = async (index: number) => {
    frontendConfig.backgrounds.splice(index, 1);
    await saveConfigItem('backgrounds');
}

const triggerFileInput = () => {
    fileInput.value?.click()
}
const previewImage = (url: string) => {
    window.open(url, '_blank')
}

// 监听器
watch(() => userStore.isLogin, (newVal) => {
    if (!newVal) {
        userStore.getStatus()
        userStore.getUser()
        userStore.$reset()
    }
})

// 生命周期
const isLoading = ref(false) // 新增加载状态

onMounted(async () => {
    try {
        isLoading.value = true;
        // 先获取用户状态和配置
        await Promise.all([
            userStore.getStatus(),
            userStore.getUser(),
            fetchConfig(),
            fetchRegisterConfig()
        ]);

        // 如果用户已登录，再获取 token
        if (userStore.isLogin) {
            const response = await fetch('/api/user/token', {
                method: 'GET',
                credentials: 'include',
                headers: {
                    'Content-Type': 'application/json'
                }
            });

            const data = await response.json();
            if (data.code === 1 && data.data?.token) {
                userToken.value = data.data.token;
            }
        }
    } catch (error) {
        console.error('初始化失败:', error);
        useToast().add({
            title: '系统提示',
            description: '当前未登录',
            color: 'red',
            timeout: 3000
        });
    } finally {
        isLoading.value = false;
    }
});
const databaseFileInput = ref<HTMLInputElement | null>(null)

const downloadBackup = async () => {
    try {
        const response = await fetch('/api/backup/download', {
            credentials: 'include'
        })
        
        if (!response.ok) {
            throw new Error('下载失败')
        }

        const blob = await response.blob()
        const url = window.URL.createObjectURL(blob)
        const a = document.createElement('a')
        a.href = url
        a.download = `noise_backup_${new Date().toISOString().slice(0,10)}.zip`
        document.body.appendChild(a)
        a.click()
        window.URL.revokeObjectURL(url)
        document.body.removeChild(a)
    } catch (error) {
        useToast().add({
            title: '错误',
            description: error.message || '备份下载失败',
            color: 'red'
        })
    }
}

const triggerDatabaseUpload = () => {
    databaseFileInput.value?.click()
}
const emit = defineEmits(['restore-success'])
const handleDatabaseUpload = async (event: Event) => {
    const files = (event.target as HTMLInputElement).files
    if (!files || !files[0]) return

    try {
        const formData = new FormData()
        formData.append('database', files[0])

        const response = await fetch('/api/backup/restore', {
            method: 'POST',
            credentials: 'include',
            body: formData
        })

        const data = await response.json()
        if (data.code === 1) {
            useToast().add({
                title: '成功',
                description: '数据库恢复成功',
                color: 'green'
            })
            emit('restore-success')
            // 添加成功后刷新页面
            setTimeout(() => {
                window.location.reload()
            }, 1500)
        } else {
            throw new Error(data.msg)
        }
    } catch (error) {
        useToast().add({
            title: '错误',
            description: error.message || '数据库恢复失败',
            color: 'red'
        })
    }

    if (databaseFileInput.value) {
        databaseFileInput.value.value = ''
    }
}
const storageEnabled = ref(false)
const storageConfig = reactive({
  provider: '',
  endpoint: '',
  region: '',
  bucket: '',
  accessKey: '',
  secretKey: '',
  usePathStyle: true,
  publicBaseURL: ''
})
const uploadURL = ref('')
const downloadURL = ref('')
const loadStorageConfig = async () => {
  try {
    const res = await fetch('/api/frontend/config', { credentials: 'include' })
    const data = await res.json()
    if (data?.code === 1) {
      storageEnabled.value = !!data.data.storageEnabled
      const sc = data.data.storageConfig || {}
      storageConfig.provider = sc.provider || ''
      storageConfig.endpoint = sc.endpoint || ''
      storageConfig.region = sc.region || ''
      storageConfig.bucket = sc.bucket || ''
      storageConfig.accessKey = sc.accessKey || ''
      storageConfig.secretKey = sc.secretKey || ''
      storageConfig.usePathStyle = !!sc.usePathStyle
      storageConfig.publicBaseURL = sc.publicBaseURL || ''
    }
  } catch {}
}
const saveStorageConfig = async () => {
  try {
    const payload: any = {
      storageEnabled: storageEnabled.value,
      storageConfig: { ...storageConfig }
    }
    const res = await fetch('/api/settings', { method: 'PUT', credentials: 'include', headers: { 'Content-Type': 'application/json' }, body: JSON.stringify(payload) })
    const data = await res.json()
    if (data?.code === 1) {
      useToast().add({ title: '已保存云存储配置', color: 'green' })
    } else {
      throw new Error(data?.msg || '保存失败')
    }
  } catch (e: any) {
    useToast().add({ title: '保存失败', description: e.message, color: 'red' })
  }
}
const uploadCloudBackup = async () => {
  try {
    const url = uploadURL.value.trim()
    if (!url) throw new Error('请填写预签名上传URL')
    const res = await postRequest<any>('backup/storage/upload', { uploadURL: url }, { credentials: 'include' })
    if (res?.code === 1) {
      useToast().add({ title: '云备份上传成功', color: 'green' })
    } else {
      throw new Error(res?.msg || '上传失败')
    }
  } catch (e: any) {
    useToast().add({ title: '上传失败', description: e.message, color: 'red' })
  }
}
const restoreCloudBackup = async () => {
  try {
    const url = downloadURL.value.trim()
    if (!url) throw new Error('请填写预签名下载URL')
    const res = await postRequest<any>('backup/storage/restore', { downloadURL: url }, { credentials: 'include' })
    if (res?.code === 1) {
      useToast().add({ title: '云备份恢复成功', color: 'green' })
    } else {
      throw new Error(res?.msg || '恢复失败')
    }
  } catch (e: any) {
    useToast().add({ title: '恢复失败', description: e.message, color: 'red' })
  }
}
const positionOptions = [
  { label: '静态', value: 'static' },
  { label: '左上', value: 'top-left' },
  { label: '右上', value: 'top-right' },
  { label: '左下', value: 'bottom-left' },
  { label: '右下', value: 'bottom-right' },
]
const themeOptions = [
  { label: '自动', value: 'auto' },
  { label: '浅色', value: 'light' },
  { label: '深色', value: 'dark' },
]
</script>

<style scoped>
.hidden {
    display: none;
}
</style>
