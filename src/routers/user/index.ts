import Router from "express-promise-router"
import { Validator } from "express-json-validator-middleware"
import {
	registerSchema,
	loginSchema,
	updateSchema,
	passwordSchema,
	scopeSchema,
} from "./schema"
import {
	register,
	hello,
	login,
	updateUser,
	listUsers,
	updatePassword,
	updateScope,
	refreshToken
} from "./fcns"
import { authorize } from "../../middleware"
import { UserScopes } from "../../entities"

const router = Router()
const { validate } = new Validator({})

/**
 * /user/hello
 */
router.get("/hello", hello)

/**
 * /user/register
 */
router.post(
	"/register",
	[authorize([UserScopes.ADMIN]), validate({ body: registerSchema })],
	register
)

/**
 * /user/login
 */
router.post("/login", validate({ body: loginSchema }), login)

/**
 * /user/update
 */
router.patch(
	"/update",
	[authorize([UserScopes.ADMIN]), validate({ body: updateSchema })],
	updateUser
)

/**
 * /user/list
 */
router.get("/list", authorize([UserScopes.ADMIN]), listUsers)

/**
 * /user/password
 */
router.patch(
	"/password",
	[authorize([UserScopes.ADMIN, UserScopes.USER]), validate({ body: passwordSchema })],
	updatePassword
)

/**
 * /user/scopes
 */
router.patch(
	"/scopes",
	[authorize([UserScopes.ADMIN]), validate({ body: scopeSchema })],
	updateScope
)

// Routes TODO
// refresh a token
router.post("/refresh-token",
[	authorize([UserScopes.USER])],
	refreshToken
)

export const UserRouter = router


