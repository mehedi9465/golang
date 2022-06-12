const axios = require("axios")
const { apiBaseUrl } = require("../BaseURL")

module.exports.createNewCampaign = async (insertionData) => {

    const res = await axios.post(`${apiBaseUrl}/campaign/new`, insertionData)
    const data = await res.data

    if (data.success) {
        return data
    }
        return data
}

module.exports.modifyCampaign = async (campaignID, modificationData) => {

    const res = await axios.put(`${apiBaseUrl}/campaign/modify/${campaignID}`, modificationData)
    const data = await res.data

    if (data.success) {
        // console.log(res.data);
        return data
    }
        return data
}

module.exports.campaignSearch= async (searchObj) => {
    const res = await axios.post(`${apiBaseUrl}/campaign/get_all`, searchObj)
    const data = await res.data

     if (data.success) {
        return data
    }
        return data
}

module.exports.campaignPopulatedSearch= async (searchObj) => {
    const res = await axios.post(`${apiBaseUrl}/campaign/get_all_populated`, searchObj)
    const data = await res.data

     if (data.success) {
        return data
    }
        return data
}

module.exports.deleteCampaign = async (campaignID) => {
    const res = await axios.delete(`${apiBaseUrl}/campaign/delete/${campaignID}`)
    const data = await res.data

     if (data.success) {
        return data
    }
        return data
}

//                       Status Change Methods

const DraftToReview = async (id) => {
    const res = await axios.put(`${apiBaseUrl}/campaign/set_status/${id}/Review`)
    return res.data
}

const ReviewToConfirmed = async (id) => {
    const res = await axios.put(`${apiBaseUrl}/campaign/set_status/${id}/Confirmed`)
    return res.data
}

const ConfirmedToStarted = async (id) => {
    const res = await axios.put(`${apiBaseUrl}/campaign/set_status/${id}/Started`)
    return res.data
}

const StartedToClosed = async (id) => {
    const res = await axios.put(`${apiBaseUrl}/campaign/set_status/${id}/Closed`)
    return res.data
}

const StartedToCanceled = async (id) => {
    const res = await axios.put(`${apiBaseUrl}/campaign/set_status/${id}/Canceled`)
    return res.data
}

module.exports.DraftToClosedStatus = async (id) => {
    let response = await DraftToReview(id)
    .then(res => console.log("Review Success", res))
    .catch(error => console.log("Review Error", error.response.data))

    response = await ReviewToConfirmed(id)
    .then(res => console.log("Confirmed Success", res))
    .catch(error => console.log("Confirmed Error", error.response.data))

    response = await ConfirmedToStarted(id)
    .then(res => console.log("Started Success", res))
    .catch(error => console.log("Started Error", error.response.data))

    response = await StartedToClosed(id)
    .then(res => console.log("Closed Success", res))
    .catch(error => console.log("Closed Error", error.response.data))
}

module.exports.DraftToCanceledStatus = async (id) => {
    let response = await DraftToReview(id)
    .then(res => console.log("Review Success", res))
    .catch(error => console.log("Review Error", error.response.data))

    response = await ReviewToConfirmed(id)
    .then(res => console.log("Confirmed Success", res))
    .catch(error => console.log("Confirmed Error", error.response.data))

    response = await ConfirmedToStarted(id)
    .then(res => console.log("Started Success", res))
    .catch(error => console.log("Started Error", error.response.data))

    response = await StartedToCanceled(id)
    .then(res => console.log("Canceled Success", res))
    .catch(error => console.log("Canceled Error", error.response.data))
}



