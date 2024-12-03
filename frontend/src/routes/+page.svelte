<script lang="ts">
	import LineChart from '$lib/components/LineChart.svelte'

	export let data
	let results = data?.results || []

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

	let nwValues = getNW(results)

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

<div class="container h-screen w-full flex justify-center items-center px-10">
	<LineChart {...{ stats, label: "Strategy" }}></LineChart>
</div>
