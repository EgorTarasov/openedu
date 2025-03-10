// Main page script
console.log('Main page script loaded');

const api = "http://127.0.0.1:8080"

// Listen for messages from background script or iframes
chrome.runtime.onMessage.addListener(function (message, sender, sendResponse) {
    // Only process messages from iframes
    if (message.fromIframe) {
        // Create a unique key for this iframe content
        const storageKey = 'iframeContent_' + message.url;

        // Check if we've already processed this iframe URL
        const previouslySent = localStorage.getItem(storageKey);

        if (previouslySent) {
            console.log('Data from this iframe URL already sent to server:', message.url);
            // Send response to indicate we've already processed this
            sendResponse({ status: 'already_sent' });
        } else {
            // This is a new iframe content we haven't processed
            const payload = {
                url: message.url,
                data: message.data,
                title: window.document.URL
            }

            console.log('Sending new data from iframe URL:', message.url);

            fetch(`${api}/collect`, {
                method: "POST",
                body: JSON.stringify(payload),
                headers: {
                    "Content-Type": "application/json"
                }
            })
                .then(res => {
                    console.log('Server response:', res);
                    // Save to localStorage to prevent duplicate requests
                    localStorage.setItem(storageKey, JSON.stringify({
                        url: message.url,
                        timestamp: new Date().toISOString()
                    }));
                    // Send response to indicate successful processing
                    sendResponse({ status: 'sent_successfully' });
                })
                .catch(err => {
                    console.error('Error sending data to server:', err);
                    sendResponse({ status: 'error', error: err.message });
                });

            // Return true to indicate we'll send response asynchronously
            return true;
        }
    }

    return true;
});