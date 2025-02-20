<script lang="ts">
	import LineChart from '$lib/components/LineChart.svelte'

	export let data
	let results = data?.results || []

	let start = data.start || '2014-01-01'
    let end = data.end || '2024-01-01'

    function handleSubmit(event: Event) {
        event.preventDefault()
		const formattedStart = start.replace(/-/g, '_')
        const formattedEnd = end.replace(/-/g, '_')
        const params = new URLSearchParams({ start: formattedStart, end: formattedEnd })
        window.location.search = params.toString()
    }

	function getNW(results: WeeklyRecords[]): number[][] {
		let res: number[][] = []

		for (const recordGroup of results) {
			if (!Array.isArray(recordGroup.Records)) {
				console.error("Invalid record structure: ", recordGroup.Records)
				continue
			}

			let strategyResults: number[] = []
			for (const record of recordGroup.Records) {
				strategyResults.push(record.NetWorth)
			}
			res.push(strategyResults)
		}

		return res
	}

	function formatDateToCustomFormat(dateString: string) {
		const date = new Date(dateString)
		const options: Intl.DateTimeFormatOptions = { month: 'short', year: '2-digit' }
		return new Intl.DateTimeFormat('en-US', options).format(date).replace(' ', " '")
	}

	function getDates(results: WeeklyRecords[]): string[] {
		let res: string[] = []

		for (const record of results[0].Records) {
				res.push(formatDateToCustomFormat(record.Time))
		}

		return res
	}

	let nwValues = getNW(results)
	let monthsValues = getDates(results)

	const colors = [
		{ border: 'rgba(75, 192, 192, 1)', background: 'rgba(75, 192, 192, 0.2)' },
		{ border: 'rgba(255, 99, 132, 1)', background: 'rgba(255, 99, 132, 0.2)' },
		{ border: 'rgba(255, 165, 0, 1)', background: 'rgba(255, 165, 0, 0.2)' },
		{ border: 'rgba(54, 162, 235, 1)', background: 'rgba(54, 162, 235, 0.2)' },
		{ border: 'rgba(153, 102, 255, 1)', background: 'rgba(153, 102, 255, 0.2)' },
	]

	let stats = results.map((result, index) => ({
		label: result.Strategy, // Use strategy name as the label
		data: nwValues[index],
		borderColor: colors[index % colors.length].border,
		backgroundColor: colors[index % colors.length].background
	}))

</script>

<div class="container h-screen w-full flex flex-col justify-center items-center px-10">
	<form on:submit={handleSubmit} class="flex flex-col ">
		<label>
			Start Date:
			<input class="input" type="date" bind:value={start} />
		</label>
		<label>
			End Date:
			<input class="input" type="date" bind:value={end} />
		</label>
		<button class="btn" type="submit">Fetch Data</button>
	</form>

	<LineChart {...{ stats, label: "" , xAxisLabels: monthsValues}}></LineChart>
</div>
