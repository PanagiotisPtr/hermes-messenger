import { parse } from "cookie";
import { decode, verify } from "jsonwebtoken";
import { NextApiHandler, NextApiRequest, NextApiResponse } from "next"
import { getPublicKeys } from "./publicKeys";

export async function tokenIsValid(): Promise<boolean> {
    return true
}

export interface AuthenticationContext {
    userUuid: string
}

interface NextApiRequestWithAuth extends NextApiRequest {
    context?: AuthenticationContext
}

type NextApiHandlerWithAuth<T> = (req: NextApiRequestWithAuth, res: NextApiResponse<T>) => unknown | Promise<unknown>;

export function withAuth<T>(handler: NextApiHandlerWithAuth<T>): NextApiHandler<T | { error: string }> {
    return async function (req: NextApiRequest, res: NextApiResponse<T | { error: string }>) {

        const cookies = parse(req.headers.cookie ?? '')
        const { accessToken } = cookies

        if (!accessToken) {
            res.status(401).json({ error: "Unauthorised" })
            return
        }

        const headers = decode(accessToken)
        if (typeof headers == 'string' || headers == null) {
            res.status(401).json({ error: "Malformed authentication token" })
            return
        }
        const { dat, exp } = headers
        if (!dat || !exp) {
            res.status(401).json({ error: "Malformed authentication token" })
            return
        }
        if (exp * 1000 < Date.now()) {
            res.status(401).json({ error: "Authentication token has expired" })
            return
        }
        const keys = await getPublicKeys()
        let isVerified = false
        for (const key of keys) {
            const formattedKey = key.replaceAll(" RSA ", " ")
            const result = await new Promise<boolean>((res, _) => {
                verify(accessToken, formattedKey, (err, _) => {
                    if (err) {
                        res(false)
                    }
                    res(true)
                })
            })
            if (result) {
                isVerified = true
                break
            }
        }
        if (!isVerified) {
            res.status(401).json({ error: "Invalid token" })
            return
        }

        const reqWithAuth: NextApiRequestWithAuth = req
        reqWithAuth.context = {
            userUuid: dat,
        }

        return handler(req, res)
    }
}