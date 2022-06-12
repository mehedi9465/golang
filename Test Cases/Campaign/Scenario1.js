const axios = require("axios")
const { apiBaseUrl } = require("../BaseURL")
const {createNewCampaign, modifyCampaign, deleteCampaign} = require("./Campaign")

// Scenario 1
// First of all a Campaign will be inserted.
// Then using the inserted id it will perform 
// delete and modify operation

// Expected outcome:
// Campaign has been created successfully!
// Campaign has been modified successfully!
// Campaign has been Deleted successfully!

const insertionData = {
    "name": "test 1",
    "desc": "Here is the desc",
    "type": "Online",
    "status": "Draft",
    "online": {
        "channelref": "624bd87694bdd567de9463cd",
        // "files": ["data:text/plain;base64,aGVsbG8gd29ybGQYTakhi=","data:text/plain;base64,aGVsbG8gd29ybGQYR="],
        "files": [],
        "filesext": [],
        "others": ""
    },
    "expectedstartingdate": "2022-03-01T05:39:30.114+00:00",
    "expectedendingdate": "2022-04-03T06:39:30.115+00:00",
    "actualstartingdate": "2022-04-16T07:39:30.116+00:00",
    "actualendingdate": "2022-05-26T08:39:30.117+00:00",
    "products":["6242d300070b9b58c39a20ed"],
    "ownerref": "62455a1c6e66a71de077af43",
    "team": ["62453ba56c02f367b298a7c7", "6242d23c070b9b58c39a20ec"]
}

const modificationData = {
        "name": "test Modified 1",
        "desc": "Here is the Modified desc",
        "type": "Event" ,
        "status": "Review",
        "event": {
            "place": "zoom",
            // "files": ["Public/Files/file_files_545_1649821654390.txt","Public/Files/file_files_527_1649821654443.csv","data:text/plain;base64,aGVsbG8gd29ybGQYR=", "data:text/plain;base64,aGVsbG8gd29ybGQYR="],
            "files": [],
            "filesext": [],
            "others": ""
        },
        "expectedstartingdate": "2022-03-01T05:39:30.114+00:00",
        "expectedendingdate": "2022-04-03T06:39:30.115+00:00",
        "actualstartingdate": "2022-04-16T07:39:30.116+00:00",
        "actualendingdate": "2022-05-26T08:39:30.117+00:00",
        "products":["6242d300070b9b58c39a20ed"],
        "team": ["62453ba56c02f367b298a7c7", "6242d23c070b9b58c39a20ec", "62455b40e5df8f1ec8441aa7"]
    }

module.exports.scenario1 = async (insertionData) => {
    const insertState = await createNewCampaign(insertionData)
    console.log(insertState.message);

    const modifyState = await modifyCampaign(insertState?.result?.InsertedID, modificationData)
    console.log(modifyState.message);

    const deleteState = await deleteCampaign(insertState?.result?.InsertedID)
    console.log(deleteState.message);
}