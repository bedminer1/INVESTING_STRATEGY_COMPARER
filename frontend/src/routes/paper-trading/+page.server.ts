export const load = async ({ url }) => {
    try {   
        const start = url.searchParams.get("start") || "2015_01_01"
        const end = url.searchParams.get("end") || "2020_01_01"
        const response = await fetch(`http://localhost:4000/prices?start=${start}&end=${end}`)

        if (!response.ok) {
            throw new Error(`Error fetching data: ${response.statusText}`)
        }
        const data: { price_data: PriceRecord[] } = await response.json()

        return {
            records: data.price_data,
        }
    } catch (err) {
        console.error("Error in load function: ", err)

        return {
            weeklyRecords: [],
        }
    }
}