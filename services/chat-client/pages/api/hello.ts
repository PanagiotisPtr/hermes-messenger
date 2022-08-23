// Next.js API route support: https://nextjs.org/docs/api-routes/introduction
import type { NextApiRequest, NextApiResponse } from 'next'
import { credentials, ServiceError } from '@grpc/grpc-js';
import { AuthenticationClient, GetPublicKeysResponse } from '../../grpc-clients/authentication'

type Data = {
  name: string
}

export default async function handler(
  req: NextApiRequest,
  res: NextApiResponse<Data>
) {
  const service = new AuthenticationClient('localhost:7777', credentials.createInsecure(), {
    'grpc.keepalive_time_ms': 120000,
    'grpc.http2.min_time_between_pings_ms': 120000,
    'grpc.keepalive_timeout_ms': 20000,
    'grpc.http2.max_pings_without_data': 0,
    'grpc.keepalive_permit_without_calls': 1,
  });

  service.getPublicKeys({}, (err: ServiceError | null, res: GetPublicKeysResponse) => {
    if (err) {
      console.log('Error: ', err)
    }

    console.log(res.PublicKeys)
  })
  res.status(200).json({ name: 'John Doe' })
}