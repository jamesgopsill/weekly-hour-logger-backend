import {
	api,
	server,
	orm,
    SendMoneyArgs,
    ListMoneyArgs

} from "../src"
import supertest from "supertest"
import { validAdminToken } from "./tokens"

beforeAll(async () => {
	// initialise the app containing the api and orm connections
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

test("GET /money/hello", async () => {
	await supertest(api)
		.get("/money/hello")
		.expect(200)
		.then((res) => {
			const json = JSON.parse(res.text)
			expect(json.data).toBe("world")
		})
})

test("POST /money/sendMoney - no previous entry", async () => {
	const args: SendMoneyArgs = 
		{
            group: 1,
            date: "Tue Jul 26 2022 08:48:28 GMT+0100 (British Summer Time)",
            userOneMoney: 100,
            userTwoMoney: 100,
            userThreeMoney: 100,
            userFourMoney: 100,
            userFiveMoney: 0,
		}
	
	await supertest(api)
		.post("/money/sendMoney")
		.set("Content-Type", "application/json")
		// .set("authorization", `Bearer ${validAdminToken}`)
		.send(args)
		.expect(200)
		.then((res) => {
			console.log(res.text)
		})
})

test("POST /money/sendMoney - no previous entry [same group, different week]", async () => {
	const args: SendMoneyArgs = 
		{
            group: 1,
            date: "Tue Aug 2 2022 08:48:28 GMT+0100 (British Summer Time)",
            userOneMoney: 100,
            userTwoMoney: 100,
            userThreeMoney: 100,
            userFourMoney: 100,
            userFiveMoney: 0,
		}
	
	await supertest(api)
		.post("/money/sendMoney")
		.set("Content-Type", "application/json")
		// .set("authorization", `Bearer ${validAdminToken}`)
		.send(args)
		.expect(200)
		.then((res) => {
			console.log(res.text)
		})
})

test("POST /money/sendMoney - previous entry exists", async () => {
	const args: SendMoneyArgs = 
		{
            group: 1,
            date: "Wed Jul 27 2022 08:48:28 GMT+0100 (British Summer Time)",
            userOneMoney: 100,
            userTwoMoney: 100,
            userThreeMoney: 100,
            userFourMoney: 100,
            userFiveMoney: 0,
		}
	
	await supertest(api)
		.post("/money/sendMoney")
		.set("Content-Type", "application/json")
		// .set("authorization", `Bearer ${validAdminToken}`)
		.send(args)
		.expect(400)
		.then((res) => {
			console.log(res.text)
		})
})

test("PATCH /money/updateMoney - previous entry exists", async () => {
	const args: SendMoneyArgs = 
		{
            group: 1,
            date: "Wed Jul 27 2022 08:48:28 GMT+0100 (British Summer Time)",
            userOneMoney: 250,
            userTwoMoney: 50,
            userThreeMoney: 50,
            userFourMoney: 50,
            userFiveMoney: 0,
		}
	
	await supertest(api)
		.patch("/money/updateMoney")
		.set("Content-Type", "application/json")
		// .set("authorization", `Bearer ${validAdminToken}`)
		.send(args)
		.expect(200)
		.then((res) => {
			console.log(res.text)
		})
})

test("PATCH /money/updateMoney - no previous entry", async () => {
	const args: SendMoneyArgs = 
		{
            group: 2,
            date: "Wed Jul 27 2022 08:48:28 GMT+0100 (British Summer Time)",
            userOneMoney: 250,
            userTwoMoney: 50,
            userThreeMoney: 50,
            userFourMoney: 50,
            userFiveMoney: 0,
		}
	
	await supertest(api)
		.patch("/money/updateMoney")
		.set("Content-Type", "application/json")
		// .set("authorization", `Bearer ${validAdminToken}`)
		.send(args)
		.expect(404)
		.then((res) => {
			console.log(res.text)
		})
})

test("LIST /money/listMoney", async () => {
	const args: ListMoneyArgs = 
		{
            group: 1,
		}
	
	await supertest(api)
		.get("/money/listMoney")
		.set("Content-Type", "application/json")
		// .set("authorization", `Bearer ${validAdminToken}`)
		.send(args)
		.expect(200)
		.then((res) => {
			console.log(res.text)
		})
})

test("GET /money/listMoney - no previous entries for that group", async () => {
	const args: ListMoneyArgs = 
		{
            group: 3,
		}
	
	await supertest(api)
		.get("/money/listMoney")
		.set("Content-Type", "application/json")
		// .set("authorization", `Bearer ${validAdminToken}`)
		.send(args)
		.expect(400)
		.then((res) => {
			console.log(res.text)
		})
})

test("GET /money/listAllMoney", async () => {
	await supertest(api)
		.get("/money/listAllMoney")
		.set("Content-Type", "application/json")
		.set("authorization", `Bearer ${validAdminToken}`)
		.send()
		.expect(200)
		.then((res) => {
			console.log(res.text)
		})
})