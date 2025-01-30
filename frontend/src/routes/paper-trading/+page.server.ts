export const load = async ({ fetch }) => {
    try {
        const response = await fetch("http://localhost:4000/user-stats?user_id=bed")
        if (!response.ok) {
            console.error("error fetching response")
        }

        const data = await response.json()
        const userID = data.user_info.user_id
        const cash = data.user_info.cash
        const position = data.user_info.position
        let portfolioHistory: PriceRecord[] = []
        for (let record of data.user_info.net_worth_history) {
            portfolioHistory.push({
                Price: record.price,
                Date: new Date(record.date)
            })
        }
        console.log("Fetched Cash:", cash, 
            "\nFetched Position: ", position)

        return {
            userID,
            cash,
            position,
            portfolioHistory
        }

    } catch (error) {
        console.error("error fetching data", error)
    }
}