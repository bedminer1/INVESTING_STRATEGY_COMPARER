

export const load = async () => {
    try {   
        const response = await fetch("http://localhost:4000/strategies?start=2014_01_01&end=2024_01_01")
        if (!response.ok) {
            throw new Error(`Error fetching data: ${response.statusText}`)
        }
        const data: { results: WeeklyRecords[] } = await response.json()

        return {
            results: data.results,
        }
    } catch (err) {
        console.error("Error in load function: ", err)

        return {
            weeklyRecords: []
        }
    }
}