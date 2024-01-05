export const initWebsocket = () => {
    const ws = new WebSocket('ws://localhost:8080/ws');
    ws.onopen = () => {
        console.log('connected');
    };
    ws.onMessage = (callback) => {
       ws.onmessage = (msg) => callback(JSON.parse(msg.data));
    }
    ws.sendMessage = (msg) => {
        ws.send(JSON.stringify(msg));
    }
    ws.onclose = () => {
        console.log('disconnected');
    };
    ws.onerror = (err) => {
        console.error(
            'Socket encountered error: ',
            err.message,
            'Closing socket'
        );
        ws.close();
    };
    window.$ws = ws;
}
