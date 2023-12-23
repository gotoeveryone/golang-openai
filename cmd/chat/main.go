package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	chat "github.com/gotoeveryone/golang-openai"
	"github.com/joho/godotenv"
)

const endpoint = "https://api.openai.com/v1/chat/completions"

var messages []chat.Message

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalln("Error loading .env file")
	}

	apiKey := os.Getenv("OPENAI_API_KEY")

	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("Ask a question: ")
		question, _ := reader.ReadString('\n')
		question = strings.TrimSpace(question)

		if question == "exit" {
			break
		}

		messages = append(messages, chat.Message{
			Role:    "user",
			Content: question,
		})

		message, err := getMessage(apiKey)
		if err != nil {
			log.Fatalln(err)
		}

		fmt.Printf("%s\n", message)
	}
}

func getMessage(apiKey string) (string, error) {
	res, err := getOpenAIResponse(apiKey)
	if err != nil {
		return "", err
	}

	if len(res.Choices) == 0 {
		return "", errors.New("has not choices from OpenAI")
	}

	message := res.Choices[0].Messages.Content
	messages = append(messages, chat.Message{
		Role:    "assistant",
		Content: message,
	})

	return message, err
}

func getOpenAIResponse(apiKey string) (*chat.OpenAIResponse, error) {
	reqJSON, _ := json.Marshal(chat.OpenAIRequest{
		Model:    "gpt-3.5-turbo",
		Messages: messages,
	})

	req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(reqJSON))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", apiKey))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var res chat.OpenAIResponse
	err = json.Unmarshal(body, &res)
	if err != nil {
		return nil, err
	}

	return &res, nil
}
