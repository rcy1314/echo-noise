import type { Response } from "~/types/models";
import { useUserStore } from "~/store/user";
import { useToast } from "#imports";

export const postRequest = async <T>(url: string, body: object | FormData, options?: { credentials?: RequestCredentials }) => {
    const BASE_API = useRuntimeConfig().public.baseApi || '/api';
    const userStore = useUserStore();
    const token = userStore.token || "null";

    try {
        const isFormData = body instanceof FormData;
        const headers = isFormData 
        ? { 'Authorization': `Bearer ${token}` }
        : { 
            'Content-Type': 'application/json',
            'Authorization': `Bearer ${token}`
          };

        const response: Response<T> = await $fetch(`${BASE_API}/${url}`, {
            method: 'POST',
            headers,
            body: isFormData ? body : JSON.stringify(body),
            credentials: options?.credentials,
            timeout: 8000,
            retry: 0
        });

        return response;
    } catch (error) {
        const toast = useToast();
        toast.add({ title: '请求失败', description: '网络异常或服务器不可用', color: 'red', timeout: 2000 });
        return { code: 0, msg: '网络异常', data: null } as any as Response<T>;
    }
};

export const getRequest = async <T>(url: string, params?: any, options?: { credentials?: RequestCredentials }) => {
    const BASE_API = useRuntimeConfig().public.baseApi || '/api';
    const userStore = useUserStore();
    const token = userStore.token || "null";

    try {
        const queryParamString = params ? "?" + Object.keys(params).map(key => key + "=" + params[key]).join("&") : "";
        
        const response: Response<T> = await $fetch(`${BASE_API}/${url}${queryParamString}`, {
            method: 'GET',
            headers: {
                'Authorization': `Bearer ${token}`,
                'Cache-Control': 'no-cache',
                'Pragma': 'no-cache'
            },
            credentials: options?.credentials,
            timeout: 8000,
            retry: 0
        });

        return response;
    } catch (error) {
        const e: any = error;
        const status = e?.response?.status || e?.status;
        if (status === 401) {
            const msg = e?.response?._data?.msg || e?.response?.statusText || '未登录或登录已过期';
            return { code: 0, msg } as any as Response<T>;
        }
        const toast = useToast();
        toast.add({ title: '请求失败', description: '网络异常或服务器不可用', color: 'red', timeout: 2000 });
        return { code: 0, msg: '网络异常', data: null } as any as Response<T>;
    }
};

export const putRequest = async <T>(url: string, body: object, options?: { credentials?: RequestCredentials }) => {
    const BASE_API = useRuntimeConfig().public.baseApi || '/api';
    const toast = useToast();
    const userStore = useUserStore();
    const token = userStore.token || "null";

    try {
        const response: Response<T> = await $fetch(`${BASE_API}/${url}`, {
            method: 'PUT',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${token}`,
            },
            body: JSON.stringify(body),
            credentials: options?.credentials,
            timeout: 8000,
            retry: 0
        });

        if (response.code !== 1) {
            toast.add({
                title: "请求失败",
                description: response.msg,
                icon: "i-fluent-error-circle-16-filled",
                color: "red",
                timeout: 2000,
            });
            return null;
        }

        return response;
    } catch (error) {
        const toast = useToast();
        toast.add({ title: '请求失败', description: '网络异常或服务器不可用', color: 'red', timeout: 2000 });
        return { code: 0, msg: '网络异常', data: null } as any as Response<T>;
    }
};

export const deleteRequest = async <T>(url: string, params?: any, options?: { credentials?: RequestCredentials }) => {
    const BASE_API = useRuntimeConfig().public.baseApi || '/api';
    const userStore = useUserStore();
    const token = userStore.token || "null";

    try {
        const queryParamString = params ? "?" + Object.keys(params).map(key => key + "=" + params[key]).join("&") : "";
        
        const response: Response<T> = await $fetch(`${BASE_API}/${url}${queryParamString}`, {
            method: 'DELETE',
            headers: {
                'Authorization': `Bearer ${token}`,
            },
            credentials: options?.credentials
        });

        return response;
    } catch (error) {
        console.error('请求失败:', error);
        throw error;
    }
};
