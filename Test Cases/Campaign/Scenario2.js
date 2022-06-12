const axios = require("axios")
const {createNewCampaign, DraftToCanceledStatus} = require("./Campaign")

// Scenario 2
// First scenario will insert a document
// in the database and then it will
// try to change the status Draft to Canceled

// Status Directions: 
// Draft => Review 
// Review => Confirmed 
// Confirmed => Started 
// Started => Canceled 

// Expected Outcome
// Campaign has been created successfully!
// Review Success {
//   message: 'Campaign Status has been modified successfully!',
//   success: true
// }
// Confirmed Success {
//   message: 'Campaign Status has been modified successfully!',
//   success: true
// }
// Started Success {
//   message: 'Campaign Status has been modified successfully!',
//   success: true
// }
// Canceled Success {
//   message: 'Campaign Status has been modified successfully!',
//   success: true
// }

module.exports.scenario2 = async (insertionData) => {
    const insertState = await createNewCampaign(insertionData)
    console.log(insertState.message, "\n");

    const changeStatus = await DraftToCanceledStatus(insertState?.result?.InsertedID)
}