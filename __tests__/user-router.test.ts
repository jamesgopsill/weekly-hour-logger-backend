import { api, server } from "../src"
import { orm } from "../src/entities"
import supertest from "supertest"
import { validAdminToken } from "./tokens"
import { LoginArgs, RegisterArgs } from "../src/routers/user/interfaces"

beforeAll(async () => {
	// initialise the app containing the api and orm connections
	//await initApp("../db/test.db", "test")
	const generator = orm.getSchemaGenerator()
	await generator.dropSchema()
	await generator.createSchema()
	await generator.updateSchema()
})

afterAll(async () => {
	// close the orm and server connections
	orm.close()
	server.close()
})

test("GET /user/hello", async () => {
	await supertest(api)
		.get("/user/hello")
		.expect(200)
		.then((res) => {
			const json = JSON.parse(res.text)
			expect(json.data).toBe("world")
		})
})

test("POST /user/register - valid", async () => {
	const args: RegisterArgs[] = [
		{
			name: "Test User",
			email: "test@test.com",
			password: "test",
		},
	]
	await supertest(api)
		.post("/user/register")
		.set("Content-Type", "application/json")
		.set("authorization", `Bearer ${validAdminToken}`)
		.send(args)
		.expect(200)
		.then((res) => {
			console.log(res.text)
		})
})

test("POST /user/register - no token", async () => {
	const args: RegisterArgs[] = [
		{
			name: "Test User",
			email: "test@test.com",
			password: "test",
		},
	]
	await supertest(api)
		.post("/user/register")
		.set("Content-Type", "application/json")
		.send(args)
		.expect(400)
})

test("POST /user/login", async () => {
	const args: LoginArgs = {
		email: "test@test.com",
		password: "test",
	}
	await supertest(api)
		.post("/user/login")
		.set("Content-Type", "application/json")
		.send(args)
		.expect(200)
		.then((res) => {
			console.log(res.text)
		})
})
