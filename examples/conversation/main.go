package main

import (
	"fmt"
	"log"
	"os"

	"github.com/eisenzopf/agentic-text/pkg/easy"
	"github.com/eisenzopf/agentic-text/pkg/llm"
)

// Sample banking conversations
var sampleConversations = []string{
	`Customer: I need to check my account balance, please.
Agent: I'd be happy to help with that. Could you please verify your account information?
Customer: Sure, my account number is 12345678.
Agent: Thank you. I can see your current balance is $1,542.67. Is there anything else you'd like to know?
Customer: No, that's all I needed. Thank you for your help.
Agent: You're welcome. Have a great day!`,

	`Customer: I'm calling because there's a charge on my account I don't recognize.
Agent: I understand your concern. Let me look into that for you. Can you tell me which charge you're referring to?
Customer: There's a charge for $89.99 from "TechSubscribe" that I never authorized.
Agent: I see that charge. I'll open a dispute right away and issue a temporary credit while we investigate.
Customer: That's great. How long will the investigation take?
Agent: It typically takes 7-10 business days. We'll notify you once it's complete.
Customer: OK, thank you for taking care of this so quickly.`,

	`Customer: I've been trying to make a transfer through your mobile app and it keeps failing.
Agent: I apologize for the inconvenience. Let me help resolve this issue. Are you getting any error message?
Customer: Yes, it says "Transfer failed - insufficient funds" but I know I have enough money in my account.
Agent: I understand that's frustrating. Let me check your account... I see the issue. There's a 24-hour hold on recent deposits. The funds should be available tomorrow.
Customer: That's ridiculous! I need to make this payment today!
Agent: I understand your frustration. As an alternative, would you like me to process this transfer for you over the phone? We can override the hold in this case.
Customer: Yes, please do that. I need to transfer $500 to my savings account immediately.`,
}

func main() {
	// Create a debug configuration to see detailed results
	config := &easy.Config{
		Provider:    llm.Google,
		Model:       "gemini-2.0-flash",
		MaxTokens:   1024,
		Temperature: 0.2,
		Debug:       false,
	}

	fmt.Println("Banking Conversation Analyzer")
	fmt.Println("============================")
	fmt.Println()

	// Check if a custom conversation was provided via command line
	args := os.Args[1:]
	conversations := sampleConversations
	if len(args) > 0 {
		// Use the provided conversation instead of samples
		conversations = []string{args[0]}
		fmt.Println("Analyzing custom conversation...")
	} else {
		fmt.Printf("Analyzing %d sample conversations...\n", len(conversations))
	}

	fmt.Println()

	// Process each conversation for both sentiment and intent
	for i, conversation := range conversations {
		fmt.Printf("Conversation #%d:\n", i+1)
		fmt.Printf("----------------%s\n", "-" /* padding to match number width */)
		fmt.Printf("Excerpt: %s...\n\n", truncateString(conversation, 80))

		// Create wrapper for sentiment analysis
		sentimentWrapper, err := easy.NewWithConfig("sentiment", config)
		if err != nil {
			log.Fatalf("Failed to create sentiment processor: %v", err)
		}

		// Process sentiment
		sentimentResult, err := sentimentWrapper.Process(conversation)
		if err != nil {
			log.Fatalf("Sentiment analysis failed: %v", err)
		}

		// Create wrapper for intent analysis
		intentWrapper, err := easy.NewWithConfig("intent", config)
		if err != nil {
			log.Fatalf("Failed to create intent processor: %v", err)
		}

		// Process intent
		intentResult, err := intentWrapper.Process(conversation)
		if err != nil {
			log.Fatalf("Intent analysis failed: %v", err)
		}

		// Display the results
		fmt.Println("SENTIMENT ANALYSIS:")
		fmt.Printf("  - Overall Sentiment: %s\n", sentimentResult["sentiment"])
		fmt.Printf("  - Score: %.2f\n", sentimentResult["score"])
		fmt.Printf("  - Confidence: %.2f\n", sentimentResult["confidence"])

		if keywords, ok := sentimentResult["keywords"].([]interface{}); ok && len(keywords) > 0 {
			fmt.Println("  - Key Sentiment Words:")
			for _, keyword := range keywords {
				if k, ok := keyword.(string); ok {
					fmt.Printf("    * %s\n", k)
				}
			}
		}

		fmt.Println()
		fmt.Println("INTENT ANALYSIS:")
		fmt.Printf("  - Primary Intent: %s\n", intentResult["label_name"])
		fmt.Printf("  - Machine Label: %s\n", intentResult["label"])
		fmt.Printf("  - Description: %s\n", intentResult["description"])
		fmt.Println()

		// Add a separator between conversations
		if i < len(conversations)-1 {
			fmt.Println("============================")
			fmt.Println()
		}
	}
}

// Helper function to truncate a string to a maximum length
func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen]
}
