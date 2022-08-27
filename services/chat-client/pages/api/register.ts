import { credentials, ServiceError } from "@grpc/grpc-js";
import type { NextApiRequest, NextApiResponse } from "next"
import { AuthenticationClient, RegisterRequest, RegisterResponse } from "../../grpc-clients/authentication";
import { defaultOptions } from "../../grpc-utils/options";

interface APIRegisterRequest {
    email: string;
    password: string;
}

function isRegisterRequest(req: any): req is APIRegisterRequest {
    return req && req.email && req.password
}

export default async function handler(
    req: NextApiRequest,
    res: NextApiResponse<any>
) {
    const body = JSON.parse(req.body)
    if (!isRegisterRequest(body)) {
        res.status(500).json({ error: "Invalid payload" })
        return
    }

    const service = new AuthenticationClient(
        process.env.AUTHENTICATION_SERVICE_ADDR ?? "",
        credentials.createInsecure(),
        defaultOptions,
    )

    try {
        const response = await new Promise<RegisterResponse>((res, rej) =>
            service.register({
                Email: body.email,
                Password: body.password,
            }, (err, resp) => err ? rej(err) : res(resp))
        )
        res.status(200).json(response)
    } catch (err: any) {
        res.status(500).json(err)
    }
}
