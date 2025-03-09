// Check if we're in an iframe using standard browser check
console.log("sub.js loaded");
if (window.self !== window.top) {
    // We are in an iframe, but we need to wait for it to load fully
    console.log('Current iframe URL:', window.location.href);

    // Function to collect and send data after content is loaded
    function collectAndSendData() {
        let currentUrl = window.location.href;
        chrome.runtime.sendMessage({ ping: true })
            .then(response => {
                let msg = {
                    sendBack: true,
                    data: document.body.innerHTML,
                    fromIframe: true,
                    url: currentUrl
                };
                return chrome.runtime.sendMessage(msg);
            })
            .catch(error => {
                console.log('Extension connection check failed:', error);
                handleDataLocally(data, currentUrl);
            });
    }

    if (document.readyState === 'complete') {
        setTimeout(collectAndSendData, 500); // Small delay even if complete
    } else {
        // Wait for the page to fully load before collecting data
        window.addEventListener('load', function () {
            // Add a delay to ensure dynamic content is also loaded
            setTimeout(collectAndSendData, 1500); // Increased timeout for more reliability
        });
    }
}