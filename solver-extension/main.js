// Main page script
console.log('Main page script loaded');

const api = "http://localhost:8080"

function extractQuestions(html) {
    // Create a DOM parser to work with the HTML string
    const parser = new DOMParser();
    const doc = parser.parseFromString(html, 'text/html');

    // Find all question paragraphs (they start with a number followed by a dot)
    const questions = [];
    const questionParagraphs = doc.querySelectorAll('p');

    // Process each question paragraph and its following wrapper-problem-response
    questionParagraphs.forEach((questionP) => {
        // Check if this paragraph contains a question (starts with a number)
        const questionText = questionP.textContent.trim();
        if (!/^\d+\./.test(questionText)) return; // Skip if not a numbered question

        // Find the response wrapper that follows this question
        const responseWrapper = questionP.nextElementSibling;
        if (!responseWrapper || !responseWrapper.classList.contains('wrapper-problem-response')) return;

        // Find all choices in this question
        const choices = [];
        const choiceElements = responseWrapper.querySelectorAll('.field');

        choiceElements.forEach((choiceEl) => {
            const input = choiceEl.querySelector('input[type="radio"]');
            const label = choiceEl.querySelector('label');

            if (!input || !label) return;

            const isCorrect = label.classList.contains('choicegroup_correct');
            const isChecked = input.hasAttribute('checked');

            choices.push({
                text: label.textContent.trim(),
                isCorrect: isCorrect,
                isSelected: isChecked
            });
        });

        // Extract the question ID from the input name
        let questionId = '';
        const firstInput = responseWrapper.querySelector('input[type="radio"]');
        if (firstInput) {
            questionId = firstInput.name;
        }

        // Find the correct answer (if marked)
        const correctAnswers = choices
            .filter(choice => choice.isCorrect)
            .map(choice => choice.text);

        questions.push({
            id: questionId,
            question: questionText,
            choices: choices,
            correctAnswers: correctAnswers,
            userSelectedAnswer: choices.find(c => c.isSelected)?.text || ''
        });
    });

    // Get the problem ID if available
    let problemId = '';
    const problemIdInput = doc.querySelector('input[name="problem_id"]');
    if (problemIdInput) {
        problemId = problemIdInput.value;
    }

    return {
        problemId: problemId,
        questions: questions
    };
}

function insertHtmlIntoUnit(html) {
    // Find the element with class 'unit'
    const unitElement = document.querySelector('.notification-tray-divider');
    
    if (unitElement) {
        // Insert the HTML content into the unit element
        unitElement.innerHTML = unitElement.innerHTML +  html;
        console.log('HTML content inserted into unit element');
    } else {
        console.error('Element with class "unit" not found');
    }
}


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
            // console.log('Sending new data from iframe URL:', message.problemDiv);
            const questionsData = extractQuestions(message.problemDiv);
            console.log('Extracted questions:', questionsData);

            // Create an array to hold all promises for fetching answers
            const fetchPromises = [];
            const answers = [];

            // Create a promise for each question
            questionsData.questions.forEach((question) => {
                console.log(question.question);
                const promise = fetch(`${api}/q?q=${question.question}`, {
                    method: "GET",
                    headers: {
                        "Content-Type": "application/json"
                    }
                })
                    .then(res => res.json())
                    .then(data => {
                        if (data.length === 0) {
                            console.log('No data found for question:', question.question);
                            return null;
                        }
                        console.log('Found data:', data);
                        const top = data[0];
                        // Save the answer with its question number for sorting later
                        const questionNumber = parseInt(question.question.match(/^(\d+)\./)?.[1] || "0");
                        answers.push({
                            questionNumber: questionNumber,
                            question: top.Question,
                            answer: top.Answer
                        });
                    })
                    .catch(err => {
                        console.error('Error fetching answer:', err);
                    });

                fetchPromises.push(promise);
            });

            // Wait for all fetches to complete
            Promise.all(fetchPromises)
                .then(() => {
                    // Sort answers by question number
                    answers.sort((a, b) => a.questionNumber - b.questionNumber);

                    // Create HTML for all answers
                    let allAnswersHtml = '<div class="openedu-answers">';
                    answers.forEach(item => {
                        allAnswersHtml += `
                            <div class="openedu-answer">
                                <h6>Question: ${item.question}</h6>
                                <p>Answer: ${item.answer}</p>
                            </div>
                        `;
                    });
                    allAnswersHtml += '</div>';

                    // Insert all answers at once
                    if (answers.length > 0) {
                        insertHtmlIntoUnit(allAnswersHtml);
                    }

                    // Mark this URL as processed
                    localStorage.setItem(storageKey, 'true');
                    sendResponse({ status: 'success' });
                });

            // Return true to indicate we'll send response asynchronously
            return true;
        }
    }

    return true;
});