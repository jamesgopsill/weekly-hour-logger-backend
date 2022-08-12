import { Request, Response } from "express"
import {
	RegisterArgs,
	LoginArgs,
	UserUpdateArgs,
	PasswordUpdateArgs,
	ScopeArgs,
} from "./interfaces"
import { ResponseFormat } from "../interfaces"
import crypto from "crypto"
import { User, UserScopes, userRepo } from "../../entities"
import jwt from "jsonwebtoken"
import { RequestWithToken } from "../../middleware/interfaces"

const secret = process.env.SECRET || "secret"

/**
 * hello returns "world"
 * @param req
 * @param res
 * @returns
 */
export const hello = async (req: Request, res: Response) => {
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
export const register = async (req: RequestWithToken, res: Response) => {
	// Get the array of users from the request
	const body: RegisterArgs[] = req.body

	// First loop to check whether any of the users you're tyring to add already exist
	for (const userToRegister of body) {
		// Check if user exists
		const query = {
			email: userToRegister.email,
		}
		// console.log("[register] querying database", query)
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
			const salt = crypto.randomBytes(16)
			const hashedPassword = crypto.pbkdf2Sync(
				userToRegister.password,
				salt,
				310000,
				32,
				"sha256"
			)
			const user = new User()
			user.name = userToRegister.name
			user.email = userToRegister.email
			user.hashedPassword = hashedPassword
			user.salt = salt
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
		const hashedPassword = crypto.pbkdf2Sync(
			body.password,
			user.salt,
			310000,
			32,
			"sha256"
		)
		if (!crypto.timingSafeEqual(user.hashedPassword, hashedPassword)) {
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

// update user information. Only admins can do this at the moment.
export const updateUser = async (req: RequestWithToken, res: Response) => {
	// Get the body data, and turn the email into a search query
	const body: UserUpdateArgs = req.body
	const query = {
		email: body.email,
	}

	// find the user with that email. Returns all for that user
	const users = await userRepo.find(query)

	// check that a user has been found
	if (users.length < 1) {
		const json: ResponseFormat = {
			error: "No user found",
			data: null,
		}
		return res.status(400).json(json)
	}

	const user = users[0]

	if (
		user.email != req.token.email &&
		!req.token.scopes.includes(UserScopes.ADMIN)
	) {
		const json: ResponseFormat = {
			error: "You do not have the permissions to perform this operation.",
			data: null,
		}
		return res.status(400).json(json)
	}

	// tries to update the database with the new user info
	try {
		user.name = body.name
		user.email = body.email
		user.group = body.group
		await userRepo.persistAndFlush(user)
	} catch (e: any) {
		console.log("Error updating user")
		const json: ResponseFormat = {
			error: "Error updating user",
			data: null,
		}
		return res.status(400).json(json)
	}

	// if we reach here, then all has happened as we want it to
	const json: ResponseFormat = {
		error: null,
		data: "success",
	}

	return res.status(200).json(json)
}

// update user password. Users can update this themselves.
export const updatePassword = async (req: RequestWithToken, res: Response) => {
	// Get the body data, and turn the email into a search query
	const body: PasswordUpdateArgs = req.body

	// find the user with that email. Returns all for that user
	const users = await userRepo.find({ email: body.email })

	// check that a user has been found
	if (users.length < 1) {
		const json: ResponseFormat = {
			error: "No user found",
			data: null,
		}
		return res.status(400).json(json)
	}

	// extract the user info from the array
	const user = users[0]

	// Check if it is the user or admin requesting the change of password.
	if (
		user.email != req.token.email &&
		!req.token.scopes.includes(UserScopes.ADMIN)
	) {
		const json: ResponseFormat = {
			error: "You do not have the permissions to perform this operation.",
			data: null,
		}
		return res.status(400).json(json)
	}

	// check the old password matches
	try {
		const hashedOldPassword = crypto.pbkdf2Sync(
			body.oldPassword,
			user.salt,
			310000,
			32,
			"sha256"
		)
		if (!crypto.timingSafeEqual(user.hashedPassword, hashedOldPassword)) {
			const json: ResponseFormat = {
				error: "Passwords do not match",
				data: null,
			}
			return res.status(404).json(json)
		}

		const salt = crypto.randomBytes(16)
		const hashedPassword = crypto.pbkdf2Sync(
			body.newPassword,
			salt,
			310000,
			32,
			"sha256"
		)

		user.hashedPassword = hashedPassword
		user.salt = salt

		await userRepo.persistAndFlush(user)
	} catch (e: any) {
		console.log("Error updating password")
		const json: ResponseFormat = {
			error: "Error updating password",
			data: null,
		}
		return res.status(400).json(json)
	}

	// if we reach here, then all has happened as we want it to
	const json: ResponseFormat = {
		error: null,
		data: "success",
	}

	return res.status(200).json(json)
}

// get a list of all users in the database
export const listUsers = async (req: RequestWithToken, res: Response) => {
	const users = await userRepo.find({ scopes: [UserScopes.USER] })

	// check that a user has been found
	if (users.length < 1) {
		const json: ResponseFormat = {
			error: "No user found",
			data: null,
		}
		return res.status(404).json(json)
	}

	const filteredUserInfo = users.map((u) => {
		delete u.hashedPassword
		delete u.salt
		return u
	})
	const json: ResponseFormat = {
		error: null,
		data: filteredUserInfo,
	}
	return res.status(200).json(json)
}

// update user scope. Requires admin.
export const updateScope = async (req: RequestWithToken, res: Response) => {
	// Get the body data, and turn the email into a search query
	const body: ScopeArgs = req.body
	const query = { email: body.email }

	// find the user with that email. Returns all for that user
	const users = await userRepo.find(query)

	// check that a user has been found
	if (users.length < 1) {
		const json: ResponseFormat = {
			error: "No user found",
			data: null,
		}
		return res.status(404).json(json)
	}

	// extract the user info from the array
	const user = users[0]

	// Check the requested scope, and if it's admin, create admin combined rights
	try {
		if (body.scope == UserScopes.ADMIN) {
			user.scopes = [UserScopes.USER, UserScopes.ADMIN]
			await userRepo.persistAndFlush(user)
			const json: ResponseFormat = {
				error: null,
				data: "success",
			}
			return res.status(200).json(json)
		}
		user.scopes = [UserScopes.USER]
		await userRepo.persistAndFlush(user)
		const json: ResponseFormat = {
			error: null,
			data: "success",
		}
		return res.status(200).json(json)
	} catch {
		console.log("Error updating user scope")
		const json: ResponseFormat = {
			error: "Error updating user scope",
			data: null,
		}
		return res.status(400).json(json)
	}
}

export const refreshToken = async (req: RequestWithToken, res: Response) => {
	const data = {
		name: req.token.name,
		email: req.token.email,
		scopes: req.token.scopes,
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
}
