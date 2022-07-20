import { AllowedSchema } from "express-json-validator-middleware"

export const RegisterSchema: AllowedSchema = {
	type: "array",
	items: {
		type: "object",
		properties: {
			name: {
				type: "string",
			},
			email: {
				type: "string",
			},
			password: {
				type: "string",
			},
		},
		required: ["name", "email", "password"],
	},
}

export const LoginSchema: AllowedSchema = {
	type: "object",
	required: ["email", "password"],
	properties: {
		email: {
			type: "string",
		},
		password: {
			type: "string",
		},
	},
}
