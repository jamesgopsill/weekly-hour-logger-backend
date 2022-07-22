import { Request, Response } from "express"
import { RegisterArgs, LoginArgs, UserUpdateArgs, PasswordUpdateArgs, ScopeArgs } from "./interfaces"
import { ResponseFormat } from "../interfaces"
import bcrypt from "bcrypt"
import { User, UserScopes, userRepo } from "../../entities"
import jwt from "jsonwebtoken"
import { RequestWithToken } from "../../middleware/interfaces"

const secret = process.env.SECRET || "secret"

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
	console.log("[register] called")

	// First check if the user is an admin
	// Get the decoded token that has been passed along with the request from our authorize function.
	const t = req.token

	if (!t.scopes.includes(UserScopes.ADMIN)) {
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

// update user information. Only admins can do this at the moment.
export const updateUser = async (req: Request, res: Response) => {

	// Get the body data, and turn the email into a search query
	const body: UserUpdateArgs = req.body
	const query = { email: body.email }

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
	const users = await userRepo.find( {email: body.email} )
	
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

	// check the old password matches
	try {
		if (!bcrypt.compareSync(body.oldPassword, user.passwordHash)) {
			const json: ResponseFormat = {
				error: "Passwords do not match",
				data: null,
			}
			return res.status(400).json(json)
		}

		// update the new password
		const saltRounds = 10
		const hash = bcrypt.hashSync(body.newPassword, saltRounds)
		user.passwordHash = hash
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

	const users = await userRepo.find( {scopes: [UserScopes.USER]} )
	console.log("I got here!")

	// check that a user has been found
	if (users.length < 1) {
		const json: ResponseFormat = {
			error: "No user found",
			data: null,
		}
		return res.status(404).json(json)
	}

	console.log(users[0])
	return res.status(200).json({users})
}

// update user scope. Requires admin.
export const updateScope = async (req: RequestWithToken, res: Response) => {

	// Get the body data, and turn the email into a search query
	const body: ScopeArgs = req.body
	const query = { email: body.email }

	// find the user with that email. Returns all for that user
	const users = await userRepo.find( query )
	
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
