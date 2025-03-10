
console.log('Background script loaded');
chrome.runtime.onMessage.addListener(function (message, sender, sendResponse) {
    // Handle ping messages specifically
    if (message.ping) {
        console.log('Received ping from content script');
        sendResponse({ pong: true });
        return true;
    }

    // Log all messages coming from any frame
    if (message.fromIframe) {
        console.log('Background received iframe content:', message);
    }

    // Forward the message to the main frame if requested
    if (message.sendBack && sender.tab) {
        console.log('Forwarding iframe message to main frame');
        chrome.tabs.sendMessage(sender.tab.id, message);
    }

    return true; // Keep the message channel open for async responses
});