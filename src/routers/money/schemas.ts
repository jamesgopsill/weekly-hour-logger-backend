import { AllowedSchema } from "express-json-validator-middleware"

export const SendMoneySchema: AllowedSchema = {
	type: "object",
	required: ["group", "date", "userOneMoney", "userTwoMoney", "userThreeMoney", "userFourMoney"],
	properties: {
		group: {
			type: "number",
		},
		date: {
			type: "string",
		},
        weekNo: {
			type: "number",
		},
        userOneMoney: {
			type: "number",
		},
        userTwoMoney: {
			type: "number",
		},
        userThreeMoney: {
			type: "number",
		},
        userFourMoney: {
			type: "number",
		},
        userFiveMoney: {
			type: "number",
		},
	},
}

export const ListMoneySchema: AllowedSchema = {
	type: "object",
	required: ["group"],
	properties: {
		group: {
			type: "number",
        }
    }
}