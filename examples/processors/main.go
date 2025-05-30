package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/eisenzopf/agentic-text/pkg/easy"
	_ "github.com/eisenzopf/agentic-text/pkg/processor/builtin"
)

// ConversationData represents a customer service conversation
type ConversationData struct {
	ID       string `json:"id"`
	Text     string `json:"text"`
	Category string `json:"category"`
	Summary  string `json:"summary"`
}

// TestResult holds the results of processor testing
type TestResult struct {
	ProcessorName string      `json:"processor_name"`
	Input         string      `json:"input"`
	Output        interface{} `json:"output"`
	Duration      string      `json:"duration"`
	Error         string      `json:"error,omitempty"`
}

var sampleConversations = []ConversationData{
	{
		ID:       "conv_001",
		Category: "fee_dispute",
		Summary:  "Customer disputing late payment fee",
		Text: `Agent: Hello, thank you for calling CustomerCare. My name is Sarah. How can I help you today?

Customer: Hi Sarah, I'm calling about this $35 late payment fee on my bill. I don't think I should have been charged this because I set up autopay last month.

Agent: I understand your concern about the late payment fee. Let me pull up your account to take a look. I can see that you did set up autopay on March 15th, but your payment due date was March 10th, so unfortunately the autopay didn't process in time for that billing cycle.

Customer: But I set it up as soon as I got the bill! This doesn't seem fair. I've been a customer for 8 years and never had any issues before.

Agent: I absolutely understand your frustration, and I really appreciate your loyalty as an 8-year customer. You're right that you acted quickly to set up the autopay. Let me see what I can do for you. Since this was due to the timing of when you set up autopay and you've been such a great customer, I'm going to go ahead and credit that $35 fee back to your account.

Customer: Really? Thank you so much! I really appreciate that.

Agent: You're very welcome! The credit will appear on your next billing statement. Is there anything else I can help you with today?

Customer: No, that takes care of everything. Thank you again for your help!

Agent: You're welcome! Thank you for being a valued customer. Have a great day!`,
	},
	{
		ID:       "conv_002",
		Category: "technical_support",
		Summary:  "Internet connectivity issues troubleshooting",
		Text: `Agent: Tech Support, this is Mike. How can I assist you today?

Customer: Hi Mike, my internet has been really slow for the past few days. I'm only getting like 10 Mbps when I'm supposed to get 100.

Agent: I'm sorry to hear you're experiencing slow speeds. Let me help you troubleshoot this. First, are you connected via WiFi or ethernet cable?

Customer: I'm using WiFi right now.

Agent: Okay, let's start by having you run a speed test while connected directly via ethernet cable if possible. This will help us determine if it's a WiFi issue or something with the connection itself.

Customer: Alright, let me plug in the cable... okay, running the test now... wow, I'm getting 95 Mbps now!

Agent: Great! That tells us your internet service is working properly, but there's an issue with your WiFi. Let's try a few things. First, can you tell me how old your router is?

Customer: I think it's about 5 years old. It's the one you guys provided when I first signed up.

Agent: That could definitely be the issue. Those older routers don't handle higher speeds as well. Let me check what newer models we have available. I can see we have a new high-speed router that would be perfect for your plan. I can send one out to you at no charge since you've been a customer for over 3 years.

Customer: That would be amazing! How long would that take?

Agent: I can have it shipped out today and you should receive it by tomorrow. I'll also include easy setup instructions.

Customer: Perfect, thank you so much for your help!`,
	},
	{
		ID:       "conv_003",
		Category: "billing_inquiry",
		Summary:  "Customer confused about sudden bill increase",
		Text: `Agent: Billing department, this is Jessica. How can I help you?

Customer: Hi, I just got my bill and it's $40 higher than usual. I don't understand what happened.

Agent: I'd be happy to help you understand the changes to your bill. Let me pull up your account. I can see there are a couple of changes this month. First, your promotional rate for internet expired, which added $25 to your bill. Second, you upgraded your cable package last month which added $15.

Customer: Oh no, I didn't realize the promotion was ending. I definitely don't remember upgrading anything though.

Agent: Let me check the details on that upgrade... I can see it was processed on March 3rd when you called in. You added the Sports package. Do you recall calling about getting some specific sports channels?

Customer: Oh wait, yes! My husband called because he wanted to watch the basketball playoffs. I didn't realize it would be that much extra.

Agent: I understand. The Sports package is $15 per month. The good news is that you can cancel it anytime if you don't want to keep it. As for the promotional rate, I can see if there are any current promotions you might qualify for.

Customer: Yes, please look into that. And I think we probably want to cancel the sports package since the playoffs are almost over.

Agent: Absolutely. I've cancelled the sports package effective at the end of this billing cycle, so you won't be charged for it next month. And I found a current promotion that can give you $20 off your internet for the next 12 months. Would you like me to apply that?

Customer: Yes, that would be great! So my bill will actually be lower than before?

Agent: Exactly! With the sports package removed and the new promotion applied, your bill will be $5 less than your original amount.

Customer: Perfect, thank you so much for explaining everything and getting this sorted out!`,
	},
	{
		ID:       "conv_004",
		Category: "service_cancellation",
		Summary:  "Customer wanting to cancel service due to relocation",
		Text: `Agent: Customer Service, this is David. How may I help you today?

Customer: Hi David, I need to cancel my service because I'm moving to a different state next month.

Agent: I'm sorry to hear you'll be leaving us. I'd be happy to help you with that. May I ask which state you're moving to? We might be able to transfer your service instead of cancelling.

Customer: I'm moving to Oregon. Do you provide service there?

Agent: Unfortunately, we don't currently provide service in Oregon, so cancellation would be the right option. Let me help make this process as smooth as possible for you. When is your move date?

Customer: April 15th. I'd like to have the service disconnected on April 14th.

Agent: Perfect. I can schedule your disconnection for April 14th. Now, I do need to let you know that since you're under contract, there will be an early termination fee of $180. However, given that this is due to moving outside our service area, I can waive that fee for you.

Customer: Oh wow, I was really worried about that fee. Thank you so much!

Agent: You're very welcome! I also want to make sure you know about our Win Back program. If you ever move back to an area where we provide service, just give us a call and we can often get you set up with a great returning customer discount.

Customer: That's good to know. One more question - I'm renting this modem from you guys. How do I return it?

Agent: Great question. I'll email you a prepaid return label. Just pack the modem in any box and drop it off at any UPS location within 30 days. If you don't have a box, any UPS store can provide one at no charge.

Customer: This has been so much easier than I expected. Thank you for all your help!

Agent: You're very welcome! Is there anything else I can help you with today?

Customer: No, that covers everything. Thanks again!`,
	},
	{
		ID:       "conv_005",
		Category: "product_inquiry",
		Summary:  "Customer interested in upgrading internet speed",
		Text: `Agent: Sales department, this is Amanda. How can I help you today?

Customer: Hi Amanda, I work from home now and my current internet speed isn't cutting it. I keep getting dropped from video calls. What options do you have for faster speeds?

Agent: I'd be happy to help you find a better plan for working from home. Let me first check what plan you currently have... I see you're on our Basic plan with 25 Mbps download speed. For reliable video conferencing, I'd recommend at least 50 Mbps, but our Premium plan with 200 Mbps would give you plenty of headroom for multiple video calls and other usage.

Customer: What would be the price difference?

Agent: Your current plan is $39.99 per month. The Premium plan is $69.99, but I have a current promotion that can bring that down to $54.99 for the first year. That's only $15 more than what you're paying now.

Customer: That sounds reasonable. How quickly could I get the upgrade?

Agent: Since you already have our service, this is just a plan change on our end. The speed increase would take effect within 2-4 hours of processing the upgrade. You might need to restart your modem, but that's it.

Customer: Perfect! And if I'm not happy with it, can I downgrade?

Agent: Absolutely. There's no contract required for the upgrade, so you can change your plan anytime. Though I'm confident you'll love the improved performance for your work calls.

Customer: Alright, let's do it. Should I restart my modem now or wait?

Agent: I'm processing the upgrade right now... okay, it's complete. Give it about an hour and then restart your modem by unplugging it for 30 seconds. You should see the faster speeds after that.

Customer: Excellent! Thank you for making this so easy.

Agent: You're welcome! Enjoy your improved internet speed, and don't hesitate to call if you need anything else.`,
	},
}

