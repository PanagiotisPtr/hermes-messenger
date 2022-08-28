import type { NextPage } from "next"
import { useEffect, useState } from "react"
import styles from "../styles/Login.module.css"
import Link from "next/link"

interface entity {
    uuid: string
    email: string
}

interface Props {
    friendUuid: string
}

const Chat: NextPage<Props> = ({ friendUuid }) => {
    const [friends, setFriends] = useState<entity[]>([])
    const [messages, setMessages] = useState<string[]>([])
    const [message, setMessage] = useState("")
    const [friend, setFriend] = useState<entity>({ uuid: friendUuid, email: "" })
    const [addFriendField, setAddFriendField] = useState("")

    useEffect(() => {
        fetch('/api/getFriends')
            .then(res => res.json())
            .then(res => { setFriends(res.friends); console.log(friends); })

        if (friend.uuid) {
            fetch("/api/getMessages", {
                method: "POST",
                body: JSON.stringify({
                    to: friend.uuid,
                })
            }).then(res => res.json()).then((res) => {
                if (res.messages) {
                    setMessages(res.messages.map((m: any) => m.content))
                }
            })
        }
    }, [])

    useEffect(() => {
        const f = friends.find(f => f.uuid === friendUuid)
        if (f) {
            setFriend(f)
        }
    }, [friends, friendUuid])

    const addFriend = async (email: string) => {
        const resp = await fetch("/api/addFriend", {
            method: "POST",
            body: JSON.stringify({
                friendEmail: email,
            })
        }).then(res => res.json())

        if (resp.error) {
            console.error(resp.error)
        }
    }

    const sendMessage = async (message: string) => {
        const resp = await fetch("/api/sendMessage", {
            method: "POST",
            body: JSON.stringify({
                to: friend.uuid,
                content: message,
            })
        }).then(res => res.json())

        if (resp.error) {
            console.error(resp.error)
        }
    }

    return (
        <div className={styles.mainContainer}>
            <div className={styles.colContainer}>
                {friends.map((friend, i) => <span key={i}>
                    <Link href={{ pathname: "chat", query: { uuid: friend.uuid } }}>{friend.email}</Link>
                </span>)}
                <form onSubmit={e => { addFriend(addFriendField); e.preventDefault() }}>
                    <input type="text" onChange={e => setAddFriendField(e.target.value)} placeholder="friend's email" />
                    <button type="submit">Add</button>
                </form>
            </div>
            {friend.email &&
                <div className={styles.colContainer}>
                    <span><b>{friend.email}</b></span>
                    {messages.map((m, i) => <span key={i}>{m}</span>)}
                    <form onSubmit={e => { setMessages([...messages, message]); sendMessage(message); setMessage(""); e.preventDefault() }}>
                        <input type="text" value={message} onChange={e => setMessage(e.target.value)} placeholder="Say something" />
                    </form>
                </div>
            }
        </div>
    )
}

Chat.getInitialProps = async ({ query }) => {
    const uuid = "" + query?.uuid

    return { friendUuid: uuid }
}

export default Chat
