import { credentials } from "@grpc/grpc-js";
import type { NextApiResponse } from "next"
import { NextApiRequestWithAuth, withAuth } from "../../auth-utils/token";
import { GetMessagesResponse, MessagingClient } from "../../grpc-clients/messaging";
import { defaultOptions } from "../../grpc-utils/options";

type SendMessageRequest = {
    to: string
}

interface Message {
    from: string
    to: string
    content: string
    timestamp: Date
}

function isSendMessageRequest(req: any): req is SendMessageRequest {
    return req && req.to
}

async function handler(
    req: NextApiRequestWithAuth,
    res: NextApiResponse<{ messages: Message[] } | { error: string }>
) {
    const userUuid = req.context?.userUuid ?? ""
    if (!userUuid) {
        res.status(401).json({ error: "Unauthorised" })
        return
    }

    const body = JSON.parse(req.body)
    if (!isSendMessageRequest(body)) {
        res.status(500).json({ error: "Invalid payload" })
        return
    }

    const messagingClient = new MessagingClient(
        process.env.MESSAGING_SERVICE_ADDR ?? "",
        credentials.createInsecure(),
        defaultOptions,
    )

    try {
        const response = await new Promise<GetMessagesResponse>((res, rej) =>
            messagingClient.getMessages({
                From: userUuid,
                To: body.to,
                Start: 0,
                End: Math.floor(Date.now() / 1000),
            }, (err, resp) => err ? rej(err) : res(resp))
        )

        const messages: Message[] = []
        for (const m of response.Messages) {
            messages.push({
                from: m.From,
                to: m.To,
                content: m.Content,
                timestamp: new Date(m.Timestamp * 1000)
            })
        }
        res.status(200).json({ messages })
    } catch (error) {
        if (error instanceof Error) {
            res.status(500).json({ error: error.message })
        } else {
            console.error(error)
            res.status(500).json({ error: "internal server error" })
        }
    }
}

export default withAuth(handler)