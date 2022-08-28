import { credentials } from "@grpc/grpc-js";
import type { NextApiResponse } from "next"
import { NextApiRequestWithAuth, withAuth } from "../../auth-utils/token";
import { AddFriendResponse, FriendsClient } from "../../grpc-clients/friends";
import { GetUserByEmailResponse, UserClient } from "../../grpc-clients/user";
import { defaultOptions } from "../../grpc-utils/options";

type AddFriendRequest = {
    friendEmail: string
}

function isAddFriendRequest(req: any): req is AddFriendRequest {
    return req && req.friendEmail
}

async function handler(
    req: NextApiRequestWithAuth,
    res: NextApiResponse<{ success: boolean } | { error: string }>
) {
    const userUuid = req.context?.userUuid ?? ""
    if (!userUuid) {
        res.status(401).json({ error: "Unauthorised" })
        return
    }

    const body = JSON.parse(req.body)
    if (!isAddFriendRequest(body)) {
        res.status(500).json({ error: "Invalid payload" })
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

    const userResp = await new Promise<GetUserByEmailResponse>((res, rej) =>
        userClient.getUserByEmail({
            Email: body.friendEmail,
        }, (err, resp) => err ? rej(err) : res(resp))
    )
    const friend = userResp.User
    if (!friend) {
        res.status(500).json({ error: `Could not find friend with email ${body.friendEmail}` })
        return
    }
    const friendsResp = await new Promise<AddFriendResponse>((res, rej) =>
        friendsClient.addFriend({
            UserUuid: userUuid,
            FriendUuid: friend.Uuid
        }, (err, resp) => err ? rej(err) : res(resp))
    )

    res.status(200).json({ success: true })
}

export default withAuth(handler)
