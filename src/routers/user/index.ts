import Router from "express-promise-router"
import { Validator } from "express-json-validator-middleware"
import {
	RegisterSchema,
	LoginSchema,
	UserUpdateSchema,
	PasswordUpdateSchema,
	ScopeSchema,
} from "./schemas"
import {
	register,
	hello,
	login,
	updateUser,
	listUsers,
	updatePassword,
	updateScope,
} from "./fcns"
import { authorize } from "../../middleware"
import { UserScopes } from "../../entities"

const router = Router()
const { validate } = new Validator({})

/**
 * @openapi
 * /user/hello:
 *   get:
 *     description: A simple hello test endpoint
 *     responses:
 *       200:
 *         description: Returns a "world".
 */
router.get("/hello", hello)

/**
 * @openapi
 * /user/register:
 *   post:
 *     description: For registering one or more users. Requires authentication
 *     responses:
 *       200:
 *         description: Returns success of the users have been added.
 */
router.post(
	"/register",
	[authorize([UserScopes.ADMIN]), validate({ body: RegisterSchema })],
	register
)

// Login a user
router.post("/login", validate({ body: LoginSchema }), login)

// Routes TODO
// perform a token check
// router.get("/token")

/**
 * @openapi
 * /user/updateUser:
 *   patch:
 *     description: To update user name, email, or group. Requires authentication.
 *     responses:
 *       200:
 *         description: Returns success of the user has been updated.
 */
router.patch(
	"/updateUser",
	[authorize([UserScopes.ADMIN]), validate({ body: UserUpdateSchema })],
	updateUser
)

// Get a list of users
router.get("/list", authorize([UserScopes.ADMIN]), listUsers)

/**
 * @openapi
 * /user/updatePassword:
 *   get:
 *     description: Update the password for a user. Requires the old password and a new password
 *     responses:
 *       200:
 *         description: Returns success of the password being updated.
 */
router.patch(
	"/updatePassword",
	[authorize([UserScopes.ADMIN]), validate({ body: PasswordUpdateSchema })],
	updatePassword
)

// update user scopes
router.patch(
	"/updateScope",
	[authorize([UserScopes.ADMIN]), validate({ body: ScopeSchema })],
	updateScope
)

export const UserRouter = router
