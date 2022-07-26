import { AllowedSchema } from "express-json-validator-middleware"
import { authHeaderSchema } from "../../schema-utils/auth-schema"

export const SendMoneySchema: AllowedSchema = {
	type: "object",
	required: [
		"group",
		"date",
		"userOneMoney",
		"userTwoMoney",
		"userThreeMoney",
		"userFourMoney",
	],
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
		},
	},
}

export const moneyRouterSchema = {
	"money/hello": {
		get: {
			description: "A simple hello test endpoint.",
			responses: {
				200: {
					description: "Returns world.",
					schema: {
						type: "object",
						properties: {
							error: {
								type: "string",
								nullable: true,
								description:
									"Displays an error if something went wrong.",
							},
							data: {
								type: "string",
								description: "Returns world",
							},
						},
					},
				},
			},
		},
	},
	"money/send": {
		post: {
			description:
				"Creates an entry for money distribution in a given week. Only one entry allowed per week number. Calculates the week number based on date of entry.",
			consumes: ["application/json"],
			parameters: [
				authHeaderSchema,
				{
					in: "body",
					schema: SendMoneySchema,
				},
			],
			responses: {
				200: {
					description: "Returns success.",
					schema: {
						type: "object",
						properties: {
							error: {
								type: "string",
								nullable: true,
								description:
									"Displays an error if something went wrong.",
							},
							data: {
								type: "string",
								description: "Returns success",
							},
						},
					},
				},
			},
		},
	},
	"money/update": {
		patch: {
			description:
				"Updates an entry for money distribution in a given week. Only one entry allowed per week numberalculates the week number based on date of entry.",
			consumes: ["application/json"],
			parameters: [
				authHeaderSchema,
				{
					in: "body",
					schema: SendMoneySchema,
				},
			],
			responses: {
				200: {
					description: "Returns success.",
					schema: {
						type: "object",
						properties: {
							error: {
								type: "string",
								nullable: true,
								description:
									"Displays an error if something went wrong.",
							},
							data: {
								type: "string",
								description: "Returns success",
							},
						},
					},
				},
			},
		},
	},
	"money/list": {
		get: {
			description:
				"Retrieves all entries for a specified group. Does NOT require admin.",
			consumes: ["application/json"],
			parameters: [
				authHeaderSchema,
				{
					in: "body",
					schema: ListMoneySchema,
				},
			],
			responses: {
				200: {
					description: "Returns success.",
					schema: {
						type: "object",
						properties: {
							error: {
								type: "string",
								nullable: true,
								description:
									"Displays an error if something went wrong.",
							},
							data: {
								type: "string",
								description: "Returns money",
							},
						},
					},
				},
			},
		},
	},
	"money/list-all": {
		get: {
			description: "Retrieves all entries for all groups. Requires admin",
			consumes: ["application/json"],
			parameters: [authHeaderSchema],
			responses: {
				200: {
					description: "Returns success.",
					schema: {
						type: "object",
						properties: {
							error: {
								type: "string",
								nullable: true,
								description:
									"Displays an error if something went wrong.",
							},
							data: {
								type: "string",
								description: "Returns money",
							},
						},
					},
				},
			},
		},
	},
}
