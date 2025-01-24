<script lang="ts">
    import LineChart from "$lib/components/LineChart.svelte"
    import { onDestroy } from "svelte";
    let records: PriceRecord[] = []

    async function fetchRecords() {
        try {
            const response = await fetch("http://localhost:4000/paper-trading")
            if (!response.ok) {
                console.error("error fetching data: ", response.statusText)
            }
            const data: { records: PriceRecord[] } = await response.json()
            records = data.records
            records.forEach(record => {
                record.Date = new Date(record.Date)
            })
        } catch (error) {
            console.error("error updated records")
        }
    }

    const interval = setInterval(fetchRecords, 2000)
    onDestroy(() => {
        clearInterval(interval)
    })

    let displayedRecords: PriceRecord[] = []
    let recordIndex = 0
   
    const months = ["JAN", "FEB", "MAR", "APR", "MAY", "JUN", "JUL", "AUG", "SEP", "OCT", "NOV", "DEC"]
    let start = records.at(-30)?.Date as Date ?? new Date(2015, 0, 27)
    let end = records.at(-1)?.Date as Date ?? new Date(2021, 1, 2)
    let currDate = records.at(-1)?.Date as Date ?? new Date(2015, 0, 27)
    let curr = formatDate(currDate)
    let dates: string[] = []
    let windowLength = 30
    
    $: {
        displayedRecords = []
        dates = []
        start = records.at(-windowLength+1)?.Date as Date ?? new Date(2015, 0, 27)
        for (let record of records) {
            if (record.Date < start) {
                continue
            }
            displayedRecords.push(record)
            dates.push((record.Date as Date).toLocaleDateString("en-GB"))
            if (record.Date > currDate) {
                break
            }
            recordIndex++
        }
        end = records.at(-1)?.Date as Date ?? new Date(2021, 1, 2)
        currDate = records.at(-1)?.Date as Date ?? new Date(2015, 0, 27) 
        curr = formatDate(currDate)
    }

    function formatDate(date: Date): string {
        const day = String(date.getDate()).padStart(2, "0")
        const month = months[date.getMonth()]
        const year = date.getFullYear()
        return `${day} ${month} ${year}`
    }

    function nextDay() {
        if (currDate >= end) {
            return
        }
        currDate.setDate(currDate.getDate() + 1)
        curr = formatDate(currDate)
        if (records[recordIndex].Date <= currDate) {
            recordIndex++
            displayedRecords = [...displayedRecords, records[recordIndex]]
            dates = [...dates, (records[recordIndex].Date as Date).toLocaleDateString("en-GB")]
            // console.log(displayedRecords.slice(-3))
        }
    }

    function previousDay() {
        if (currDate <= start) {
            return
        }
        currDate.setDate(currDate.getDate() - 1)
        curr = formatDate(currDate)
        if (recordIndex > 0 && records[recordIndex - 1].Date > currDate) {
            recordIndex--
            displayedRecords = displayedRecords.slice(0, -1)
            dates = dates.slice(0, -1)
            // console.log(displayedRecords.slice(-3))
        }
    }

    $: graphPriceData = displayedRecords.map((record) => record.Price)



    // BUYING AND SELLING
    let cash = 100000
    let position = 0
    let quantity: number
    $: currValue = displayedRecords.at(-1)?.Price!
    $: marketValue = position * currValue
    $: netWorth = cash + marketValue

    // $:console.log(netWorth, cash, marketValue, position)

    function buy() {
        cash -= currValue * quantity
        marketValue += currValue * quantity
        position += quantity
    }

    function sell() {
        cash += currValue * quantity
        marketValue -= currValue * quantity
        position -= quantity
    }
</script>

<div class="w-full flex flex-col items-center justify-center h-screen">
    <p class="text-3xl mb-4">Date: {curr}</p>
    <!-- <div class="flex gap-2 mb-3">
        <button on:click={previousDay} class="btn variant-ghost-primary" disabled={curr === formatDate(start)}>Previous</button>
        <button on:click={nextDay} class="btn variant-ghost-primary" disabled={curr === formatDate(end)}>Next</button>
    </div> -->
    <div class="flex gap-2 mb-3">
        <button on:click={() => windowLength = 7} class="btn variant-ghost-primary">Week</button>
        <button on:click={() => windowLength = 30} class="btn variant-ghost-primary">Month</button>
        <button on:click={() => windowLength = 365} class="btn variant-ghost-primary">Year</button>
        <button on:click={() => windowLength = 365*5} class="btn variant-ghost-primary">5 Years</button>
        <button on:click={() => windowLength = 365*5*5} class="btn variant-ghost-primary">All</button>
    </div>
    <div class="flex justify-center items-center text-center w-full">
        <div class="w-full">
            <LineChart {...{ stats: [{
                label: "S&P 500", 
                data: graphPriceData,
                borderColor: 'rgba(54, 162, 235, 1)',
                backgroundColor: 'rgba(54, 162, 235, 0.2)'
            }], label: "" , xAxisLabels: dates}}/>
        </div>
    </div>
    <div class="flex gap-2 mb-3">
        <input type="number" class="input" bind:value={quantity} placeholder="quantity">
        <button on:click={buy} class="btn variant-ghost-primary">Buy</button>
        <button on:click={sell} class="btn variant-ghost-primary">Sell</button>
    </div>
</div>