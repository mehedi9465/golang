const { scenario1 } = require("./Scenario1")
const { scenario2 } = require("./Scenario2")
const { scenario3 } = require("./Scenario3")

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

    
const nameSearchObj = {
    "name": "2",
    "nameisused": true
}

const statusSearchObj = {
    "status": "Draft",
    "statusisused": true
}

const main = async () => {
   await console.log("\n\n Scenario 1: \n\n");
  await  scenario1(insertionData);
   await console.log("\n\n Scenario 2: \n\n");
   await scenario2(insertionData);
   await console.log("\n\n Scenario 3: \n\n");
   await scenario3(nameSearchObj, statusSearchObj);
}

main()