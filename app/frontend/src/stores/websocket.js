import {defineStore} from 'pinia';
import {computed, ref} from 'vue';
import {ElMessage} from 'element-plus';

export const useWebsocket = defineStore('websocket', () => {
    const data = ref(null);
    const instance = new WebSocket('ws://127.0.0.1:8001/ws');
    // const instance = new WebSocket('ws://1.94.126.139:8080/ws');

    const send = (message) => {
        instance.send(JSON.stringify(message));
    };

    instance.onmessage = (event) => {
        data.value = event.data;
    };

    instance.onopen = (event) => {
        ElMessage.success('服务已连接');
    };

    instance.onerror = (error) => {
        console.log('WebSocket error', error);
    };

    instance.onclose = (event) => {
        ElMessage.error('服务已断开');
        console.log('Connection closed', event);
    };
    return {instance,data,send};
});
