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

export const UserUpdateSchema: AllowedSchema = {
	type: "object",
	required: ["name", "email", "group"],
	properties: {
		name: {
			type: "string",
		},
		email: {
			type: "string",
		},
		group: {
			type: "string",
		},
	},
}

export const PasswordUpdateSchema: AllowedSchema = {
	type: "object",
	required: ["email", "oldPassword", "newPassword"],
	properties: {
		email: {
			type: "string",
		},
		oldPassword: {
			type: "string",
		},
		newPassword: {
			type: "string",
		},
	},
}

export const ScopeSchema: AllowedSchema = {
	type: "object",
	required: ["email", "scope"],
	properties: {
		email: {
			type: "string",
		},
		scope: {
			type: "string",
		},
	},
}
