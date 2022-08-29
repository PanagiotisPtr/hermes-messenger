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

export async function tokenIsValid(token: string): Promise<[string, Error | null]> {
    if (!token) {
        return ["", new Error("No token provided")]
    }

    const headers = decode(token)
    if (typeof headers == 'string' || headers == null) {
        return ["", new Error("Malformed authentication token")]
    }
    const { dat, exp, kid } = headers
    if (!dat || !exp || !kid) {
        return ["", new Error("Malformed authentication token")]
    }
    if (exp * 1000 < Date.now()) {
        return ["", new Error("Authentication token has expired")]
    }
    const keys = await getPublicKeys()
    const key = keys[kid]
    const formattedKey = key.replaceAll(" RSA ", " ")
    const result = await new Promise<boolean>((res, _) => {
        verify(token, formattedKey, (err, _) => {
            if (err) {
                res(false)
            }
            res(true)
        })
    })
    if (!result) {
        return ["", new Error("Invalid token")]
    }

    return [dat, null]
}

export function withAuth<T>(handler: NextApiHandlerWithAuth<T>): NextApiHandler<T | { error: string }> {
    return async function (req: NextApiRequest, res: NextApiResponse<T | { error: string }>) {
        const cookies = parse(req.headers.cookie ?? '')
        const { accessToken } = cookies

        const [dat, err] = await tokenIsValid(accessToken)
        if (err != null) {
            res.status(401).json({ error: err.message })
            return
        }

        const reqWithAuth: NextApiRequestWithAuth = req
        reqWithAuth.context = {
            userUuid: dat,
        }

        return handler(req, res)
    }
}