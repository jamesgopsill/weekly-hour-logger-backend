import express, { Request, Response } from "express"
import cors from "cors"
import { UserRouter, MoneyRouter } from "./routers"
import { ValidationError } from "express-json-validator-middleware"
import { orm } from "./entities"
import swaggerUi from "swagger-ui-express"
import { combinedSchema } from "./schema-utils/combined-schema"

export { orm, UserScopes } from "./entities"
export * from "./routers/user/interfaces"
export * from "./routers/money/interfaces"

export const api = express()

api.use(express.json())
api.use(cors())

api.use("/user", UserRouter)
api.use("/money", MoneyRouter) // ## added this ##

// Covers the return of any validation errors
api.use((error: any, req: Request, res: Response, next: any) => {
	console.log("Validation Error")
	if (error instanceof ValidationError) {
		const errorJSON = {
			error: error.validationErrors,
			data: null,
		}
		res.status(400).send(errorJSON)
		next()
	} else {
		next(error)
	}
})

api.get("/docs.json", function (req: Request, res: Response) {
	res.setHeader("Content-Type", "application/json")
	res.send(combinedSchema)
})

api.use("/docs", swaggerUi.serve, swaggerUi.setup(combinedSchema))

const port = process.env.PORT || 3000

export const server = api.listen(port, () => {
	console.log(`API start running on http://localhost:${port}`)
})

process.on("SIGINT", () => {
	console.log("Shutting down server and database connection")
	server.close()
	orm.close()
	process.exit()
})
