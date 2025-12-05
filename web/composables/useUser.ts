import type { UserToRegister, UserToLogin } from '~/types/models';
import { useUserStore } from '~/store/user';

export const useUser = () => {
    const userStore = useUserStore();
    const toast = useToast()
    const router = useRouter()

    // 添加状态检查函数
    const checkLoginStatus = async () => {
        const status = await userStore.getStatus();
        if (status) {
            userStore.setUserStatus(status);
        }
        return status;
    }

    const register = async (userToRegister: UserToRegister) => {
        const response = await userStore.register(userToRegister);
        if (response) {
            toast.add({
                title: '注册成功',
                description: '欢迎使用！请登录',
                icon: 'i-fluent-checkmark-starburst-16-filled',
                color: 'green',
                timeout: 1000,
            });
        }
    }

    const login = async (userToLogin: UserToLogin) => {
        const response = await userStore.login(userToLogin);
        if (response) {
            await checkLoginStatus();
            toast.add({
                title: '登录成功',
                description: '欢迎回来！',
                icon: 'i-fluent-checkmark-starburst-16-filled',
                color: 'green',
                timeout: 1000,
            });
            // 页面自行决定跳转（登录页处理 redirect->/status）
        }
    }

    const logout = async () => {
        await userStore.logout()
        // 退出后清除状态
        userStore.clearUserStatus();
        toast.add({
            title: '注销成功',
            description: '欢迎再次使用！',
            icon: 'i-fluent-checkmark-starburst-16-filled',
            color: 'green',
            timeout: 1000,
        });
        router.push('/')
    }

    const getStatus = async () => {
        return await checkLoginStatus();
    }

    // 初始化时检查登录状态
    onMounted(async () => {
        await checkLoginStatus();
    })

    return {
        register,
        login,
        logout,
        getStatus,
    }
}
