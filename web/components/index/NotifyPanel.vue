<template>
    <div class="rounded-lg p-4 mb-6">
        <div class="flex justify-between items-center mb-4">
            <h2 class="text-xl font-semibold" :class="text">推送渠道配置</h2>
        </div>

        <!-- 配置列表 -->
        <div class="space-y-4" :class="props.disabled ? 'opacity-60 pointer-events-none' : ''">
            <!-- Webhook配置 -->
            <div class="rounded p-3" :class="subtleBg">
                <div class="flex justify-between items-center mb-2">
                    <span :class="mutedText">Webhook推送</span>
                    <div class="flex items-center gap-2">
                        <USwitch v-model="localConfig.webhookEnabled" />
                        <UButton
                            size="sm"
                            @click="editItem.webhook = !editItem.webhook"
                            :color="editItem.webhook ? 'gray' : 'primary'"
                            variant="soft"
                        >
                            {{ editItem.webhook ? '取消' : '编辑' }}
                        </UButton>
                    </div>
                </div>
                <div v-if="editItem.webhook" class="space-y-2">
                    <UInput
                        v-model="localConfig.webhookURL"
                        placeholder="Webhook URL"
                    />
                    <div class="flex justify-end gap-2">
                        <UButton @click="resetConfigItem('webhook')" variant="ghost" color="gray">
                            重置
                        </UButton>
                        <UButton @click="saveConfigItem('webhook')" color="primary">
                            保存
                        </UButton>
                    </div>
                </div>
            </div>

            <!-- Telegram配置 -->
            <div class="rounded p-3" :class="subtleBg">
                <div class="flex justify-between items-center mb-2">
                    <span :class="mutedText">Telegram推送</span>
                    <div class="flex items-center gap-2">
                        <USwitch v-model="localConfig.telegramEnabled" />
                        <UButton
                            size="sm"
                            @click="editItem.telegram = !editItem.telegram"
                            :color="editItem.telegram ? 'gray' : 'primary'"
                            variant="soft"
                        >
                            {{ editItem.telegram ? '取消' : '编辑' }}
                        </UButton>
                    </div>
                </div>
                <div v-if="editItem.telegram" class="space-y-2">
                    <UInput
                        v-model="localConfig.telegramToken"
                        placeholder="Bot Token"
                    />
                    <UInput
                        v-model="localConfig.telegramChatID"
                        placeholder="Chat ID"
                    />
                    <div class="flex justify-end gap-2">
                        <UButton @click="resetConfigItem('telegram')" variant="ghost" color="gray">
                            重置
                        </UButton>
                        <UButton @click="saveConfigItem('telegram')" color="primary">
                            保存
                        </UButton>
                    </div>
                </div>
            </div>

            <!-- 企业微信配置 -->
            <div class="rounded p-3" :class="subtleBg">
                <div class="flex justify-between items-center mb-2">
                    <span :class="mutedText">企业微信推送</span>
                    <div class="flex items-center gap-2">
                        <USwitch v-model="localConfig.weworkEnabled" />
                        <UButton
                            size="sm"
                            @click="editItem.wework = !editItem.wework"
                            :color="editItem.wework ? 'gray' : 'primary'"
                            variant="soft"
                        >
                            {{ editItem.wework ? '取消' : '编辑' }}
                        </UButton>
                    </div>
                </div>
                <div v-if="editItem.wework" class="space-y-2">
                    <UInput
                        v-model="localConfig.weworkKey"
                        placeholder="WebHook Key"
                    />
                    <div class="flex justify-end gap-2">
                        <UButton @click="resetConfigItem('wework')" variant="ghost" color="gray">
                            重置
                        </UButton>
                        <UButton @click="saveConfigItem('wework')" color="primary">
                            保存
                        </UButton>
                    </div>
                </div>
            </div>

            <!-- 飞书配置 -->
            <div class="rounded p-3" :class="subtleBg">
                <div class="flex justify-between items-center mb-2">
                    <span :class="mutedText">飞书推送</span>
                    <div class="flex items-center gap-2">
                        <USwitch v-model="localConfig.feishuEnabled" />
                        <UButton
                            size="sm"
                            @click="editItem.feishu = !editItem.feishu"
                            :color="editItem.feishu ? 'gray' : 'primary'"
                            variant="soft"
                        >
                            {{ editItem.feishu ? '取消' : '编辑' }}
                        </UButton>
                    </div>
                </div>
                <div v-if="editItem.feishu" class="space-y-2">
                    <UInput
                        v-model="localConfig.feishuWebhook"
                        placeholder="WebHook URL"
                    />
                    <UInput
                        v-model="localConfig.feishuSecret"
                        placeholder="签名密钥"
                        type="password"
                    />
                    <div class="flex justify-end gap-2">
                        <UButton @click="resetConfigItem('feishu')" variant="ghost" color="gray">
                            重置
                        </UButton>
                        <UButton @click="saveConfigItem('feishu')" color="primary">
                            保存
                        </UButton>
                    </div>
                </div>
            </div> <!-- 这里补上飞书配置块的闭合 -->
            <!-- Twitter配置 -->
            <div class="rounded p-3" :class="subtleBg">
                <div class="flex justify-between items-center mb-2">
                    <span :class="mutedText">Twitter推送</span>
                    <div class="flex items-center gap-2">
                        <USwitch v-model="localConfig.twitterEnabled" />
                        <UButton
                            size="sm"
                            @click="editItem.twitter = !editItem.twitter"
                            :color="editItem.twitter ? 'gray' : 'primary'"
                            variant="soft"
                        >
                            {{ editItem.twitter ? '取消' : '编辑' }}
                        </UButton>
                    </div>
                </div>
                <div v-if="editItem.twitter" class="space-y-2">
                    <UInput v-model="localConfig.twitterApiKey" placeholder="API Key" />
                    <UInput v-model="localConfig.twitterApiSecret" placeholder="API Secret" />
                    <UInput v-model="localConfig.twitterAccessToken" placeholder="Access Token" />
                    <UInput v-model="localConfig.twitterAccessTokenSecret" placeholder="Access Token Secret" />
                    <div class="flex justify-end gap-2">
                        <UButton @click="resetConfigItem('twitter')" variant="ghost" color="gray">重置</UButton>
                        <UButton @click="saveConfigItem('twitter')" color="primary">保存</UButton>
                    </div>
                </div>
            </div>

            <!-- 自定义HTTP配置 -->
            <div class="rounded p-3" :class="subtleBg">
                <div class="flex justify-between items-center mb-2">
                    <span :class="mutedText">自定义HTTP推送</span>
                    <div class="flex items-center gap-2">
                        <USwitch v-model="localConfig.customHttpEnabled" />
                        <UButton
                            size="sm"
                            @click="editItem.customHttp = !editItem.customHttp"
                            :color="editItem.customHttp ? 'gray' : 'primary'"
                            variant="soft"
                        >
                            {{ editItem.customHttp ? '取消' : '编辑' }}
                        </UButton>
                    </div>
                </div>
                <div v-if="editItem.customHttp" class="space-y-2">
                    <UInput v-model="localConfig.customHttpUrl" placeholder="请求URL" />
                    <USelect v-model="localConfig.customHttpMethod" :options="['POST', 'PUT', 'PATCH']" />
                    <UTextarea v-model="localConfig.customHttpHeaders" placeholder="请求Headers（JSON格式）" />
                    <UTextarea v-model="localConfig.customHttpBody" placeholder='请求Body模板（如：{"content":"{{content}}"})' />
                    <div class="flex justify-end gap-2">
                        <UButton @click="resetConfigItem('customHttp')" variant="ghost" color="gray">重置</UButton>
                        <UButton @click="saveConfigItem('customHttp')" color="primary">保存</UButton>
                    </div>
                </div>
            </div>
        </div>

        <!-- 测试按钮区域 -->
        <div class="mt-4 flex flex-wrap gap-2" :class="props.disabled ? 'opacity-60 pointer-events-none' : ''">
    <UButton
        v-for="type in notifyTypes"
        :key="type"
        @click="testNotify(type)"
        class="flex-grow sm:flex-grow-0"
        variant="solid"
        :disabled="props.disabled"
    >
        测试{{ getNotifyTypeName(type) }}
    </UButton>
