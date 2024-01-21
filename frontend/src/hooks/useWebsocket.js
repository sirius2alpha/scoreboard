import { ElMessage } from 'element-plus';
import { ref, onMounted, onUnmounted } from 'vue';

export default function useWebSocket(url) {
    const data = ref(null);
    let socket = null;

    const send = (message) => {
        socket.send(JSON.stringify(message));
    };

    onMounted(() => {
        socket = new WebSocket(url);

        socket.onmessage = (event) => {
            data.value = event.data;
        };

        socket.onopen = (event) => {
            ElMessage.success('服务已连接');
        };

        socket.onerror = (error) => {
            console.log("WebSocket error", error);
        };

        socket.onclose = (event) => {
            console.log("Connection closed", event);
        };
    });

    onUnmounted(() => {
        if (socket) {
            socket.close();
        }
    });

    return { data, send };
}
