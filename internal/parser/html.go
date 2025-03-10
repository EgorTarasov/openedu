package parser

import (
	"openedu/internal/models"
	"strings"

	"golang.org/x/net/html"
)

// ParseContent parses HTML content and extracts problems
func ParseContent(htmlContent string) []models.Problem {
	doc, err := html.Parse(strings.NewReader(htmlContent))
	if err != nil {
		return nil
	}

	var problems []models.Problem
	findProblemElements(doc, &problems)

	return problems
}

// findProblemElements recursively searches for elements with id starting with "problem_"
func findProblemElements(n *html.Node, problems *[]models.Problem) {
	if n.Type == html.ElementNode {
		for _, attr := range n.Attr {
			if attr.Key == "id" && strings.HasPrefix(attr.Val, "problem_") {
				extractedProblems := extractProblemData(n, attr.Val)
				*problems = append(*problems, extractedProblems...)
				break
			}
		}
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		findProblemElements(c, problems)
	}
}

// extractProblemData extracts problem data from an HTML element
func extractProblemData(n *html.Node, id string) []models.Problem {
	var problems []models.Problem

	// Extract problem title (question number)
	titleNode := findNodeByClass(n, "problem-header")
	mainTitle := ""
	if titleNode != nil {
		mainTitle = extractTextContent(titleNode)
	}

	// Extract course info if available
	dataProblemID := ""
	for _, attr := range n.Attr {
		if attr.Key == "data-problem-id" {
			dataProblemID = attr.Val
			break
		}
	}

	course := ""
	if dataProblemID != "" && strings.Contains(dataProblemID, "course-v1:") {
		course = ParseURL(dataProblemID)
	}

	// Find all question wrappers
	wrapperNodes := findNodesByClass(n, "wrapper-problem-response")

	if len(wrapperNodes) == 0 {
		// No question wrappers found, treat as a single problem
		problem := models.Problem{
			ID:     id,
			Title:  mainTitle,
			Course: course,
		}

		// Extract question text (the legend in the fieldset)
		legendNodes := findNodesByTagName(n, "legend")
		if len(legendNodes) > 0 {
			for _, legendNode := range legendNodes {
				if hasClass(legendNode, "response-fieldset-legend") || hasClass(legendNode, "field-group-hd") {
					problem.Question = extractTextContent(legendNode)
					break
				}
			}
		}

		// Try to find question text in p tags if no legend was found
		if problem.Question == "" {
			paragraphs := findNodesByTagName(n, "p")
			if len(paragraphs) > 0 {
				problem.Question = extractTextContent(paragraphs[0])
			}
		}

		// Extract choices and identify correct answers
		var choices []models.Choice
		fieldNodes := findNodesByClass(n, "field")

		for i, fieldNode := range fieldNodes {
			labelNode := findNodeByTagName(fieldNode, "label")
			if labelNode != nil {
				choiceText := extractTextContent(labelNode)
				isCorrect := hasClass(labelNode, "choicegroup_correct")

				choice := models.Choice{
					ID:        i,
					Text:      choiceText,
					IsCorrect: isCorrect,
				}

				choices = append(choices, choice)

				if isCorrect {
					problem.Answer = append(problem.Answer, choiceText)
				}
			}
		}

		problem.Choices = choices
		problems = append(problems, problem)
	} else {
		// Multiple questions in one problem
		// Find all paragraphs in this node
		allParagraphs := findNodesByTagName(n, "p")

		// Map each wrapper to its preceding paragraph
		wrapperToParagraph := make(map[*html.Node]*html.Node)

		// For each wrapper, find the closest preceding paragraph
		for _, wrapper := range wrapperNodes {
			// Find parent div that contains both wrappers and paragraphs
			parent := findFirstCommonParent(wrapper, allParagraphs)
			if parent != nil {
				// Create a map of nodes to their order in the DOM
				nodeOrder := make(map[*html.Node]int)
				orderIndex := 0

				// Walk through parent's children and record their order
				walkDOM(parent, func(node *html.Node) {
					nodeOrder[node] = orderIndex
					orderIndex++
				})

				// Find the closest paragraph that comes before this wrapper
				closestParagraph := findClosestPrecedingParagraph(wrapper, allParagraphs, nodeOrder)
				if closestParagraph != nil {
					wrapperToParagraph[wrapper] = closestParagraph
				}
			}
		}

		// Create a problem for each wrapper
		for i, wrapperNode := range wrapperNodes {
			problem := models.Problem{
				ID:     id + "_" + string(rune('a'+i)),
				Title:  mainTitle,
				Course: course,
			}

			// Set question text from the associated paragraph
			if paragraph, ok := wrapperToParagraph[wrapperNode]; ok && paragraph != nil {
				problem.Question = extractTextContent(paragraph)
			} else if ariaLabel := getAttrValue(wrapperNode, "aria-label"); ariaLabel != "" {
				// Fall back to aria-label if no paragraph found
				problem.Question = ariaLabel
			}

			// Extract choices and identify correct answers
			var choices []models.Choice
			fieldNodes := findNodesByClass(wrapperNode, "field")

			for j, fieldNode := range fieldNodes {
				labelNode := findNodeByTagName(fieldNode, "label")
				if labelNode != nil {
					choiceText := extractTextContent(labelNode)
					isCorrect := hasClass(labelNode, "choicegroup_correct")

					choice := models.Choice{
						ID:        j,
						Text:      choiceText,
						IsCorrect: isCorrect,
					}

					choices = append(choices, choice)

					if isCorrect {
						problem.Answer = append(problem.Answer, choiceText)
					}
				}
			}

			problem.Choices = choices
			problems = append(problems, problem)
		}
	}

	return problems
}