</div>
    </div>
</template>

<script setup lang="ts">
import { useToast } from '#ui/composables/useToast'
import { ref, reactive, watch, onMounted, computed } from 'vue'

// 添加 emit 定义
const emit = defineEmits(['save'])

// 添加 props 定义
const props = defineProps<{
    config: NotifyConfig;
    immediate?: boolean;
    subtleBg?: string;
    text?: string;
    mutedText?: string;
    disabled?: boolean;
}>();

// 添加类型定义
interface NotifyConfig {
    webhookEnabled: boolean;
    webhookURL: string;
    telegramEnabled: boolean;
    telegramToken: string;
    telegramChatID: string;
    weworkEnabled: boolean;
    weworkKey: string;
    feishuEnabled: boolean;
    feishuWebhook: string;
    feishuSecret: string;
    twitterEnabled: boolean;
    twitterApiKey: string;
    twitterApiSecret: string;
    twitterAccessToken: string;
    twitterAccessTokenSecret: string;
    customHttpEnabled: boolean;
    customHttpUrl: string;
    customHttpMethod: string;
    customHttpHeaders: string;
    customHttpBody: string;
}

const localConfig = reactive<NotifyConfig>({
    webhookEnabled: false,
    webhookURL: '',
    telegramEnabled: false,
    telegramToken: '',
    telegramChatID: '',
    weworkEnabled: false,
    weworkKey: '',
    feishuEnabled: false,
    feishuWebhook: '',
    feishuSecret: '',
    twitterEnabled: false,
    twitterApiKey: '',
    twitterApiSecret: '',
    twitterAccessToken: '',
    twitterAccessTokenSecret: '',
    customHttpEnabled: false,
    customHttpUrl: '',
    customHttpMethod: 'POST',
    customHttpHeaders: '',
    customHttpBody: '{"content":"{{content}}"}'
})

