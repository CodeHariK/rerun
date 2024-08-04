import { socket, oriIrameUrl, setIframeUrl } from '/ui/ws.js';

const statusPane = document.getElementById('statusPane');
const pageUrlInputElement = document.getElementById('pageUrlInput');
const commandInputElement = document.getElementById('commandInput');
const terminal = document.getElementById('terminal');
const logsPane = document.getElementById('logsPane');

const logsUrl = 'http://localhost:9753/logs/';

export function statusRender(message) {
    logsPane.innerHTML = ""
    statusPane.innerHTML = ""

    let matches = message.match(/\d+/);
    if (matches) {
        let num = Number(matches[0]);

        for (let i = 1; i <= num; i++) {
            const button = document.createElement('button');
            button.innerText = i;
            button.addEventListener('click', async function (event) {
                try {
                    const response = await fetch(logsUrl + i, { method: 'GET' });
                    let body = await response.text()
                    if (response.ok) {

                        const logs = JSON.parse(body);
                        const logsPane = document.getElementById('logsPane');
                        logsPane.innerHTML = logs.map(log => {
                            if (log.log.includes("Error")) {
                                return "<p class='Output Error'>" + log.log.replace("Error:", "").trim() + "</p>"
                            } else {
                                return "<p class='Output'>" + log.log.replace("Output:", "").trim() + "</p>"
                            }
                        }).join('');
                    }
                } catch (error) { }
            });
            statusPane.appendChild(button);
        }
    }
}

pageUrlInputElement.addEventListener('keypress', function (event) {
    if (event.key === 'Enter') {
        let pageurl = pageUrlInputElement.value.replace(oriIrameUrl, "")
        let iframeUrl = oriIrameUrl + pageurl
        setIframeUrl(iframeUrl)
        document.getElementById('contentFrame').src = iframeUrl
        localStorage.setItem('iframeUrl', iframeUrl);
    }
});

commandInputElement.addEventListener('keypress', function (event) {
    if (event.key === 'Enter') {
        const command = commandInputElement.value;
        if (socket.readyState === WebSocket.OPEN) {
            append(command, "", terminal, false)
            socket.send("SPIDER:Console:" + command);
            commandInputElement.value = '';
        } else {
            console.log("WebSocket is not open.");
        }
    }
});

export function append(message, code, box, clear) {
    if (clear) {
        box.innerHTML = ""
    }
    if (message.startsWith(code)) {
        if (code.includes("Error")) {
            box.innerHTML += "<p class='Output Error'>" + message.replace(code, "").trim() + "</p>"
        } else {
            box.innerHTML += "<p class='Output'>" + message.replace(code, "").trim() + "</p>"
        }
    }
    box.scrollTop = box.scrollHeight;
}

function cancelTerminal() {
    if (socket.readyState === WebSocket.OPEN) {
        socket.send("SPIDER:Console:Cancel");
        console.log("SPIDER:Console:Cancel");
    } else {
        console.log("WebSocket is not open.");
    }
}
window.cancelTerminal = cancelTerminal;
