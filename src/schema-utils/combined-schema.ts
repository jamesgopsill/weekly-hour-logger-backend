import { moneyRouterSchema } from "../routers/money/schema"
import { userRouterSchema } from "../routers/user/schema"

const info = {
	info: {
		title: "Gopsill and Sniders API",
		version: "0.1.0",
	},
	swagger: "2.0",
}

export const combinedSchema = {
	...info,
	paths: {
		...userRouterSchema,
		...moneyRouterSchema,
	},
	definitions: {},
	responses: {},
	parameters: {},
	securityDefinitions: {},
	tags: {},
}
