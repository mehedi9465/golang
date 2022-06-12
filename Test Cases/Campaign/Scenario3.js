const axios = require("axios")
const { campaignSearch, campaignPopulatedSearch } = require("./Campaign")

// Scenario 3
// In this scenario the system will check
// the search by name and  status

// Expected Outcome

module.exports.scenario3 = async (nameSearchObj, statusSearchObj) => {
    const searchtState = await campaignSearch(nameSearchObj)
    console.log("Search: \n", searchtState?.result?.length);
    
    const populatedSearchState = await campaignPopulatedSearch(statusSearchObj)
    console.log("\nPopulated Search: \n",populatedSearchState?.results?.length);
}
