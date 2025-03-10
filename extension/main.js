// Main page script
console.log('Main page script loaded');

const api = "http://127.0.0.1:8080"

// Listen for messages from background script or iframes
chrome.runtime.onMessage.addListener(function (message, sender, sendResponse) {

    if (message.fromIframe) {
        const last = localStorage.getItem('iframeContent_' + message.url);
        if (!(last === null)) {
            console.log('data already send to server');
            return;
        } else if (message.fromIframe){
            const payload = {
                url: message.url,
                data: message.data,
                title: window.document.URL
            }
            console.log('sending new data from url: ', payload);
            fetch(`${api}/collect`, {
                method: "POST",
                body: JSON.stringify(payload),
                headers: {
                    "Content-Type": "application/json"
                }
            }).then(res => {
                console.log(res);
                localStorage.setItem('iframeContent_' + message.url, JSON.stringify({
                    url: message.url,
                    content: message.data,
                    timestamp: new Date().toISOString()
                }));
            }).catch(err => {
                console.log(err);
            })
        }
    }

    return true;
});