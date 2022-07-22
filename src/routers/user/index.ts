import Router from "express-promise-router"
import { Validator } from "express-json-validator-middleware"
import { RegisterSchema, LoginSchema, UserUpdateSchema, PasswordUpdateSchema } from "./schemas"
import { register, hello, login, updateUser, updatePassword } from "./fcns"
import { authorize } from "../../middleware"

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
	[authorize, validate({ body: RegisterSchema })],
	register
)

// Login a user
router.post("/login", validate({ body: LoginSchema }), login)

// Routes TODO
// perform a token check
// router.get("/token")

// update user details
router.patch("/updateUser", [authorize, validate({body: UserUpdateSchema})], updateUser)

// update users
// router.patch("/list")

// update user password
router.patch("/updatePassword", [authorize, validate({ body: PasswordUpdateSchema }) ], updatePassword)

// update user scopes
// router.patch("/admin/scopes")

export const UserRouter = router
