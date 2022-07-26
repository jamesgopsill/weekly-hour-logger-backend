import jwt from "jsonwebtoken"
import { UserScopes } from "../src"

const secret = "test"

const adminDetails = {
	name: "admin",
	email: "admin@test.com",
	scopes: [UserScopes.ADMIN, UserScopes.USER],
	iat: Math.floor(Date.now() / 1000) + 60 * 60 * 24,
}

export const validAdminToken = jwt.sign(adminDetails, secret, {
	algorithm: "HS256",
})

const userDetails = {
	name: "user",
	email: "user@test.com",
	scopes: [UserScopes.USER],
	iat: Math.floor(Date.now() / 1000) + 60 * 60 * 24,
}

export const validUserToken = jwt.sign(userDetails, secret, {
	algorithm: "HS256",
})
