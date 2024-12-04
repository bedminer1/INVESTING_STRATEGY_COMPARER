

export const load = async ({ url }) => {
    try {   
        const start = url.searchParams.get("start") || "2014_01_01"
        const end = url.searchParams.get("end") || "2024_01_01"
        const response = await fetch(`http://localhost:4000/strategies?start=${start}&end=${end}`)
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