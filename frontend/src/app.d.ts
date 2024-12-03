// See https://kit.svelte.dev/docs/types#app
// for information about these interfaces
// and what to do when importing types
declare namespace App {
	// interface Locals {}
	// interface PageData {}
	// interface Error {}
	// interface Platform {}
}

interface DataSet {
	label: string
	data: number[]
	borderColor: string
	backgroundColor: string
}

interface WeeklyRecord {
	Time: string
	NetWorth: number
	SnpValue: number
	Shares: number
	Reserves: number
}

interface WeeklyRecords {
	Strategy: string
	Records: WeeklyRecord[]
}