import { credentials } from '@grpc/grpc-js';
import type { NextApiResponse } from 'next'
import { NextApiRequestWithAuth, withAuth } from '../../auth-utils/token';
import { GetUserResponse, UserClient } from '../../grpc-clients/user';
import { defaultOptions } from '../../grpc-utils/options';

type Data = {
  name: string
}

async function handler(
  req: NextApiRequestWithAuth,
  res: NextApiResponse<Data | { error: string }>
) {
  const userUuid = req.context?.userUuid ?? ""
  if (!userUuid) {
    res.status(401).json({ error: "Unauthorised" })
    return
  }

  const service = new UserClient(
    process.env.USER_SERVICE_ADDR ?? "",
    credentials.createInsecure(),
    defaultOptions,
  )

  const response = await new Promise<GetUserResponse>((res, rej) =>
    service.getUser({
      Uuid: userUuid,
    }, (err, resp) => err ? rej(err) : res(resp))
  )
  const user = response.User
  if (!user) {
    res.status(401).json({ error: "Unauthorised" })
    return
  }

  res.status(200).json({ name: user.Email })
}

export default withAuth(handler)