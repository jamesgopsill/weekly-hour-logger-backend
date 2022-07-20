import express, { Request, Response } from "express"
import cors from "cors"
import { UserRouter } from "./routers"
// import { hello } from "./routers"
import { ValidationError } from "express-json-validator-middleware"

const port = process.env.PORT || 3000

export const api = express()

api.use(express.json())
api.use(cors())

api.use("/user", UserRouter)

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

export const server = api.listen(port, () => {
	console.log(`API start running on http://localhost:${port}`)
})
