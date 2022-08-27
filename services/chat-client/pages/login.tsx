import type { NextPage } from "next"
import { useState } from "react"
import styles from "../styles/Login.module.css"

const Login: NextPage = () => {
    const [email, setEmail] = useState("")
    const [password, setPassword] = useState("")
    const login = async () => {
        const resp = await fetch("/api/login", {
            method: "POST",
            body: JSON.stringify({
                email,
                password
            })
        }).then(res => res.json())
    }

    return (
        <div className={styles.mainContainer}>
            <div className={styles.colContainer}>
                <input type="email" onChange={e => setEmail(e.target.value)} placeholder="Email" />
                <input type="password" onChange={e => setPassword(e.target.value)} placeholder="Password" />
                <button type="submit" onClick={login}>Login</button>
            </div>
        </div>
    )
}

export default Login
