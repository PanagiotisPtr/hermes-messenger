import { credentials } from "@grpc/grpc-js"
import { AuthenticationClient, GetPublicKeysResponse } from "../grpc-clients/authentication"
import { defaultOptions } from "../grpc-utils/options"

type KeyMap = { [key: string]: string }

let keys: KeyMap = {}

export async function getPublicKeys(): Promise<KeyMap> {
    if (Object.keys(keys).length) {
        return keys
    }

    const service = new AuthenticationClient(
        process.env.AUTHENTICATION_SERVICE_ADDR ?? "",
        credentials.createInsecure(),
        defaultOptions,
    )

    const response = await new Promise<GetPublicKeysResponse>((res, rej) =>
        service.getPublicKeys({}, (err, resp) => err ? rej(err) : res(resp))
    )

    for (const key of response.PublicKeys) {
        keys[key.Uuid] = key.Value
    }

    return keys
}