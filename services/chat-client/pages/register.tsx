import type { NextPage } from "next"
import Router from "next/router"
import { useState } from "react"
import styles from "../styles/Login.module.css"

const Register: NextPage = () => {
    const [email, setEmail] = useState("")
    const [password, setPassword] = useState("")
    const [verifyPassword, setVerifyPassword] = useState("")
    const [displayMessage, setDisplayMessage] = useState("")

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

        if (resp.success) {
            setDisplayMessage("Successfully registered user! You will be redirected to the login page...")
            setTimeout(() => {
                setDisplayMessage("")
                Router.push("/login")
            }, 2500)
        }
    }

    return (
        <div className={styles.mainContainer}>
            <div className={styles.colContainer}>
                {displayMessage && <span><b>{displayMessage}</b></span>}
                <input type="email" onChange={e => setEmail(e.target.value)} placeholder="Email" />
                <input type="password" onChange={e => setPassword(e.target.value)} placeholder="Password" />
                <input type="password" onChange={e => setVerifyPassword(e.target.value)} placeholder="Verify password" />
                <button type="submit" onClick={register}>Register</button>
            </div>
        </div>
    )
}

export default Register
