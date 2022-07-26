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
import { User, UserScopes } from "../../entities"

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

/**
 * @openapi
 * /user/updateUser:
 *   patch:
 *     description: To update user name, email, or group. Requires authentication.
 *     responses:
 *       200:
 *         description: Returns success of the user has been updated.
 */
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

/**
 * @openapi
 * /user/list:
 *   patch:
 *     description: Retrieves all users in the database. Requires admin.
 *     responses:
 *       200:
 *         description: Returns user list. Password hashes have been removed.
 */
router.get("/list", authorize([UserScopes.ADMIN]), listUsers)

/**
 * @openapi
 * /user/updatePassword:
 *   patch:
 *     description: Updates user password. Users are able to do this themselves.
 *     responses:
 *       200:
 *         description: Returns success of the user password being updated.
 */
router.patch(
	"/updatePassword",
	[validate({ body: PasswordUpdateSchema })],
	updatePassword
)

/**
 * @openapi
 * /user/updateScope:
 *   patch:
 *     description: To update user scopes. Requires admin.
 *     responses:
 *       200:
 *         description: Returns success of the user scope being updated.
 */
router.patch(
	"/updateScope",
	[authorize([UserScopes.ADMIN]), validate({ body: ScopeSchema })],
	updateScope
)

export const UserRouter = router