const editItem = reactive({
    webhook: false,
    telegram: false,
    wework: false,
    feishu: false,
    // 新增
    twitter: false,
    customHttp: false
})

const notifyTypes = ['webhook', 'telegram', 'wework', 'feishu', 'twitter', 'customHttp']

const getNotifyTypeName = (type) => {
    const names = {
        webhook: 'Webhook',
        telegram: 'Telegram',
        wework: '企业微信',
        feishu: '飞书',
        twitter: 'Twitter',
        customHttp: '自定义HTTP'
    }
    return names[type] || type
}

// 重置配置项
const resetConfigItem = (type: string) => {
    const prefix = type.toLowerCase()
    Object.keys(localConfig).forEach(key => {
        if (key.toLowerCase().startsWith(prefix)) {
            localConfig[key] = props.config[key]
        }
    })
    editItem[type] = false
}

// 修改保存配置项方法
// 在 script setup 部分添加配置验证函数
// 确保 isConfigValid 函数正确验证配置
const isConfigValid = (type: string) => {
    const config = localConfig;
    switch(type) {
        case 'webhook':
            return config.webhookEnabled && !!config.webhookURL?.trim();
        case 'telegram':
            return config.telegramEnabled && 
                   !!config.telegramToken?.trim() && 
                   !!config.telegramChatID?.trim();
        case 'wework':
            return config.weworkEnabled && !!config.weworkKey?.trim();
        case 'feishu':
            return config.feishuEnabled && !!config.feishuWebhook?.trim();
        case 'twitter':
            return config.twitterEnabled &&
                   !!config.twitterApiKey?.trim() &&
                   !!config.twitterApiSecret?.trim() &&
                   !!config.twitterAccessToken?.trim() &&
                   !!config.twitterAccessTokenSecret?.trim();
        case 'customHttp':
            return config.customHttpEnabled && !!config.customHttpUrl?.trim();
        default:
            return false;
    }
}