// walkDOM traverses the DOM tree and calls the provided function for each node
func walkDOM(n *html.Node, fn func(*html.Node)) {
	if n == nil {
		return
	}

	fn(n)

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		walkDOM(c, fn)
	}
}

// findFirstCommonParent finds the closest ancestor that contains both the wrapper and at least one paragraph
func findFirstCommonParent(wrapper *html.Node, paragraphs []*html.Node) *html.Node {
	// Start with the wrapper's parent
	for parent := wrapper.Parent; parent != nil; parent = parent.Parent {
		// Check if this parent contains any paragraphs
		for _, p := range paragraphs {
			if isAncestor(parent, p) {
				return parent
			}
		}
	}
	return nil
}

// isAncestor checks if node is an ancestor of potentialChild
func isAncestor(node *html.Node, potentialChild *html.Node) bool {
	for parent := potentialChild.Parent; parent != nil; parent = parent.Parent {
		if parent == node {
			return true
		}
	}
	return false
}

// findClosestPrecedingParagraph finds the paragraph that precedes the wrapper in document order
func findClosestPrecedingParagraph(wrapper *html.Node, paragraphs []*html.Node, nodeOrder map[*html.Node]int) *html.Node {
	wrapperOrder := nodeOrder[wrapper]
	closestParagraph := (*html.Node)(nil)
	closestDistance := -1

	for _, paragraph := range paragraphs {
		paragraphOrder := nodeOrder[paragraph]

		// Only consider paragraphs that come before the wrapper
		if paragraphOrder < wrapperOrder {
			distance := wrapperOrder - paragraphOrder

			// Update if this is the first paragraph or the closest one so far
			if closestParagraph == nil || distance < closestDistance {
				closestParagraph = paragraph
				closestDistance = distance
			}
		}
	}

	return closestParagraph
}

// getAttrValue gets attribute value for a given node and key
func getAttrValue(n *html.Node, key string) string {
	for _, attr := range n.Attr {
		if attr.Key == key {
			return attr.Val
		}
	}
	return ""
}

// findNodeByTagName finds the first node with the specified tag name
func findNodeByTagName(n *html.Node, tagName string) *html.Node {
	if n.Type == html.ElementNode && n.Data == tagName {
		return n
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if found := findNodeByTagName(c, tagName); found != nil {
			return found
		}
	}
	return nil
}

// findNodesByTagName finds all nodes with the specified tag name
func findNodesByTagName(n *html.Node, tagName string) []*html.Node {
	var nodes []*html.Node

	if n.Type == html.ElementNode && n.Data == tagName {
		nodes = append(nodes, n)
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		found := findNodesByTagName(c, tagName)
		nodes = append(nodes, found...)
	}

	return nodes
}

// findNodeByClass finds the first node with the specified class
func findNodeByClass(n *html.Node, className string) *html.Node {
	if hasClass(n, className) {
		return n
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if found := findNodeByClass(c, className); found != nil {
			return found
		}
	}
	return nil
}

// findNodesByClass finds all nodes with the specified class
func findNodesByClass(n *html.Node, className string) []*html.Node {
	var nodes []*html.Node

	if hasClass(n, className) {
		nodes = append(nodes, n)
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		foundNodes := findNodesByClass(c, className)
		nodes = append(nodes, foundNodes...)
	}

	return nodes
}

// hasClass checks if a node has a specific class
func hasClass(n *html.Node, className string) bool {
	if n.Type != html.ElementNode {
		return false
	}

	for _, attr := range n.Attr {
		if attr.Key == "class" {
			classes := strings.Fields(attr.Val)
			for _, class := range classes {
				if class == className {
					return true
				}
			}
		}
	}
	return false
}

// extractTextContent extracts all text from a node and its children
func extractTextContent(n *html.Node) string {
	if n == nil {
		return ""
	}

	if n.Type == html.TextNode {
		return n.Data
	}

	var result string
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		result += extractTextContent(c)
	}

	return strings.TrimSpace(result)
}

func ParseURL(url string) (course string) {
	if strings.Contains(url, "course-v1:") {
		parts := strings.Split(url, "course-v1:")
		if len(parts) > 1 {
			courseParts := strings.Split(parts[1], "+")
			if len(courseParts) > 0 {
				return courseParts[0]
			}
		}
	}
	return ""
}
