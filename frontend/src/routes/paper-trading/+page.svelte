<script lang="ts">
    export let data: {
        records: PriceRecord[]
    }
   const records = data.records
   let displayedRecords: PriceRecord[] = []
   records.forEach(record => {
    record.Date = new Date(record.Date)
   })
   let recordIndex = 0
   
   const months = ["JAN", "FEB", "MAR", "APR", "MAY", "JUN", "JUL", "AUG", "SEP", "OCT", "NOV", "DEC"]
   const start = new Date(2015, 0, 27)
   const end = new Date(2021, 1, 2)
   let currDate = new Date(2015, 0, 27)
   let curr = formatDate(currDate)
   

    for (let record of records) {
        displayedRecords.push(record)
        if (record.Date > currDate) {
            break
        }
        recordIndex++
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
            displayedRecords.push(records[recordIndex])
            console.log(displayedRecords.slice(-3))
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
            displayedRecords.pop()
            console.log(displayedRecords.slice(-3))
        }
    }

</script>

<div class="w-full flex flex-col items-center justify-center h-screen">
    <p class="text-3xl mb-4">Date: {curr}</p>
    <div class="flex gap-2 mb-3">
        <button on:click={previousDay} class="btn variant-ghost-primary" disabled={curr === formatDate(start)}>Previous</button>
        <button on:click={nextDay} class="btn variant-ghost-primary" disabled={curr === formatDate(end)}>Next</button>
    </div>
    <div class="flex justify-center items-center text-center">
        {JSON.stringify(displayedRecords.slice(-3))}
    </div>
</div>