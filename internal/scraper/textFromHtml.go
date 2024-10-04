package scraper

import (
	"fmt"
	"net/http"
	"strings"

	"golang.org/x/net/html"
)

// Helper function to print the HTML node structure (for debugging)
func printHTMLStructure(node *html.Node, indent string) {
	if node.Type == html.ElementNode {
		fmt.Printf("%s<%s>\n", indent, node.Data)
	} else if node.Type == html.TextNode {
		fmt.Printf("%sText: %q\n", indent, strings.TrimSpace(node.Data))
	}

	for c := node.FirstChild; c != nil; c = c.NextSibling {
		printHTMLStructure(c, indent+"  ")
	}
}

// Function to extract text from HTML nodes
func extractTextFromHTML(node *html.Node) string {
	if node.Type == html.TextNode {
		// Return trimmed text content from text nodes
		return strings.TrimSpace(node.Data)
	}
	if node.Type != html.ElementNode {
		return ""
	}

	// For element nodes, recursively extract text from child nodes
	var text string
	for c := node.FirstChild; c != nil; c = c.NextSibling {
		childText := extractTextFromHTML(c)
		if len(childText) > 0 {
			// Add a space between the concatenated text to preserve words
			text += childText + " "
		}
	}

	// Return the full text for this node
	return strings.TrimSpace(text)
}

// Function to fetch the HTML content from a URL
func FetchHTML(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("status code error: %d %s", resp.StatusCode, resp.Status)
	}

	doc, err := html.Parse(resp.Body)
	if err != nil {
		return "", err
	}

	// Debugging: Print the structure of the HTML document
	printHTMLStructure(doc, "")

	// Extract text from the parsed HTML document
	return extractTextFromHTML(doc), nil
}
