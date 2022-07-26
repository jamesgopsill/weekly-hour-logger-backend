import Router from "express-promise-router"
import { Validator } from "express-json-validator-middleware"
import { SendMoneySchema, ListMoneySchema } from "./schema"
import { hello, sendMoney, updateMoney, listMoney, listAllMoney } from "./fcns"
import { authorize } from "../../middleware"
import { UserScopes } from "../../entities"

const router = Router()
const { validate } = new Validator({})

/**
 * /money/hello
 */
router.get("/hello", hello)

/**
 * /money/send-money
 */
router.post(
	"/send",
	[authorize([UserScopes.USER]), validate({ body: SendMoneySchema })],
	sendMoney
)

/**
 * /money/update-money:
 */
router.patch(
	"/update",
	[authorize([UserScopes.USER]), validate({ body: SendMoneySchema })],
	updateMoney
)

/**
 * @openapi
 * /money/listMoney:
 *   get:
 *     description: Retrieves all entries for a specified group. Does NOT require admin.
 *     responses:
 *       200:
 *         description: Returns list of entries
 */
router.get(
	"/list",
	[authorize([UserScopes.USER]), validate({ body: ListMoneySchema })],
	listMoney
)

/**
 * @openapi
 * /money/listAllMoney:
 *   get:
 *     description: Retrieves all entries for all groups. Requires admin.
 *     responses:
 *       200:
 *         description: Returns list of entries
 */
router.get("/list-all", authorize([UserScopes.ADMIN]), listAllMoney)

export const MoneyRouter = router
