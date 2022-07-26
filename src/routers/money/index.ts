import Router from "express-promise-router"
import { Validator } from "express-json-validator-middleware"
import { 
	SendMoneySchema,
	ListMoneySchema,

} from "./schemas"
import { 
    hello,
    sendMoney,
	updateMoney,
	listMoney,
	listAllMoney
} from "./fcns"
import { authorize } from "../../middleware"
import { User, UserScopes, Money } from "../../entities"

const router = Router()
const { validate } = new Validator({})

/**
 * @openapi
 * /money/hello:
 *   get:
 *     description: A simple hello test endpoint
 *     responses:
 *       200:
 *         description: Returns a "world".
 */
 router.get("/hello", hello)

// TODO
/**
 * @openapi
 * /money/sendMoney:
 *   get:
 *     description: Creates an entry for money distribution in a given week. 
 * 		Only one entry allowed per week number.
 * 		Calculates the week number based on date of entry.
 *     responses:
 *       200:
 *         description: Returns success of the data entry being created
 */
router.post(
	"/sendMoney",
	[validate({ body: SendMoneySchema })],
	sendMoney
)

/**
 * @openapi
 * /money/updateMoney:
 *   get:
 *     description: Updates an entry for money distribution in a given week. 
 * 		Only one entry allowed per week number.
 * 		Calculates the week number based on date of entry.
 *     responses:
 *       200:
 *         description: Returns success of the data entry being updated
 */
router.patch(
	"/updateMoney",
	[validate({ body: SendMoneySchema })],
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
	"/listMoney", 
	[validate({ body: ListMoneySchema })],
	listMoney)

/**
 * @openapi
 * /money/listAllMoney:
 *   get:
 *     description: Retrieves all entries for all groups. Requires admin.
 *     responses:
 *       200:
 *         description: Returns list of entries
 */
 router.get(
	"/listAllMoney", 
	authorize([UserScopes.ADMIN]),
	listAllMoney)


export const MoneyRouter = router