import { Request } from "express"
import { UserScopes } from "../entities/user"
import jwt from "jsonwebtoken"

if (!process.env.SECRET) {
	console.warn("Using default secret. Should only be used for testing.")
}
const secret = process.env.SECRET || "test"

export interface UserToken {
	name: string
	email: string
	scopes: UserScopes[]
	iat: number
}

export const verifyToken = (
	req: Request
): {
	token: UserToken
	error: string | null
} => {
	if (!req.headers["authorization"]) {
		return {
			token: {
				name: "",
				email: "",
				scopes: [],
				iat: 0,
			},
			error: "No Authorization Header",
		}
	}

	if (!req.headers["authorization"].includes("Bearer")) {
		return {
			token: {
				name: "",
				email: "",
				scopes: [],
				iat: 0,
			},
			error: "Token formatting error",
		}
	}

	const token = req.headers["authorization"].split(" ").pop()

	if (!token) {
		return {
			token: {
				name: "",
				email: "",
				scopes: [],
				iat: 0,
			},
			error: "Token formatting error",
		}
	}

	try {
		const decoded = jwt.verify(token, secret)
		const decodedToken = decoded as UserToken
		return {
			token: decodedToken,
			error: null,
		}
	} catch (e: any) {
		return {
			token: {
				name: "",
				email: "",
				scopes: [],
				iat: 0,
			},
			error: e.message,
		}
	}
}
