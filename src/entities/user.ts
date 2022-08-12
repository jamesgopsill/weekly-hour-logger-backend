import { Entity, PrimaryKey, Property } from "@mikro-orm/core"
import { v4 as uuidv4 } from "uuid"

@Entity()
export class User {
	@PrimaryKey()
	uuid: string = uuidv4()

	@Property()
	name!: string

	@Property()
	email!: string

	@Property()
	group: string = ""

	@Property()
	hashedPassword!: Buffer

	@Property()
	salt!: Buffer

	@Property()
	scopes: UserScopes[] = [UserScopes.USER]
}

export enum UserScopes {
	USER = "user",
	ADMIN = "admin",
}
