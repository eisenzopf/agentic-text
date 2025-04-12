package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/eisenzopf/agentic-text/pkg/easy"
)

func main() {
	// Step 1: Define the questions to get required attributes
	questions := `
What are the top reasons customer's cancel?
How often do agents try to save?
What save offers work the best?
`

	// Step 2: Use required_attributes processor to identify needed attributes
	requiredAttrsResult, err := easy.ProcessText(questions, "required_attributes")
	if err != nil {
		log.Fatalf("Failed to process required attributes: %v", err)
	}

	// Pretty print the required attributes
	attributesJSON, _ := json.MarshalIndent(requiredAttrsResult, "", "  ")
	fmt.Printf("Required Attributes:\n%s\n\n", string(attributesJSON))

	// Step 3: Define a mock customer service conversation
	conversation := `
Customer: Hi, I'd like to cancel my subscription.
Agent: I'm sorry to hear that. May I ask why you're considering cancellation?
Customer: It's just too expensive for what I'm getting.
Agent: I understand your concern about the price. We do have a promotional offer that could reduce your monthly fee by 30% for the next 6 months.
Customer: That's interesting, but I'm still not sure.
Agent: I can also add our premium package at no extra cost, which includes additional features that might provide more value for you.
Customer: That sounds better. How would that work?
Agent: I'll apply both discounts to your account right now, so your next bill will reflect the 30% reduction, and you'll immediately have access to the premium features.
Customer: OK, I'll give it another try with those changes.
Agent: Excellent! I've applied those changes to your account. Is there anything else I can help you with today?
`

	// Step 4: Combine the required attributes and conversation for extraction
	combinedInput := fmt.Sprintf("%s\n\n%s", string(attributesJSON), conversation)

	// Step 5: Use get_attributes processor to extract attribute values
	attributesResult, err := easy.ProcessText(combinedInput, "get_attributes")
	if err != nil {
		log.Fatalf("Failed to extract attributes: %v", err)
	}

	// Step 6: Display the extracted attributes
	extractedJSON, _ := json.MarshalIndent(attributesResult, "", "  ")
	fmt.Printf("Extracted Attributes:\n%s\n", string(extractedJSON))
}