// 修改 saveConfigItem 方法中的调试日志
const saveConfigItem = async (type: string) => {
    try {
        // 根据类型设置启用状态
        switch(type) {
            case 'webhook':
                if (localConfig.webhookEnabled && !localConfig.webhookURL) {
                    throw new Error('Webhook URL不能为空');
                }
                break;
            case 'telegram':
                if (localConfig.telegramEnabled && (!localConfig.telegramToken || !localConfig.telegramChatID)) {
                    throw new Error('Telegram配置不完整');
                }
                break;
            case 'wework':
                if (localConfig.weworkEnabled && !localConfig.weworkKey) {
                    throw new Error('企业微信Key不能为空');
                }
                break;
            case 'feishu':
                if (localConfig.feishuEnabled && !localConfig.feishuWebhook) {
                    throw new Error('飞书Webhook不能为空');
                }
                break;
            case 'twitter':
                if (localConfig.twitterEnabled && (
                    !localConfig.twitterApiKey ||
                    !localConfig.twitterApiSecret ||
                    !localConfig.twitterAccessToken ||
                    !localConfig.twitterAccessTokenSecret
                )) {
                    throw new Error('Twitter配置不完整');
                }
                break;
            case 'customHttp':
                if (localConfig.customHttpEnabled && !localConfig.customHttpUrl) {
                    throw new Error('自定义HTTP URL不能为空');
                }
                break;
        }

        // 修复：补全所有推送渠道的配置字段
        const configData = {
            webhookEnabled: localConfig.webhookEnabled,
            webhookURL: localConfig.webhookURL,
            telegramEnabled: localConfig.telegramEnabled,
            telegramToken: localConfig.telegramToken,
            telegramChatID: localConfig.telegramChatID,
            weworkEnabled: localConfig.weworkEnabled,
            weworkKey: localConfig.weworkKey,
            feishuEnabled: localConfig.feishuEnabled,
            feishuWebhook: localConfig.feishuWebhook,
            feishuSecret: localConfig.feishuSecret,
            // 新增 Twitter 配置字段
            twitterEnabled: localConfig.twitterEnabled,
            twitterApiKey: localConfig.twitterApiKey,
            twitterApiSecret: localConfig.twitterApiSecret,
            twitterAccessToken: localConfig.twitterAccessToken,
            twitterAccessTokenSecret: localConfig.twitterAccessTokenSecret,
            // 新增自定义HTTP配置字段
            customHttpEnabled: localConfig.customHttpEnabled,
            customHttpUrl: localConfig.customHttpUrl,
            customHttpMethod: localConfig.customHttpMethod,
            customHttpHeaders: localConfig.customHttpHeaders,
            customHttpBody: localConfig.customHttpBody
        };

        const baseApi = useRuntimeConfig().public.baseApi || '/api'
        const response = await fetch(`${baseApi}/notify/config`, {
            method: 'PUT',
            headers: {
                'Content-Type': 'application/json'
            },
            credentials: 'include',
            body: JSON.stringify(configData)
        });
        
        const data = await response.json();
        if (data.code === 1) {
            editItem[type] = false;
            
            // 强制更新所有配置字段，确保启用状态同步
            if (data.data) {
                Object.keys(data.data).forEach(key => {
                    // 修复：所有 xxxEnabled 字段都强制转为布尔值
                    if (key.endsWith('Enabled')) {
                        localConfig[key] = !!data.data[key] && data.data[key] !== 'false';
                    } else {
                        localConfig[key] = data.data[key] || '';
                    }
                });
            }
            
            emit('save', {...localConfig});
            useToast().add({
                title: '成功',
                description: `${getNotifyTypeName(type)}配置已保存`,
                color: 'green'
            });
        } else {
            throw new Error(data.msg || '保存失败');
        }
    } catch (error) {
        console.error('保存配置失败:', error);
        useToast().add({
            title: '错误',
            description: error.message || '保存失败',
            color: 'red'
        });
    }
};

// 修改测试推送方法
const testNotify = async (type: string) => {
    try {
        // 检查对应类型的配置是否启用和有效
        if (!isConfigValid(type)) {
            throw new Error(`请先完成并启用${getNotifyTypeName(type)}配置`);
        }

        const response = await fetch('/api/notify/test', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            credentials: 'include',
            body: JSON.stringify({ type })
        });
        
        const data = await response.json();
        if (data.code === 1) {
            useToast().add({
                title: '成功',
                description: `${getNotifyTypeName(type)}测试消息已发送`,
                color: 'green'
            });
        } else {
            throw new Error(data.msg || '测试失败');
        }
    } catch (error) {
        console.error('测试失败:', error);
        useToast().add({
            title: '错误',
            description: error.message || '测试失败',
            color: 'red'
        });
    }
};


// 初始化配置
watch(() => props.config, (newConfig) => {
    if (newConfig) {
        Object.keys(localConfig).forEach(key => {
            if (newConfig[key] !== undefined) {
                // 修复：所有 xxxEnabled 字段都强制转为布尔值
                if (key.endsWith('Enabled')) {
                    localConfig[key] = !!newConfig[key] && newConfig[key] !== 'false';
                } else {
                    localConfig[key] = newConfig[key];
                }
            }
        });
    }
}, { deep: true, immediate: true });

// 添加页面加载时获取配置的逻辑
onMounted(async () => {
    try {
        const response = await fetch('/api/notify/config', {
            method: 'GET',
            credentials: 'include',
            headers: {
                'Content-Type': 'application/json'
            }
        });
        const data = await response.json();
        if (data.code === 1 && data.data) {
            const config = data.data;
            Object.keys(config).forEach(key => {
                // 修复：所有 xxxEnabled 字段都强制转为布尔值
                if (key.endsWith('Enabled')) {
                    localConfig[key] = !!config[key] && config[key] !== 'false';
                } else {
                    localConfig[key] = config[key] || '';
                }
            });
            console.log('加载的配置:', localConfig);
        }
    } catch (error) {
        console.error('获取配置失败:', error);
        useToast().add({
            title: '错误',
            description: '获取推送配置失败',
            color: 'red'
        });
    }
});
const subtleBg = computed(() => props.subtleBg || 'bg-gray-800')
const text = computed(() => props.text || 'text-white')
const mutedText = computed(() => props.mutedText || 'text-gray-300')
</script>
