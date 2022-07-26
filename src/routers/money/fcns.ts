import { Request, Response } from "express"
import { SendMoneyArgs, ListMoneyArgs } from "./interfaces"
import { ResponseFormat } from "../interfaces"
import { Money, moneyRepo } from "../../entities"
import { RequestWithToken } from "../../middleware/interfaces"

// code snippet - produces the week number of the given week.
function getWeekNumber(d: Date) {
	// Copy date so don't modify original
	d = new Date(Date.UTC(d.getFullYear(), d.getMonth(), d.getDate()))

	// Set to nearest Thursday: current date + 4 - current day number
	// Make Sunday's day number 7
	d.setUTCDate(d.getUTCDate() + 4 - (d.getUTCDay() || 7))

	// Get first day of year
	var yearStart = new Date(Date.UTC(d.getUTCFullYear(), 0, 1))

	// Calculate full weeks to nearest Thursday
	var weekNo = Math.ceil(
		((d.valueOf() - yearStart.valueOf()) / 86400000 + 1) / 7
	)

	// Return week number
	return weekNo
}

// simple hello function to test endpoint
export const hello = async (req: Request, res: Response) => {
	console.log("[Saying Hello]")
	return res.status(200).json({
		error: null,
		data: "world",
	})
}

// Creates an entry with the money distribution for a given week.
// Checks that no entry exists for that group that week first.
export const sendMoney = async (req: RequestWithToken, res: Response) => {
	console.log("[Sending Money]")

	// retreive body data
	const body: SendMoneyArgs = req.body

	// retrieve week number for passed week. Must take a date object.
	const weekNo = getWeekNumber(new Date(body.date))
	// console.log(weekNo)

	const query = { group: body.group, weekNo: weekNo }
	const entry = await moneyRepo.find(query)
	// console.log(entry)

	// searches for an entry from this group this week and fails if one is found

	if (entry.length > 0) {
		const json: ResponseFormat = {
			error: "Entry already exists for this group this week",
			data: null,
		}
		return res.status(400).json(json)
	}

	// add the entry to the database
	try {
		const money = new Money()
		money.group = body.group
		money.date = body.date
		money.weekNo = weekNo
		money.userOneMoney = body.userOneMoney
		money.userTwoMoney = body.userTwoMoney
		money.userThreeMoney = body.userThreeMoney
		money.userFourMoney = body.userFourMoney
		money.userFiveMoney = body.userFiveMoney
		await moneyRepo.persistAndFlush(money)
	} catch (e: any) {
		console.log("Error sending money")
		const json: ResponseFormat = {
			error: "Error sending money",
			data: null,
		}
		return res.status(400).json(json)
	}

	// if we reach here, then all has happened as we want it to
	const json: ResponseFormat = {
		error: null,
		data: "Money sent",
	}

	return res.status(200).json(json)
}

// Updates an entry with the money distribution for a given week.
// Checks that an entry exists to update first.
export const updateMoney = async (req: RequestWithToken, res: Response) => {
	console.log("[Updating Money]")

	// retreive body data
	const body: SendMoneyArgs = req.body

	// retrieve week number for passed week. Must take a date object.
	const weekNo = getWeekNumber(new Date(body.date))
	// console.log(weekNo)

	const query = { group: body.group, weekNo: weekNo }
	const entries = await moneyRepo.find(query)
	// console.log(entries)

	// searches for an entry from this group this week and fails if one is found

	if (entries.length < 1) {
		const json: ResponseFormat = {
			error: "No entry found for this group this week. Create one instead",
			data: null,
		}
		return res.status(404).json(json)
	}

	const entry = entries[0]

	// update the entry in the database
	try {
		entry.group = body.group
		entry.date = body.date
		entry.weekNo = weekNo
		entry.userOneMoney = body.userOneMoney
		entry.userTwoMoney = body.userTwoMoney
		entry.userThreeMoney = body.userThreeMoney
		entry.userFourMoney = body.userFourMoney
		entry.userFiveMoney = body.userFiveMoney
		await moneyRepo.persistAndFlush(entry)
	} catch (e: any) {
		// console.log("Error updating money")
		const json: ResponseFormat = {
			error: "Error updating money",
			data: null,
		}
		return res.status(400).json(json)
	}

	// if we reach here, then all has happened as we want it to
	const json: ResponseFormat = {
		error: null,
		data: "Money updated",
	}

	return res.status(200).json(json)
}

// get a list of all money entries for a group
export const listMoney = async (req: RequestWithToken, res: Response) => {
	const body: ListMoneyArgs = req.body
	console.log("[Retrieving list of entries for group " + body.group + "]")

	const query = { group: body.group }
	const entries = await moneyRepo.find(query)

	// check that entries have been found
	if (entries.length < 1) {
		const json: ResponseFormat = {
			error: "No entries found for this group",
			data: null,
		}
		return res.status(400).json(json)
	}

	return res.status(200).json({ entries })
}

// get a list of all money entries for all groups
export const listAllMoney = async (req: RequestWithToken, res: Response) => {
	console.log("[Retrieving list of entries for all groups]")

	const query = {}
	const entries = await moneyRepo.find(query)

	// check that entries have been found
	if (entries.length < 1) {
		const json: ResponseFormat = {
			error: "No entries found",
			data: null,
		}
		return res.status(400).json(json)
	}

	return res.status(200).json({ entries })
}
