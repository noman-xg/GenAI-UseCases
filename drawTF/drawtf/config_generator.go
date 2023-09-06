package drawtf

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	openai "github.com/sashabaranov/go-openai"
)

// Constants used in the package.
const (
	PythonInterpreter = "python3"                  // Path to the Python interpreter.
	EmbeddingsScript  = "./scripts/embeddings.py"  // Path to the embeddings Python script.
	GradioLauncher    = "./scripts/chat-portal.py" // Path to the Gradio launcher script.
	FilePath          = "./initialConf.txt"        // Path to a file used for temporary storage.
	TerrFile          = "./main.tf"                // Path to a Terraform configuration file.
	Delimiter         = "####"                     // Delimiter used in messages.
	Step4Start        = "Step 4:"                  // Start marker for a specific step in a response.
	Step4End          = "{4-END}"                  // End marker for a specific step in a response.
)

// Variables used in the package.
var (
	openAIKey = os.Getenv("OPENAI_API_KEY") // Retrieve the OpenAI API key from the environment variable.
	client    = openai.NewClient(openAIKey) // Create an OpenAI client.
)

// integratedHandler handles the integration of various functionalities.
func configGenerator(userInput json.RawMessage, docsPath string, isRag bool) (string, error) {
	// Check if the OpenAI API key is set.
	if openAIKey == "" {
		log.Fatal("OPENAI_API_KEY not found.")
	}

	// Generate a prompt from the JSON.
	userQuery, err := generateUserQuery(userInput)
	if err != nil {
		return "", fmt.Errorf(" Error building initial user query: %w", err)
	}

	var initialResponse string
	// If isRag is true, fetch a configuration from VectorStore, else use OpenAI API for the config generation.
	if isRag {
		_, err = fetchConfigFromVecStore(userQuery, docsPath, isRag)
		if err != nil {
			return "", fmt.Errorf("error getting response from VectorStore: %w", err)
		}

		// Read the output from a file.
		initialResponse, err = readOutput(FilePath)
		if err != nil {
			return "", fmt.Errorf("error reading python process output: %w", err)
		}
	} else {
		// Build an initial Terraform configuration.
		initialResponse, err = initialConfig(userQuery)
		if err != nil {
			return "", fmt.Errorf("error building initial terraform configuration: %w", err)
		}
	}

	// Refine the embedding response.
	refinedResponse, err := refineEmbeddingResponse(json.RawMessage([]byte(initialResponse)), string(userInput))
	if err != nil {
		return "", fmt.Errorf("error refining the response: %w", err)
	}

	// Save the refined response to a file.
	err = saveToFile([]byte(refinedResponse), TerrFile)
	if err != nil {
		return "", fmt.Errorf("error writing to file: %w", err)
	}

	return refinedResponse, nil
}

// generateUserQuery generates a user query.
func generateUserQuery(usr_msg json.RawMessage) (string, error) {
	var assistant_mesages []string
	user_message := string(usr_msg)
	system_message := Prompts("system")

	// Set the tone of the response.
	response, err := setTone(client, system_message, user_message, false)
	if err != nil {
		return "", err
	}
	assistant_mesages = append(assistant_mesages, response)

	// Find and extract a specific section of the response.
	startIndex := strings.Index(response, Step4Start)
	endIndex := strings.Index(response, Step4End)

	if startIndex == -1 || endIndex == -1 {
		fmt.Println("Step 4 not found.")
		return "", errors.New("delimiters not found")
	}

	step4Query := strings.TrimSpace(response[startIndex+len("Step 4:") : endIndex])
	return strings.ReplaceAll(strings.TrimSpace(step4Query), "\\", ""), nil
}

// setTone sets the tone of the response using OpenAI.
func setTone(client *openai.Client, sys_msg, usr_msg string, isRefine bool) (string, error) {
	model := openai.GPT3Dot5Turbo16K

	// Create a chat completion request to generate a response.
	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model:       model,
			Temperature: 0.5,
			MaxTokens:   3000,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleSystem,
					Content: sys_msg,
				},
				{
					Role:    openai.ChatMessageRoleUser,
					Content: usr_msg,
				},
			},
		},
	)

	if err != nil {
		return "", err
	}
	return fmt.Sprintf(resp.Choices[0].Message.Content), nil
}

// refineEmbeddingResponse refines the embedding response.
func refineEmbeddingResponse(usr_msg json.RawMessage, userInput string) (string, error) {
	user_message := userInput + "\n" + Delimiter + "\n" + string(usr_msg) + Delimiter
	system_message := Prompts("refinement")

	// Set the tone for refining the response.
	response, err := setTone(client, system_message, user_message, true)
	if err != nil {
		return "", err
	}
	return response, nil
}

// initialConfig builds an initial Terraform configuration.
func initialConfig(query string) (string, error) {
	sys_msg := Prompts("initial")

	// Set the tone for the initial configuration.
	response, _ := setTone(client, sys_msg, query, false)
	return strings.Split(response, "```")[0], nil
}

// fetchConfigFromVecStore fetches a configuration from VectorStore.
func fetchConfigFromVecStore(query, path string, isRag bool) (string, error) {
	rag := fmt.Sprintf("%v", isRag)
	pythonCmd := exec.Command(PythonInterpreter, EmbeddingsScript, query, path, rag)
	pythonCmd.Env = append(os.Environ(), openAIKey)
	output, err := pythonCmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("running Python script: %v", string(output))
	}

	return string(output), nil
}
