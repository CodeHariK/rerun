import { append, statusRender } from '/ui/panel.js';
import { fetchFile } from '/ui/monaco.js';
import { fileTree, generateFileTreeHTML } from '/ui/tree.js';

const fileTreeContainer = document.getElementById('file-tree');

export let socket;
let retryCount = 0;
const maxRetries = 5;

export let oriIrameUrl = 'http://localhost:9753/redirect';
let iframeUrl = oriIrameUrl;
export function setIframeUrl(url) {
    iframeUrl = url
}

const storedIframeUrl = localStorage.getItem('iframeUrl');
if (storedIframeUrl) {
    iframeUrl = storedIframeUrl
}

let pageurl = iframeUrl.replace(oriIrameUrl, "")
document.getElementById('pageUrlInput').value = pageurl;

async function checkServerStatus() {
    let attempts = 0;
    let interval = 100;

    const tryRequest = async () => {
        try {
            const response = await fetch(iframeUrl, { method: 'HEAD' });
            if (response.status == 200) {
                console.log('Server is alive');
                let iframeElement = document.getElementById('contentFrame')

                iframeElement.onload = function () {
                    Array.from(document.getElementById('contentFrame').contentDocument.getElementsByTagName("a"))
                        .forEach((a) => {
                            a.href = a.href.replace(9753, "9753/redirect")
                            a.addEventListener("click", (e) => {
                                iframeUrl = a.href
                                localStorage.setItem('iframeUrl', iframeUrl);
                                let pageurl = iframeUrl.replace(oriIrameUrl, "")
                                document.getElementById('pageUrlInput').value = pageurl;
                            })
                        })
                }

                iframeElement.src = iframeUrl

                return;
            }
        } catch (error) {
            // Handle errors here, if needed
        }

        attempts++;
        if (attempts < 8) {
            interval *= 2;
            setTimeout(tryRequest, interval);
        } else {
            console.log('Server did not respond, Attempt:', attempts);
        }
    };

    tryRequest();
}

connect()
function connect() {
    socket = new WebSocket("ws://localhost:9753/ws");

    socket.onopen = function (event) {
        console.log("Connected to WebSocket spider.");
        socket.send("SPIDER:PWD")
        retryCount = 0;
        checkServerStatus();
    };

    socket.onmessage = function (event) {
        let message = event.data

        if (message.startsWith("SPIDER:PWD:")) {
            fileTreeContainer.innerHTML = ""
            fileTreeContainer.appendChild(
                generateFileTreeHTML(
                    JSON.parse(
                        message.replace("SPIDER:PWD:", "").trim())));

            fileTree('file-tree');
            return
        }

        console.log(message)
        if (message.startsWith("SPIDER:ReRun")) {
            console.log(iframeUrl)
            checkServerStatus();
            statusRender(message)
            fetchFile(window.currentFile)
        }

        append(message, "SPIDER:Console:Output:", terminal, false)
        append(message, "SPIDER:Console:Error:", terminal, false)
        append(message, "SPIDER:Logs:Output:", logsPane, false)
        append(message, "SPIDER:Logs:Error:", logsPane, false)
    };

    socket.onclose = function (event) {
        socket.close()
        console.log("Disconnected from Spider." + event);
        if (retryCount < maxRetries) {
            const retryDelay = Math.pow(3, retryCount) * 400;
            console.log('Retrying in ' + retryDelay + 'ms...');
            setTimeout(connect, retryDelay);
            retryCount++;
        } else {
            console.log('Max retries reached. No further attempts will be made.');
        }
    };

    socket.onerror = function (event) {
        console.error('Spider error:', event);
        socket.close()
    };
}
