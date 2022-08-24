import { credentials } from '@grpc/grpc-js';
import type { NextApiRequest, NextApiResponse } from 'next'
import { serialize } from "cookie"
import { AuthenticateResponse, AuthenticationClient } from '../../grpc-clients/authentication';
import { defaultOptions } from '../../grpc-utils/options';

interface APILoginRequest {
  email: string;
  password: string;
}

function isLoginRequest(req: any): req is APILoginRequest {
  return req && req.email && req.password
}

export default async function handler(
  req: NextApiRequest,
  res: NextApiResponse<any>
) {
  const body = JSON.parse(req.body)
  if (!isLoginRequest) {
    res.status(500).json({ error: "Invalid payload" })
    return
  }

  const service = new AuthenticationClient(
    process.env.AUTHENTICATION_SERVICE_ADDR ?? "",
    credentials.createInsecure(),
    defaultOptions,
  )

  try {
    const response = await new Promise<AuthenticateResponse>((res, rej) =>
      service.authenticate({
        Email: body.email,
        Password: body.password,
      }, (err, resp) => err ? rej(err) : res(resp))
    )

    const refreshExpiry = new Date()
    refreshExpiry.setDate(refreshExpiry.getDate() + 1)

    const accessExpiry = new Date()
    accessExpiry.setHours(accessExpiry.getHours() + 1)

    res.setHeader(
      "Set-Cookie",
      [
        serialize(
          "refreshToken",
          response.RefreshToken,
          {
            httpOnly: true,
            secure: true,
            path: "/api/refresh",
            expires: refreshExpiry,
          }
        ),
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
