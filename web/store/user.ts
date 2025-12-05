import type { User, Status, UserToLogin, UserToRegister, Response } from "~/types/models"

export const useUserStore = defineStore("userStore", () => {
    // 状态
    const user = ref<User | null>(null);
    const status = ref<Status | null>(null);
    const isLogin = ref<boolean>(false);
    const token = ref<string>("");
    const toast = useToast()

    // 设置用户状态
    const setUserStatus = (newStatus: Status) => {
        status.value = newStatus as any;
        const list = (newStatus as any).users || (newStatus as any).Users || []
        const u = list.find((it: any) => (it.user_id ?? it.ID) === (user.value as any)?.userid)
        if (u) {
            user.value = {
                userid: u.user_id ?? u.ID,
                username: u.username ?? u.Username,
                is_admin: u.is_admin ?? u.IsAdmin,
                avatar_url: u.avatar_url ?? u.AvatarURL,
                total_messages: (newStatus as any).total_messages ?? 0
            } as any
            isLogin.value = true
        }
    }

    // 清除用户状态
    const clearUserStatus = () => {
        status.value = null;
        user.value = null;
        isLogin.value = false;
        token.value = "";
    }

    // 注册
    const register = async (userToRegister: UserToRegister) => {
        const response = await postRequest<any>("register", userToRegister, {
            credentials: 'include'
        });
        if (!response || response.code !== 1) {
            console.log("注册失败");
            toast.add({
                title: "注册失败",
                description: response?.msg,
                icon: "i-fluent-error-circle-16-filled",
                color: "red",
                timeout: 2000,
            });
            return false;
        }

        return response.code === 1;
    };

    // 登录
    const login = async (userToLogin: UserToLogin) => {
        const response = await postRequest<User>("login", userToLogin, {
            credentials: 'include'
        });
        if (!response || response.code !== 1) {
            console.log("登录失败");
            toast.add({
                title: "登录失败",
                description: response?.msg,
                icon: "i-fluent-error-circle-16-filled",
                color: "red",
                timeout: 2000,
            });
            return false;
        }

        if (response && response.code === 1 && response.data) {
            const u: any = response.data as any
            user.value = {
                userid: u.id ?? u.ID ?? u.user_id,
                username: u.username ?? u.Username,
                is_admin: u.is_admin ?? u.IsAdmin,
                avatar_url: u.avatar_url ?? u.AvatarURL,
                description: u.description ?? u.Description
            } as any
            token.value = u.token ?? u.Token ?? token.value
            isLogin.value = true;
            await getStatus();
            return true;
        }

        return false;
    }

    // 获取状态
    const getStatus = async () => {
        const response = await getRequest<Status>("status", undefined, {
            credentials: 'include'
        });
        if (!response || response.code !== 1) {
            console.log("获取系统信息失败");
            toast.add({
                title: "获取系统信息失败",
                description: response?.msg,
                icon: "i-fluent-error-circle-16-filled",
                color: "red",
                timeout: 2000,
            });
            return false;
        }

        if (response && response.code === 1 && response.data) {
            setUserStatus(response.data);
            return response.data;
        }
        return null;
    }

    // 获取当前登录用户信息
    const getUser = async (showToast: boolean = false) => {
        const response = await getRequest<User>("user", undefined, {
            credentials: 'include'
        });
        if (!response || response.code !== 1) {
            if (showToast) {
                console.log("获取用户信息失败");
                toast.add({
                    title: "当前用户未登录",
                    description: response?.msg,
                    icon: "i-fluent-error-circle-16-filled",
                    color: "red",
                    timeout: 2000,
                });
            }
            clearUserStatus();
            return false;
        }

        if (response && response.code === 1 && response.data) {
            const u: any = response.data as any
            user.value = {
                userid: u.id ?? u.ID ?? u.user_id,
                username: u.username ?? u.Username,
                is_admin: u.is_admin ?? u.IsAdmin,
                avatar_url: u.avatar_url ?? u.AvatarURL,
                description: u.description ?? u.Description
            } as any
            isLogin.value = true;
            await getStatus();
            return true;
        }
        return false;
    }

    // 退出登录
    const logout = async () => {
        const response = await postRequest("logout", {}, {
            credentials: 'include'
        });
        
        clearUserStatus();
        return true;
    }

    const checkLoginStatus = async () => {
        try {
            // 先尝试获取用户信息
            const userResult = await getUser();
            if (userResult) {
                return true;
            }
    
            // 如果获取用户信息失败，尝试获取状态
            const userStatus = await getStatus();
            const list = (userStatus as any)?.users || (userStatus as any)?.Users
            if (userStatus && Array.isArray(list)) {
                const currentUser = list.find((u: any) => (u.user_id ?? u.ID) === (user.value as any)?.userid)
                if (currentUser) {
                    user.value = {
                        userid: currentUser.user_id ?? currentUser.ID,
                        username: currentUser.username ?? currentUser.Username,
                        is_admin: currentUser.is_admin ?? currentUser.IsAdmin,
                        total_messages: (userStatus as any).total_messages ?? 0
                    } as any
                    isLogin.value = true
                    return true
                }
            }
    
            // 如果都失败了，清除状态
            clearUserStatus();
            return false;
        } catch (error) {
            console.error('检查登录状态失败:', error);
            clearUserStatus();
            return false;
        }
    }

    return {
        user,
        status,
        isLogin,
        token,
        register,
        login,
        getStatus,
        logout,
        getUser,
        setUserStatus,
        clearUserStatus,
        checkLoginStatus
    }
})
