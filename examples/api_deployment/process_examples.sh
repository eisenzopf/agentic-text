#!/bin/bash

# Define the base URL for the API
BASE_URL="http://localhost:8080/api/process" # Change port if needed

# --- Example 1: Sentiment Analysis ---
echo "--- Running Sentiment Analysis ---"
curl -X POST "$BASE_URL" \
  -H "Content-Type: application/json" \
  -d '{"text": "I am absolutely thrilled with this service! It exceeded all expectations.", "processor": "sentiment"}'
echo # Add a newline for readability
echo

# --- Example 2: Intent Detection ---
echo "--- Running Intent Detection ---"
curl -X POST "$BASE_URL" \
  -H "Content-Type: application/json" \
  -d '{"text": "Hi, I need to update my billing address and also ask about the return policy.", "processor": "intent"}'
echo
echo

# --- Example 3: Keyword Extraction ---
echo "--- Running Keyword Extraction ---"
curl -X POST "$BASE_URL" \
  -H "Content-Type: application/json" \
  -d '{"text": "Agentic Text is a Go library designed for building applications that leverage Large Language Models for advanced text processing tasks.", "processor": "keyword_extraction"}'
echo
echo

# --- Example 4: Speech Act Classification ---
echo "--- Running Speech Act Classification ---"
curl -X POST "$BASE_URL" \
  -H "Content-Type: application/json" \
  -d '{"text": "Hello there. Could you please tell me the status of order #12345? Thanks!", "processor": "speech_act"}'
echo
echo

# --- Example 5: Get Attributes ---
echo "--- Running Get Attributes ---"
curl -X POST "$BASE_URL" \
  -H "Content-Type: application/json" \
  -d '{"text": "The customer is Jane Smith, phone number 555-1234, reporting a broken screen on model X.", "processor": "get_attributes"}'
echo
echo

# --- Example 6: Required Attributes ---
echo "--- Running Required Attributes ---"
curl -X POST "$BASE_URL" \
  -H "Content-Type: application/json" \
  -d '{"text": "What data fields are necessary to identify the most valuable customers based on purchase history and engagement?", "processor": "required_attributes"}'
echo
echo

echo "--- All examples finished ---"

