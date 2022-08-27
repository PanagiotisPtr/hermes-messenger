export async function autoRefresh() {
    setInterval(async () => {
        const last = localStorage.getItem('lastRefresh')
        const oneHour = 1000 * 60 * 60
        if (last == null || (Date.now() - Number(last)) >= oneHour) {
            await fetch("/api/refresh", {})
            localStorage.setItem('lastRefresh', '' + Date.now())
        }
    }, 5000)
}