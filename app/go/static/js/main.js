
document.addEventListener("DOMContentLoaded", function() {
    console.log('ready!');

    getData()
 });

function getData() {
    fetch('./ui/getData', {
                    method: "GET",
                    credentials: "include",
                    mode: "cors",
                })
                .then((response) => {
                    return response.json()
                })
                .then((json) => {
                    if (json.status && json.status != 200) {
                        console.log("Error getting data, ", json)
                    } else {
                        document.getElementById("latestValue").innerText = json[0].data.data
                    }
                })
                .catch((error) => {
                    console.log("Error getting data,  ", error)
                })
}