func main() {
	var (
		interactive = flag.Bool("interactive", false, "Run in interactive mode")
		processor   = flag.String("processor", "", "Specific processor to test")
		output      = flag.String("output", "", "Save results to JSON file")
		verbose     = flag.Bool("verbose", false, "Enable verbose output")
		testAll     = flag.Bool("test-all", false, "Test all processors with all sample conversations")
	)
	flag.Parse()

	// Check for API key
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		apiKey = os.Getenv("GEMINI_API_KEY")
		if apiKey == "" {
			log.Fatal("OPENAI_API_KEY or GEMINI_API_KEY environment variable is required")
		}
	}

	// Get available processors
	availableProcessors := easy.ListAvailableProcessors()

	fmt.Println("🔬 Processor Testing Application")
	fmt.Println("=================================")
	fmt.Printf("Available processors: %s\n", strings.Join(availableProcessors, ", "))
	fmt.Printf("Sample conversations: %d\n\n", len(sampleConversations))

	var results []TestResult

	if *testAll {
		results = testAllProcessors(*verbose)
	} else if *interactive {
		results = runInteractiveMode()
	} else if *processor != "" {
		results = testSpecificProcessor(*processor, *verbose)
	} else {
		flag.Usage()
		fmt.Println("\nExamples:")
		fmt.Println("  go run main.go -interactive")
		fmt.Println("  go run main.go -processor sentiment")
		fmt.Println("  go run main.go -test-all -output results.json")
		return
	}

	// Save results if requested
	if *output != "" {
		saveResults(results, *output)
	}
}

