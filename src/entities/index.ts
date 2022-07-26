import { MikroORM } from "@mikro-orm/core"
import { User } from "./user"
import { Money } from "./money"
import path from "path"
import { fileURLToPath } from "url"

export * from "./user"
export * from "./money" // added this...

const __filename = fileURLToPath(import.meta.url)
const __dirname = path.dirname(__filename)

const dbName = process.env.DBNAME || "test.db"
const dbPath = __dirname + "/../../db/" + dbName

export const orm = await MikroORM.init({
	dbName: dbPath,
	type: "sqlite",
	entities: [User, Money],
	debug: false,
	allowGlobalContext: true,
})

// const em = orm.em
export const userRepo = orm.em.getRepository(User)
export const moneyRepo = orm.em.getRepository(Money) // added this...
