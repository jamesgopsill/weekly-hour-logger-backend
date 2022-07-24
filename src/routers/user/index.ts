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
 *         schema:
 *           type: object
 *           properties:
 *             error:
 *               type: string
 *               nullable: true
 *               description: "Displays an error if something went wrong."
 *             data:
 *               type: string
 *               description: "Returns world"
 */
router.get("/hello", hello)

/**
 * @openapi
 * /user/register:
 *   post:
 *     description: For registering one or more users. Requires authentication
 *     consumes:
 *       - application/json
 *     parameters:
 *       - in: header
 *         name: authorization
 *         description: A valid JWT token with admin scope.
 *         required: true
 *         schema:
 *           type: string
 *       - in: body
 *         schema:
 *           type: array
 *           items:
 *             type: object
 *             required:
 *               - name
 *               - email
 *               - password
 *             properties:
 *               name:
 *                 type: string
 *               email:
 *                 type: string
 *               password:
 *                 type: string
 *     responses:
 *       200:
 *         description: Returns success of the users have been added.
 *         schema:
 *           type: object
 *           properties:
 *             error:
 *               type: string
 *               nullable: true
 *               description: "Displays an error if something went wrong."
 *             data:
 *               type: string
 *               description: "Returns success"
 *       400:
 *         description: There was an error
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

// update user password
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
