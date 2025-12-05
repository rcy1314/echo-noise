import type { MessageToSave, Message } from "~/types/models"
import { useMessageStore } from "~/store/message"

export const useMessage = () => {
    const message = useMessageStore()
    const toast = useToast()

    // save
    const save = async (messageToSave: MessageToSave) => {
        try {
            const response = await postRequest<Message>('messages', messageToSave, { credentials: 'include' });
            if (!response || response.code !== 1) {
                console.log(response?.code);
                toast.add({
                    title: "添加失败",
                    description: String(response?.msg),
                    icon: "i-fluent-error-circle-16-filled",
                    color: "red",
                    timeout: 2000,
                });
                return null;
            }

            // 更新笔记列表
            const newMessage = response.data;
            message.messages.unshift(newMessage); // 将新消息添加到列表的开头

            // 提示成功
            toast.add({
                title: "保存成功",
                description: "消息已保存",
                icon: "i-fluent-checkmark-starburst-16-filled",
                color: "green",
                timeout: 1000,
            })

            return response.data;
        } catch (error) {
            console.error(error);
        }
    }

    const deleteMessage = async (id: number) => {
        try {
            const response = await message.deleteMessage(id);
            if (response && response.code === 1) {
                toast.add({
                    title: "删除成功",
                    description: "笔记已删除",
                    icon: "i-fluent-checkmark-starburst-16-filled",
                    color: "green",
                    timeout: 1000,
                });
                return true;
            }
        } catch (error) {
            console.error(error);
        }
    }

    // 获取笔记列表
    // const getAllMessages = async () => {
    //     try {
    //         const response = await getRequest<Message[]>('messages');
    //         if (!response) {
    //             throw new Error('Get messages request failed');
    //         }

    //         return response.data;
    //     } catch (error) {
    //         console.error(error);
    //     }
    // }

    // 图片上传
    const uploadImage = async (file: File): Promise<string | null> => {
        try {
            const formData = new FormData();
            formData.append('image', file);

            const response = await postRequest<string>('images/upload', formData, { credentials: 'include' });

            if (!response || response.code !== 1) {
                console.log(response?.code);
                toast.add({
                    title: "上传失败",
                    description: String(response?.msg),
                    icon: "i-fluent-error-circle-16-filled",
                    color: "red",
                    timeout: 2000,
                });
                return null;
            }

            return response.data; // 返回图片的URL
        } catch (error) {
            console.error('上传图片错误:', error);
            return null; // 出现错误时返回null
        }
    };




    return {
        save, uploadImage, deleteMessage,
    }
}
