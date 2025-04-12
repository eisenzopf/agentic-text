# Banking Conversation Analyzer

This example demonstrates how to use the agentic-text library to analyze banking customer service conversations for both sentiment and intent.

## Overview

The example:
1. Analyzes sample banking conversations (or custom input)
2. Processes each conversation through sentiment and intent analyzers
3. Displays the results in a human-readable format

## Features

- **Sentiment Analysis**: Identifies whether the conversation has a positive, negative, or neutral tone
- **Intent Analysis**: Determines the customer's primary purpose for contacting the bank
- **Multiple Sample Conversations**: Includes three typical banking customer service scenarios
- **Custom Input Support**: Allows analyzing your own conversation text

## Usage

From the project root, run:

```bash
go run examples/banking_sentiment/main.go
```

This will analyze the sample conversations included in the code.

To analyze your own conversation text:

```bash
go run examples/banking_sentiment/main.go "Your conversation text here"
```

## Sample Output

The output will look similar to this:

```
Banking Conversation Analyzer
============================

Analyzing 3 sample conversations...

Conversation #1:
-----------------
Excerpt: Customer: I need to check my account balance, please.
Agent: I'd be happy to...

SENTIMENT ANALYSIS:
  - Overall Sentiment: positive
  - Score: 0.65
  - Confidence: 0.85
  - Key Sentiment Words:
    * happy
    * help
    * thank you

INTENT ANALYSIS:
  - Primary Intent: Balance Inquiry
  - Machine Label: balance_inquiry
  - Description: The customer contacted the bank to check their account balance.

============================
...
```

## Customization

You can modify the code to:
- Change the LLM provider (Google, OpenAI, etc.)
- Adjust model parameters
- Add additional analysis processors
- Customize the output format

See the `easy.Config` struct in the code for configuration options. 