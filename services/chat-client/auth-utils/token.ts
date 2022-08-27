import { parse } from "cookie";
import { decode, verify } from "jsonwebtoken";
import { NextApiHandler, NextApiRequest, NextApiResponse } from "next"
import { getPublicKeys } from "./publicKeys";

export interface AuthenticationContext {
    userUuid: string
}

export interface NextApiRequestWithAuth extends NextApiRequest {
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
        const { dat, exp, kid } = headers
        if (!dat || !exp || !kid) {
            res.status(401).json({ error: "Malformed authentication token" })
            return
        }
        if (exp * 1000 < Date.now()) {
            res.status(401).json({ error: "Authentication token has expired" })
            return
        }
        const keys = await getPublicKeys()
        const key = keys[kid]
        const formattedKey = key.replaceAll(" RSA ", " ")
        const result = await new Promise<boolean>((res, _) => {
            verify(accessToken, formattedKey, (err, _) => {
                if (err) {
                    res(false)
                }
                res(true)
            })
        })
        if (!result) {
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