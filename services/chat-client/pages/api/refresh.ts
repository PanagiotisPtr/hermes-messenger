// Next.js API route support: https://nextjs.org/docs/api-routes/introduction
import type { NextApiRequest, NextApiResponse } from 'next'
import { parse } from "cookie"

type Data = {
    name: string
}

export default async function handler(
    req: NextApiRequest,
    res: NextApiResponse<Data>
) {
    const cookies = parse(req.headers.cookie ?? '')
    console.log(cookies.refreshToken)

    res.status(200).json({ name: "Refresh" })
}
