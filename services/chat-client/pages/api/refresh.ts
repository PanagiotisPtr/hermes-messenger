import type { NextApiRequest, NextApiResponse } from 'next'
import { parse, serialize } from "cookie"
import { AuthenticationClient, RefreshResponse } from '../../grpc-clients/authentication'
import { credentials } from '@grpc/grpc-js'
import { defaultOptions } from '../../grpc-utils/options'

export default async function handler(
    req: NextApiRequest,
    res: NextApiResponse<any>
) {
    const cookies = parse(req.headers.cookie ?? '')
    const { refreshToken } = cookies

    if (!refreshToken) {
        res.status(401).json({ error: "Unauthorised" })
        return
    }

    const service = new AuthenticationClient(
        process.env.AUTHENTICATION_SERVICE_ADDR ?? "",
        credentials.createInsecure(),
        defaultOptions,
    )

    try {
        const response = await new Promise<RefreshResponse>((res, rej) =>
            service.refresh({
                RefreshToken: refreshToken,
            }, (err, resp) => err ? rej(err) : res(resp))
        )

        const accessExpiry = new Date()
        accessExpiry.setHours(accessExpiry.getHours() + 1)
        res.setHeader(
            "Set-Cookie",
            [
                serialize(
                    "accessToken",
                    response.AccessToken,
                    {
                        httpOnly: true,
                        secure: true,
                        path: "/",
                        expires: accessExpiry,
                    }
                ),
            ]
        )

        res.status(200).json({
            success: true,
        })
    } catch (err: any) {
        res.status(500).json(err)
    }
}
