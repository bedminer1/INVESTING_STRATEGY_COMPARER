<script lang="ts">
    let curr = "13 JAN 2001"
    let start = "13 JAN 2001"
    let end = "21 JAN 2001"

    function nextDay() {
        const [day, month, year] = curr.split(" ")
        const months = ["JAN", "FEB", "MAR", "APR", "MAY", "JUN", "JUL", "AUG", "SEP", "OCT", "NOV", "DEC"]
        const monthIndex = months.indexOf(month)
        const currDate = new Date(Number(year), monthIndex, Number(day))

        // Increment
        currDate.setDate(currDate.getDate() + 1)

        // Format
        const nextDay = String(currDate.getDate()).padStart(2, "0")
        const nextMonth = months[currDate.getMonth()]
        const nextYear = currDate.getFullYear()
        const nextDate = `${nextDay} ${nextMonth} ${nextYear}`

        // Check bounds
        const endDate = new Date(Number(end.split(" ")[2]), months.indexOf(end.split(" ")[1]), Number(end.split(" ")[0]));
        if (currDate <= endDate) {
            curr = nextDate
        }
    }

    function previousDay() {
        const [day, month, year] = curr.split(" ")
        const months = ["JAN", "FEB", "MAR", "APR", "MAY", "JUN", "JUL", "AUG", "SEP", "OCT", "NOV", "DEC"]
        const monthIndex = months.indexOf(month)
        const currDate = new Date(Number(year), monthIndex, Number(day))

        // Decrement
        currDate.setDate(currDate.getDate() - 1)

        // Format
        const prevDay = String(currDate.getDate()).padStart(2, "0")
        const prevMonth = months[currDate.getMonth()]
        const prevYear = currDate.getFullYear()
        const prevDate = `${prevDay} ${prevMonth} ${prevYear}`

        // Check bounds
        const startDate = new Date(Number(start.split(" ")[2]), months.indexOf(start.split(" ")[1]), Number(start.split(" ")[0]));
        if (currDate >= startDate) {
            curr = prevDate
        }
    }
</script>

<div class="w-full flex flex-col items-center justify-center h-screen">
    <p class="text-3xl mb-4">Date: {curr}</p>
    <div class="flex gap-2">
        <button on:click={previousDay} class="btn variant-ghost-primary" disabled={curr === start}>Previous</button>
        <button on:click={nextDay} class="btn variant-ghost-primary" disabled={curr === end}>Next</button>
    </div>
</div>