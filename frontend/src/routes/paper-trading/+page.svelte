<script lang="ts">
    import LineChart from "$lib/components/LineChart.svelte"
    import { onDestroy } from "svelte";
    import Card from "$lib/components/Card.svelte";

    let records: PriceRecord[] = []
    let frequencyMode = "fast-paper-trading"

    async function fetchRecords() {
        try {
            const response = await fetch(`http://localhost:4000/${frequencyMode}`)
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

    const interval = setInterval(fetchRecords, 1000)
    onDestroy(() => {
        clearInterval(interval)
    })

    let displayedRecords: PriceRecord[] = []
    let recordIndex = 0
   
    const months = ["JAN", "FEB", "MAR", "APR", "MAY", "JUN", "JUL", "AUG", "SEP", "OCT", "NOV", "DEC"]
    let start = records.at(-30)?.Date as Date ?? new Date(2015, 0, 27)
    let currDate = records.at(-1)?.Date as Date ?? new Date(2015, 0, 27)
    let curr = formatDate(currDate)
    let windowLength = 30
    
    $: {
        displayedRecords = []
        start = records.at(-windowLength+1)?.Date as Date ?? new Date(2015, 0, 27)
        for (let record of records) {
            if (record.Date < start) {
                continue
            }
            displayedRecords.push(record)
            if (record.Date > currDate) {
                break
            }
            recordIndex++
        }
        currDate = records.at(-1)?.Date as Date ?? new Date(2015, 0, 27) 
        curr = formatDate(currDate)
    }

    function formatDate(date: Date): string {
        const day = String(date.getDate()).padStart(2, "0")
        const month = months[date.getMonth()]
        const year = date.getFullYear()
        return `${day} ${month} ${year}`
    }


    // BUYING AND SELLING
    let cash = 100000
    let position = 0
    let quantity: number
    $: currValue = displayedRecords.at(-1)?.Price!
    $: orderValue = quantity * currValue
    $: marketValue = position * currValue
    $: netWorth = cash + marketValue
    
    let orderModeIsBuy = true
    let popUpOpen = false
    $: newPosition = orderModeIsBuy ? position + quantity : position - quantity
    $: newCashBalance = orderModeIsBuy ? cash - orderValue : cash + orderValue

    function buy() {
        orderModeIsBuy = true
        popUpOpen = true
    }

    function sell() {
        orderModeIsBuy = false
        popUpOpen = true
    }

    function executeOrder() {
        if (orderModeIsBuy) {
            cash -= currValue * quantity
            marketValue += currValue * quantity
            position += quantity
        } else {
            cash += currValue * quantity
            marketValue -= currValue * quantity
            position -= quantity
        }
        popUpOpen = false
    }

    // PORTFOLIO HISTORY
    let showHistory = false
    let portfolioHistory: PriceRecord[] = []
    let displayedPortfolioHistory: PriceRecord[] = []
    let performance: number = 0
    $: {
        if (displayedRecords && !isNaN(netWorth)) {
            portfolioHistory = [...portfolioHistory, ({
                Price: netWorth,
                Date: displayedRecords.at(-1)?.Date!
            })]
            displayedPortfolioHistory = portfolioHistory.slice(-windowLength)
            performance = ((displayedPortfolioHistory.at(-1)?.Price! - displayedPortfolioHistory[1]?.Price!) / displayedPortfolioHistory.at(-1)?.Price!) * 100
        }
    }

    // TODO

    // show portfolio performance history, metrics of gains and losses
    // allow for custom orders to execute when price is met
</script>

<div class="w-full flex flex-col items-center justify-center h-screen">
    <p class="text-3xl mb-4">Date: {curr}</p>
    <div class="flex gap-2 mb-3">
        <button on:click={() => windowLength = 7} class="btn variant-ghost-primary">Week</button>
        <button on:click={() => windowLength = 30} class="btn variant-ghost-primary">Month</button>
        <button on:click={() => windowLength = 365} class="btn variant-ghost-primary">Year</button>
        <button on:click={() => windowLength = 365*5} class="btn variant-ghost-primary">5 Years</button>
        <button on:click={() => windowLength = 365*5*5} class="btn variant-ghost-primary">All</button>
    </div>
    <div class="flex flex-col justify-center items-center text-center w-full">
        <div class="w-full">
            {#if showHistory}
            <LineChart {...{ stats: [{
                label: "Portfolio Value", 
                data: displayedPortfolioHistory.map(record => record.Price),
                borderColor: 'rgba(54, 162, 235, 1)',
                backgroundColor: 'rgba(54, 162, 235, 0.2)'
            }], label: "" , xAxisLabels: displayedPortfolioHistory.map(record => (record.Date as Date).toLocaleDateString("en-GB"))}}/>
            {:else}
            <LineChart {...{ stats: [{
                label: "S&P 500", 
                data: displayedRecords.map((record) => record.Price),
                borderColor: 'rgba(54, 162, 235, 1)',
                backgroundColor: 'rgba(54, 162, 235, 0.2)'
            }], label: "" , xAxisLabels: displayedRecords.map(record => (record.Date as Date).toLocaleDateString("en-GB"))}}/>
            {/if}
        </div>
    </div>
    <div class="flex gap-2 mb-3 w-1/2 justify-center">
        <input type="number" class="input w-1/5" bind:value={quantity} placeholder="quantity">
        <button on:click={buy} class="btn variant-ghost-primary" disabled={!quantity}>Buy</button>
        <button on:click={sell} class="btn variant-ghost-primary" disabled={!quantity}>Sell</button>
        <button on:click={() => {showHistory = !showHistory}} class="btn variant-ghost-primary">Performance</button>
    </div>
    <div hidden={!popUpOpen} class="w-full mb-3">
        <div class="flex flex-col justify-center items-center">
            <div class="w-1/4 flex justify-between mb-4 border-2 border-dotted p-4">
                <div class="flex flex-col w-3/4">                    
                    <p class="w-full">Order Value: {orderValue.toFixed(2)}</p>
                    <p class="w-full">New Position: {newPosition.toFixed(2)}</p>
                    <p class="w-full">Cash Balance: {newCashBalance.toFixed(2)}</p>
                </div>
                <button class="w-20 btn variant-ghost-error h-12" on:click={() => {popUpOpen = false}}>Cancel</button>
            </div>
            <button on:click={executeOrder} class="btn variant-ghost-primary w-1/4">Submit Order</button>
        </div>
    </div>
    <div class="flex gap-4">
        <Card 
            {...{
                title: "Net Worth",
                body: netWorth.toFixed(2),
                subtitle: "Total Assets Value",
                icon: "&#9814;"
            }}
        />
        <Card 
            {...{
                title: "Position",
                body: position.toFixed(2),
                subtitle: "No. of Stocks Owned",
                icon: "&#9814;"
            }}
        />
        <Card 
            {...{
                title: "Cash",
                body: cash.toFixed(2),
                subtitle: "Liquid Cash Available",
                icon: "&#9814;"
            }}
        />
        <Card 
            {...{
                title: "Performace",
                body: performance.toFixed(2) + "%",
                subtitle: "Performance over Time",
                icon: "&#9814;"
            }}
        />
    </div>
</div>