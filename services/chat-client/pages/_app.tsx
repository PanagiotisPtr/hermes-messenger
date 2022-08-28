import '../styles/globals.css'
import type { AppProps } from 'next/app'
import { autoRefresh } from '../auth-utils/autoRefresh'
import { useEffect } from 'react';

function MyApp({ Component, pageProps }: AppProps) {
    let isRunning = false
    useEffect(() => {
        if (isRunning) {
            return
        }
        isRunning = true
        autoRefresh()

        const socket = new WebSocket('ws://localhost:3000/api/ws')
        socket.onopen = (e) => console.log('CONNECTED TO WEBSOCKETS!')
        socket.onclose = (e) => console.log('DISCONNECTED FROM WEBSOCKET')
        socket.onmessage = (e) => {
            const data = JSON.parse(e.data)
            if (data.type == "message") {
                const event = new CustomEvent("message-event", { detail: data.data })
                window.dispatchEvent(event)
            }
        }
    }, []);
    return <Component {...pageProps} />
}

export default MyApp