func testAllProcessors(verbose bool) []TestResult {
	fmt.Println("🧪 Testing all processors with all sample conversations...")
	fmt.Println()

	availableProcessors := easy.ListAvailableProcessors()
	var allResults []TestResult

	for i, procName := range availableProcessors {
		fmt.Printf("[%d/%d] Testing processor: %s\n", i+1, len(availableProcessors), procName)
		fmt.Println(strings.Repeat("-", 50))

		for j, conv := range sampleConversations {
			fmt.Printf("  Conversation %d (%s): ", j+1, conv.Category)

			result := testProcessorWithConversation(procName, conv, verbose)
			allResults = append(allResults, result)

			if result.Error != "" {
				fmt.Printf("❌ Error\n")
				if verbose {
					fmt.Printf("    Error: %s\n", result.Error)
				}
			} else {
				fmt.Printf("✅ Success (%s)\n", result.Duration)
				if verbose {
					printProcessorResult(result)
				}
			}
		}
		fmt.Println()
	}

	// Print summary
	successCount := 0
	for _, result := range allResults {
		if result.Error == "" {
			successCount++
		}
	}

	fmt.Printf("📊 Test Summary: %d/%d successful (%.1f%%)\n",
		successCount, len(allResults), float64(successCount)/float64(len(allResults))*100)

	return allResults
}

func testSpecificProcessor(procName string, verbose bool) []TestResult {
	availableProcessors := easy.ListAvailableProcessors()
	found := false
	for _, p := range availableProcessors {
		if p == procName {
			found = true
			break
		}
	}

	if !found {
		fmt.Printf("❌ Unknown processor: %s\n", procName)
		fmt.Printf("Available processors: %s\n", strings.Join(availableProcessors, ", "))
		return nil
	}

	fmt.Printf("🎯 Testing processor: %s\n", procName)
	fmt.Println(strings.Repeat("-", 50))

	var results []TestResult

	for i, conv := range sampleConversations {
		fmt.Printf("Conversation %d (%s): ", i+1, conv.Category)

		result := testProcessorWithConversation(procName, conv, verbose)
		results = append(results, result)

		if result.Error != "" {
			fmt.Printf("❌ Error\n")
			if verbose {
				fmt.Printf("  Error: %s\n", result.Error)
			}
		} else {
			fmt.Printf("✅ Success (%s)\n", result.Duration)
			if verbose {
				fmt.Println()
				printProcessorResult(result)
			}
		}
		fmt.Println()
	}

	return results
}

