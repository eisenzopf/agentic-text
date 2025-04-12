package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/eisenzopf/agentic-text/pkg/data"
	"github.com/eisenzopf/agentic-text/pkg/llm"
	"github.com/eisenzopf/agentic-text/pkg/processor"
)

// Server holds the API server configuration
type Server struct {
	provider llm.Provider
}

// ProcessRequest represents a text processing request
type ProcessRequest struct {
	Text      string `json:"text"`
	Processor string `json:"processor"`
}

// ProcessResponse represents a text processing response
type ProcessResponse struct {
	Original string      `json:"original"`
	Result   interface{} `json:"result"`
	Success  bool        `json:"success"`
	Error    string      `json:"error,omitempty"`
}

// NewServer creates a new API server
func NewServer(provider llm.Provider) *Server {
	return &Server{
		provider: provider,
	}
}

// HandleProcess handles text processing requests
func (s *Server) HandleProcess(w http.ResponseWriter, r *http.Request) {
	// Only accept POST requests
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse the request body
	var req ProcessRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Get the processor
	proc, err := processor.Create(req.Processor, s.provider, processor.Options{})
	if err != nil {
		respondWithError(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Create a ProcessItem from the text
	item := data.NewTextProcessItem("api-request", req.Text, nil)

	// Process the text
	result, err := proc.Process(r.Context(), item)
	if err != nil {
		respondWithError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Extract processing results
	var processorResult interface{}
	if result.ProcessingInfo != nil {
		// Get the processing info from the processor
		for _, info := range result.ProcessingInfo {
			processorResult = info
			break
		}
	}

	// Send the response
	response := ProcessResponse{
		Original: req.Text,
		Result:   processorResult,
		Success:  true,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// HandleListProcessors lists available processors
func (s *Server) HandleListProcessors(w http.ResponseWriter, r *http.Request) {
	// Only accept GET requests
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get the list of processors
	processors := processor.ListProcessors()

	// Send the response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"processors": processors,
		"count":      len(processors),
	})
}

// respondWithError sends an error response
func respondWithError(w http.ResponseWriter, message string, status int) {
	response := ProcessResponse{
		Success: false,
		Error:   message,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(response)
}

func main() {
	// Get API key from environment
	apiKey := os.Getenv("LLM_API_KEY")
	if apiKey == "" {
		apiKey = "your-api-key" // For demo purposes only
		fmt.Println("Warning: Using demo API key. Set LLM_API_KEY environment variable for production.")
	}

	// Initialize LLM provider
	config := llm.Config{
		APIKey:      apiKey,
		Model:       "gemini-pro",
		MaxTokens:   1024,
		Temperature: 0.2,
	}

	provider, err := llm.NewProvider(llm.Google, config)
	if err != nil {
		log.Fatalf("Failed to initialize provider: %v", err)
	}

	// Create and start the server
	server := NewServer(provider)

	// Register routes
	http.HandleFunc("/api/process", server.HandleProcess)
	http.HandleFunc("/api/processors", server.HandleListProcessors)

	// Start the server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Printf("Server starting on port %s...\n", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}

/*
Example curl commands:

# List processors
curl -X GET http://localhost:8080/api/processors

# Process text
curl -X POST http://localhost:8080/api/process \
  -H "Content-Type: application/json" \
  -d '{"text": "I really enjoyed this product!", "processor": "sentiment"}'
*/
