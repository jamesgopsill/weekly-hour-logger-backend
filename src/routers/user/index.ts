import Router from "express-promise-router"
import { Validator } from "express-json-validator-middleware"
import { RegisterSchema, LoginSchema } from "./schemas"
import { register, hello, login } from "./fcns"

const router = Router()
const { validate } = new Validator({})

router.get("/hello", hello)

// Register a user
router.post("/register", validate({ body: RegisterSchema }), register)

// Login a user
router.post("/login", validate({ body: LoginSchema }), login)

// Routes TODO
// perform a token check
// router.get("/token")

// update user details
// router.patch("/")

// update users
// router.patch("/list")

// update user password
// router.patch("/password")

// update user scopes
// router.patch("/admin/scopes")

export const UserRouter = router