func runInteractiveMode() []TestResult {
	fmt.Println("🎮 Interactive Mode")
	fmt.Println("Type 'help' for commands, 'quit' to exit")
	fmt.Println()

	scanner := bufio.NewScanner(os.Stdin)
	var results []TestResult
	availableProcessors := easy.ListAvailableProcessors()

	for {
		fmt.Print("processor-test> ")
		if !scanner.Scan() {
			break
		}

		input := strings.TrimSpace(scanner.Text())
		if input == "" {
			continue
		}

		parts := strings.Fields(input)
		command := parts[0]

		switch command {
		case "help":
			printHelp()

		case "list":
			fmt.Println("Available processors:")
			for i, proc := range availableProcessors {
				fmt.Printf("  %d. %s\n", i+1, proc)
			}

		case "conversations", "convs":
			fmt.Println("Sample conversations:")
			for i, conv := range sampleConversations {
				fmt.Printf("  %d. %s (%s)\n", i+1, conv.Category, conv.ID)
				fmt.Printf("     %s\n", conv.Summary)
			}

		case "test":
			if len(parts) < 2 {
				fmt.Println("Usage: test <processor> [conversation_id]")
				continue
			}

			procName := parts[1]
			if !contains(availableProcessors, procName) {
				fmt.Printf("Unknown processor: %s\n", procName)
				continue
			}

			var targetConvs []ConversationData
			if len(parts) >= 3 {
				convID := parts[2]
				conv := findConversationByID(convID)
				if conv == nil {
					fmt.Printf("Unknown conversation ID: %s\n", convID)
					continue
				}
				targetConvs = []ConversationData{*conv}
			} else {
				targetConvs = sampleConversations
			}

			for _, conv := range targetConvs {
				fmt.Printf("\nTesting %s with conversation %s...\n", procName, conv.ID)
				result := testProcessorWithConversation(procName, conv, true)
				results = append(results, result)

				if result.Error != "" {
					fmt.Printf("❌ Error: %s\n", result.Error)
				} else {
					fmt.Printf("✅ Success (%s)\n", result.Duration)
					printProcessorResult(result)
				}
			}

		case "custom":
			if len(parts) < 2 {
				fmt.Println("Usage: custom <processor>")
				continue
			}

			procName := parts[1]
			if !contains(availableProcessors, procName) {
				fmt.Printf("Unknown processor: %s\n", procName)
				continue
			}

			fmt.Println("Enter your text (end with '###' on a new line):")
			var lines []string
			for scanner.Scan() {
				line := scanner.Text()
				if line == "###" {
					break
				}
				lines = append(lines, line)
			}

			customText := strings.Join(lines, "\n")
			if customText == "" {
				fmt.Println("No text provided")
				continue
			}

			customConv := ConversationData{
				ID:   "custom",
				Text: customText,
			}

			fmt.Printf("\nTesting %s with custom text...\n", procName)
			result := testProcessorWithConversation(procName, customConv, true)
			results = append(results, result)

			if result.Error != "" {
				fmt.Printf("❌ Error: %s\n", result.Error)
			} else {
				fmt.Printf("✅ Success (%s)\n", result.Duration)
				printProcessorResult(result)
			}

		case "quit", "exit":
			fmt.Println("Goodbye!")
			return results

		default:
			fmt.Printf("Unknown command: %s\nType 'help' for available commands\n", command)
		}

		fmt.Println()
	}

	return results
}

func testProcessorWithConversation(procName string, conv ConversationData, verbose bool) TestResult {
	start := time.Now()

	// Prepare input text based on processor type
	inputText := prepareInputForProcessor(procName, conv)

	// Create options for verbose mode
	options := map[string]interface{}{}
	if verbose {
		options["debug"] = true
	}

	// Process using easy library
	result, err := easy.ProcessTextWithOptions(inputText, procName, options)
	duration := time.Since(start)

	testResult := TestResult{
		ProcessorName: procName,
		Input:         conv.ID,
		Duration:      duration.String(),
	}

	if err != nil {
		testResult.Error = err.Error()
	} else {
		testResult.Output = result
	}

	return testResult
}

