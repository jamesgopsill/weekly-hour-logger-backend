import { Request } from "express"
import { UserScopes } from "../entities/user"

export interface UserToken {
	name: string
	email: string
	scopes: UserScopes[]
	iat: number
}

export interface RequestWithToken extends Request {
	token: UserToken
}
