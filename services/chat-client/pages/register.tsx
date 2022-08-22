import type { NextPage } from "next"
import { useState } from "react"
import styles from "../styles/Login.module.css"

const Register: NextPage = () => {
    const [email, setEmail] = useState("")
    const [password, setPassword] = useState("")
    const [verifyPassword, setVerifyPassword] = useState("")
    const register = async () => {
        if (password != verifyPassword) {
            console.warn("Passwords do not much")
            return
        }
        const resp = await fetch("/api/register", {
            method: "POST",
            body: JSON.stringify({
                email,
                password
            })
        }).then(res => res.json())

        console.log(resp)
    }

    return (
        <div className={styles.mainContainer}>
            <div className={styles.colContainer}>
                <input type="email" onChange={e => setEmail(e.target.value)} placeholder="Email" />
                <input type="password" onChange={e => setPassword(e.target.value)} placeholder="Password" />
                <input type="password" onChange={e => setVerifyPassword(e.target.value)} placeholder="Verify password" />
                <button type="submit" onClick={register}>Register</button>
            </div>
        </div>
    )
}

export default Register
