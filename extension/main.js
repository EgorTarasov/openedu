// Main page script
console.log('Main page script loaded');

// Listen for messages from background script or iframes
chrome.runtime.onMessage.addListener(function (message, sender, sendResponse) {
    if (typeof message === 'string') {
        console.log('Direct iframe message received:', message.data);
        alert('string received from iframe: ');
    } else if (message.fromIframe) {
        const last = localStorage.getItem('iframeContent_' + message.url);
        if (!last === null) {
            alert('data already send to server');
            return;
        } else if (message.fromIframe){
            const payload = {
                url: message.url,
                data: message
            };
            alert('sending new data from url: ', message.url);
            fetch("https://openedu.larek.tech/collect", {
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