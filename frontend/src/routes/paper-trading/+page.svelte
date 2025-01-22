<script lang="ts">
    const months = ["JAN", "FEB", "MAR", "APR", "MAY", "JUN", "JUL", "AUG", "SEP", "OCT", "NOV", "DEC"]
    const start = new Date(2001, 0, 13)
    const end = new Date(2001, 0, 21)
    let currDate = new Date(2001, 0, 13)
    let curr = formatDate(currDate)

    function formatDate(date: Date): string {
        const day = String(date.getDate()).padStart(2, "0")
        const month = months[date.getMonth()]
        const year = date.getFullYear()
        return `${day} ${month} ${year}`
    }

    function nextDay() {
        if (currDate < end) {
            currDate.setDate(currDate.getDate() + 1)
            curr = formatDate(currDate)
        }
    }

    function previousDay() {
        if (currDate > start) {
            currDate.setDate(currDate.getDate() - 1)
            curr = formatDate(currDate)
        }
    }
</script>

<div class="w-full flex flex-col items-center justify-center h-screen">
    <p class="text-3xl mb-4">Date: {curr}</p>
    <div class="flex gap-2">
        <button on:click={previousDay} class="btn variant-ghost-primary" disabled={curr === formatDate(start)}>Previous</button>
        <button on:click={nextDay} class="btn variant-ghost-primary" disabled={curr === formatDate(end)}>Next</button>
    </div>
</div>