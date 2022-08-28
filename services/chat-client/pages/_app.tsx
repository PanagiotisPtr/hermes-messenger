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
    }, []);
    return <Component {...pageProps} />
}

export default MyApp
