import { Entity, PrimaryKey, Property } from "@mikro-orm/core"
import { v4 as uuidv4 } from "uuid"

// group, user1money, user2money, user3money, user4money, week

@Entity()
export class Money {
	@PrimaryKey()
	uuid: string = uuidv4()

	@Property()
	group!: number

	@Property()
	date!: string

	@Property()
	weekNo!: number

	@Property()
	userOneMoney!: number

	@Property()
	userTwoMoney!: number

	@Property()
	userThreeMoney!: number

	@Property()
	userFourMoney!: number

	@Property()
	userFiveMoney: number = 0
}
