import { MikroORM } from "@mikro-orm/core"
import { User } from "./user"
import path from "path"
import { fileURLToPath } from "url"

export * from "./user"

const __filename = fileURLToPath(import.meta.url)
const __dirname = path.dirname(__filename)

const dbName = process.env.DBNAME || "test.db"
const dbPath = __dirname + "/../../db/" + dbName

export const orm = await MikroORM.init({
	dbName: dbPath,
	type: "sqlite",
	entities: [User],
	debug: true,
	allowGlobalContext: true,
})

// const em = orm.em
export const userRepo = orm.em.getRepository(User)
