export default {
	preset: "ts-jest",
	testEnvironment: "node",
	verbose: true,
	testRegex: "\.test\.ts",
	extensionsToTreatAsEsm: [".ts"],
	globals: {
		"ts-jest": {
			useESM: true,
		},
	},
}