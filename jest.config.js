export default {
	preset: "ts-jest",
	testEnvironment: "node",
	verbose: false,
	testRegex: "\.test\.ts",
	extensionsToTreatAsEsm: [".ts"],
	maxWorkers: 1,
	globals: {
		"ts-jest": {
			useESM: true,
		},
	},
}