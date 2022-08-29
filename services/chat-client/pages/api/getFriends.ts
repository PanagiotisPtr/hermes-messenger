import { credentials } from '@grpc/grpc-js';
import type { NextApiResponse } from 'next'
import { NextApiRequestWithAuth, withAuth } from '../../auth-utils/token';
import { FriendsClient, GetFriendsResponse } from '../../grpc-clients/friends';
import { GetUserResponse, UserClient } from '../../grpc-clients/user';
import { defaultOptions } from '../../grpc-utils/options';

type Friend = {
    uuid: string;
    email: string;
}

async function handler(
    req: NextApiRequestWithAuth,
    res: NextApiResponse<{ friends: Friend[] } | { error: string }>
) {
    const userUuid = req.context?.userUuid ?? ""
    if (!userUuid) {
        res.status(401).json({ error: "Unauthorised" })
        return
    }

    const friendsClient = new FriendsClient(
        process.env.FRIENDS_SERVICE_ADDR ?? "",
        credentials.createInsecure(),
        defaultOptions,
    )
    const userClient = new UserClient(
        process.env.USER_SERVICE_ADDR ?? "",
        credentials.createInsecure(),
        defaultOptions,
    )

    const friendsResp = await new Promise<GetFriendsResponse>((res, rej) =>
        friendsClient.getFriends({
            UserUuid: userUuid,
        }, (err, resp) => err ? rej(err) : res(resp))
    )
    const friends: Friend[] = []
    for (const f of friendsResp.Friends) {
        if (f.Status != 'accepted') {
            continue
        }
        const userResp = await new Promise<GetUserResponse>((res, rej) =>
            userClient.getUser({
                Uuid: f.FriendUuid,
            }, (err, resp) => err ? rej(err) : res(resp))
        )
        const u = userResp.User
        if (!u) {
            continue
        }
        friends.push({ uuid: u.Uuid, email: u.Email })
    }

    res.status(200).json({ friends })
}

export default withAuth(handler)