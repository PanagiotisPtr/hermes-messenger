import type { NextPage } from "next"
import { useEffect, useState } from "react"
import styles from "../styles/Login.module.css"
import { useRouter } from "next/router"
import Link from "next/link"

interface entity {
    uuid: string
    name: string
}

interface Props {
    friendUuid: string
}

const Chat: NextPage<Props> = ({ friendUuid }) => {
    const [friends, setFriends] = useState<entity[]>([])
    const [messages, setMessages] = useState<string[]>([])
    const [message, setMessage] = useState("")
    const [friend, setFriend] = useState<entity>({ uuid: friendUuid, name: "" })

    useEffect(() => {
        setFriends([
            { uuid: "1", name: "Bob" },
            { uuid: "2", name: "Alice" },
        ])

        setMessages([
            "Hello",
            "How are you",
        ])
    }, [])

    useEffect(() => {
        const f = friends.find(f => f.uuid === friendUuid)
        if (f) {
            setFriend(f)
        }
    }, [friends, friendUuid])

    return (
        <div className={styles.mainContainer}>
            <div className={styles.colContainer}>
                {friends.map((friend, i) => <span key={i}>
                    <Link href={{ pathname: "chat", query: { uuid: friend.uuid } }}>{friend.name}</Link>
                </span>)}
            </div>
            {friend.name &&
                <div className={styles.colContainer}>
                    <span><b>{friend.name}</b></span>
                    {messages.map((m, i) => <span key={i}>{m}</span>)}
                    <form onSubmit={e => { setMessages([...messages, message]); setMessage(""); e.preventDefault() }}>
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
