export const load = async ({ url }) => {
    try {   
        const response = await fetch(`http://localhost:4000/paper-trading`)

        if (!response.ok) {
            throw new Error(`Error fetching data: ${response.statusText}`)
        }
        const data: { records: PriceRecord[] } = await response.json()

        return {
            records: data.records,
        }
    } catch (err) {
        console.error("Error in load function: ", err)

        return {
            weeklyRecords: [],
        }
    }
}