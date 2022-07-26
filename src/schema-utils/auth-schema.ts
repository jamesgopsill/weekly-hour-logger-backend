export const authHeaderSchema = {
	in: "header",
	name: "authorization",
	description: "A valid JWT token.",
	required: true,
	schema: {
		type: "string",
	},
}
