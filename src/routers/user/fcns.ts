import { Request, Response } from "express"
import { RegisterArgs, LoginArgs } from "./interfaces"
import { ResponseFormat, verifyToken } from "../../utils"
import bcrypt from "bcrypt"
import { User, UserScopes } from "../../entities"
import jwt from "jsonwebtoken"
import { userRepo } from "../../entities"

const secret = process.env.SECRET || "secret"

export const hello = async (req: Request, res: Response) => {
	console.log("hello")
	return res.status(200).json({
		error: null,
		data: "world",
	})
}

/**
 * register takes an array of users that one wishes to add and adds them to the database. Only admins can register new users.
 * @param req
 * @param res
 * @returns
 */
export const register = async (req: Request, res: Response) => {
	// First check if the user is an admin
	console.log("[register] called")

	const t = verifyToken(req)
	if (t.error || !t.token.scopes.includes(UserScopes.ADMIN)) {
		const json: ResponseFormat = {
			error: "You do not have the permissions to perform this operation",
			data: null,
		}
		return res.status(400).json(json)
	}

	console.log("[register] token valid")

	// Get the array of users from the request
	const body: RegisterArgs[] = req.body

	// First loop to check whether any of the users you're tyring to add already exist
	for (const userToRegister of body) {
		// Check if user exists
		const query = {
			email: userToRegister.email,
		}
		console.log("[register] querying database", query)
		try {
			const user = await userRepo.findOne(query)
			if (user) {
				const json: ResponseFormat = {
					error: "User already exists",
					data: null,
				}
				return res.status(404).json(json)
			}
		} catch (e: any) {
			console.log(e)
			const json: ResponseFormat = {
				error: "Internal Server Error",
				data: null,
			}
			return res.status(500).json(json)
		}
	}

	// Now loop through again and add them to the database
	for (const userToRegister of body) {
		try {
			const saltRounds = 10
			const hash = bcrypt.hashSync(userToRegister.password, saltRounds)
			const user = new User()
			user.name = userToRegister.name
			user.email = userToRegister.email
			user.passwordHash = hash
			await userRepo.persistAndFlush(user)
		} catch (e: any) {
			console.log("Error creating user")
			return res.status(400).json({
				error: e.message,
				data: null,
			})
		}
	}

	// If we reach here then all is good. Return a successful result.
	const json: ResponseFormat = {
		error: null,
		data: "success",
	}

	return res.status(200).json(json)
}

export const login = async (req: Request, res: Response) => {
	const body: LoginArgs = req.body

	const users = await userRepo.find({ email: body.email })

	if (users.length < 1) {
		const json: ResponseFormat = {
			error: "User does not exist",
			data: null,
		}
		return res.status(404).json(json)
	}

	const user = users[0]

	try {
		if (!bcrypt.compareSync(body.password, user.passwordHash)) {
			const json: ResponseFormat = {
				error: "Passwords do not match",
				data: null,
			}
			return res.status(404).json(json)
		}

		const data = {
			name: user.name,
			email: user.email,
			scopes: user.scopes,
			// 24-hour expiry date
			iat: Math.floor(Date.now() / 1000) + 60 * 60 * 24,
		}

		const token = jwt.sign(data, secret, {
			algorithm: "HS256",
		})

		const json: ResponseFormat = {
			error: null,
			data: `Bearer ` + token,
		}
		return res.status(200).json(json)
	} catch (e: any) {
		const json: ResponseFormat = {
			error: e.message,
			data: null,
		}
		return res.status(400).json(json)
	}
}