func prepareInputForProcessor(procName string, conv ConversationData) string {
	switch procName {
	case "required_attributes":
		// For required_attributes, provide research questions
		return fmt.Sprintf(`Research Questions about this conversation:
1. What was the customer's main concern or issue?
2. How did the agent resolve the customer's problem?
3. What was the customer's satisfaction level?
4. What type of service interaction was this?
5. Were there any escalation points in the conversation?

Conversation:
%s`, conv.Text)

	case "data_analyzer":
		// For data_analyzer, provide structured analysis questions
		return fmt.Sprintf(`Conversation Data: %s

Research Questions:
1. What was the primary customer issue in this conversation?
2. What resolution steps were taken by the agent?
3. What was the outcome of the interaction?
4. What patterns can be identified in the customer service approach?
5. How would you rate the overall customer experience?`, conv.Text)

	case "recommendation_engine":
		// For recommendation_engine, provide analysis context
		return fmt.Sprintf(`Customer Service Analysis:

Conversation Category: %s
Summary: %s

Full Conversation:
%s

Please provide recommendations for improving customer service based on this interaction.`,
			conv.Category, conv.Summary, conv.Text)

	case "quality_reviewer":
		// For quality_reviewer, provide content to review
		return fmt.Sprintf(`Please review this customer service conversation for quality:

Conversation Summary: %s
Category: %s

Conversation Content:
%s

Evaluate the agent's performance, customer satisfaction, and overall interaction quality.`,
			conv.Summary, conv.Category, conv.Text)

	case "attribute_matcher":
		// For attribute_matcher, provide required and available attributes
		return fmt.Sprintf(`Required Attributes:
- customer_issue: The main problem or concern raised by the customer
- resolution_method: How the agent addressed the customer's issue  
- customer_satisfaction: The customer's level of satisfaction with the resolution
- interaction_type: The category of customer service interaction

Available Attributes:
- main_concern: Primary issue mentioned by customer
- agent_solution: Steps taken by agent to help
- customer_mood: Customer's emotional state during call
- service_category: Type of service request

Conversation Context:
%s`, conv.Text)

	case "question_generator":
		// For question_generator, provide business context
		return fmt.Sprintf(`Business Context: Customer Service Operations
Conversation Type: %s
Sample Interaction Summary: %s

Generate research questions that would help analyze and improve customer service operations based on conversation data like this.`,
			conv.Category, conv.Summary)

	default:
		// For other processors, use the conversation text directly
		return conv.Text
	}
}

func printProcessorResult(result TestResult) {
	output, err := json.MarshalIndent(result.Output, "    ", "  ")
	if err != nil {
		fmt.Printf("    Error formatting output: %v\n", err)
		return
	}

	fmt.Printf("    Result:\n%s\n", string(output))
}

func printHelp() {
	fmt.Println("Available commands:")
	fmt.Println("  help                     - Show this help message")
	fmt.Println("  list                     - List available processors")
	fmt.Println("  conversations, convs     - List sample conversations")
	fmt.Println("  test <processor> [conv]  - Test processor with conversation(s)")
	fmt.Println("  custom <processor>       - Test processor with custom text")
	fmt.Println("  quit, exit               - Exit interactive mode")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  test sentiment           - Test sentiment with all conversations")
	fmt.Println("  test intent conv_001     - Test intent with specific conversation")
	fmt.Println("  custom data_analyzer     - Test data_analyzer with custom input")
}

func saveResults(results []TestResult, filename string) {
	file, err := os.Create(filename)
	if err != nil {
		fmt.Printf("Error creating output file: %v\n", err)
		return
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")

	if err := encoder.Encode(results); err != nil {
		fmt.Printf("Error encoding results: %v\n", err)
		return
	}

	fmt.Printf("Results saved to %s\n", filename)
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

func findConversationByID(id string) *ConversationData {
	for _, conv := range sampleConversations {
		if conv.ID == id {
			return &conv
		}
	}
	return nil
}
