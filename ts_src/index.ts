// transpiling from TypeScript
const acctIn = <HTMLInputElement>document.querySelector("#acct")
const amtIn = <HTMLInputElement>document.querySelector("#amt")
const depBtn = <HTMLButtonElement>document.querySelector("#depBtn")
const witBtn = <HTMLButtonElement>document.querySelector("#witBtn")
const statusDiv = <HTMLDivElement>document.querySelector("#status")

depBtn.addEventListener("click", () => {
    sendAmt(parseInt(amtIn.value))
})

witBtn.addEventListener("click", () => {
    sendAmt(parseInt(amtIn.value) * -1)
})


// sendAmt posts the deposit/withdrawl amount to the server api
function sendAmt(n: number) {
    fetch("/api", {
        method: "POST",
        body: JSON.stringify({ amt: n })
    }).then((res) => {
        if (res.status == 201) {
            // reset value
            amtIn.value = null
        }
    }).catch((e) => {
        statusDiv.innerText = "error: " + e
    })
}

const ws = new WebSocket("ws://" + window.location.host + "/socket");

ws.onopen = function () {
    amtIn.disabled = false;
    depBtn.disabled = false;
    witBtn.disabled = false;
    statusDiv.innerText = "value is in sync\n[websocket open]"
    ws.send(ip)
}

ws.onmessage = function (e) {
    statusDiv.innerText = "value is in sync\n[websocket active]"
    acctIn.value = e.data
    return false
}

ws.onerror = function (event) {
    statusDiv.innerText = "error: " + event
}

ws.onclose = function () {
    statusDiv.innerText = "value might be out of sinc\n[websocket closed]"
}

var ip

function getIP(jsn){
    ip = jsn.ip
}