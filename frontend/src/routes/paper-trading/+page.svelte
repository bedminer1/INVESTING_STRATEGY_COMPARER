<script lang="ts">
    let curr = "13 JAN 2001"
    let end = "13 JAN 2020"

    function nextDay() {
        const [day, month, year] = curr.split(" ")
        const months = ["JAN", "FEB", "MAR", "APR", "MAY", "JUN", "JUL", "AUG", "SEP", "OCT", "NOV", "DEC"]
        const monthIndex = months.indexOf(month)
        const currDate = new Date(Number(year), monthIndex, Number(day))

        // Increment by one day
        currDate.setDate(currDate.getDate() + 1)

        // Format the date back to "DD MON YYYY"
        const nextDay = String(currDate.getDate()).padStart(2, "0")
        const nextMonth = months[currDate.getMonth()]
        const nextYear = currDate.getFullYear()
        const nextDate = `${nextDay} ${nextMonth} ${nextYear}`

        // Check if the new date is within bounds
        const endDate = new Date(Number(end.split(" ")[2]), months.indexOf(end.split(" ")[1]), Number(end.split(" ")[0]));
        if (currDate <= endDate) {
            curr = nextDate
        }
    }
</script>

<div class="w-full flex flex-col items-center justify-center h-screen">
    <p class="text-3xl mb-4">Date: {curr}</p>
    <button on:click={nextDay} class="btn variant-ghost-primary">Next</button>
</div>