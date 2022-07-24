import { Request, Response, NextFunction } from "express"
import { UserToken } from "./interfaces"
import { ResponseFormat } from "../routers"
import jwt from "jsonwebtoken"
import { UserScopes } from "../entities"

if (!process.env.SECRET) {
	console.warn("Using default secret. Should only be used for testing.")
}
const secret = process.env.SECRET || "test"

export const authorize = (scopes: UserScopes[]) => {
	return (req: Request, res: Response, next: NextFunction) => {
		if (!req.headers["authorization"]) {
			const json: ResponseFormat = {
				error: "You do not have the permissions to perform this operation",
				data: null,
			}
			return res.status(400).json(json)
		}

		if (!req.headers["authorization"].includes("Bearer")) {
			const json: ResponseFormat = {
				error: "Incorrect token format",
				data: null,
			}
			return res.status(400).json(json)
		}

		const token = req.headers["authorization"].split(" ").pop()

		if (!token) {
			const json: ResponseFormat = {
				error: "Incorrect token format",
				data: null,
			}
			return res.status(400).json(json)
		}

		try {
			const decoded = jwt.verify(token, secret)
			const decodedToken = decoded as UserToken
			let permitted = false
			for (const scope of scopes) {
				if (decodedToken.scopes.includes(scope)) {
					permitted = true
				}
			}
			if (!permitted) {
				const json: ResponseFormat = {
					error: "Not authorised",
					data: null,
				}
				return res.status(400).json(json)
			}
			req["token"] = decodedToken
			next()
		} catch (e: any) {
			const json: ResponseFormat = {
				error: "Error decoding token",
				data: null,
			}
			return res.status(400).json(json)
		}
	}
}


