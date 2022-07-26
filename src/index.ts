import express, { Request, Response } from "express"
import cors from "cors"
import { UserRouter, MoneyRouter } from "./routers"
import { ValidationError } from "express-json-validator-middleware"
import { orm } from "./entities"
import swaggerUi from "swagger-ui-express"
import swaggerJSDoc from "swagger-jsdoc"

export { orm, UserScopes } from "./entities"
export * from "./routers/user/interfaces"
export * from "./routers/money/interfaces"

const swaggerOptions = {
	swaggerDefinition: {
		info: {
			title: "Gopsill and Sniders Api", // Title (required)
			version: "0.1.0", // Version (required)
		},
	},
	apis: ["./**/*.js"], // Path to the API docs
}
const swaggerSpec = swaggerJSDoc(swaggerOptions)
console.log(swaggerSpec)

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

api.get("/api-docs.json", function (req: Request, res: Response) {
	res.setHeader("Content-Type", "application/json")
	res.send(swaggerSpec)
})

api.use("/api-docs", swaggerUi.serve, swaggerUi.setup(swaggerSpec))

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
