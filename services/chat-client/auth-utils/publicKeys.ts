import { credentials } from "@grpc/grpc-js"
import { AuthenticationClient, GetPublicKeysResponse } from "../grpc-clients/authentication"
import { defaultOptions } from "../grpc-utils/options"

const service = new AuthenticationClient(
    process.env.AUTHENTICATION_SERVICE_ADDR ?? "",
    credentials.createInsecure(),
    defaultOptions,
)

let keys: string[] = []

export async function getPublicKeys(): Promise<string[]> {
    if (keys.length) {
        return keys
    }

    const response = await new Promise<GetPublicKeysResponse>((res, rej) =>
        service.getPublicKeys({}, (err, resp) => err ? rej(err) : res(resp))
    )

    // todo - add KID
    // todo - set expiry on keys when service exits
    keys = response.PublicKeys

    return keys
}