<script lang="ts">
    import LineChart from "$lib/components/LineChart.svelte"
    import { onMount, onDestroy } from "svelte";
    import Card from "$lib/components/Card.svelte";

    export let data: {
        userID: string
        cash: number,
        position: number
        portfolioHistory: PriceRecord[]
    }
    let records: PriceRecord[] = []
    let frequencyMode = "fast-paper-trading"

    async function fetchRecords() {
        try {
            const response = await fetch(`http://localhost:4000/${frequencyMode}`)
            if (!response.ok) {
                console.error("error fetching data: ", response.statusText)
            }
            const data: { records: PriceRecord[] } = await response.json()
            if (!data.records) {
                console.error("No records found in response")
            }
            let unparsedRecords = data.records
            unparsedRecords.forEach(record => {
                record.Date = new Date(record.Date)
            })
            records = unparsedRecords
        } catch (error) {
            console.error("error updated records")
        }
    }

    onMount(() => {
        fetchRecords()

        const interval = setInterval(fetchRecords, 1000)

        onDestroy(() => {
            clearInterval(interval)
        })
    })

    let displayedRecords: PriceRecord[] = []

   
    const months = ["JAN", "FEB", "MAR", "APR", "MAY", "JUN", "JUL", "AUG", "SEP", "OCT", "NOV", "DEC"]
    let currDate = records.at(-1)?.Date as Date ?? new Date(2015, 0, 27)
    let curr = formatDate(currDate)
    let windowLength = 30

    function formatDate(date: Date): string {
        const day = String(date.getDate()).padStart(2, "0")
        const month = months[date.getMonth()]
        const year = date.getFullYear()
        return `${day} ${month} ${year}`
    }


    // ASSETS 
    let cash = data.cash ?? 100000
    let position = data.position ?? 0
    let quantity: number
    $: currValue = displayedRecords.at(-1)?.Price ?? 0
    $: orderValue = quantity * targetPrice
    $: marketValue = position * currValue
    $: netWorth = cash + marketValue

    // PORTFOLIO HISTORY
    let showHistory = false
    let showPerformanceHistory = false
    let portfolioHistory: PriceRecord[] = data.portfolioHistory ?? []
    let displayedPortfolioHistory: PriceRecord[] = []
    let performance: number = 0
    let performanceHistory: PriceRecord[] = []
    let displayedPerformanceHistory: PriceRecord[] = []
    
    // UPDATE stock data and portfolio stats
    $: {
        displayedRecords = records.slice(-windowLength)
        currDate = records.at(-1)?.Date ?? new Date(2015, 0, 27) 
        curr = formatDate(currDate)
    }

    // UPDATE portfolio history
    $: {
        if (displayedRecords.length > 0 && !isNaN(netWorth)) {
            const newRecord = {
                Price: netWorth,
                Date: currDate
            }
            portfolioHistory = [...portfolioHistory, newRecord]
            displayedPortfolioHistory = portfolioHistory.slice(-windowLength)
        }

        saveMetric()
    }

    // UPDATE performance history
    $: {
        if (portfolioHistory.length > 1) {
            const latestPrice = portfolioHistory.at(-1)?.Price ?? 0
            const initialPrice = portfolioHistory[0]?.Price ?? 1
            performance = ((latestPrice - initialPrice) / initialPrice) * 100
            performanceHistory = [...performanceHistory, { Price: performance, Date: currDate }]
            displayedPerformanceHistory = performanceHistory.slice(-windowLength)
        }
    }

    $: {
        for (let order of orders) {
            if ((order.IsBuyOrder && order.Price >= currValue) 
            || (!order.IsBuyOrder && order.Price <= currValue)) {
                executeOrder(order)
                let index = orders.indexOf(order)
                orders.splice(index, 1)
            }
        }
    }

    // CUSTOM ORDERS
    let orders: OrderRecord[] = []
    let orderModeIsBuy = true
    let popUpOpen = false
    let targetPrice: number
    $: newPosition = orderModeIsBuy ? position + quantity : position - quantity
    $: newCashBalance = orderModeIsBuy ? cash - orderValue : cash + orderValue

    function buy() {
        targetPrice = currValue
        orderModeIsBuy = true
        popUpOpen = true
    }

    function sell() {
        targetPrice = currValue
        orderModeIsBuy = false
        popUpOpen = true
    }

    function addOrder() {
        let order = {
            Price: targetPrice,
            Date: currDate,
            Quantity: quantity,
            IsBuyOrder: orderModeIsBuy
        }
        orders = [...orders, order]
        console.log("Added Order: ", order)
        popUpOpen = false
    }

    function executeOrder(order: OrderRecord) {
        const orderValue = currValue * order.Quantity
        if (order.IsBuyOrder) {
            cash -= orderValue
            marketValue += orderValue
            position += order.Quantity
        } else {
            cash += orderValue
            marketValue -= orderValue
            position -= order.Quantity
        }
        console.log("Executed order: ", order, "Curr Value: ", currValue)
    }

    const USERID = "bed"

    async function saveMetric() {
        const userMetrics = {
            user_id: USERID,
            cash,
            position,
            net_worth_history: portfolioHistory.map(record => {
                return {
                    price: record.Price,
                    date: record.Date,
                    user_id: USERID
                }
            })
        }

        try {
            const response = await fetch("http://localhost:4000/save-metrics", {
                method: "POST",
                headers: {
                    "Content-Type": "application/json",
                },
                body: JSON.stringify(userMetrics)
            })
        } catch (error) {
            console.error("Error saving metrics")
        }
    }

    // TODO

    // add feature to cancel orders, view orders
    // add state to the app that's persisted beyond refresh
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
        <div class="w-full p-4">
            {#if !showHistory}
            <LineChart {...{ stats: [{
                label: "S&P 500", 
                data: displayedRecords.map((record) => record.Price),
                borderColor: 'rgba(54, 162, 235, 1)',
                backgroundColor: 'rgba(54, 162, 235, 0.2)'
            }], label: "Value (USD)" , xAxisLabels: displayedRecords.map(record => (record.Date as Date).toLocaleDateString("en-GB"))}}/>
            {:else if showPerformanceHistory}
            <LineChart {...{ stats: [{
                label: "Portfolio Performance", 
                data: displayedPerformanceHistory.map(record => record.Price),
                borderColor: 'rgba(54, 162, 235, 1)',
                backgroundColor: 'rgba(54, 162, 235, 0.2)'
            }], label: "%" , xAxisLabels: displayedPerformanceHistory.map(record => (record.Date as Date).toLocaleDateString("en-GB"))}}/>
            {:else}
            <LineChart {...{ stats: [{
                label: "Portfolio Value", 
                data: displayedPortfolioHistory.map(record => record.Price),
                borderColor: 'rgba(54, 162, 235, 1)',
                backgroundColor: 'rgba(54, 162, 235, 0.2)'
            }], label: "Value (USD)" , xAxisLabels: displayedPortfolioHistory.map(record => (record.Date as Date).toLocaleDateString("en-GB"))}}/>
            {/if}
        </div>
    </div>
    <div class="flex gap-2 mb-3 w-2/3 justify-center">
        <button on:click={buy} class="btn variant-ghost-primary">Buy</button>
        <button on:click={sell} class="btn variant-ghost-primary">Sell</button>
        <button on:click={() => {showHistory = !showHistory}} class="btn variant-ghost-primary">Stock/Portfolio</button>
        <button on:click={() => {showPerformanceHistory = !showPerformanceHistory}} disabled={showHistory === false} class="btn variant-ghost-primary">Value/Performance</button>
    </div>
    <div hidden={!popUpOpen} class="w-full mb-3">
        <div class="flex flex-col justify-center items-center">
            <div class="w-1/2 flex justify-between mb-4 border-2 border-dotted p-4">
                <div class="flex flex-col w-3/4">
                    <input type="number" class="input" bind:value={quantity} placeholder="quantity">
                    <input type="number" class="input" bind:value={targetPrice} placeholder="order price">                    
                    <p class="w-full">Order Value: {orderValue.toFixed(2)}</p>
                    <p class="w-full">New Position: {newPosition.toFixed(2)}</p>
                    <p class="w-full">Cash Balance: {newCashBalance.toFixed(2)}</p>
                </div>
                <button class="w-20 btn variant-ghost-error h-12" on:click={() => {popUpOpen = false}}>Cancel</button>
            </div>
            <button on:click={addOrder} class="btn variant-ghost-primary w-1/4">Submit Order</button>
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
                title: "P&L",
                body: performance.toFixed(2) + "%",
                subtitle: "Performance over Time",
                icon: "&#9814;"
            }}
        />
    </div>
</div>