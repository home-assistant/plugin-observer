async function updateLog() {
    const logElement = document.querySelector("#log")
    const connectionElement = document.querySelector("#connection")

    if (connectionElement.innerText == "Disconnected" && logElement.innerText == "") {
        logElement.innerText = "Loading Supervisor logs..."

        const logEntries = await fetch("/logs")
        logElement.innerText = await logEntries.text()
    }
}

window.addEventListener('DOMContentLoaded', () => {
    updateLog()
});