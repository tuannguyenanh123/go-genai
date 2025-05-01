package main

import (
	"bufio"
	"context"
	"log"
	"os"

	"github.com/google/generative-ai-go/genai"
	"github.com/joho/godotenv"
	"google.golang.org/api/option"
)

const GenaiModel = "gemini-1.5-flash"
const KEY = "GENAI_API_KEY"
const KEY_GEMI = "GEMINI_API_KEY"

type App struct {
	client *genai.Client
	model  *genai.GenerativeModel
	cs     *genai.ChatSession
}

func main() {
	var err error
	err = godotenv.Load()

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// đọc dữ liệu từ người nhập
	genaiApp := &App{}
	reader := bufio.NewReader(os.Stdin)
	genaiApp.client, err = genai.NewClient(context.Background(), option.WithAPIKey(os.Getenv(KEY)))

	if err != nil {
		log.Fatal(err)
	}

	genaiApp.model = NewModel(genaiApp.client, GenaiModel)
	genaiApp.model.Tools = []*genai.Tool{FileTools}
	genaiApp.cs = genaiApp.model.StartChat()

	for {
		input, _ := reader.ReadString('\n')
		input = input[:len(input)-1]

		res, err := genaiApp.cs.SendMessage(context.Background(), genai.Text(input))
		if err != nil {
			log.Println("Error sending message", err)
			return
		}

		resString := buildResponse(res, genaiApp.cs)
		log.Println(resString)
	}
}

func buildResponse(resp *genai.GenerateContentResponse, cs *genai.ChatSession) string {
	funcResp := make(map[string]interface{})

	for _, v := range resp.Candidates[0].Content.Parts {
		functionCall, ok := v.(genai.FunctionCall)
		if ok {
			log.Println("Function call: ", functionCall.Name)
			switch functionCall.Name {
			case "file_write":
				//file_write tool has been called
				fileName, okFileName := functionCall.Args["filename"].(string)
				content, okContent := functionCall.Args["filename"].(string)

				if !okFileName || fileName == "" {
					funcResp["error"] = "Expected non-empty string at key filename."
					break
				}
				if !okContent || content == "" {
					funcResp["error"] = "Expected non-empty string at content."
					break
				}

				err := WriteDesktop(fileName, content)
				if err != nil {
					funcResp["error"] = "Could not write file."
				} else {
					funcResp["result"] = "File successfully written."
				}
			default:
				funcResp["error"] = "Unknown function call."
				//unrecognized function call
			}

		}
	}

	if len(funcResp) > 0 {
		resp, err := cs.SendMessage(context.Background(), genai.FunctionResponse{
			Name:     "Function_call",
			Response: funcResp,
		})

		if err != nil {
			log.Println("Error sending message", err.Error())
		}
		funcResp = nil
		return buildResponse(resp, cs)
	}

	for _, cand := range resp.Candidates {
		if cand.Content != nil {
			for _, part := range cand.Content.Parts {
				res, ok := part.(genai.Text)
				if ok {
					return string(res)
				}
			}
		}
	}

	return ""
}
