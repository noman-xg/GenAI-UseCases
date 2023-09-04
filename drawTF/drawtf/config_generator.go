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
	"time"

	openai "github.com/sashabaranov/go-openai"
)

const (
	PythonInterpreter = "python3"
	EmbeddingsScript  = "./scripts/embeddings.py"
	GradioLauncher    = "./scripts/chat-portal.py"
	FilePath          = "./file.txt"
	TerrFile          = "./main.tf"
	Delimiter         = "####"
	Step4Start        = "Step 4:"
	Step4End          = "{4-END}"
)

var (
	openAIKey = os.Getenv("OPENAI_API_KEY")
	client    = openai.NewClient(openAIKey)
)

func integratedHandler(userInput json.RawMessage) (string, error) {

	if openAIKey == "" {
		log.Fatal("OPENAI_API_KEY not found.")
	}

	start := time.Now()
	userQuery, err := generateUserQuery(userInput)
	if err != nil {
		return "", fmt.Errorf(" Error building initial user query: %w", err)
	}
	timeDelta(start, "generateUserQuery()")

	// start = time.Now()
	// _, err = fetchConfigFromVecStore(userQuery)
	// if err != nil {
	// 	return "", fmt.Errorf("error getting response from VectorStore: %w", err)
	// }
	// timeDelta(start, "fetchConfigFromVecStore()")

	start = time.Now()
	response, err := initialConfig(userQuery)
	if err != nil {
		return "", fmt.Errorf("error building initial terraform configuration: %w", err)
	}
	timeDelta(start, "experiment()")

	// start = time.Now()
	// response, err := readOutput(FilePath)
	// if err != nil {
	// 	return "", fmt.Errorf("error reading python process output: %w", err)
	// }
	// timeDelta(start, "readOutput()")

	start = time.Now()
	refinedResponse, err := refineEmbeddingResponse(json.RawMessage([]byte(response)), string(userInput))
	if err != nil {
		return "", fmt.Errorf("error refining the response: %w", err)
	}
	timeDelta(start, "refineEmbeddingResponse()")

	err = saveToFile([]byte(refinedResponse), TerrFile)
	if err != nil {
		return "", fmt.Errorf("error writing to file: %w", err)
	}

	return refinedResponse, nil
}

func generateUserQuery(usr_msg json.RawMessage) (string, error) {
	var assistant_mesages []string
	user_message := string(usr_msg)
	system_message := Prompts("system")

	response, err := setTone(client, system_message, user_message, false)
	if err != nil {
		return "", err
	}
	assistant_mesages = append(assistant_mesages, response)

	startIndex := strings.Index(response, Step4Start)
	endIndex := strings.Index(response, Step4End)

	if startIndex == -1 || endIndex == -1 {
		fmt.Println("Step 4 not found.")
		return "", errors.New("delimiters not found")
	}

	step4Query := strings.TrimSpace(response[startIndex+len("Step 4:") : endIndex])
	return strings.ReplaceAll(strings.TrimSpace(step4Query), "\\", ""), nil
}

func setTone(client *openai.Client, sys_msg, usr_msg string, isRefine bool) (string, error) {
	model := openai.GPT3Dot5Turbo16K

	// if isRefine {
	// 	model = openai.GPT4
	// }

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

func refineEmbeddingResponse(usr_msg json.RawMessage, userInput string) (string, error) {
	user_message := userInput + "\n" + Delimiter + "\n" + string(usr_msg) + Delimiter
	system_message := Prompts("refinement")
	response, err := setTone(client, system_message, user_message, true)
	if err != nil {
		return "", err
	}
	return response, nil
}

func initialConfig(query string) (string, error) {
	sys_msg := Prompts("experiment")
	response, _ := setTone(client, sys_msg, query, false)
	fmt.Println(response)
	return strings.Split(response, "```")[0], nil
}

func fetchConfigFromVecStore(query, path string) (string, error) {
	pythonCmd := exec.Command(PythonInterpreter, EmbeddingsScript, query, path)
	pythonCmd.Env = append(os.Environ(), openAIKey)
	output, err := pythonCmd.CombinedOutput()
	fmt.Println(query)
	if err != nil {
		return "", fmt.Errorf("running Python script: %v", string(output))
	}

	return string(output), nil
}
