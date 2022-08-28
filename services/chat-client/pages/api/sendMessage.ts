import { credentials } from "@grpc/grpc-js";
import type { NextApiResponse } from "next"
import { NextApiRequestWithAuth, withAuth } from "../../auth-utils/token";
import { MessagingClient, SendMessageResponse } from "../../grpc-clients/messaging";
import { defaultOptions } from "../../grpc-utils/options";

type SendMessageRequest = {
    to: string
    content: string
}

function isSendMessageRequest(req: any): req is SendMessageRequest {
    return req && req.to && req.content
}

async function handler(
    req: NextApiRequestWithAuth,
    res: NextApiResponse<{ success: boolean } | { error: string }>
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
        const response = await new Promise<SendMessageResponse>((res, rej) =>
            messagingClient.sendMessage({
                From: userUuid,
                To: body.to,
                Content: body.content,
            }, (err, resp) => err ? rej(err) : res(resp))
        )
        res.status(200).json({ success: true })
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
