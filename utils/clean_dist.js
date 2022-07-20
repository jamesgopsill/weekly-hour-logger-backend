import fs from "fs"
import path from "path"
import { fileURLToPath } from "url"

const __filename = fileURLToPath(import.meta.url)
const __dirname = path.dirname(__filename)

const dir = __dirname + "/../dist"

// Clear the dist directory
if (fs.existsSync(dir)) {
	fs.rmSync(dir, { recursive: true })
	console.log("Directory Cleaned!"); 
} else {
	console.log("./dist does not exist")
}