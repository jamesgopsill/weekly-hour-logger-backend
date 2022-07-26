import { AllowedSchema } from "express-json-validator-middleware"
import { authHeaderSchema } from "../../schema-utils/auth-schema"

export const registerSchema: AllowedSchema = {
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

export const loginSchema: AllowedSchema = {
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

export const updateSchema: AllowedSchema = {
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

export const passwordSchema: AllowedSchema = {
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

export const scopeSchema: AllowedSchema = {
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

export const userRouterSchema = {
	"user/hello": {
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
	"user/register": {
		post: {
			description:
				"For registering one or more users. Requires authentication.",
			consumes: ["application/json"],
			parameters: [
				authHeaderSchema,
				{
					in: "body",
					schema: registerSchema,
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
	"user/login": {
		post: {
			description: "Logs a user in",
			consumes: ["application/json"],
			parameters: [
				{
					in: "body",
					schema: loginSchema,
				},
			],
			responses: {
				200: {
					description: "Returns a Bearer token.",
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
								description: "Returns a Bearer token",
							},
						},
					},
				},
			},
		},
	},
	"user/update": {
		patch: {
			description: "Updates a user's details",
			consumes: ["application/json"],
			parameters: [
				authHeaderSchema,
				{
					in: "body",
					schema: updateSchema,
				},
			],
			responses: {
				200: {
					description: "Returns success.",
				},
			},
		},
	},
	"user/list": {
		get: {
			description: "List users",
			consumes: ["application/json"],
			parameters: [authHeaderSchema],
			responses: {
				200: {
					description: "Returns a list of users.",
				},
			},
		},
	},
	"user/password": {
		patch: {
			description: "Updates a users password",
			consumes: ["application/json"],
			parameters: [
				authHeaderSchema,
				{
					in: "body",
					schema: passwordSchema,
				},
			],
			responses: {
				200: {
					description: "Returns success if the password has changed.",
				},
			},
		},
	},
	"user/refresh-token": {
		post: {
			description: "Refreshes the expiry date of a current valid token.",
			consumes: ["application/json"],
			parameters: [authHeaderSchema],
			responses: {
				200: {
					description: "Returns the new Bearer token",
				},
			},
		},
	},
}
