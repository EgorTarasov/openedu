// Check if we're in an iframe using standard browser check
console.log("sub.js loaded");
if (window.self !== window.top) {
    // We are in an iframe, but we need to wait for it to load fully
    console.log('Current iframe URL:', window.location.href);

    // Function to collect and send data after content is loaded
    function collectAndSendData() {
        let currentUrl = window.location.href;

        // First check if the extension is available
        chrome.runtime.sendMessage({ ping: true })
            .then(response => {
                // Prepare the message with iframe content
                const problemDiv = document.querySelector('.problem');

                let msg = {
                    sendBack: true,
                    fromIframe: true,
                    problemDiv: problemDiv.innerHTML,
                };

                // Send the message and handle the response
                return chrome.runtime.sendMessage(msg);
            })
            .then(response => {
                // Handle response from main.js
                if (response && response.status) {
                    console.log(`Iframe data processing status: ${response.status}`);
                }
            })
            .catch(error => {
                console.log('Extension communication error:', error);
            });
    }

    // Ensure we only collect data once the page is fully loaded
    if (document.readyState === 'complete') {
        setTimeout(collectAndSendData, 500); // Small delay even if complete
    } else {
        // Wait for the page to fully load before collecting data
        window.addEventListener('load', function () {
            // Add a delay to ensure dynamic content is also loaded
            setTimeout(collectAndSendData, 1500);
        });
    }
}